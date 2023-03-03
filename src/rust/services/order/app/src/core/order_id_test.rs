#[cfg(test)]
mod tests {
    use test_case::test_case;
    use time::{format_description::well_known::Rfc3339, OffsetDateTime};

    use crate::core::order_id::{generate_order_id, OrderId, Region};

    #[test_case(Region::NONE,  String::from("1"),  OrderId(String::from("eBdsGMJzvCr-NONE-*2YpETWfUnA")))]
    #[test_case(Region::NONE,  String::from("101"),  OrderId(String::from("*kNJT1sDk5G-NONE-seQgT4znBAw")))]
    #[test_case(Region::NONE,  String::from("10101"),  OrderId(String::from("nZtuGALKpL5-NONE-ARtCIu9N*FQ")))]
    #[test_case(Region::NONE,  String::from("1010101"),  OrderId(String::from("h8z56svkr5L-NONE-6QpfDV4oO9Q")))]
    #[test_case(Region::EU, String::from("1"),  OrderId(String::from("4fgM*2Mxlj4-EU-46VPq2kWqrw")))]
    #[test_case(Region::EU, String::from("101"),  OrderId(String::from("K3CYLwUu!hl-EU-wKbosSfnxzQ")))]
    #[test_case(Region::EU, String::from("10101"),  OrderId(String::from("O1L084PzJk1-EU-MkUf9QVu93Q")))]
    #[test_case(Region::EU, String::from("1010101"),  OrderId(String::from("hCxXuqKLbL7-EU-lO9r4PJKV6A")))]
    #[test_case(Region::US, String::from("1"),  OrderId(String::from("snlOa30GO*g-US-0kIrItux5DQ")))]
    #[test_case(Region::US, String::from("101"),  OrderId(String::from("q7rjY3nW0TP-US-sP*Kfimr66Q")))]
    #[test_case(Region::US, String::from("10101"),  OrderId(String::from("Rjh5MuK9KSK-US-JbiPhobhv2Q")))]
    #[test_case(Region::US, String::from("1010101"),  OrderId(String::from("flyKlyvBVMO-US-Eg7Q9e8q52g")))]
    fn test_generate_order_id(
        region: Region,
        salt: String,
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
