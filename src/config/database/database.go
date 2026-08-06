package database

import (
	"books/config"
	"fmt"
)

func Dsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host(),
		user(),
		password(),
		name(),
		port(),
	)
}

func host() string {
	return config.Env("DB_HOST", "")
}

func name() string {
	return config.Env("DB_DATABASE", "")
}

func password() string {
	return config.Env("DB_PASSWORD", "")
}

func user() string {
	return config.Env("DB_USERNAME", "")
}

func port() string {
	return config.Env("DB_PORT", "5432")
}
