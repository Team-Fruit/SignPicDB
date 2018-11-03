CREATE USER `signpic`@`%`;
GRANT INSERT,SELECT,UPDATE,DELETE ON `signpic_db`.* TO `signpic`@`%`;

CREATE DATABASE IF NOT EXISTS `signpic_db`;

CREATE TABLE IF NOT EXISTS `signpic_db`.`user` (
    `uuid`              CHAR(32)        NOT NULL,
    `username`          TINYTEXT        NOT NULL,
    `ip`                VARCHAR(255)    NOT NULL,
    `version_mod`       TINYTEXT        NOT NULL,
    `version_mod_mc`    TINYTEXT        NOT NULL,
    `version_mod_forge` TINYTEXT        NOT NULL,
    `version_mc`        TINYTEXT        NOT NULL,
    `version_forge`     TINYTEXT        NOT NULL,
    `message`           TEXT            NOT NULL,
    `created_at`        TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at`        TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ( `uuid` ),
);

