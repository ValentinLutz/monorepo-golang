use std::sync::Arc;

use async_trait::async_trait;

use crate::core::{model::order::Order, service::order_id::OrderId};

pub type DynOrderRepository = Arc<dyn OrderRepository + Send + Sync>;

#[async_trait]
pub trait OrderRepository {
    async fn find_all_orders(&self, offset: i32, limit: i32) -> Result<Vec<Order>, String>;
    async fn find_order_by_id(&self, order_id: OrderId) -> Result<Order, String>;
    async fn save_order(&self, order: Order) -> Result<Order, String>;
}
