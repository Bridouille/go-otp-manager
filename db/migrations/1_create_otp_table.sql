-- +migrate Up
CREATE TABLE `otp` (
	`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	`name` varchar(255) NOT NULL,
	`key` varchar(255) NOT NULL,
	`time` tinyint(4) NOT NULL,
	`digits` tinyint(4) NOT NULL
);

-- +migrate Down
DROP TABLE `otp`;