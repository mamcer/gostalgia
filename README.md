# Nostalgia

## Entities

nfile  
ndirectory  
nscan  
nfile_nscan  
nerror  

## MySQL

    docker pull mysql:latest
    docker run -p 3306:3306 --name nostalgia -e MYSQL_ROOT_PASSWORD=root -d mysql:latest

    docker exec -it nostalgia mysql -uroot -p
    create database nostalgia;

## MyCli

    mycli -h localhost -u root -D nostalgia -P 3306

## MySql connections

    SHOW PROCESSLIST;
    show global status;
    SHOW STATUS WHERE `variable_name` = 'Max_used_connections';
    show status where variable_name = 'threads_connected';

## Stash size

mac format Mib

    -- unique files size (size in disk)
    select sum(size)/1000/1000 from nfile

    -- repeated files size
    select sum(n.size)/1000/1000 from nfile_nscan as nfs, nfile as n where nfs.nfile_id = n.id