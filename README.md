# Avito Pixel

Avito Pixel — это продукт для учёта количества посетителей сайта и мобильных приложений.
Сервис предоставляет аггрегированную по дням статистику по уникальным посетителям.

Решение является OpenSource. Предполагается наличие технических навыков у пользователей продукта, т.к. раскладка на сервера, запуск приложений выполняется пользователем самостоятельно.

Сервис состоит из следующих частей:
- Серверное приложение на Go
  - Инструкции по запуску
  - Сборщик событий
  - Аггрегатор результатов
- База данных Clickhouse
- Клиент для отправки событий на сервер
  - [Web (Javascript)](https://github.com/avito-tech/avito-pixel-web-client)

## Инструкция по запуску

Серверная часть состоит из 2-х методов: hit и inspect.

- inspect — для формирования отчётов

### hit
`hit` предназначен для сбора событий с клиентов


### report
`report` предназначен для сбора событий с клиентов


## Пример

Для демонстрации того, как запустить avito-pixel, подготовлен пример. См. директорию `./example`.
Для запуска примера понадобится docker compose, убедитесь, что он установлен на вашей машине.

Для запуска выполните команду
```
docker compose up --build
```

Команда запустит clickhouse и сервис для сбора статистики

Далее необходимо выполнить миграции, которые создадут неоюходимые таблицы
```
docker exec -it clickhouse clickhouse-client --queries-file /db/changelog/master/1-init-20-11-2023.sql
```

Для того, чтобы протестировать функциональность работы сервиса, необходимо послеать несколько запросов

Эмулируем событие с клиентов
```
curl -X POST http://localhost:3000/hit/ \
-H 'Content-Type: application/json' \
--data-raw '{ "type": "load" }'
```

Можно выполнить данный запрос несколько раз, затем проверить, что данные сохранились в базу и готовы к выдаче.
```
curl --location -X GET 'http://localhost:3000/report/json/?metric=visitors&from=2024-01-01&to=2024-01-31&interval=1'
```
или откройте страницу http://localhost:3000/report/html

**Для подключения скрипта на ваш сайт используйте инструкцию из [README.MD avito-pixel-web-client](https://github.com/avito-tech/avito-pixel-web-client)**
