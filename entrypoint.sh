#!/bin/sh

echo "[INIT] Fixing permissions for /app/logs..."
chown -R "$(id -u)" /app/logs 2>/dev/null || true

exec "$@"
