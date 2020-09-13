mod configs;

fn main() {
    let config = configs::load_configs("./config/config.toml");
    let offset_hash = configs::get_fields_offset_info(&config);
    for (key, val) in offset_hash.iter() {
        println!("key: {}, value: {:?}", key, val);
    }
    let filter_offsets = configs::get_filters_offset_info(&offset_hash, &config.filterinfos);
    println!("{:?}", filter_offsets);
}
