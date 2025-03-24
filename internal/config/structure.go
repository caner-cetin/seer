package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// RootConfig holds the application configuration settings
type RootConfig struct {
	DB                    dbConfig `mapstructure:"db" yaml:"db"`
	DisplayAsciiArtOnHelp bool     `mapstructure:"display_ascii_art_on_help" yaml:"display_ascii_art_on_help"`
	Path                  string   `mapstructure:"-" yaml:"-"`
}

type dbConfig struct {
	URL  string       `mapstructure:"-" yaml:"-"`
	Auth dbAuthConfig `mapstructure:"auth" yaml:"auth"`
}

type dbAuthConfig struct {
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Database string `mapstructure:"database" yaml:"database"`
}

// Config is Config. how helpful.
var Config RootConfig

// SetDBURL constructs and sets the database URL using the authentication configuration
func (c *RootConfig) SetDBURL() {
	if c.DB.Auth.Host == "" {
		log.Fatal().Msg("db host is required")
	}
	if c.DB.Auth.Port == 0 {
		log.Fatal().Msg("db port is required")
	}
	if c.DB.Auth.User == "" {
		log.Fatal().Msg("db user is required")
	}
	if c.DB.Auth.Password == "" {
		log.Fatal().Msg("db password is required")
	}
	if c.DB.Auth.Database == "" {
		log.Fatal().Msg("db database is required")
	}
	c.DB.URL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.DB.Auth.User,
		c.DB.Auth.Password,
		c.DB.Auth.Host,
		c.DB.Auth.Port,
		c.DB.Auth.Database,
	)
}
