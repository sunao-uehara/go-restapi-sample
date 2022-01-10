CREATE TABLE sample (
	id int(11) unsigned NOT NULL AUTO_INCREMENT,
	foo varchar(255) DEFAULT NULL,
	int_val int(11) DEFAULT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;	