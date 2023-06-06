CREATE TABLE IF NOT EXISTS `orders` (
  `id` BIGINT auto_increment NOT NULL,
  `diner_id` BIGINT NOT NULL,
  `menu_id` BIGINT NOT NULL,
  `quantity` int NOT NULL DEFAULT 1,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `orders_diner_id_IDX` (`diner_id`),
  CONSTRAINT uc_orders_diner_id_and_menu_id UNIQUE (`diner_id`, `menu_id`),
  FOREIGN KEY (menu_id) REFERENCES menus (id) ON DELETE CASCADE,
  FOREIGN KEY (diner_id) REFERENCES diners (id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
