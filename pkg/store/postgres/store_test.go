package ttmpostgres

import (
	"strings"
	"testing"
	"ttm/pkg/config"
)

func TestDSNUsesConfiguredConnectionSettings(t *testing.T) {
	store := NewStore(config.PostgresConfig{
		Host:     "db.example.com",
		Port:     5433,
		User:     "ttm user",
		Password: "password with spaces",
		DBName:   "tasks",
		SSLMode:  "require",
	})

	dsn, err := store.dsn()
	if err != nil {
		t.Fatalf("dsn() error = %v", err)
	}

	want := "postgres://ttm%20user:password%20with%20spaces@db.example.com:5433/tasks?sslmode=require"
	if dsn != want {
		t.Errorf("dsn() = %q, want %q", dsn, want)
	}
}

func TestDSNUsesPasswordEnvironmentVariable(t *testing.T) {
	t.Setenv("TTM_POSTGRES_PASSWORD", "secret")
	store := NewStore(config.PostgresConfig{
		Host:        "localhost",
		Port:        5432,
		User:        "ttm",
		PasswordEnv: "TTM_POSTGRES_PASSWORD",
		DBName:      "ttmdb",
		SSLMode:     "disable",
	})

	dsn, err := store.dsn()
	if err != nil {
		t.Fatalf("dsn() error = %v", err)
	}
	if !strings.Contains(dsn, "ttm:secret@") {
		t.Errorf("dsn() = %q, expected password from environment", dsn)
	}
}

func TestDSNRequiresConfiguredPasswordEnvironmentVariable(t *testing.T) {
	store := NewStore(config.PostgresConfig{PasswordEnv: "MISSING_POSTGRES_PASSWORD"})

	_, err := store.dsn()
	if err == nil {
		t.Fatal("dsn() error = nil, want missing password environment variable error")
	}
}
