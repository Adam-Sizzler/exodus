#!/usr/bin/env bash
set -Eeuo pipefail

BASELINE_NAME="20260506000000_initial_schema"
LEGACY_MIGRATIONS=(
  "20260302031100_init_schema"
  "20260309163000_probe_users_created_at_index"
  "20260312230000_srs_lists"
  "20260313002000_srs_lists_enabled_drop_type"
  "20260316000100_add_modules_settings"
  "20260316003000_create_modules_settings"
  "20260319170000_rename_xray_columns_to_singbox"
  "20260322030000_hosts_mux_per_core_and_selector_toggle"
  "20260325010000_preserve_json_key_order"
  "20260325223000_subscription_connections"
  "20260326000100_sub_nodes_table"
  "20260326010100_sub_nodes_runtime_columns"
  "20260326023000_sub_nodes_panel_api_token"
  "20260326220000_sub_nodes_subpage_config_uuid"
  "20260326233000_sub_nodes_grpc_auth_token"
  "20260326234000_drop_sub_nodes_panel_api_token"
  "20260327000100_sub_nodes_cleanup_and_keygen_grpc_token"
  "20260327020000_sub_nodes_subpage_join_and_cleanup"
  "20260327030000_drop_sub_nodes_to_subpage_timestamps"
  "20260413190000_rename_cerberus_settings_to_exodus_settings"
  "20260413201500_drop_hosts_xhttp_extra_params"
  "20260419170000_sub_nodes_public_domain"
)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXODUS_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
PROJECT_DIR="${PROJECT_DIR:-$(cd "${EXODUS_DIR}/.." && pwd)}"
COMPOSE_FILE="${COMPOSE_FILE:-${PROJECT_DIR}/docker-compose.yml}"
DB_SERVICE="${DB_SERVICE:-exodus-db}"
BACKUP_DIR="${BACKUP_DIR:-${PROJECT_DIR}/backups}"
MIGRATION_SQL="${EXODUS_DIR}/backend/internal/db/prisma/migrations/${BASELINE_NAME}/migration.sql"

COMPOSE=(docker compose -f "${COMPOSE_FILE}" --project-directory "${PROJECT_DIR}")
DB_USER=""
DB_NAME=""

usage() {
  cat <<EOF
Usage:
  $0 prepare
      Convert an old Exodus/Cerberus DB migration history to the new baseline.

  $0 export [archive.tar.gz]
      Prepare the source DB and create a portable PostgreSQL archive.

  $0 import <archive.tar.gz> [--force]
      Restore the archive into the target compose DB. --force drops public schema first.

Environment:
  PROJECT_DIR=/home/docker/projectSB
  COMPOSE_FILE=/home/docker/projectSB/docker-compose.yml
  DB_SERVICE=exodus-db
  BACKUP_DIR=/home/docker/projectSB/backups
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
  [[ -f "$1" ]] || die "file not found: $1"
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
  require_file "${MIGRATION_SQL}"

  log "Starting ${DB_SERVICE} if needed"
  "${COMPOSE[@]}" up -d "${DB_SERVICE}" >/dev/null

  log "Waiting for PostgreSQL"
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

baseline_checksum() {
  sha256sum "${MIGRATION_SQL}" | awk '{print $1}'
}

contains_item() {
  local needle="$1"
  shift
  local item
  for item in "$@"; do
    [[ "${item}" == "${needle}" ]] && return 0
  done
  return 1
}

backup_legacy_migration_history() {
  mkdir -p "${BACKUP_DIR}"
  local stamp
  stamp="$(date -u +%Y%m%dT%H%M%SZ)"
  local out="${BACKUP_DIR}/schema_migrations-legacy-${stamp}.sql"
  compose_exec pg_dump -U "${DB_USER}" -d "${DB_NAME}" --data-only --table=public.schema_migrations >"${out}"
  log "Legacy migration history backup: ${out}"
}

prepare_db() {
  ensure_db_up

  local checksum
  checksum="$(baseline_checksum)"

  local has_table
  has_table="$(psql_query "select case when to_regclass('public.schema_migrations') is null then 'no' else 'yes' end")"
  [[ "${has_table}" == "yes" ]] || die "public.schema_migrations not found; this does not look like an old Exodus DB"

  mapfile -t applied < <(psql_query "select migration_name from public.schema_migrations order by migration_name" | sed '/^[[:space:]]*$/d')
  [[ "${#applied[@]}" -gt 0 ]] || die "public.schema_migrations is empty"

  if [[ "${#applied[@]}" -eq 1 && "${applied[0]}" == "${BASELINE_NAME}" ]]; then
    psql_run -v baseline_name="${BASELINE_NAME}" -v baseline_checksum="${checksum}" <<'SQL'
UPDATE public.schema_migrations
SET checksum = :'baseline_checksum',
    finished_at = COALESCE(finished_at, now()),
    applied_steps_count = 1,
    applied_at = COALESCE(applied_at, now())
WHERE migration_name = :'baseline_name';
SQL
    log "Database already uses baseline ${BASELINE_NAME}"
    return 0
  fi

  local missing=()
  local migration
  for migration in "${LEGACY_MIGRATIONS[@]}"; do
    if ! contains_item "${migration}" "${applied[@]}"; then
      missing+=("${migration}")
    fi
  done
  if [[ "${#missing[@]}" -gt 0 ]]; then
    printf '%s\n' "${missing[@]}" >&2
    die "old database is not at the supported final legacy migration state"
  fi

  local unknown=()
  for migration in "${applied[@]}"; do
    if [[ "${migration}" != "${BASELINE_NAME}" ]] && ! contains_item "${migration}" "${LEGACY_MIGRATIONS[@]}"; then
      unknown+=("${migration}")
    fi
  done
  if [[ "${#unknown[@]}" -gt 0 ]]; then
    printf '%s\n' "${unknown[@]}" >&2
    die "unknown migration history detected; review it manually before baseline conversion"
  fi

  backup_legacy_migration_history

  log "Converting legacy migration history to baseline ${BASELINE_NAME}"
  psql_run -v baseline_name="${BASELINE_NAME}" -v baseline_checksum="${checksum}" <<'SQL'
BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;

DO $$
DECLARE
  missing_tables text[];
BEGIN
  SELECT ARRAY(
    SELECT table_name
    FROM (VALUES
      ('admin'),
      ('admin_sessions'),
      ('api_tokens'),
      ('config_profiles'),
      ('exodus_settings'),
      ('external_squads'),
      ('hosts'),
      ('hwid_user_devices'),
      ('infra_providers'),
      ('internal_squads'),
      ('keygen'),
      ('modules_settings'),
      ('nodes'),
      ('passkeys'),
      ('srs_lists'),
      ('sub_nodes'),
      ('subscription_page_config'),
      ('subscription_settings'),
      ('subscription_templates'),
      ('users')
    ) AS required(table_name)
    WHERE to_regclass(format('public.%I', table_name)) IS NULL
  ) INTO missing_tables;

  IF array_length(missing_tables, 1) IS NOT NULL THEN
    RAISE EXCEPTION 'required tables are missing: %', array_to_string(missing_tables, ', ');
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'hosts' AND column_name = 'singbox_mux_params'
  ) THEN
    RAISE EXCEPTION 'hosts.singbox_mux_params is missing';
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'hosts' AND column_name = 'clash_mux_params'
  ) THEN
    RAISE EXCEPTION 'hosts.clash_mux_params is missing';
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'hosts' AND column_name = 'selector_nodes_first'
  ) THEN
    RAISE EXCEPTION 'hosts.selector_nodes_first is missing';
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'keygen' AND column_name = 'grpc_auth_token'
  ) THEN
    RAISE EXCEPTION 'keygen.grpc_auth_token is missing';
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'sub_nodes' AND column_name = 'public_domain'
  ) THEN
    RAISE EXCEPTION 'sub_nodes.public_domain is missing';
  END IF;
END $$;

DROP TABLE public.schema_migrations;
CREATE TABLE public.schema_migrations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  migration_name TEXT NOT NULL UNIQUE,
  checksum TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  finished_at TIMESTAMPTZ,
  applied_steps_count INTEGER NOT NULL DEFAULT 1,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO public.schema_migrations (
  migration_name, checksum, started_at, finished_at, applied_steps_count, applied_at
) VALUES (
  :'baseline_name', :'baseline_checksum', now(), now(), 1, now()
);

COMMIT;
SQL

  log "Database is ready for the new baseline"
}

export_db() {
  local out="${1:-}"
  if [[ -z "${out}" ]]; then
    mkdir -p "${BACKUP_DIR}"
    out="${BACKUP_DIR}/exodus-db-$(date -u +%Y%m%dT%H%M%SZ).tar.gz"
  fi
  mkdir -p "$(dirname "${out}")"
  out="$(cd "$(dirname "${out}")" && pwd)/$(basename "${out}")"

  prepare_db

  local work_dir
  work_dir="$(mktemp -d)"

  log "Creating PostgreSQL custom dump"
  compose_exec pg_dump -U "${DB_USER}" -d "${DB_NAME}" -Fc --no-owner --no-acl >"${work_dir}/postgres.dump"

  cat >"${work_dir}/manifest.env" <<EOF
EXODUS_DB_ARCHIVE_VERSION=1
CREATED_AT=$(date -u +%Y-%m-%dT%H:%M:%SZ)
BASELINE_NAME=${BASELINE_NAME}
DB_NAME=${DB_NAME}
DB_SERVICE=${DB_SERVICE}
EOF

  cat >"${work_dir}/README.restore.txt" <<EOF
Restore on the target machine:

  cd /home/docker/projectSB/exodus
  ./scripts/exodus-db-migration.sh import ${out} --force

Then start the panel from /home/docker/projectSB:

  docker compose up -d --build exodus-db exodus-redis exodus
EOF

  tar -czf "${out}" -C "${work_dir}" .
  rm -rf "${work_dir}"
  log "Portable DB archive created: ${out}"
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
  archive="$(cd "$(dirname "${archive}")" && pwd)/$(basename "${archive}")"

  ensure_db_up

  local work_dir
  work_dir="$(mktemp -d)"
  tar -xzf "${archive}" -C "${work_dir}"
  [[ -f "${work_dir}/postgres.dump" ]] || die "archive does not contain postgres.dump"

  local table_count
  table_count="$(psql_query "select count(*) from information_schema.tables where table_schema = 'public' and table_type = 'BASE TABLE'")"
  if [[ "${table_count}" != "0" && "${force}" != "true" ]]; then
    die "target database is not empty (${table_count} public tables); rerun with --force to replace it"
  fi

  if [[ "${force}" == "true" ]]; then
    log "Dropping target public schema"
    psql_run <<'SQL'
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO public;
SQL
  fi

  log "Restoring PostgreSQL dump"
  compose_exec pg_restore -U "${DB_USER}" -d "${DB_NAME}" --no-owner --no-acl --clean --if-exists <"${work_dir}/postgres.dump"

  prepare_db
  rm -rf "${work_dir}"
  log "Import completed"
}

main() {
  local command="${1:-help}"
  shift || true

  case "${command}" in
    prepare)
      prepare_db "$@"
      ;;
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
