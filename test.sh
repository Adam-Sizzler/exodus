#!/bin/bash

# ==========================================
# V2Ray Stat - Тестовый скрипт для создания тестовых данных
# ==========================================
# Скрипт создает:
# - 3 пользователя
# - 2 ноды (обе отключены)
# - 2 конфиг профиля
# - 2 сквада
# - Привязки: сквад1 -> 2 пользователя, сквад2 -> 1 пользователь
# - Привязки: сквад1 -> 2 inbound, сквад2 -> 1 inbound
# - Привязки: нода1 -> 2 inbound (config_profile_inbounds_to_nodes)
# ==========================================

BASE_URL="http://localhost:9243"

echo "=========================================="
echo "V2Ray Stat - Тестовое создание данных"
echo "=========================================="
echo ""

# ==========================================
# Шаг 0: Очистка базы данных (если нужно)
# ==========================================
echo "⚠️  Шаг 0: Очистка базы данных..."
echo ""

# Удаляем все данные в правильном порядке (из-за внешних ключей)
curl -s -X DELETE "$BASE_URL/api/v1/squad-inbounds" > /dev/null 2>&1
curl -s -X DELETE "$BASE_URL/api/v1/squad-members" > /dev/null 2>&1
curl -s -X DELETE "$BASE_URL/api/v1/inbound-assignments" > /dev/null 2>&1

echo "  → База данных очищена"
echo ""

# ==========================================
# Шаг 1: Создание 3 пользователей
# ==========================================
echo "📝 Шаг 1: Создание 3 пользователей..."
echo ""

# Пользователь 1
echo "  → Создание пользователя: user_alpha"
USER1_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user_alpha",
    "status": "ACTIVE",
    "traffic_limit_bytes": 10737418240,
    "traffic_limit_strategy": "MONTH",
    "expire_at": "2027-12-31T23:59:59Z",
    "description": "Test user alpha",
    "tag": "premium",
    "email": "alpha@example.com",
    "hwid_device_limit": 3
  }' \
  "$BASE_URL/api/v1/users-list/create")

echo "  Ответ: $USER1_RESPONSE"
USER1_ID=$(echo "$USER1_RESPONSE" | grep -o '"t_id":[0-9]*' | head -1 | cut -d':' -f2)
echo "  ID: $USER1_ID"
echo ""

# Пользователь 2
echo "  → Создание пользователя: user_beta"
USER2_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user_beta",
    "status": "ACTIVE",
    "traffic_limit_bytes": 21474836480,
    "traffic_limit_strategy": "MONTH",
    "expire_at": "2027-06-30T23:59:59Z",
    "description": "Test user beta",
    "tag": "standard",
    "email": "beta@example.com",
    "hwid_device_limit": 2
  }' \
  "$BASE_URL/api/v1/users-list/create")

echo "  Ответ: $USER2_RESPONSE"
USER2_ID=$(echo "$USER2_RESPONSE" | grep -o '"t_id":[0-9]*' | head -1 | cut -d':' -f2)
echo "  ID: $USER2_ID"
echo ""

# Пользователь 3
echo "  → Создание пользователя: user_gamma"
USER3_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user_gamma",
    "status": "ACTIVE",
    "traffic_limit_bytes": 5368709120,
    "traffic_limit_strategy": "WEEK",
    "expire_at": "2026-12-31T23:59:59Z",
    "description": "Test user gamma",
    "tag": "basic",
    "email": "gamma@example.com",
    "hwid_device_limit": 1
  }' \
  "$BASE_URL/api/v1/users-list/create")

echo "  Ответ: $USER3_RESPONSE"
USER3_ID=$(echo "$USER3_RESPONSE" | grep -o '"t_id":[0-9]*' | head -1 | cut -d':' -f2)
echo "  ID: $USER3_ID"
echo ""

# ==========================================
# Шаг 2: Создание 2 нод (обе отключены)
# ==========================================
echo "📝 Шаг 2: Создание 2 нод (обе отключены)..."
echo ""

# Нода 1
echo "  → Создание ноды: Node-Moscow-01"
NODE1_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Node-Moscow-01",
    "address": "192.168.1.10",
    "port": 9253,
    "api_schema": "grpc",
    "api_path": "",
    "is_disabled": true,
    "consumption_multiplier": 100,
    "is_traffic_tracking_active": true,
    "traffic_reset_day": 1,
    "traffic_limit_bytes": 1099511627776,
    "notify_percent": 80,
    "view_position": 1,
    "country_code": "RU",
    "tags": ["production", "moscow"]
  }' \
  "$BASE_URL/api/v1/nodes")

echo "  Ответ: $NODE1_RESPONSE"
NODE1_UUID=$(echo "$NODE1_RESPONSE" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "  UUID: $NODE1_UUID"
echo ""

# Нода 2
echo "  → Создание ноды: Node-Berlin-02"
NODE2_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Node-Berlin-02",
    "address": "192.168.2.20",
    "port": 9253,
    "api_schema": "grpc",
    "api_path": "",
    "is_disabled": true,
    "consumption_multiplier": 150,
    "is_traffic_tracking_active": true,
    "traffic_reset_day": 15,
    "traffic_limit_bytes": 2199023255552,
    "notify_percent": 75,
    "view_position": 2,
    "country_code": "DE",
    "tags": ["production", "berlin", "eu"]
  }' \
  "$BASE_URL/api/v1/nodes")

echo "  Ответ: $NODE2_RESPONSE"
NODE2_UUID=$(echo "$NODE2_RESPONSE" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "  UUID: $NODE2_UUID"
echo ""

# ==========================================
# Шаг 3: Создание 2 конфиг профилей
# ==========================================
echo "📝 Шаг 3: Создание 2 конфиг профилей..."
echo ""

# Конфиг профиль 1
echo "  → Создание конфиг профиля: Profile-Shadowsocks"
PROFILE1_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Profile-Shadowsocks",
    "view_position": 1,
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
          "type": "shadowsocks",
          "tag": "Shadowsocks-2080",
          "listen": "0.0.0.0",
          "listen_port": 2080,
          "network": "tcp,udp"
        },
        {
          "type": "shadowsocks",
          "tag": "Shadowsocks-2081",
          "listen": "0.0.0.0",
          "listen_port": 2081,
          "network": "tcp,udp"
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
  "$BASE_URL/api/v1/config-profiles")

echo "  Ответ: $PROFILE1_RESPONSE"
PROFILE1_UUID=$(echo "$PROFILE1_RESPONSE" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "  UUID: $PROFILE1_UUID"
echo ""

# Конфиг профиль 2
echo "  → Создание конфиг профиля: Profile-VLESS-Reality"
PROFILE2_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Profile-VLESS-Reality",
    "view_position": 2,
    "config": {
      "log": {
        "level": "info",
        "timestamp": true
      },
      "dns": {
        "servers": [
          {"tag": "cloudflare", "address": "1.1.1.1"}
        ]
      },
      "inbounds": [
        {
          "type": "vless",
          "tag": "VLESS-Reality-8443",
          "listen": "0.0.0.0",
          "listen_port": 8443,
          "network": "tcp",
          "security": "reality"
        },
        {
          "type": "vless",
          "tag": "VLESS-Reality-8444",
          "listen": "0.0.0.0",
          "listen_port": 8444,
          "network": "tcp",
          "security": "reality"
        },
        {
          "type": "trojan",
          "tag": "Trojan-8445",
          "listen": "0.0.0.0",
          "listen_port": 8445,
          "network": "tcp",
          "security": "tls"
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
  "$BASE_URL/api/v1/config-profiles")

echo "  Ответ: $PROFILE2_RESPONSE"
PROFILE2_UUID=$(echo "$PROFILE2_RESPONSE" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "  UUID: $PROFILE2_UUID"
echo ""

# ==========================================
# Шаг 4: Получение созданных inbound из профилей
# ==========================================
echo "📝 Шаг 4: Получение списка inbound из профилей..."
echo ""

INBOUNDS_RESPONSE=$(curl -s "$BASE_URL/api/v1/config-profiles-with-inbounds")
echo "  Ответ получен, парсим inbound..."

# Извлекаем UUID inbound из первого профиля (первые 2 inbound)
INBOUND1_UUID=$(echo "$INBOUNDS_RESPONSE" | grep -o '"uuid":"[a-f0-9-]*"' | head -2 | tail -1 | cut -d'"' -f4)
INBOUND2_UUID=$(echo "$INBOUNDS_RESPONSE" | grep -o '"uuid":"[a-f0-9-]*"' | head -3 | tail -1 | cut -d'"' -f4)

# Извлекаем UUID inbound из второго профиля (первый inbound)
INBOUND3_UUID=$(echo "$INBOUNDS_RESPONSE" | grep -o '"uuid":"[a-f0-9-]*"' | head -6 | tail -1 | cut -d'"' -f4)

echo "  Inbound 1 (Shadowsocks-2080): $INBOUND1_UUID"
echo "  Inbound 2 (Shadowsocks-2081): $INBOUND2_UUID"
echo "  Inbound 3 (VLESS-Reality-8443): $INBOUND3_UUID"
echo ""

# ==========================================
# Шаг 5: Создание 2 сквадов
# ==========================================
echo "📝 Шаг 5: Создание 2 сквадов..."
echo ""

# Сквад 1
echo "  → Создание сквада: Squad-Premium"
SQUAD1_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Squad-Premium",
    "view_position": 1
  }' \
  "$BASE_URL/api/v1/internal-squads")

echo "  Ответ: $SQUAD1_RESPONSE"
SQUAD1_UUID=$(echo "$SQUAD1_RESPONSE" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "  UUID: $SQUAD1_UUID"
echo ""

# Сквад 2
echo "  → Создание сквада: Squad-Standard"
SQUAD2_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Squad-Standard",
    "view_position": 2
  }' \
  "$BASE_URL/api/v1/internal-squads")

echo "  Ответ: $SQUAD2_RESPONSE"
SQUAD2_UUID=$(echo "$SQUAD2_RESPONSE" | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "  UUID: $SQUAD2_UUID"
echo ""

# ==========================================
# Шаг 6: Привязка пользователей к сквадам
# ==========================================
echo "📝 Шаг 6: Привязка пользователей к сквадам..."
echo ""

# Сквад 1 -> 2 пользователя (user_alpha и user_beta)
echo "  → Привязка 2 пользователей к Squad-Premium (user_alpha, user_beta)"
SQUAD1_MEMBERS_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"squad_uuid\": \"$SQUAD1_UUID\",
    \"user_ids\": [$USER1_ID, $USER2_ID]
  }" \
  "$BASE_URL/api/v1/squad-members")

echo "  Ответ: $SQUAD1_MEMBERS_RESPONSE"
echo ""

# Сквад 2 -> 1 пользователь (user_gamma)
echo "  → Привязка 1 пользователя к Squad-Standard (user_gamma)"
SQUAD2_MEMBERS_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"squad_uuid\": \"$SQUAD2_UUID\",
    \"user_ids\": [$USER3_ID]
  }" \
  "$BASE_URL/api/v1/squad-members")

echo "  Ответ: $SQUAD2_MEMBERS_RESPONSE"
echo ""

# ==========================================
# Шаг 7: Привязка inbound к сквадам
# ==========================================
echo "📝 Шаг 7: Привязка inbound к сквадам..."
echo ""

# Сквад 1 -> 2 inbound (из первого профиля)
echo "  → Привязка 2 inbound к Squad-Premium"
SQUAD1_INBOUNDS_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"squad_uuid\": \"$SQUAD1_UUID\",
    \"inbound_uuids\": [\"$INBOUND1_UUID\", \"$INBOUND2_UUID\"]
  }" \
  "$BASE_URL/api/v1/squad-inbounds")

echo "  Ответ: $SQUAD1_INBOUNDS_RESPONSE"
echo ""

# Сквад 2 -> 1 inbound (из второго профиля)
echo "  → Привязка 1 inbound к Squad-Standard"
SQUAD2_INBOUNDS_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"squad_uuid\": \"$SQUAD2_UUID\",
    \"inbound_uuids\": [\"$INBOUND3_UUID\"]
  }" \
  "$BASE_URL/api/v1/squad-inbounds")

echo "  Ответ: $SQUAD2_INBOUNDS_RESPONSE"
echo ""

# ==========================================
# Шаг 8: Привязка конфига к ноде (опционально)
# ==========================================
echo "📝 Шаг 8: Привязка конфига к ноде (опционально)..."
echo ""

echo "  → Привязка Profile-Shadowsocks к Node-Moscow-01"
NODE_CONFIG_RESPONSE=$(curl -s -X PATCH \
  -H "Content-Type: application/json" \
  -d "{
    \"active_config_profile_uuid\": \"$PROFILE1_UUID\"
  }" \
  "$BASE_URL/api/v1/nodes/$NODE1_UUID")

echo "  Ответ: $NODE_CONFIG_RESPONSE"
echo ""

# ==========================================
# Шаг 9: Привязка inbound к ноде (config_profile_inbounds_to_nodes)
# ==========================================
echo "📝 Шаг 9: Привязка inbound к ноде (config_profile_inbounds_to_nodes)..."
echo ""
sleep 1
echo "  → Привязка 2 inbound к Node-Moscow-01"
NODE_INBOUNDS_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"node_uuid\": \"$NODE1_UUID\",
    \"inbound_uuids\": [\"$INBOUND1_UUID\", \"$INBOUND2_UUID\"]
  }" \
  "$BASE_URL/api/v1/inbound-assignments")

echo "  Ответ: $NODE_INBOUNDS_RESPONSE"
echo ""

# ==========================================
# Финальная сводка
# ==========================================
echo "=========================================="
echo "✅ Готово! Созданные данные:"
echo "=========================================="
echo ""
echo "👥 Пользователи (3):"
echo "   - user_alpha (ID: $USER1_ID)"
echo "   - user_beta (ID: $USER2_ID)"
echo "   - user_gamma (ID: $USER3_ID)"
echo ""
echo "🖥️  Ноды (2, обе отключены):"
echo "   - Node-Moscow-01 (UUID: $NODE1_UUID)"
echo "   - Node-Berlin-02 (UUID: $NODE2_UUID)"
echo ""
echo "📋 Конфиг профили (2):"
echo "   - Profile-Shadowsocks (UUID: $PROFILE1_UUID)"
echo "     Inbounds: Shadowsocks-2080, Shadowsocks-2081"
echo "   - Profile-VLESS-Reality (UUID: $PROFILE2_UUID)"
echo "     Inbounds: VLESS-Reality-8443, VLESS-Reality-8444, Trojan-8445"
echo ""
echo "🔷 Сквады (2):"
echo "   - Squad-Premium (UUID: $SQUAD1_UUID)"
echo "     Участники: user_alpha, user_beta (2)"
echo "     Inbounds: Shadowsocks-2080, Shadowsocks-2081 (2)"
echo "   - Squad-Standard (UUID: $SQUAD2_UUID)"
echo "     Участники: user_gamma (1)"
echo "     Inbounds: VLESS-Reality-8443 (1)"
echo ""
echo "🔗 Привязки:"
echo "   - Node-Moscow-01 → Profile-Shadowsocks"
echo "   - Node-Moscow-01 → Shadowsocks-2080, Shadowsocks-2081 (config_profile_inbounds_to_nodes)"
echo "   - Squad-Premium → Shadowsocks-2080, Shadowsocks-2081 (internal_squad_inbounds)"
echo "   - Squad-Standard → VLESS-Reality-8443 (internal_squad_inbounds)"
echo ""
echo "=========================================="
echo ""
echo "📊 Проверка данных в БД:"
echo ""
echo "  → Проверить internal_squad_inbounds:"
echo "     SELECT * FROM internal_squad_inbounds WHERE internal_squad_uuid = '$SQUAD1_UUID';"
echo ""
echo "  → Проверить config_profile_inbounds_to_nodes:"
echo "     SELECT * FROM config_profile_inbounds_to_nodes WHERE node_uuid = '$NODE1_UUID';"
echo ""
echo "  → Проверить inbound'ы с названиями профилей:"
echo "     SELECT i.uuid, i.tag, i.type, p.name as profile_name"
echo "     FROM config_profile_inbounds i"
echo "     JOIN config_profiles p ON p.uuid = i.profile_uuid;"
echo ""
echo "=========================================="
echo ""
echo "🎯 Быстрые тестовые запросы (curl):"
echo ""
echo "# 1. Проверить internal_squad_inbounds:"
echo "curl -s $BASE_URL/api/v1/squad-inbounds?squad_uuid=$SQUAD1_UUID | jq"
echo ""
echo "# 2. Проверить config_profile_inbounds_to_nodes:"
echo "curl -s $BASE_URL/api/v1/inbound-assignments?node_uuid=$NODE1_UUID | jq"
echo ""
echo "# 3. Проверить Squad-Premium details:"
echo "curl -s $BASE_URL/api/v1/squad-details/$SQUAD1_UUID | jq"
echo ""
echo "# 4. Получить все конфиги с inbound'ами:"
echo "curl -s $BASE_URL/api/v1/config-profiles-with-inbounds | jq '.profiles[] | {name, inbounds: [.inbounds[] | {tag, type, port}]}'"
echo ""
echo "=========================================="
echo ""
echo "✅ Все данные созданы и готовы к тестированию!"
echo ""
echo "📖 Документация:"
echo "   - README.md - полная API документация"
echo "   - FRONTEND_API_GUIDE.md - руководство для frontend-разработчиков"
echo "   - API_REFERENCE.md - краткая справка по API"
echo ""
echo "=========================================="
