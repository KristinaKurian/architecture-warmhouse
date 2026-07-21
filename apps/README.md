# Smart Home Sensor Management API

Проект состоит из трёх сервисов:

- `app` — монолит Smart Home API на Go, порт `8080`;
- `temperature-api` — отдельный сервис случайной температуры, порт `8081`;
- `postgres` — PostgreSQL для хранения датчиков, порт `5432`.

## Требования

- Docker;
- Docker Compose.

## Запуск

```bash
./init.sh
```

Или напрямую:

```bash
docker compose up --build -d
```

При первом запуске PostgreSQL выполняет скрипт `./smart_home/init.sql`.
Если ранее уже создавался volume и требуется повторная инициализация базы:

```bash
docker compose down -v
docker compose up --build -d
```

## Проверка temperature-api

Запрос по названию комнаты:

```bash
curl "http://localhost:8081/temperature?location=Living%20Room"
```

Запрос по идентификатору датчика:

```bash
curl "http://localhost:8081/temperature?sensorId=2"
```

Соответствие значений:

| sensorId | location |
|---|---|
| `1` | `Living Room` |
| `2` | `Bedroom` |
| `3` | `Kitchen` |
| другое | `Unknown` |

Если `sensorId` не передан, он определяется по `location`. Для неизвестной комнаты используется `sensorId=0`.
Каждый последовательный запрос для одного датчика возвращает новое случайное значение от `15.0` до `30.0 °C`.

## Проверка Smart Home API

Импортируйте `smarthome-api.postman_collection.json` в Postman и выполните:

1. `Create Sensor`;
2. `Get All Sensors` несколько раз.

Поле `value` температурного датчика будет обновляться при каждом запросе списка.

Основные endpoints:

- `GET /health` — health check;
- `GET /api/v1/sensors` — все датчики с актуальной температурой;
- `GET /api/v1/sensors/:id` — датчик по ID;
- `POST /api/v1/sensors` — создать датчик;
- `PUT /api/v1/sensors/:id` — изменить датчик;
- `DELETE /api/v1/sensors/:id` — удалить датчик;
- `PATCH /api/v1/sensors/:id/value` — изменить сохранённое значение;
- `GET /api/v1/sensors/temperature/:location` — запрос температуры через монолит.

## Остановка

```bash
docker compose down
```
