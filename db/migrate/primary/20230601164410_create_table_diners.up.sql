CREATE TABLE IF NOT EXISTS `diners` (
  `id` BIGINT auto_increment NOT NULL,
  `table_no` int NOT NULL,
  `name` varchar(250) NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `diners_name_IDX` (`name`),
  INDEX `diners_table_no_IDX` (`table_no`),
  CONSTRAINT uc_diners_tableno_name UNIQUE (`table_no`, `name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
