use serde_derive::Deserialize;
use std::fs;

#[derive(Deserialize, Debug)]
pub struct Config {
    localpath: Vec<String>,
    remotepath: Vec<String>,
    remoteupload: Vec<String>,
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

    pub fn get_local_domains(&self) -> Vec<String> {
        self.localpath.clone()
    }

    pub fn get_remote_domains(&self) -> Vec<String> {
        self.remotepath.clone()
    }

    pub fn get_remote_upload_paths(&self) -> Vec<String> {
        self.remoteupload.clone()
    }
}
