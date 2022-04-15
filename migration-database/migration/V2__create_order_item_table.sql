CREATE TABLE golang_reference_project.order_item
(
    id            INT GENERATED ALWAYS AS IDENTITY,
    order_id      varchar     NOT NULL,
    creation_date TIMESTAMPTZ NOT NULL,
    item_name     varchar     NOT NULL,
    CONSTRAINT order_id FOREIGN KEY (order_id) REFERENCES "order" (id),
    PRIMARY KEY (id)
);
