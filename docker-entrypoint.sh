#!/bin/sh
set -e

echo "Starting entrypoint script..."
export INTERNAL_JWT_SECRET=$(hexdump -vn 64 -e '1/1 "%02x"' /dev/urandom 2>/dev/null || od -An -N64 -tx1 /dev/urandom | tr -d ' \n')


echo "Entrypoint script completed."
exec "$@"
