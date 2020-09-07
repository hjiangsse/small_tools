use serde_derive::Deserialize;
use std::fs;

#[derive(Deserialize, Debug)]
pub struct Config {
    remote: Remote,
    command: Command,
}

#[derive(Deserialize, Debug)]
struct Remote {
    ip: String,
    port: u32,
    username: String,
    passwd: String,
}

#[derive(Deserialize, Debug)]
struct Command {
    cphashs: Vec<(String, String)>,
}

pub fn get_configs(cfgpath: &str) -> Config {
    let config_str = fs::read_to_string(cfgpath).unwrap();
    let config: Config = toml::from_str(config_str.as_str()).unwrap();
    config
}

impl Config {
    pub fn get_remote_ip(&self) -> &str {
        &self.remote.ip
    }

    pub fn get_remote_port(&self) -> u32 {
        self.remote.port
    }

    pub fn get_remote_username(&self) -> &str {
        &self.remote.username
    }

    pub fn get_remote_passwd(&self) -> &str {
        &self.remote.passwd
    }

    pub fn get_cp_commands(&self) -> Vec<(String, String)> {
        self.command.cphashs.clone()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_config_load() {
        let config = get_configs("./conf/fastftp.toml");
        assert!(config.get_remote_ip() == "127.0.0.1");
        assert_eq!(config.get_remote_port(), 21);
        assert!(config.get_remote_username() == "hjiang");
        assert!(config.get_remote_passwd() == "jiang186212");
    }
}
