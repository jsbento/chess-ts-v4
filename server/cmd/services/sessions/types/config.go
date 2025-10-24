package types

import "errors"

type Config struct {
	PostgresDSN string
}

func (c *Config) Validate() error {
	if c.PostgresDSN == "" {
		return errors.New("POSTGRES_DSN is required")
	}
	return nil
}
