package config

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
)

const minJWTSecretLength = 32

func LoadJWTSecret(cfg *Config) ([]byte, error) {
	if cfg.JWTSecret == "" {
		if cfg.IsProd() {
			return nil, errors.New("JWT_SECRET is required in prod")
		}
		log.Println("JWT_SECRET not set: using a temporary secret, tokens will be invalid after restart")
		secret := make([]byte, 64)
		_, err := rand.Read(secret)
		return secret, err
	}

	if len(cfg.JWTSecret) < minJWTSecretLength {
		return nil, fmt.Errorf("JWT_SECRET must be at least %d characters", minJWTSecretLength)
	}
	return []byte(cfg.JWTSecret), nil
}
