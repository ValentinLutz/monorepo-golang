use strum_macros::{Display, EnumString};
use time::OffsetDateTime;

use crate::core::service::order_id::OrderId;

#[derive(Debug, Display, EnumString)]
pub enum OrderStatus {
    OrderPlaced,
    // OrderInProgress,
    // OrderCancelled,
    // OrderCompleted,
}

pub struct Order {
    pub order_id: OrderId,
    pub creation_date: OffsetDateTime,
    pub status: OrderStatus,
    pub items: Vec<OrderItem>,
}

pub struct OrderItem {
    pub order_item_id: i32,
    pub name: String,
    pub creation_date: OffsetDateTime,
}
