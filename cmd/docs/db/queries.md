# queries

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
    `modified`,
    `parent_id`,
    `size`,
    `count`)
    VALUES
    (
    "$",
    "/",
    now(),
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
