CREATE TABLE IF NOT EXISTS `menus` (
  `id` BIGINT auto_increment NOT NULL,
  `name` varchar(120) NOT NULL,
  `description` varchar(255) NOT NULL,
  `price` int NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `menus_name_IDX` (`name`),
  CONSTRAINT uc_menus_name UNIQUE (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
