# Entities

## NDirectory
	id
	name
	path
	[parent]
	size
	count

## NFile
	id
	name
	extension
	path
	modified
	size
	hash
	[directory]

# Database

## MySQL

    docker pull mysql:5.7.30
    docker run -p 3306:3306 --name nostalgia -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7.30

    docker exec -it nostalgia mysql -uroot -p
    create database nostalgia;

## MyCli

    mycli -h localhost -u root -D nostalgia -P 3306

## SQL

    use nostalgia;

    CREATE TABLE ndirectory (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL,
        `path` VARCHAR(260) NOT NULL,
        `parent_id` BIGINT UNSIGNED NOT NULL,
        `size` DECIMAL(15,4) NOT NULL,
        `count` BIGINT NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

    ALTER TABLE `ndirectory` 
    ADD CONSTRAINT fk_ndirectory_parent_directory 
    FOREIGN KEY (`parent_id`) 
    REFERENCES ndirectory(`id`);

    CREATE INDEX `idx_ndirectory_name` ON `ndirectory`(`name`);

    CREATE TABLE nfile (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL,
        `extension` VARCHAR(50) NOT NULL,
        `path` VARCHAR(260) NOT NULL,
        `modified` DATETIME NOT NULL,
        `size` DECIMAL(15,4) NOT NULL,
        `hash` VARCHAR(40) NOT NULL,
        `ndirectory_id` BIGINT UNSIGNED NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE=INNODB AUTO_INCREMENT=1540 DEFAULT CHARSET=utf8;

    ALTER TABLE `nfile` 
    ADD CONSTRAINT `fk_nfile_ndirectory` 
    FOREIGN KEY (`ndirectory_id`) 
    REFERENCES `ndirectory`(`id`);

    CREATE INDEX `idx_nfile_name` ON `nfile`(`name`);
    CREATE INDEX `idx_nfile_extension` ON `nfile`(`extension`);
    CREATE INDEX `idx_nfile_modified` ON `nfile`(`modified`);
