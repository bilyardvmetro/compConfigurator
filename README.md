# План этапов (roadmap)

## Этап 0. Каркас и договорённости

- выбор стека (router, config, логирование, миграции, DB access)

- структура проекта + соглашения (слои, naming, ошибки)

- healthcheck + swagger/openapi заготовка (по желанию)

## Этап 1. Инфраструктура проекта

- cmd/api/main.go

- конфиг (env) + валидация

- подключение к PostgreSQL

- миграции (goose / migrate)

- graceful shutdown

- базовый middleware (request-id, logging, recover)

## Этап 2. Модуль авторизации (минимально)

- roles, users

- регистрации/логин (JWT или session token)

- middleware авторизации + RBAC (USER/ADMIN)

## Этап 3. Каталог комплектующих (read-only API)

- CRUD (минимум GET-list/GET-by-id) по: CPU, GPU, MB, RAM, Drive, PSU, Case, Cooler

- фильтры под совместимость (socket, form-factor, ram_type/freq, nvme/sata, габариты)

- пагинация + сортировка

## Этап 4. Сборки

- CRUD по assemblies

- добавить/удалить RAM и диски (таблицы связей)

- эндпоинт “проверить совместимость” (дергает is_assembly_compatible)

- эндпоинт “пересчитать цену” (дергает recalc_assembly_total_price)

## Этап 5. Магазины и офферы

- CRUD для shops

- офферы product_offers

- эндпоинт “минимальная цена на компонент”

- “цена сборки по минимальным офферам”

## Этап 6. Админка и справочники

- управление сокетами/форм-факторами

- связи cooler_sockets, case_form_factor_support

- RBAC на админские роуты

## Этап 7. Тесты и качество

- unit для сервисов

- интеграционные тесты с testcontainers

- линтеры, gofmt, CI