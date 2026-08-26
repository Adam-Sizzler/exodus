#!/command/with-contenv sh

echo "[init-env] preparing runtime environment..."

mkdir -p /run /opt/app/singbox /var/log/singbox

SINGBOX_CORE_VERSION=$(/usr/local/bin/sing-box version 2>/dev/null | head -n 1 || echo "N/A")
echo "[init-env] Sing-box version: $SINGBOX_CORE_VERSION"

echo "[init-env] done."
exit 0
