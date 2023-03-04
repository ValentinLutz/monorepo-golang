use crate::core::{model::order::Order, service::order_id::OrderId};

pub struct OrderError;

pub trait OrdersGetter {
    fn get_orders(&self, offset: isize, limit: isize) -> Result<Vec<Order>, OrderError>;
}

pub trait OrderPlacer {
    fn place_order(&self, order: Order) -> Result<Order, OrderError>;
}

pub trait OrderGetter {
    fn get_order(&self, id: OrderId) -> Result<Order, OrderError>;
}
