-- +goose Up
-- +goose StatementBegin
CREATE TABLE `traders`
(
    `uid`             VARCHAR(32)    NOT NULL,
    `pnl`             DECIMAL(10, 2) NOT NULL,
    `pnl_weekly`      DECIMAL(10, 2) NOT NULL,
    `pnl_monthly`     DECIMAL(10, 2) NOT NULL,
    `pnl_yearly`      DECIMAL(10, 2) NOT NULL,
    `roi`             DECIMAL(10, 2) NOT NULL,
    `roi_weekly`      DECIMAL(10, 2) NOT NULL,
    `roi_monthly`     DECIMAL(10, 2) NOT NULL,
    `roi_yearly`      DECIMAL(10, 2) NOT NULL,
    `position_shared` TINYINT(1) UNSIGNED NOT NULL,
    `publisher`       TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    `published_at`    TIMESTAMP NULL DEFAULT NULL,
    `created_at`      TIMESTAMP      NOT NULL,
    `updated_at`      TIMESTAMP      NOT NULL,
    PRIMARY KEY (`uid`),
    KEY               `position_shared_idx` (`position_shared`),
    KEY               `publisher_idx` (`publisher`)
) ENGINE=InnoDB;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `traders`;
-- +goose StatementEnd
