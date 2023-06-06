INSERT INTO orders (diner_id, menu_id, quantity) VALUES (1, 2, 1)ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO orders (diner_id, menu_id, quantity) VALUES (1, 3, 2)ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO orders (diner_id, menu_id, quantity) VALUES (2, 1, 1)ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO orders (diner_id, menu_id, quantity) VALUES (2, 2, 1)ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO orders (diner_id, menu_id, quantity) VALUES (2, 3, 1)ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO orders (diner_id, menu_id, quantity) VALUES (3, 1, 2)ON DUPLICATE KEY UPDATE updated_at=NOW();
