package db

import (
	"testing"

	"github.com/nantapop-kj/go-auth-service/config"
)

func TestDSNEscapesSpecialCharacters(t *testing.T) {
	got := dsn(config.PostgresConfig{Host: "go-auth-db", Port: "5432", User: "admin", Password: "p@ss word/#?", Database: "go-auth-db", Schema: "public", SSLMode: "disable"})
	want := "postgres://admin:p%40ss%20word%2F%23%3F@go-auth-db:5432/go-auth-db?search_path=public&sslmode=disable"
	if got != want {
		t.Fatalf("dsn:\n got %s\nwant %s", got, want)
	}
}

func TestValidateReportsMissingEnv(t *testing.T) {
	if err := validate(config.PostgresConfig{}); err == nil {
		t.Fatal("expected error for empty config")
	}
}
