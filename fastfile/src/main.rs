mod config_load;
mod user_interface;

use std::process;

fn main() {
    let config_structure = config_load::load_config_from_path("/etc/fastfile.toml");

    //list all available address
    let config_addrs = config_structure.get_addrs();
    user_interface::list_all_remote_addresses(&config_addrs);

    //let user choose the address
    let choose_index = user_interface::promote_user_input_address_index();

    //interact with remote mechine
    if choose_index <= config_addrs.len() {
        let exec_res =
            user_interface::exec_ssh_remote_command(&config_addrs[choose_index - 1], "ls -l");
        println!("{}", exec_res);
    } else {
        println!("You choose the wrong address, bye!");
        process::exit(0);
    }
}
