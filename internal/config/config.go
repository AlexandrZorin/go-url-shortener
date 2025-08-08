package config

import "flag"

type Config struct {
	ServerAddress string
	URL           string
}

func CreateConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Адрес сервера HTTP (адрес:порт)")
	flag.StringVar(&cfg.URL, "b", "http://localhost:8080", "URL для коротких ссылок")
	flag.Parse()
	return cfg
}
