-- +goose Up
CREATE TABLE IF NOT EXISTS `ntag` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS `nfile` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `extension` VARCHAR(50),
    `path` VARCHAR(4096) NOT NULL,
    `date_modified` DATETIME(6),
    `size` BIGINT,
    `hash` VARCHAR(255) UNIQUE
);

CREATE TABLE IF NOT EXISTS `ndirectory` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `date_modified` DATETIME(6),
    `parent_directory_id` BIGINT,
    `full_path` VARCHAR(4096) NOT NULL,
    `size` BIGINT,
    `file_count` BIGINT,
    `directory_count` BIGINT,
    `is_source` BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (`parent_directory_id`) REFERENCES `ndirectory`(`id`)
);

CREATE TABLE IF NOT EXISTS `nscan` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `status` INT,
    `date_created` DATETIME(6),
    `existing_file_repeated_count` BIGINT,
    `internal_file_repeated_count` BIGINT,
    `file_count` BIGINT,
    `directory_count` BIGINT,
    `root_directory_id` BIGINT,
    FOREIGN KEY (`root_directory_id`) REFERENCES `ndirectory`(`id`)
);

CREATE TABLE IF NOT EXISTS `ntag_ndirectory` (
    `ntag_id` BIGINT,
    `ndirectory_id` BIGINT,
    PRIMARY KEY (`ntag_id`, `ndirectory_id`),
    FOREIGN KEY (`ntag_id`) REFERENCES `ntag`(`id`),
    FOREIGN KEY (`ndirectory_id`) REFERENCES `ndirectory`(`id`)
);

CREATE TABLE IF NOT EXISTS `ntag_nfile` (
    `ntag_id` BIGINT,
    `nfile_id` BIGINT,
    PRIMARY KEY (`ntag_id`, `nfile_id`),
    FOREIGN KEY (`ntag_id`) REFERENCES `ntag`(`id`),
    FOREIGN KEY (`nfile_id`) REFERENCES `nfile`(`id`)
);

CREATE TABLE IF NOT EXISTS `nfilenode` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255),
    `ndirectory_id` BIGINT,
    `nfile_id` BIGINT,
    `nscan_id` BIGINT,
    FOREIGN KEY (`ndirectory_id`) REFERENCES `ndirectory`(`id`),
    FOREIGN KEY (`nfile_id`) REFERENCES `nfile`(`id`),
    FOREIGN KEY (`nscan_id`) REFERENCES `nscan`(`id`)
);

-- +goose Down
DROP TABLE IF EXISTS `nfilenode`;
DROP TABLE IF EXISTS `ntag_nfile`;
DROP TABLE IF EXISTS `ntag_ndirectory`;
DROP TABLE IF EXISTS `nscan`;
DROP TABLE IF EXISTS `ndirectory`;
DROP TABLE IF EXISTS `nfile`;
DROP TABLE IF EXISTS `ntag`;
