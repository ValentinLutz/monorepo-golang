use actix_web::{web, HttpResponse, Responder};
use rand::Rng;
use time::OffsetDateTime;

use crate::{
    core::{
        order,
        order_id::{generate_order_id, Region},
    },
    DatabasePool,
};

#[actix_web::get("/orders")]
pub async fn get_orders(database_pool: web::Data<DatabasePool>) -> impl Responder {
    let orders = vec!["order1", "order2"];
    return HttpResponse::Ok().json(orders);
}

#[actix_web::post("/orders")]
pub async fn post_orders(database_pool: web::Data<DatabasePool>) -> impl Responder {
    let timestamp = OffsetDateTime::now_utc();
    let order_id = generate_order_id(
        Region::NONE,
        timestamp,
        String::from(rand::thread_rng().gen()).as_str(),
    );

    let orders = order::Order {
        order_id: order_id,
        creation_date: timestamp,
        status: order::OrderStatus::OrderPlaced,
        workflow: String::from("default"),
        items: vec![],
    };

    let order = order::place_order(database_pool.get_ref(), orders)
        .await
        .unwrap();

    return HttpResponse::Ok().json(orders);
}

#[actix_web::get("/orders/{order_id}")]
pub async fn get_order(database_pool: web::Data<DatabasePool>) -> impl Responder {
    let orders = vec!["order1", "order2"];
    return HttpResponse::Ok().json(orders);
}
