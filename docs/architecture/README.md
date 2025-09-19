# Архитектурная документация

Документация по архитектуре микросервисного чат-приложения с использованием Go, Temporal, Kafka и PostgreSQL.

## 📊 [Диаграммы архитектуры](./diagrams/)

## 🏗️ Ключевые архитектурные принципы

### Service per Bounded Context
- **Gateway**: HTTP/WS + маршрутизация + JWT
- **Auth Service**: аутентификация + JWT + JWKS
- **User Service**: профили + настройки + статусы
- **Chat Service**: чаты + участники + права доступа
- **Message Service**: сообщения + история + доставка
- **Temporal**: workflow оркестрация + саги

### Database per Service
- **auth_schema**: пользователи + токены
- **user_schema**: профили + настройки
- **chat_schema**: чаты + участники
- **message_schema**: сообщения + outbox
- **temporal_schema**: workflows + executions

### Event-driven Architecture
- **Synchronous**: gRPC для запрос-ответ
- **Asynchronous**: Kafka для событий и команд
- **Orchestration**: Temporal для сложных workflows
- **Choreography**: Kafka для fire-and-forget событий

## 🔄 Паттерны интеграции

### Temporal Workflows (Оркестрация)
- **CreateChatWithMembersWorkflow**: создание чата + участники + welcome сообщение
- **UserRegistrationWorkflow**: регистрация + профиль + приветствие
- **ChatMemberManagementWorkflow**: управление участниками + уведомления

### Kafka Choreography (Хореография)
- **Real-time события**: typing, presence, notifications
- **Analytics события**: метрики, KPI, пользовательская активность
- **Audit события**: логирование, compliance, аудит
- **System события**: мониторинг, алерты, health checks

## 🛠️ Технологический стек

### Backend
- **Go 1.24+**: основной язык программирования
- **gRPC**: межсервисная коммуникация
- **Protocol Buffers**: сериализация данных
- **PostgreSQL**: основная база данных
- **Redis**: кэширование и сессии

### Event Streaming
- **Apache Kafka**: event streaming платформа
- **KRaft mode**: Kafka без ZooKeeper
- **Redpanda**: альтернативная реализация Kafka API

### Workflow Engine
- **Temporal**: workflow оркестрация
- **Workflow-as-Code**: бизнес-логика в Go
- **Saga Pattern**: компенсации и откаты

### Observability
- **OpenTelemetry**: distributed tracing
- **Jaeger**: анализ производительности
- **Prometheus + Grafana**: метрики и дашборды
- **ELK Stack**: централизованные логи

### Security
- **JWT RS256**: аутентификация
- **JWKS**: публикация ключей
- **TLS**: шифрование соединений
- **Zero Trust**: проверка каждого запроса

## 🚀 Развёртывание

### Локальная разработка
```bash
# Запуск всей системы
docker-compose -f docker-compose.micro.yml up

# Сервисы:
# - kafka (KRaft mode)
# - kafka-ui
# - postgres
# - temporal-server
# - temporal-ui
# - jaeger
# - gateway
# - auth-service
# - user-service
# - chat-service
# - message-service
```

### Production
- **Kubernetes**: оркестрация контейнеров
- **Service Mesh**: Istio для управления трафиком
- **Multi-region**: активный-активный режим
- **Auto-scaling**: на основе метрик

## 📈 Мониторинг и метрики

### Business Metrics
- **Active users**: количество активных пользователей
- **Messages per second**: пропускная способность
- **Chat creation rate**: создание чатов
- **User engagement**: вовлечённость пользователей

### Technical Metrics
- **Response time**: время ответа API
- **Error rate**: процент ошибок
- **Throughput**: RPS по сервисам
- **Resource utilization**: CPU, memory, disk

### Alerts
- **Critical**: < 5 минут response time
- **Warning**: < 30 минут response time
- **Escalation**: автоматическое эскалирование

## 🔧 Инструменты разработки

### Code Generation
- **protoc**: генерация gRPC клиентов/серверов
- **mockgen**: моки для тестирования
- **sqlc**: генерация type-safe SQL кода

### Quality Assurance
- **golangci-lint**: статический анализ
- **go test**: unit и integration тесты
- **go mod**: управление зависимостями

### Development Workflow
- **Git**: версионирование кода
- **CI/CD**: автоматическое тестирование и деплой
- **Code review**: обязательный review всех изменений
- **Feature flags**: постепенное включение фичей

## 📚 Дополнительные ресурсы

- [Temporal Documentation](https://docs.temporal.io/)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)