use serde_derive::Deserialize;
use ssh2::Session;
use std::fs;
use std::io::prelude::*;
use std::io::stdin;
use std::net::TcpStream;
use std::process;

#[derive(Deserialize, Debug)]
struct Config {
    addrs: Vec<String>,
}

fn main() {
    let contents = fs::read_to_string("/etc/fastfile.toml").unwrap();
    let config_structure: Config = toml::from_str(&contents).unwrap();

    //list all available address
    println!("The remote addresses are here: ");
    for (index, addr) in config_structure.addrs.iter().enumerate() {
        println!("{}-{}", index + 1, addr);
    }
    println!("please choose the one you want to interact with: ");

    //let user choose the address
    let mut user_choose = String::new();
    stdin()
        .read_line(&mut user_choose)
        .expect("User not enter a corret thing");
    let choose_index: usize = user_choose.trim().parse().unwrap();

    //interact with remote mechine
    if choose_index <= config_structure.addrs.len() {
        let (user, hostaddr, passwd) =
            parse_addr_elements(&config_structure.addrs[choose_index - 1]);

        println!("-----------Remote Mechine Info---------------");
        println!("User: {}", user);
        println!("HostAddr: {}", hostaddr);
        println!("PassWord: {}", passwd);
        println!("---------------------------------------------");

        let mut addr_port = String::new();
        addr_port.push_str(&hostaddr);
        addr_port.push_str(":22");

        let tcp = TcpStream::connect(&addr_port).unwrap();
        let mut sess = Session::new().unwrap();

        sess.handshake(&tcp).unwrap();
        sess.userauth_password(&user, &passwd).unwrap();

        let mut channel = sess.channel_session().unwrap();
        channel.exec("find . -name *abba*").unwrap();
        let mut s = String::new();
        channel.read_to_string(&mut s).unwrap();
        println!("{}", s);
        println!("{}", channel.exit_status().unwrap());
    } else {
        println!("You choose the wrong address, bye!");
        process::exit(0);
    }
}

fn parse_addr_elements(addr: &str) -> (String, String, String) {
    let v: Vec<&str> = addr.split(':').collect();
    (String::from(v[0]), String::from(v[1]), String::from(v[2]))
}
