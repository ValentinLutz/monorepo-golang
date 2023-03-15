use time::OffsetDateTime;

pub struct OrderEntity {
    pub order_id: String,
    pub creation_date: OffsetDateTime,
    pub status: String,
    pub order_status: String,
}

pub struct OrderItemEntity {
    pub order_item_id: isize,
    pub name: String,
    pub creation_date: OffsetDateTime,
}