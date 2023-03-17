use serde::Serialize;
use time::OffsetDateTime;


#[derive(Serialize)]
pub struct OrderResponse {
    pub order_id: String,
    pub creation_date: OffsetDateTime,
    pub status: String,
    pub items: Vec<OrderItemResponse>,
}

#[derive(Serialize)]
pub struct OrderItemResponse {
    pub name: String,
}
