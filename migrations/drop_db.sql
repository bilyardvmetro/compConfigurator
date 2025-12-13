-- Аккуратно гасим соединения и дропаем базу
-- Запускать из-под postgres

-- Сначала отключаем всех клиентов от нужной БД
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'pc_configurator_db'
  AND pid <> pg_backend_pid();

-- Теперь можно удалить
DROP DATABASE IF EXISTS pc_configurator_db;
