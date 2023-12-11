use criterion::{black_box, criterion_group, criterion_main, Criterion};
use lazy_static::lazy_static;
use std::{cell::RefCell, env, rc::Rc, sync::Arc, time::Duration};
use tvix_castore::{
    blobservice::{BlobService, MemoryBlobService},
    directoryservice::{DirectoryService, MemoryDirectoryService},
};
use tvix_glue::{
    builtins::add_derivation_builtins, configure_nix_path, known_paths::KnownPaths,
    tvix_store_io::TvixStoreIO,
};
use tvix_store::pathinfoservice::{MemoryPathInfoService, PathInfoService};

lazy_static! {
    static ref BLOB_SERVICE: Arc<dyn BlobService> = Arc::new(MemoryBlobService::default());
    static ref DIRECTORY_SERVICE: Arc<dyn DirectoryService> =
        Arc::new(MemoryDirectoryService::default());
    static ref PATH_INFO_SERVICE: Arc<dyn PathInfoService> = Arc::new(MemoryPathInfoService::new(
        BLOB_SERVICE.clone(),
        DIRECTORY_SERVICE.clone(),
    ));
    static ref TOKIO_RUNTIME: tokio::runtime::Runtime = tokio::runtime::Runtime::new().unwrap();
}

fn interpret(code: &str) {
    // TODO: this is a bit annoying.
    // It'd be nice if we could set this up once and then run evaluate() with a
    // piece of code. b/262
    let mut eval = tvix_eval::Evaluation::new_impure(code, None);

    let known_paths: Rc<RefCell<KnownPaths>> = Default::default();
    add_derivation_builtins(&mut eval, known_paths.clone());
    configure_nix_path(
        &mut eval,
        // The benchmark requires TVIX_BENCH_NIX_PATH to be set, so barf out
        // early, rather than benchmarking tvix returning an error.
        &Some(env::var("TVIX_BENCH_NIX_PATH").expect("TVIX_BENCH_NIX_PATH must be set")),
    );

    eval.io_handle = Box::new(tvix_glue::tvix_io::TvixIO::new(
        known_paths.clone(),
        TvixStoreIO::new(
            BLOB_SERVICE.clone(),
            DIRECTORY_SERVICE.clone(),
            PATH_INFO_SERVICE.clone(),
            TOKIO_RUNTIME.handle().clone(),
        ),
    ));

    let result = eval.evaluate();

    assert!(result.errors.is_empty());
}

fn eval_nixpkgs(c: &mut Criterion) {
    c.bench_function("hello outpath", |b| {
        b.iter(|| {
            interpret(black_box("(import <nixpkgs> {}).hello.outPath"));
        })
    });
}

criterion_group!(
    name = benches;
    config = Criterion::default().measurement_time(Duration::from_secs(30)).sample_size(10);
    targets = eval_nixpkgs
);
criterion_main!(benches);