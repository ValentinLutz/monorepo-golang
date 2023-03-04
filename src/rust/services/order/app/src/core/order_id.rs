use base64::{engine::general_purpose, Engine};
use serde::{Serialize, Deserialize};
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

enum Environment {
    LOCAL,
    CONATAINER,
    TEST,
    PROD,
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct OrderId(pub String);

pub fn generate_order_id(region: Region, timestamp: OffsetDateTime, salt: &str) -> OrderId {
    let value_to_hash = region.to_string() + timestamp.to_string().borrow() + salt;
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

#[cfg(test)]
mod tests {
    use test_case::test_case;
    use time::{format_description::well_known::Rfc3339, OffsetDateTime};
    use crate::core::order_id::{generate_order_id, OrderId, Region};
    
    #[test_case(Region::NONE,  "1",  OrderId(String::from("eBdsGMJzvCr-NONE-*2YpETWfUnA")))]
    #[test_case(Region::NONE,  "101",  OrderId(String::from("*kNJT1sDk5G-NONE-seQgT4znBAw")))]
    #[test_case(Region::NONE,  "10101",  OrderId(String::from("nZtuGALKpL5-NONE-ARtCIu9N*FQ")))]
    #[test_case(Region::NONE,  "1010101",  OrderId(String::from("h8z56svkr5L-NONE-6QpfDV4oO9Q")))]
    #[test_case(Region::EU, "1",  OrderId(String::from("4fgM*2Mxlj4-EU-46VPq2kWqrw")))]
    #[test_case(Region::EU, "101",  OrderId(String::from("K3CYLwUu!hl-EU-wKbosSfnxzQ")))]
    #[test_case(Region::EU, "10101",  OrderId(String::from("O1L084PzJk1-EU-MkUf9QVu93Q")))]
    #[test_case(Region::EU, "1010101",  OrderId(String::from("hCxXuqKLbL7-EU-lO9r4PJKV6A")))]
    #[test_case(Region::US, "1",  OrderId(String::from("snlOa30GO*g-US-0kIrItux5DQ")))]
    #[test_case(Region::US, "101",  OrderId(String::from("q7rjY3nW0TP-US-sP*Kfimr66Q")))]
    #[test_case(Region::US, "10101",  OrderId(String::from("Rjh5MuK9KSK-US-JbiPhobhv2Q")))]
    #[test_case(Region::US, "1010101",  OrderId(String::from("flyKlyvBVMO-US-Eg7Q9e8q52g")))]
    fn test_generate_order_id(
        region: Region,
        salt: &str,
        expected: OrderId,
    ) {
        // GIVEN
        let timestamp = OffsetDateTime::parse("1980-01-01T00:00:00+00:00", &Rfc3339).unwrap();
         
        // WHEN
        let actual = generate_order_id(region, timestamp, salt);

        // THEN
        assert_eq!(actual, expected)
    }
}

