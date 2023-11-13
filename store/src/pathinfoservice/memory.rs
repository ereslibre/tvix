use super::PathInfoService;
use crate::{nar::calculate_size_and_sha256, proto::PathInfo};
use futures::{stream::iter, Stream};
use std::{
    collections::HashMap,
    pin::Pin,
    sync::{Arc, RwLock},
};
use tonic::async_trait;
use tvix_castore::proto as castorepb;
use tvix_castore::Error;
use tvix_castore::{blobservice::BlobService, directoryservice::DirectoryService};

pub struct MemoryPathInfoService {
    db: Arc<RwLock<HashMap<[u8; 20], PathInfo>>>,

    blob_service: Arc<dyn BlobService>,
    directory_service: Arc<dyn DirectoryService>,
}

impl MemoryPathInfoService {
    pub fn new(
        blob_service: Arc<dyn BlobService>,
        directory_service: Arc<dyn DirectoryService>,
    ) -> Self {
        Self {
            db: Default::default(),
            blob_service,
            directory_service,
        }
    }

    /// Constructs a [MemoryPathInfoService] from the passed [url::Url]:
    /// - scheme has to be `memory://`
    /// - there may not be a host.
    /// - there may not be a path.
    pub fn from_url(
        url: &url::Url,
        blob_service: Arc<dyn BlobService>,
        directory_service: Arc<dyn DirectoryService>,
    ) -> Result<Self, Error> {
        if url.scheme() != "memory" {
            return Err(Error::StorageError("invalid scheme".to_string()));
        }

        if url.has_host() || !url.path().is_empty() {
            return Err(Error::StorageError("invalid url".to_string()));
        }

        Ok(Self::new(blob_service, directory_service))
    }
}

#[async_trait]
impl PathInfoService for MemoryPathInfoService {
    async fn get(&self, digest: [u8; 20]) -> Result<Option<PathInfo>, Error> {
        let db = self.db.read().unwrap();

        match db.get(&digest) {
            None => Ok(None),
            Some(path_info) => Ok(Some(path_info.clone())),
        }
    }

    async fn put(&self, path_info: PathInfo) -> Result<PathInfo, Error> {
        // Call validate on the received PathInfo message.
        match path_info.validate() {
            Err(e) => Err(Error::InvalidRequest(format!(
                "failed to validate PathInfo: {}",
                e
            ))),

            // In case the PathInfo is valid, and we were able to extract a NixPath, store it in the database.
            // This overwrites existing PathInfo objects.
            Ok(nix_path) => {
                let mut db = self.db.write().unwrap();
                db.insert(*nix_path.digest(), path_info.clone());

                Ok(path_info)
            }
        }
    }

    async fn calculate_nar(
        &self,
        root_node: &castorepb::node::Node,
    ) -> Result<(u64, [u8; 32]), Error> {
        calculate_size_and_sha256(
            root_node,
            self.blob_service.clone(),
            self.directory_service.clone(),
        )
        .await
        .map_err(|e| Error::StorageError(e.to_string()))
    }

    fn list(&self) -> Pin<Box<dyn Stream<Item = Result<PathInfo, Error>> + Send>> {
        let db = self.db.read().unwrap();

        // Copy all elements into a list.
        // This is a bit ugly, because we can't have db escape the lifetime
        // of this function, but elements need to be returned owned anyways, and this in-
        // memory impl is only for testing purposes anyways.
        let items: Vec<_> = db.iter().map(|(_k, v)| Ok(v.clone())).collect();

        Box::pin(iter(items))
    }
}

#[cfg(test)]
mod tests {
    use crate::tests::utils::gen_blob_service;
    use crate::tests::utils::gen_directory_service;

    use super::MemoryPathInfoService;

    /// This uses a wrong scheme.
    #[test]
    fn test_invalid_scheme() {
        let url = url::Url::parse("http://foo.example/test").expect("must parse");

        assert!(
            MemoryPathInfoService::from_url(&url, gen_blob_service(), gen_directory_service())
                .is_err()
        );
    }

    /// This correctly sets the scheme, and doesn't set a path.
    #[test]
    fn test_valid_scheme() {
        let url = url::Url::parse("memory://").expect("must parse");

        assert!(
            MemoryPathInfoService::from_url(&url, gen_blob_service(), gen_directory_service())
                .is_ok()
        );
    }

    /// This sets the host to `foo`
    #[test]
    fn test_invalid_host() {
        let url = url::Url::parse("memory://foo").expect("must parse");

        assert!(
            MemoryPathInfoService::from_url(&url, gen_blob_service(), gen_directory_service())
                .is_err()
        );
    }

    /// This has the path "/", which is invalid.
    #[test]
    fn test_invalid_has_path() {
        let url = url::Url::parse("memory:///").expect("must parse");

        assert!(
            MemoryPathInfoService::from_url(&url, gen_blob_service(), gen_directory_service())
                .is_err()
        );
    }

    /// This has the path "/foo", which is invalid.
    #[test]
    fn test_invalid_path2() {
        let url = url::Url::parse("memory:///foo").expect("must parse");

        assert!(
            MemoryPathInfoService::from_url(&url, gen_blob_service(), gen_directory_service())
                .is_err()
        );
    }
}
