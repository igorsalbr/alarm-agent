#!/bin/bash

set -e

MIGRATE_PATH="./db/migrations"
DB_URL="${POSTGRES_DSN:-postgres://alarm_user:alarm_pass@localhost:5432/alarm_agent?sslmode=disable}"

command="$1"

case $command in
  "up")
    echo "Applying migrations..."
    migrate -path $MIGRATE_PATH -database "$DB_URL" up
    echo "Migrations applied successfully!"
    ;;
  "down")
    echo "Rolling back migrations..."
    migrate -path $MIGRATE_PATH -database "$DB_URL" down 1
    echo "Migration rolled back successfully!"
    ;;
  "force")
    version="$2"
    if [ -z "$version" ]; then
      echo "Usage: $0 force <version>"
      exit 1
    fi
    echo "Forcing migration to version $version..."
    migrate -path $MIGRATE_PATH -database "$DB_URL" force $version
    echo "Migration forced to version $version!"
    ;;
  "version")
    echo "Current migration version:"
    migrate -path $MIGRATE_PATH -database "$DB_URL" version
    ;;
  "create")
    name="$2"
    if [ -z "$name" ]; then
      echo "Usage: $0 create <migration_name>"
      exit 1
    fi
    echo "Creating new migration: $name"
    migrate create -ext sql -dir $MIGRATE_PATH -seq $name
    echo "Migration files created successfully!"
    ;;
  *)
    echo "Usage: $0 {up|down|force <version>|version|create <name>}"
    echo ""
    echo "Commands:"
    echo "  up           Apply all pending migrations"
    echo "  down         Rollback one migration"
    echo "  force <ver>  Force migration to specific version"
    echo "  version      Show current migration version"
    echo "  create <name> Create new migration files"
    exit 1
    ;;
esac