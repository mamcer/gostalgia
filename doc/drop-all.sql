use nostalgia;

-- drop all
SET FOREIGN_KEY_CHECKS=0;

DROP TABLE `nfile`;
DROP TABLE `ndirectory`;
DROP TABLE `nscan`;
DROP TABLE `nfile_nscan`;
DROP TABLE `nerror`;
DROP TABLE `ntag`;

SET FOREIGN_KEY_CHECKS=1;