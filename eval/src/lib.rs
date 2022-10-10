mod builtins;
mod chunk;
mod compiler;
mod errors;
mod eval;
pub mod observer;
mod opcode;
mod source;
mod spans;
mod upvalues;
mod value;
mod vm;
mod warnings;

mod nix_search_path;
#[cfg(test)]
mod properties;
#[cfg(test)]
mod test_utils;
#[cfg(test)]
mod tests;

// Re-export the public interface used by other crates.
pub use crate::builtins::global_builtins;
pub use crate::compiler::compile;
pub use crate::errors::EvalResult;
pub use crate::eval::{interpret, Options};
pub use crate::source::SourceCode;
pub use crate::value::Value;
pub use crate::vm::run_lambda;
