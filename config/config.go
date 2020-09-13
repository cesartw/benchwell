package config

import (
	"bytes"
	"database/sql"
	"html/template"
	"io"
	"os"
	"time"

	"github.com/gotk3/gotk3/gtk"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/benchwell/assets"
)

const AppID = "io.benchwell"

type Var struct {
	ID      int64
	Key     string
	Value   string
	Enabled bool
}

func (e *Var) Name() string      { return e.Key }
func (e *Var) Val() string       { return e.Value }
func (e *Var) IsEnabled() bool   { return e.Enabled }
func (e *Var) SetName(v string)  { e.Key = v }
func (e *Var) SetVal(v string)   { e.Value = v }
func (e *Var) SetEnabled(v bool) { e.Enabled = v }

type Env struct {
	ID        int64
	Name      string
	Variables []*EnvVar
}

func (e Env) Interpolate(s string) string {
	funcs := template.FuncMap{}
	for _, v := range e.Variables {
		value := v.Value
		funcs[v.Key] = func() string {
			return value
		}
	}

	t := template.New("")
	t.Funcs(funcs)
	t, err := t.Parse(s)
	if err != nil {
		return s
	}

	buff := bytes.NewBuffer(nil)
	err = t.Execute(buff, s)
	if err != nil {
		return s
	}

	return buff.String()
}

func (e *Env) Save() error {
	if e.ID == 0 {
		sql := `INSERT INTO environments(name)
				VALUES(?)`
		result, err := db.Exec(sql, e.Name)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		e.ID = id
	} else {
		sql := `UPDATE environments
					SET name = ?
				WHERE ID = ?`
		_, err := db.Exec(sql,
			e.Name,
			e.ID)
		if err != nil {
			return err
		}
	}

	for _, v := range e.Variables {
		v.EnvID = e.ID
		err := v.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

type EnvVar struct {
	Var
	EnvID int64
}

func (e *EnvVar) Save() error {
	if e.ID == 0 {
		sql := `INSERT INTO environment_variables(key, value, enabled, environment_id)
				VALUES(?, ?, ?, ?)`
		result, err := db.Exec(sql, e.Key, e.Value, e.Enabled, e.EnvID)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		e.ID = id
	} else {
		sql := `UPDATE environment_variables(key, value, enabled, environment_id)
				set key = ?, value = ?, enabled = ?, environment_id = ?
				WHERE ID = ?`
		_, err := db.Exec(sql,
			e.Key, e.Value, e.Enabled, e.EnvID,
			e.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

var (
	Version        = "dev"
	EncryptionMode = &Setting{}
	Editor         = struct {
		WordWrap *Setting
	}{}
	GUI = struct {
		CellWidth   *Setting
		TabPosition *Setting
		PageSize    *Setting
		DarkMode    *Setting
	}{}
	Home         string
	Connections  []*Connection
	Collections  []*HTTPCollection
	Environments []*Env
	ActiveWindow *gtk.ApplicationWindow

	loadedSettings map[string]*Setting
	logFile        string
	db             *sql.DB
	logger         *logrus.Logger
)

func Init() {
	logger = logrus.New()

	var err error
	userHome, _ := os.UserConfigDir()
	benchwellHome := userHome + "/benchwell"
	dbFile := benchwellHome + "/config.db"
	logger.Infof("Using %s", dbFile)

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(assets.DEFAULT_CONFIG)
	if err != nil {
		panic(err)
	}

	Home = benchwellHome
	loadedSettings = map[string]*Setting{}

	Editor.WordWrap = getSetting("gui.editor.word_wrap")
	GUI.CellWidth = getSetting("gui.cell_width")
	GUI.TabPosition = getSetting("gui.table_tab_position")
	GUI.PageSize = getSetting("gui.page_size")
	GUI.DarkMode = getSetting("gui.dark_mode")
	EncryptionMode = getSetting("encryption_mode")

	initKeyChain(EncryptionMode.String())
	err = loadConnections()
	if err != nil {
		panic(err)
	}

	err = loadCollections()
	if err != nil {
		logger.Error(err)
	}

	err = loadEnvironments()
	if err != nil {
		logger.Error(err)
	}
}

func getSetting(s string) *Setting {
	setting, ok := loadedSettings[s]
	if ok {
		return setting
	}
	setting = &Setting{subs: map[uint]*SettingUpdater{}}

	row, err := db.Query("SELECT * FROM settings WHERE name = ? LIMIT 1", s)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	for row.Next() {
		if err := row.Scan(&setting.id, &setting.name, &setting.value); err != nil {
			panic(err)
		}
	}
	loadedSettings[s] = setting

	return setting
}

// name, adapter, type, database, host, options, user, password, port, encrypted
func loadConnections() error {
	Connections = nil
	rows, err := db.Query("SELECT * FROM db_connections")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		conn := &Connection{}
		err := rows.Scan(&conn.ID, &conn.Name, &conn.Adapter, &conn.Type, &conn.Database,
			&conn.Host, &conn.Options, &conn.User, &conn.Password, &conn.Port, &conn.Encrypted,
			&conn.Socket, &conn.File, &conn.SshHost, &conn.SshAgent)
		if err != nil {
			return err
		}
		Connections = append(Connections, conn)
	}

	return nil
}

func loadCollections() error {
	Collections = nil
	rows, err := db.Query("SELECT * FROM http_collections")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		collection := &HTTPCollection{}
		err := rows.Scan(&collection.ID, &collection.Count, &collection.Name)
		if err != nil {
			return err
		}
		Collections = append(Collections, collection)
	}

	return nil
}

func loadEnvironments() error {
	Environments = nil

	rows, err := db.Query("SELECT * FROM environments")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		env := &Env{}
		err := rows.Scan(&env.ID, &env.Name)
		if err != nil {
			return err
		}
		Environments = append(Environments, env)
	}

	rows, err = db.Query("SELECT * FROM environment_variables")
	if err != nil {
		return err
	}
	defer rows.Close()

	variables := []*EnvVar{}
	for rows.Next() {
		envar := &EnvVar{}
		err := rows.Scan(&envar.ID, &envar.Key, &envar.Value, &envar.Enabled, &envar.EnvID)
		if err != nil {
			return err
		}
		variables = append(variables, envar)
	}

	envmap := map[int64]*Env{}
	for _, env := range Environments {
		envmap[env.ID] = env
	}

	for _, envar := range variables {
		envmap[envar.EnvID].Variables = append(envmap[envar.EnvID].Variables, envar)
	}

	return nil
}

func LoadQueries() error {
	rows, err := db.Query("SELECT * FROM db_queries")
	if err != nil {
		return err
	}
	defer rows.Close()

	connMap := map[int64]*Connection{}
	for _, conn := range Connections {
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

	return nil
}

func EditorTheme() string {
	theme := "benchwell-dark"
	if !GUI.DarkMode.Bool() {
		theme = "benchwell-light"
	}
	return theme
}

func SaveSetting(s *Setting) error {
	return nil
}

func SaveQuery(query *Query) error {
	if query.ID == 0 {
		sql := `INSERT INTO db_queries(name, query, connections_id)
				VALUES(?, ?, ?)`
		result, err := db.Exec(sql, query.Name, query.Query, query.ConnectionID)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		query.ID = id

		for _, conn := range Connections {
			if conn.ID == query.ConnectionID {
				conn.Queries = append(conn.Queries, query)
				break
			}
		}
	} else {
		sql := `UPDATE db_queries
					SET name = ?, query = ?
				WHERE ID = ?`
		_, err := db.Exec(sql,
			query.Name, query.Query, query.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func CSS() string {
	style := ""
	if GUI.DarkMode.Bool() {
		style = assets.THEME_DARK + assets.BRAND + assets.BRAND_DARK
	} else {
		style = assets.THEME_LIGHT + assets.BRAND + assets.BRAND_LIGHT
	}
	return style
}

type ctxlogger struct {
	*logrus.Entry
	start time.Time
}

func (c ctxlogger) Done() {
	c.WithField("duration", time.Since(c.start)).Debug("END")
}

func LogStart(fname string, args map[string]interface{}) func() {
	fields := logrus.Fields{"func": fname}
	for k, v := range args {
		fields[k] = v
	}

	e := logger.WithFields(fields)
	e.Debug("START")

	return ctxlogger{start: time.Now(), Entry: e}.Done
}

func SaveConnection(conn *Connection) error {
	err := conn.Encrypt()
	if err != nil {
		return err
	}

	if conn.ID == 0 {
		sql := `INSERT INTO db_connections(adapter, type, name, socket, file, host, port,
					user, password, database, sshhost, sshagent, options, encrypted)
				VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := db.Exec(sql,
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
		Connections = append(Connections, conn)
	} else {
		sql := `UPDATE db_connections
					SET adapter = ?, type = ?, name = ?, socket = ?, file = ?, host = ?, port = ?,
					user = ?, password = ?, database = ?, sshhost = ?, sshagent = ?, options = ?, encrypted = ?
				WHERE ID = ?`
		_, err := db.Exec(sql,
			conn.Adapter, conn.Type, conn.Name, conn.Socket, conn.File,
			conn.Host, conn.Port, conn.User, conn.Password, conn.Database,
			conn.SshHost, conn.SshAgent, conn.Options, conn.Encrypted, conn.ID)
		if err != nil {
			return err
		}
	}
	conn.Decrypt()

	return nil
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func SetLevel(l logrus.Level) {
	logger.SetLevel(l)
}

func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}
