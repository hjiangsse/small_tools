use std::io::stdin;

#[derive(Debug)]
pub struct AccountItem {
    pub item: String,
    pub cost: f32,
}

impl AccountItem {
    pub fn new() -> AccountItem {
        let mut m_item = String::new();
        println!("Please input your transaction item: ");
        stdin()
            .read_line(&mut m_item)
            .expect("Invalid input.(I want a string!)");

        let mut m_cost = String::new();
        loop {
            m_cost.clear();
            println!("Please input your cost for this item: ");
            stdin()
                .read_line(&mut m_cost)
                .expect("Invalid input.(I want a string!)");
            match m_cost.trim().parse::<f32>() {
                Ok(cnum) => {
                    break AccountItem {
                        item: m_item.trim().to_string(),
                        cost: cnum,
                    }
                }
                Err(_) => println!("Sorry! Invalid f32 number!"),
            }
        }
    }
}
