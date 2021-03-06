mod config_load;
mod user_interface;

use std::process;
use structopt::StructOpt;
extern crate dirs;
use std::path::{Path, PathBuf};
use walkdir::WalkDir;

extern crate subprocess;
use subprocess::Exec;

#[derive(Debug, StructOpt)]
struct Opt {
    #[structopt(
        short = "c",
        long = "config",
        default_value = "~/.fastfile.toml",
        help = "the config file path"
    )]
    config: String,
    #[structopt(
        short = "n",
        long = "filename",
        default_value = "cpp",
        help = "the default file name you want to send of download"
    )]
    filename: String,
    #[structopt(short = "s", long = "send", help = "send local file to remote mechine")]
    send: bool,
    #[structopt(
        short = "d",
        long = "download",
        help = "download file from remote mechine"
    )]
    download: bool,
    #[structopt(
        short = "r",
        long = "reconfig",
        help = "change the default config file"
    )]
    reconf: bool,
}

fn main() {
    let opt = Opt::from_args();

    let config_pathbuf = expand_tilde(&opt.config);
    let config_path = config_pathbuf.unwrap();

    if opt.reconf {
        let exit_status = Exec::cmd("vim").arg(config_path).join().unwrap();
        assert!(exit_status.success());
        process::exit(0);
    }

    let config_structure = config_load::load_config_from_path(config_path.to_str().unwrap());

    //list all available address
    let config_addrs = config_structure.get_addrs();
    user_interface::list_all_remote_addresses(&config_addrs);

    //let user choose the address
    let choose_index = user_interface::promote_user_input_index(
        "please choose the one you want to interact with: ",
    );

    let local_search_domains = config_structure.get_local_domains();
    let remote_search_domains = config_structure.get_remote_domains();
    let remote_upload_domains = config_structure.get_remote_upload_paths();

    if choose_index <= config_addrs.len() {
        //send file to remote mechine
        if opt.send && !opt.download {
            let local_search_domain =
                user_interface::promote_user_select_one_domain(&local_search_domains, "local");

            let file_list = get_local_files_under_some_path(&local_search_domain, &opt.filename);
            //if not find file under local path, exit
            if file_list.len() == 0 {
                println!("Sorry, can't find the file under local path. Bye!");
                process::exit(0);
            }

            println!("-----------------------------------------------------");
            for (pos, e) in file_list.iter().enumerate() {
                println!("{}: {}", pos + 1, e);
            }

            //promote user specify the file send to remote
            let choose_file_index = user_interface::promote_user_input_index(
                "Input the file you want send to remote: ",
            );

            if choose_file_index <= file_list.len() && choose_file_index > 0 {
                let remote_upload_domain = user_interface::promote_user_select_one_domain(
                    &remote_upload_domains,
                    "upload",
                );
                user_interface::upload_file_to_remote_mechine(
                    &config_addrs[choose_index - 1],
                    &file_list[choose_file_index - 1],
                    &remote_upload_domain,
                )
            }
        }

        //download file from remote mechine
        if opt.download && !opt.send {
            //search file in remote mechine
            let remote_search_domain =
                user_interface::promote_user_select_one_domain(&remote_search_domains, "remote");
            let remote_file_list_str = user_interface::search_remote_files_under_some_path(
                &config_addrs[choose_index - 1],
                &remote_search_domain,
                &opt.filename,
            );

            let rmt_files: Vec<&str> = remote_file_list_str.split('\n').collect();
            //if not find file under remote path, exit
            if rmt_files.len() == 0 {
                println!("Sorry, can't find file under remote path. Bye!");
                process::exit(0);
            }
            println!("-----------------------------------------------------");
            for (pos, e) in rmt_files.iter().enumerate() {
                if !(e).eq(&"") {
                    println!("{}: {}", pos + 1, e);
                }
            }

            //promote user specify the file want to download from remote mechine
            let choose_file_index = user_interface::promote_user_input_index(
                "Input the file you want download from remote: ",
            );

            if choose_file_index <= rmt_files.len() && choose_file_index > 0 {
                let local_search_domain =
                    user_interface::promote_user_select_one_domain(&local_search_domains, "local");
                user_interface::download_file_from_remote_mechine(
                    &config_addrs[choose_index - 1],
                    &rmt_files[choose_file_index - 1],
                    &local_search_domain,
                )
            }
        }
    } else {
        println!("You choose the wrong address, bye!");
        process::exit(0);
    }
}

fn expand_tilde<P: AsRef<Path>>(path_user_input: P) -> Option<PathBuf> {
    let p = path_user_input.as_ref();
    if !p.starts_with("~") {
        return Some(p.to_path_buf());
    }
    if p == Path::new("~") {
        return dirs::home_dir();
    }
    dirs::home_dir().map(|mut h| {
        if h == Path::new("/") {
            // Corner case: `h` root directory;
            // don't prepend extra `/`, just drop the tilde.
            p.strip_prefix("~").unwrap().to_path_buf()
        } else {
            h.push(p.strip_prefix("~/").unwrap());
            h
        }
    })
}

fn get_local_files_under_some_path(path: &str, file_name: &str) -> Vec<String> {
    let mut files = Vec::new();

    for entry in WalkDir::new(path)
        .follow_links(true)
        .into_iter()
        .filter_map(|e| e.ok())
    {
        let f_name = entry.file_name().to_string_lossy();

        if f_name.contains(file_name) {
            files.push(String::from(entry.path().to_str().unwrap()));
        }
    }

    files
}
