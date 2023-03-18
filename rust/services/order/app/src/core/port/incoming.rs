use std::sync::Arc;

use async_trait::async_trait;

use crate::core::{model::order::Order, service::order_id::OrderId};

pub type DynOrderService = Arc<dyn OrderService + Send + Sync>;

#[async_trait]
pub trait OrderService {
    async fn get_orders(&self, offset: i64, limit: i64) -> Result<Vec<Order>, String>;
    async fn place_order(&self, item_names: Vec<String>) -> Result<Order, String>;
    async fn get_order(&self, order_id: OrderId) -> Result<Order, String>;
}
