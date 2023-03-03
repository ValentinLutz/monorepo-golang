use actix_web::{HttpResponse, Responder};

#[actix_web::get("/orders")]
pub async fn get_orders() -> impl Responder {
    let orders = vec!["order1", "order2"];
    return HttpResponse::Ok().json(orders);
}

#[actix_web::post("/orders")]
pub async fn post_orders() -> impl Responder {
    let orders = vec!["order1", "order2"];
    return HttpResponse::Ok().json(orders);
}

#[actix_web::get("/orders/{order_id}")]
pub async fn get_order() -> impl Responder {
    let orders = vec!["order1", "order2"];
    return HttpResponse::Ok().json(orders);
}
