-- +goose Up
ALTER TABLE `nfile` 
ADD COLUMN `captured_at` DATETIME(6) AFTER `date_modified`,
ADD COLUMN `importance` TINYINT DEFAULT 3 AFTER `size`,
ADD COLUMN `metadata` JSON AFTER `hash`,
ADD COLUMN `created_at` DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),
ADD COLUMN `updated_at` DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6);

-- +goose Down
ALTER TABLE `nfile` 
DROP COLUMN `captured_at`,
DROP COLUMN `importance`,
DROP COLUMN `metadata`,
DROP COLUMN `created_at`,
DROP COLUMN `updated_at`;
