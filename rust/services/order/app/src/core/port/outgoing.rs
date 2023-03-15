use crate::core::{model::order::Order, service::order_id::OrderId};

pub trait OrderRepository {
    fn find_all_orders(&self, offset: i32, limit: i32) -> Result<Vec<Order>, String>;
    fn find_order_by_id(&self, order_id: OrderId) -> Result<Order, String>;
    fn save_order(&self, order: Order) -> Result<Order, String>;
}
