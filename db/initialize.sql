CREATE USER `signpic`@`%` IDENTIFIED BY 'password';
GRANT INSERT,SELECT,UPDATE,DELETE ON `signpic_db`.* TO `signpic`@`%`;

CREATE DATABASE IF NOT EXISTS `signpic_db`;

CREATE TABLE IF NOT EXISTS `signpic_db`.`user` (
    `uuid`              CHAR(32)     NOT NULL,
    `username`          TINYTEXT     NOT NULL,
    `ip`                VARCHAR(255) NOT NULL,
    `message`           TEXT         NOT NULL,
    `created_at`        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_count`     INT UNSIGNED NOT NULL DEFAULT 1,

    PRIMARY KEY ( `uuid` )
);

CREATE TABLE IF NOT EXISTS `signpic_db`.`user__version_mc__version_mod` (
    `uuid`              CHAR(32)     NOT NULL,
    `version_mod`       VARCHAR(255) NOT NULL,
    `version_mod_mc`    TINYTEXT     NOT NULL,
    `version_mod_forge` TINYTEXT     NOT NULL,
    `version_mc`        VARCHAR(255) NOT NULL,
    `version_forge`     TINYTEXT     NOT NULL,
    `created_at`        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_count`     INT UNSIGNED NOT NULL DEFAULT 1,

    PRIMARY KEY ( `uuid`, `version_mc`, `version_mod` )
);

