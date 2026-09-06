# Exodus Node

`exodus-node` - легковесная gRPC-нода для панели Exodus. Она запускает sing-box под `s6-overlay`, отдает панели статистику через gRPC, принимает задачи деплоя конфигурации sing-box, а также обновляет runtime cache пользователей HAProxy при включенном модуле.

Это не старый HTTP `v2ray-stat` API. Управление пользователями внутри ядра отключено: пользователи и конфиги собираются на стороне панели, а нода применяет готовый sing-box JSON.


Основные настройки:

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

## Локальные Проверки

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

- `GetApiStats` - разовый сбор stats из sing-box;
- `StreamNodeData` - streaming stats с дефолтным интервалом 15 секунд;
- `SubmitTask(operation=deploy_config)` - применить sing-box config;

Операции управления пользователями (`AddUsers`, `DeleteUsers`, `SetUserEnabled`, `ListUsers`) намеренно возвращают `Unimplemented`, потому что нода работает в stats/deploy режиме.
