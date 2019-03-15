package config

import kitlog "github.com/go-kit/kit/log"

// Config is a shared config object we can pass around to configure components.
// Contains attributes that the server needs in order to operate.
type Config struct {
	ServerAddr         string
	ConnStr            string
	EncryptionPassword string
	Verbose            bool
	Domains            []string
	DashboardURL       string
	ClientTimeout      int
	Logger             kitlog.Logger
}
