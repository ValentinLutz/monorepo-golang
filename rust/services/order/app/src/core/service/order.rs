use super::order_id::{generate_order_id, OrderId, Region};
use crate::core::{
    model::order::{Order, OrderItem, OrderStatus},
    port::{incoming::OrderService, outgoing::DynOrderRepository},
};
use async_trait::async_trait;
use rand::{thread_rng, Rng};
use time::OffsetDateTime;

#[derive(Clone)]
pub struct OrderServiceImpl {
    order_repository: DynOrderRepository,
    region: Region,
}

impl OrderServiceImpl {
    pub fn new(order_repository: DynOrderRepository, region: Region) -> Self {
        return OrderServiceImpl {
            order_repository: order_repository,
            region: region,
        };
    }
}

#[async_trait]
impl OrderService for OrderServiceImpl {
    async fn get_orders(&self, offset: i64, limit: i64) -> Result<Vec<Order>, String> {
        return self.order_repository.find_all_orders(offset, limit).await;
    }

    async fn place_order(&self, item_names: Vec<String>) -> Result<Order, String> {
        let creation_date = OffsetDateTime::now_utc();

        let order_id = generate_order_id(
            self.region,
            creation_date,
            thread_rng().gen::<i32>().to_string().as_str(),
        );

        let order = Order {
            order_id: order_id,
            creation_date: creation_date,
            status: OrderStatus::OrderPlaced,
            items: item_names
                .iter()
                .map(|item_name| {
                    return OrderItem {
                        order_item_id: 0,
                        name: item_name.to_string(),
                        creation_date: creation_date,
                    };
                })
                .collect(),
        };

        return self.order_repository.save_order(order).await;
    }

    async fn get_order(&self, order_id: OrderId) -> Result<Order, String> {
        return self.order_repository.find_order_by_id(order_id).await;
    }
}
