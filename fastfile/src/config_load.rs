use serde_derive::Deserialize;
use std::fs;

#[derive(Deserialize, Debug)]
pub struct Config {
    localpath: String,
    remotepath: String,
    addrs: Vec<String>,
}

pub fn load_config_from_path(path: &str) -> Config {
    let contents = fs::read_to_string(path).unwrap();
    let config_structure: Config = toml::from_str(&contents).unwrap();
    config_structure
}

impl Config {
    fn new() -> Config {
        Config {
            localpath: String::from(""),
            remotepath: String::from(""),
            addrs: vec![String::from("")],
        }
    }

    pub fn get_addrs(&self) -> Vec<String> {
        self.addrs.clone()
    }

    pub fn get_local_domain(&self) -> String {
        self.localpath.clone()
    }

    pub fn get_remote_domain(&self) -> String {
        self.remotepath.clone()
    }
}
