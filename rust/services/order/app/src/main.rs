use crate::core::port::incoming::DynOrderService;
use crate::core::service::order::OrderServiceImpl;
use crate::core::service::order_id::Region;
use axum::routing::{get, post};
use axum::Router;
use incoming::openapi::routes::{open_api_specs, swagger_files};
use incoming::orderapi::routes::{get_order, get_orders, post_orders};
use outgoing::postgres::PostgresOrderRepository;
use sqlx::postgres::PgPoolOptions;
use std::net::SocketAddr;
use std::sync::Arc;

mod core;
mod incoming;
mod outgoing;

#[tokio::main]
async fn main() {
    let postgres_pool = PgPoolOptions::new()
        .max_connections(5)
        .connect("postgres://test:test@localhost:5432/test")
        .await
        .expect("failed to build postgres connection pool");

    let order_service = Arc::new(OrderServiceImpl::new(
        Arc::new(PostgresOrderRepository::new(Arc::new(
            postgres_pool.clone(),
        ))),
        Region::NONE,
    )) as DynOrderService;

    let order_api = Router::new()
        .route("/orders", get(get_orders))
        .route("/orders", post(post_orders))
        .route("/orders/:order_id", get(get_order))
        .with_state(order_service);

    let router = Router::new()
        .nest_service("/swagger/:file", get(swagger_files))
        .nest_service("/openapi/:file", get(open_api_specs))
        .nest("/api", order_api);

    let socket_address = SocketAddr::from(([0, 0, 0, 0], 8080));
    axum::Server::bind(&socket_address)
        .serve(router.into_make_service())
        .await
        .unwrap();
}
