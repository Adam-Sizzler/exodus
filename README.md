# Exodus Node

`exodus-node` — легковесная gRPC-нода для панели Exodus. Она запускает sing-box под `s6-overlay`, отдаёт панели статистику через gRPC, принимает задачи деплоя конфигурации sing-box, а также управляет файлом авторизации и runtime-кэшем HAProxy при включенном модуле `haproxyAuth`.

Управление пользователями внутри ядра sing-box отключено: пользователи и конфигурации собираются на стороне панели, а нода применяет готовый sing-box JSON.

## Основные настройки

```env
LOG_LEVEL=info
NODE_GRPC_ADDRESS=0.0.0.0
NODE_GRPC_PORT=2222
NODE_GRPC_PATH=/node/
```

## Сборка

```bash
cd /home/docker/projectSB/exodus-node
docker build -t exodus-node:local .
```

Для другой версии sing-box:

```bash
docker build \
  --build-arg SINGBOX_VERSION=v1.13.5 \
  -t exodus-node:local .
```

## Локальные проверки и тесты

```bash
cd /home/docker/projectSB/exodus-node
go test ./...
go build ./...
```

Проверить версию бинарника:

```bash
go run . --version
```

## Операции gRPC

Панель использует:

- `GetApiStats` — разовый сбор статистики из sing-box;
- `StreamNodeData` — streaming stats с дефолтным интервалом 15 секунд;
- `SubmitTask(operation=deploy_config)` — применение конфигурации sing-box и сопутствующих модулей;

Операции управления пользователями (`AddUsers`, `DeleteUsers`, `SetUserEnabled`, `ListUsers`) намеренно возвращают `Unimplemented`, так как нода работает в stats/deploy режиме.

---

## Модуль интеграции с HAProxy

Когда в плагине ноды включен модуль `haproxyAuth`, нода автоматически выполняет подготовку файла учетных записей и горячее перечитывание HAProxy:

1. **Генерация `/opt/app/haproxy/data/users.csv`**:
   - Формат строго двухколоночный: `<username>,<credential>`.
   - Поддерживаемые протоколы:
     - **VLESS**: `<username>,<vless_uuid>`
     - **Trojan**: `<username>,<sha224_hash>`
     - **AnyTLS**: `<username>,<sha256_hash>`
     - **NaiveProxy**: `<username>,basic:<base64(username:naive_password)>`
   - Если учетные данные не изменились по сравнению с файлом на диске, запись пропускается (идемпотентность).

2. **Горячая перезагрузка (Hot Reload)**:
   - При обнаружении изменений файла нода отправляет команду:
     ```text
     lua reload users\n
     ```
     в UNIX-сокет `/var/run/haproxy/haproxy.sock`.
   - HAProxy моментально обновляет свой кэш в оперативной памяти без перезапуска контейнера и без сброса активных сессий клиентов.
