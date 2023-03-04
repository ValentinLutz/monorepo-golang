use time::OffsetDateTime;

use crate::core::{
    model::order::{Order, OrderStatus},
    port::incoming::{OrderError, OrderGetter, OrderPlacer, OrdersGetter},
};

use super::order_id::OrderId;

pub struct OrderService;

impl OrdersGetter for OrderService {
    fn get_orders(&self, offset: isize, limit: isize) -> Result<Vec<Order>, OrderError> {
        return Ok(vec![]);
    }
}

impl OrderPlacer for OrderService {
    fn place_order(&self, order: Order) -> Result<Order, OrderError> {
        return Ok(order);
    }
}

impl OrderGetter for OrderService {
    fn get_order(&self, id: OrderId) -> Result<Order, OrderError> {
        return Ok(Order {
            order_id: id,
            creation_date: OffsetDateTime::now_utc(),
            status: OrderStatus::OrderPlaced,
            workflow: String::from(""),
            items: vec![],
        });
    }
}
