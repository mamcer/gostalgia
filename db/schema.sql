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
    `date_modified` DATETIME NOT NULL,
    `size` BIGINT UNSIGNED NOT NULL,
    `file_count` INT UNSIGNED NOT NULL,
    `directory_count` INT UNSIGNED NOT NULL,
    `parent_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

CREATE INDEX `idx_ndirectory_name` ON `ndirectory`(`name`);
CREATE INDEX `idx_ndirectory_date_modified` ON `ndirectory`(`date_modified`);

-- nscan
CREATE TABLE `nscan` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `date_created` DATETIME NOT NULL,
    `duration` INT UNSIGNED NULL,
    `file_count` INT UNSIGNED NULL,
    `directory_count` INT UNSIGNED NULL,
    `file_repeated_count` INT UNSIGNED NULL,
    `status` INT NOT NULL,                  -- inprogress = 0, done = 1, error = 2
    `root_directory_id` BIGINT UNSIGNED NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

-- nfile_ndirectory
CREATE TABLE `nfile_ndirectory` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `nfile_id` BIGINT UNSIGNED NOT NULL,
    `ndirectory_id` BIGINT UNSIGNED NOT NULL,
    `nscan_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

-- fk ndirectory
ALTER TABLE `ndirectory` 
ADD CONSTRAINT `fk_ndirectory_ndirectory` 
FOREIGN KEY (`parent_id`) 
REFERENCES `ndirectory`(`id`);

-- fk nscan
ALTER TABLE `nscan` 
ADD CONSTRAINT `fk_nscan_ndirectory` 
FOREIGN KEY (`root_directory_id`) 
REFERENCES `ndirectory`(`id`);

-- fk nfile_ndirectory
ALTER TABLE `nfile_ndirectory` 
ADD CONSTRAINT `fk_nfile_ndirectory_nfile` 
FOREIGN KEY (`nfile_id`) 
REFERENCES `nfile`(`id`);

ALTER TABLE `nfile_ndirectory` 
ADD CONSTRAINT `fk_nfile_ndirectory_ndirectory` 
FOREIGN KEY (`ndirectory_id`) 
REFERENCES `ndirectory`(`id`);

ALTER TABLE `nfile_ndirectory` 
ADD CONSTRAINT `fk_nfile_ndirectory_nscan` 
FOREIGN KEY (`nscan_id`) 
REFERENCES `nscan`(`id`);

-- insert default values

SET FOREIGN_KEY_CHECKS=0;

ALTER TABLE `nostalgia`.`nfile` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`ndirectory` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`nscan` AUTO_INCREMENT = 0;
ALTER TABLE `nostalgia`.`nfile_ndirectory` AUTO_INCREMENT = 0;

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`date_modified`,
`size`,
`file_count`,
`directory_count`,
`parent_id`)
VALUES
(
"/",
"/",
now(),
0,
0,
0,
0);

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`date_modified`,
`size`,
`file_count`,
`directory_count`,
`parent_id`)
VALUES
(
"cd",
"/cd",
now(),
0,
0,
0,
1);

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`date_modified`,
`size`,
`file_count`,
`directory_count`,
`parent_id`)
VALUES
(
"iso",
"/iso",
now(),
0,
0,
0,
1);

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`date_modified`,
`size`,
`file_count`,
`directory_count`,
`parent_id`)
VALUES
(
"year",
"/year",
now(),
0,
0,
0,
1);

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`date_modified`,
`size`,
`file_count`,
`directory_count`,
`parent_id`)
VALUES
(
"music",
"/music",
now(),
0,
0,
0,
1);

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`date_modified`,
`size`,
`file_count`,
`directory_count`,
`parent_id`)
VALUES
(
"picture",
"/picture",
now(),
0,
0,
0,
1);

INSERT INTO `nostalgia`.`ndirectory`
(`name`,
`path`,
`date_modified`,
`size`,
`file_count`,
`directory_count`,
`parent_id`)
VALUES
(
"video",
"/video",
now(),
0,
0,
0,
1);

SET FOREIGN_KEY_CHECKS=1;
