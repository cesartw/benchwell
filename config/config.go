package config

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const AppID = "com.sqlaid"

var Env = &Config{}

func init() {
	Env.Log = logrus.New()
}

var hasher = md5.New()
var configPath = os.Getenv("HOME") + "/.config/sqlhero/config.toml"

// Config ...
type Config struct {
	Version string `json:"-"`

	Connections []*Connection

	GUI struct {
		TabPosition    string `mapstructure:"tab_position"`
		SubTabPosition string `mapstructure:"sub_tab_position"`
		Editor         struct {
			WordWrap string `mapstructure:"word_wrap"`
		}
		PageSize int  `mapstructure:"page_size"`
		DarkMode bool `mapstructure:"page_size"`
	}

	Log *logrus.Logger `json:"-"`

	logFile string `json:"-"`
}

// Connection ...
type Connection struct {
	Adapter   string
	Type      string
	Name      string
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	Options   string
	Encrypted bool
}

// GetDSN ...
func (c Connection) GetDSN() string {
	b := bytes.NewBuffer([]byte{})
	b.WriteString("mysql://")

	if c.User != "" {
		b.WriteString(c.User)
	}

	if c.Password != "" && c.User != "" {
		b.WriteString(":" + c.Password)
	}

	if c.User != "" {
		b.WriteString("@")
	}

	b.WriteString("tcp(" + c.Host)
	if c.Port != 0 {
		b.WriteString(fmt.Sprintf(":%d", c.Port))
	}
	b.WriteString(")")

	b.WriteString("/" + c.Database)

	if c.Options != "" {
		b.WriteString("?")
		b.WriteString(c.Options)
	}

	return b.String()
}

func (c Connection) Valid() bool {
	if c.Host == "" {
		return false
	}
	if c.User == "" {
		return false
	}

	return true
}

// Save current configuration
func (c *Config) Save() error {
	var err error

	for _, conn := range c.Connections {
		err := conn.Encrypt()
		if err != nil {
			return err
		}
	}

	d, err := json.Marshal(c)
	if err != nil {
		c.Log.Error(err)
		return err
	}

	viper.MergeConfig(bytes.NewReader(d))

	for _, conn := range c.Connections {
		err := conn.Decrypt()
		if err != nil {
			return err
		}
	}

	return viper.WriteConfig()
}

func (c *Connection) Encrypt() error {
	if c.Encrypted {
		return nil
	}

	keys := map[string]string{
		"name":     c.Name,
		"host":     c.Host,
		"user":     c.User,
		"database": c.Database,
		"port":     fmt.Sprintf("%d", c.Port),
	}

	path, err := Keychain.Set(keys, c.Password)
	if err != nil {
		return err
	}
	c.Password = path
	c.Encrypted = true

	return nil
}

func (c *Connection) Decrypt() error {
	if !c.Encrypted {
		return nil
	}

	pass, err := Keychain.Get(c.Password)
	if err == nil {
		c.Password = pass
	} else {
		Env.Log.Error(err)
	}
	c.Encrypted = false

	return nil
}
