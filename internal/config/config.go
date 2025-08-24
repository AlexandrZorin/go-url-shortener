package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	URL           string
}

func CreateConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес сервера HTTP (адрес:порт)")
	flag.StringVar(&cfg.URL, "b", "http://localhost:8080", "URL для коротких ссылок")
	flag.Parse()

	serverAddressFlag := cfg.ServerAddress
	baseURLFlag := cfg.URL

	serverAddressEnv := os.Getenv("SERVER_ADDRESS")
	baseURLEnv := os.Getenv("BASE_URL")

	if serverAddressEnv != "" {
		cfg.ServerAddress = serverAddressEnv
	} else {
		cfg.ServerAddress = serverAddressFlag
	}

	if baseURLEnv != "" {
		cfg.URL = baseURLEnv
	} else {
		cfg.URL = baseURLFlag
	}

	return cfg
}
