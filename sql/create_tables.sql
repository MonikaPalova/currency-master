CREATE DATABASE IF NOT EXISTS `currency-master`;

USE `currency-master`;

CREATE TABLE IF NOT EXISTS `USERS` (
    `username` VARCHAR(36) NOT NULL PRIMARY KEY,
    `password` VARCHAR(36) NOT NULL,
    `email` VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS `USER_ASSETS` (
    `username` VARCHAR(36) NOT NULL,
    `asset_id` VARCHAR(10) NOT NULL,
    `name` VARCHAR(36) NOT NULL,
    `quantity` FLOAT NOT NULL,
    FOREIGN KEY (username) REFERENCES USERS(username),
    CONSTRAINT PK_USER_ASSET PRIMARY KEY (username,asset_id)
);

CREATE TABLE IF NOT EXISTS `ACQUISITIONS` (
    `username` VARCHAR(36) NOT NULL,
    `asset_id` VARCHAR(10) NOT NULL,
    `quantity` FLOAT NOT NULL,
    `price_usd` FLOAT NOT NULL,
    `created` DATETIME NOT NULL,
    FOREIGN KEY (username) REFERENCES USERS(username),
    CONSTRAINT PK_USER_ASSET PRIMARY KEY (username,asset_id,created)
);

-- CREATE TABLE IF NOT EXISTS `SESSIONS` (
--     `id` VARCHAR(128) NOT NULL PRIMARY KEY,
--     `owner` VARCHAR(36) NOT NULL,
--     `expiration` DATETIME NOT NULL,
--     FOREIGN KEY (owner) REFERENCES USERS(username)
-- );