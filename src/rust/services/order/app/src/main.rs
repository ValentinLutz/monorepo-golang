use actix_web::{web, App, HttpServer};
use incoming::{
    openapi::routes::{dist, index, spec},
    orderapi::routes::{get_order, get_orders, post_orders},
};
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};

mod incoming;
mod core;

pub type DatabasePool = Pool<Postgres>;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let postgres_pool = PgPoolOptions::new()
        .max_connections(5)
        .connect("postgres://test:test@localhost:9432/dev_db")
        .await
        .expect("failed to build postgres connection pool");

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(postgres_pool.clone()))
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
