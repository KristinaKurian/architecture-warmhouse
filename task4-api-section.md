# Задание 4. Создание и документирование API

## 1. Тип API

Для взаимодействия компонентов системы используется **гибридный подход**:

### Синхронное взаимодействие — REST API

REST API по HTTPS с передачей данных в JSON используется для операций, результат которых необходим вызывающей стороне сразу:

- регистрация и авторизация пользователей;
- управление домами, комнатами и зонами;
- регистрация и настройка устройств;
- создание и редактирование сценариев;
- получение текущего состояния и истории телеметрии;
- просмотр истории команд и уведомлений;
- работа с каталогом, камерами и видеоархивом.

Клиентские приложения обращаются к REST API через API Gateway/BFF. Внутренние синхронные вызовы между микросервисами также могут выполняться по HTTP, например для проверки существования дома или capabilities устройства.

REST выбран для этих операций, потому что он предоставляет понятную модель запрос–ответ, хорошо поддерживается веб- и мобильными клиентами, позволяет использовать стандартные HTTP-коды и легко документируется с помощью OpenAPI/Swagger.

### Асинхронное взаимодействие — Event-driven API

Для команд устройствам, телеметрии и доменных событий используется брокер сообщений **Kafka**.

Асинхронный обмен применяется для следующих процессов:

- Device Command Service публикует команды устройствам;
- Partner Integration Service доставляет команды и публикует результат выполнения;
- Partner Integration Service передаёт телеметрию устройств;
- Telemetry Service публикует события изменения состояния;
- Scenario Service реагирует на изменение состояния и запускает автоматизацию;
- Notification Service реагирует на события и отправляет уведомления.

Такое взаимодействие уменьшает связанность микросервисов, не блокирует пользовательский запрос при медленном или недоступном устройстве и позволяет независимо масштабировать обработчики телеметрии, команд и уведомлений.

Физические устройства могут использовать MQTT, HTTP, CoAP или протокол производителя. Эти различия скрываются внутри Partner Integration Service. Для внутреннего взаимодействия микросервисов используется единый событийный контракт Kafka.

### Итоговое решение

| Взаимодействие | Технология | Документация |
|---|---|---|
| Web/Mobile UI → API Gateway → микросервисы | HTTPS, REST, JSON | OpenAPI 3.1.1 |
| Синхронные запросы между микросервисами | HTTP, REST, JSON | OpenAPI 3.1.1 |
| Команды, телеметрия и доменные события | Kafka | AsyncAPI 3.1.0 |
| Partner Integration Service → устройства | MQTT, HTTP, CoAP, API производителя | Документация адаптера |

## 2. Документация API

### OpenAPI — синхронные REST API

- [01-user-service.openapi.yaml](./openapi/01-user-service.openapi.yaml)
- [02-house-service.openapi.yaml](./openapi/02-house-service.openapi.yaml)
- [03-device-management-service.openapi.yaml](./openapi/03-device-management-service.openapi.yaml)
- [04-device-command-service.openapi.yaml](./openapi/04-device-command-service.openapi.yaml)
- [05-telemetry-service.openapi.yaml](./openapi/05-telemetry-service.openapi.yaml)
- [06-scenario-service.openapi.yaml](./openapi/06-scenario-service.openapi.yaml)
- [07-partner-integration-service.openapi.yaml](./openapi/07-partner-integration-service.openapi.yaml)
- [08-notification-service.openapi.yaml](./openapi/08-notification-service.openapi.yaml)
- [09-video-service.openapi.yaml](./openapi/09-video-service.openapi.yaml)
- [10-catalog-service.openapi.yaml](./openapi/10-catalog-service.openapi.yaml)

### AsyncAPI — асинхронные события

- [Smart Home Event API](./asyncapi/smart-home-events.asyncapi.yaml)

### Основные Kafka topics

| Topic | Producer | Consumer |
|---|---|---|
| `smart-home.device.commands.v1` | Device Command Service | Partner Integration Service |
| `smart-home.device.command-status.v1` | Partner Integration Service | Device Command Service |
| `smart-home.telemetry.readings.v1` | Partner Integration Service | Telemetry Service |
| `smart-home.device.state-changed.v1` | Telemetry Service | Scenario Service, Notification Service |
| `smart-home.device.status-changed.v1` | Device Management Service | Notification Service |
| `smart-home.scenario.executed.v1` | Scenario Service | Notification Service |

### Общие правила контрактов

- REST API версионируются через префикс `/v1`.
- Для идентификаторов используется UUID.
- Даты и время передаются в формате ISO 8601 / RFC 3339.
- Для авторизации используется JWT Bearer Token.
- Успешное принятие асинхронной команды возвращает `202 Accepted`.
- Команды поддерживают `idempotencyKey`, чтобы повторный запрос не создавал дубликат.
- Каждое событие содержит идентификатор, тип, версию, время возникновения, сервис-источник и trace ID.
- Изменения, нарушающие обратную совместимость, публикуются в новой версии REST API или новом Kafka topic.
