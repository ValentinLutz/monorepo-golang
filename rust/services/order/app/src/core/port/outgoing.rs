use crate::core::{model::order::Order, service::order_id::OrderId};
use async_trait::async_trait;
use std::sync::Arc;

pub type DynOrderRepository = Arc<dyn OrderRepository + Send + Sync>;

#[async_trait]
pub trait OrderRepository {
    async fn find_all_orders(&self, offset: i64, limit: i64) -> Result<Vec<Order>, String>;
    async fn find_order_by_id(&self, order_id: OrderId) -> Result<Order, String>;
    async fn save_order(&self, order: Order) -> Result<Order, String>;
}
