@echo off
echo Resetting database...

echo Running migrations down...
go run cmd/migrate/main.go down
if %ERRORLEVEL% NEQ 0 (
    echo Failed to run migrations down
    exit /b %ERRORLEVEL%
)

echo Running migrations up...
go run cmd/migrate/main.go up
if %ERRORLEVEL% NEQ 0 (
    echo Failed to run migrations up
    exit /b %ERRORLEVEL%
)

echo Seeding database...
go run cmd/seed/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to seed database
    exit /b %ERRORLEVEL%
)

echo Database reset complete.
