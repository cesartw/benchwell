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
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
)

var hasher = md5.New()
var configPath = os.Getenv("HOME") + "/.config/sqlhero/config.toml"

// Config ...
type Config struct {
	Connection []*Connection
	Debug      struct {
		Level int
	}

	phrase string
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
func New(phrase string) (*Config, error) {
	conf := &Config{}
	conf.phrase = phrase

	if err := conf.load(); err != nil {
		return nil, err
	}

	err := conf.decrypt()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// Load reads $HOME/.config/sqlhero/config.toml
func (c *Config) load() error {
	if _, err := toml.DecodeFile(configPath, c); err != nil {
		return err
	}

	return nil
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

func (c *Config) encrypt() error {
	if c.phrase == "" {
		return nil
	}

	for _, conn := range c.Connection {
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

	for _, conn := range c.Connection {
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
