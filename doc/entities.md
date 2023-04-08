# Entities

## remove all data

SET FOREIGN_KEY_CHECKS=0;

TRUNCATE TABLE `nostalgia`.`nfile`;
TRUNCATE TABLE `nostalgia`.`ndirectory`;

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


## pictures

total 126487 files in 957 directories
process finished: 4h17m49.674226637s
