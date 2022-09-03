CREATE TABLE golang_reference_project.order_item
(
    order_item_id INT GENERATED ALWAYS AS IDENTITY,
    order_id      varchar     NOT NULL,
    creation_date TIMESTAMPTZ NOT NULL,
    item_name     varchar     NOT NULL,
    CONSTRAINT order_id FOREIGN KEY (order_id) REFERENCES "order" (order_id),
    PRIMARY KEY (order_item_id)
);
