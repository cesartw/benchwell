package config

import (
	"database/sql"
	"os"

	"github.com/gotk3/gotk3/gtk"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/benchwell/assets"
)

var version = "dev"

const AppID = "io.benchwell"

// Config ...
type Config struct {
	db *sql.DB
	*logrus.Logger

	Version     string
	Home        string
	Connections []*Connection
	Collections []*HTTPCollection

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

func Init() *Config {
	userHome, _ := os.UserConfigDir()
	benchwellHome := userHome + "/benchwell"

	if version == "dev" {
		userHome = "./assets/data/"
		benchwellHome = userHome
	}

	db, err := sql.Open("sqlite3", benchwellHome+"/config.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(assets.DEFAULT_CONFIG)
	if err != nil {
		panic(err)
	}

	c := &Config{
		Logger:         logrus.New(),
		Home:           benchwellHome,
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
		c.SaveConnection(nil, &Connection{Name: "New connection", Type: "tcp", Adapter: "mysql", Port: 3306, Config: c})
	}

	err = c.loadCollections()
	if err != nil {
		c.Error(err)
	}

	return c
}

// name, adapter, type, database, host, options, user, password, port, encrypted
func (c *Config) loadConnections() error {
	c.Connections = nil
	rows, err := c.db.Query("SELECT * FROM db_connections")
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

	return nil
}

func (c *Config) loadCollections() error {
	c.Collections = nil
	rows, err := c.db.Query("SELECT * FROM http_collections")
	if err != nil {
		return err
	}

	for rows.Next() {
		collection := &HTTPCollection{Config: c}
		err := rows.Scan(&collection.ID, &collection.Count, &collection.Name)
		if err != nil {
			return err
		}
		c.Collections = append(c.Collections, collection)
	}

	return nil
}

func (c *Config) LoadQueries() error {
	rows, err := c.db.Query("SELECT * FROM db_queries")
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
	theme := "benchwell-dark"
	if !c.GUI.DarkMode.Bool() {
		theme = "benchwell-light"
	}
	return theme
}

func (c *Config) SaveSetting(s *Setting) error {
	return nil
}

func (c *Config) SaveConnection(w *gtk.ApplicationWindow, conn *Connection) error {
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
	conn.Decrypt(w)

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

func (c *Config) SaveHTTPItem(item *HTTPItem) error {
	if item.ID == 0 {
		sql := `INSERT INTO http_items(name, description, parent_id, is_folder, sort, http_collections_id, external_data, method, url, body, mime)
				VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := c.db.Exec(sql, item.Name,
			item.Description, item.ParentID, item.IsFolder,
			item.Sort, item.HTTPCollectionID, "", item.Method,
			item.URL, item.Body, item.Mime)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		item.ID = id
	} else {
		sql := `UPDATE http_items
					SET name = ?, description = ?, parent_id = ?, is_folder = ?,
					sort = ?, http_collections_id = ?, external_data = ?,
					method = ?, url = ?, body = ?, mime = ?)
				WHERE ID = ?`
		_, err := c.db.Exec(sql,
			item.Name, item.Description, item.ParentID,
			item.IsFolder, item.Sort, item.HTTPCollectionID,
			"", item.Method, item.URL, item.Body, item.Mime,
			item.ID)
		return err
	}

	for _, kv := range item.Params {
		kv.Config = c
		kv.HTTPItemID = item.ID
		err := kv.Save()
		if err != nil {
			return err
		}
	}
	for _, kv := range item.Headers {
		kv.Config = c
		kv.HTTPItemID = item.ID
		err := kv.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) SaveHTTPKV(kv *HTTPKV) error {
	if kv.ID == 0 {
		sql := `INSERT INTO http_kvs(key, value,  type, sort, enabled, http_items_id)
				VALUES(?, ?, ?, ?, ?, ?)`
		result, err := c.db.Exec(sql, kv.Key, kv.Value, kv.Type, kv.Sort, kv.Enabled, kv.HTTPItemID)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		kv.ID = id
	} else {
		sql := `update http_kvs(key, value,  type, sort, enabled, http_items_id)
				VALUES(?, ?, ?, ?, ?, ?)
				where id = ?`
		_, err := c.db.Exec(sql, kv.Key, kv.Value, kv.Type, kv.Sort, kv.Enabled, kv.HTTPItemID, kv.ID)
		return err
	}

	return nil
}

func (c *Config) SaveHTTPCollection(collection *HTTPCollection) error {
	if collection.ID == 0 {
		sql := `INSERT INTO http_collections(name)
				VALUES(?)`
		result, err := c.db.Exec(sql, collection.Name)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		collection.ID = id
	} else {
		sql := `UPDATE http_collections
					SET name = ?)
				WHERE ID = ?`
		_, err := c.db.Exec(sql,
			collection.Name, collection.ID)
		return err
	}

	return nil
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
