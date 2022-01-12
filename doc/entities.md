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

## find duplicates

    SELECT id, name, hash, size 
    FROM nostalgia.nfile as n 
    WHERE hash IN (SELECT hash FROM nostalgia.nfile WHERE id != n.id)

With Size

    select c.hash, c.s, n.size
    from
    (
    SELECT hash, sum(size) as s
    FROM nostalgia.nfile as n
    WHERE hash IN (SELECT hash FROM nostalgia.nfile WHERE id != n.id) group by hash) as c, nfile as n
    WHERE
    n.hash = c.hash
    group by hash, s, size
    order by size desc

> only pictures, 20321 distinct files repeated. total 251.54GB

