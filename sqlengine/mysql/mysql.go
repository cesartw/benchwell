package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	// mysql implementation
	_ "github.com/go-sql-driver/mysql"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/sqlengine/driver"
)

type mysqlDriver struct {
	cfgCon config.Connection
}

type mysqlConn struct {
	cfgCon    config.Connection
	database  string
	driver    *mysqlDriver
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

func (d *mysqlDriver) Connect(ctx context.Context, cfg config.Connection) (driver.Connection, error) {
	d.cfgCon = cfg
	return d.connect(ctx)
}

func (d *mysqlDriver) dsn() string {
	colonS := strings.Split(d.cfgCon.GetDSN(), ":")
	dsn := strings.TrimPrefix(d.cfgCon.GetDSN(), colonS[0]+"://")

	if strings.Index(dsn, "?") == -1 {
		dsn += "?"
	} else {
		dsn += "&"
	}
	dsn += "multiStatements=true"

	return dsn
}

func (d *mysqlDriver) connect(ctx context.Context) (*mysqlConn, error) {
	t, ok := ctx.Deadline()
	if !ok {
		t = time.Now().Add(time.Minute)
	}

	var db *sql.DB

	c := make(chan error, 1)
	go func() {
		var err error
		db, err = sql.Open("mysql", d.dsn())
		if err != nil {
			c <- err
			return
		}

		err = db.Ping()
		if err != nil {
			c <- err
			return
		}
		c <- nil
	}()

	select {
	case <-time.After(time.Until(t)):
		return nil, errors.New("context timeout")
	case <-ctx.Done():
		return nil, errors.New("context done")
	case err := <-c:
		if err != nil {
			return nil, err
		}
		close(c)
		return &mysqlConn{cfgCon: d.cfgCon, driver: d, db: db}, nil
	}
}

func (d *mysqlDriver) useDatabase(ctx context.Context, dbName string) (*sql.DB, error) {
	db, err := d.connect(ctx)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("USE %s", dbName)
	driver.Log(ctx, query)
	_, err = db.db.ExecContext(ctx, query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return nil, err
	}

	return db.db, nil
}

func (c *mysqlConn) Reconnect(ctx context.Context) error {
	c.db.Close()

	db, err := sql.Open("mysql", c.driver.dsn())
	if err != nil {
		return err
	}

	c.db = db

	return nil
}

func (c *mysqlConn) UseDatabase(ctx context.Context, db string) (driver.Database, error) {
	sqldb, err := c.driver.useDatabase(ctx, db)
	if err != nil {
		driver.Log(ctx, err.Error())
		return nil, err
	}

	return &mysqlDb{db: sqldb, name: db}, nil
}

// Disconnect ...
func (c *mysqlConn) Disconnect(ctx context.Context) error {
	return c.db.Close()
}

// LastError ...
func (c *mysqlConn) LastError() error {
	return c.lastError
}

// Databases ...
func (c *mysqlConn) Databases(ctx context.Context) ([]string, error) {
	query := "SHOW databases"
	driver.Log(ctx, query)
	rows, err := c.db.Query(query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return nil, err
	}

	dbs := make([]string, 0)
	for rows.Next() {
		var n string
		err := rows.Scan(&n)
		if err != nil {
			return nil, err
		}

		dbs = append(dbs, n)
	}

	return dbs, nil
}

func (d *mysqlDb) Name() string {
	return d.name
}

func (d *mysqlDb) Tables(ctx context.Context) ([]driver.TableDef, error) {
	query := "SHOW FULL TABLES"
	driver.Log(ctx, query)
	rows, err := d.db.Query(query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return nil, err
	}

	tables := []driver.TableDef{}
	for rows.Next() {
		var name, tableT string
		if err := rows.Scan(&name, &tableT); err != nil {
			return nil, err
		}

		def := driver.TableDef{}
		def.Name = name
		switch tableT {
		case "BASE TABLE":
			def.Type = driver.TableTypeRegular
		case "VIEW":
			def.Type = driver.TableTypeView
		default:
			def.Type = driver.TableTypeRegular
		}

		tables = append(tables, def)
	}

	return tables, nil
}

func (d *mysqlDb) TableDefinition(ctx context.Context, tableName string) ([]driver.ColDef, error) {
	query := "DESCRIBE " + tableName
	//driver.Log(ctx, query)

	sqlRows, err := d.db.Query(query)
	if err != nil {
		driver.Log(ctx, err.Error())
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

func (d *mysqlDb) Query(ctx context.Context, query string) (columnNames []string, data [][]interface{}, err error) {
	driver.Log(ctx, query)

	var sqlRows *sql.Rows
	sqlRows, err = d.db.Query(query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return nil, nil, err
	}

	defer sqlRows.Close()

	columnNames, err = sqlRows.Columns()
	if err != nil {
		return nil, nil, err
	}

	// insert or update
	if len(columnNames) == 0 {
		sqlRows.Close()
		sqlRows, err = d.db.Query("SELECT ROW_COUNT() AS affected_rows, LAST_INSERT_ID() AS last_inserted_id")
		if err != nil {
			return nil, nil, err
		}
		defer sqlRows.Close()
	}

	return d.loadData(sqlRows)
}

func (c *mysqlDb) loadData(sqlRows *sql.Rows) ([]string, [][]interface{}, error) {
	columnNames, err := sqlRows.Columns()
	if err != nil {
		return nil, nil, err
	}

	data := make([][]interface{}, 0)
	for sqlRows.Next() {
		row := make([]interface{}, len(columnNames))

		for ci := range columnNames {
			row[ci] = &row[ci]
		}

		if err := sqlRows.Scan(row...); err != nil {
			return nil, nil, err
		}

		for i, col := range row {
			if b, ok := col.([]byte); ok {
				row[i] = string(b)
			}
			if b, ok := col.([]uint8); ok {
				row[i] = string(b)
			}
		}

		data = append(data, row)
	}

	if err := sqlRows.Err(); err != nil {
		return nil, nil, err
	}

	return columnNames, data, nil
}

func (d *mysqlDb) Execute(ctx context.Context, query string) (string, int64, error) {
	driver.Log(ctx, query)

	result, err := d.db.Exec(query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return "", 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return "", 0, err
	}

	return fmt.Sprintf("%d", id), count, nil
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
	case driver.ColTypeBoolean:
		return strings.EqualFold(value, "true") || value == "1"
	case driver.ColTypeFloat:
		v, _ := strconv.ParseFloat(value, 64)
		return v
	case driver.ColTypeInt:
		v, _ := strconv.ParseInt(value, 10, 64)
		return v
	}

	return value
}

var typerg = regexp.MustCompile(`([a-z ]+)(\((.+)\))?\s?(unsigned)?`)

func (d *mysqlDb) parseType(mysqlStringType string) (driver.ColType, int, []string, bool) {
	matches := typerg.FindStringSubmatch(mysqlStringType)
	t := matches[1] // type
	s := matches[3] // size/precision
	u := matches[4] // unsigned

	switch t {
	case "enum":
		return driver.ColTypeList, 0, strings.Split(s, ","), false
	case "text", "mediumtext", "longtext", "blob", "mediumblob", "longblob":
		return driver.ColTypeLongString, 0, nil, false
	case "varchar", "tinytext":
		si, _ := strconv.Atoi(s)
		return driver.ColTypeString, si, nil, false
	case "int", "smallint", "mediumint", "bigint":
		si, _ := strconv.Atoi(s)
		return driver.ColTypeInt, si, nil, u == "unsigned"
	case "tinyint":
		if s == "1" {
			return driver.ColTypeBoolean, 0, nil, false
		}

		si, _ := strconv.Atoi(s)
		return driver.ColTypeInt, si, nil, u == "unsigned"
	case "double precision", "double", "float", "decimal":
		return driver.ColTypeFloat, 0, nil, u == "unsigned"
	case "time", "datetime":
		return driver.ColTypeDate, 0, nil, false
	}

	return driver.ColTypeString, 0, nil, true
}

type FetchTableOptions []driver.SortOption

func (o FetchTableOptions) SQL(tableName string) string {
	if len(o) == 0 {
		return ""
	}

	orderby := []string{}
	for _, sort := range o {
		s := fmt.Sprintf("`%s`.`%s` ", tableName, sort.Column.Name)
		if sort.Direction == driver.SortDirectionAsc {
			s += "DESC"
		} else {
			s += "ASC"
		}

		orderby = append(orderby, s)
	}

	return "ORDER BY " + strings.Join(orderby, ", ")
}

func (d *mysqlDb) FetchTable(
	ctx context.Context,
	tableName string,
	opts driver.FetchTableOptions,
) (
	colDef []driver.ColDef,
	rows [][]interface{},
	err error,
) {
	var sqlRows *sql.Rows

	wheres := []string{}
	for _, cond := range opts.Conditions {
		if cond.Op == "" || cond.Field.Name == "" {
			continue
		}
		//args = append(args, cond.Value)
		switch cond.Op {
		case driver.IsNull:
			wheres = append(wheres, fmt.Sprintf("`%s` IS NULL", cond.Field.Name))
		case driver.IsNotNull:
			wheres = append(wheres, fmt.Sprintf("`%s` IS NOT NULL", cond.Field.Name))
		case driver.Nin:
			v := []string{}
			for _, i := range strings.Split(cond.Value.(string), ",") {
				v = append(v, fmt.Sprintf("%#v", i))
			}
			wheres = append(wheres, fmt.Sprintf("`%s` NOT IN (%s)",
				cond.Field.Name, strings.Join(v, ", ")))
		case driver.In:
			v := []string{}
			for _, i := range strings.Split(cond.Value.(string), ",") {
				v = append(v, fmt.Sprintf("%#v", i))
			}
			wheres = append(wheres, fmt.Sprintf("`%s` IN (%s)",
				cond.Field.Name, strings.Join(v, ", ")))
		default:
			wheres = append(wheres, fmt.Sprintf("`%s` %s %#v",
				cond.Field.Name, string(cond.Op), cond.Value))
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = "WHERE " + strings.Join(wheres, " AND ")
	}
	where = strings.Replace(where, "%", "%%", -1)

	//args = append(args, opts.Offset, opts.Limit)

	query := fmt.Sprintf(`SELECT * FROM %s %s %s LIMIT %d, %d`,
		tableName, where, FetchTableOptions(opts.Sort).SQL(tableName), opts.Offset, opts.Limit)

	driver.Log(ctx, query)

	sqlRows, err = d.db.Query(query)
	if err != nil {
		driver.Log(ctx, err.Error())
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
		for i, col := range row {
			if b, ok := col.([]byte); ok {
				row[i] = string(b)
			}
			if b, ok := col.([]uint8); ok {
				row[i] = string(b)
			}
		}

		rows = append(rows, row)
	}
	if err := sqlRows.Err(); err != nil {
		return nil, nil, err
	}

	colDef, err = d.TableDefinition(ctx, tableName)
	if err != nil {
		return nil, nil, err
	}

	return colDef, rows, err
}

func (d *mysqlDb) DeleteRecord(ctx context.Context, tableName string, cols []driver.ColDef, args []interface{}) error {
	if len(cols) == 0 {
		return errors.New("cols is empty")
	}

	wheres := []string{}
	for i := range cols {
		wheres = append(wheres, fmt.Sprintf("`%s` = %#v", cols[i].Name, args[i]))
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE %s`, tableName, strings.Join(wheres, " AND "))

	driver.Log(ctx, query)

	_, err := d.db.Exec(query)
	if err != nil {
		driver.Log(ctx, err.Error())
	}

	return err
}

func (d *mysqlDb) UpdateRecord(
	ctx context.Context,
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
			ID = oldValues[i]
			if reflect.DeepEqual(oldValues[i], values[i]) {
				continue
			}
		}
		sets = append(sets, fmt.Sprintf("%s = ?", cols[i].Name))
		args = append(args, values[i])
	}

	query := `UPDATE %s SET %s WHERE %s = ?`

	query = fmt.Sprintf(query, tableName, strings.Join(sets, ", "), pk.Name, ID)
	//args = append(args, ID)

	driver.Log(ctx, query)

	result, err := d.db.Exec(query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (d *mysqlDb) UpdateField(
	ctx context.Context,
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
	wheres := []string{}

	for i := 0; i <= len(cols)-2; i++ {
		wheres = append(wheres, fmt.Sprintf("`%s` = %#v", cols[i].Name, values[i]))
		//args = append(args, values[i])
	}

	query := fmt.Sprintf("UPDATE `%s` SET `%s` = %#v WHERE %s",
		tableName,
		cols[lastIndex].Name,
		values[lastIndex],
		strings.Join(wheres, " AND "))

	driver.Log(ctx, query)

	result, err := d.db.Exec(query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (d *mysqlDb) UpdateFields(
	ctx context.Context,
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
			wheres = append(wheres, fmt.Sprintf("`%s` = %#v", cols[i].Name, values[i]))
		} else {
			sets = append(sets, fmt.Sprintf("`%s` = %#v", cols[i].Name, values[i]))
		}
	}

	if len(sets) == 0 {
	}

	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s",
		tableName,
		sets,
		strings.Join(wheres, " AND "))

	driver.Log(ctx, query)
	result, err := d.db.Exec(query)
	if err != nil {
		driver.Log(ctx, err.Error())
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (d *mysqlDb) InsertRecord(
	ctx context.Context,
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
		args[i] = d.ParseValue(cols[i], values[i].(string))
		qm[i] = "?"
	}

	query := `INSERT INTO %s(%s) VALUES (%s)`

	query = fmt.Sprintf(query, tableName, strings.Join(collist, ","), strings.Join(qm, ","))

	driver.Log(ctx, query)
	result, err := d.db.Exec(query, args...)
	if err != nil {
		driver.Log(ctx, err.Error())
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return d.fetchRecord(tableName, cols, id)
}

func (d *mysqlDb) GetCreateTable(
	ctx context.Context, tableName string,
) (
	string, error,
) {
	var (
		sqlRows *sql.Rows
		err     error
	)
	// NOTE: ? doesn't work here
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	driver.Log(ctx, query)
	sqlRows, err = d.db.Query(query)
	if err != nil {
		return "", err
	}

	defer sqlRows.Close()

	columns, err := sqlRows.Columns()
	if err != nil {
		return "", err
	}

	rows := make([][]interface{}, 0)

	for sqlRows.Next() {
		row := make([]interface{}, len(columns))

		for ci := range columns {
			row[ci] = &row[ci]
		}

		if err := sqlRows.Scan(row...); err != nil {
			return "", err
		}

		rows = append(rows, row)
	}
	if err := sqlRows.Err(); err != nil {
		return "", err
	}

	return string(rows[0][1].([]uint8)), err
}

func (d *mysqlDb) GetInsertStatement(
	ctx context.Context,
	tableName string,
	cols []driver.ColDef,
	v []interface{},
) (string, error) {
	if len(cols) == 0 {
		return "", errors.New("empty")
	}

	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("INSERT INTO `%s` ", tableName))
	colnames := []string{}
	values := []string{}
	for i, col := range cols {
		pv := d.ParseValue(col, v[i].(string))

		colnames = append(colnames, "`"+col.Name+"`")

		switch t := pv.(type) {
		case bool:
			if t {
				values = append(values, "1")
			} else {
				values = append(values, "0")
			}
		case float64:
			values = append(values, fmt.Sprintf("%s", v[i].(string)))
		case int64:
			values = append(values, fmt.Sprintf("%s", v[i].(string)))
		case nil:
			values = append(values, "NULL")
		default:
			values = append(values, fmt.Sprintf("'%s'", v[i].(string)))
		}
	}

	b.WriteString("(" + strings.Join(colnames, ", ") + ") VALUES ")
	b.WriteString("(" + strings.Join(values, ", ") + ");")

	return b.String(), nil
}

func (d *mysqlDb) GetSelectStatement(
	ctx context.Context,
	table driver.TableDef,
) (string, error) {
	switch table.Type {
	case driver.TableTypeDummy:
		return table.Query, nil
	default:
		return fmt.Sprintf(`SELECT * FROM %s`, table.Name), nil
	}
}

func (d *mysqlDb) DeleteTable(
	ctx context.Context,
	table driver.TableDef,
) error {
	switch table.Type {
	case driver.TableTypeDummy:
		return nil
	default:
		query := fmt.Sprintf(`DROP TABLE %s`, table.Name)
		driver.Log(ctx, query)

		_, err := d.db.ExecContext(ctx, query)
		if err != nil {
			driver.Log(ctx, err.Error())
		}

		return nil
	}
}

func (d *mysqlDb) TruncateTable(
	ctx context.Context,
	table driver.TableDef,
) error {
	switch table.Type {
	case driver.TableTypeDummy:
		return nil
	default:
		query := fmt.Sprintf(`TRUNCATE TABLE %s`, table.Name)
		driver.Log(ctx, query)

		_, err := d.db.ExecContext(ctx, query)
		if err != nil {
			driver.Log(ctx, err.Error())
		}
		return nil
	}
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

	query := `SELECT * FROM %s WHERE %s = ?`

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
