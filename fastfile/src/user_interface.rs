use ssh2::Session;
use std::io::stdin;
use std::io::Read;
use std::net::TcpStream;

pub fn list_all_remote_addresses(addrs: &[String]) {
    println!("The remote addresses are here: ");
    for (index, addr) in addrs.iter().enumerate() {
        println!("{}-{}", index + 1, addr);
    }
    println!("please choose the one you want to interact with: ");
}

pub fn promote_user_input_address_index() -> usize {
    let mut user_choose = String::new();
    stdin()
        .read_line(&mut user_choose)
        .expect("User not enter a corret thing");
    let choose_index: usize = user_choose.trim().parse().unwrap();
    choose_index
}

pub fn exec_ssh_remote_command(addrinfo: &str, cmd: &str) -> String {
    let (user, hostaddr, passwd) = parse_addr_elements(addrinfo);

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
    channel.exec(cmd).unwrap();

    let mut s = String::new();
    channel.read_to_string(&mut s).unwrap();
    s
}

fn parse_addr_elements(addr: &str) -> (String, String, String) {
    let v: Vec<&str> = addr.split(':').collect();
    (String::from(v[0]), String::from(v[1]), String::from(v[2]))
}
