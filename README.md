# Exodus

Exodus - self-hosted панель для управления нодами, пользователями, подписками, роутингом по хостам, шаблонами, response rules, HWID-лимитами, SRS-списками и инфраструктурной биллинговой статистикой.

Backend написан на Go и отдает API вместе с собранным React/Vite frontend. PostgreSQL хранит состояние панели, Valkey/Redis используется для фоновых задач. Полный стек с node/subscription сервисами запускается из compose в `/home/docker/projectSB`.

## Быстрые Команды

Сброс пароля admin через CLI:

```bash
cd /home/docker/projectSB
docker exec -it exodus exodus
```

В меню выберите `Reset admin password`.

Подготовить старую БД к новому baseline без переноса:

```bash
cd /home/docker/projectSB/exodus
./scripts/exodus-db-migration.sh prepare
```

Создать готовый архив БД для переноса на другую машину:

```bash
cd /home/docker/projectSB/exodus
./scripts/exodus-db-migration.sh export /home/docker/projectSB/backups/exodus-db.tar.gz
```

Развернуть архив на новой машине:

```bash
cd /home/docker/projectSB/exodus
./scripts/exodus-db-migration.sh import /path/to/exodus-db.tar.gz --force
```

## Структура Репозитория

```text
backend/
  main.go                     тонкая точка входа backend
  cmd/exodus/                 lifecycle, servers, static UI и emergency CLI
  internal/                   приватные backend-пакеты Go
  internal/db/prisma/         Prisma schema и baseline-миграция релиза
scripts/                      one-shot миграция и перенос PostgreSQL
frontend/                     React/Vite панель
Dockerfile                    multi-stage сборка frontend + Go backend
docker-compose.yml            локальный compose сервиса
compose.example.yml           минимальный пример compose для панели
../docker-compose.yml         рекомендуемый compose стек projectSB
```

Go module находится в `backend/go.mod`. Go-команды нужно запускать из `backend`, а не из корня репозитория.

## Политика Миграций

Этот релиз намеренно не поддерживает старую цепочку миграций. В поставке осталась одна baseline-миграция:

```text
backend/internal/db/prisma/migrations/20260506000000_initial_schema/migration.sql
```

Релиз рассчитан на чистую БД или на БД, уже подготовленную под этот baseline. Если в `schema_migrations` есть старые миграции вида `202603...`, backend быстро завершит запуск с явной ошибкой `unsupported legacy migration history`.

Для перехода со старой БД используйте one-shot скрипт:

```bash
cd /home/docker/projectSB/exodus
./scripts/exodus-db-migration.sh prepare
```

Скрипт проверяет, что старая БД находится на последней поддерживаемой legacy-миграции, сохраняет старую таблицу `schema_migrations` в `/home/docker/projectSB/backups/schema_migrations-legacy-*.sql`, затем заменяет историю миграций на baseline `20260506000000_initial_schema`. Данные пользователей, нод, hosts, subscription templates и других таблиц не удаляются.

Для переноса между машинами:

```bash
# На старой машине
cd /home/docker/projectSB/exodus
./scripts/exodus-db-migration.sh export /home/docker/projectSB/backups/exodus-db.tar.gz

# На новой машине
cd /home/docker/projectSB/exodus
./scripts/exodus-db-migration.sh import /path/to/exodus-db.tar.gz --force
```

Для чистой переустановки удаляйте только те volumes, данные в которых можно потерять:

```bash
cd /home/docker/projectSB
docker compose down
docker volume rm exodus-data panel-data
docker compose up -d --build
```

Не удаляйте production volumes без backup.

## Сборка И Запуск

Основной путь сборки из корня `projectSB`:

```bash
cd /home/docker/projectSB
docker compose build exodus
docker compose up -d exodus-db exodus-redis exodus
curl -fsS http://127.0.0.1:${APP_PORT:-3000}/api/health
```

Полный локальный стек:

```bash
cd /home/docker/projectSB
docker compose up -d --build
```

Остановка:

```bash
cd /home/docker/projectSB
docker compose down
```

## Compose

Минимальный compose для панели есть в `compose.example.yml`. Он оформлен в том же стиле, что и основной `/home/docker/projectSB/docker-compose.yml`: общие настройки вынесены в YAML anchors `x-common`, `x-logging`, `x-env`, все сервисы сидят в одной сети `exodus-network`.

Runtime-переменные не задаются через `environment`; они берутся из `.env`. Базовый `APP_PATH=/`, поэтому UI и API открываются от корня. Если нужен префикс вроде `/panel`, его нужно явно задать в `.env`.

```yaml
name: exodus

x-common: &common
  restart: always
  networks:
    - exodus-network

x-logging: &logging
  logging:
    driver: json-file
    options:
      max-size: 100m
      max-file: 5

x-env: &env
  env_file:
    - path: ./.env
      required: false

services:
  exodus:
    build:
      context: ./exodus
      dockerfile: Dockerfile
    image: exodus:latest
    container_name: exodus
    hostname: exodus
    <<: [*common, *logging, *env]
    ports:
      - "127.0.0.1:${APP_PORT:-3000}:${APP_PORT:-3000}"
      - "127.0.0.1:${METRICS_PORT:-3001}:${METRICS_PORT:-3001}"
    volumes:
      - panel-data:/app/data
    depends_on:
      exodus-db:
        condition: service_healthy
      exodus-redis:
        condition: service_healthy

  exodus-db:
    image: postgres:18
    container_name: exodus-db
    hostname: exodus-db
    <<: [*common, *logging, *env]
    ports:
      - "127.0.0.1:6868:5432"
    volumes:
      - exodus-data:/var/lib/postgresql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U \"$${POSTGRES_USER:-postgres}\" -d \"$${POSTGRES_DB:-postgres}\""]
      interval: 3s
      timeout: 10s
      retries: 10

  exodus-redis:
    image: valkey/valkey:9-alpine
    container_name: exodus-redis
    hostname: exodus-redis
    <<: [*common, *logging]
    command: >
      valkey-server
      --save ""
      --appendonly no
      --maxmemory-policy noeviction
      --loglevel warning
    healthcheck:
      test: ["CMD", "valkey-cli", "ping"]
      interval: 3s
      timeout: 3s
      retries: 10

networks:
  exodus-network:
    name: exodus-network
    driver: bridge
    external: false

volumes:
  panel-data:
    name: panel-data
  exodus-data:
    name: exodus-data
```

Полный `/home/docker/projectSB/docker-compose.yml` дополнительно поднимает `exodus-node`, `exodus-subscription`, HAProxy и nginx.

## Сброс Пароля Админа

В Exodus rescue CLI доступен короткой командой внутри контейнера, аналогично Remnawave:

```bash
docker exec -it exodus exodus
```

В меню выберите `Reset admin password`. Команда интерактивно спросит username, новый пароль и подтверждение. После этого она:

- обновит пароль существующего admin;
- включит password authentication в `exodus_settings`;
- удалит активные сессии этого admin.

Прямой non-menu вариант тоже доступен:

```bash
cd /home/docker/projectSB
docker compose exec exodus exodus --reset-admin-password
```

Если контейнер `exodus` остановлен, сначала поднимите PostgreSQL и Redis, затем запустите одноразовый контейнер:

```bash
cd /home/docker/projectSB
docker compose up -d exodus-db exodus-redis
docker compose run --rm --no-deps --entrypoint /app/exodus exodus --reset-admin-password
```

Admin должен уже существовать. На новой чистой БД первый admin создается через web login page.

## Сравнение С Remnawave

В Remnawave backend - Node/NestJS сервис с отдельным rescue CLI. В его `Reset superadmin` flow первый admin удаляется, кэш настроек очищается, а новый admin затем регистрируется через web UI.

Типичный flow Remnawave:

```bash
docker compose exec remnawave remnawave
# В интерактивном меню выбрать "Reset superadmin".
```

В Exodus модель проще: один Go-бинарник, который напрямую меняет пароль существующего admin и не удаляет аккаунт.

```bash
docker exec -it exodus exodus
```

## Проверки Для Разработки

Backend:

```bash
cd /home/docker/projectSB/exodus/backend
go test ./...
go build -trimpath -o /tmp/exodus-backend .
```

Frontend audit:

```bash
cd /home/docker/projectSB/exodus/frontend
npm audit --omit=dev
```

Docker image:

```bash
cd /home/docker/projectSB
docker compose build exodus
```

## Основные Возможности

- Панель и API на одном `APP_PORT`.
- Emergency CLI для сброса пароля admin.
- PostgreSQL baseline schema встроена в Go-бинарник.
- Мониторинг нод и subscription нод через gRPC/mTLS.
- Шаблоны подписок и response rules для Xray/Base64, Clash, Mihomo, Stash и Sing-box.
- Управление hosts, config profiles, internal squads и external squads.
- HWID device tracking и история запросов подписок.
- Prometheus-compatible metrics endpoint.
