# Exodus Node

`exodus-node` - легковесная gRPC-нода для панели Exodus. Она запускает sing-box под `s6-overlay`, отдает панели статистику через gRPC, принимает задачи деплоя конфигурации sing-box, а также обновляет runtime cache пользователей HAProxy при включенном модуле.

Это не старый HTTP `v2ray-stat` API. Управление пользователями внутри ядра отключено: пользователи и конфиги собираются на стороне панели, а нода применяет готовый sing-box JSON.

## Что Делает Нода

- поднимает gRPC сервер на `NODE_GRPC_ADDRESS:NODE_GRPC_PORT`;
- поддерживает два режима подключения: `gRPC + mTLS` через `SECRET_KEY` или `gRPC + TLS + token` через nginx/reverse proxy и `NODE_GRPC_TOKEN`;
- подключается к sing-box `experimental.v2ray_api` на `127.0.0.1:10085`;
- стримит inbound/outbound/user traffic stats в панель;
- добавляет/обновляет `experimental.v2ray_api` в полученном sing-box config;
- пишет config в `/opt/app/singbox/config.json`;
- управляет жизненным циклом sing-box через супервизор `s6-overlay` (`s6-svc`);
- обновляет `/opt/app/haproxy/data/users.csv` и reload users через HAProxy runtime socket;
- предоставляет утилиты `slogs` (полный живой лог) и `serrors` (фильтр ошибок и предупреждений).

## Структура

```text
main.go                    точка входа
config/                    env config, SECRET_KEY payload, token auth
grpcserver/                gRPC HTTP2 сервер, mTLS/token interceptors
server/                    NodeService handlers, deploy, stream, HAProxy
api/                       facade над sing-box stats API
sdk/                       sing-box v2ray_api gRPC client
proto/                     Exodus node gRPC contract
deploy/                    s6-overlay сервисы и конфигурация
Dockerfile                 multi-stage image: Go node + custom sing-box
```

## Runtime Переменные

Для `gRPC + mTLS`:

```env
SECRET_KEY=<base64-json-payload-from-panel>
```

Payload должен содержать:

```json
{
  "caCertPem": "-----BEGIN CERTIFICATE-----\\n...",
  "nodeCertPem": "-----BEGIN CERTIFICATE-----\\n...",
  "nodeKeyPem": "-----BEGIN PRIVATE KEY-----\\n...",
  "jwtPublicKey": "-----BEGIN PUBLIC KEY-----\\n..."
}
```

Для `gRPC + TLS + token` за nginx/reverse proxy:

```env
NODE_GRPC_TOKEN=<grpc-token-from-panel>
```

В token-режиме TLS завершается на nginx, а сама нода слушает h2c внутри приватной сети или на localhost.

Основные настройки:

```env
LOG_LEVEL=info
NODE_GRPC_ADDRESS=0.0.0.0
NODE_GRPC_PORT=2222
NODE_GRPC_PATH=/node/
SINGBOX_API_PORT=10085
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
- `StreamNodeData` - streaming stats с дефолтным интервалом 20 секунд;
- `SubmitTask(operation=deploy_config)` - применить sing-box config;

Операции управления пользователями (`AddUsers`, `DeleteUsers`, `SetUserEnabled`, `ListUsers`) намеренно возвращают `Unimplemented`, потому что нода работает в stats/deploy режиме.
