package config

import (
	"bytes"
	"fmt"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
)

var configPath = os.Getenv("HOME") + "/.config/sqlhero/config.toml"

// Config ...
type Config struct {
	Connection []*Connection
	Debug      struct {
		Level int
	}
}

// Connection ...
type Connection struct {
	Adapter  string // "mysql"
	Type     string // "ssh"
	Name     string // "localhost"
	Host     string // "locahost"
	Port     int
	Username string
	Password string
	Database string
}

// New read $HOME/.config/sqlhero/config.toml
func New() (*Config, error) {
	conf := &Config{}

	if err := conf.Load(); err != nil {
		return nil, err
	}

	return conf, nil
}

// Load reads $HOME/.config/sqlhero/config.toml
func (c *Config) Load() error {
	if _, err := toml.DecodeFile(configPath, c); err != nil {
		return err
	}

	return nil
}

// Save current configuration
func (c *Config) Save() error {
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := toml.NewEncoder(f)

	return dec.Encode(c)
}

// DSN ...
func (c Connection) DSN() string {
	b := bytes.NewBuffer([]byte{})
	b.WriteString("mysql://")

	if c.Username != "" {
		b.WriteString(c.Username)
	}

	if c.Password != "" && c.Username != "" {
		b.WriteString(":" + c.Password)
	}

	if c.Username != "" {
		b.WriteString("@")
	}

	b.WriteString("tcp(" + c.Host)
	if c.Port != 0 {
		b.WriteString(fmt.Sprintf(":%d", c.Port))
	}
	b.WriteString(")")

	b.WriteString("/" + c.Database)
	if c.Database != "" {
		b.WriteString(c.Database)
	}

	u, _ := url.Parse(b.String())
	return u.String()
}
