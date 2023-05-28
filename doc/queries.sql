use nostalgia;

-- drop all
set foreign_key_checks=0;

drop table `nfile`;
drop table `ndirectory`;
drop table `nscan`;
drop table `nfile_ndirectory`;
drop table `nerror`;

set foreign_key_checks=1;

-- count all 

select count(*) as nfile from `nfile`;
select count(*) as ndirectory from `ndirectory`;
select count(*) as nscan from `nscan`;
select count(*) as nfile_directory from `nfile_ndirectory`;
select count(*) as nerror from `nerror`;

-- select top 5 all 

select * from `nfile` limit 5;  
select * from `ndirectory` limit 5;
select * from `nscan` limit 5;
select * from `nfile_ndirectory` limit 5;
select * from `nerror` limit 5;

-- schema 
 
\dt `nfile`;
\dt `ndirectory`;
\dt `nscan`;
\dt `nfile_ndirectory`;
\dt `nerror`;
