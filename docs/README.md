# nostalgia

## entities

`nfile`  
`ndirectory`  
`nscan`  
`nfile_ndirectory`  

## database

### MySQL

```bash
docker pull mysql:latest
docker run -p 3306:3306 --name nostalgia -e MYSQL_ROOT_PASSWORD=root -d mysql:latest
docker exec -it nostalgia mysql -uroot -p
create database nostalgia;
// run schema.sql
```

### MyCli

```bash
mycli -h localhost -u root -D nostalgia -P 3306
```

### MySql connections

```sql
SHOW PROCESSLIST;
SHOW global status;
SHOW STATUS WHERE `variable_name` = 'Max_used_connections';
SHOW status where variable_name = 'threads_connected';
```

## link

```bash
ln -s /media/darkforce/stash/ stash

ln -s /mnt/homunculus/pictures stash/2
```

## times

2022-01-22

```bash
docs
process finished: 4h56m38.687302278s
scan_id: 1, files: 274887, directories: 41753, existing files: 122007 (44%), errors: 1
updating directory size...[ok]
total size: 897316160587 bytes, 897.3 GB
```

2023-04-07

```bash
process finished: 4h29m23.65938061s
scan_id: 1, files: 329518, directories: 44109, existing files: 134574 (40%), errors: 1
updating directory size...[ok]
total size: 1241799735052 bytes, 1241.8 GB
```

pictures

```bash
process finished: 4h50m42.951757549s
scan_id: 2, files: 133896, directories: 983, existing files: 51306 (38%), errors: 0
updating directory size...[ok]
total size: 1115321065715 bytes, 1115.3 GB
```

2023-04-07

```bash
process finished: 3h50m38.861497854s
scan_id: 2, files: 145397, directories: 1155, existing files: 58606 (40%), errors: 0
updating directory size...[ok]
total size: 1184495799147 bytes, 1184.5 GB
```

2024-07-11 (nuc ssd - all homunculus)

```bash
OK (1h38m52.748935232s)

updating file size...OK (7.213911ms)

checking existing files...OK (1m50.100899339s)
file repeated count: 0 (0%)

scan process finished: 1h40m45.509446558s
persist changes...OK (1h37m37.018295924s)
```