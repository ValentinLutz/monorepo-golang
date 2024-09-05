INSERT INTO order_service.order
    (order_id, customer_id, order_workflow, creation_date, order_status)
VALUES ('01J71WWQ9K1A2-NONE-1JB0VJ47X1W0Q', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_placed'),

       ('01J71WXMDHMWH-NONE-PVHE8560H324S', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_in_progress'),

       ('01J71X1VZRWF2-NONE-4Q25SM23WRP1G', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_canceled'),

       ('01J71WYVYCDDA-NONE-44EJWDYHQHCHZ', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_completed');

INSERT INTO order_service.order_item
    (order_id, creation_date, order_item_name)
VALUES ('01J71WWQ9K1A2-NONE-1JB0VJ47X1W0Q', '1970-01-01 00:00:00 +00:00', 'orange'),
       ('01J71WWQ9K1A2-NONE-1JB0VJ47X1W0Q', '1970-01-01 00:00:00 +00:00', 'banana'),

       ('01J71WXMDHMWH-NONE-PVHE8560H324S', '1970-01-01 00:00:00 +00:00', 'chocolate'),

       ('01J71X1VZRWF2-NONE-4Q25SM23WRP1G', '1970-01-01 00:00:00 +00:00', 'marshmallow'),

       ('01J71WYVYCDDA-NONE-44EJWDYHQHCHZ', '1970-01-01 00:00:00 +00:00', 'apple');