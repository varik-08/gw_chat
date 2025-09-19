# Микросервисы и их границы

## Обзор декомпозиции монолита

Наш монолитный чат-сервис разбивается на следующие микросервисы по принципу **Domain-Driven Design** и **bounded contexts**:

## Сервисы

### gateway (API Gateway)
**Ответственность:**
- Единая точка входа для HTTP/REST запросов и WebSocket соединений
- Валидация входящих запросов и JWT токенов  
- Маршрутизация запросов в соответствующие микросервисы
- Проксирование gRPC/HTTP вызовов
- Fanout для WebSocket соединений (подписчик Kafka событий)
- Rate limiting и CORS

**Взаимодействие:**
- Внешний: HTTP/REST + WebSocket клиенты
- Внутренний: gRPC вызовы к другим сервисам, Kafka consumer

### auth-service (Аутентификация)
**Ответственность:**
- Регистрация новых пользователей
- Аутентификация (логин/logout)
- Управление JWT токенами (access/refresh)
- JWKS endpoint (`/.well-known/jwks.json`)
- Валидация учётных данных

**База данных:** `auth_schema`
- `users_auth` (id, email, password_hash, created_at, updated_at)
- `refresh_tokens` (token_id, user_id, token_hash, expires_at)

**Взаимодействие:**
- gRPC сервер для других сервисов
- HTTP endpoint для JWKS
- События: `user.registered`, `user.authenticated`

### user-service (Управление пользователями)
**Ответственность:**
- CRUD операции с профилями пользователей
- Управление пользовательскими настройками
- Поиск пользователей
- Статусы пользователей (онлайн/оффлайн)

**База данных:** `user_schema`
- `users` (id, username, email, display_name, avatar_url, status, created_at, updated_at)
- `user_settings` (user_id, settings_json)

**Взаимодействие:**
- gRPC сервер
- Kafka consumer для команд Temporal workflow
- События: `user.updated`, `user.status_changed`

### chat-service (Управление чатами)
**Ответственность:**
- Создание чатов (приватные/публичные)
- Управление участниками чатов
- Настройки чатов (название, описание, аватар)
- Права доступа и роли в чатах

**База данных:** `chat_schema`
- `chats` (id, name, description, type, created_by, created_at, updated_at)
- `chat_members` (chat_id, user_id, role, joined_at)
- `chat_settings` (chat_id, settings_json)

**Взаимодействие:**
- gRPC сервер
- Kafka producer/consumer (outbox pattern)
- События: `chat.created`, `chat.member_added`, `chat.member_removed`

### message-service (Управление сообщениями)
**Ответственность:**
- Создание и сохранение сообщений
- Получение истории сообщений
- Обеспечение идемпотентности через `client_message_id`
- Outbox pattern для надёжной публикации событий

**База данных:** `message_schema`
- `messages` (id, chat_id, user_id, content, client_message_id, created_at, updated_at)
- `message_outbox` (id, event_type, payload, published_at, saga_id)

**Взаимодействие:**
- gRPC сервер
- Kafka producer (outbox worker)
- События: `message.posted`, `message.delivered`

### temporal-server (Workflow оркестрация)
**Ответственность:**
- Координация сложных бизнес-процессов
- Обеспечение consistency через Temporal workflows
- Компенсации при ошибках (Saga pattern)
- Долгоживущие процессы

**Workflows:**
- `CreateChatWithMembersWorkflow`: создание чата + добавление участников + начальное сообщение
- `UserRegistrationWorkflow`: регистрация + создание профиля + отправка welcome сообщения
- `ChatMemberManagementWorkflow`: добавление/удаление участников с уведомлениями

**Взаимодействие:**
- Temporal server с Workers в каждом сервисе
- Kafka для событий choreography
- Вызовы activities в других сервисах

## Принципы декомпозиции

### Database-per-Service
- Каждый сервис владеет своими данными
- Отдельная схема PostgreSQL для каждого сервиса
- Никаких shared databases или прямых обращений к чужим БД

### Bounded Context
- **Auth Context**: аутентификация и авторизация
- **User Context**: профили и настройки пользователей  
- **Chat Context**: чаты и участники
- **Message Context**: сообщения и их доставка
- **Gateway Context**: маршрутизация и протоколы

### Communication Patterns
- **Sync**: gRPC для запрос-ответ операций
- **Async**: Kafka для событий и команд
- **Orchestration**: Temporal для сложных workflows
- **Choreography**: Kafka для fire-and-forget событий (typing, статусы)

### Data Consistency
- **Strong consistency**: внутри сервиса (ACID)
- **Eventual consistency**: между сервисами
- **Saga pattern**: через Temporal workflows для критичных операций
- **Outbox pattern**: для атомарности БД + события
