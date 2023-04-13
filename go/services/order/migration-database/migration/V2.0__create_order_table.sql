CREATE TABLE IF NOT EXISTS order_service.order
(
    order_id       VARCHAR     NOT NULL UNIQUE,
    customer_id    UUID        NOT NULL,
    creation_date  TIMESTAMPTZ NOT NULL,
    modified_date  TIMESTAMPTZ NOT NULL DEFAULT now(),
    order_workflow VARCHAR     NOT NULL,
    order_status   VARCHAR     NOT NULL,
    PRIMARY KEY (order_id)
);

CREATE INDEX order_creation_date_idx
    ON order_service.order (creation_date);

CREATE TRIGGER update_order_modified_date
    BEFORE UPDATE
    ON order_service.order
    FOR EACH ROW
EXECUTE PROCEDURE order_service.update_modified_date();
