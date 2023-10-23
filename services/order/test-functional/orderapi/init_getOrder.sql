INSERT INTO order_service.order
    (order_id, customer_id, order_workflow, creation_date, order_status)
VALUES ('fdCDxjV9o!O-NONE-ZCTH5i6fWcA', '18bcf290-a61a-4338-808f-5759839b2056', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_placed');

INSERT INTO order_service.order_item
    (order_id, creation_date, order_item_name)
VALUES ('fdCDxjV9o!O-NONE-ZCTH5i6fWcA', '1970-01-01 00:00:00 +00:00', 'orange'),
       ('fdCDxjV9o!O-NONE-ZCTH5i6fWcA', '1970-01-01 00:00:00 +00:00', 'banana');