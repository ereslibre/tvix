use std::io;

use auto_impl::auto_impl;
use tonic::async_trait;

use crate::composition::{Registry, ServiceBuilder};
use crate::proto::stat_blob_response::ChunkMeta;
use crate::B3Digest;

mod chunked_reader;
mod combinator;
mod from_addr;
mod grpc;
mod memory;
mod object_store;

#[cfg(test)]
pub mod tests;

pub use self::chunked_reader::ChunkedReader;
pub use self::combinator::{CombinedBlobService, CombinedBlobServiceConfig};
pub use self::from_addr::from_addr;
pub use self::grpc::{GRPCBlobService, GRPCBlobServiceConfig};
pub use self::memory::{MemoryBlobService, MemoryBlobServiceConfig};
pub use self::object_store::{ObjectStoreBlobService, ObjectStoreBlobServiceConfig};

/// The base trait all BlobService services need to implement.
/// It provides functions to check whether a given blob exists,
/// a way to read (and seek) a blob, and a method to create a blobwriter handle,
/// which will implement a writer interface, and also provides a close funtion,
/// to finalize a blob and get its digest.
#[async_trait]
#[auto_impl(&, &mut, Arc, Box)]
pub trait BlobService: Send + Sync {
    /// Check if the service has the blob, by its content hash.
    /// On implementations returning chunks, this must also work for chunks.
    async fn has(&self, digest: &B3Digest) -> io::Result<bool>;

    /// Request a blob from the store, by its content hash.
    /// On implementations returning chunks, this must also work for chunks.
    async fn open_read(&self, digest: &B3Digest) -> io::Result<Option<Box<dyn BlobReader>>>;

    /// Insert a new blob into the store. Returns a [BlobWriter], which
    /// implements [tokio::io::AsyncWrite] and a [BlobWriter::close] to finalize
    /// the blob and get its digest.
    async fn open_write(&self) -> Box<dyn BlobWriter>;

    /// Return a list of chunks for a given blob.
    /// There's a distinction between returning Ok(None) and Ok(Some(vec![])).
    /// The former return value is sent in case the blob is not present at all,
    /// while the second one is sent in case there's no more granular chunks (or
    /// the backend does not support chunking).
    /// A default implementation checking for existence and then returning it
    /// does not have more granular chunks available is provided.
    async fn chunks(&self, digest: &B3Digest) -> io::Result<Option<Vec<ChunkMeta>>> {
        if !self.has(digest).await? {
            return Ok(None);
        }
        // default implementation, signalling the backend does not have more
        // granular chunks available.
        Ok(Some(vec![]))
    }
}

/// A [tokio::io::AsyncWrite] that the user needs to close() afterwards for persist.
/// On success, it returns the digest of the written blob.
#[async_trait]
pub trait BlobWriter: tokio::io::AsyncWrite + Send + Unpin {
    /// Signal there's no more data to be written, and return the digest of the
    /// contents written.
    ///
    /// Closing a already-closed BlobWriter is a no-op.
    async fn close(&mut self) -> io::Result<B3Digest>;
}

/// BlobReader is a [tokio::io::AsyncRead] that also allows seeking.
pub trait BlobReader: tokio::io::AsyncRead + tokio::io::AsyncSeek + Send + Unpin + 'static {}

/// A [`io::Cursor<Vec<u8>>`] can be used as a BlobReader.
impl BlobReader for io::Cursor<&'static [u8]> {}
impl BlobReader for io::Cursor<&'static [u8; 0]> {}
impl BlobReader for io::Cursor<Vec<u8>> {}
impl BlobReader for io::Cursor<bytes::Bytes> {}
impl BlobReader for tokio::fs::File {}

/// Registers the builtin BlobService implementations with the registry
pub(crate) fn register_blob_services(reg: &mut Registry) {
    reg.register::<Box<dyn ServiceBuilder<Output = dyn BlobService>>, super::blobservice::ObjectStoreBlobServiceConfig>("objectstore");
    reg.register::<Box<dyn ServiceBuilder<Output = dyn BlobService>>, super::blobservice::MemoryBlobServiceConfig>("memory");
    reg.register::<Box<dyn ServiceBuilder<Output = dyn BlobService>>, super::blobservice::CombinedBlobServiceConfig>("combined");
    reg.register::<Box<dyn ServiceBuilder<Output = dyn BlobService>>, super::blobservice::GRPCBlobServiceConfig>("grpc");
}
