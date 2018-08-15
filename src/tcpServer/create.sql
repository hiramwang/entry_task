use mysql;
CREATE TABLE `shopee_test` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(40) NOT NULL,
  `password` varchar(40) NOT NULL,
  `nickname` varchar(20) DEFAULT NULL,
  `extend` varchar(40) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `account` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=8113496 DEFAULT CHARSET=utf8; 
insert into mysql.user(Host,User,Password) values("localhost","hiram",password("hiram"));
grant all privileges on mysql.* to hiram@localhost Identified by "hiram";
flush privileges;
