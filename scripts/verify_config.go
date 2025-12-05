package main

import (
	"fmt"
	"os"
	"github.com/hospital-emr/backend/internal/common/config"
)

func main() {
	// Test case 1: DATABASE_URL set
	expected := "postgres://user:pass@host:5432/db?sslmode=require"
	os.Setenv("DATABASE_URL", expected)
	
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL: expected,
		},
	}

	if dsn := cfg.GetDSN(); dsn != expected {
		fmt.Printf("FAIL: expected %s, got %s\n", expected, dsn)
	} else {
		fmt.Println("PASS: DATABASE_URL set")
	}

	// Test case 2: DATABASE_URL not set
	os.Unsetenv("DATABASE_URL")
	cfg = &config.Config{
		Database: config.DatabaseConfig{
			Host:    "localhost",
			Port:    "5432",
			User:    "user",
			Password: "password",
			Name:    "dbname",
			SSLMode: "disable",
		},
	}

	expected = "host=localhost port=5432 user=user password=password dbname=dbname sslmode=disable"
	if dsn := cfg.GetDSN(); dsn != expected {
		fmt.Printf("FAIL: expected %s, got %s\n", expected, dsn)
	} else {
		fmt.Println("PASS: DATABASE_URL not set")
	}
}
