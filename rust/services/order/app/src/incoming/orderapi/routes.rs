use actix_web::{HttpResponse, Responder};

#[actix_web::get("/orders")]
pub async fn get_orders() -> impl Responder {
    let orders = vec!["order1", "order2"];
    return HttpResponse::Ok().json(orders);
}

#[actix_web::post("/orders")]
pub async fn post_orders() -> impl Responder {
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
    return HttpResponse::Ok().json(orders);
}

#[actix_web::get("/orders/{order_id}")]
pub async fn get_order() -> impl Responder {
    let orders = vec!["order1", "order2"];
    return HttpResponse::Ok().json(orders);
}
