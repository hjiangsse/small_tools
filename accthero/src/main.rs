mod base_connect;
mod data_collect;

use crate::base_connect::ConnectData;
use crate::data_collect::AccountItem;

use structopt::StructOpt;

#[derive(Debug, StructOpt)]
struct Opt {
    #[structopt(short = "i", long = "insert", help = "insert item into account_book")]
    insert: bool,
    #[structopt(short = "a", long = "showall", help = "show all items in account_book")]
    showall: bool,
}

fn main() {
    let conn_data = ConnectData::new("localhost", "postgres", "jiang186212", "postlab");
    let mut client = conn_data.init_connection();

    let opt = Opt::from_args();
    if opt.insert {
        let item = AccountItem::new();
        let _ = client
            .query(
                "insert into account_book(item, cost) values ($1, $2)",
                &[&item.item, &item.cost],
            )
            .unwrap();
    }

    if opt.showall {
        for row in client
            .query("select id, item, cost, time from account_book", &[])
            .unwrap()
        {
            let cur_index: i32 = row.get(0);
            let cur_item: &str = row.get(1);
            let cur_cost: f32 = row.get(2);

            println!("{}:{:20} -- {}()", cur_index, cur_item, cur_cost);
        }
    }
}
