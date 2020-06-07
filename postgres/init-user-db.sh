#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER "$PG_APP_USER" WITH PASSWORD '$PG_APP_PASSWORD';
    CREATE DATABASE "$PG_APP_DB";
    GRANT ALL ON DATABASE "$PG_APP_DB" TO "$PG_APP_USER";
EOSQL

psql -v ON_ERROR_STOP=1 \
  --username "$PG_APP_USER" \
  --dbname "$PG_APP_DB" \
  -f /docker-entrypoint-initdb.d/schema.sql