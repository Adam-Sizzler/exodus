#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${EXODUS_BASE_URL:-https://localhost}"
TOKEN="${EXODUS_AUTH_TOKEN:-}"
COUNT="${BENCH_USERS:-1000}"
CONCURRENCY="${BENCH_CONCURRENCY:-32}"
READ_REQUESTS="${BENCH_READ_REQUESTS:-200}"
DELETE_CHUNK_SIZE="${BENCH_DELETE_CHUNK_SIZE:-500}"
DELETE_CONCURRENCY="${BENCH_DELETE_CONCURRENCY:-4}"
RUN_ID="${BENCH_RUN_ID:-$(date +%s)}"
RESULT_DIR="${BENCH_RESULT_DIR:-.bench/results}"
INSECURE="${BENCH_INSECURE:-1}"

if [[ -z "$TOKEN" ]]; then
    echo "EXODUS_AUTH_TOKEN is required" >&2
    exit 2
fi

if ! command -v curl >/dev/null 2>&1; then
    echo "curl is required" >&2
    exit 2
fi

if ! command -v jq >/dev/null 2>&1; then
    echo "jq is required" >&2
    exit 2
fi

mkdir -p "$RESULT_DIR"

uuid_file="$(mktemp)"
delete_dir="$(mktemp -d)"
trap 'rm -f "$uuid_file"; rm -rf "$delete_dir"' EXIT

now_ms() {
    date +%s%3N
}

elapsed_ms() {
    local start="$1"
    local end
    end="$(now_ms)"
    echo $((end - start))
}

create_one() {
    local i="$1"
    local username="b${RUN_ID}_${i}"
    local expire_at="2030-01-01T00:00:00Z"
    local payload
    local tls_flags=()
    if [[ "${INSECURE}" == "1" ]]; then
        tls_flags=(-k)
    fi

    payload="$(jq -nc --arg username "$username" --arg expireAt "$expire_at" '{
        username: $username,
        status: "DISABLED",
        expireAt: $expireAt,
        trafficLimitBytes: 0,
        trafficLimitStrategy: "NO_RESET"
    }')"

    curl -fsS -g "${tls_flags[@]}" \
        -H "Authorization: Bearer ${TOKEN}" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -X POST "${BASE_URL}/api/users" \
        --data "$payload" | jq -r '.response.uuid'
}

read_one() {
    local i="$1"
    local start=$(( (i - 1) % 50 ))
    local tls_flags=()
    if [[ "${INSECURE}" == "1" ]]; then
        tls_flags=(-k)
    fi

    curl -fsS -g "${tls_flags[@]}" \
        -H "Authorization: Bearer ${TOKEN}" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        "${BASE_URL}/api/users?start=${start}&size=25&filters=[]&filterModes={}&sorting=[]" >/dev/null
}

delete_chunk() {
    local file="$1"
    local payload
    local tls_flags=()
    if [[ "${INSECURE}" == "1" ]]; then
        tls_flags=(-k)
    fi

    payload="$(jq -R -s '{uuids: (split("\n") | map(select(length > 0)))}' "$file")"
    curl -fsS -g "${tls_flags[@]}" \
        -H "Authorization: Bearer ${TOKEN}" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -X POST "${BASE_URL}/api/users/bulk/delete" \
        --data "$payload" >/dev/null
}

export BASE_URL TOKEN RUN_ID INSECURE
export -f create_one read_one delete_chunk now_ms elapsed_ms

echo "benchmark run_id=${RUN_ID} count=${COUNT} concurrency=${CONCURRENCY} base_url=${BASE_URL}" >&2

start_total="$(now_ms)"

start_create="$(now_ms)"
seq 1 "$COUNT" | xargs -n 1 -P "$CONCURRENCY" bash -c 'create_one "$1"' _ > "$uuid_file"
create_ms="$(elapsed_ms "$start_create")"
created_count="$(grep -c '^[0-9a-fA-F-]\{36\}$' "$uuid_file" || true)"

start_read="$(now_ms)"
seq 1 "$READ_REQUESTS" | xargs -n 1 -P "$CONCURRENCY" bash -c 'read_one "$1"' _
read_ms="$(elapsed_ms "$start_read")"

split -l "$DELETE_CHUNK_SIZE" "$uuid_file" "${delete_dir}/chunk_"

start_delete="$(now_ms)"
find "$delete_dir" -type f -name 'chunk_*' -print0 \
    | xargs -0 -n 1 -P "$DELETE_CONCURRENCY" bash -c 'delete_chunk "$1"' _
delete_ms="$(elapsed_ms "$start_delete")"

total_ms="$(elapsed_ms "$start_total")"
result_file="${RESULT_DIR}/api-load-${RUN_ID}.json"

jq -nc \
    --arg runId "$RUN_ID" \
    --arg baseUrl "$BASE_URL" \
    --argjson requested "$COUNT" \
    --argjson created "$created_count" \
    --argjson concurrency "$CONCURRENCY" \
    --argjson readRequests "$READ_REQUESTS" \
    --argjson createMs "$create_ms" \
    --argjson readMs "$read_ms" \
    --argjson deleteMs "$delete_ms" \
    --argjson totalMs "$total_ms" \
    '{
        runId: $runId,
        baseUrl: $baseUrl,
        requestedUsers: $requested,
        createdUsers: $created,
        concurrency: $concurrency,
        readRequests: $readRequests,
        timingsMs: {
            create: $createMs,
            read: $readMs,
            delete: $deleteMs,
            total: $totalMs
        },
        ratesPerSecond: {
            create: (if $createMs > 0 then ($created * 1000 / $createMs) else 0 end),
            read: (if $readMs > 0 then ($readRequests * 1000 / $readMs) else 0 end),
            delete: (if $deleteMs > 0 then ($created * 1000 / $deleteMs) else 0 end)
        }
    }' | tee "$result_file"

echo "result saved to ${result_file}" >&2
