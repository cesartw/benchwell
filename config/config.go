package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var Env = &Config{}

func init() {
	Env.Log = logrus.New()
}

var hasher = md5.New()
var configPath = os.Getenv("HOME") + "/.config/sqlhero/config.toml"

// Config ...
type Config struct {
	Version string

	Connections []*Connection

	GUI struct {
		TabPosition    string `mapstructure:"tab_position"`
		SubTabPosition string `mapstructure:"sub_tab_position"`
		Editor         struct {
			WordWrap string `mapstructure:"word_wrap"`
		}
	}

	Log    *logrus.Logger
	phrase string
}

// Connection ...
type Connection struct {
	Adapter  string
	Type     string
	Name     string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Options  string
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

	if c.Database != "" {
		b.WriteString("/" + c.Database)
	}

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
	err := c.encrypt()
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := toml.NewEncoder(f)
	err = dec.Encode(c)
	if err != nil {
		return err
	}

	return c.decrypt()
}

func (c *Config) encrypt() error {
	if c.phrase == "" {
		return nil
	}

	for _, conn := range c.Connections {
		if conn.Password == "" {
			continue
		}

		var err error
		conn.Password, err = encrypt(c.phrase, conn.Password)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) decrypt() error {
	if c.phrase == "" {
		return nil
	}

	for _, conn := range c.Connections {
		if conn.Password == "" {
			continue
		}

		result, err := decrypt(c.phrase, conn.Password)
		if err != nil {
			return err
		}

		conn.Password = result
	}

	return nil
}

func encrypt(phrase, text string) (string, error) {
	textBytes := []byte(text)
	block, err := aes.NewCipher(hash(phrase))
	if err != nil {
		return "", err
	}

	b := base64.StdEncoding.EncodeToString(textBytes)
	ciphertext := make([]byte, aes.BlockSize+len(b))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(phrase, text string) (string, error) {
	if text == "" {
		return "", nil
	}

	textBytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	if len(textBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short: " + text)
	}

	block, err := aes.NewCipher(hash(phrase))
	if err != nil {
		return "", err
	}

	iv := textBytes[:aes.BlockSize]
	textBytes = textBytes[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(textBytes, textBytes)

	data, err := base64.StdEncoding.DecodeString(string(textBytes))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func hash(phrase string) []byte {
	defer func() {
		hasher.Reset()
	}()

	hasher.Write([]byte(phrase))
	return hasher.Sum(nil)
}
