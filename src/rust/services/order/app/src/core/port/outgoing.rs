use crate::core::{model::order::Order, service::order_id::OrderId};

pub trait OrderRepoitory {
    fn find_all_orders(&self, offset: isize, limit: isize) -> Vec<Order>;
    fn find_order_by_id(&self, id: OrderId) -> Option<Order>;
    fn save_order(&self, order: Order) -> Order;
}
