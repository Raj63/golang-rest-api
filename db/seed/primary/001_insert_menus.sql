INSERT INTO menus (`name`, `description`, price) VALUES ('HCDB', 'Hyderabadi Chicken Dum Briyani', '20000')ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO menus (`name`, `description`, price) VALUES ('HMDB', 'Hyderabadi Mutton Dum Briyani', '28000')ON DUPLICATE KEY UPDATE updated_at=NOW();

INSERT INTO menus (`name`, `description`, price) VALUES ('MVB', 'Muglai Veg Briyani', '18050')ON DUPLICATE KEY UPDATE updated_at=NOW();
