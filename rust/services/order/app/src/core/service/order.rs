use time::OffsetDateTime;

use crate::core::{
    model::order::{Order, OrderItem, OrderStatus},
    port::{incoming, outgoing::OrderRepository},
};

use rand::{thread_rng, Rng};

use super::order_id::{generate_order_id, OrderId, Region};

pub struct OrderService {
    order_repository: Box<dyn OrderRepository>,
    region: Region,
}

impl OrderService {
    pub fn new(order_repository: Box<dyn OrderRepository>, region: Region) -> Self {
        return OrderService {
            order_repository: order_repository,
            region: region,
        };
    }
}

impl incoming::OrderService for OrderService {
    fn get_orders(&self, offset: i32, limit: i32) -> Result<Vec<Order>, String> {
        return self.order_repository.find_all_orders(offset, limit);
    }

    fn place_order(&self, item_names: Vec<String>) -> Result<Order, String> {
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

        return self.order_repository.save_order(order);
    }

    fn get_order(&self, order_id: OrderId) -> Result<Order, String> {
        return self.order_repository.find_order_by_id(order_id);
    }
}
