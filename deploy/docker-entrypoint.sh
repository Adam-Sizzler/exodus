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

installed_singbox_version() {
    if [ ! -x /usr/local/bin/sing-box ]; then
        return 0
    fi

    local line version
    line="$(/usr/local/bin/sing-box version 2>/dev/null | head -n 1 || true)"
    version="$(printf '%s' "$line" | awk '{print $NF}')"
    normalize_version "$version"
}

detect_singbox_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            printf '%s' "amd64"
            ;;
        aarch64|arm64)
            printf '%s' "arm64"
            ;;
        *)
            return 1
            ;;
    esac
}

ensure_singbox_version() {
    local desired desired_tag desired_norm current_norm sb_arch url tmp_file validate_output
    local current_output has_working_current

    current_output="$(/usr/local/bin/sing-box version 2>&1 || true)"
    has_working_current=0
    if printf '%s\n' "$current_output" | grep -qi '^sing-box version'; then
        has_working_current=1
    fi

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

    if ! sb_arch="$(detect_singbox_arch)"; then
        if [ "$has_working_current" -eq 1 ]; then
            echo "[Entrypoint] Unsupported architecture for sing-box auto-download, continuing with current sing-box: $(uname -m)" >&2
            return 0
        fi
        echo "[Entrypoint] Unsupported architecture for sing-box auto-download: $(uname -m)" >&2
        exit 1
    fi

    url="https://github.com/Adam-Sizzler/sing-box-v2ray-api/releases/download/${desired_tag}/sing-box-linux-${sb_arch}"
    tmp_file="/usr/local/bin/.sing-box-${desired_tag}-${sb_arch}.tmp"

    echo "[Entrypoint] Downloading sing-box ${desired_tag} for ${sb_arch}"
    if ! curl -fL "$url" -o "$tmp_file"; then
        rm -f "$tmp_file"
        if [ "$has_working_current" -eq 1 ]; then
            echo "[Entrypoint] Failed to download ${desired_tag}, continuing with current sing-box: $(printf '%s\n' "$current_output" | head -n 1)" >&2
            return 0
        fi
        echo "[Entrypoint] Failed to download ${desired_tag} and no working current sing-box is available" >&2
        exit 1
    fi
    chmod +x "$tmp_file"

    validate_output="$("$tmp_file" version 2>&1 || true)"
    if [ -z "$validate_output" ]; then
        if [ "$has_working_current" -eq 1 ]; then
            echo "[Entrypoint] Downloaded sing-box is invalid (empty output), continuing with current sing-box: $(printf '%s\n' "$current_output" | head -n 1)" >&2
            rm -f "$tmp_file"
            return 0
        fi
        echo "[Entrypoint] Downloaded sing-box is invalid (empty output): ${url}" >&2
        rm -f "$tmp_file"
        exit 1
    fi
    if ! printf '%s\n' "$validate_output" | grep -qi '^sing-box version'; then
        if [ "$has_working_current" -eq 1 ]; then
            echo "[Entrypoint] Downloaded sing-box failed validation, continuing with current sing-box: $(printf '%s\n' "$current_output" | head -n 1)" >&2
            echo "[Entrypoint] Validation output: ${validate_output}" >&2
            rm -f "$tmp_file"
            return 0
        fi
        echo "[Entrypoint] Downloaded sing-box failed validation: ${url}" >&2
        echo "[Entrypoint] Validation output: ${validate_output}" >&2
        rm -f "$tmp_file"
        exit 1
    fi

    if ! mv "$tmp_file" /usr/local/bin/sing-box; then
        rm -f "$tmp_file"
        if [ "$has_working_current" -eq 1 ]; then
            echo "[Entrypoint] Failed to replace sing-box binary, continuing with current sing-box: $(printf '%s\n' "$current_output" | head -n 1)" >&2
            return 0
        fi
        echo "[Entrypoint] Failed to replace sing-box binary and no working current sing-box is available" >&2
        exit 1
    fi
    echo "[Entrypoint] Updated sing-box to: $(printf '%s\n' "$validate_output" | head -n 1)"
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
echo "[Entrypoint] exodus-node command: $*"

exec "$@"
