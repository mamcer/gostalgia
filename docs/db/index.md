# db

## remove all data

    SET FOREIGN_KEY_CHECKS=0;
    
    TRUNCATE TABLE `nostalgia`.`nfile`;
    TRUNCATE TABLE `nostalgia`.`ndirectory`;
    TRUNCATE TABLE `nostalgia`.`nfile_ndirectory`;
    TRUNCATE TABLE `nostalgia`.`nerror`;
    TRUNCATE TABLE `nostalgia`.`nscan`;
    
    SET FOREIGN_KEY_CHECKS=1;

## seed data

    SET FOREIGN_KEY_CHECKS=0;

    DELETE FROM `nostalgia`.`ndirectory`;

    ALTER TABLE `nostalgia`.`ndirectory` AUTO_INCREMENT = 0;

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
    "$",
    "/",
    now(),
    0,
    0,
    0,
    0);

    SET FOREIGN_KEY_CHECKS=1;

## drop all

    use nostalgia;

    set foreign_key_checks=0;

    drop table `nfile`;
    drop table `ndirectory`;
    drop table `nscan`;
    drop table `nfile_ndirectory`;
    drop table `nerror`;

    set foreign_key_checks=1;

## count all 

    select count(*) as nfile from `nfile`;
    select count(*) as ndirectory from `ndirectory`;
    select count(*) as nscan from `nscan`;
    select count(*) as nfile_directory from `nfile_ndirectory`;
    select count(*) as nerror from `nerror`;

## select top 5 all 

    select * from `nfile` limit 5;  
    select * from `ndirectory` limit 5;
    select * from `nscan` limit 5;
    select * from `nfile_ndirectory` limit 5;
    select * from `nerror` limit 5;

## schema 

> mycli version

    \dt `nfile`;
    \dt `ndirectory`;
    \dt `nscan`;
    \dt `nfile_ndirectory`;
    \dt `nerror`;

## stash size

MAC format

    -- unique files size (size in disk in TiB)
    select sum(size)/1024/1024/1024/1024 from nfile

Size of repeated files in GB

    select sum(s.r) as duplicated_GB 
    from (select fd.nfile_id as id, (count(fd.nfile_id)-1)*f.size/1000/1000/1000 as r 
            from nfile_ndirectory as fd, nfile as f 
            where fd.nfile_id = f.id  
            group by id order by r desc) as s


## db schema

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
            `name` VARCHAR(260) NOT NULL,
            `date_created` DATETIME NOT NULL,
            `duration` INT UNSIGNED NULL,
            `file_count` INT UNSIGNED NULL,
            `directory_count` INT UNSIGNED NULL,
            `file_repeated_count` INT UNSIGNED NULL,
            `status` INT NOT NULL,                  -- done = 1, inprogress = 2, error = 3
            `root_directory_id` BIGINT UNSIGNED NULL,
            `retry_count` TINYINT UNSIGNED NOT NULL,
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

        -- nerror
        CREATE TABLE `nerror` (
            `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
            `description` TEXT NOT NULL,
            `nscan_id` BIGINT UNSIGNED NOT NULL,
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

        -- fk nerror
        ALTER TABLE `nerror` 
        ADD CONSTRAINT `fk_nerror_nscan` 
        FOREIGN KEY (`nscan_id`) 
        REFERENCES `nscan`(`id`);

        -- insert default values

        SET FOREIGN_KEY_CHECKS=0;

        ALTER TABLE `nostalgia`.`nfile` AUTO_INCREMENT = 0;
        ALTER TABLE `nostalgia`.`ndirectory` AUTO_INCREMENT = 0;
        ALTER TABLE `nostalgia`.`nscan` AUTO_INCREMENT = 0;
        ALTER TABLE `nostalgia`.`nfile_ndirectory` AUTO_INCREMENT = 0;
        ALTER TABLE `nostalgia`.`nerror` AUTO_INCREMENT = 0;

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

        SET FOREIGN_KEY_CHECKS=1;