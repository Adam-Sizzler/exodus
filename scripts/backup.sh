#!/usr/bin/env bash
set -Eeuo pipefail

# Берем текущую рабочую директорию (откуда запущен скрипт)
PROJECT_DIR="${PROJECT_DIR:-$(pwd)}"

# Если COMPOSE_FILE не задан явно, ищем стандартные имена.
if [[ -z "${COMPOSE_FILE:-}" ]]; then
  if [[ -f "${PROJECT_DIR}/compose.yml" ]]; then
    COMPOSE_FILE="${PROJECT_DIR}/compose.yml"
  elif [[ -f "${PROJECT_DIR}/docker-compose.yml" ]]; then
    COMPOSE_FILE="${PROJECT_DIR}/docker-compose.yml"
  else
    COMPOSE_FILE="${PROJECT_DIR}/compose.yml"
  fi
fi

DB_SERVICE="${DB_SERVICE:-exodus-db}"
BACKUP_DIR="${BACKUP_DIR:-${PROJECT_DIR}/backups}"

COMPOSE_OVERRIDE_FILE=""
COMPOSE=()

DB_USER=""
DB_NAME=""
DB_HOST=""
DB_PORT=""
DB_PASSWORD=""
DB_CONNECTION_MODE=""

TELEGRAM_BOT_TOKEN=""
TELEGRAM_NOTIFY_SERVICE=""

WORK_DIRS=()

usage() {
  cat <<EOF
Usage:
  $0 export [archive.tar.gz]
      Create a portable PostgreSQL database backup archive
      and send it to Telegram.

  $0 import <archive.tar.gz> [--force]
      Restore the database archive into the target compose DB.
      Use --force to drop existing public schema before import.

Environment:
  PROJECT_DIR=${PROJECT_DIR}
  COMPOSE_FILE=${COMPOSE_FILE}
  DB_SERVICE=${DB_SERVICE}
  BACKUP_DIR=${BACKUP_DIR}
EOF
}

log() {
  printf '[exodus-db] %s\n' "$*"
}

warn() {
  printf '[exodus-db] WARN: %s\n' "$*" >&2
}

die() {
  printf '[exodus-db] ERROR: %s\n' "$*" >&2
  exit 1
}

require_file() {
  [[ -f "$1" ]] || die "file not found: $1"
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || die "required command is missing: $1"
}

cleanup() {
  local dir

  if [[ -n "${COMPOSE_OVERRIDE_FILE:-}" && -f "${COMPOSE_OVERRIDE_FILE}" ]]; then
    rm -f "${COMPOSE_OVERRIDE_FILE}"
  fi

  for dir in "${WORK_DIRS[@]:-}"; do
    if [[ -n "${dir}" && -d "${dir}" ]]; then
      rm -rf "${dir}"
    fi
  done
}

trap cleanup EXIT

prepare_compose() {
  require_file "${COMPOSE_FILE}"

  COMPOSE_OVERRIDE_FILE="$(mktemp "${TMPDIR:-/tmp}/exodus-db-compose.XXXXXX.yml")"

  cat >"${COMPOSE_OVERRIDE_FILE}" <<'YAML'
services: {}

volumes:
  postgres-socket:
    name: postgres-socket
    driver: local

  valkey-socket:
    name: valkey-socket
    driver: local
YAML

  COMPOSE=(
    docker
    compose
    -f "${COMPOSE_FILE}"
    -f "${COMPOSE_OVERRIDE_FILE}"
    --project-directory "${PROJECT_DIR}"
  )
}

compose_exec() {
  "${COMPOSE[@]}" exec -T "${DB_SERVICE}" "$@"
}

db_exec() {
  if [[ -n "${DB_PASSWORD}" ]]; then
    compose_exec env "PGPASSWORD=${DB_PASSWORD}" "$@"
  else
    compose_exec "$@"
  fi
}

db_psql_args() {
  local -n _args="$1"

  _args=(
    psql
    -X
    -v
    ON_ERROR_STOP=1
  )

  if [[ -n "${DB_HOST}" ]]; then
    _args+=(-h "${DB_HOST}")
  fi

  if [[ -n "${DB_PORT}" ]]; then
    _args+=(-p "${DB_PORT}")
  fi

  _args+=(
    -U "${DB_USER}"
    -d "${DB_NAME}"
  )
}

db_pg_dump_args() {
  local -n _args="$1"

  _args=(
    pg_dump
  )

  if [[ -n "${DB_HOST}" ]]; then
    _args+=(-h "${DB_HOST}")
  fi

  if [[ -n "${DB_PORT}" ]]; then
    _args+=(-p "${DB_PORT}")
  fi

  _args+=(
    -U "${DB_USER}"
    -d "${DB_NAME}"
    -Fc
    --no-owner
    --no-acl
  )
}

db_pg_restore_args() {
  local -n _args="$1"

  _args=(
    pg_restore
  )

  if [[ -n "${DB_HOST}" ]]; then
    _args+=(-h "${DB_HOST}")
  fi

  if [[ -n "${DB_PORT}" ]]; then
    _args+=(-p "${DB_PORT}")
  fi

  _args+=(
    -U "${DB_USER}"
    -d "${DB_NAME}"
    --no-owner
    --no-acl
    --clean
    --if-exists
  )
}

load_db_config() {
  local database_url
  local postgres_user
  local postgres_db
  local postgres_password

  database_url="$(compose_exec sh -lc 'printf "%s" "${DATABASE_URL:-}"')"
  postgres_user="$(compose_exec sh -lc 'printf "%s" "${POSTGRES_USER:-}"')"
  postgres_db="$(compose_exec sh -lc 'printf "%s" "${POSTGRES_DB:-}"')"
  postgres_password="$(compose_exec sh -lc 'printf "%s" "${POSTGRES_PASSWORD:-}"')"

  [[ -n "${database_url}" ]] || die \
    "DATABASE_URL is not set in ${DB_SERVICE}"

  case "${database_url}" in
    postgresql://*|postgres://*)
      ;;
    *)
      die "unsupported DATABASE_URL scheme; expected postgresql:// or postgres://"
      ;;
  esac

  local url_without_scheme
  local authority
  local path_and_query
  local userinfo
  local hostport

  url_without_scheme="${database_url#*://}"

  if [[ "${url_without_scheme}" == */* ]]; then
    authority="${url_without_scheme%%/*}"
    path_and_query="${url_without_scheme#*/}"
  else
    authority="${url_without_scheme}"
    path_and_query=""
  fi

  if [[ "${authority}" == *"@"* ]]; then
    userinfo="${authority%@*}"
    hostport="${authority#*@}"
  else
    userinfo=""
    hostport="${authority}"
  fi

  #
  # Пользователь и пароль из DATABASE_URL.
  #
  if [[ -n "${userinfo}" ]]; then
    if [[ "${userinfo}" == *":"* ]]; then
      DB_USER="${userinfo%%:*}"
      DB_PASSWORD="${userinfo#*:}"
    else
      DB_USER="${userinfo}"
    fi
  fi

  #
  # Явные POSTGRES_* имеют приоритет.
  #
  if [[ -n "${postgres_user}" ]]; then
    DB_USER="${postgres_user}"
  fi

  if [[ -n "${postgres_password}" ]]; then
    DB_PASSWORD="${postgres_password}"
  fi

  #
  # Имя базы.
  #
  if [[ -n "${postgres_db}" ]]; then
    DB_NAME="${postgres_db}"
  else
    DB_NAME="${path_and_query%%\?*}"
    DB_NAME="${DB_NAME%%#*}"
  fi

  DB_USER="${DB_USER:-postgres}"
  DB_NAME="${DB_NAME:-postgres}"

  #
  # Определяем TCP или Unix socket.
  #
  if [[ -z "${hostport}" ]]; then
    DB_CONNECTION_MODE="socket"
    DB_HOST=""
    DB_PORT=""
  else
    DB_CONNECTION_MODE="tcp"

    #
    # IPv6:
    #   [::1]:5432
    #
    if [[ "${hostport}" == \[*\]* ]]; then
      DB_HOST="${hostport%%]*}"
      DB_HOST="${DB_HOST#\[}"

      if [[ "${hostport}" == *"]:"* ]]; then
        DB_PORT="${hostport##*:}"
      else
        DB_PORT="5432"
      fi

    elif [[ "${hostport}" == *":"* ]]; then
      DB_HOST="${hostport%:*}"
      DB_PORT="${hostport##*:}"

    else
      DB_HOST="${hostport}"
      DB_PORT="5432"
    fi

    [[ -n "${DB_HOST}" ]] || die \
      "could not determine PostgreSQL host from DATABASE_URL"

    [[ -n "${DB_PORT}" ]] || DB_PORT="5432"
  fi

  log "PostgreSQL connection mode: ${DB_CONNECTION_MODE}"

  if [[ "${DB_CONNECTION_MODE}" == "tcp" ]]; then
    log "PostgreSQL endpoint: ${DB_HOST}:${DB_PORT}"
  else
    log "PostgreSQL endpoint: Unix socket"
  fi

  log "PostgreSQL database: ${DB_NAME}"
  log "PostgreSQL user: ${DB_USER}"
}

load_telegram_config() {
  TELEGRAM_BOT_TOKEN="$(
    compose_exec sh -lc 'printf "%s" "${TELEGRAM_BOT_TOKEN:-}"' 2>/dev/null || true
  )"

  TELEGRAM_NOTIFY_SERVICE="$(
    compose_exec sh -lc 'printf "%s" "${TELEGRAM_NOTIFY_SERVICE:-}"' 2>/dev/null || true
  )"

  if [[ -z "${TELEGRAM_BOT_TOKEN}" || -z "${TELEGRAM_NOTIFY_SERVICE}" ]]; then
    return 1
  fi
  return 0
}

parse_telegram_destination() {
  local destination="$1"

  TELEGRAM_CHAT_ID=""
  TELEGRAM_THREAD_ID=""

  if [[ "${destination}" == *":"* ]]; then
    TELEGRAM_CHAT_ID="${destination%%:*}"
    TELEGRAM_THREAD_ID="${destination#*:}"
  else
    TELEGRAM_CHAT_ID="${destination}"
  fi

  if [[ -z "${TELEGRAM_CHAT_ID}" ]]; then
    warn "invalid TELEGRAM_NOTIFY_SERVICE: chat_id is empty"
    return 1
  fi

  if [[ -n "${TELEGRAM_THREAD_ID}" ]]; then
    if [[ ! "${TELEGRAM_THREAD_ID}" =~ ^[0-9]+$ ]]; then
      warn "invalid Telegram thread_id: ${TELEGRAM_THREAD_ID}"
      return 1
    fi
  fi
  return 0
}

send_to_telegram() {
  local archive="$1"

  if ! command -v curl >/dev/null 2>&1; then
    warn "curl command is missing, skipping Telegram upload"
    return 0
  fi

  if ! load_telegram_config; then
    warn "Telegram configuration (TELEGRAM_BOT_TOKEN / TELEGRAM_NOTIFY_SERVICE) not set, skipping Telegram upload"
    return 0
  fi

  if ! parse_telegram_destination "${TELEGRAM_NOTIFY_SERVICE}"; then
    warn "failed to parse Telegram destination, skipping Telegram upload"
    return 0
  fi

  local api_url
  api_url="https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendDocument"

  log "Sending backup archive to Telegram..."

  local response
  local curl_args=(
    --silent
    --show-error
    --max-time 300
    -X POST
    "${api_url}"
    -F "chat_id=${TELEGRAM_CHAT_ID}"
    -F "document=@${archive}"
    -F "caption=Exodus PostgreSQL backup: $(basename "${archive}")"
  )

  if [[ -n "${TELEGRAM_THREAD_ID}" ]]; then
    curl_args+=(
      -F "message_thread_id=${TELEGRAM_THREAD_ID}"
    )
  fi

  if ! response="$(curl "${curl_args[@]}" 2>&1)"; then
    warn "failed to send backup archive to Telegram: ${response}"
    return 0
  fi

  if ! printf '%s' "${response}" | grep -q '"ok":true'; then
    warn "Telegram rejected the backup archive: ${response}"
    return 0
  fi

  log "Backup archive sent to Telegram"
}

ensure_db_up() {
  require_command docker
  require_file "${COMPOSE_FILE}"

  prepare_compose

  log "Using compose file: ${COMPOSE_FILE}"

  log "Starting ${DB_SERVICE} if needed..."

  "${COMPOSE[@]}" up -d "${DB_SERVICE}" >/dev/null

  load_db_config

  log "Waiting for PostgreSQL..."

  local psql_args=()
  db_psql_args psql_args

  for _ in $(seq 1 90); do
    if db_exec "${psql_args[@]}" -tAc 'select 1' >/dev/null 2>&1; then
      return 0
    fi

    sleep 1
  done

  die "PostgreSQL is not ready in service ${DB_SERVICE}"
}

psql_query() {
  local query="$1"
  local psql_args=()

  db_psql_args psql_args

  db_exec "${psql_args[@]}" -tAc "${query}"
}

psql_run() {
  local psql_args=()

  db_psql_args psql_args

  db_exec "${psql_args[@]}" "$@"
}

export_db() {
  local out="${1:-}"

  if [[ -z "${out}" ]]; then
    mkdir -p "${BACKUP_DIR}"
    out="${BACKUP_DIR}/exodus-db-$(date -u +%Y%m%dT%H%M%SZ).tar.gz"
  else
    mkdir -p "$(dirname "${out}")"
  fi

  ensure_db_up

  local work_dir
  work_dir="$(mktemp -d)"
  WORK_DIRS+=("${work_dir}")

  log "Creating PostgreSQL custom dump..."

  local pg_dump_args=()
  db_pg_dump_args pg_dump_args

  db_exec "${pg_dump_args[@]}" >"${work_dir}/postgres.dump"

  cat >"${work_dir}/manifest.env" <<EOF
EXODUS_DB_ARCHIVE_VERSION=1
CREATED_AT=$(date -u +%Y-%m-%dT%H:%M:%SZ)
DB_NAME=${DB_NAME}
DB_SERVICE=${DB_SERVICE}
DB_CONNECTION_MODE=${DB_CONNECTION_MODE}
EOF

  cat >"${work_dir}/README.restore.txt" <<EOF
Restore on the target machine:

  $0 import ${out} --force

The script will automatically determine TCP or Unix socket
connection mode from DATABASE_URL in .env.

Then restart the services:

  docker compose up -d --build exodus-db exodus-redis exodus
EOF

  tar -czf "${out}" -C "${work_dir}" .

  rm -rf "${work_dir}"

  log "Backup archive created: ${out}"

  send_to_telegram "${out}"
}

import_db() {
  local archive="${1:-}"

  [[ -n "${archive}" ]] || die "archive path is required"

  shift || true

  local force="false"

  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      --force)
        force="true"
        ;;
      *)
        die "unknown import option: $1"
        ;;
    esac

    shift
  done

  [[ -f "${archive}" ]] || die "archive not found: ${archive}"

  ensure_db_up

  local work_dir
  work_dir="$(mktemp -d)"
  WORK_DIRS+=("${work_dir}")

  tar -xzf "${archive}" -C "${work_dir}"

  [[ -f "${work_dir}/postgres.dump" ]] || \
    die "archive does not contain postgres.dump"

  local table_count

  table_count="$(
    psql_query \
      "select count(*) from information_schema.tables where table_schema = 'public' and table_type = 'BASE TABLE'"
  )"

  if [[ "${table_count}" != "0" && "${force}" != "true" ]]; then
    die "target database is not empty (${table_count} public tables); rerun with --force to replace it"
  fi

  if [[ "${force}" == "true" ]]; then
    log "Dropping target public schema..."

    psql_run <<'SQL'
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO public;
SQL
  fi

  log "Restoring PostgreSQL dump..."

  local pg_restore_args=()
  db_pg_restore_args pg_restore_args

  if ! db_exec "${pg_restore_args[@]}" <"${work_dir}/postgres.dump"; then
    warn "pg_restore completed with warnings or notices (check output above)"
  fi

  rm -rf "${work_dir}"

  log "Import completed successfully"
}

main() {
  local command="${1:-help}"

  shift || true

  case "${command}" in
    export)
      export_db "$@"
      ;;
    import)
      import_db "$@"
      ;;
    help|-h|--help)
      usage
      ;;
    *)
      usage >&2
      die "unknown command: ${command}"
      ;;
  esac
}

main "$@"
