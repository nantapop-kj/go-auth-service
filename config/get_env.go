package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key string, defaultValue string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return defaultValue
	}
	return val
}

func GetEnvInt(key string, defaultValue int) int {
	return parseEnv(key, defaultValue, strconv.Atoi)
}

func GetEnvBool(key string, defaultValue bool) bool {
	return parseEnv(key, defaultValue, strconv.ParseBool)
}

func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	return parseEnv(key, defaultValue, time.ParseDuration)
}

func GetEnvList(key string) []string {
	var out []string
	for item := range strings.SplitSeq(GetEnv(key, ""), ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}

func parseEnv[T any](key string, defaultValue T, parse func(string) (T, error)) T {
	val := GetEnv(key, "")
	if val == "" {
		return defaultValue
	}
	parsed, err := parse(val)
	if err != nil {
		log.Printf("invalid %s=%q, using default %v", key, val, defaultValue)
		return defaultValue
	}
	return parsed
}
