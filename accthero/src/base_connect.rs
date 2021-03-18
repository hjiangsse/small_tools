use postgres::{Client, NoTls};

#[derive(Debug)]
pub struct ConnectData {
    host: String,
    user: String,
    password: String,
    dbname: String,
}

impl ConnectData {
    pub fn new(host: &str, user: &str, password: &str, dbname: &str) -> ConnectData {
        ConnectData {
            host: String::from(host),
            user: String::from(user),
            password: String::from(password),
            dbname: String::from(dbname),
        }
    }

    pub fn init_connection(&self) -> Client {
        let client = Client::connect(
            &format!(
                "host={} user={} password={} dbname={}",
                self.host,
                self.user,
                self.password,
                self.dbname
            ),
            NoTls,
        ).unwrap();
        client
    }
}
