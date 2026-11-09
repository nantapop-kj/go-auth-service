package db

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/nantapop-kj/go-auth-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.PostgresConfig) (*gorm.DB, error) {
	if err := validate(cfg); err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn(cfg)), &gorm.Config{
		TranslateError: true,
		Logger: logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	return db, nil
}

func dsn(cfg config.PostgresConfig) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.Database,
		RawQuery: url.Values{
			"sslmode":     {cfg.SSLMode},
			"search_path": {cfg.Schema},
		}.Encode(),
	}
	return u.String()
}

func validate(cfg config.PostgresConfig) error {
	var missing []string
	for _, env := range []struct{ name, value string }{
		{"PG_USER", cfg.User},
		{"PG_PASSWORD", cfg.Password},
		{"PG_DATABASE", cfg.Database},
	} {
		if env.value == "" {
			missing = append(missing, env.name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing env: %s (is .env loaded? run the app with docker compose)", strings.Join(missing, ", "))
	}
	return nil
}
