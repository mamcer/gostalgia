# db

## remove all data

```sql
SET FOREIGN_KEY_CHECKS=0;

TRUNCATE TABLE `nostalgia`.`nfile`;
DELETE FROM `nostalgia`.`ndirectory` where is_source = 0x00;
TRUNCATE TABLE `nostalgia`.`nfile_ndirectory`;
TRUNCATE TABLE `nostalgia`.`nscan`;
TRUNCATE TABLE `nostalgia`.`ntag`;
TRUNCATE TABLE `nostalgia`.`ntag_nfile`;
TRUNCATE TABLE `nostalgia`.`ntag_ndirectory`;

SET FOREIGN_KEY_CHECKS=1;
```

## drop all

```sql
use nostalgia;

set foreign_key_checks=0;

drop table `nfile`;
drop table `ndirectory`;
drop table `nscan`;
drop table `nfile_ndirectory`;
drop table `ntag`;
drop table `ntag_nfile`;
drop table `ntag_ndirectory`;

set foreign_key_checks=1;
```

## count all 

```sql
select count(*) as nfile from `nfile`;
select count(*) as ndirectory from `ndirectory`;
select count(*) as nscan from `nscan`;
select count(*) as nfile_directory from `nfile_ndirectory`;
```

## select top 20 all 

```sql
select * from `nfile` limit 20;  
select * from `ndirectory` limit 20;
select * from `nscan` limit 20;
select * from `nfile_ndirectory` limit 20;
```

## schema 

> mycli version

```bash
\dt `nfile`;
\dt `ndirectory`;
\dt `nscan`;
\dt `nfile_ndirectory`;
```

## stash size

Unique file size 

```sql
select sum(size)/1000/1000/1000 as size_GB from nfile
```

Size of repeated files in GB

```sql
select sum(s.r) as duplicated_GB 
from (select fd.nfile_id as id, (count(fd.nfile_id)-1)*f.size/1000/1000/1000 as r 
        from nfile_ndirectory as fd, nfile as f 
        where fd.nfile_id = f.id  
        group by id order by r desc) as s
```

Top 10 repeated files

```bash
select fd.nfile_id as id, f.name as name, count(fd.nfile_id)-1 as repeated_count, (count(fd.nfile_id)-1)*f.size/1000/1000/1000 as size_GB 
        from nfile_ndirectory as fd, nfile as f 
        where fd.nfile_id = f.id  
        group by id,name order by size_GB desc limit 10
```

 Duplicated files & szie

```sql
select nfd.nfile_id as i, count(nfd.id) as c, sum(nf.size) as s 
from nfile_ndirectory as nfd, nfile as nf
where nf.id = nfd.nfile_id
group by i 
order by c desc       
```

## remove last scan

```sql
SET FOREIGN_KEY_CHECKS=0;
select @scanid = 2
delete from nfile where id in (select nfile_id from nfile_ndirectory where nscan_id = @scanid)
delete from ndirectory where id in (select ndirectory_id from nfile_ndirectory where nscan_id = @scanid)
delete from nfile_ndirectory where nscan_id = @scanid
delete from nscan where id = @scanid
SET FOREIGN_KEY_CHECKS=1;
-- could leave zombie directories (directories that don't have files won't have an entry in nfile_ndirectory, no way to )
```