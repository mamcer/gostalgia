
    -- nfile
    CREATE TABLE `nfile` (
        `id` BIGINT UNSIGNED NOT NULL,
        `name` VARCHAR(255) NOT NULL,
        `extension` VARCHAR(50) NOT NULL,
        `path` VARCHAR(260) NOT NULL,
        `date_modified` DATETIME NOT NULL,
        `size` BIGINT UNSIGNED NOT NULL,
        `hash` VARCHAR(40) NOT NULL,
        `ndirectory_id` BIGINT UNSIGNED NOT NULL,
        `nscan_id` BIGINT UNSIGNED NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

    CREATE INDEX `idx_nfile_name` ON `nfile`(`name`);
    CREATE INDEX `idx_nfile_extension` ON `nfile`(`extension`);
    CREATE INDEX `idx_nfile_date_modified` ON `nfile`(`date_modified`);
    CREATE INDEX `idx_nfile_hash` ON `nfile`(`hash`);

    -- ndirectory
    CREATE TABLE `ndirectory` (
        `id` BIGINT UNSIGNED NOT NULL,
        `name` VARCHAR(255) NOT NULL,
        `path` VARCHAR(260) NOT NULL,
        `size` DECIMAL(15,4) NOT NULL,
        `file_count` INT UNSIGNED NOT NULL,
        `parent_id` BIGINT UNSIGNED NOT NULL,
        `nscan_id` BIGINT UNSIGNED NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

    CREATE INDEX `idx_nfile_name` ON `ndirectory`(`name`);

    -- nscan
    CREATE TABLE `nscan` (
        `id` BIGINT UNSIGNED NOT NULL,
        `date_created` DATETIME NOT NULL,
        `duration` INT UNSIGNED NOT NULL,
        `file_count` INT UNSIGNED NOT NULL,
        `directory_count` INT UNSIGNED NOT NULL,
        `status` VARCHAR(10) NOT NULL,                  -- done, inprogress, error
        `root_directory_path` VARCHAR(260) NOT NULL,
        `root_directory_id` BIGINT UNSIGNED NOT NULL,
        `retry_count` TINYINT UNSIGNED NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

    -- nfile_nscan
    CREATE TABLE `nfile_nscan` (
        `id` BIGINT UNSIGNED NOT NULL,
        `nfile_id` BIGINT UNSIGNED NOT NULL,
        `ndirectory_id` BIGINT UNSIGNED NOT NULL,
        `nscan_id` BIGINT UNSIGNED NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

    -- nerror
    CREATE TABLE `nerror` (
        `id` BIGINT UNSIGNED NOT NULL,
        `description` TEXT NOT NULL,
        `nscan_id` BIGINT UNSIGNED NOT NULL,
        `retry_count` TINYINT UNSIGNED NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

    -- fk nfile
    ALTER TABLE `nfile` 
    ADD CONSTRAINT `fk_nfile_ndirectory` 
    FOREIGN KEY (`ndirectory_id`) 
    REFERENCES `ndirectory`(`id`);

    ALTER TABLE `nfile` 
    ADD CONSTRAINT `fk_nfile_nscan` 
    FOREIGN KEY (`nscan_id`) 
    REFERENCES `nscan`(`id`);

    -- fk ndirectory
    ALTER TABLE `ndirectory` 
    ADD CONSTRAINT `fk_ndirectory_ndirectory` 
    FOREIGN KEY (`parent_id`) 
    REFERENCES `ndirectory`(`id`);

    ALTER TABLE `ndirectory` 
    ADD CONSTRAINT `fk_ndirectory_nscan` 
    FOREIGN KEY (`nscan_id`) 
    REFERENCES `nscan`(`id`);

    -- fk nscan
    ALTER TABLE `nscan` 
    ADD CONSTRAINT `fk_nscan_ndirectory` 
    FOREIGN KEY (`root_directory_id`) 
    REFERENCES `ndirectory`(`id`);

    -- fk nfile_nscan
    ALTER TABLE `nfile_nscan` 
    ADD CONSTRAINT `fk_nfile_nscan_nfile` 
    FOREIGN KEY (`nfile_id`) 
    REFERENCES `nfile`(`id`);

    ALTER TABLE `nfile_nscan` 
    ADD CONSTRAINT `fk_nfile_nscan_ndirectory` 
    FOREIGN KEY (`ndirectory_id`) 
    REFERENCES `ndirectory`(`id`);

    ALTER TABLE `nfile_nscan` 
    ADD CONSTRAINT `fk_nfile_nscan_nscan` 
    FOREIGN KEY (`nscan_id`) 
    REFERENCES `nscan`(`id`);

    -- fk nerror
    ALTER TABLE `nerror` 
    ADD CONSTRAINT `fk_nerror_nscan` 
    FOREIGN KEY (`nscan_id`) 
    REFERENCES `nscan`(`id`);
