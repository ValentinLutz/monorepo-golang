use super::order::OrderItemEntity;
use crate::{
    core::{
        model::order::{Order, OrderItem, OrderStatus},
        port::outgoing::OrderRepository,
        service::order_id::OrderId,
    },
    outgoing::order::OrderEntity,
};
use async_trait::async_trait;
use sqlx::PgPool;
use std::{collections::HashMap, str::FromStr, sync::Arc};

pub struct PostgresOrderRepository {
    postgres_pool: Arc<PgPool>,
}

impl PostgresOrderRepository {
    pub fn new(postgres_pool: Arc<PgPool>) -> Self {
        return PostgresOrderRepository {
            postgres_pool: postgres_pool,
        };
    }
}

#[async_trait]
impl OrderRepository for PostgresOrderRepository {
    async fn find_all_orders(&self, offset: i64, limit: i64) -> Result<Vec<Order>, String> {
        let order_entitites = sqlx::query_as!(
            OrderEntity,
            "SELECT order_id, creation_date, order_status, modified_date FROM order_service.order ORDER BY creation_date OFFSET $1 LIMIT $2",
            offset,
            limit
        )
        .fetch_all(self.postgres_pool.as_ref())
        .await
        .unwrap();

        let order_ids = order_entitites
            .iter()
            .map(|order_entity| order_entity.order_id.clone())
            .collect::<Vec<String>>();

        let order_item_entitites = sqlx::query_as!(
            OrderItemEntity,
            "SELECT order_item_id, order_id, creation_date, item_name, modified_date FROM order_service.order_item WHERE order_id = ANY($1)",
            order_ids.as_slice()
        )
        .fetch_all(self.postgres_pool.as_ref())
        .await
        .unwrap();

        return Ok(new_orders(order_entitites, order_item_entitites));
    }

    async fn find_order_by_id(&self, order_id: OrderId) -> Result<Order, String> {
        let order_entity = sqlx::query_as!(
            OrderEntity,
            "SELECT order_id, creation_date, order_status, modified_date FROM order_service.order WHERE order_id = $1",
            order_id,
            )
            .fetch_one(self.postgres_pool.as_ref())
            .await
            .unwrap();

        let order_item_entitites = sqlx::query_as!(
            OrderItemEntity,
            "SELECT order_item_id, order_id, creation_date, item_name, modified_date FROM order_service.order_item WHERE order_id = $1",
            order_id,
            )
            .fetch_all(self.postgres_pool.as_ref())
            .await
            .unwrap();

        return Ok(new_order(order_entity, order_item_entitites));
    }

    async fn save_order(&self, order: Order) -> Result<Order, String> {
        let order_entity = OrderEntity {
            order_id: order.order_id.to_string(),
            creation_date: order.creation_date,
            order_status: order.status.to_string(),
            modified_date: order.creation_date,
        };

        let order_item_entitites: Vec<_> = order
            .items
            .iter()
            .map(|order_item| {
                return OrderItemEntity {
                    order_item_id: order_item.order_item_id,
                    order_id: order.order_id.to_string(),
                    item_name: order_item.name.to_string(),
                    creation_date: order_item.creation_date,
                    modified_date: order_item.creation_date,
                };
            })
            .collect();

        let mut transaction = self.postgres_pool.begin().await.unwrap();

        sqlx::query!(
            "INSERT INTO order_service.order (order_id, creation_date, order_status) VALUES ($1, $2, $3)",
            order_entity.order_id,
            order_entity.creation_date,
            order_entity.order_status
        ).execute(&mut transaction)
        .await
        .unwrap();

        for order_item_entity in order_item_entitites {
            sqlx::query!(
                "INSERT INTO order_service.order_item (order_id, item_name, creation_date) VALUES ($1, $2, $3)",
                order_item_entity.order_id,
                order_item_entity.item_name,
                order_item_entity.creation_date
            ).execute(&mut transaction)
            .await
            .unwrap();
        }

        transaction.commit().await.unwrap();

        return Ok(order);
    }
}

fn new_order(order_entity: OrderEntity, order_item_entities: Vec<OrderItemEntity>) -> Order {
    let mut order_items: Vec<OrderItem> = Vec::new();
    for order_item_entity in order_item_entities {
        order_items.push(OrderItem {
            order_item_id: order_item_entity.order_item_id,
            name: order_item_entity.item_name,
            creation_date: order_item_entity.creation_date,
        });
    }

    return Order {
        order_id: order_entity.order_id,
        creation_date: order_entity.creation_date,
        status: OrderStatus::from_str(order_entity.order_status.as_str()).unwrap(),
        items: order_items,
    };
}

fn new_orders(
    order_entities: Vec<OrderEntity>,
    order_item_entities: Vec<OrderItemEntity>,
) -> Vec<Order> {
    let mut order_id_to_order_items: HashMap<String, Vec<OrderItem>> = HashMap::new();
    for order_item_entity in order_item_entities {
        order_id_to_order_items
            .entry(order_item_entity.order_id)
            .or_insert(Vec::new())
            .push(OrderItem {
                order_item_id: order_item_entity.order_item_id,
                name: order_item_entity.item_name,
                creation_date: order_item_entity.creation_date,
            });
    }

    let mut orders: Vec<Order> = Vec::new();
    for order_entity in order_entities {
        orders.push(Order {
            order_id: order_entity.order_id.clone(),
            creation_date: order_entity.creation_date,
            status: OrderStatus::from_str(order_entity.order_status.as_str()).unwrap(),
            items: order_id_to_order_items
                .remove(&order_entity.order_id)
                .unwrap(),
        });
    }

    return orders;
}
