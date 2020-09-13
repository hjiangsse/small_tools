use serde_derive::Deserialize;
use std::collections::HashMap;
use std::fs;

#[derive(Deserialize, Debug)]
pub struct Config {
    pub metainfo: Metainfo,
    pub fieldinfos: Fieldinfos,
    pub filterinfos: Filterinfos,
}

#[derive(Deserialize, Debug)]
pub struct Metainfo {
    filetype: String,
    ifmindexname: String,
    ifmindexnum: u32,
    ficherindexname: String,
    ficherindexnum: u32,
}

#[derive(Deserialize, Debug)]
pub struct Fieldinfos {
    fields: Vec<(String, String, u32)>,
}

#[derive(Deserialize, Debug)]
pub struct Filterinfos {
    pub filters: Vec<(String, u32, u32)>,
}

//load config file, return a Big Config structure out
pub fn load_configs(path: &str) -> Config {
    let configs_str = fs::read_to_string(path).unwrap();
    let config: Config = toml::from_str(configs_str.as_str()).unwrap();
    config
}

//get fileds offset info as a hash table
pub fn get_fields_offset_info(config: &Config) -> HashMap<String, (u32, u32)> {
    let mut res = HashMap::new();
    let mut start_offset = 0;
    let mut end_offset;
    for info in config.fieldinfos.fields.iter() {
        end_offset = start_offset + info.2;
        res.insert(info.0.clone(), (start_offset, end_offset));
        start_offset = end_offset + 1;
    }
    res
}

//get filter offset info
pub fn get_filters_offset_info(
    fields_info_map: &HashMap<String, (u32, u32)>,
    filterinfos: &Filterinfos,
) -> Vec<(String, u32, u32)> {
    let mut res = Vec::new();

    for filter in filterinfos.filters.iter() {
        match fields_info_map.get(&filter.0) {
            Some(&offset_pair) => {
                if offset_pair.0 + filter.1 > offset_pair.1 {
                    res.push((filter.0.clone(), offset_pair.0, offset_pair.1));
                } else if offset_pair.0 + filter.1 + filter.2 > offset_pair.1 {
                    res.push((filter.0.clone(), offset_pair.0 + filter.1, offset_pair.1));
                } else {
                    res.push((
                        filter.0.clone(),
                        offset_pair.0 + filter.1,
                        offset_pair.0 + filter.1 + filter.2,
                    ));
                }
            }
            _ => println!("can not find filed {}", filter.0),
        }
    }
    res
}

mod tests {
    //use super::*;

    #[test]
    fn test_config_load() {
        //let config = load_configs("./config/config.toml");
    }
}
