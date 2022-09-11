CREATE TABLE golang_reference_project.order_item
(
    order_item_id INT GENERATED ALWAYS AS IDENTITY,
    creation_date TIMESTAMPTZ NOT NULL,
    modified_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    order_id      VARCHAR     NOT NULL,
    item_name     VARCHAR     NOT NULL,
    PRIMARY KEY (order_item_id),
    CONSTRAINT order_id FOREIGN KEY (order_id) REFERENCES golang_reference_project.order (order_id)
);

CREATE TRIGGER update_order_item_modified_date
    BEFORE UPDATE
    ON golang_reference_project.order
    FOR EACH ROW
EXECUTE PROCEDURE golang_reference_project.update_modified_date();
