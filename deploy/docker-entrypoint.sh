#!/bin/sh
set -eu

rm -f /run/supervisord-*.sock 2>/dev/null || true
rm -f /run/supervisord-*.pid 2>/dev/null || true

generate_random() {
    local length="${1:-64}"
    tr -dc 'a-zA-Z0-9' < /dev/urandom | head -c "$length"
}

normalize_version() {
    printf '%s' "${1:-}" | sed 's/^v//'
}

resolve_sb_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        armv7l|armv6l|arm) echo "arm" ;;
        i386|i686|386) echo "386" ;;
        *)
            echo "[Entrypoint] Unsupported CPU architecture: $(uname -m)" >&2
            exit 1
            ;;
    esac
}

installed_singbox_version() {
    if [ ! -x /usr/local/bin/sing-box ]; then
        return 0
    fi

    local line version
    line="$(/usr/local/bin/sing-box version 2>/dev/null | head -n 1 || true)"
    version="$(printf '%s' "$line" | awk '{print $NF}')"
    normalize_version "$version"
}

ensure_singbox_version() {
    local desired desired_tag desired_norm current_norm sb_arch url tmp_file

    desired="${SINGBOX_VERSION:-}"
    if [ -z "$desired" ]; then
        echo "[Entrypoint] SINGBOX_VERSION is not set, using bundled sing-box"
        return 0
    fi

    desired_tag="$desired"
    case "$desired_tag" in
        v*) ;;
        *) desired_tag="v${desired_tag}" ;;
    esac
    desired_norm="$(normalize_version "$desired_tag")"
    current_norm="$(installed_singbox_version)"

    if [ -n "$current_norm" ] && [ "$current_norm" = "$desired_norm" ]; then
        echo "[Entrypoint] sing-box version matches requested ${desired_tag}, skipping download"
        return 0
    fi

    sb_arch="${SB_ARCH:-$(resolve_sb_arch)}"
    url="https://github.com/Adam-Sizzler/sing-box-v2ray-api/releases/download/${desired_tag}/sing-box-linux-${sb_arch}"
    tmp_file="/tmp/sing-box-${desired_tag}-${sb_arch}.tmp"

    echo "[Entrypoint] Downloading sing-box ${desired_tag} for ${sb_arch}"
    curl -fL "$url" -o "$tmp_file"
    chmod +x "$tmp_file"

    if ! "$tmp_file" version >/dev/null 2>&1; then
        echo "[Entrypoint] Downloaded sing-box is invalid: ${url}" >&2
        rm -f "$tmp_file"
        exit 1
    fi

    mv "$tmp_file" /usr/local/bin/sing-box
    echo "[Entrypoint] Updated sing-box to: $(/usr/local/bin/sing-box version | head -n 1)"
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

ensure_singbox_version

echo "[Entrypoint] supervisord version: $(supervisord --version | head -n 1)"
supervisord -c /etc/supervisord.conf
sleep 1

echo "[Entrypoint] sing-box version: $(/usr/local/bin/sing-box version | head -n 1)"
echo "[Entrypoint] cerberus-node command: $*"

exec "$@"
