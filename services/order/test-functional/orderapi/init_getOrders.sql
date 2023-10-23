INSERT INTO order_service.order
    (order_id, customer_id, order_workflow, creation_date, order_status)
VALUES ('IsQah2TkaqS-NONE-JewgL0Ye73g', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_placed'),

       ('Fs2VoM7ZhrK-NONE-vzTf7kaHbRA', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_in_progress'),

       ('sgy1K3*SXcv-NONE-eVbldUAYXnA', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_canceled'),

       ('F2P!criGu2L-NONE-fJ7bBFx1vHg', '44bd6239-7e3d-4d4a-90a0-7d4676a00f5c', 'default_workflow',
        '1970-01-01 00:00:00 +00:00', 'order_completed');

INSERT INTO order_service.order_item
    (order_id, creation_date, order_item_name)
VALUES ('IsQah2TkaqS-NONE-JewgL0Ye73g', '1970-01-01 00:00:00 +00:00', 'orange'),
       ('IsQah2TkaqS-NONE-JewgL0Ye73g', '1970-01-01 00:00:00 +00:00', 'banana'),

       ('Fs2VoM7ZhrK-NONE-vzTf7kaHbRA', '1970-01-01 00:00:00 +00:00', 'chocolate'),

       ('sgy1K3*SXcv-NONE-eVbldUAYXnA', '1970-01-01 00:00:00 +00:00', 'marshmallow'),

       ('F2P!criGu2L-NONE-fJ7bBFx1vHg', '1970-01-01 00:00:00 +00:00', 'apple');