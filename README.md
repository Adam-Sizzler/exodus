# Exodus

Exodus - self-hosted панель для управления нодами

## Быстрые Команды

Подготовить старую БД к новому baseline без переноса:

```bash
./scripts/exodus-db-migration.sh prepare
```

Создать готовый архив БД для переноса на другую машину:

```bash
./scripts/exodus-db-migration.sh export /home/docker/projectSB/backups/exodus-db.tar.gz
```

Развернуть архив на новой машине:

```bash
./scripts/exodus-db-migration.sh import /path/to/exodus-db.tar.gz --force
```

Для перехода со старой БД используйте one-shot скрипт:

```bash
cd /home/docker/projectSB/exodus
./scripts/exodus-db-migration.sh prepare
```

Для переноса между машинами:

```bash
# На старой машине
./scripts/exodus-db-migration.sh export /home/docker/projectSB/backups/exodus-db.tar.gz

# На новой машине
./scripts/exodus-db-migration.sh import /path/to/exodus-db.tar.gz --force
```

Для чистой переустановки удаляйте только те volumes, данные в которых можно потерять:

```bash
docker compose down
docker volume rm exodus-data panel-data
docker compose up -d --build
```

## Сброс Пароля Админа

В Exodus rescue CLI доступен короткой командой внутри контейнера:

```bash
docker exec -it exodus exodus
```

В меню выберите `Reset admin password`. Команда интерактивно спросит username, новый пароль и подтверждение. После этого она:

- обновит пароль существующего admin;
- включит password authentication в `exodus_settings`;
- удалит активные сессии этого admin.

Прямой non-menu вариант тоже доступен:

```bash
docker compose exec exodus exodus --reset-admin-password
```

## Проверки Для Разработки

Backend:

```bash
go test ./...
go build -trimpath -o /tmp/exodus-backend .
```

Frontend audit:

```bash
npm audit --omit=dev
```
