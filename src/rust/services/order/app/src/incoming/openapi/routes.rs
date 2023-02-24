use actix_web::{HttpResponse, Responder, web};
use rust_embed::RustEmbed;

#[derive(RustEmbed)]
#[folder = "src/incoming/openapi/swagger-ui"]
struct SwaggerUI;

#[actix_web::get("/swagger/")]
pub(crate) async fn index() -> impl Responder {
    match SwaggerUI::get("index.html") {
        Some(content) => HttpResponse::Ok()
            .body(content.data.into_owned()),
        None => HttpResponse::NotFound().body("404 Not Found"),
    }
}

#[actix_web::get("/swagger/{file}")]
pub(crate) async fn dist(path: web::Path<String>) -> impl Responder {
    match SwaggerUI::get(path.as_str()) {
        Some(content) => HttpResponse::Ok()
            .body(content.data.into_owned()),
        None => HttpResponse::NotFound().body("404 Not Found"),
    }
}


#[derive(RustEmbed)]
#[folder = "src/incoming/openapi/spec"]
struct OpenApiSpecs;

#[actix_web::get("/openapi/{file}")]
pub(crate) async fn spec(path: web::Path<String>) -> impl Responder {
    match OpenApiSpecs::get(path.as_str()) {
        Some(content) => HttpResponse::Ok()
            .body(content.data.into_owned()),
        None => HttpResponse::NotFound().body("404 Not Found"),
    }
}