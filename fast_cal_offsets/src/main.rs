use calamine::{open_workbook, Reader, Xlsx};

fn main() {
    let mut excel: Xlsx<_> = open_workbook("sheets.xlsx").unwrap();

    if let Some(Ok(r)) = excel.worksheet_range("sheet1") {
        let mut cur_start_offset = 1;
        let mut cur_end_offset = 1;

        for row in r.rows() {
            let length_field = row[1].clone();
            let length_str = &length_field.to_string()[1..];
            let length_num = length_str.to_string().parse::<i32>().unwrap();

            cur_end_offset = cur_start_offset + length_num;
            println!(
                "field: {}, start offset: {}, end offset: {}",
                row[0], cur_start_offset, cur_end_offset
            );
            cur_start_offset = cur_end_offset + 1;
        }
    }
}
