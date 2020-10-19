use super::schedule_item;
use chrono::{Datelike, Utc};
use std::env;
use std::path::Path;

/// add new tasks
pub fn insert_new_item(workitem: &str, priority: &str, bookpath: &str, datestr: &str) {
    //find the schedule file according to datestr
    let schedule_file = get_schedule_filename(bookpath, datestr);
    let mut new_item = schedule_item::Item::new();
    let mut sch_items;

    if Path::new(schedule_file.as_str()).exists() {
        //the schedule file exists
        sch_items = schedule_item::load_schedule_items_from_file(schedule_file.as_str());
        new_item.set_index((sch_items.len() + 1) as u32);
        new_item.set_workitem(workitem);
        new_item.set_priority(priority);
        sch_items.push(new_item);
    } else {
        //the schedule file not exists
        sch_items = Vec::new();
        new_item.set_workitem(workitem);
        new_item.set_priority(priority);
        sch_items.push(new_item);
    }

    schedule_item::dump_schedule_items_to_file(&sch_items, schedule_file.as_str());
}

// delete item by index
pub fn delete_item_by_index(index: u32, bookpath: &str, datestr: &str) {
    let schedule_file = get_schedule_filename(bookpath, datestr);

    if Path::new(schedule_file.as_str()).exists() {
        let mut sch_items = schedule_item::load_schedule_items_from_file(schedule_file.as_str());
        if index <= (sch_items.len() as u32) {
            sch_items.remove((index - 1) as usize);
        }

        if index < (sch_items.len() as u32) {
            for i in index - 1..(sch_items.len() as u32) {
                let cur_index = sch_items[i as usize].get_index();
                sch_items[i as usize].set_index(cur_index - 1);
            }
        }

        schedule_item::dump_schedule_items_to_file(&sch_items, schedule_file.as_str());
    }
}

/// priority change
pub fn change_priority_by_index(index: u32, priority: &str, bookpath: &str, datestr: &str) {
    let schedule_file = get_schedule_filename(bookpath, datestr);
    let mut sch_items;

    if Path::new(schedule_file.as_str()).exists() {
        //the schedule file exists
        sch_items = schedule_item::load_schedule_items_from_file(schedule_file.as_str());
        for item in sch_items.iter_mut() {
            if item.index == index {
                item.set_priority(priority);
            }
        }
        schedule_item::dump_schedule_items_to_file(&sch_items, schedule_file.as_str());
    }
}

/// progress change
pub fn change_progress_by_index(index: u32, progress: &str, bookpath: &str, datestr: &str) {
    let schedule_file = get_schedule_filename(bookpath, datestr);
    let mut sch_items;

    if Path::new(schedule_file.as_str()).exists() {
        //the schedule file exists
        sch_items = schedule_item::load_schedule_items_from_file(schedule_file.as_str());
        for item in sch_items.iter_mut() {
            if item.index == index {
                item.set_progress(progress);
            }
        }
        schedule_item::dump_schedule_items_to_file(&sch_items, schedule_file.as_str());
    }
}

/// list all the work items of a specific date
pub fn list_items_of_spcific_date(bookpath: &str, date: &str) {
    let book_path = get_schedule_filename(bookpath, date);
    if Path::new(book_path.as_str()).exists() {
        let sch_items = schedule_item::load_schedule_items_from_file(book_path.as_str());
        for item in sch_items.iter() {
            println!(
                "{}|{}|{}|{}",
                item.get_index(),
                item.get_workitem(),
                item.get_priority(),
                item.get_progress()
            );
        }
    }
}

/// get current date in YYYY-MM-DD
pub fn get_current_date() -> String {
    let now = Utc::now();
    let (_, year) = now.year_ce();
    format!("{}-{:02}-{:02}", year, now.month(), now.day())
}

/// get bookkeep filename according to a DATA string
fn get_schedule_filename(bookpath: &str, datestr: &str) -> String {
    let mut pure_path = String::from("");
    match env::home_dir() {
        Some(homepath) => {
            pure_path.push_str(&String::from(
                homepath.join(bookpath).join(datestr).to_str().unwrap(),
            ));
            pure_path.push_str(".txt")
        }
        _ => {}
    }

    pure_path
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;

    #[test]
    fn test_get_schedule_filename() {
        assert_eq!(
            get_schedule_filename("./bookkeep", "2020-01-01").as_str(),
            "./bookkeep/2020-01-01.txt"
        );
    }

    #[test]
    fn test_insert_new_item() {
        let sch_file_name = get_schedule_filename("./bookkeep/", "2020-01-01");

        if Path::new(sch_file_name.as_str()).exists() {
            fs::remove_file(sch_file_name).unwrap();
        }

        insert_new_item("this is the first item", "low", "./bookkeep/", "2020-01-01");
        insert_new_item(
            "this is the second item",
            "middle",
            "./bookkeep",
            "2020-01-01",
        );
        insert_new_item(
            "this is the third item",
            "high",
            "./bookkeep/",
            "2020-01-01",
        );
    }

    #[test]
    fn test_change_priority_by_index() {
        let sch_file_name = get_schedule_filename("./bookkeep/", "2020-01-01");

        if Path::new(sch_file_name.as_str()).exists() {
            fs::remove_file(sch_file_name).unwrap();
        }

        insert_new_item("this is the first item", "low", "./bookkeep/", "2020-01-01");
        insert_new_item(
            "this is the second item",
            "middle",
            "./bookkeep/",
            "2020-01-01",
        );
        insert_new_item(
            "this is the third item",
            "high",
            "./bookkeep/",
            "2020-01-01",
        );

        change_priority_by_index(2, "high", "./bookkeep/", "2020-01-01");
    }

    #[test]
    fn test_get_current_date() {}
}
