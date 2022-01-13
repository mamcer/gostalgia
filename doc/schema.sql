use nostalgia;

-- nfile
CREATE TABLE `nfile` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
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
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `path` VARCHAR(260) NOT NULL,
    `size` BIGINT UNSIGNED NOT NULL,
    `file_count` INT UNSIGNED NOT NULL,
    `parent_id` BIGINT UNSIGNED NOT NULL,
    `nscan_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

CREATE INDEX `idx_nfile_name` ON `ndirectory`(`name`);

-- nscan
CREATE TABLE `nscan` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `date_created` DATETIME NOT NULL,
    `duration` INT UNSIGNED NULL,
    `file_count` INT UNSIGNED NULL,
    `directory_count` INT UNSIGNED NULL,
    `status` INT NOT NULL,                  -- done = 0, inprogress = 1, error = 2
    `root_directory_path` VARCHAR(260) NOT NULL,
    `root_directory_id` BIGINT UNSIGNED NULL,
    `retry_count` TINYINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

-- nfile_nscan
CREATE TABLE `nfile_nscan` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `nfile_id` BIGINT UNSIGNED NOT NULL,
    `ndirectory_id` BIGINT UNSIGNED NOT NULL,
    `nscan_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

-- nerror
CREATE TABLE `nerror` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `description` TEXT NOT NULL,
    `nscan_id` BIGINT UNSIGNED NOT NULL,
    `retry_count` TINYINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

-- ntag
CREATE TABLE `ntag` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `nfile_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

CREATE INDEX `idx_ntag_name` ON `ntag`(`name`);

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

-- fk nfiletag
ALTER TABLE `ntag` 
ADD CONSTRAINT `fk_ntag_nfile` 
FOREIGN KEY (`nfile_id`) 
REFERENCES `nfile`(`id`);

-- insert default values

SET FOREIGN_KEY_CHECKS=0;

ALTER TABLE `nostalgia`.`nfile` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`ndirectory` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`nscan` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`nfile_nscan` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`nerror` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`ntag` AUTO_INCREMENT = 0;

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`size`,
`file_count`,
`parent_id`,
`nscan_id`)
VALUES
(
"/",
"/",
0,
0,
0,
0);

SET FOREIGN_KEY_CHECKS=1;