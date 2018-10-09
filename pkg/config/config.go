package config

import kitlog "github.com/go-kit/kit/log"

// Config is a shared config object we can pass around to configure components.
type Config struct {
	ServerAddr string
	ConnStr    string
	Verbose    bool
	Logger     kitlog.Logger
}
