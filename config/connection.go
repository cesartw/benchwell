package config

import (
	"bytes"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// Connection ...
type Connection struct {
	ID        int64
	Adapter   string
	Type      string
	Name      string
	Socket    string
	File      string
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	SshHost   string
	SshAgent  string
	Options   string
	Encrypted bool
	Queries   []*Query

	Config *Config
}

type Query struct {
	ID           int64
	Name         string
	Query        string
	ConnectionID int64
}

// GetDSN ...
func (c Connection) GetDSN() string {
	switch c.Adapter {
	case "mysql":
		return c.mysqlDSN()
	case "sqlite":
		return c.sqliteDSN()
	}

	return ""
}

func (c Connection) mysqlDSN() string {
	b := bytes.NewBuffer([]byte{})
	b.WriteString("mysql://")

	if c.User != "" {
		b.WriteString(c.User)
	}

	if c.Password != "" && c.User != "" {
		c.Decrypt(nil)
		b.WriteString(":" + c.Password)
	}

	if c.User != "" {
		b.WriteString("@")
	}

	switch c.Type {
	case "tcp":
		b.WriteString("tcp(" + c.Host)
		if c.Port != 0 {
			b.WriteString(fmt.Sprintf(":%d", c.Port))
		}
	case "socket":
		b.WriteString("unix(" + c.Socket)
	case "ssh":
		b.WriteString("ssh(" + c.Host)

		if c.Port != 0 {
			b.WriteString(fmt.Sprintf(":%d", c.Port))
		}
		b.WriteString("," + c.SshHost + ";" + c.SshAgent)
	default:
		return ""
	}

	b.WriteString(")")

	b.WriteString("/" + c.Database)

	if c.Options != "" {
		b.WriteString("?")
		b.WriteString(c.Options)
	}

	return b.String()
}

func (c Connection) sqliteDSN() string {
	b := bytes.NewBuffer([]byte{})
	b.WriteString("file:")

	b.WriteString(c.File)

	if c.Options != "" {
		b.WriteString("?")
		b.WriteString(c.Options)
	}

	return b.String()
}

func (c *Connection) LoadQueries() error {
	rows, err := c.Config.db.Query("SELECT * FROM db_queries WHERE connections_id = ?", c.ID)
	if err != nil {
		return err
	}

	c.Queries = nil
	for rows.Next() {
		query := &Query{}
		err := rows.Scan(&query.ID, &query.Name, &query.Query, &query.ConnectionID)
		if err != nil {
			return err
		}

		c.Queries = append(c.Queries, query)
	}

	return nil
}

func (c *Connection) DeleteQuery(name string) error {
	_, err := c.Config.db.Exec("DELETE FROM db_queries WHERE connections_id = ? AND name = ?", c.ID, name)
	if err != nil {
		return err
	}

	for i, q := range c.Queries {
		if q.Name != name {
			continue
		}

		c.Queries = append(c.Queries[:i], c.Queries[i+1:]...)
	}

	return nil
}

func (c *Connection) Encrypt(w *gtk.ApplicationWindow) error {
	if c.Encrypted {
		return nil
	}

	keys := map[string]string{
		"id": fmt.Sprintf("%d", c.ID),
	}

	path, err := Keychain.Set(nil, keys, c.Password)
	if err != nil {
		return err
	}
	c.Password = path
	c.Encrypted = true

	return nil
}

func (c *Connection) Decrypt(w *gtk.ApplicationWindow) error {
	if !c.Encrypted {
		return nil
	}

	pass, err := Keychain.Get(&w.Window, c.Password)
	if err != nil {
		return err
	}

	c.Password = pass
	c.Encrypted = false

	return nil
}
