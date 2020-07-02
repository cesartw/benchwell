package config

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"bitbucket.org/goreorto/sqlaid/assets"
)

const AppID = "com.sqlaid"

var Env = &Config{}

func init() {
	Env.Log = logrus.New()
}

var hasher = md5.New()
var configPath = os.Getenv("HOME") + "/.config/sqlhero/config.toml"

type theme struct {
	Comment                  string `mapstructure:"comment"`
	CommentHashbang          string `mapstructure:"commentHashbang"`
	CommentMultiline         string `mapstructure:"commentMultiline"`
	CommentPreproc           string `mapstructure:"commentPreproc"`
	CommentSingle            string `mapstructure:"commentSingle"`
	CommentSpecial           string `mapstructure:"commentSpecial"`
	Generic                  string `mapstructure:"generic"`
	GenericDeleted           string `mapstructure:"genericDeleted"`
	GenericEmph              string `mapstructure:"genericEmph"`
	GenericError             string `mapstructure:"genericError"`
	GenericHeading           string `mapstructure:"genericHeading"`
	GenericInserted          string `mapstructure:"genericInserted"`
	GenericOutput            string `mapstructure:"genericOutput"`
	GenericPrompt            string `mapstructure:"genericPrompt"`
	GenericStrong            string `mapstructure:"genericStrong"`
	GenericSubheading        string `mapstructure:"genericSubheading"`
	GenericTraceback         string `mapstructure:"genericTraceback"`
	GenericUnderline         string `mapstructure:"genericUnderline"`
	Error                    string `mapstructure:"error"`
	Keyword                  string `mapstructure:"keyword"`
	KeywordConstant          string `mapstructure:"keywordConstant"`
	KeywordDeclaration       string `mapstructure:"keywordDeclaration"`
	KeywordNamespace         string `mapstructure:"keywordNamespace"`
	KeywordPseudo            string `mapstructure:"keywordPseudo"`
	KeywordReserved          string `mapstructure:"keywordReserved"`
	KeywordType              string `mapstructure:"keywordType"`
	Literal                  string `mapstructure:"literal"`
	LiteralDate              string `mapstructure:"literalDate"`
	Name                     string `mapstructure:"name"`
	NameAttribute            string `mapstructure:"nameAttribute"`
	NameBuiltin              string `mapstructure:"nameBuiltin"`
	NameBuiltinPseudo        string `mapstructure:"nameBuiltinPseudo"`
	NameClass                string `mapstructure:"nameClass"`
	NameConstant             string `mapstructure:"nameConstant"`
	NameDecorator            string `mapstructure:"nameDecorator"`
	NameEntity               string `mapstructure:"nameEntity"`
	NameException            string `mapstructure:"nameException"`
	NameFunction             string `mapstructure:"nameFunction"`
	NameLabel                string `mapstructure:"nameLabel"`
	NameNamespace            string `mapstructure:"nameNamespace"`
	NameOther                string `mapstructure:"nameOther"`
	NameTag                  string `mapstructure:"nameTag"`
	NameVariable             string `mapstructure:"nameVariable"`
	NameVariableClass        string `mapstructure:"nameVariableClass"`
	NameVariableGlobal       string `mapstructure:"nameVariableGlobal"`
	NameVariableInstance     string `mapstructure:"nameVariableInstance"`
	LiteralNumber            string `mapstructure:"literalNumber"`
	LiteralNumberBin         string `mapstructure:"literalNumberBin"`
	LiteralNumberFloat       string `mapstructure:"literalNumberFloat"`
	LiteralNumberHex         string `mapstructure:"literalNumberHex"`
	LiteralNumberInteger     string `mapstructure:"literalNumberInteger"`
	LiteralNumberIntegerLong string `mapstructure:"literalNumberIntegerLong"`
	LiteralNumberOct         string `mapstructure:"literalNumberOct"`
	Operator                 string `mapstructure:"operator"`
	OperatorWord             string `mapstructure:"operatorWord"`
	Other                    string `mapstructure:"other"`
	Punctuation              string `mapstructure:"punctuation"`
	LiteralString            string `mapstructure:"literalString"`
	LiteralStringBacktick    string `mapstructure:"literalStringBacktick"`
	LiteralStringChar        string `mapstructure:"literalStringChar"`
	LiteralStringDoc         string `mapstructure:"literalStringDoc"`
	LiteralStringDouble      string `mapstructure:"literalStringDouble"`
	LiteralStringEscape      string `mapstructure:"literalStringEscape"`
	LiteralStringHeredoc     string `mapstructure:"literalStringHeredoc"`
	LiteralStringInterpol    string `mapstructure:"literalStringInterpol"`
	LiteralStringOther       string `mapstructure:"literalStringOther"`
	LiteralStringRegex       string `mapstructure:"literalStringRegex"`
	LiteralStringSingle      string `mapstructure:"literalStringSingle"`
	LiteralStringSymbol      string `mapstructure:"literalStringSymbol"`
	Text                     string `mapstructure:"text"`
	TextWhitespace           string `mapstructure:"textWhitespace"`
	Background               string `mapstructure:"background"`
}

// Config ...
type Config struct {
	Version string `json:"-"`

	Connections []*Connection `mapstructure:"connections"`

	GUI struct {
		CellWidth             int    `mapstructure:"initial_cell_width"`
		ConnectionTabPosition string `mapstructure:"connection_tab_position"`
		TableTabPosition      string `mapstructure:"table_tab_position"`
		Editor                struct {
			WordWrap string `mapstructure:"word_wrap"`
			Theme    struct {
				Dark  theme `mapstructure:"dark"`
				Light theme `mapstructure:"light"`
			} `mapstructure:"theme"`
		}
		PageSize int  `mapstructure:"page_size"`
		DarkMode bool `mapstructure:"dark_mode"`
	}
	EncryptMode string `mapstructure:"encryptMode"`

	Log *logrus.Logger `mapstructure:"-" json:"-"`

	logFile string `mapstructure:"-" json:"-"`
}

// Connection ...
type Connection struct {
	Adapter   string  `mapstructure:"-"`
	Type      string  `mapstructure:"type"`
	Name      string  `mapstructure:"name"`
	Socket    string  `mapstructure:"socket"`
	Host      string  `mapstructure:"host"`
	Port      int     `mapstructure:"port"`
	User      string  `mapstructure:"user"`
	Password  string  `mapstructure:"password"`
	Database  string  `mapstructure:"database"`
	SshHost   string  `mapstructure:"sshhost"`
	SshAgent  string  `mapstructure:"sshsocket"`
	Options   string  `mapstructure:"options"`
	Encrypted bool    `mapstructure:"encrypted"`
	Queries   []Query `mapstructure:"queries"`
}

type Query struct {
	Name  string `mapstructure:"name"`
	Query string `mapstructure:"query"`
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

func (c Connection) Valid() bool {
	switch c.Type {
	case "tcp":
		if c.Host == "" {
			return false
		}
		if c.User == "" {
			return false
		}
	case "socket":
		if c.Socket == "" {
			return false
		}
	}

	return true
}

// Save current configuration
func (c *Config) Save(w *gtk.ApplicationWindow) error {
	var err error

	for _, conn := range c.Connections {
		err := conn.Encrypt(w)
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
		err := conn.Decrypt(w)
		if err != nil {
			return err
		}
	}

	return viper.WriteConfig()
}

func (c *Config) CSS() string {
	style := ""
	if c.GUI.DarkMode {
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
		"name":     c.Name,
		"host":     c.Host,
		"user":     c.User,
		"database": c.Database,
		"port":     fmt.Sprintf("%d", c.Port),
	}

	path, err := Keychain.Set(&w.Window, keys, c.Password)
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
