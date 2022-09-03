CREATE TABLE golang_reference_project.order
(
    order_id      varchar     NOT NULL UNIQUE,
    workflow      varchar     NOT NULL,
    creation_date TIMESTAMPTZ NOT NULL,
    order_status  varchar     NOT NULL,
    PRIMARY KEY (order_id)
);
