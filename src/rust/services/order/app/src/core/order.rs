use std::fmt;

use serde::{Deserialize, Serialize};
use time::OffsetDateTime;

use crate::DatabasePool;

use super::order_id::OrderId;

#[derive(Debug, Serialize, Deserialize)]
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

#[derive(Debug, Serialize, Deserialize)]
pub struct Order {
    pub order_id: OrderId,
    pub creation_date: OffsetDateTime,
    pub status: OrderStatus,
    pub workflow: String,
    pub items: Vec<OrderItem>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct OrderItem {
    pub order_item_id: isize,
    pub name: String,
    pub creation_date: OffsetDateTime,
}

pub struct OrderError;

pub async fn get_orders(
    database_pool: &DatabasePool,
    offset: isize,
    limit: isize,
) -> Result<Vec<Order>, OrderError> {
    return sqlx::query_as!(
        Order,
        r#"
        SELECT order_id, creation_date, status, workflow, items
        FROM orders
        OFFSET $1
        LIMIT $2
        "#,
        offset,
        limit
    )
    .fetch_all(&database_pool)
    .await?;
}

pub async fn place_order(database_pool: &DatabasePool, order: Order) -> Result<Order, OrderError> {
    return Ok(order);
}

pub async fn get_order(database_pool: &DatabasePool, id: OrderId) -> Result<Order, OrderError> {
    return Ok(Order {
        order_id: id,
        creation_date: OffsetDateTime::now_utc(),
        status: OrderStatus::OrderPlaced,
        workflow: String::from(""),
        items: vec![],
    });
}
