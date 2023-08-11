CREATE DATABASE `rttys` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;


-- rttys.account definition

CREATE TABLE `account` (
  `username` varchar(512) NOT NULL,
  `password` text NOT NULL,
  `admin` int(11) DEFAULT NULL,
  `tenant` varchar(512) DEFAULT NULL,
  `token` varchar(512) NOT NULL,
  `description` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- rttys.config definition

CREATE TABLE `config` (
  `name` varchar(512) NOT NULL,
  `value` text NOT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- rttys.device definition

CREATE TABLE `device` (
  `id` varchar(512) NOT NULL,
  `description` text NOT NULL,
  `online` datetime NOT NULL,
  `username` text NOT NULL,
  `tenant` text DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- rttys.tenant definition

CREATE TABLE `tenant` (
  `name` varchar(512) NOT NULL,
  `owner` varchar(512) DEFAULT NULL,
  `description` text DEFAULT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


INSERT INTO rttys.account (username, password, admin, tenant, token, description) VALUES('admin', 'admin', 2, NULL, 'tokenstring', 'abc');
