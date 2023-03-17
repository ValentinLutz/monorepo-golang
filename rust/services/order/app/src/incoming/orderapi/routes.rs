use axum::{extract::State, http::StatusCode, Json};

use crate::{
    core::port::incoming::DynOrderService,
    incoming::orderapi::models::{OrderItemResponse, OrderResponse},
};

pub async fn get_orders(
    State(order_service): State<DynOrderService>,
) -> (StatusCode, Json<Vec<OrderResponse>>) {
    let orders: Vec<OrderResponse> = order_service
        .get_orders(0, 10)
        .await
        .unwrap()
        .iter()
        .map(|order| OrderResponse {
            order_id: order.order_id.to_string(),
            creation_date: order.creation_date,
            status: order.status.to_string(),
            items: order
                .items
                .iter()
                .map(|order_item| OrderItemResponse {
                    name: order_item.name.to_string(),
                })
                .collect(),
        })
        .collect();
    return (StatusCode::OK, Json(orders));
}

pub async fn post_orders() -> (StatusCode, Json<Vec<&'static str>>) {
    // let timestamp = OffsetDateTime::now_utc();
    // let order_id = generate_order_id(
    //     Region::NONE,
    //     timestamp,
    //     String::from(rand::thread_rng().gen()).as_str(),
    // );

    // let orders = Order {
    //     order_id: order_id,
    //     creation_date: timestamp,
    //     status: OrderStatus::OrderPlaced,
    //     items: vec![],
    // };

    // let order = place_order(orders)
    //     .await
    //     .unwrap();

    // return HttpResponse::Ok().json(orders);

    let orders = vec!["order1", "order2"];
    return (StatusCode::OK, Json(orders));
}

pub async fn get_order() -> (StatusCode, Json<Vec<&'static str>>) {
    let orders = vec!["order1", "order2"];
    return (StatusCode::OK, Json(orders));
}
