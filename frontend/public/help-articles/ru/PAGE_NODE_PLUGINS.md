# Плагины для нод (Node Plugins)

Плагины – это дополнительные модули, которые можно активировать на Exodus Node для расширения функционала: **Ingress Filter**, **Egress Filter**, **HAProxy Auth**, **Shared Lists** (Общие списки) и команды **Executor**.

> **Важно:** Отнеситесь к настройке плагинов с осторожностью. Плагины Ingress Filter и Egress Filter работают напрямую с сетевым фаерволом (`nftables`) на вашем сервере.

---

## Требования

Для корректной работы плагинов требуется:
1. Директива `cap_add: NET_ADMIN` в `docker-compose.yml` ноды.
2. Версия ядра Linux **5.7** или выше (`uname -r`).
3. Утилита `nftables` на хостовой машине (`nft --version`).

```yaml title="docker-compose.yml"
services:
    exodus-node:
        container_name: exodus-node
        image: exodus/node:latest
        network_mode: host
        restart: always
        cap_add:
            - NET_ADMIN
        ulimits:
            nofile:
                soft: 1048576
                hard: 1048576
```

---

## Таблицы nftables

При включении плагинов фильтрации трафика Exodus Node автоматически формирует следующие таблицы правил в `nftables`:

```nftables
table ip exodus {
        counter processed {
                packets 32 bytes 2060
        }

        counter ingress-filter-ip {
                packets 0 bytes 0
        }

        counter egress-filter-ip {
                packets 0 bytes 0
        }

        counter egress-filter-port {
                packets 0 bytes 0
        }

        set ingress-filter-ip {
                type ipv4_addr
                flags timeout
                counter
        }

        set egress-filter-ip {
                type ipv4_addr
                flags timeout
                counter
        }

        set egress-filter-port {
                type inet_proto . inet_service
                flags timeout
                counter
        }

        chain input {
                type filter hook input priority filter - 10; policy accept;
                counter name "processed"
                ip saddr @ingress-filter-ip log prefix "ingress-filter-ip: " counter name "ingress-filter-ip" drop
        }

        chain forward {
                type filter hook forward priority filter - 10; policy accept;
                counter name "processed"
                ip saddr @ingress-filter-ip log prefix "ingress-filter-ip: " counter name "ingress-filter-ip" drop
        }

        chain output {
                type filter hook output priority filter - 10; policy accept;
                ip daddr @egress-filter-ip counter name "egress-filter-ip" drop
                meta l4proto . th dport @egress-filter-port counter name "egress-filter-port" drop
        }
}

table ip6 exodus6 {
        counter processed {
                packets 0 bytes 0
        }

        counter ingress-filter-ip6 {
                packets 0 bytes 0
        }

        counter egress-filter-ip6 {
                packets 0 bytes 0
        }

        counter egress-filter-port6 {
                packets 0 bytes 0
        }

        set ingress-filter-ip6 {
                type ipv6_addr
                flags timeout
                counter
        }

        set egress-filter-ip6 {
                type ipv6_addr
                flags timeout
                counter
        }

        set egress-filter-port6 {
                type inet_proto . inet_service
                flags timeout
                counter
        }

        chain input {
                type filter hook input priority filter - 10; policy accept;
                counter name "processed"
                ip6 saddr @ingress-filter-ip6 log prefix "ingress-filter-ip: " counter name "ingress-filter-ip6" drop
        }

        chain forward {
                type filter hook forward priority filter - 10; policy accept;
                counter name "processed"
                ip6 saddr @ingress-filter-ip6 log prefix "ingress-filter-ip: " counter name "ingress-filter-ip6" drop
        }

        chain output {
                type filter hook output priority filter - 10; policy accept;
                ip6 daddr @egress-filter-ip6 counter name "egress-filter-ip6" drop
                meta l4proto . th dport @egress-filter-port6 counter name "egress-filter-port6" drop
        }
}
```

---

## Поддержка CIDR

CIDR-нотация (IPv4 и IPv6) поддерживается в плагинах **Ingress Filter**, **Egress Filter** и **Shared Lists**.

Примеры допустимых CIDR-значений: `192.168.1.1`, `10.0.0.0/8`, `172.16.0.0/12`, `2001:db8::1`, `2001:db8::/32`.

---

## Структура конфигурации JSON

```json
{
    "ingressFilter": {
        "enabled": false,
        "blockedIps": []
    },
    "egressFilter": {
        "enabled": false,
        "blockedIps": [],
        "blockedPorts": []
    },
    "haproxyAuth": {
        "enabled": false,
        "inboundTags": ["*"]
    }
}
```

---

## Ingress Filter (Входящий фильтр)

Ingress Filter перманентно блокирует входящий трафик от указанных IP-адресов или CIDR-подсетей в фаерволе `nftables`.

```json
"ingressFilter": {
    "enabled": true,
    "blockedIps": ["192.168.1.1", "10.0.0.0/8", "ext:blocked-ips"]
}
```

---

## Egress Filter (Исходящий фильтр)

Egress Filter блокирует исходящие подключения с ноды к определённым IP-адресам или портам назначения (например, закрыть почтовые порты 25, 465, 587 или доступ к внутренним сетям).

```json
"egressFilter": {
    "enabled": true,
    "blockedIps": ["10.0.0.1", "ext:internal-networks"],
    "blockedPorts": [25, 465, 587]
}
```

---

## Модуль: HAProxy Auth

Модуль формирует таблицу учетных данных пользователей для быстрой авторизации в HAProxy через Lua-скрипты и автоматически обновляет ее при изменении конфигурации или пользователей.

```json
"haproxyAuth": {
    "enabled": true,
    "inboundTags": ["*"]
}
```

- В поле `inboundTags` можно указать `["*"]` для включения всех поддерживаемых инбаундов ноды или перечислить конкретные теги (например, `["vless-inbound", "trojan-inbound"]`).
- Создает и актуализирует файл `/app/haproxy/data/users.csv` внутри контейнера ноды.
- Записывает формат: `1,username,credential`.
- Сохраняет `VLESS UUID` и `Trojan credential` (SHA-224 hash) для всех пользователей ноды.

---

## Shared Lists (Общие списки)

Общие списки настраиваются в отдельном меню **Shared Lists** и позволяют создавать глобальные переиспользуемые коллекции IP-адресов или номеров автономных систем (ASN).

В плагинах они подключаются по имени в формате: `ext:имя_списка`.

Пример использования в Ingress Filter:
```json
"ingressFilter": {
    "enabled": true,
    "blockedIps": ["ext:my-blocked-subnets"]
}
```

---

## Executor (Команды выполнения)

Кнопка **Executor** в шапке раздела позволяет отправлять прямые команды на выбранные ноды:
- **Block IPs**: Временная блокировка IP-адресов в фаерволе `nftables` с возможностью указания таймаута в секундах.
- **Unblock IPs**: Снятие временной блокировки с IP-адресов.
- **Reset nftables**: Принудительный сброс и пересоздание таблиц правил `nftables`.
