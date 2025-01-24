use clap::Parser;
use itertools::Itertools;
use serde_json::Value;
use std::error::Error;
use std::fs::File;
use std::io;
use std::io::{BufReader, BufWriter, Read, Write};

/// Converts JSON object array to CSV
#[derive(Parser)]
#[command(version)]
struct Args {
    /// The input JSON file
    #[arg(short, long, value_name = "JSON")]
    input: Option<String>,

    /// The output CSV file
    #[arg(short, long, value_name = "CSV")]
    output: Option<String>,
}

fn main() -> Result<(), Box<dyn Error>> {
    let args = Args::parse();

    let mut reader: Box<dyn Read> = if let Some(path) = args.input {
        Box::new(BufReader::new(File::open(path)?))
    } else {
        Box::new(BufReader::new(io::stdin()))
    };

    let json: Value = serde_json::from_reader(&mut reader)?;

    let mut writer: Box<dyn Write> = if let Some(path) = args.output {
        Box::new(BufWriter::new(File::create(path)?))
    } else {
        Box::new(BufWriter::new(io::stdout()))
    };

    to_csv(&mut writer, &json)?;
    writer.flush()?;

    Ok(())
}

#[inline]
fn json_type(json: &Value) -> &'static str {
    match json {
        Value::Null => "Null",
        Value::Bool(_) => "Bool",
        Value::Number(_) => "Number",
        Value::String(_) => "String",
        Value::Array(_) => "Array",
        Value::Object(_) => "Object",
    }
}

fn to_csv(writer: impl Write, json: &Value) -> Result<(), Box<dyn Error>> {
    let mut writer = csv::Writer::from_writer(writer);

    let array = json
        .as_array()
        .ok_or(format!("expected array, got {}", json_type(json)))?;

    for value in array {
        if !matches!(value, Value::Object(_) | Value::Null) {
            return Err(format!("expected object or null, got {:?}", json_type(json)).into());
        }
    }

    let array = array
        .into_iter()
        .filter_map(Value::as_object)
        .collect::<Vec<_>>();

    let keys = array
        .iter()
        .flat_map(|object| object.keys())
        .unique()
        .collect::<Vec<_>>();

    writer.write_record(keys.iter())?;

    for object in array {
        writer.write_record(
            keys.iter()
                .map(|&key| object.get(key).map_or_else(String::new, Value::to_string)),
        )?;
    }

    Ok(())
}
