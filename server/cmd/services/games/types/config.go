package types

import "errors"

type Config struct {
	PostgresDSN   string
	JWTKeyPath    string
	JWTSecretPath string
}

func (c *Config) Validate() error {
	if c.PostgresDSN == "" {
		return errors.New("POSTGRES_DSN is required")
	}
	if c.JWTKeyPath == "" {
		return errors.New("JWT_KEY_PATH is required")
	}
	if c.JWTSecretPath == "" {
		return errors.New("JWT_SECRET_PATH is required")
	}
	return nil
}
