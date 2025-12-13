-- Скрипт создания БД (запускается из-под postgres / суперпользователя)

DROP DATABASE IF EXISTS pc_configurator_db;

CREATE DATABASE pc_configurator_db
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    TEMPLATE = template0;
