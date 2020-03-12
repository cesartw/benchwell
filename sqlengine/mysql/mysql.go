package driver

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	// mysql implementation
	_ "github.com/go-sql-driver/mysql"

	"bitbucket.org/goreorto/sqlhero/sqlengine/driver"
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

		defs = append(defs, driver.ColDef{
			Name: name.String,
			PK:   key.String == "PRI",
		})
	}

	if err := sqlRows.Err(); err != nil {
		return nil, err
	}

	return defs, nil
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

func (d *mysqlDb) DeleteRecord(tableName string, defs []driver.ColDef, values []*string) error {
	var ID string
	var priCol driver.ColDef

	for i := range values {
		if defs[i].PK {
			ID = *values[i]
			priCol = defs[i]
			break
		}
	}

	if ID == "" {
		return errors.New("table doesn't have a primary key")
	}

	query := fmt.Sprintf(`
	DELETE FROM %s
	WHERE %s = ?
	`, tableName, priCol.Name)

	_, err := d.db.Exec(query, ID)

	return err
}

func (d *mysqlDb) UpdateRecord(
	tableName string,
	cols []driver.ColDef,
	values, oldValues []*string,
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
	var ID string

	for i := range values {
		if cols[i].PK {
			ID = *oldValues[i]
			if oldValues[i] == values[i] {
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

func (d *mysqlDb) InsertRecord(
	tableName string,
	cols []driver.ColDef,
	values []*string,
) ([]*string, error) {
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
	[]*string,
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

	rows := make([][]*string, 0)

	for sqlRows.Next() {
		row := make([]*string, len(columns))
		irow := make([]interface{}, len(columns))

		for ci := range columns {
			irow[ci] = &row[ci]
		}

		if err := sqlRows.Scan(irow...); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}
	if err := sqlRows.Err(); err != nil {
		return nil, err
	}

	return rows[0], err
}
