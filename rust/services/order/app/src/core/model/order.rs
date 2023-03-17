use serde::Serialize;
use strum_macros::Display;
use time::OffsetDateTime;

use crate::core::service::order_id::OrderId;

#[derive(Debug, Display)]
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
    pub order_item_id: isize,
    pub name: String,
    pub creation_date: OffsetDateTime,
}
