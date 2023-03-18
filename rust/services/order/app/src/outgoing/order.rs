use time::OffsetDateTime;

pub struct OrderEntity {
    pub order_id: String,
    pub creation_date: OffsetDateTime,
    pub modified_date: OffsetDateTime,
    pub order_status: String,
}

pub struct OrderItemEntity {
    pub order_item_id: i32,
    pub order_id: String,
    pub item_name: String,
    pub creation_date: OffsetDateTime,
    pub modified_date: OffsetDateTime,
}
