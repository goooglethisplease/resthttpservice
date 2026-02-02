package config

import (
	"fmt"
	"os"
)

const (
	envHTTPPort = "HTTP_PORT"

	envDBHost     = "DB_HOST"
	envDBPort     = "DB_PORT"
	envDBUser     = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBName     = "DB_NAME"
	envDBSSLMode  = "DB_SSL_MODE"
)

const (
	defaultHTTPPort  = "8080"
	defaultDBHost    = "db"
	defaultDBPort    = "5432"
	defaultDBUser    = "postgres"
	defaultDBPass    = "postgres"
	defaultDBName    = "subscriptions"
	defaultDBSSLMode = "disable"
)

type Config struct {
	HTTPPort string
	DSN      string
}

func getenv(key, df string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return df
}

func MustLoad() *Config {
	port := getenv(envHTTPPort, defaultHTTPPort)

	dbhost := getenv(envDBHost, defaultDBHost)
	dbport := getenv(envDBPort, defaultDBPort)
	dbuser := getenv(envDBUser, defaultDBUser)
	dbpass := getenv(envDBPassword, defaultDBPass)
	dbname := getenv(envDBName, defaultDBName)
	dbsslmode := getenv(envDBSSLMode, defaultDBSSLMode)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s", dbuser, dbpass, dbhost, dbport, dbname, dbsslmode,
	)

	return &Config{
		HTTPPort: port,
		DSN:      dsn,
	}
}
