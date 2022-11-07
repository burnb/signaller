-- +goose Up
CREATE TABLE `traders`
(
    `uid`             VARCHAR(64)    NOT NULL,
    `pnl`             DECIMAL(10, 2) NOT NULL,
    `pnl_weekly`      DECIMAL(10, 2) NOT NULL,
    `pnl_monthly`     DECIMAL(10, 2) NOT NULL,
    `pnl_yearly`      DECIMAL(10, 2) NOT NULL,
    `roi`             DECIMAL(10, 2) NOT NULL,
    `roi_weekly`      DECIMAL(10, 2) NOT NULL,
    `roi_monthly`     DECIMAL(10, 2) NOT NULL,
    `roi_yearly`      DECIMAL(10, 2) NOT NULL,
    `position_shared` TINYINT(1) UNSIGNED NOT NULL,
    `publisher`       TINYINT(1) UNSIGNED NOT NULL,
    `published_at`    TIMESTAMP NULL DEFAULT NULL,
    `created_at`      TIMESTAMP      NOT NULL,
    `updated_at`      TIMESTAMP      NOT NULL,
    PRIMARY KEY (`uid`),
    KEY               `position_shared` (`position_shared`),
    KEY               `published_at` (`published_at`),
    KEY               `publisher` (`publisher`)
) ENGINE=InnoDB;

CREATE TABLE `positions`
(
    `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `trader_uid`  VARCHAR(64)    NOT NULL,
    `symbol`      VARCHAR(45)    NOT NULL,
    `entry_price` FLOAT          NOT NULL,
    `pnl`         DECIMAL(10, 2) NOT NULL,
    `roe`         DECIMAL(10, 2) NOT NULL,
    `amount`      FLOAT          NOT NULL,
    `exchange`    TINYINT(1) UNSIGNED NOT NULL,
    `long`        TINYINT(1) UNSIGNED NOT NULL,
    `leverage`    TINYINT(1) UNSIGNED NOT NULL,
    `margin_mode` TINYINT(1) UNSIGNED NOT NULL,
    `hedged`      TINYINT(1) UNSIGNED NOT NULL,
    `closed_at`   TIMESTAMP NULL DEFAULT NULL,
    `created_at`  TIMESTAMP      NOT NULL,
    `updated_at`  TIMESTAMP      NOT NULL,
    PRIMARY KEY (`id`),
    KEY           `trader_uid_closed_at` (`trader_uid`, `closed_at`),
    KEY           `symbol` (`symbol`),
    KEY           `long` (`long`),
    CONSTRAINT `trader_uid_traders_uid` FOREIGN KEY (`trader_uid`) REFERENCES `traders` (`uid`)
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE `positions`;
DROP TABLE `traders`;
