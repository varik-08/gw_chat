# Технологический стек

## Основные технологии

### Язык программирования
- **Go 1.24+**: основной язык для всех микросервисов
- **Строгая типизация**: использование Go generics для type-safe кода
- **Стандартная библиотека**: максимальное использование stdlib

### Транспорт и протоколы

#### Внешние интерфейсы
- **HTTP/REST**: стандартизированный API для веб-клиентов
- **WebSocket**: двунаправленная связь в реальном времени
- **TLS**: шифрование всех внешних соединений

#### Внутренние интерфейсы
- **gRPC**: синхронная межсервисная коммуникация
- **Protocol Buffers**: сериализация данных для gRPC
- **HTTP/2**: транспорт для gRPC

#### Асинхронная коммуникация
- **Apache Kafka**: event streaming платформа
- **KRaft mode**: Kafka без ZooKeeper
- **Совместимость с Redpanda**: альтернативная реализация Kafka API

## Оркестрация и workflow

### Temporal
- **Temporal Server**: workflow engine для сложных бизнес-процессов
- **Temporal Workers**: в каждом микросервисе для выполнения activities
- **Workflow-as-Code**: бизнес-логика в виде кода Go
- **Temporal UI**: веб-интерфейс для мониторинга workflows

#### Temporal особенности
- **Durable Execution**: workflows переживают перезапуски
- **Versioning**: безопасное обновление workflow кода
- **Saga Pattern**: встроенная поддержка компенсаций
- **Visibility**: детальное логирование выполнения

### Kafka для событий
- **Event Sourcing**: события как источник истины
- **CQRS**: разделение команд и запросов
- **Choreography**: fire-and-forget события
- **Partitioning**: по ключам `saga_id`, `chat_id`, `user_id`

#### Event-driven Choreography (технические детали)
**Назначение:** Децентрализованная координация между сервисами без центрального оркестратора.

**Kafka топики для хореографии:**
```yaml
# Real-time события
realtime.typing:        # typing indicators
  partitions: 20
  key: chat_id
  consumers: [gateway]

realtime.presence:      # user online/offline
  partitions: 10  
  key: user_id
  consumers: [gateway, analytics-service]

# Analytics события  
analytics.events:       # метрики и бизнес-аналитика
  partitions: 50
  key: user_id | chat_id
  consumers: [analytics-service, bi-service]

# Audit события
audit.events:          # логирование для compliance
  partitions: 30
  key: entity_id
  consumers: [audit-service, compliance-service]

# Уведомления
notifications.events:   # push, email, in-app notifications
  partitions: 20
  key: user_id
  consumers: [notification-service, email-service, push-service]

# Системные события
system.events:         # мониторинг, alerting
  partitions: 10
  key: service_name
  consumers: [monitoring-service, alerting-service]
```

**Consumer Groups:**
- **gateway-fanout**: WebSocket broadcast в real-time
- **analytics-collectors**: сбор метрик и KPI
- **audit-loggers**: запись аудитных событий
- **notification-processors**: обработка уведомлений
- **monitoring-agents**: системные метрики и алерты

**Технические особенности:**
- **At-least-once delivery**: Kafka producer `acks=all`
- **Idempotent consumers**: обработка дублей через unique constraints
- **Graceful degradation**: падение одного consumer не влияет на других
- **Schema evolution**: backward compatibility для event formats
- **Dead Letter Queues**: для stuck/failed messages

**Vs Temporal оркестрация:**
| Аспект | Choreography (Kafka) | Orchestration (Temporal) |
|--------|---------------------|-------------------------|
| **Latency** | < 10ms | 100-500ms |
| **Complexity** | Простые flows | Сложные workflows |
| **Guarantees** | At-least-once | Exactly-once + ACID |
| **Debugging** | Distributed logs | Temporal UI |
| **Use cases** | Real-time, analytics | Business transactions |

## Хранение данных

### PostgreSQL
- **Одна БД, множество схем**: логическое разделение по сервисам
- **ACID транзакции**: строгая консистентность внутри сервиса
- **Connection pooling**: pgbouncer или встроенные пулы
- **Миграции**: golang-migrate или аналог

#### Схемы БД
```
├── auth_schema (auth-service)
├── user_schema (user-service)  
├── chat_schema (chat-service)
├── message_schema (message-service)
└── temporal_schema (temporal-server)
```

### Паттерны работы с данными
- **Repository Pattern**: абстракция над БД
- **Outbox Pattern**: атомарность БД + события
- **Saga Pattern**: через Temporal workflows
- **Идемпотентность**: unique constraints и upserts

## Безопасность

### JWT и криптография
- **RS256**: асимметричное подписание токенов
- **JWKS**: публикация публичных ключей
- **Key rotation**: периодическое обновление ключей
- **Refresh tokens**: безопасное обновление доступа

### Межсервисная безопасность
- **Zero Trust**: проверка каждого запроса
- **Service-to-service auth**: JWT в заголовках gRPC
- **mTLS**: взаимная аутентификация (optional)

## Наблюдаемость (Observability)

### OpenTelemetry
- **Distributed Tracing**: сквозная трассировка запросов
- **Metrics**: бизнес и технические метрики
- **Logs**: структурированное логирование
- **OTLP экспорт**: в Jaeger, Prometheus, etc.

#### Компоненты
- **Jaeger**: трассировка и анализ производительности
- **Prometheus + Grafana**: метрики и дашборды
- **ELK Stack**: централизованные логи (optional)

### Корреляция
- **Trace ID**: передача через HTTP headers и Kafka headers
- **Span context**: для связи операций
- **Baggage**: дополнительные метаданные

## Конфигурация и развёртывание

### Конфигурация
- **Environment variables**: основной способ конфигурации
- **Feature flags**: `ENABLE_TEMPORAL`, `ENABLE_OUTBOX`, etc.
- **Config validation**: проверка при старте сервиса
- **Hot reload**: для некритичных настроек

### Контейнеризация
- **Docker**: контейнеры для всех сервисов
- **Multi-stage builds**: оптимизация размера образов
- **Health checks**: встроенные проверки здоровья
- **Graceful shutdown**: корректное завершение работы

### Оркестрация
- **Docker Compose**: для локальной разработки
- **Kubernetes**: для production (планы)
- **Service mesh**: Istio для продвинутого управления трафиком (future)

## Инструменты разработки

### Генерация кода
- **protoc**: генерация gRPC клиентов/серверов
- **mockgen**: моки для тестирования
- **sqlc**: генерация type-safe SQL кода

### Качество кода
- **golangci-lint**: статический анализ
- **gofmt**: форматирование кода
- **go mod**: управление зависимостями
- **Go test**: unit и integration тесты

### Local Development
```yaml
# docker-compose.micro.yml
services:
  - kafka (KRaft mode)
  - kafka-ui  
  - postgres
  - temporal-server
  - temporal-ui
  - jaeger
  - gateway
  - auth-service
  - user-service
  - chat-service
  - message-service
```

## Производительность и масштабирование

### Горизонтальное масштабирование
- **Stateless services**: все сервисы без состояния
- **Load balancing**: nginx или cloud LB
- **Consumer groups**: Kafka для распределения нагрузки
- **Connection pooling**: для БД соединений

### Оптимизации
- **Caching**: Redis для часто запрашиваемых данных
- **CDN**: для статических ресурсов
- **Compression**: gzip для HTTP, snappy для Kafka
- **Batching**: групповая обработка событий

### Надёжность
- **Circuit breakers**: защита от каскадных сбоев  
- **Retries**: с exponential backoff
- **Timeouts**: на всех внешних вызовах
- **Dead Letter Queues**: для необработанных сообщений
