package jobbr

import "github.com/caarlos0/env"

type EnvConfig struct {
	URL  string `env:"JOBS_API_URL" envDefault:"http://127.0.0.1"`
	Port int    `env:"JOBS_API_PORT" envDefault:"8080"`
}

// ParseEnv parses the environment variables and returns a new EnvConfig
func ParseEnv() (*EnvConfig, error) {
	cfg := &EnvConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
