## Модуль: HAPROXY

Кратко: модуль формирует таблицу учетных данных для Lua-авторизации в HAProxy и поддерживает ее в актуальном состоянии при деплое ноды с поддержкой горячей перезагрузки через сокет.

### Что делает модуль

- Создает и обновляет **один** файл пользователей `users.csv`.
- Добавляет в файл активных пользователей текущей ноды.
- Для каждого пользователя пишет `VLESS UUID`.
- Для каждого пользователя пишет `Trojan credential` как SHA-224 hash (56 символов hex).
- Для каждого пользователя пишет `AnyTLS credential` как SHA-256 hash (64 символа hex).
- Для каждого пользователя пишет `NaiveProxy credential` в формате `basic:<base64(username:password)>`.
- Отправляет команду `lua reload users` в UNIX-сокет HAProxy для мгновенного обновления кэша в оперативной памяти без рестарта контейнера.

### Какой файл создается

- Относительный путь в логике модуля: `haproxy/data/users.csv`
- Абсолютный путь внутри контейнера ноды: `/opt/app/haproxy/data/users.csv`
- Путь внутри контейнера HAProxy: `/etc/haproxy/data/users.csv`

### Формат строк

Строго двухколоночный формат без префиксов 1/0:

`<username>,<credential>`

### Пример содержимого

```text
pablo_escobar,90cdd126-819d-5f2f-a1dd-4cbf085182b9
pablo_escobar,basic:cGFibG9fZXNjb2JhcjpteXBhc3N3b3Jk
alvarez,9f37a85cafb66ab9c36f866e5c787b2e42a997716c02337f532a259c
guest,8f49d237686bee189aaedaba1ad434c5580356d8e8b5a3d702bd0228
```

### Для чего это нужно

- `haproxy/lua/auth.lua` читает пользователей из `/etc/haproxy/data/users.csv`.
- L4: Trojan определяется по credential, включая проверку MUX.
- L4: VLESS определяется по UUID из бинарного рукопожатия.
- L4: AnyTLS определяется по SHA-256 хэшу из бинарного рукопожатия.
- L7: NaiveProxy проверяется через sample-fetch `lua.auth_naive` по HTTP-заголовку `Proxy-Authorization: Basic <token>`.
- `haproxy/lua/users_loader.lua` разбирает CSV и раскладывает значения в таблицы:
  - `vless[uuid] -> username`
  - `trojan[sha224] -> username`
  - `anytls[sha256] -> username`
  - `naive[token] -> username`
- Конфигурация `haproxy.cfg` направляет валидный трафик в соответствующие бэкенды (VLESS, Trojan, AnyTLS, Naive) или передаёт веб-запросы в Nginx.

### Важно

- Модуль не изменяет `haproxy.cfg` и Lua-скрипты.
- Модуль управляет файлом `users.csv` и триггерит runtime reload через сокет.
