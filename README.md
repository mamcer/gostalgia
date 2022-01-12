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

