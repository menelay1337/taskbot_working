package config

import (
	"flag"
	"log"
)

type Config struct {
	Token  string
	Dsn		string
}

func MustLoad() Config {
	token := flag.String(
		"token",
		"",
		"token for access to telegram bot",
	)

	dsn := flag.String(
		"dsn",
		"",
		"connection string for MySQL",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	if *dsn == "" {
		log.Fatal("MySQL connection string is not specified")
	}

	return Config{
		Token:            *token,
		Dsn: *dsn,
	}
}
