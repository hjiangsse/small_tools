use std::fs::File;
use std::io::{self, BufRead, Write};
use std::path::Path;

#[derive(Debug)]
pub enum Priority {
    Low(String),
    Middle(String),
    High(String),
}

#[derive(Debug)]
pub enum Progress {
    NotStart(String),
    InProcess(String),
    Finished(String),
}

#[derive(Debug)]
pub struct Item {
    pub index: u32,
    workitem: String,
    priority: Priority,
    progress: Progress,
}

impl Item {
    pub fn new() -> Item {
        Item {
            index: 1,
            workitem: String::from(""),
            priority: Priority::Low(String::from("low")),
            progress: Progress::NotStart(String::from("notstart")),
        }
    }

    pub fn set_index(&mut self, idx: u32) {
        self.index = idx
    }

    pub fn get_index(&self) -> u32 {
        self.index
    }

    pub fn set_workitem(&mut self, item: &str) {
        self.workitem = String::from(item);
    }

    pub fn get_workitem(&self) -> String {
        String::from(self.workitem.as_str())
    }

    pub fn set_priority(&mut self, priority: &str) {
        let priostr = String::from(priority.to_lowercase().trim());
        if priostr.starts_with("l") {
            self.priority = Priority::Low(String::from("low"));
        }

        if priostr.starts_with("m") {
            self.priority = Priority::Middle(String::from("middle"));
        }

        if priostr.starts_with("h") {
            self.priority = Priority::High(String::from("high"));
        }
    }

    pub fn get_priority(&self) -> String {
        match &self.priority {
            Priority::Low(str) => String::from(str.as_str()),
            Priority::Middle(str) => String::from(str.as_str()),
            Priority::High(str) => String::from(str.as_str()),
        }
    }

    pub fn set_progress(&mut self, progress: &str) {
        let progstr = String::from(progress.to_lowercase().trim());
        if progstr.starts_with("n") {
            self.progress = Progress::NotStart(String::from("notstart"));
        }

        if progstr.starts_with("i") {
            self.progress = Progress::InProcess(String::from("inprocess"));
        }

        if progstr.starts_with("f") {
            self.progress = Progress::Finished(String::from("finished"));
        }
    }

    pub fn get_progress(&self) -> String {
        match &self.progress {
            Progress::NotStart(str) => String::from(str.as_str()),
            Progress::InProcess(str) => String::from(str.as_str()),
            Progress::Finished(str) => String::from(str.as_str()),
        }
    }

    pub fn dump_work_item(&self) -> String {
        let index = self.index;
        let workitem = &self.workitem;
        let priority = match &self.priority {
            Priority::Low(l) => l,
            Priority::Middle(m) => m,
            Priority::High(h) => h,
        };
        let progress = match &self.progress {
            Progress::NotStart(n) => n,
            Progress::InProcess(i) => i,
            Progress::Finished(f) => f,
        };

        format!("{}|{}|{}|{}", index, workitem, priority, progress)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_schedule_filename() {}
}

fn load_item_from_string(record: &str) -> Item {
    let mut item = Item::new();
    let v: Vec<&str> = record.split('|').collect();

    item.set_index(v[0].parse::<u32>().unwrap());
    item.set_workitem(v[1]);
    item.set_priority(v[2]);
    item.set_progress(v[3]);

    item
}

fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where
    P: AsRef<Path>,
{
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}

pub fn load_schedule_items_from_file(filepath: &str) -> Vec<Item> {
    let mut res: Vec<Item> = vec![];
    if let Ok(lines) = read_lines(filepath) {
        for line in lines {
            if let Ok(schedule_line) = line {
                res.push(load_item_from_string(&schedule_line.trim()))
            }
        }
    }
    res
}

pub fn dump_schedule_items_to_file(items: &Vec<Item>, filepath: &str) {
    let mut file = File::create(filepath).unwrap();
    for item in items.iter() {
        file.write_all(item.dump_work_item().as_bytes())
            .expect("write dump work item failed!");
        file.write_all("\n".as_bytes())
            .expect("write new line failed!");
    }
    drop(file);
}
