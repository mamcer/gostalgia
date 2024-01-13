# Nostalgia

## Entities

nfile  
ndirectory  
nscan  
nfile_ndirectory  
nerror  

## Database

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

### Stash size

mac format Mib

    -- unique files size (size in disk)
    select sum(size)/1024/1024/1024/1024 from nfile

    -- repeated files size
    select sum(n.size)/1000/1000 from nfile_nscan as nfs, nfile as n where nfs.nfile_id = n.id

    select sum(n.size)/1024/1024/1024 from nfile_nscan as nfs, nfile as n where nfs.nfile_id = n.id

## Link

    ln -s /media/darkforce/stash/ stash

    ln -s /mnt/homunculus/pictures stash/2

## Times

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

## Search

type:

    image   photos & images
    doc     documents
    sheet   spreadsheets
    audio   audios
    video   videos
    zip     compressed files
    [any]

has the words:

date modified: 
between 
[any]

## Api

### ping 

returns 'pong', health status

GET 

request example

    ping

response example

    {
      "message": "pong"
    }

### search

Search for a resource, query, type, date after & before and paging

GET 

request example

    search?q=a&type=[image|doc|sheet|audio|video|zip|any]&after=1000-01-01&before=9999-12-31&page=1&per_page=50

response example

    search?q=a&type=doc&after=2010-01-01&before=2020-12-31&page=1&per_page=5
    {
      "directories": null,
      "files": [
        {
          "id": 2,
          "name": "Design Patterns Elements of Reusable Object-Oriented Software.pdf",
          "extension": "pdf",
          "path": "1/doc",
          "date_modified": "02-08-2012",
          "size": "4.3 MB",
          "hash": "8865aeb8efaa49a1700230e2cb1dee4c157800c8"
        },
        {
          "id": 1,
          "name": "Building_Maintainable_Software_SIG.pdf",
          "extension": "pdf",
          "path": "1/doc",
          "date_modified": "11-10-2016",
          "size": "6.5 MB",
          "hash": "3cf2bebbdadfe1a9fb6112c102553db0f1d7ed9b"
        }
      ],
      "page": 1,
      "per_page": 5,
      "query": "a",
      "total_directories": 0,
      "total_files": 2
    }

error codes

    404 if the query return no results

### files 

files representation

GET

request example 

    files/[id]

response example

    {
      "id": 3,
      "name": "coreutils.pdf",
      "extension": "pdf",
      "path": "1/doc",
      "date_modified": "04-01-2023",
      "size": "1.2 MB",
      "hash": "40877fd288bc8c6118518d6c5fe565d67658d24e"
    }

error codes

    404 if there is no file with id [id]

### filescount

return the total database file count

GET 

request example

    filescount

response example

    {
      "count": 5
    }

## directories

directories representation

GET

request example

    directories/[id]

response example

    {
      "id": 3,
      "name": "doc",
      "path": "1/doc",
      "date_modified": "21-05-2023",
      "size": "12.0 MB",
      "file_count": 3,
      "directory_count": 0,
      "parent_id": 2,
      "nscan_id": 1
    }

error codes

    404 if there is no directory with id [id]
 