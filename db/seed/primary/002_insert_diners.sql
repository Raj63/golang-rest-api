INSERT INTO diners (table_no, name) VALUES (1, 'Mr. Smith')ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO diners (table_no, name) VALUES (2, 'Ms. Sunita Chaudary')ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO diners (table_no, name) VALUES (3, 'Dr. Giridhari Reddy')ON DUPLICATE KEY UPDATE updated_at=NOW();
