CREATE DATABASE IF NOT EXISTS friendfy;

USE friendfy;

DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id int auto_increment PRIMARY KEY,
    name varchar(50) NOT NULL,
    nick varchar(50) NOT NULL unique,
    email varchar(50) NOT NULL unique,
    password varchar(80) NOT NULL,
    created_at timestamp NOT NULL default current_timestamp()
) ENGINE=InnoDB;