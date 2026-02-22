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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9952/api/v1/users-list
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
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001
```

---

### Создать пользователя

**POST** `/api/v1/users-list/create`

Создаёт нового пользователя в таблице `users`.

**Пример запроса:**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"description": ""}' \
  http://localhost:9952/api/v1/users-list/550e8400-e29b-41d4-a716-446655440001
```

---

## Config Profiles API

API для управления конфигурациями sing-box в формате JSON.

### Получить все конфигурации

**GET** `/api/v1/config-profiles`

Возвращает список всех конфигураций с их JSON данными.

**Пример запроса:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9952/api/v1/config-profiles
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
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9952/api/v1/config-profiles/550e8400-e29b-41d4-a716-446655440000
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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
  -H "Authorization: Bearer YOUR_TOKEN" \
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

---

## Статистика

### Получить статистику пользователей

**GET** `/api/v1/users`

Возвращает пользователей, сгруппированных по нодам.

**Пример запроса:**
```bash
curl http://localhost:9952/api/v1/users
```

**Пример ответа:**
```json
[
  {
    "node_name": "node1",
    "address": "192.168.1.10",
    "users": [
      {
        "user": "testuser1",
        "inbounds": [
          {
            "inbound_tag": "vless-in",
            "id": "550e8400-e29b-41d4-a716-446655440000"
          }
        ],
        "rate": "200",
        "enabled": "true",
        "uplink": 300,
        "downlink": 400
      }
    ]
  }
]
```

---

### Получить статистику сервера

**GET** `/api/v1/server_stats`

**Пример запроса:**
```bash
curl http://localhost:9952/api/v1/server_stats?node=node1,node2&sort_by=rate&sort_order=desc
```

---

### Получить статистику клиентов

**GET** `/api/v1/client_stats`

**Пример запроса:**
```bash
curl http://localhost:9952/api/v1/client_stats?user=user1&sort_by=user&sort_order=asc
```

---

### Получить DNS статистику

**GET** `/api/v1/dns_stats`

**Пример запроса:**
```bash
curl "http://localhost:9952/api/v1/dns_stats?node=node1&domain=example.com&count=50"
```

---

## Управление пользователями

### Добавить пользователей на ноды

**POST** `/api/v1/add_user`

**Пример запроса:**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "users": ["user1", "user2"],
    "inbound_tag": "vless-in",
    "nodes": ["node1", "node2"]
  }' \
  http://localhost:9952/api/v1/add_user
```

---

### Удалить пользователей с нод

**POST** `/api/v1/delete_user`

**Пример запроса:**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "users": ["user1"],
    "inbound_tag": "vless-in",
    "nodes": ["node1"]
  }' \
  http://localhost:9952/api/v1/delete_user
```

---

### Включить/выключить пользователей

**PATCH** `/api/v1/set_user_enabled`

**Пример запроса:**
```bash
curl -X PATCH \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "users": ["user1"],
    "enabled": true,
    "inbound_tag": "vless-in",
    "nodes": ["node1", "node2"]
  }' \
  http://localhost:9952/api/v1/set_user_enabled
```

---

## Сброс статистики

### Сброс DNS статистики

**POST** `/api/v1/reset_dns_stats`

**Пример запроса:**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:9952/api/v1/reset_dns_stats?nodes=node1,node2"
```

---

### Сброс трафика

**POST** `/api/v1/reset_bound_traffic`

**Пример запроса:**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:9952/api/v1/reset_bound_traffic?nodes=node1"
```

---

### Сброс статистики клиентов

**POST** `/api/v1/reset_user_traffic`

**Пример запроса:**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:9952/api/v1/reset_user_traffic?nodes=node1,node2"
```

---

## Коды ответов

| Код | Описание |
|-----|----------|
| 200 | Успешный запрос |
| 201 | Ресурс создан |
| 400 | Неверный запрос (валидация) |
| 401 | Неверный токен |
| 404 | Ресурс не найден |
| 405 | Метод не разрешён |
| 409 | Конфликт (например, username занят) |
| 500 | Внутренняя ошибка сервера |

---

## Примечания

- Все даты в формате ISO 8601 (RFC3339): `2026-12-31T23:59:59Z`
- UUID в формате: `550e8400-e29b-41d4-a716-446655440000`
- PATCH запросы обновляют только указанные поля
- Пустое строковое значение (`""`) в PATCH запросе устанавливает поле в NULL
