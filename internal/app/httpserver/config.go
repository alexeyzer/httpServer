package httpserver

import "github.com/alexeyzer/httpServer/internal/app/store"

type config struct {
	BindAddr string  `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store *store.Config
}

func NewConfig() *config {
	return &config{
		BindAddr: "8080",
		LogLevel: "debug",
		Store: store.NewConfig(),
	}
}