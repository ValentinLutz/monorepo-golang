use std::fmt;

use time::OffsetDateTime;

use crate::core::service::order_id::OrderId;

#[derive(Debug)]
pub enum OrderStatus {
    OrderPlaced,
    OrderInProgress,
    OrderCancelled,
    OrderCompleted,
}

impl fmt::Display for OrderStatus {
    fn fmt(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            OrderStatus::OrderPlaced => write!(formatter, "order_placed"),
            OrderStatus::OrderInProgress => write!(formatter, "order_in_progress"),
            OrderStatus::OrderCancelled => write!(formatter, "order_canceled"),
            OrderStatus::OrderCompleted => write!(formatter, "order_completed"),
        }
    }
}

pub struct Order {
    pub order_id: OrderId,
    pub creation_date: OffsetDateTime,
    pub status: OrderStatus,
    pub workflow: String,
    pub items: Vec<OrderItem>,
}

pub struct OrderItem {
    pub order_item_id: isize,
    pub name: String,
    pub creation_date: OffsetDateTime,
}
