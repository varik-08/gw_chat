# Архитектурные принципы

## Основные принципы

### Service per Bounded Context
**Принцип:** Каждый микросервис соответствует одному ограниченному контексту из Domain-Driven Design.

**Реализация:**
- `auth-service`: контекст аутентификации и авторизации
- `user-service`: контекст пользователей и профилей
- `chat-service`: контекст чатов и участников
- `message-service`: контекст сообщений и доставки
- `gateway`: контекст протоколов и маршрутизации

**Преимущества:**
- Чёткие границы ответственности
- Независимая разработка и развёртывание
- Естественная декомпозиция по команды

### API Gateway Pattern
**Принцип:** Единая точка входа для всех внешних запросов.

**Реализация:**
- HTTP/REST API для веб-клиентов
- WebSocket для real-time коммуникации
- JWT валидация на gateway уровне
- Маршрутизация в соответствующие микросервисы
- Rate limiting и CORS

**Преимущества:**
- Упрощение клиентской логики
- Централизованная безопасность
- Версионирование API
- Мониторинг и логирование

### Database per Service
**Принцип:** Каждый сервис владеет своими данными и не имеет прямого доступа к данным других сервисов.

**Реализация:**
- Отдельная схема PostgreSQL для каждого сервиса
- Обмен данными только через API или события

**Преимущества:**
- Независимость развития схемы БД
- Невозможность создания tight coupling
- Масштабирование данных по сервисам

## Паттерны интеграции

### Synchronous over gRPC
**Принцип:** Синхронная коммуникация через gRPC для запрос-ответ операций.

**Когда использовать:**
- Получение данных для отображения пользователю
- Валидация данных перед операцией
- Критичные операции требующие немедленного ответа

**Реализация:**
- Protocol Buffers для сериализации
- HTTP/2 для транспорта
- Circuit breakers для надёжности
- Timeout'ы на все вызовы

### Event-driven Architecture
**Принцип:** Асинхронная коммуникация через события для слабо связанных операций.

**Когда использовать:**
- Уведомления о изменениях
- Обновление read-models
- Fire-and-forget операции
- Integration между bounded contexts

**Реализация:**
- Apache Kafka как event store
- Партиционирование по ключам
- Idempotent consumers
- Event versioning

### Temporal Workflows (Saga Orchestration)
**Принцип:** Координация сложных бизнес-процессов через Temporal workflows.

**Когда использовать:**
- Многошаговые операции
- Операции требующие компенсации
- Долгоживущие процессы
- Cross-service транзакции

**Реализация:**
- Workflow-as-Code на Go
- Automatic retries и error handling
- Durable execution
- Compensation patterns

### Choreography для простых событий
**Принцип:** Децентрализованная координация через события без центрального координатора.

**Описание паттерна:**
В отличие от оркестрации (где Temporal координирует действия), в хореографии каждый сервис знает, на какие события реагировать и какие события публиковать. Нет центрального "дирижёра" - сервисы танцуют вместе, следуя общим правилам.

**Когда использовать:**
- **Typing indicators**: пользователь печатает → все участники чата видят индикатор
- **Status updates**: пользователь онлайн/оффлайн → обновление во всех чатах  
- **Message notifications**: новое сообщение → push уведомления
- **Analytics events**: любое действие → сбор метрик
- **Audit logging**: изменения данных → запись в audit log

**Примеры хореографии в системе:**
```
1. User sends message:
   message-service → publishes → `evt.message.posted` 
   ↓
   gateway → receives → sends via WebSocket to chat members
   analytics-service → receives → updates metrics
   notification-service → receives → sends push notifications
   audit-service → receives → logs action

2. User comes online:
   user-service → publishes → `evt.user.online`
   ↓
   gateway → receives → updates presence in all user's chats
   chat-service → receives → updates last_seen timestamps
   analytics-service → receives → tracks user activity

3. User starts typing:
   frontend → sends typing → gateway
   ↓
   gateway → publishes → `evt.chat.typing_start`
   ↓
   gateway → receives → fanout to other chat members via WebSocket
```

**Реализация:**
- **Kafka topics**: `realtime.*`, `analytics.*`, `audit.*`
- **Multiple consumers**: каждый заинтересованный сервис подписывается
- **Fire-and-forget**: нет ожидания ответа или подтверждения
- **Idempotent consumers**: обработка дублей безопасна
- **No business logic coordination**: только простые реакции на события

**Преимущества хореографии:**
- **Loose coupling**: сервисы не знают друг о друге напрямую
- **Scalability**: легко добавить новых потребителей событий
- **Resilience**: падение одного consumer не влияет на других
- **Simplicity**: нет сложной логики координации

**Недостатки хореографии:**
- **No transaction guarantees**: нет атомарности между сервисами
- **Difficult debugging**: сложно отследить полный flow
- **Event proliferation**: может быть много мелких событий
- **No compensation**: если что-то пошло не так, нет автоматического rollback

**Хореография vs Оркестрация:**

| Аспект | Хореография | Оркестрация (Temporal) |
|--------|-------------|----------------------|
| **Координация** | Децентрализованная | Центральная |
| **Сложность** | Простые flows | Сложные business процессы |
| **Transactionality** | Нет | Есть (saga pattern) |
| **Compensation** | Ручная | Автоматическая |
| **Debugging** | Сложно | Легко (Temporal UI) |
| **Use cases** | Notifications, analytics | Critical business operations |

**Когда НЕ использовать хореографию:**
- Критичные бизнес-операции (создание чата + добавление участников)
- Операции требующие компенсации при ошибках
- Сложные multi-step workflows
- Операции требующие строгого порядка выполнения

## Паттерны надёжности

### Outbox Pattern
**Принцип:** Атомарность между изменением данных и публикацией события.

**Реализация:**
- События сохраняются в outbox таблицу в той же транзакции
- Отдельный worker читает outbox и публикует в Kafka
- Идемпотентность через unique constraints
- At-least-once delivery гарантии

### Idempotency
**Принцип:** Повторные операции должны быть безопасными.

**Реализация:**
- `saga_id` для workflow operations
- `client_message_id` для user actions
- Unique constraints в БД
- Idempotent consumers в Kafka

### Circuit Breaker
**Принцип:** Защита от каскадных сбоев при недоступности downstream сервисов.

**Реализация:**
- Hystrix-like patterns в Go
- Fallback strategies
- Health checks
- Graceful degradation

## Паттерны согласованности данных

### Eventual Consistency
**Принцип:** Данные между сервисами могут быть временно несогласованными.

**Реализация:**
- Event-driven updates
- Reconciliation processes
- Conflict resolution strategies
- Monitoring для detection

### CQRS (Command Query Responsibility Segregation)
**Принцип:** Разделение write и read models для оптимизации.

**Реализация:**
- Commands через workflows/events
- Read models из событий
- Separate databases для reads (optional)
- Event sourcing для audit trail

### Saga Pattern через Temporal
**Принцип:** Distributed transactions через compensating actions.

**Реализация:**
- Temporal workflows как saga orchestrator
- Compensation activities для rollback
- Timeout handling
- Retry policies

## Принципы безопасности

### Zero Trust Architecture
**Принцип:** Не доверяй, проверяй каждый запрос.

**Реализация:**
- JWT validation на каждом hop'е
- Service-to-service authentication
- Network segmentation
- Principle of least privilege

### Defense in Depth
**Принцип:** Множественные уровни защиты.

**Реализация:**
- Gateway level security (rate limiting, CORS)
- Application level (JWT, RBAC)
- Infrastructure level (TLS, VPC)
- Data level (encryption at rest)

## Принципы наблюдаемости

### Distributed Tracing
**Принцип:** Сквозная трассировка запросов через все сервисы.

**Реализация:**
- OpenTelemetry instrumentation
- Trace context propagation
- Jaeger для анализа
- Performance bottleneck detection

### Structured Logging
**Принцип:** Логи как структурированные данные.

**Реализация:**
- JSON format
- Correlation IDs
- Log levels
- Centralized aggregation

### Metrics and Alerting
**Принцип:** Проактивный мониторинг системы.

**Реализация:**
- Business metrics (messages/sec, users online)
- Technical metrics (latency, errors, throughput)
- SLI/SLO definition
- Automated alerting

## Принципы разработки

### Clean Architecture внутри сервисов
**Принцип:** Слоение с чёткими зависимостями.

**Структура:**
```
controller → service → repository
     ↓         ↓         ↓
   HTTP    Business   Database
  Layer     Logic      Layer
```

### Dependency Injection
**Принцип:** Инверсия зависимостей для тестируемости.

**Реализация:**
- Interfaces для abstractions
- DI containers (wire, dig)
- Mock generation для tests
- Configuration injection

### Test Pyramid
**Принцип:** Больше unit tests, меньше integration tests.

**Реализация:**
- Unit tests для business logic
- Integration tests для databases
- Contract tests для API
- End-to-end tests для критичных flows
