#!/bin/bash
set -e

echo "Resetting database..."

echo "Running migrations down..."
go run cmd/migrate/main.go down

echo "Running migrations up..."
go run cmd/migrate/main.go up

echo "Seeding database..."
go run cmd/seed/main.go

echo "Database reset complete."
