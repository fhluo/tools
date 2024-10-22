use std::error::Error;
use std::fs::{FileTimes, OpenOptions};
use std::path::Path;
use std::time::SystemTime;

use clap::Parser;

#[derive(Parser, Debug)]
struct Args {
    #[arg(required = true)]
    filenames: Vec<String>,
}

fn touch(filenames: Vec<String>) -> Result<(), Box<dyn Error>> {
    for filename in filenames {
        if !Path::new(&filename).exists() {
            OpenOptions::new().create(true).write(true).open(filename)?;
        } else {
            OpenOptions::new()
                .write(true)
                .open(filename)?
                .set_times({
                    let now = SystemTime::now();
                    FileTimes::new().set_accessed(now).set_modified(now)
                })?;
        }
    }

    Ok(())
}

fn main() {
    let args = Args::parse();

    if let Err(err) = touch(args.filenames) {
        eprintln!("{err}");
    }
}
