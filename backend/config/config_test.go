package config

import "testing"

func TestLoadDatabaseConfig(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_USER", "app")
	t.Setenv("DB_PASS", "secret")
	t.Setenv("DB_NAME", "expenses")

	config := LoadDatabaseConfig()

	if config.Host != "db.example.com" {
		t.Errorf("Host = %q, want %q", config.Host, "db.example.com")
	}
	if config.Port != "3307" {
		t.Errorf("Port = %q, want %q", config.Port, "3307")
	}
	if config.User != "app" {
		t.Errorf("User = %q, want %q", config.User, "app")
	}
	if config.Password != "secret" {
		t.Errorf("Password = %q, want %q", config.Password, "secret")
	}
	if config.Name != "expenses" {
		t.Errorf("Name = %q, want %q", config.Name, "expenses")
	}
}

func TestLoadDatabaseConfigEmptyWhenUnset(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASS", "")
	t.Setenv("DB_NAME", "")

	config := LoadDatabaseConfig()

	if config != (DatabaseConfig{}) {
		t.Errorf("LoadDatabaseConfig() = %+v, want empty config", config)
	}
}
