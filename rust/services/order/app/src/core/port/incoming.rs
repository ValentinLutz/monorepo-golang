use crate::core::{model::order::Order, service::order_id::OrderId};

pub trait OrderService {
    fn get_orders(&self, offset: i32, limit: i32) -> Result<Vec<Order>, String>;
    fn place_order(&self, item_names: Vec<String>) -> Result<Order, String>;
    fn get_order(&self, order_id: OrderId) -> Result<Order, String>;
}
