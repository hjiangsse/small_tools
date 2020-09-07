mod configs;

use ftp::FtpStream;
use std::fs::File;
use std::io::BufReader;
use std::path::Path;

fn main() -> std::io::Result<()> {
    let config = configs::get_configs("./conf/fastftp.toml");

    //create ftp connection stream
    let remote_ip = config.get_remote_ip();
    let remote_port = config.get_remote_port();
    let remote_addrs = String::from(remote_ip) + ":" + remote_port.to_string().as_str();
    let mut ftp_stream = FtpStream::connect(remote_addrs).unwrap();

    //login to the remote mechine
    let user_name = config.get_remote_username();
    let pass_word = config.get_remote_passwd();
    let _ = ftp_stream.login(user_name, pass_word).unwrap();

    //get cp commands
    let cp_commands_vec = config.get_cp_commands();
    for cmd in cp_commands_vec.iter() {
        if Path::new(&cmd.0).is_file() {
            let _ = ftp_stream.cwd(&cmd.1).unwrap();

            let file = File::open(&cmd.0).unwrap();
            let mut buf_reader = BufReader::new(file);

            let file_name = Path::new(&cmd.0).file_name().unwrap();

            let _ = ftp_stream.put(file_name.to_str().unwrap(), &mut buf_reader);
            println!("remote copy {} to {} success!", cmd.0, cmd.1);
        }
    }

    Ok(())
}
