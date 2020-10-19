mod schedule_item;
mod user_interface;

use std::process;
use structopt::StructOpt;

#[derive(Debug, StructOpt)]
struct Opt {
    #[structopt(short = "l", long = "list")]
    check: bool,
    #[structopt(short = "c", long = "change")]
    change: bool,
    #[structopt(short = "a", long = "add", default_value = "This is a new work item")]
    item: String,
    #[structopt(short = "p", long = "priority")]
    priority: Option<String>,
    #[structopt(short = "d", long = "date", default_value = "today")]
    date: String,
    #[structopt(short = "g", long = "progress")]
    progress: Option<String>,
    #[structopt(short = "i", long = "index", default_value = "1")]
    index: u32,
    #[structopt(short = "b", long = "bookpath", default_value = "bookkeep")]
    bookpath: String,
    #[structopt(short = "r", long = "remove")]
    remove: bool,
}

fn main() {
    let opt = Opt::from_args();
    let cur_date;

    if opt.date == "today" {
        cur_date = user_interface::get_current_date();
    } else {
        cur_date = opt.date;
    }

    //just list all the working items of a specific date
    if opt.check {
        user_interface::list_items_of_spcific_date(opt.bookpath.as_str(), cur_date.as_str());
        process::exit(0);
    }

    //change the priority or progress of a specific workitem
    if opt.change {
        if let Some(prio) = opt.priority {
            user_interface::change_priority_by_index(
                opt.index,
                prio.as_str(),
                opt.bookpath.as_str(),
                cur_date.as_str(),
            );
        }

        if let Some(prog) = opt.progress {
            user_interface::change_progress_by_index(
                opt.index,
                prog.as_str(),
                opt.bookpath.as_str(),
                cur_date.as_str(),
            );
        }
        process::exit(0);
    }

    if opt.remove {
        user_interface::delete_item_by_index(opt.index, opt.bookpath.as_str(), cur_date.as_str());
        process::exit(0);
    }

    match opt.priority {
        Some(prio) => user_interface::insert_new_item(
            opt.item.as_str(),
            prio.as_str(),
            opt.bookpath.as_str(),
            cur_date.as_str(),
        ),
        None => user_interface::insert_new_item(
            opt.item.as_str(),
            "low",
            opt.bookpath.as_str(),
            cur_date.as_str(),
        ),
    }
}
