# nostalgia

## entities

nfile  
ndirectory  
nscan  
nfile_ndirectory  
nerror  

## database

### MySQL

    docker pull mysql:latest
    docker run -p 3306:3306 --name nostalgia -e MYSQL_ROOT_PASSWORD=root -d mysql:latest

    docker exec -it nostalgia mysql -uroot -p
    create database nostalgia;

    // run schema.sql

### MyCli

    mycli -h localhost -u root -D nostalgia -P 3306

### MySql connections

    SHOW PROCESSLIST;
    show global status;
    SHOW STATUS WHERE `variable_name` = 'Max_used_connections';
    show status where variable_name = 'threads_connected';

## link

    ln -s /media/darkforce/stash/ stash

    ln -s /mnt/homunculus/pictures stash/2

## times

2022-01-22

    docs
    process finished: 4h56m38.687302278s
    scan_id: 1, files: 274887, directories: 41753, existing files: 122007 (44%), errors: 1
    updating directory size...[ok]
    total size: 897316160587 bytes, 897.3 GB

2023-04-07

    process finished: 4h29m23.65938061s
    scan_id: 1, files: 329518, directories: 44109, existing files: 134574 (40%), errors: 1
    updating directory size...[ok]
    total size: 1241799735052 bytes, 1241.8 GB

pictures

    process finished: 4h50m42.951757549s
    scan_id: 2, files: 133896, directories: 983, existing files: 51306 (38%), errors: 0
    updating directory size...[ok]
    total size: 1115321065715 bytes, 1115.3 GB

2023-04-07

    process finished: 3h50m38.861497854s
    scan_id: 2, files: 145397, directories: 1155, existing files: 58606 (40%), errors: 0
    updating directory size...[ok]
    total size: 1184495799147 bytes, 1184.5 GB
