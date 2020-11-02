use ssh2::Session;
use std::fs;
use std::fs::File;
use std::io::stdin;
use std::io::Read;
use std::io::Write;
use std::net::TcpStream;
use std::path::Path;

pub fn list_all_remote_addresses(addrs: &[String]) {
    println!("The remote addresses: ");
    for (index, addr) in addrs.iter().enumerate() {
        println!("{}-{}", index + 1, addr);
    }
}

pub fn promote_user_input_index(promote: &str) -> usize {
    println!("{}", promote);

    let mut user_choose = String::new();
    stdin()
        .read_line(&mut user_choose)
        .expect("User not enter a corret thing");
    let choose_index: usize = user_choose.trim().parse().unwrap();
    choose_index
}

pub fn search_remote_files_under_some_path(
    remote_addr: &str,
    remote_path: &str,
    file_name: &str,
) -> String {
    let mut exec_cmd = String::new();
    exec_cmd.push_str("find -L ");
    exec_cmd.push_str(remote_path);
    exec_cmd.push_str(" -name *");
    exec_cmd.push_str(file_name);
    exec_cmd.push_str("*");

    exec_ssh_remote_command(remote_addr, &exec_cmd)
        .trim_end()
        .to_string()
}

pub fn download_file_from_remote_mechine(
    addrinfo: &str,
    remote_file_path: &str,
    local_file_dir: &str,
) {
    let (user, hostaddr, passwd) = parse_addr_elements(addrinfo);

    let mut addr_port = String::new();
    addr_port.push_str(&hostaddr);
    addr_port.push_str(":22");

    let tcp = TcpStream::connect(&addr_port).unwrap();
    let mut sess = Session::new().unwrap();

    sess.handshake(&tcp).unwrap();
    sess.userauth_password(&user, &passwd).unwrap();

    let remote_path = Path::new(remote_file_path);
    let (mut remote_file, _) = sess.scp_recv(remote_path).unwrap();

    let mut contents = Vec::new();
    remote_file.read_to_end(&mut contents).unwrap();

    let local_file = Path::new(local_file_dir).join(remote_path.file_name().unwrap());

    let mut local_write = File::create(local_file).unwrap();
    local_write.write(&contents).unwrap();

    println!(
        "Down load file *{}* from *{}* finish!",
        remote_file_path, addrinfo
    );
}

pub fn upload_file_to_remote_mechine(addrinfo: &str, local_file_path: &str, remote_file_dir: &str) {
    let (user, hostaddr, passwd) = parse_addr_elements(addrinfo);

    let mut addr_port = String::new();
    addr_port.push_str(&hostaddr);
    addr_port.push_str(":22");

    let tcp = TcpStream::connect(&addr_port).unwrap();
    let mut sess = Session::new().unwrap();

    sess.handshake(&tcp).unwrap();
    sess.userauth_password(&user, &passwd).unwrap();

    let local_path = Path::new(local_file_path);
    let local_file_name = local_path.file_name().unwrap().to_str().unwrap();

    let data = fs::read(local_path).expect("unable to read file");

    let mut remote_file = sess
        .scp_send(
            &Path::new(remote_file_dir).join(local_file_name),
            0o644,
            data.len() as u64,
            None,
        )
        .unwrap();

    remote_file.write(&data).unwrap();

    println!(
        "Upload *{}* to remote *{}* finish!",
        local_file_path, addrinfo
    );
}

pub fn exec_ssh_remote_command(addrinfo: &str, cmd: &str) -> String {
    let (user, hostaddr, passwd) = parse_addr_elements(addrinfo);

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
