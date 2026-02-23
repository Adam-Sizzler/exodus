# v2ray-stat API Documentation

API для управления пользователями, нодами, конфигурациями и статистикой в проекте v2ray-stat.

**Base URL:** `http://127.0.0.1:9952`

---

## Содержание

1. [Аутентификация](#аутентификация)
2. [Nodes API](#nodes-api)
3. [Users API](#users-api)
4. [Config Profiles API](#config-profiles-api)
5. [Статистика](#статистика)
6. [Управление пользователями](#управление-пользователями)
7. [Сброс статистики](#сброс-статистики)

---

## Аутентификация

Большинство endpoints требуют API токен в заголовке:

```bash
Authorization: Bearer YOUR_API_TOKEN
```

---

## Nodes API

### Получить все ноды

**GET** `/api/v1/nodes`

Возвращает список всех нод с их данными.

**Пример запроса:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9952/api/v1/nodes
```

**Пример ответа:**
```json
{
  "nodes": [
    {
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "id": 1,
      "name": "Node Moscow",
      "address": "192.168.1.100",
      "port": 443,
      "is_connected": true,
      "is_connecting": false,
      "is_disabled": false,
      "last_status_change": "2025-02-22T10:00:00Z",
      "last_status_message": "Connected",
      "xray_version": "1.8.16",
      "node_version": "1.0.0",
      "xray_uptime": "7d 12h 30m",
      "users_online": 15,
      "consumption_multiplier": 1000000000,
      "is_traffic_tracking_active": true,
      "traffic_reset_day": 1,
      "traffic_limit_bytes": 107374182400,
      "traffic_used_bytes": 53687091200,
      "notify_percent": 80,
      "provider_uuid": "660e8400-e29b-41d4-a716-446655440001",
      "view_position": 1,
      "country_code": "RU",
      "tags": ["production", "moscow"],
      "cpu_count": 4,
      "cpu_model": "Intel Xeon",
      "total_ram": "8GB",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-02-22T10:30:00Z"
    }
  ],
  "count": 1
}
```

---

### Получить ноду по UUID

**GET** `/api/v1/nodes/{uuid}`

Возвращает данные одной ноды по UUID.

**Пример запроса:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9952/api/v1/nodes/550e8400-e29b-41d4-a716-446655440000
```

**Пример ответа:**
```json
{
  "node": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "id": 1,
    "name": "Node Moscow",
    "address": "192.168.1.100",
    "port": 443,
    "is_connected": true,
    "is_disabled": false,
    "country_code": "RU",
    "tags": ["production", "moscow"],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-02-22T10:30:00Z"
  }
}
```

---

### Обновить ноду (частичное обновление)

**PATCH** `/api/v1/nodes/{uuid}`

Обновляет только указанные поля ноды.

**Пример запроса:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Node Moscow Updated",
    "is_disabled": false,
    "tags": ["production", "moscow", "updated"],
    "country_code": "RU"
  }' \
  http://localhost:9952/api/v1/nodes/550e8400-e29b-41d4-a716-446655440000
```

**Пример ответа:**
```json
{
  "node": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Node Moscow Updated",
    "is_disabled": false,
    "tags": ["production", "moscow", "updated"],
    ...
  },
  "message": "node updated successfully"
}
```

**Доступные поля для обновления:**
| Поле | Тип | Описание |
|------|-----|----------|
| `name` | string | Имя ноды |
| `address` | string | IP адрес или домен |
| `port` | int | Порт (1-65535) |
| `is_disabled` | boolean | Статус блокировки |
| `consumption_multiplier` | int64 | Множитель потребления |
| `is_traffic_tracking_active` | boolean | Отслеживание трафика |
| `traffic_reset_day` | int | День сброса трафика (1-31) |
| `traffic_limit_bytes` | int64 | Лимит трафика |
| `notify_percent` | int | Процент уведомления (0-100) |
| `provider_uuid` | string | UUID провайдера |
| `view_position` | int | Позиция отображения |
| `country_code` | string | Код страны (2 буквы) |
| `tags` | string[] | Теги |

**Примеры обновлений:**

1. **Обновить настройки трафика:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "is_traffic_tracking_active": true,
    "traffic_reset_day": 15,
    "traffic_limit_bytes": 107374182400,
    "notify_percent": 80
  }' \
  http://localhost:9952/api/v1/nodes/550e8400-e29b-41d4-a716-446655440000
```

2. **Сбросить поле на NULL:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{"provider_uuid": ""}' \
  http://localhost:9952/api/v1/nodes/550e8400-e29b-41d4-a716-446655440000
```

---

## Users API

### Получить всех пользователей

**GET** `/api/v1/users-list`

Возвращает список всех пользователей из таблицы `users`.

**Пример запроса:**
```bash
curl "http://localhost:9952/api/v1/users-list"
```

**Пример ответа:**
```json
{
  "users": [
    {
      "t_id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440001",
      "short_uuid": "abc123",
      "username": "john_doe",
      "status": "ACTIVE",
      "traffic_limit_bytes": 10737418240,
      "traffic_limit_strategy": "MONTH",
      "expire_at": "2026-12-31T23:59:59Z",
      "sub_last_user_agent": "Mozilla/5.0",
      "sub_last_opened_at": "2025-02-22T10:00:00Z",
      "last_traffic_reset_at": "2025-02-01T00:00:00Z",
      "trojan_password": "password123",
      "vless_uuid": "660e8400-e29b-41d4-a716-446655440002",
      "ss_password": "sspass456",
      "description": "Test user",
      "tag": "premium",
      "telegram_id": 123456789,
      "email": "john@example.com",
      "hwid_device_limit": 3,
      "external_squad_uuid": null,
      "last_triggered_threshold": 0,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-02-23T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

### Получить пользователя по UUID

**GET** `/api/v1/users-list/{uuid}`

**Пример запроса:**
```bash
curl "http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001"
```

---

### Создать пользователя

**POST** `/api/v1/users-list/create`

Создаёт нового пользователя в таблице `users`.

**Пример запроса:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "username": "new_user",
    "status": "ACTIVE",
    "traffic_limit_bytes": 10737418240,
    "traffic_limit_strategy": "MONTH",
    "expire_at": "2026-12-31T23:59:59Z",
    "description": "New test user",
    "tag": "premium",
    "email": "newuser@example.com",
    "hwid_device_limit": 3
  }' \
  http://localhost:9952/api/v1/users-list/create
```

**Пример ответа (201 Created):**
```json
{
  "user": {
    "t_id": 2,
    "uuid": "770e8400-e29b-41d4-a716-446655440003",
    "short_uuid": "def456",
    "username": "new_user",
    "status": "ACTIVE",
    "traffic_limit_bytes": 10737418240,
    "traffic_limit_strategy": "MONTH",
    "expire_at": "2026-12-31T23:59:59Z",
    "trojan_password": "a1b2c3d4e5f6g7h8",
    "vless_uuid": "880e8400-e29b-41d4-a716-446655440004",
    "ss_password": "x9y8z7w6v5u4t3s2",
    "description": "New test user",
    "tag": "premium",
    "email": "newuser@example.com",
    "hwid_device_limit": 3,
    "last_triggered_threshold": 0,
    "created_at": "2025-02-23T12:00:00Z",
    "updated_at": "2025-02-23T12:00:00Z"
  },
  "message": "user created successfully"
}
```

**Обязательные поля:**
- `username` — уникальное имя (буквы, цифры, `_`, `-`)
- `status` — один из: `ACTIVE`, `DISABLED`, `LIMITED`, `EXPIRED`
- `expire_at` — дата истечения в формате ISO 8601 (RFC3339)

**Опциональные поля:**
- `uuid` — будет сгенерирован автоматически, если не указан
- `trojan_password` — будет сгенерирован автоматически
- `vless_uuid` — будет сгенерирован автоматически
- `ss_password` — будет сгенерирован автоматически
- `traffic_limit_bytes` — по умолчанию 0
- `traffic_limit_strategy` — по умолчанию `NO_RESET`
- `description`, `tag`, `telegram_id`, `email`, `hwid_device_limit`, `external_squad_uuid`
- `last_triggered_threshold` — по умолчанию 0

**Примеры ошибок:**
```json
// 400 Bad Request - неверный username
{
  "error": "username can only contain letters, numbers, underscores, and hyphens"
}

// 400 Bad Request - неверный статус
{
  "error": "status must be one of: ACTIVE, DISABLED, LIMITED, EXPIRED"
}

// 409 Conflict - username занят
{
  "error": "username already exists"
}
```

---

### Обновить пользователя (частичное обновление)

**PATCH** `/api/v1/users-list/{uuid}`

**Пример запроса:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "status": "DISABLED",
    "traffic_limit_bytes": 21474836480,
    "expire_at": "2027-06-30T23:59:59Z"
  }' \
  http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001
```

**Пример ответа:**
```json
{
  "user": {
    "t_id": 1,
    "uuid": "550e8400-e29b-41d4-a716-446655440001",
    "username": "john_doe",
    "status": "DISABLED",
    "traffic_limit_bytes": 21474836480,
    "expire_at": "2027-06-30T23:59:59Z",
    ...
  },
  "message": "user updated successfully"
}
```

**Доступные поля для обновления:**
| Поле | Тип | Описание |
|------|-----|----------|
| `status` | string | `ACTIVE`, `DISABLED`, `LIMITED`, `EXPIRED` |
| `traffic_limit_bytes` | int64 | Лимит трафика (>= 0) |
| `traffic_limit_strategy` | string | `NO_RESET`, `DAY`, `WEEK`, `MONTH` |
| `expire_at` | string | Дата в ISO 8601 |
| `description` | string | Описание |
| `tag` | string | Тег |
| `telegram_id` | int64 | Telegram ID |
| `email` | string | Email |
| `hwid_device_limit` | int | Лимит устройств (>= 0) |
| `external_squad_uuid` | string | UUID сквада |
| `last_triggered_threshold` | int | Процент (0-100) |

**Примеры обновлений:**

1. **Обновить дату истечения и email:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "expire_at": "2026-12-31T23:59:59Z",
    "email": "updated@example.com"
  }' \
  http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001
```

2. **Обновить стратегию сброса трафика:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "traffic_limit_strategy": "MONTH",
    "hwid_device_limit": 5
  }' \
  http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001
```

3. **Сбросить поле на NULL:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{"description": ""}' \
  http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001
```

---

### Удалить пользователя

**DELETE** `/api/v1/users-list/{uuid}`

Удаляет пользователя по UUID. Связанные записи в других таблицах удаляются автоматически через CASCADE.

**Пример запроса:**
```bash
curl -X DELETE \
  http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001
```

**Пример ответа:**
```json
{
  "message": "user deleted successfully",
  "uuid": "550e8400-e29b-41d4-a716-446655440001",
  "username": "john_doe",
  "t_id": 1
}
```

**Примеры ошибок:**
```json
// 404 Not Found - пользователь не найден
{
  "error": "user not found"
}

// 500 Internal Server Error - ошибка при удалении
{
  "error": "failed to delete user"
}
```

---

## Config Profiles API

API для управления конфигурациями sing-box в формате JSON.

### Автоматическая синхронизация Inbounds

При создании или обновлении конфигурации через API, система **автоматически извлекает все inbounds** из JSON конфигурации и сохраняет их в таблицу `config_profile_inbounds`.

**Как это работает:**
1. Вы отправляете JSON конфигурации с массивом `inbounds`
2. Система парсит каждый inbound и извлекает:
   - `tag` — уникальный идентификатор inbound в рамках профиля (обязательно)
   - `type` или `protocol` — тип протокола
   - `network` — сеть (tcp, udp, и т.д.)
   - `security` — безопасность (tls, reality, и т.д.)
   - `port` или `listen_port` — порт
   - Полный JSON inbound сохраняется в `raw_inbound`
3. Все inbounds с одинаковым `tag` в рамках одного профиля автоматически заменяются
4. При удалении конфигурации все связанные inbounds удаляются автоматически через CASCADE

**Пример JSON конфигурации:**
```json
{
  "log": {"level": "info"},
  "inbounds": [
    {
      "tag": "Shadowsocks",
      "type": "shadowsocks",
      "listen": "0.0.0.0",
      "listen_port": 1234,
      "network": "tcp,udp"
    },
    {
      "tag": "VLESS-Reality",
      "type": "vless",
      "listen": "0.0.0.0",
      "listen_port": 1235,
      "network": "tcp",
      "security": "tls"
    }
  ],
  "outbounds": [{"type": "direct", "tag": "direct"}]
}
```

После создания конфигурации с таким JSON, в таблице `config_profile_inbounds` появятся две записи с соответствующими данными.

---

### Получить все конфигурации

**GET** `/api/v1/config-profiles`

Возвращает список всех конфигураций с их JSON данными.

**Пример запроса:**
```bash
curl "http://localhost:9952/api/v1/config-profiles"
```

**Пример ответа:**
```json
{
  "profiles": [
    {
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "view_position": 0,
      "name": "default-profile",
      "config": "{\"log\":{\"level\":\"info\"},\"dns\":{\"servers\":[\"8.8.8.8\"]},\"inbounds\":[{\"type\":\"mixed\",\"tag\":\"mixed-in\",\"listen\":\"0.0.0.0\",\"listen_port\":2080}],\"outbounds\":[{\"type\":\"direct\",\"tag\":\"direct\"}],\"route\":{\"rules\":[]}}",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-02-23T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

### Получить конфигурацию по UUID

**GET** `/api/v1/config-profiles/{uuid}`

Возвращает данные одной конфигурации по UUID.

**Пример запроса:**
```bash
curl "http://localhost:9952/api/v1/config-profiles/550e8400-e29b-41d4-a716-446655440000"
```

**Пример ответа:**
```json
{
  "profile": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "view_position": 0,
    "name": "default-profile",
    "config": "{\"log\":{\"level\":\"info\"},\"dns\":{\"servers\":[\"8.8.8.8\"]}}",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-02-23T10:00:00Z"
  }
}
```

---

### Создать конфигурацию

**POST** `/api/v1/config-profiles`

Создаёт новую конфигурацию sing-box в таблице `config_profiles`.

**Пример запроса:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "default-profile",
    "view_position": 0,
    "config": {
      "log": {
        "level": "info",
        "timestamp": true
      },
      "dns": {
        "servers": [
          {"tag": "google", "address": "8.8.8.8"},
          {"tag": "local", "address": "223.5.5.5", "detour": "direct"}
        ]
      },
      "inbounds": [
        {
          "type": "mixed",
          "tag": "mixed-in",
          "listen": "0.0.0.0",
          "listen_port": 2080
        }
      ],
      "outbounds": [
        {"type": "direct", "tag": "direct"}
      ],
      "route": {
        "rules": [],
        "auto_detect_interface": true
      }
    }
  }' \
  http://localhost:9952/api/v1/config-profiles
```

**Пример ответа (201 Created):**
```json
{
  "message": "config profile created",
  "uuid": "770e8400-e29b-41d4-a716-446655440003"
}
```

**Обязательные поля:**
- `name` — уникальное имя конфигурации
- `config` — валидный JSON объект конфигурации sing-box

**Опциональные поля:**
- `view_position` — позиция отображения (по умолчанию 0)

**Примеры ошибок:**
```json
// 400 Bad Request - неверный JSON
{
  "error": "config must be valid JSON"
}

// 400 Bad Request - отсутствует имя
{
  "error": "name is required"
}

// 400 Bad Request - пустая конфигурация
{
  "error": "config is required"
}

// 409 Conflict - имя занято
{
  "error": "name already exists"
}
```

---

### Обновить конфигурацию (частичное обновление)

**PATCH** `/api/v1/config-profiles/{uuid}`

Обновляет только указанные поля конфигурации.

**Пример запроса:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "name": "updated-profile",
    "view_position": 5
  }' \
  http://localhost:9952/api/v1/config-profiles/550e8400-e29b-41d4-a716-446655440000
```

**Пример ответа:**
```json
{
  "profile": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "view_position": 5,
    "name": "updated-profile",
    "config": "{\"log\":{\"level\":\"info\"},\"dns\":{\"servers\":[\"8.8.8.8\"]}}",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-02-23T12:00:00Z"
  }
}
```

**Доступные поля для обновления:**
| Поле | Тип | Описание |
|------|-----|----------|
| `name` | string | Имя конфигурации (уникальное) |
| `view_position` | int | Позиция отображения |
| `config` | object | sing-box конфигурация JSON |

**Примеры обновлений:**

1. **Обновить JSON конфигурацию:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "log": {"level": "debug"},
      "dns": {"servers": ["1.1.1.1", "8.8.8.8"]},
      "inbounds": [{"type": "mixed", "tag": "mixed-in", "listen": "0.0.0.0", "listen_port": 2080}],
      "outbounds": [{"type": "direct", "tag": "direct"}],
      "route": {"rules": []}
    }
  }' \
  http://localhost:9952/api/v1/config-profiles/550e8400-e29b-41d4-a716-446655440000
```

2. **Обновить позицию отображения:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{"view_position": 10}' \
  http://localhost:9952/api/v1/config-profiles/550e8400-e29b-41d4-a716-446655440000
```

---

### Удалить конфигурацию

**DELETE** `/api/v1/config-profiles/{uuid}`

Удаляет конфигурацию по UUID.

**Пример запроса:**
```bash
curl -X DELETE \
  http://localhost:9952/api/v1/config-profiles/550e8400-e29b-41d4-a716-446655440000
```

**Пример ответа:**
```json
{
  "message": "config profile deleted",
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "name": "default-profile"
}
```

---

## Internal Squads API

API для управления внутренними сквадами (internal squads).

### Получить все сквады

**GET** `/api/v1/internal-squads`

Возвращает список всех внутренних сквадов.

**Пример запроса:**
```bash
curl "http://localhost:9952/api/v1/internal-squads"
```

**Пример ответа:**
```json
{
  "squads": [
    {
      "uuid": "2ef69520-9dc4-44f7-a492-c709fc34c41f",
      "view_position": 1,
      "name": "Default-Squad",
      "created_at": "2026-02-23T14:11:28.074Z",
      "updated_at": "2026-02-23T14:11:28.074Z"
    },
    {
      "uuid": "b542c59d-f22d-4b90-921f-5759797dfd52",
      "view_position": 2,
      "name": "squad-not-default",
      "created_at": "2026-02-23T14:27:19.332Z",
      "updated_at": "2026-02-23T14:27:19.331Z"
    }
  ],
  "count": 2
}
```

---

### Получить сквад по UUID

**GET** `/api/v1/internal-squads/{uuid}`

Возвращает данные одного сквада по UUID.

**Пример запроса:**
```bash
curl "http://localhost:9952/api/v1/internal-squads/2ef69520-9dc4-44f7-a492-c709fc34c41f"
```

**Пример ответа:**
```json
{
  "squad": {
    "uuid": "2ef69520-9dc4-44f7-a492-c709fc34c41f",
    "view_position": 1,
    "name": "Default-Squad",
    "created_at": "2026-02-23T14:11:28.074Z",
    "updated_at": "2026-02-23T14:11:28.074Z"
  }
}
```

---

### Создать сквад

**POST** `/api/v1/internal-squads`

Создаёт новый внутренний сквад.

**Пример запроса:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-squad",
    "view_position": 3
  }' \
  http://localhost:9952/api/v1/internal-squads
```

**Пример ответа (201 Created):**
```json
{
  "message": "internal squad created",
  "uuid": "c7f8a9b0-1234-5678-9abc-def012345678"
}
```

**Обязательные поля:**
- `name` — уникальное имя сквада

**Опциональные поля:**
- `view_position` — позиция отображения (по умолчанию 0)

**Примеры ошибок:**
```json
// 400 Bad Request - отсутствует имя
{
  "error": "name is required"
}

// 409 Conflict - имя занято
{
  "error": "name already exists"
}
```

---

### Обновить сквад (частичное обновление)

**PATCH** `/api/v1/internal-squads/{uuid}`

Обновляет только указанные поля сквада.

**Пример запроса:**
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "name": "updated-squad",
    "view_position": 5
  }' \
  http://localhost:9952/api/v1/internal-squads/2ef69520-9dc4-44f7-a492-c709fc34c41f
```

**Пример ответа:**
```json
{
  "squad": {
    "uuid": "2ef69520-9dc4-44f7-a492-c709fc34c41f",
    "view_position": 5,
    "name": "updated-squad",
    "created_at": "2026-02-23T14:11:28.074Z",
    "updated_at": "2026-02-23T14:30:00Z"
  }
}
```

**Доступные поля для обновления:**
| Поле | Тип | Описание |
|------|-----|----------|
| `name` | string | Имя сквада (уникальное) |
| `view_position` | int | Позиция отображения |

---

### Удалить сквад

**DELETE** `/api/v1/internal-squads/{uuid}`

Удаляет сквад по UUID. Связанные записи в `internal_squad_inbounds` и `internal_squad_members` удаляются автоматически через CASCADE.

**Пример запроса:**
```bash
curl -X DELETE \
  http://localhost:9952/api/v1/internal-squads/2ef69520-9dc4-44f7-a492-c709fc34c41f
```

**Пример ответа:**
```json
{
  "message": "internal squad deleted",
  "uuid": "2ef69520-9dc4-44f7-a492-c709fc34c41f",
  "name": "Default-Squad"
}
```

**Примеры ошибок:**
```json
// 404 Not Found - сквад не найден
{
  "error": "internal squad not found"
}

// 500 Internal Server Error - ошибка при удалении
{
  "error": "failed to delete internal squad"
}
```

---

## Структура sing-box конфигурации

Базовая структура конфигурации sing-box:

```json
{
  "log": {
    "level": "info",
    "timestamp": true
  },
  "dns": {
    "servers": [
      {"tag": "google", "address": "8.8.8.8"},
      {"tag": "local", "address": "223.5.5.5", "detour": "direct"}
    ],
    "rules": [
      {"outbound": "any", "server": "google"}
    ]
  },
  "inbounds": [
    {
      "type": "mixed",
      "tag": "mixed-in",
      "listen": "0.0.0.0",
      "listen_port": 2080,
      "users": []
    }
  ],
  "outbounds": [
    {"type": "direct", "tag": "direct"},
    {"type": "block", "tag": "block"}
  ],
  "route": {
    "rules": [
      {"ip_is_private": true, "outbound": "direct"}
    ],
    "auto_detect_interface": true
  },
  "experimental": {
    "cache_file": {"enabled": true}
  }
}
```

**Основные секции:**
- `log` — настройки логирования
- `dns` — DNS серверы и правила
- `inbounds` — входящие соединения (mixed, socks, http, shadowsocks, vmess, vless, trojan)
- `outbounds` — исходящие соединения (direct, block, socks, http, shadowsocks, vmess, vless, trojan)
- `route` — правила маршрутизации
- `experimental` — экспериментальные функции (cache_file, clash_api)
