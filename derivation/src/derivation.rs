use crate::nix_hash;
use crate::output::Output;
use crate::write;
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use std::{collections::BTreeMap, fmt, fmt::Write, iter::FromIterator};
use tvix_store::nixbase32::NIXBASE32;
use tvix_store::nixpath::{NixPath, ParseNixPathError, STORE_DIR};

#[derive(Clone, Debug, Eq, PartialEq, Serialize, Deserialize)]
pub struct Derivation {
    #[serde(rename = "args")]
    pub arguments: Vec<String>,

    pub builder: String,

    #[serde(rename = "env")]
    pub environment: BTreeMap<String, String>,

    #[serde(rename = "inputDrvs")]
    pub input_derivations: BTreeMap<String, Vec<String>>,

    #[serde(rename = "inputSrcs")]
    pub input_sources: Vec<String>,

    pub outputs: BTreeMap<String, Output>,

    pub system: String,
}

/// This returns a store path, either of a derivation or a regular output.
/// The path_hash is compressed to 20 bytes, and nixbase32-encoded (32 characters)
fn build_store_path(
    is_derivation: bool,
    path_hash: &[u8],
    name: &str,
) -> Result<NixPath, ParseNixPathError> {
    let compressed = nix_hash::compress_hash(path_hash, 20);
    if is_derivation {
        NixPath::from_string(
            format!(
                "{}-{}{}",
                NIXBASE32.encode(&compressed),
                name,
                write::DOT_FILE_EXT,
            )
            .as_str(),
        )
    } else {
        NixPath::from_string(format!("{}-{}", NIXBASE32.encode(&compressed), name,).as_str())
    }
}

impl Derivation {
    pub fn serialize(&self, writer: &mut impl Write) -> Result<(), fmt::Error> {
        writer.write_str(write::DERIVATION_PREFIX)?;
        writer.write_char(write::PAREN_OPEN)?;

        write::write_outputs(writer, &self.outputs)?;
        write::write_input_derivations(writer, &self.input_derivations)?;
        write::write_input_sources(writer, &self.input_sources)?;
        write::write_system(writer, &self.system)?;
        write::write_builder(writer, &self.builder)?;
        write::write_arguments(writer, &self.arguments)?;
        write::write_enviroment(writer, &self.environment)?;

        writer.write_char(write::PAREN_CLOSE)?;

        Ok(())
    }

    /// Returns the fixed output path and its hash
    // (if the Derivation is fixed output),
    /// or None if there is no fixed output.
    /// This takes some shortcuts in case more than one output exists, as this
    /// can't be a valid fixed-output Derivation.
    pub fn get_fixed_output(&self) -> Option<(&String, &Hash)> {
        if self.outputs.len() != 1 {
            return None;
        }

        match self.outputs.get("out") {
            #[allow(clippy::manual_map)]
            Some(out_output) => match &out_output.hash {
                Some(out_output_hash) => Some((&out_output.path, out_output_hash)),
                // There has to be a hash, otherwise it would not be FOD
                None => None,
            },
            None => None,
        }
    }

    /// Returns the drv path of a Derivation struct.
    ///
    /// The drv path is calculated like this:
    ///   - Write the fingerprint of the Derivation to the sha256 hash function.
    ///     This is: `text:`,
    ///     all d.InputDerivations and d.InputSources (sorted, separated by a `:`),
    ///     a `:`,
    ///     a `sha256:`, followed by the sha256 digest of the ATerm representation (hex-encoded)
    ///     a `:`,
    ///     the storeDir, followed by a `:`,
    ///     the name of a derivation,
    ///     a `.drv`.
    ///   - Write the .drv A-Term contents to a hash function
    ///   - Take the digest, run hash.CompressHash(digest, 20) on it.
    ///   - Encode it with nixbase32
    ///   - Use it (and the name) to construct a NixPath.
    pub fn calculate_derivation_path(&self, name: &str) -> Result<NixPath, ParseNixPathError> {
        let mut hasher = Sha256::new();

        // collect the list of paths from input_sources and input_derivations
        // into a sorted list, and join them by :
        hasher.update(write::TEXT_COLON);

        let concat_inputs: Vec<String> = {
            let mut inputs = self.input_sources.clone();
            let input_derivation_keys: Vec<String> =
                self.input_derivations.keys().cloned().collect();
            inputs.extend(input_derivation_keys);
            inputs.sort();
            inputs
        };

        for input in concat_inputs {
            hasher.update(input);
            hasher.update(write::COLON);
        }

        // calculate the sha256 hash of the ATerm representation, and represent
        // it as a hex-encoded string (prefixed with sha256:).
        hasher.update(write::SHA256_COLON);

        let digest = {
            let mut derivation_hasher = Sha256::new();
            derivation_hasher.update(self.to_string());
            derivation_hasher.finalize()
        };

        hasher.update(format!("{:x}", digest));
        hasher.update(write::COLON);
        hasher.update(STORE_DIR);
        hasher.update(write::COLON);

        hasher.update(name);
        hasher.update(write::DOT_FILE_EXT);

        build_store_path(true, &hasher.finalize(), name)
    }
}

impl fmt::Display for Derivation {
    /// Formats the Derivation in ATerm representation.
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        self.serialize(f)
    }
}
