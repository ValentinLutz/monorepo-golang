use super::models::{OrderItemResponse, OrderRequest, OrderResponse};
use crate::core::port::incoming::DynOrderService;
use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
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

pub async fn post_orders(
    State(order_service): State<DynOrderService>,
    Json(order_request): Json<OrderRequest>,
) -> (StatusCode, Json<OrderResponse>) {
    let order_items: Vec<String> = order_request
        .items
        .iter()
        .map(|order_item_request| order_item_request.name.clone())
        .collect();

    let order = order_service.place_order(order_items).await.unwrap();
    let order_response = OrderResponse {
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
    };

    return (StatusCode::CREATED, Json(order_response));
}

pub async fn get_order(
    State(order_service): State<DynOrderService>,
    Path(order_id): Path<String>,
) -> (StatusCode, Json<OrderResponse>) {
    let order = order_service.get_order(order_id).await.unwrap();

    let order_response = OrderResponse {
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
    };

    return (StatusCode::OK, Json(order_response));
}
