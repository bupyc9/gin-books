package config

import "os"

func Env(key string, def string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return def
	}

	return value
}
