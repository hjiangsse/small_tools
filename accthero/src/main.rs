mod data_collect;
mod base_connect;

use crate::data_collect::AccountItem;
use crate::base_connect::ConnectData;

fn main() {
    let item = AccountItem::new();
    println!("{:?}", item);
    println!("{}", item.item);
    println!("{}", item.cost);

    let conn_data = ConnectData::new("localhost", "postgres", "jiang186212", "postlab");
    println!("{:?}", conn_data);

    let mut client = conn_data.init_connection();
    for row in client.query("select item, cost from account_table", &[]).unwrap() {
        let cur_item: &str = row.get(0);
        let cur_cost: f32 = row.get(1);

        println!("{} {}", cur_item, cur_cost);
    }
}
