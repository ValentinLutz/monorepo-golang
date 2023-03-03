use base64::{engine::general_purpose, Engine};
use std::{borrow::Borrow, fmt};
use time::OffsetDateTime;

#[derive(Debug)]
pub enum Region {
    NONE,
    EU,
    US,
}

impl fmt::Display for Region {
    fn fmt(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            Region::NONE => write!(formatter, "NONE"),
            Region::EU => write!(formatter, "EU"),
            Region::US => write!(formatter, "US"),
        }
    }
}

#[derive(Debug)]
enum Environment {
    LOCAL,
    CONATAINER,
    TEST,
    PROD,
}

#[derive(Debug, PartialEq)]
pub struct OrderId(pub String);

pub fn generate_order_id(region: Region, timestamp: OffsetDateTime, salt: String) -> OrderId {
    let value_to_hash = region.to_string() + timestamp.to_string().borrow() + salt.borrow();
    let md5_sum = md5::compute(value_to_hash);

    let base64_string = general_purpose::URL_SAFE_NO_PAD.encode(md5_sum.as_ref());
    let base64_without_underscore_and_hyphen = base64_string.replace("-", "!").replace("_", "*");

    let region_identifier = format!("-{}-", region);

    let base64_string_half_length = base64_without_underscore_and_hyphen.len() / 2;

    return OrderId(
        base64_without_underscore_and_hyphen[..base64_string_half_length].to_string()
            + region_identifier.borrow()
            + base64_without_underscore_and_hyphen[base64_string_half_length..]
                .to_string()
                .borrow(),
    );
}
