struct Memory {
    orders: Vec<Order>,
}

impl OrderRepoitory for Memory {
    async fn find_all_orders(&self, offset: isize, limit: isize) -> Vec<Order> {
        return self.orders.clone();
    }

    async fn find_order_by_id(&self, id: OrderId) -> Option<Order> {
        return self
            .orders
            .iter()
            .find(|order| order.order_id == id)
            .cloned();
    }

    async fn save_order(&self, order: Order) -> Order {
        self.orders.push(order.clone());
        return order;
    }
}
