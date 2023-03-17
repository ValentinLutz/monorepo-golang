use axum::{
    body::{boxed, Full},
    extract::Path,
    http::{header::CONTENT_TYPE, StatusCode},
    response::Response,
};
use rust_embed::RustEmbed;

#[derive(RustEmbed)]
#[folder = "src/incoming/openapi/spec"]
struct OpenApiSpecs;

pub async fn open_api_specs(Path(file): Path<String>) -> Response {
    match OpenApiSpecs::get(file.as_str()) {
        Some(content) => {
            let body = boxed(Full::from(content.data));
            let mime = mime_guess::from_path(file).first_or_octet_stream();

            Response::builder()
                .header(CONTENT_TYPE, mime.as_ref())
                .body(body)
                .unwrap()
        }
        None => Response::builder()
            .status(StatusCode::NOT_FOUND)
            .body(boxed(Full::from("404 Not Found")))
            .unwrap(),
    }
}

#[derive(RustEmbed)]
#[folder = "src/incoming/openapi/swagger-ui"]
struct SwaggerUI;

pub async fn swagger_files(Path(file): Path<String>) -> Response {
    match SwaggerUI::get(file.as_str()) {
        Some(content) => {
            let body = boxed(Full::from(content.data));
            let mime = mime_guess::from_path(file).first_or_octet_stream();

            Response::builder()
                .header(CONTENT_TYPE, mime.as_ref())
                .body(body)
                .unwrap()
        }
        None => Response::builder()
            .status(StatusCode::NOT_FOUND)
            .body(boxed(Full::from("404 Not Found")))
            .unwrap(),
    }
}
