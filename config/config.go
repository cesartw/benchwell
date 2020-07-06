package config

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/sqlaid/assets"
)

const AppID = "com.sqlaid"

type Setting struct {
	id    int64
	name  string
	value string
}

func (s *Setting) SetBool(b bool) {
	if b {
		s.value = "1"
	} else {
		s.value = "0"
	}
}

func (s Setting) Bool() bool {
	return s.value == "1" || strings.EqualFold(s.value, "1")
}

func (s Setting) String() string {
	return s.value
}

func (s Setting) Int() int {
	i, _ := strconv.ParseInt(s.value, 10, 64)
	return int(i)
}

func (s Setting) Int64() int64 {
	i, _ := strconv.ParseInt(s.value, 10, 64)
	return i
}

// Config ...
type Config struct {
	db *sql.DB
	*logrus.Logger
	Version     string
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
		Theme    struct {
			Dark  theme
			Light theme
		}
	}
	EncryptionMode *Setting

	logFile string
}

func Init(path string) *Config {
	loadInitial := false
	if _, err := os.Stat(path); err != nil && !os.IsExist(err) {
		loadInitial = true
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	if loadInitial {
		_, err := db.Exec(assets.DEFAULT_CONFIG)
		if err != nil {
			panic(err)
		}
	}

	c := &Config{
		Logger:         logrus.New(),
		db:             db,
		loadedSettings: map[string]*Setting{},
	}

	c.Editor.Theme.Dark.SchemaName = "dark"
	c.Editor.Theme.Light.SchemaName = "light"

	c.loadStyle(&c.Editor.Theme.Light)
	c.loadStyle(&c.Editor.Theme.Dark)
	registerStyle(c)

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

	return c
}

func (c *Config) loadStyle(t *theme) {
	prefix := "gui.editor.theme." + t.SchemaName + "."
	t.Comment = c.Get(prefix + "comment")
	t.CommentHashbang = c.Get(prefix + "commentHashbang")
	t.CommentMultiline = c.Get(prefix + "commentMultiline")
	t.CommentPreproc = c.Get(prefix + "commentPreproc")
	t.CommentSingle = c.Get(prefix + "commentSingle")
	t.CommentSpecial = c.Get(prefix + "commentSpecial")
	t.Generic = c.Get(prefix + "generic")
	t.GenericDeleted = c.Get(prefix + "genericDeleted")
	t.GenericEmph = c.Get(prefix + "genericEmph")
	t.GenericError = c.Get(prefix + "genericError")
	t.GenericHeading = c.Get(prefix + "genericHeading")
	t.GenericInserted = c.Get(prefix + "genericInserted")
	t.GenericOutput = c.Get(prefix + "genericOutput")
	t.GenericPrompt = c.Get(prefix + "genericPrompt")
	t.GenericStrong = c.Get(prefix + "genericStrong")
	t.GenericSubheading = c.Get(prefix + "genericSubheading")
	t.GenericTraceback = c.Get(prefix + "genericTraceback")
	t.GenericUnderline = c.Get(prefix + "genericUnderline")
	t.Error = c.Get(prefix + "error")
	t.Keyword = c.Get(prefix + "keyword")
	t.KeywordConstant = c.Get(prefix + "keywordConstant")
	t.KeywordDeclaration = c.Get(prefix + "keywordDeclaration")
	t.KeywordNamespace = c.Get(prefix + "keywordNamespace")
	t.KeywordPseudo = c.Get(prefix + "keywordPseudo")
	t.KeywordReserved = c.Get(prefix + "keywordReserved")
	t.KeywordType = c.Get(prefix + "keywordType")
	t.Literal = c.Get(prefix + "literal")
	t.LiteralDate = c.Get(prefix + "literalDate")
	t.Name = c.Get(prefix + "name")
	t.NameAttribute = c.Get(prefix + "nameAttribute")
	t.NameBuiltin = c.Get(prefix + "nameBuiltin")
	t.NameBuiltinPseudo = c.Get(prefix + "nameBuiltinPseudo")
	t.NameClass = c.Get(prefix + "nameClass")
	t.NameConstant = c.Get(prefix + "nameConstant")
	t.NameDecorator = c.Get(prefix + "nameDecorator")
	t.NameEntity = c.Get(prefix + "nameEntity")
	t.NameException = c.Get(prefix + "nameException")
	t.NameFunction = c.Get(prefix + "nameFunction")
	t.NameLabel = c.Get(prefix + "nameLabel")
	t.NameNamespace = c.Get(prefix + "nameNamespace")
	t.NameOther = c.Get(prefix + "nameOther")
	t.NameTag = c.Get(prefix + "nameTag")
	t.NameVariable = c.Get(prefix + "nameVariable")
	t.NameVariableClass = c.Get(prefix + "nameVariableClass")
	t.NameVariableGlobal = c.Get(prefix + "nameVariableGlobal")
	t.NameVariableInstance = c.Get(prefix + "nameVariableInstance")
	t.LiteralNumber = c.Get(prefix + "literalNumber")
	t.LiteralNumberBin = c.Get(prefix + "literalNumberBin")
	t.LiteralNumberFloat = c.Get(prefix + "literalNumberFloat")
	t.LiteralNumberHex = c.Get(prefix + "literalNumberHex")
	t.LiteralNumberInteger = c.Get(prefix + "literalNumberInteger")
	t.LiteralNumberIntegerLong = c.Get(prefix + "literalNumberIntegerLong")
	t.LiteralNumberOct = c.Get(prefix + "literalNumberOct")
	t.Operator = c.Get(prefix + "operator")
	t.OperatorWord = c.Get(prefix + "operatorWord")
	t.Other = c.Get(prefix + "other")
	t.Punctuation = c.Get(prefix + "punctuation")
	t.LiteralString = c.Get(prefix + "literalString")
	t.LiteralStringBacktick = c.Get(prefix + "literalStringBacktick")
	t.LiteralStringChar = c.Get(prefix + "literalStringChar")
	t.LiteralStringDoc = c.Get(prefix + "literalStringDoc")
	t.LiteralStringDouble = c.Get(prefix + "literalStringDouble")
	t.LiteralStringEscape = c.Get(prefix + "literalStringEscape")
	t.LiteralStringHeredoc = c.Get(prefix + "literalStringHeredoc")
	t.LiteralStringInterpol = c.Get(prefix + "literalStringInterpol")
	t.LiteralStringOther = c.Get(prefix + "literalStringOther")
	t.LiteralStringRegex = c.Get(prefix + "literalStringRegex")
	t.LiteralStringSingle = c.Get(prefix + "literalStringSingle")
	t.LiteralStringSymbol = c.Get(prefix + "literalStringSymbol")
	t.Text = c.Get(prefix + "text")
	t.TextWhitespace = c.Get(prefix + "textWhitespace")
	t.Background = c.Get(prefix + "background")
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

	return nil
}

func (c *Config) Get(s string) *Setting {
	setting, ok := c.loadedSettings[s]
	if ok {
		return setting
	}
	setting = &Setting{}

	row, err := c.db.Query("SELECT * FROM config WHERE name = ? LIMIT 1", s)
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

type theme struct {
	SchemaName string

	Comment                  *Setting
	CommentHashbang          *Setting
	CommentMultiline         *Setting
	CommentPreproc           *Setting
	CommentSingle            *Setting
	CommentSpecial           *Setting
	Generic                  *Setting
	GenericDeleted           *Setting
	GenericEmph              *Setting
	GenericError             *Setting
	GenericHeading           *Setting
	GenericInserted          *Setting
	GenericOutput            *Setting
	GenericPrompt            *Setting
	GenericStrong            *Setting
	GenericSubheading        *Setting
	GenericTraceback         *Setting
	GenericUnderline         *Setting
	Error                    *Setting
	Keyword                  *Setting
	KeywordConstant          *Setting
	KeywordDeclaration       *Setting
	KeywordNamespace         *Setting
	KeywordPseudo            *Setting
	KeywordReserved          *Setting
	KeywordType              *Setting
	Literal                  *Setting
	LiteralDate              *Setting
	Name                     *Setting
	NameAttribute            *Setting
	NameBuiltin              *Setting
	NameBuiltinPseudo        *Setting
	NameClass                *Setting
	NameConstant             *Setting
	NameDecorator            *Setting
	NameEntity               *Setting
	NameException            *Setting
	NameFunction             *Setting
	NameLabel                *Setting
	NameNamespace            *Setting
	NameOther                *Setting
	NameTag                  *Setting
	NameVariable             *Setting
	NameVariableClass        *Setting
	NameVariableGlobal       *Setting
	NameVariableInstance     *Setting
	LiteralNumber            *Setting
	LiteralNumberBin         *Setting
	LiteralNumberFloat       *Setting
	LiteralNumberHex         *Setting
	LiteralNumberInteger     *Setting
	LiteralNumberIntegerLong *Setting
	LiteralNumberOct         *Setting
	Operator                 *Setting
	OperatorWord             *Setting
	Other                    *Setting
	Punctuation              *Setting
	LiteralString            *Setting
	LiteralStringBacktick    *Setting
	LiteralStringChar        *Setting
	LiteralStringDoc         *Setting
	LiteralStringDouble      *Setting
	LiteralStringEscape      *Setting
	LiteralStringHeredoc     *Setting
	LiteralStringInterpol    *Setting
	LiteralStringOther       *Setting
	LiteralStringRegex       *Setting
	LiteralStringSingle      *Setting
	LiteralStringSymbol      *Setting
	Text                     *Setting
	TextWhitespace           *Setting
	Background               *Setting
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
	Queries   []Query
}

type Query struct {
	Name  string `mapstructure:"name"`
	Query string `mapstructure:"query"`
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
