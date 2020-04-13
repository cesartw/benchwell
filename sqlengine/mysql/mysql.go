package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	// mysql implementation
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type mysqlDriver struct{}

type mysqlConn struct {
	dsn       string
	db        *sql.DB
	lastError error
}

type mysqlDb struct {
	db   *sql.DB
	name string
}

func init() {
	driver.RegisterDriver("mysql", &mysqlDriver{})
}

func (d *mysqlDriver) Connect(dsn string) (driver.Connection, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &mysqlConn{dsn: dsn, db: db}, nil
}

func (c *mysqlConn) Reconnect() error {
	c.db.Close()

	db, err := sql.Open("mysql", c.dsn)
	if err != nil {
		return err
	}

	c.db = db

	return nil
}

func (c *mysqlConn) UseDatabase(db string) error {
	if db == "" {
		return errors.New("database name empty")
	}

	_, err := c.db.Exec(fmt.Sprintf("USE %s", db))
	return err
}

// Disconnect ...
func (c *mysqlConn) Disconnect() error {
	return c.db.Close()
}

// LastError ...
func (c *mysqlConn) LastError() error {
	return c.lastError
}

// Databases ...
func (c *mysqlConn) Databases() ([]driver.Database, error) {
	rows, err := c.db.Query(`SHOW databases`)
	if err != nil {
		return nil, err
	}

	dbs := make([]driver.Database, 0)
	for rows.Next() {
		var n string
		err := rows.Scan(&n)
		if err != nil {
			return nil, err
		}

		dbs = append(dbs, driver.Database(&mysqlDb{c.db, n}))
	}

	return dbs, nil
}

func (d *mysqlDb) Name() string {
	return d.name
}

func (d *mysqlDb) Tables() ([]string, error) {
	rows, err := d.db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}

	tables := []string{}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (d *mysqlDb) TableDefinition(tableName string) ([]driver.ColDef, error) {
	sqlRows, err := d.db.Query("DESCRIBE " + tableName)
	if err != nil {
		return nil, err
	}

	defer sqlRows.Close()

	defs := []driver.ColDef{}
	for sqlRows.Next() {
		var name, ftype, nullable, key, dflt, other sql.NullString
		if err := sqlRows.Scan(&name, &ftype, &nullable, &key, &dflt, &other); err != nil {
			return nil, err
		}

		t, precision, vv, unsigned := d.parseType(ftype.String)
		defs = append(defs, driver.ColDef{
			Name:      name.String,
			Type:      t,
			PK:        key.String == "PRI",
			Nullable:  nullable.String == "YES",
			Unsigned:  unsigned,
			Values:    vv,
			Precision: precision,
		})
	}

	if err := sqlRows.Err(); err != nil {
		return nil, err
	}

	return defs, nil
}

func (d *mysqlDb) Query(query string) (columnNames []string, data [][]interface{}, err error) {
	data = make([][]interface{}, 0)

	var sqlRows *sql.Rows
	sqlRows, err = d.db.Query(query)
	if err != nil {
		return nil, nil, err
	}

	defer sqlRows.Close()

	columnNames, err = sqlRows.Columns()
	if err != nil {
		return nil, nil, err
	}

	for sqlRows.Next() {
		row := make([]interface{}, len(columnNames))

		for ci := range columnNames {
			row[ci] = &row[ci]
		}

		if err := sqlRows.Scan(row...); err != nil {
			return nil, nil, err
		}

		data = append(data, row)
	}
	if err := sqlRows.Err(); err != nil {
		return nil, nil, err
	}

	return columnNames, data, err
}

func (d *mysqlDb) ParseValue(def driver.ColDef, value string) interface{} {
	if value == driver.NULL_PATTERN {
		if def.Nullable {
			return nil
		}
	}
	if def.PK && value == driver.NULL_PATTERN {
		return nil
	}

	switch def.Type {
	case driver.TYPE_BOOLEAN:
		return strings.EqualFold(value, "true") || value == "1"
	case driver.TYPE_FLOAT:
		v, _ := strconv.ParseFloat(value, 64)
		return v
	case driver.TYPE_INT:
		v, _ := strconv.ParseInt(value, 10, 64)
		return v
	}

	return value
}

var typerg = regexp.MustCompile(`([a-z ]+)(\((.+)\))?\s?(unsigned)?`)

func (d *mysqlDb) parseType(mysqlStringType string) (driver.TYPE, int, []string, bool) {
	matches := typerg.FindStringSubmatch(mysqlStringType)
	t := matches[1] // type
	s := matches[3] // size/precision
	u := matches[4] // unsigned

	switch t {
	case "enum":
		return driver.TYPE_LIST, 0, strings.Split(s, ","), false
	case "text":
		return driver.TYPE_STRING, 0, nil, false
	case "varchar":
		si, _ := strconv.Atoi(s)
		return driver.TYPE_STRING, si, nil, false
	case "int", "smallint", "mediumint", "bigint":
		si, _ := strconv.Atoi(s)
		return driver.TYPE_INT, si, nil, u == "unsigned"
	case "tinyint":
		if s == "1" {
			return driver.TYPE_BOOLEAN, 0, nil, false
		}

		si, _ := strconv.Atoi(s)
		return driver.TYPE_INT, si, nil, u == "unsigned"
	case "double precision", "double", "float", "decimal":
		return driver.TYPE_FLOAT, 0, nil, u == "unsigned"
	case "time", "datetime":
		return driver.TYPE_DATE, 0, nil, false
	}

	return driver.TYPE_STRING, 0, nil, true
}

func (d *mysqlDb) FetchTable(
	tableName string, limit, offset int64,
) (
	colDef []driver.ColDef, rows [][]interface{}, err error,
) {
	var sqlRows *sql.Rows
	sqlRows, err = d.db.Query(fmt.Sprintf(`
SELECT *
FROM %s
LIMIT ?, ?
	`, tableName), limit, offset)
	if err != nil {
		return nil, nil, err
	}

	defer sqlRows.Close()

	columns, err := sqlRows.Columns()
	if err != nil {
		return nil, nil, err
	}

	rows = make([][]interface{}, 0)

	for sqlRows.Next() {
		row := make([]interface{}, len(columns))

		for ci := range columns {
			row[ci] = &row[ci]
		}

		if err := sqlRows.Scan(row...); err != nil {
			return nil, nil, err
		}

		rows = append(rows, row)
	}
	if err := sqlRows.Err(); err != nil {
		return nil, nil, err
	}

	colDef, err = d.TableDefinition(tableName)
	if err != nil {
		return nil, nil, err
	}

	return colDef, rows, err
}

func (d *mysqlDb) DeleteRecord(tableName string, cols []driver.ColDef, args []interface{}) error {

	if len(cols) == 0 {
		return errors.New("cols is empty")
	}

	wheres := []string{}
	for i := range cols {
		wheres = append(wheres, fmt.Sprintf("%s = ?", cols[i].Name))
	}

	query := fmt.Sprintf(`
	DELETE FROM %s
	WHERE %s
	`, tableName, strings.Join(wheres, " AND "))

	config.Env.Log.WithFields(logrus.Fields{"query": query, "args": args}).Debug("DeleteRecord")
	_, err := d.db.Exec(query, args...)

	return err
}

func (d *mysqlDb) UpdateRecord(
	tableName string,
	cols []driver.ColDef,
	values, oldValues []interface{},
) (string, error) {
	if len(cols) != len(values) || len(values) != len(oldValues) {
		return "", errors.New("columns and values count doesn't match")
	}

	var pk *driver.ColDef
	for _, def := range cols {
		if def.PK {
			pk = &def
			break
		}
	}
	if pk == nil {
		return "", errors.New("table doesn't have a primary key")
	}

	sets := []string{}
	args := []interface{}{}
	var ID interface{}

	for i := range values {
		if cols[i].PK {
			config.Env.Log.Debug(cols[i].Name, oldValues[i])
			ID = oldValues[i]
			if reflect.DeepEqual(oldValues[i], values[i]) {
				continue
			}
		}
		sets = append(sets, fmt.Sprintf("%s = ?", cols[i].Name))
		args = append(args, values[i])
	}

	query :=
		`
UPDATE %s
SET %s
WHERE %s = ?`

	query = fmt.Sprintf(query, tableName, strings.Join(sets, ", "), pk.Name)
	args = append(args, ID)

	config.Env.Log.WithFields(logrus.Fields{"query": query, "args": args}).Debug("InsertRecord")
	result, err := d.db.Exec(query, args...)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (d *mysqlDb) UpdateField(
	tableName string,
	cols []driver.ColDef,
	values []interface{},
) (string, error) {
	if len(cols) != len(values) {
		return "", errors.New("columns and values count doesn't match")
	}

	if len(cols) == 1 {
		return "", errors.New("keys or changes are not present")
	}

	lastIndex := len(cols) - 1
	args := []interface{}{values[lastIndex]}

	wheres := []string{}

	for i := 0; i <= len(cols)-2; i++ {
		wheres = append(wheres, fmt.Sprintf("`%s` = ?", cols[i].Name))
		args = append(args, values[i])
	}

	query := fmt.Sprintf("UPDATE `%s` SET `%s` = ? WHERE %s",
		tableName,
		cols[lastIndex].Name,
		strings.Join(wheres, " AND "))

	config.Env.Log.WithFields(logrus.Fields{"query": query, "args": args}).Debug("UpdateField")
	result, err := d.db.Exec(query, args...)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (d *mysqlDb) UpdateFields(
	tableName string,
	cols []driver.ColDef,
	values []interface{},
	keycount int,
) (string, error) {
	if len(cols) != len(values) {
		return "", errors.New("columns and values count doesn't match")
	}

	if len(cols) >= keycount {
		return "", errors.New("keys or changes are not present")
	}

	wheres := []string{}
	sets := []string{}

	for i := range values {
		if cols[i].PK && i <= keycount {
			wheres = append(wheres, fmt.Sprintf("`%s` = ?", cols[i].Name))
		} else {
			sets = append(sets, fmt.Sprintf("`%s` = ?", cols[i].Name))
		}
	}

	if len(sets) == 0 {
	}

	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s",
		tableName,
		sets,
		strings.Join(wheres, " AND "))

	config.Env.Log.WithFields(logrus.Fields{"query": query, "args": values}).Debug("UpdateFields")
	result, err := d.db.Exec(query, values...)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (d *mysqlDb) InsertRecord(
	tableName string,
	cols []driver.ColDef,
	values []interface{},
) ([]interface{}, error) {
	if len(cols) != len(values) {
		return nil, errors.New("columns and values count doesn't match")
	}

	collist := make([]string, len(cols))
	qm := make([]string, len(cols))
	args := make([]interface{}, len(cols))
	for i := range cols {
		collist[i] = cols[i].Name
		args[i] = values[i]
		qm[i] = "?"
	}

	query :=
		`
INSERT INTO %s(%s)
VALUES (%s)`

	query = fmt.Sprintf(query, tableName, strings.Join(collist, ","), strings.Join(qm, ","))

	config.Env.Log.WithFields(logrus.Fields{"query": query, "args": args}).Debug("InsertRecord")
	result, err := d.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return d.fetchRecord(tableName, cols, id)
}

func (d *mysqlDb) fetchRecord(
	tableName string,
	cols []driver.ColDef,
	id int64,
) (
	[]interface{},
	error,
) {
	var pk *driver.ColDef
	for _, def := range cols {
		if def.PK {
			pk = &def
			break
		}
	}
	if pk == nil {
		return nil, errors.New("table doesn't have a primary key")
	}

	query :=
		`
SELECT *
FROM %s
WHERE %s = ?`

	query = fmt.Sprintf(query, tableName, pk.Name)

	sqlRows, err := d.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer sqlRows.Close()

	columns, err := sqlRows.Columns()
	if err != nil {
		return nil, err
	}

	rows := make([][]interface{}, 0)

	for sqlRows.Next() {
		row := make([]interface{}, len(columns))

		for ci := range row {
			row[ci] = &row[ci]
		}

		if err := sqlRows.Scan(row...); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}
	if err := sqlRows.Err(); err != nil {
		return nil, err
	}

	return rows[0], err
}
