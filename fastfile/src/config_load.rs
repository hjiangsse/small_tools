use serde_derive::Deserialize;
use std::fs;

#[derive(Deserialize, Debug)]
pub struct Config {
    localpath: String,
    remotepath: String,
    remoteupload: String,
    addrs: Vec<String>,
}

pub fn load_config_from_path(path: &str) -> Config {
    let contents = fs::read_to_string(path).unwrap();
    let config_structure: Config = toml::from_str(&contents).unwrap();
    config_structure
}

impl Config {
    pub fn get_addrs(&self) -> Vec<String> {
        self.addrs.clone()
    }

    pub fn get_local_domain(&self) -> String {
        self.localpath.clone()
    }

    pub fn get_remote_domain(&self) -> String {
        self.remotepath.clone()
    }

    pub fn get_remote_upload_path(&self) -> String {
        self.remoteupload.clone()
    }
}
