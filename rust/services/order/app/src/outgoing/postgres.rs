use async_trait::async_trait;
use sqlx::PgPool;

use crate::core::{
    model::order::Order, port::outgoing::OrderRepository, service::order_id::OrderId,
};

pub struct PostgresOrderRepository {
    postgres_pool: PgPool,
}

impl PostgresOrderRepository {
    pub fn new(postgres_pool: PgPool) -> Self {
        return PostgresOrderRepository {
            postgres_pool: postgres_pool,
        };
    }
}

#[async_trait]
impl OrderRepository for PostgresOrderRepository {
    async fn find_all_orders(&self, offset: i32, limit: i32) -> Result<Vec<Order>, String> {
        todo!()
    }

    async fn find_order_by_id(&self, order_id: OrderId) -> Result<Order, String> {
        todo!()
    }

    async fn save_order(&self, order: Order) -> Result<Order, String> {
        todo!()
    }
}

// pub async fn find_all_orders(
//     database_pool: &DatabasePool,
//     offset: i64,
//     limit: i64,
// ) -> Result<Vec<Order>, OrderError> {
//     return sqlx::query_as!(
//         OrderEntity,
//         "SELECT order_id, creation_date, order_status FROM order_service.order ORDER BY creation_date OFFSET $1 LIMIT $2",
//         offset,
//         limit
//     )
//     .fetch_all(database_pool)
//     .await
//     .map_err(|err| OrderError::from(err));
// }

// pub async fn save_order(database_pool: &DatabasePool, order: Order) -> Result<Order, OrderError> {
//     return Ok(order);
// }

// pub async fn find_order_by_id(
//     database_pool: &DatabasePool,
//     order_id: OrderId,
// ) -> Result<Order, OrderError> {
//     return sqlx::query_as!(
//         OrderEntity,
//         "SELECT order_id, creation_date, order_status FROM order_service.order WHERE order_id = $1",
//         order_id,
//     )
//     .fetch_one(database_pool)
//     .await
//     .map(|order_entity| map_order_entity_to_order(order_entity))
//     .map_err(|err| OrderError::from(err));
// }

// fn map_order_entity_to_order(order_entity: OrderEntity, order_item_entities: Vec<OrderItemEntity>) -> Order {
//     return Order {
//         order_id: order_entity.order_id,
//         creation_date: order_entity.creation_date,
//         status: OrderStatus::OrderPlaced,
//         order_status: order_entity.order_status,
//         items: vec![],
//     };
// }
