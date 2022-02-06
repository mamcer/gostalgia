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
    select sum(size)/1024/1024/1024/1024 from nfile

    -- repeated files size
    select sum(n.size)/1000/1000 from nfile_nscan as nfs, nfile as n where nfs.nfile_id = n.id


    SELECT id, name, hash, size
    FROM nostalgia.nfile as n
    WHERE hash IN (SELECT hash FROM nostalgia.nfile WHERE id != n.id)
    
    select c.hash, c.s, n.size
    from
    (
    SELECT hash, sum(size) as s
    FROM nostalgia.nfile as n
    WHERE hash IN (SELECT hash FROM nostalgia.nfile WHERE id != n.id) group by hash) as c, nfile as n
    WHERE
    n.hash = c.hash
    order by s desc

## Link

    ln -s /media/darkforce/stash/ stash

    select sum(n.size)/1024/1024/1024 from nfile_nscan as nfs, nfile as n where nfs.nfile_id = n.id

## Times

2022-01-22

docs
process finished: 4h56m38.687302278s
scan_id: 1, files: 274887, directories: 41753, existing files: 122007 (44%), errors: 1
updating directory size...[ok]
total size: 897316160587 bytes, 897.3 GB

pictures
process finished: 4h50m42.951757549s
scan_id: 2, files: 133896, directories: 983, existing files: 51306 (38%), errors: 0
updating directory size...[ok]
total size: 1115321065715 bytes, 1115.3 GB
