package config

import (
	"os"
	"testing"
)

func TestGetDSN(t *testing.T) {
	// Test case 1: DATABASE_URL is set
	t.Run("DATABASE_URL set", func(t *testing.T) {
		expected := "postgres://user:pass@host:5432/db?sslmode=require"
		os.Setenv("DATABASE_URL", expected)
		defer os.Unsetenv("DATABASE_URL")

		cfg := &Config{
			Database: DatabaseConfig{
				URL: expected,
			},
		}

		if dsn := cfg.GetDSN(); dsn != expected {
			t.Errorf("expected %s, got %s", expected, dsn)
		}
	})

	// Test case 2: DATABASE_URL not set, use individual fields
	t.Run("DATABASE_URL not set", func(t *testing.T) {
		cfg := &Config{
			Database: DatabaseConfig{
				Host:    "localhost",
				Port:    "5432",
				User:    "user",
				Password: "password",
				Name:    "dbname",
				SSLMode: "disable",
			},
		}

		expected := "host=localhost port=5432 user=user password=password dbname=dbname sslmode=disable"
		if dsn := cfg.GetDSN(); dsn != expected {
			t.Errorf("expected %s, got %s", expected, dsn)
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("Valid with DATABASE_URL", func(t *testing.T) {
		cfg := &Config{
			Database: DatabaseConfig{
				URL: "postgres://user:pass@host:5432/db",
			},
			JWT: JWTConfig{
				Secret: "secure_secret",
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Valid with Password", func(t *testing.T) {
		cfg := &Config{
			Database: DatabaseConfig{
				Password: "password",
			},
			JWT: JWTConfig{
				Secret: "secure_secret",
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Invalid without URL or Password", func(t *testing.T) {
		cfg := &Config{
			Database: DatabaseConfig{},
			JWT: JWTConfig{
				Secret: "secure_secret",
			},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
