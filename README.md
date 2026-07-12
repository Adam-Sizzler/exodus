# Exodus

Exodus - self-hosted панель для управления нодами

## Быстрые Команды

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
