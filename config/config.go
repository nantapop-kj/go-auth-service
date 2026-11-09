package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/nantapop-kj/go-auth-service/enums"
)

type Config struct {
	AppName     string
	AppEnv      enums.AppEnv
	ServerPort  string
	CORSOrigins string
	ProxyHeader string

	JWTIssuer       string
	JWTAudience     []string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	LoginMaxAttempts  int
	LoginLockDuration time.Duration

	CookieEnabled     bool
	CookieName        string
	CookieRefreshName string
	CookieDomain      string
	CookieSameSite    string
	CookieSecure      bool

	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Schema   string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func (c *Config) IsProd() bool { return c.AppEnv == enums.EnvProd }

func Load() *Config {
	_ = godotenv.Load()

	appName := GetEnv("APP_NAME", "go-auth-service")
	appEnv := enums.AppEnv(GetEnv("APP_ENV", string(enums.EnvDev)))
	cookieName := GetEnv("COOKIE_NAME", "access_token")

	return &Config{
		AppName:     appName,
		AppEnv:      appEnv,
		ServerPort:  GetEnv("SERVER_PORT", "3001"),
		CORSOrigins: GetEnv("CORS_ORIGINS", "http://localhost:3000"),
		ProxyHeader: GetEnv("PROXY_HEADER", ""),

		JWTIssuer:       GetEnv("JWT_ISSUER", appName),
		JWTAudience:     GetEnvList("JWT_AUDIENCE"),
		JWTSecret:       GetEnv("JWT_SECRET", ""),
		AccessTokenTTL:  GetEnvDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: GetEnvDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour),

		LoginMaxAttempts:  GetEnvInt("LOGIN_MAX_ATTEMPTS", 5),
		LoginLockDuration: GetEnvDuration("LOGIN_LOCK_DURATION", 15*time.Minute),

		CookieEnabled:     GetEnvBool("COOKIE_ENABLED", true),
		CookieName:        cookieName,
		CookieRefreshName: GetEnv("COOKIE_REFRESH_NAME", cookieName+"_refresh"),
		CookieDomain:      GetEnv("COOKIE_DOMAIN", ""),
		CookieSameSite:    strings.ToLower(GetEnv("COOKIE_SAME_SITE", "lax")),
		CookieSecure:      GetEnvBool("COOKIE_SECURE", appEnv == enums.EnvProd),

		Postgres: PostgresConfig{
			Host:     GetEnv("PG_HOST", "localhost"),
			Port:     GetEnv("PG_PORT", "5432"),
			User:     GetEnv("PG_USER", ""),
			Password: GetEnv("PG_PASSWORD", ""),
			Database: GetEnv("PG_DATABASE", ""),
			Schema:   GetEnv("PG_SCHEMA", "public"),
			SSLMode:  GetEnv("PG_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Address:  GetEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       GetEnvInt("REDIS_DB", 0),
		},
	}
}
