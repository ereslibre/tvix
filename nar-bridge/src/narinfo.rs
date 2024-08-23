use axum::{http::StatusCode, response::IntoResponse};
use bytes::Bytes;
use nix_compat::{narinfo::NarInfo, nix_http, nixbase32};
use prost::Message;
use tracing::{instrument, warn, Span};
use tvix_castore::proto::{self as castorepb};
use tvix_store::proto::PathInfo;

use crate::AppState;

/// The size limit for NARInfo uploads nar-bridge receives
const NARINFO_LIMIT: usize = 2 * 1024 * 1024;

#[instrument(skip(path_info_service))]
pub async fn head(
    axum::extract::Path(narinfo_str): axum::extract::Path<String>,
    axum::extract::State(AppState {
        path_info_service, ..
    }): axum::extract::State<AppState>,
) -> Result<impl IntoResponse, StatusCode> {
    let digest = nix_http::parse_narinfo_str(&narinfo_str).ok_or(StatusCode::NOT_FOUND)?;
    Span::current().record("path_info.digest", &narinfo_str[0..32]);

    if path_info_service
        .get(digest)
        .await
        .map_err(|e| {
            warn!(err=%e, "failed to get PathInfo");
            StatusCode::INTERNAL_SERVER_ERROR
        })?
        .is_some()
    {
        Ok(([("content-type", nix_http::MIME_TYPE_NARINFO)], ""))
    } else {
        warn!("PathInfo not found");
        Err(StatusCode::NOT_FOUND)
    }
}

#[instrument(skip(path_info_service))]
pub async fn get(
    axum::extract::Path(narinfo_str): axum::extract::Path<String>,
    axum::extract::State(AppState {
        path_info_service, ..
    }): axum::extract::State<AppState>,
) -> Result<impl IntoResponse, StatusCode> {
    let digest = nix_http::parse_narinfo_str(&narinfo_str).ok_or(StatusCode::NOT_FOUND)?;
    Span::current().record("path_info.digest", &narinfo_str[0..32]);

    // fetch the PathInfo
    let path_info = path_info_service
        .get(digest)
        .await
        .map_err(|e| {
            warn!(err=%e, "failed to get PathInfo");
            StatusCode::INTERNAL_SERVER_ERROR
        })?
        .ok_or(StatusCode::NOT_FOUND)?;

    let store_path = path_info.validate().map_err(|e| {
        warn!(err=%e, "invalid PathInfo");
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    let mut narinfo = path_info.to_narinfo(store_path.as_ref()).ok_or_else(|| {
        warn!(path_info=?path_info, "PathInfo contained no NAR data");
        StatusCode::INTERNAL_SERVER_ERROR
    })?;

    // encode the (unnamed) root node in the NAR url itself.
    // We strip the name from the proto node before sending it out.
    // It's not needed to render the NAR, it'll make the URL shorter, and it
    // will make caching these requests easier.
    let (_, root_node) = path_info
        .node
        .as_ref()
        .expect("invalid pathinfo")
        .to_owned()
        .into_name_and_node()
        .expect("invalid pathinfo");

    let url = format!(
        "nar/tvix-castore/{}?narsize={}",
        data_encoding::BASE64URL_NOPAD
            .encode(&castorepb::Node::from_name_and_node("".into(), root_node).encode_to_vec()),
        narinfo.nar_size,
    );

    narinfo.url = &url;

    Ok((
        [("content-type", nix_http::MIME_TYPE_NARINFO)],
        narinfo.to_string(),
    ))
}

#[instrument(skip(path_info_service, root_nodes, request))]
pub async fn put(
    axum::extract::Path(narinfo_str): axum::extract::Path<String>,
    axum::extract::State(AppState {
        path_info_service,
        root_nodes,
        ..
    }): axum::extract::State<AppState>,
    request: axum::extract::Request,
) -> Result<&'static str, StatusCode> {
    let _narinfo_digest = nix_http::parse_narinfo_str(&narinfo_str).ok_or(StatusCode::UNAUTHORIZED);
    Span::current().record("path_info.digest", &narinfo_str[0..32]);

    let narinfo_bytes: Bytes = axum::body::to_bytes(request.into_body(), NARINFO_LIMIT)
        .await
        .map_err(|e| {
            warn!(err=%e, "unable to fetch body");
            StatusCode::BAD_REQUEST
        })?;

    // Parse the narinfo from the body.
    let narinfo_str = std::str::from_utf8(narinfo_bytes.as_ref()).map_err(|e| {
        warn!(err=%e, "unable decode body as string");
        StatusCode::BAD_REQUEST
    })?;

    let narinfo = NarInfo::parse(narinfo_str).map_err(|e| {
        warn!(err=%e, "unable to parse narinfo");
        StatusCode::BAD_REQUEST
    })?;

    // Extract the NARHash from the PathInfo.
    Span::current().record("path_info.nar_info", nixbase32::encode(&narinfo.nar_hash));

    // populate the pathinfo.
    let mut pathinfo = PathInfo::from(&narinfo);

    // Lookup root node with peek, as we don't want to update the LRU list.
    // We need to be careful to not hold the RwLock across the await point.
    let maybe_root_node: Option<tvix_castore::Node> =
        root_nodes.read().peek(&narinfo.nar_hash).cloned();

    match maybe_root_node {
        Some(root_node) => {
            // Set the root node from the lookup.
            // We need to rename the node to the narinfo storepath basename, as
            // that's where it's stored in PathInfo.
            pathinfo.node = Some(castorepb::Node::from_name_and_node(
                narinfo.store_path.to_string().into(),
                root_node,
            ));

            // Persist the PathInfo.
            path_info_service.put(pathinfo).await.map_err(|e| {
                warn!(err=%e, "failed to persist the PathInfo");
                StatusCode::INTERNAL_SERVER_ERROR
            })?;

            Ok("")
        }
        None => {
            warn!("received narinfo with unknown NARHash");
            Err(StatusCode::BAD_REQUEST)
        }
    }
}
