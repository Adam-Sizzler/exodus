#!/bin/sh
set -eu

rm -f /run/supervisord-*.sock 2>/dev/null || true
rm -f /run/supervisord-*.pid 2>/dev/null || true

generate_random() {
    local length="${1:-64}"
    tr -dc 'a-zA-Z0-9' < /dev/urandom | head -c "$length"
}

RNDSTR="$(generate_random 10)"

SUPERVISORD_USER="${SUPERVISORD_USER:-$(generate_random 32)}"
SUPERVISORD_PASSWORD="${SUPERVISORD_PASSWORD:-$(generate_random 48)}"
SUPERVISORD_SOCKET_PATH="${SUPERVISORD_SOCKET_PATH:-/run/supervisord-${RNDSTR}.sock}"
SUPERVISORD_PID_PATH="${SUPERVISORD_PID_PATH:-/run/supervisord-${RNDSTR}.pid}"

export SUPERVISORD_USER
export SUPERVISORD_PASSWORD
export SUPERVISORD_SOCKET_PATH
export SUPERVISORD_PID_PATH

mkdir -p /run /var/log/supervisor /app/singbox /app/logs /app/certs

if [ ! -s /app/singbox/config.json ]; then
    cp /app/singbox/config.default.json /app/singbox/config.json
fi

echo "[Entrypoint] supervisord version: $(supervisord --version | head -n 1)"
supervisord -c /etc/supervisord.conf
sleep 1

echo "[Entrypoint] sing-box version: $(/usr/local/bin/sing-box version | head -n 1)"
echo "[Entrypoint] cerberus-node command: $*"

exec "$@"
