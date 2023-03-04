use core::port::outgoing::OrderRepoitory;

use actix_web::{web, App, HttpServer};

use incoming::openapi::routes::{dist, index, spec};
use incoming::orderapi::routes::{get_order, get_orders, post_orders};

mod core;
mod incoming;
mod outgoing;

struct AppState {
    order_repository: dyn OrderRepoitory,
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .app_data(web::Data::new(AppState {
                order_repository: String::from("Order Service"),
            }))
            .service(index)
            .service(dist)
            .service(spec)
            .service(
                web::scope("/api")
                    .service(get_orders)
                    .service(post_orders)
                    .service(get_order),
            )
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
