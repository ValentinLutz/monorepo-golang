use serde::{Deserialize, Serialize};
use time::OffsetDateTime;

#[derive(Serialize)]
pub struct OrderResponse {
    pub order_id: String,
    #[serde(with = "time::serde::rfc3339")]
    pub creation_date: OffsetDateTime,
    pub status: String,
    pub items: Vec<OrderItemResponse>,
}

#[derive(Serialize)]
pub struct OrderItemResponse {
    pub name: String,
}

#[derive(Deserialize)]
pub struct OrderRequest {
    pub items: Vec<OrderItemRequest>,
}

#[derive(Deserialize)]
pub struct OrderItemRequest {
    pub name: String,
}
