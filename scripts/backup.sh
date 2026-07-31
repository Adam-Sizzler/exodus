#!/usr/bin/env bash
set -Eeuo pipefail

# Определяем пути от текущей рабочей директории
PROJECT_DIR="${PROJECT_DIR:-$(pwd)}"
ENV_FILE="${ENV_FILE:-${PROJECT_DIR}/.env}"
DB_SERVICE="${DB_SERVICE:-exodus-db}"
BACKUP_DIR="${BACKUP_DIR:-${PROJECT_DIR}/backups}"

# Поиск файла compose (compose.yml, compose.yaml, docker-compose.yml, docker-compose.yaml)
find_compose_file() {
  if [[ -n "${COMPOSE_FILE:-}" ]]; then
    echo "${COMPOSE_FILE}"
    return 0
  fi

  local candidate
  for candidate in "compose.yml" "compose.yaml" "docker-compose.yml" "docker-compose.yaml"; do
    if [[ -f "${PROJECT_DIR}/${candidate}" ]]; then
      echo "${PROJECT_DIR}/${candidate}"
      return 0
    fi
  done

  echo "${PROJECT_DIR}/compose.yml"
}

COMPOSE_FILE="$(find_compose_file)"
COMPOSE=(docker compose -f "${COMPOSE_FILE}" --project-directory "${PROJECT_DIR}")
DB_USER=""
DB_NAME=""

usage() {
  cat <<EOF
Usage:
  $0 export [archive.tar.gz]
      Create a portable PostgreSQL database backup archive and upload to Telegram.

  $0 import <archive.tar.gz> [--force | --force-drop | -f]
      Restore the database archive into the target compose DB.
      Use --force-drop (or -f) to terminate active connections and drop existing schema before import.

Environment:
  PROJECT_DIR=${PROJECT_DIR}
  COMPOSE_FILE=${COMPOSE_FILE}
  ENV_FILE=${ENV_FILE}
  DB_SERVICE=${DB_SERVICE}
  BACKUP_DIR=${BACKUP_DIR}
EOF
}

log() {
  printf '[exodus-db] %s\n' "$*"
}

die() {
  printf '[exodus-db] ERROR: %s\n' "$*" >&2
  exit 1
}

require_file() {
  [[ -f "$1" ]] || die "compose file not found in ${PROJECT_DIR} (checked compose.yml / docker-compose.yml)"
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || die "required command is missing: $1"
}

compose_exec() {
  "${COMPOSE[@]}" exec -T "${DB_SERVICE}" "$@"
}

ensure_db_up() {
  require_command docker
  require_file "${COMPOSE_FILE}"

  log "Using compose file: ${COMPOSE_FILE}"
  log "Starting ${DB_SERVICE} if needed..."
  "${COMPOSE[@]}" up -d "${DB_SERVICE}" >/dev/null

  log "Waiting for PostgreSQL..."
  for _ in $(seq 1 90); do
    if compose_exec sh -lc 'pg_isready -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-postgres}"' >/dev/null 2>&1; then
      DB_USER="$(compose_exec sh -lc 'printf "%s" "${POSTGRES_USER:-postgres}"')"
      DB_NAME="$(compose_exec sh -lc 'printf "%s" "${POSTGRES_DB:-postgres}"')"
      return 0
    fi
    sleep 1
  done

  die "PostgreSQL is not ready in service ${DB_SERVICE}"
}

psql_query() {
  compose_exec psql -X -v ON_ERROR_STOP=1 -U "${DB_USER}" -d "${DB_NAME}" -tAc "$1"
}

psql_run() {
  compose_exec psql -X -v ON_ERROR_STOP=1 -U "${DB_USER}" -d "${DB_NAME}" "$@"
}

send_to_telegram() {
  local archive_path="$1"

  if [[ ! -f "${ENV_FILE}" ]]; then
    log "Notice: .env file not found at ${ENV_FILE}, skipping Telegram send."
    return 0
  fi

  local is_enabled bot_token notify_service
  is_enabled=$(grep -E '^IS_TELEGRAM_NOTIFICATIONS_ENABLED=' "${ENV_FILE}" | cut -d'=' -f2- | tr -d '"'\''\r' || true)
  bot_token=$(grep -E '^TELEGRAM_BOT_TOKEN=' "${ENV_FILE}" | cut -d'=' -f2- | tr -d '"'\''\r' || true)
  notify_service=$(grep -E '^TELEGRAM_NOTIFY_SERVICE=' "${ENV_FILE}" | cut -d'=' -f2- | tr -d '"'\''\r' || true)

  if [[ "${is_enabled}" != "true" ]]; then
    log "Telegram notifications disabled in .env (IS_TELEGRAM_NOTIFICATIONS_ENABLED!=true)"
    return 0
  fi

  if [[ -z "${bot_token}" || -z "${notify_service}" ]]; then
    log "Telegram bot token or TELEGRAM_NOTIFY_SERVICE missing in .env"
    return 0
  fi

  log "Uploading backup archive to Telegram..."

  local chat_id thread_id
  chat_id="${notify_service%%:*}"
  thread_id="${notify_service#*:}"

  if [[ "${thread_id}" == "${notify_service}" ]]; then
    thread_id=""
  fi

  local filename
  filename="$(basename "${archive_path}")"

  local curl_cmd=(
    curl -s -S -X POST "https://api.telegram.org/bot${bot_token}/sendDocument"
    -F "chat_id=${chat_id}"
    -F "document=@${archive_path}"
    -F "caption=📦 *Exodus DB Backup*
📄 \`${filename}\`
🕒 \`$(date -u '+%Y-%m-%d %H:%M:%S UTC')\`"
    -F "parse_mode=Markdown"
  )

  if [[ -n "${thread_id}" ]]; then
    curl_cmd+=(-F "message_thread_id=${thread_id}")
  fi

  local response
  if response=$("${curl_cmd[@]}"); then
    if echo "${response}" | grep -q '"ok":true'; then
      log "Backup archive successfully uploaded to Telegram!"
    else
      log "ERROR sending to Telegram: ${response}"
    fi
  else
    log "ERROR: Failed to connect to Telegram API."
  fi
}

export_db() {
  ensure_db_up

  local out="${1:-}"
  if [[ -z "${out}" ]]; then
    mkdir -p "${BACKUP_DIR}"
    out="${BACKUP_DIR}/exodus-db-$(date -u +%Y%m%dT%H%M%SZ).tar.gz"
  else
    mkdir -p "$(dirname "${out}")"
  fi

  local work_dir
  work_dir="$(mktemp -d)"

  log "Creating PostgreSQL custom dump..."
  compose_exec pg_dump -U "${DB_USER}" -d "${DB_NAME}" -Fc --no-owner --no-acl >"${work_dir}/postgres.dump"

  cat >"${work_dir}/manifest.env" <<EOF
EXODUS_DB_ARCHIVE_VERSION=1
CREATED_AT=$(date -u +%Y-%m-%dT%H:%M:%SZ)
DB_NAME=${DB_NAME}
DB_SERVICE=${DB_SERVICE}
EOF

  cat >"${work_dir}/README.restore.txt" <<EOF
Restore on the target machine:

  $0 import ${out} --force-drop

Then restart the services:

  docker compose up -d --build exodus-db exodus-redis exodus
EOF

  tar -czf "${out}" -C "${work_dir}" .
  rm -rf "${work_dir}"
  log "Backup archive created: ${out}"

  send_to_telegram "${out}"
}

import_db() {
  local archive=""
  local force="false"

  # Парсим аргументы
  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      --force|--force-drop|-f)
        force="true"
        ;;
      -*)
        die "unknown import option: $1"
        ;;
      *)
        if [[ -z "${archive}" ]]; then
          archive="$1"
        else
          die "unexpected argument: $1"
        fi
        ;;
    esac
    shift
  done

  [[ -n "${archive}" ]] || die "archive path is required. Usage: $0 import <path_to_archive.tar.gz> [--force-drop]"
  [[ -f "${archive}" ]] || die "archive not found: ${archive}"

  ensure_db_up

  local work_dir
  work_dir="$(mktemp -d)"
  tar -xzf "${archive}" -C "${work_dir}"
  [[ -f "${work_dir}/postgres.dump" ]] || die "archive does not contain postgres.dump"

  local table_count
  table_count="$(psql_query "select count(*) from information_schema.tables where table_schema = 'public' and table_type = 'BASE TABLE'")"

  if [[ "${table_count}" != "0" && "${force}" != "true" ]]; then
    rm -rf "${work_dir}"
    die "target database is not empty (${table_count} public tables); rerun with --force-drop (or -f) to replace it"
  fi

  if [[ "${force}" == "true" ]]; then
    log "Terminating active database connections and dropping public schema..."
    psql_run <<'SQL'
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = current_database()
  AND pid <> pg_backend_pid();

DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO public;
SQL
  fi

  log "Restoring PostgreSQL dump..."
  compose_exec pg_restore -U "${DB_USER}" -d "${DB_NAME}" --no-owner --no-acl --clean --if-exists <"${work_dir}/postgres.dump"

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
