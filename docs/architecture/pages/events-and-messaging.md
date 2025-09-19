# События и система сообщений

## Kafka топики и события

### Temporal Commands (Workflows → Services)
Команды от Temporal workflows к микросервисам для выполнения activities.

#### Chat Commands
- **Topic**: `temporal.cmd.chat`
- **Key**: `saga_id`
- **События**:
  ```json
  {
    "command": "create_chat",
    "saga_id": "saga_123",
    "data": {
      "name": "Project Discussion",
      "type": "group",
      "created_by": "user_456"
    }
  }
  
  {
    "command": "add_members",
    "saga_id": "saga_123", 
    "data": {
      "chat_id": "chat_789",
      "user_ids": ["user_101", "user_102"]
    }
  }
  
  {
    "command": "remove_members",
    "saga_id": "saga_123",
    "data": {
      "chat_id": "chat_789", 
      "user_ids": ["user_103"]
    }
  }
  ```

#### Message Commands  
- **Topic**: `temporal.cmd.message`
- **Key**: `saga_id`
- **События**:
  ```json
  {
    "command": "post_message",
    "saga_id": "saga_124",
    "data": {
      "chat_id": "chat_789",
      "user_id": "user_456",
      "content": "Hello everyone!",
      "client_message_id": "client_msg_001"
    }
  }
  
  {
    "command": "delete_message", 
    "saga_id": "saga_125",
    "data": {
      "message_id": "msg_999",
      "reason": "compensation"
    }
  }
  ```

#### User Commands
- **Topic**: `temporal.cmd.user` 
- **Key**: `saga_id`
- **События**:
  ```json
  {
    "command": "create_profile",
    "saga_id": "saga_126",
    "data": {
      "user_id": "user_789",
      "username": "john_doe", 
      "email": "john@example.com"
    }
  }
  
  {
    "command": "update_status",
    "saga_id": "saga_127", 
    "data": {
      "user_id": "user_789",
      "status": "online"
    }
  }
  ```

### Domain Events (Services → Temporal & Other Services)
События от микросервисов обратно в Temporal workflows и другие подписчики.

#### Chat Events
- **Topic**: `domain.evt.chat`
- **Key**: `saga_id` (для Temporal) или `chat_id` (для choreography)
- **События**:
  ```json
  {
    "event": "chat_created",
    "saga_id": "saga_123",
    "timestamp": "2024-01-15T10:30:00Z",
    "data": {
      "chat_id": "chat_789",
      "name": "Project Discussion",
      "created_by": "user_456"
    }
  }
  
  {
    "event": "members_added",
    "saga_id": "saga_123",
    "timestamp": "2024-01-15T10:31:00Z", 
    "data": {
      "chat_id": "chat_789",
      "added_users": ["user_101", "user_102"],
      "added_by": "user_456"
    }
  }
  
  {
    "event": "member_removed",
    "chat_id": "chat_789",
    "timestamp": "2024-01-15T10:32:00Z",
    "data": {
      "chat_id": "chat_789", 
      "removed_user": "user_103",
      "removed_by": "user_456"
    }
  }
  ```

#### Message Events
- **Topic**: `domain.evt.message`
- **Key**: `saga_id` или `chat_id`
- **События**:
  ```json
  {
    "event": "message_posted",
    "saga_id": "saga_124",
    "chat_id": "chat_789",
    "timestamp": "2024-01-15T10:33:00Z",
    "data": {
      "message_id": "msg_001",
      "chat_id": "chat_789",
      "user_id": "user_456", 
      "content": "Hello everyone!",
      "client_message_id": "client_msg_001"
    }
  }
  
  {
    "event": "message_delivered",
    "chat_id": "chat_789",
    "timestamp": "2024-01-15T10:33:01Z",
    "data": {
      "message_id": "msg_001",
      "delivered_to": ["user_101", "user_102"]
    }
  }
  ```

#### User Events
- **Topic**: `domain.evt.user`
- **Key**: `user_id`
- **События**:
  ```json
  {
    "event": "user_registered",
    "saga_id": "saga_126", 
    "timestamp": "2024-01-15T10:34:00Z",
    "data": {
      "user_id": "user_789",
      "username": "john_doe",
      "email": "john@example.com"
    }
  }
  
  {
    "event": "user_status_changed",
    "user_id": "user_789",
    "timestamp": "2024-01-15T10:35:00Z",
    "data": {
      "user_id": "user_789",
      "old_status": "offline",
      "new_status": "online"
    }
  }
  ```

### Real-time Events (Choreography)
Fire-and-forget события для real-time функциональности. В хореографии нет центрального координатора - каждый сервис реагирует на события самостоятельно.

#### Принципы хореографии в системе:
- **Decoupled**: сервисы не знают друг о друге
- **Reactive**: реакция на события, а не активные вызовы  
- **Resilient**: падение одного consumer не ломает других
- **Scalable**: легко добавлять новых потребителей

#### Typing Events
- **Topic**: `realtime.typing`
- **Key**: `chat_id`
- **События**:
  ```json
  {
    "event": "typing_start",
    "chat_id": "chat_789",
    "user_id": "user_456",
    "timestamp": "2024-01-15T10:36:00Z"
  }
  
  {
    "event": "typing_stop", 
    "chat_id": "chat_789",
    "user_id": "user_456",
    "timestamp": "2024-01-15T10:36:05Z"
  }
  ```

#### Presence Events
- **Topic**: `realtime.presence`
- **Key**: `user_id`
- **События**:
  ```json
  {
    "event": "user_online",
    "user_id": "user_456", 
    "timestamp": "2024-01-15T10:37:00Z",
    "data": {
      "connection_id": "conn_123",
      "device": "web"
    }
  }
  
  {
    "event": "user_offline",
    "user_id": "user_456",
    "timestamp": "2024-01-15T10:45:00Z",
    "data": {
      "connection_id": "conn_123",
      "last_seen": "2024-01-15T10:45:00Z"
    }
  }
  ```

#### Analytics Events (Choreography)
- **Topic**: `analytics.events`
- **Key**: `user_id` или `chat_id`
- **Consumers**: analytics-service, business-intelligence
- **События**:
  ```json
  {
    "event": "message_sent",
    "user_id": "user_456",
    "chat_id": "chat_789", 
    "timestamp": "2024-01-15T10:38:00Z",
    "data": {
      "message_length": 25,
      "message_type": "text",
      "device": "web"
    }
  }
  
  {
    "event": "chat_opened",
    "user_id": "user_456",
    "chat_id": "chat_789",
    "timestamp": "2024-01-15T10:39:00Z",
    "data": {
      "source": "notification",
      "session_id": "sess_123"
    }
  }
  ```

#### Audit Events (Choreography)
- **Topic**: `audit.events`
- **Key**: `entity_id` (chat_id, user_id, message_id)
- **Consumers**: audit-service, compliance-service
- **События**:
  ```json
  {
    "event": "data_modified",
    "entity_type": "chat",
    "entity_id": "chat_789",
    "modified_by": "user_456",
    "timestamp": "2024-01-15T10:40:00Z",
    "data": {
      "operation": "member_added",
      "old_value": null,
      "new_value": "user_102",
      "ip_address": "192.168.1.100"
    }
  }
  
  {
    "event": "data_accessed",
    "entity_type": "message",
    "entity_id": "msg_001", 
    "accessed_by": "user_456",
    "timestamp": "2024-01-15T10:41:00Z",
    "data": {
      "operation": "read",
      "chat_id": "chat_789"
    }
  }
  ```

#### Notification Events (Choreography)
- **Topic**: `notifications.events`
- **Key**: `user_id`
- **Consumers**: notification-service, email-service, push-service
- **События**:
  ```json
  {
    "event": "notification_triggered",
    "user_id": "user_102",
    "timestamp": "2024-01-15T10:42:00Z",
    "data": {
      "type": "new_message",
      "chat_id": "chat_789",
      "message_id": "msg_001",
      "sender": "user_456",
      "urgency": "normal"
    }
  }
  
  {
    "event": "mention_received",
    "user_id": "user_102", 
    "timestamp": "2024-01-15T10:43:00Z",
    "data": {
      "chat_id": "chat_789",
      "message_id": "msg_002",
      "mentioned_by": "user_456",
      "urgency": "high"
    }
  }
  ```

#### System Events (Choreography)
- **Topic**: `system.events`
- **Key**: `service_name`
- **Consumers**: monitoring-service, alerting-service
- **События**:
  ```json
  {
    "event": "service_health_changed",
    "service_name": "message-service",
    "timestamp": "2024-01-15T10:44:00Z",
    "data": {
      "old_status": "healthy",
      "new_status": "degraded",
      "reason": "high_latency",
      "metric_value": 2500
    }
  }
  
  {
    "event": "rate_limit_exceeded",
    "user_id": "user_456",
    "timestamp": "2024-01-15T10:45:00Z", 
    "data": {
      "endpoint": "/api/messages",
      "current_rate": 1500,
      "limit": 1000,
      "window": "1m"
    }
  }
  ```

#### Примеры хореографических flow:

**Flow 1: Пользователь отправляет сообщение**
```
1. message-service → `domain.evt.message.posted`
   ↓
2. gateway (consumer) → WebSocket fanout to chat members
3. analytics-service (consumer) → update metrics  
4. notification-service (consumer) → push notifications
5. audit-service (consumer) → log message creation
```

**Flow 2: Пользователь заходит онлайн**
```
1. gateway → `realtime.presence.user_online`
   ↓
2. user-service (consumer) → update user status in DB
3. analytics-service (consumer) → track user activity
4. chat-service (consumer) → update last_seen in all user's chats
5. gateway (consumer) → notify friends about status change
```

**Flow 3: Превышение rate limit**
```
1. gateway → `system.events.rate_limit_exceeded`
   ↓
2. monitoring-service (consumer) → update dashboards
3. alerting-service (consumer) → send alert to ops team
4. analytics-service (consumer) → track abuse patterns  
5. user-service (consumer) → temporarily flag user (if severe)
```

## Temporal Workflows

### CreateChatWithMembersWorkflow
**Назначение:** Создание чата с добавлением участников и отправкой welcome сообщения.

**Активности:**
1. `CreateChatActivity` (chat-service)
2. `AddMembersActivity` (chat-service)  
3. `PostWelcomeMessageActivity` (message-service)
4. `NotifyMembersActivity` (gateway)

**Компенсации:**
- `DeleteChatActivity`
- `RemoveMembersActivity`
- `DeleteMessageActivity`

### UserRegistrationWorkflow  
**Назначение:** Полная регистрация пользователя с созданием профиля.

**Активности:**
1. `CreateAuthRecordActivity` (auth-service)
2. `CreateUserProfileActivity` (user-service)
3. `SendWelcomeEmailActivity` (notification-service)

**Компенсации:**
- `DeleteAuthRecordActivity`
- `DeleteUserProfileActivity`

### ChatMemberManagementWorkflow
**Назначение:** Управление участниками чата с уведомлениями.

**Активности:**
1. `ValidateChatPermissionsActivity` (chat-service)
2. `AddRemoveMembersActivity` (chat-service)
3. `PostSystemMessageActivity` (message-service) 
4. `NotifyAffectedUsersActivity` (gateway)

## Партиционирование и Consumer Groups

### Партиционирование стратегии
- **saga_id**: для Temporal workflows (порядок выполнения)
- **chat_id**: для сообщений в чате (порядок доставки)
- **user_id**: для пользовательских событий

### Consumer Groups
- **temporal-workers**: обработка команд от workflows
- **gateway-fanout**: отправка WebSocket событий
- **analytics**: сбор метрик и аналитики
- **audit**: логирование для аудита

## Форматы сообщений

### Envelope Format
Общий формат для всех событий:
```json
{
  "id": "evt_12345",
  "type": "domain.evt.message.posted",
  "source": "message-service",
  "specversion": "1.0",
  "time": "2024-01-15T10:33:00Z",
  "subject": "chat_789",
  "datacontenttype": "application/json",
  "data": {
    // event-specific payload
  },
  "correlationid": "saga_124",
  "traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
}
```

### Headers
Kafka message headers:
- `eventType`: тип события
- `correlationId`: для связи с workflow/saga
- `traceparent`: для distributed tracing  
- `source`: источник события
- `version`: версия схемы события

## Обеспечение надёжности

### At-least-once Delivery
- Kafka producer acknowledgments: `acks=all`
- Consumer manual commits после обработки
- Outbox pattern для атомарности

### Идемпотентность
- Unique constraints по `saga_id` + `event_type`
- Upsert operations в consumers
- Idempotency keys для external APIs

### Dead Letter Queues
- Отдельные топики для failed messages
- Retry logic с exponential backoff
- Manual intervention для stuck messages

### Schema Evolution
- Backward compatible changes
- Schema registry для validation
- Versioned event types
