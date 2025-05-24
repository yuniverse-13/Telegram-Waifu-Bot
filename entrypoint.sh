#!/bin/sh
set -e

echo "Running database migrations..."

if [ -z "$DATABASE_URL" ]; then
  echo "Error: DATABASE_URL environment variable is not set."
  exit 1
fi

migrate -database "$DATABASE_URL" -path ./migrations up

echo "Migrations complete."

exec "$@"