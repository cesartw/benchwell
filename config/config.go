package config

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/gtk"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/benchwell/assets"
)

const AppID = "io.benchwell"

type Setting struct {
	id    int64
	name  string
	value string

	m    sync.Mutex
	l    uint
	subs map[uint]*SettingUpdater
}

type SettingUpdater struct {
	p *Setting
	l uint
	f func(interface{})
	c chan interface{}
}

func (s *Setting) notify() {
	for _, updater := range s.subs {
		go func(c chan interface{}) {
			c <- s.value
		}(updater.c)
	}
}

func (s *SettingUpdater) Unsubscribe() {
	close(s.c)
	s.p.unsubscribe(s.l)
}

func (s *Setting) SetBool(b bool) {
	if b {
		s.value = "1"
	} else {
		s.value = "0"
	}
	s.notify()
}

func (s *Setting) SetString(v string) {
	s.value = v
	s.notify()
}

func (s *Setting) Bool() bool {
	return s.value == "1" || strings.EqualFold(s.value, "1")
}

func (s *Setting) String() string {
	return s.value
}

func (s *Setting) Int() int {
	i, _ := strconv.ParseInt(s.value, 10, 64)
	return int(i)
}

func (s *Setting) Int64() int64 {
	i, _ := strconv.ParseInt(s.value, 10, 64)
	return i
}

func (s *Setting) unsubscribe(l uint) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.subs, l)
}

func (s *Setting) Subscribe(f func(interface{})) *SettingUpdater {
	s.m.Lock()
	defer s.m.Unlock()

	s.l++
	u := &SettingUpdater{
		l: s.l,
		f: f,
		c: make(chan interface{}, 1),
	}

	go func() {
		for {
			select {
			case v := <-u.c:
				f(v)
			}
		}
	}()

	s.subs[s.l] = u
	return u
}

func (p *Setting) Settinglisher(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	for _, s := range p.subs {
		s.c <- v
	}
}

// Config ...
type Config struct {
	db *sql.DB
	*logrus.Logger

	Version     string
	Home        string
	Connections []*Connection

	loadedSettings map[string]*Setting

	GUI struct {
		CellWidth             *Setting
		ConnectionTabPosition *Setting
		TableTabPosition      *Setting
		PageSize              *Setting
		DarkMode              *Setting
	}
	Editor struct {
		WordWrap *Setting
	}
	EncryptionMode *Setting

	logFile string
}

func Init(path string) *Config {
	db, err := sql.Open("sqlite3", path+"/config.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(assets.DEFAULT_CONFIG)
	if err != nil {
		panic(err)
	}

	c := &Config{
		Logger:         logrus.New(),
		Home:           path,
		db:             db,
		loadedSettings: map[string]*Setting{},
	}

	c.Editor.WordWrap = c.Get("gui.editor.word_wrap")
	c.GUI.CellWidth = c.Get("gui.cell_width")
	c.GUI.ConnectionTabPosition = c.Get("gui.connection_tab_position")
	c.GUI.TableTabPosition = c.Get("gui.table_tab_position")
	c.GUI.PageSize = c.Get("gui.page_size")
	c.GUI.DarkMode = c.Get("gui.dark_mode")
	c.EncryptionMode = c.Get("encryption_mode")

	initKeyChain(c.EncryptionMode.String())
	err = c.loadConnections()
	if err != nil {
		panic(err)
	}

	if len(c.Connections) == 0 {
		c.SaveConnection(&Connection{Name: "New connection", Type: "tcp", Adapter: "mysql", Port: 3306})
	}

	return c
}

// name, adapter, type, database, host, options, user, password, port, encrypted
func (c *Config) loadConnections() error {
	c.Connections = nil
	rows, err := c.db.Query("SELECT * FROM connections")
	if err != nil {
		return err
	}

	for rows.Next() {
		conn := &Connection{}
		err := rows.Scan(&conn.ID, &conn.Name, &conn.Adapter, &conn.Type, &conn.Database,
			&conn.Host, &conn.Options, &conn.User, &conn.Password, &conn.Port, &conn.Encrypted,
			&conn.Socket, &conn.File, &conn.SshHost, &conn.SshAgent)
		if err != nil {
			return err
		}
		c.Connections = append(c.Connections, conn)
	}

	for _, conn := range c.Connections {
		conn.Decrypt(nil)
	}

	return c.loadQueries()
}

func (c *Config) loadQueries() error {
	rows, err := c.db.Query("SELECT * FROM queries")
	if err != nil {
		return err
	}

	connMap := map[int64]*Connection{}
	for _, conn := range c.Connections {
		connMap[conn.ID] = conn
	}

	for rows.Next() {
		query := &Query{}
		err := rows.Scan(&query.ID, &query.Name, &query.Query, &query.ConnectionID)
		if err != nil {
			return err
		}

		connMap[query.ConnectionID].Queries = append(connMap[query.ConnectionID].Queries, query)
	}

	for _, conn := range c.Connections {
		conn.Decrypt(nil)
	}

	return nil
}

func (c *Config) Get(s string) *Setting {
	setting, ok := c.loadedSettings[s]
	if ok {
		return setting
	}
	setting = &Setting{subs: map[uint]*SettingUpdater{}}

	row, err := c.db.Query("SELECT * FROM settings WHERE name = ? LIMIT 1", s)
	if err != nil {
		panic(err)
	}

	for row.Next() {
		if err := row.Scan(&setting.id, &setting.name, &setting.value); err != nil {
			panic(err)
		}
	}
	c.loadedSettings[s] = setting

	return setting
}

func (c *Config) EditorTheme() string {
	theme := "sqlaid-dark"
	if !c.GUI.DarkMode.Bool() {
		theme = "sqlaid-light"
	}
	return theme
}

func (c *Config) SaveSetting(s *Setting) error {
	return nil
}

func (c *Config) SaveConnection(conn *Connection) error {
	err := conn.Encrypt(nil)
	if err != nil {
		return err
	}

	if conn.ID == 0 {
		sql := `INSERT INTO connections(adapter, type, name, socket, file, host, port,
					user, password, database, sshhost, sshagent, options, encrypted)
				VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := c.db.Exec(sql,
			conn.Adapter, conn.Type, conn.Name, conn.Socket, conn.File,
			conn.Host, conn.Port, conn.User, conn.Password, conn.Database,
			conn.SshHost, conn.SshAgent, conn.Options, conn.Encrypted)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		conn.ID = id
		c.Connections = append(c.Connections, conn)
	} else {
		sql := `UPDATE connections
					SET adapter = ?, type = ?, name = ?, socket = ?, file = ?, host = ?, port = ?,
					user = ?, password = ?, database = ?, sshhost = ?, sshagent = ?, options = ?, encrypted = ?
				WHERE ID = ?`
		_, err := c.db.Exec(sql,
			conn.Adapter, conn.Type, conn.Name, conn.Socket, conn.File,
			conn.Host, conn.Port, conn.User, conn.Password, conn.Database,
			conn.SshHost, conn.SshAgent, conn.Options, conn.Encrypted, conn.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) SaveQuery(query *Query) error {
	if query.ID == 0 {
		sql := `INSERT INTO queries(name, query, connections_id)
				VALUES(?, ?, ?)`
		result, err := c.db.Exec(sql, query.Name, query.Query, query.ConnectionID)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		query.ID = id

		for _, conn := range c.Connections {
			if conn.ID == query.ConnectionID {
				conn.Queries = append(conn.Queries, query)
				break
			}
		}
	} else {
		sql := `UPDATE queries
					SET name = ?, query = ?
				WHERE ID = ?`
		_, err := c.db.Exec(sql,
			query.Name, query.Query, query.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) DeleteConnection(conn *Connection) error {
	if conn.ID != 0 {
		sql := `DELETE FROM connections WHERE id = ?`
		_, err := c.db.Exec(sql, conn.ID)
		if err != nil {
			return err
		}

		for i, co := range c.Connections {
			if co.ID != conn.ID {
				continue
			}

			c.Connections = append(c.Connections[:i], c.Connections[i+1:]...)
			break
		}
	}

	return nil
}

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

func (c *Config) CSS() string {
	style := ""
	if c.GUI.DarkMode.Bool() {
		style = assets.THEME_DARK + assets.BRAND + assets.BRAND_DARK
	} else {
		style = assets.THEME_LIGHT + assets.BRAND + assets.BRAND_LIGHT
	}
	return style
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

	pass, err := Keychain.Get(nil, c.Password)
	if err != nil {
		return err
	}

	c.Password = pass
	c.Encrypted = false

	return nil
}
