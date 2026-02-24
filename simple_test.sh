#!/bin/bash

# Простой тестовый скрипт для проверки API привязок

BASE_URL="http://localhost:9243"
OUTPUT_FILE="/tmp/api_test_output.txt"

echo "=========================================="
echo "Тест API привязок"
echo "=========================================="
echo ""

# Шаг 1: Создаём тестовые данные
echo "1. Создаём тестовые данные..."

# Пользователь
curl -s -X POST "$BASE_URL/api/v1/users-list/create" \
  -H 'Content-Type: application/json' \
  -d '{"username":"api_test_user","status":"ACTIVE","expire_at":"2027-12-31T23:59:59Z"}' \
  -o /tmp/user_response.json

cat /tmp/user_response.json
echo ""

# Нода
curl -s -X POST "$BASE_URL/api/v1/nodes" \
  -H 'Content-Type: application/json' \
  -d '{"name":"API-Test-Node","address":"192.168.1.100","port":9253,"country_code":"US"}' \
  -o /tmp/node_response.json

cat /tmp/node_response.json
echo ""

# Конфиг с inbound
curl -s -X POST "$BASE_URL/api/v1/config-profiles" \
  -H 'Content-Type: application/json' \
  -d '{"name":"API-Test-Config","config":{"inbounds":[{"type":"shadowsocks","tag":"SS-Test","listen":"0.0.0.0","listen_port":2080}],"outbounds":[{"type":"direct","tag":"direct"}]}}' \
  -o /tmp/config_response.json

cat /tmp/config_response.json
echo ""

# Сквад
curl -s -X POST "$BASE_URL/api/v1/internal-squads" \
  -H 'Content-Type: application/json' \
  -d '{"name":"API-Test-Squad"}' \
  -o /tmp/squad_response.json

cat /tmp/squad_response.json
echo ""

echo ""
echo "2. Получаем созданные UUID..."

# Получаем конфиги с inbound
curl -s "$BASE_URL/api/v1/config-profiles-with-inbounds" -o /tmp/inbounds_response.json
echo "Inbounds response saved to /tmp/inbounds_response.json"

# Извлекаем UUID из файлов (используем grep и sed)
grep -o '"uuid":"[^"]*"' /tmp/node_response.json | head -1 | sed 's/"uuid":"//;s/"//' > /tmp/node_uuid.txt
grep -o '"uuid":"[^"]*"' /tmp/squad_response.json | head -1 | sed 's/"uuid":"//;s/"//' > /tmp/squad_uuid.txt
grep -o '"uuid":"[^"]*"' /tmp/inbounds_response.json | head -2 | tail -1 | sed 's/"uuid":"//;s/"//' > /tmp/inbound_uuid.txt
grep -o '"t_id":[0-9]*' /tmp/user_response.json | cut -d':' -f2 > /tmp/user_id.txt

echo "Node UUID: $(cat /tmp/node_uuid.txt)"
echo "Squad UUID: $(cat /tmp/squad_uuid.txt)"
echo "Inbound UUID: $(cat /tmp/inbound_uuid.txt)"
echo "User ID: $(cat /tmp/user_id.txt)"
echo ""

# Шаг 3: Тестируем привязку inbound к скваду
echo "3. Тестируем привязку inbound к скваду..."

NODE_UUID=$(cat /tmp/node_uuid.txt)
SQUAD_UUID=$(cat /tmp/squad_uuid.txt)
INBOUND_UUID=$(cat /tmp/inbound_uuid.txt)
USER_ID=$(cat /tmp/user_id.txt)

curl -s -X POST "$BASE_URL/api/v1/squad-inbounds" \
  -H 'Content-Type: application/json' \
  -d "{\"squad_uuid\":\"$SQUAD_UUID\",\"inbound_uuids\":[\"$INBOUND_UUID\"]}" \
  -o /tmp/squad_inbounds_response.json

echo "Squad inbounds response:"
cat /tmp/squad_inbounds_response.json
echo ""

# Шаг 4: Тестируем привязку пользователя к скваду
echo "4. Тестируем привязку пользователя к скваду..."

curl -s -X POST "$BASE_URL/api/v1/squad-members" \
  -H 'Content-Type: application/json' \
  -d "{\"squad_uuid\":\"$SQUAD_UUID\",\"user_ids\":[$USER_ID]}" \
  -o /tmp/squad_members_response.json

echo "Squad members response:"
cat /tmp/squad_members_response.json
echo ""

# Шаг 5: Тестируем привязку inbound к ноде
echo "5. Тестируем привязка inbound к ноде..."

curl -s -X POST "$BASE_URL/api/v1/inbound-assignments" \
  -H 'Content-Type: application/json' \
  -d "{\"node_uuid\":\"$NODE_UUID\",\"inbound_uuids\":[\"$INBOUND_UUID\"]}" \
  -o /tmp/node_inbounds_response.json

echo "Node inbounds response:"
cat /tmp/node_inbounds_response.json
echo ""

# Шаг 6: Проверяем результат
echo "6. Проверка результата..."

echo ""
echo "Squad inbounds:"
curl -s "$BASE_URL/api/v1/squad-inbounds?squad_uuid=$SQUAD_UUID" | head -20
echo ""

echo "Squad members:"
curl -s "$BASE_URL/api/v1/squad-members?squad_uuid=$SQUAD_UUID" | head -20
echo ""

echo "Node inbounds:"
curl -s "$BASE_URL/api/v1/inbound-assignments?node_uuid=$NODE_UUID" | head -20
echo ""

echo "Squad details:"
curl -s "$BASE_URL/api/v1/squad-details/$SQUAD_UUID" | head -30
echo ""

echo "=========================================="
echo "✅ Тест завершён!"
echo "=========================================="
echo ""
echo "Проверка в БД (sqlite3):"
echo "  sqlite3 data.db \"SELECT * FROM internal_squad_inbounds;\""
echo "  sqlite3 data.db \"SELECT * FROM config_profile_inbounds_to_nodes;\""
