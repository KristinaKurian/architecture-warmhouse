# Задание 4. Создание и документирование API

## 1. Тип API

Для взаимодействия компонентов системы используется **гибридный подход**, объединяющий синхронные REST API и асинхронное событийное взаимодействие.

### Синхронное взаимодействие — REST API

REST API с передачей данных в формате JSON используется для операций, результат которых необходим вызывающей стороне сразу:

* регистрация и авторизация пользователей;
* управление пользовательскими профилями, ролями и сессиями;
* управление домами, комнатами, зонами и участниками дома;
* регистрация, настройка, привязка и отвязка устройств;
* создание и редактирование сценариев;
* получение текущего состояния и истории телеметрии;
* создание команд и просмотр истории их выполнения;
* просмотр истории уведомлений;
* работа с каталогом, комплектами, камерами и видеоархивом.

Клиентские приложения обращаются к REST API микросервисов через API Gateway/BFF.

REST выбран для этих операций:
* редоставляет понятную модель «запрос — ответ»
* поддерживается веб- и мобильными клиентами
* позволяет использовать стандартные HTTP-коды и документируется с помощью OpenAPI.

### Асинхронное взаимодействие — Event-driven API

Для команд устройствам, телеметрии и доменных событий используется брокер сообщений **Kafka**.

Асинхронное взаимодействие применяется в следующих процессах:

* Device Command Service сохраняет команду и публикует её в Kafka;
* Partner Integration Service получает команду и доставляет её физическому устройству;
* Partner Integration Service публикует результат доставки и выполнения команды;
* Device Command Service получает события изменения статуса и обновляет историю команды;
* Partner Integration Service публикует телеметрические данные устройств;
* Telemetry Service получает и сохраняет телеметрию;
* Telemetry Service публикует события изменения состояния устройства;
* Scenario Service получает события изменения состояния, проверяет условия сценариев и создаёт команды через REST API Device Command Service;
* Device Management Service публикует события изменения доступности устройств;
* Notification Service реагирует на события устройств, команд, телеметрии и сценариев и отправляет пользователям уведомления;
* Scenario Service публикует события о результатах выполнения сценариев.

Такое взаимодействие 
* уменьшает связанность микросервисов
* не блокирует пользовательский запрос при недоступном устройстве
* позволяет независимо масштабировать обработчики телеметрии, команд, сценариев и уведомлений.

Физические устройства могут использовать MQTT, HTTP, CoAP или протокол конкретного производителя. Особенности внешних протоколов скрываются внутри Partner Integration Service.

Для внутреннего асинхронного взаимодействия используется единый набор версионируемых событийных контрактов Kafka, описанный с помощью AsyncAPI.

### Итоговое решение

| Взаимодействие                             | Технология                          | Документация          |
| ------------------------------------------ | ----------------------------------- | --------------------- |
| Web/Mobile UI → API Gateway → микросервисы | HTTPS, REST, JSON                   | OpenAPI 3.1.1         |
| Синхронные запросы между микросервисами    | HTTP/HTTPS, REST, JSON              | OpenAPI 3.1.1         |
| Команды, телеметрия и доменные события     | Kafka                               | AsyncAPI 3.1.0        |
| Partner Integration Service → устройства   | MQTT, HTTP, CoAP, API производителя | Документация адаптера |

## 2. Документация API

### OpenAPI

* [User Service](./diagrams/openapi/01-user-service.openapi.yaml)
* [House Service](./diagrams/openapi/02-house-service.openapi.yaml)
* [Device Management Service](./diagrams/openapi/03-device-management-service.openapi.yaml)
* [Device Command Service](./diagrams/openapi/04-device-command-service.openapi.yaml)
* [Telemetry Service](./diagrams/openapi/05-telemetry-service.openapi.yaml)
* [Scenario Service](./diagrams/openapi/06-scenario-service.openapi.yaml)
* [Partner Integration Service](./diagrams/openapi/07-partner-integration-service.openapi.yaml)
* [Notification Service](./diagrams/openapi/08-notification-service.openapi.yaml)
* [Video Service](./diagrams/openapi/09-video-service.openapi.yaml)
* [Catalog Service](./diagrams/openapi/10-catalog-service.openapi.yaml)

### AsyncAPI

* [Smart Home Event API](./diagrams/asyncapi/smart-home-events.asyncapi.yaml)

### Основные Kafka topics

| Topic                                 | Producer                    | Consumer                                     |
| ------------------------------------- | --------------------------- | -------------------------------------------- |
| `smart-home.device.commands.v1`       | Device Command Service      | Partner Integration Service                  |
| `smart-home.device.command-status.v1` | Partner Integration Service | Device Command Service, Notification Service |
| `smart-home.telemetry.readings.v1`    | Partner Integration Service | Telemetry Service                            |
| `smart-home.device.state-changed.v1`  | Telemetry Service           | Scenario Service, Notification Service       |
| `smart-home.device.status-changed.v1` | Device Management Service   | Notification Service                         |
| `smart-home.scenario.executed.v1`     | Scenario Service            | Notification Service                         |

### Общие правила контрактов

* REST API версионируются через префикс `/v1`.
* Для идентификаторов используется UUID.
* Для авторизации используется JWT Bearer Token.
* Успешное принятие команды для асинхронного выполнения возвращает HTTP `202 Accepted`.
* Команды поддерживают `idempotencyKey`, чтобы повторный запрос не создавал дубликат.
* Каждое событие содержит идентификатор, тип, версию, время возникновения, сервис-источник и `traceId`.