public class Benchwell.SQL.MysqlDB : Benchwell.SQL.Driver {
	public Benchwell.SQL.Connection connect (Benchwell.SQL.ConnectionInfo c) throws Benchwell.SQL.ErrorConnection {
		return new MysqlConnection (c);
	}

	public static bool validate_connection (Benchwell.SQL.ConnectionInfo c) {
		if ( c == null ) {
			return false;
		}

		switch (c.ttype) {
			case "tcp":
				if (c.name == "") {
					return false;
				}
				if (c.host == "") {
					return false;
				}
				if (c.port == 0) {
					return false;
				}
				if (c.user == "") {
					return false;
				}
				if (c.adapter == "") {
					return false;
				}
				return true;
			case "socket":
				if (c.name == "") {
					return false;
				}
				if (c.socket == "") {
					return false;
				}
				if (c.user == "") {
					return false;
				}
				if (c.adapter == "") {
					return false;
				}
				return true;
		}

		return false;
	}
}

public class Benchwell.SQL.MysqlConnection : Benchwell.SQL.Connection, Object {
	private Mysql.Database db;
	public Benchwell.SQL.ConnectionInfo info;

	public MysqlConnection (Benchwell.SQL.ConnectionInfo c) throws Benchwell.SQL.ErrorConnection {
		info = c;
		db = new Mysql.Database ();

		Mysql.ClientFlag cflag = Mysql.ClientFlag.MULTI_STATEMENTS;
		var isConnected = db.real_connect (info.host, info.user, info.password, info.database, info.port, null, cflag);
		if ( ! isConnected ) {
			throw new Benchwell.SQL.ErrorConnection.CODE_1(@"$(db.errno()): $(db.error())");
		}
	}

	public List<string> databases () throws Benchwell.SQL.ErrorQuery {
		var result = db.list_dbs ();
		var databases = new List<string> ();

		string[] row;
		while ( ( row = result.fetch_row () ) != null ) {
			databases.append ( row[0] );
		}

		return databases;
	}

	public void use_database (string name) throws Benchwell.SQL.ErrorQuery {
		var rc = db.select_db (name);
		if ( rc != 0 ) {
			throw new Benchwell.SQL.ErrorQuery.CODE_1("failed query");
		}
	}

	public bool disconnect () {
		return true;
	}

	public void reconnect () throws Benchwell.SQL.ErrorConnection {
	}

	public TableDef[] tables () throws Benchwell.SQL.ErrorQuery {
		var result = db.list_tables ();

		string[] row;
		TableDef[] tables = {};
		while ( ( row = result.fetch_row () ) != null ) {
			var def = new TableDef () {
				name = row[0]
			};

			if ( row[0] == "VIEW" ) {
				def.ttype = Benchwell.SQL.TableType.View;
			}

			if ( row[0] == "BASE TABLE" ) {
				def.ttype = Benchwell.SQL.TableType.Regular;
			}

			tables += def;
		}

		return tables;
	}

	public owned ColDef[] table_definition(string name) throws Benchwell.SQL.ErrorQuery {
		var query = @"DESCRIBE $name";

		string[] row;

		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new ErrorQuery.CODE_1("failed query");
		}

		Benchwell.SQL.ColDef[] cols = {};
		var result = db.use_result ();
		while ( ( row = result.fetch_row () ) != null ) {
			Benchwell.SQL.ColDef col = new Benchwell.SQL.ColDef ();
			col.name = row[0];
			col.nullable = row[2] == "YES";
			col.pk = row[3] == "PRI";

			var coltype = Benchwell.SQL.ColType.String;
			int precision = 0;
			bool unsigned = false;
			string[] options = null;
			parse_type(row[1],  ref coltype, ref precision, ref options, ref unsigned);

			col.precision = precision;
			col.unsigned = unsigned;
			col.values = options;
			col.ttype = coltype;

			cols += col;
		}

		return cols;
	}

	public void delete_table(Benchwell.SQL.TableDef def) throws Benchwell.SQL.ErrorQuery {
		switch ( def.ttype ) {
			case Benchwell.SQL.TableType.Dummy:
				// TODO: delete from config;
				break;
			default:
				var query = @"DROP TABLE $(def.name)";
				var rc = db.query(query);
				if ( rc != 0 ) {
					throw new Benchwell.SQL.ErrorQuery.CODE_1("failed to drop table");
				}
				break;
		}
	}

	public void truncate_table(Benchwell.SQL.TableDef def) throws Benchwell.SQL.ErrorQuery {
		switch ( def.ttype ) {
			case Benchwell.SQL.TableType.Dummy:
				// TODO: delete from config;
				break;
			default:
				var query = @"TRUNCATE TABLE $(def.name)";
				var rc = db.query(query);
				if ( rc != 0 ) {
					throw new Benchwell.SQL.ErrorQuery.CODE_1("failed to truncate table");
				}
				break;
		}
	}

	public List<List<string?>> fetch_table (
		string name,
		Benchwell.SQL.CondStmt[]? conditions,
		Benchwell.SQL.SortOption[]? sorts,
		int limit,
		int offset
	) throws Benchwell.SQL.ErrorQuery {
		string[] wheres = {};
		int i = 0;

		foreach (Benchwell.SQL.CondStmt cond in conditions) {
			if ( cond.field.name == "" ) {
				continue;
			}

			switch ( cond.op ){
				case Benchwell.SQL.Operator.IsNull:
					wheres += @"`$(cond.field.name)` IS NULL";
					break;
				case Benchwell.SQL.Operator.IsNotNull:
					wheres += @"`$(cond.field.name)` IS NOT NULL";
					break;
				case Benchwell.SQL.Operator.Nin:
					var values = @"$(string.joinv ("\",\"", cond.val.split(",")))";
					wheres += @"`$(cond.field.name)` NOT IN $(values)";
					break;
				case Benchwell.SQL.Operator.In:
					var values = @"$(string.joinv ("\",\"", cond.val.split(",")))";
					wheres += @"`$(cond.field.name)` IN $(values)";
					break;
				default:
					// TODO: fix. This probably doesn't work
					var chunk = "";
					db.real_escape_string (chunk, cond.val, cond.val.length);
					wheres += @"`$(cond.field.name)` $(cond.op) '$(chunk)'";
					break;
			}
			i++;
		};

		string whereStmt = "";
		if ( wheres.length > 0 ) {
			whereStmt = @"WHERE $(string.joinv (" AND ", wheres))";
		}

		string[] _sorts = {};
		i = 0;
		foreach (Benchwell.SQL.SortOption sort in sorts) {
			var dir = "ASC";
			if ( sort.dir == Benchwell.SQL.SortType.Desc ) {
				dir = "DESC";
			}
			_sorts += @"`$(name)`.`$(sort.column.name)` $(dir)";
		};

		string sortStmt = "";
		if ( _sorts.length > 0 ) {
			sortStmt = @"ORDER BY $(string.joinv (", ", _sorts))";
		}

		var query = @"SELECT * FROM $name $whereStmt $sortStmt LIMIT $(limit) OFFSET $(offset)";

		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.SQL.ErrorQuery.CODE_1 (db.error());
		}

		var rows = new List<List<string?>> ();
		var result = db.use_result ();
		string[] row;
		while ((row = result.fetch_row () ) != null) {
			List<string> rowl = null;
			foreach (string s in row) {
				rowl.append (s);
			}

			rows.append ((owned) rowl);
		}

		return rows;
	}

	public void update_field (string table, ColDef[] columns, string[] row) throws Benchwell.SQL.ErrorQuery
		requires(columns.length == row.length)
		requires(columns.length > 1)
	{
		string[] wheres = {};
		for (var i = 0; i < columns.length - 1; i++) {
			wheres += @"`$(columns[i].name)` = '$(row[i])'";
		}

		var query = @"UPDATE `$table` SET `$(columns[columns.length -1].name)` = '$(row[row.length - 1])' WHERE $(string.joinv (" AND ", wheres))";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.SQL.ErrorQuery.CODE_1 (db.error());
		}
	}

	public string[] insert_record(string name, ColDef[] columns, string[] row) throws ErrorQuery
		requires (name != "")
		requires (row.length > 0)
		requires (columns.length == row.length)
	{
		var names = new string[row.length];
		var values = new string[row.length];
		for (var i = 0; i < row.length; i++) {
			names[i] = @"`$(columns[i].name)`";
			if (row[i] == Benchwell.null_string) {
				values[i] = "NULL";
			} else {
				values[i] = @"\"$(row[i])\"";
			}
		}

		var query = @"INSERT INTO `$name`($(string.joinv (", ", names))) VALUES ($(string.joinv (", ", values)))";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.SQL.ErrorQuery.CODE_1 (db.error());
		}

		var id = db.insert_id ();

		string? pk = null;
		foreach (var column in columns) {
			if (column.pk) {
				pk = column.name;
				break;
			}
		}

		if (pk == null || pk == "") {
			return row;
		}

		query = @"SELECT * FROM `$name` WHERE `$pk` = $id";

		rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.SQL.ErrorQuery.CODE_1 (db.error());
		}

		var result = db.use_result ();
		while ((row = result.fetch_row () ) != null) {
			return row;
		}

		return null;
	}

	private void parse_type(
		string t,
		ref Benchwell.SQL.ColType coltype,
		ref int precision,
		ref string[] options,
		ref bool unsigned
	) throws GLib.RegexError {
		var regex = new Regex ("([a-z ]+)(\\((.+)\\))?\\s?(unsigned)?");
		var parts = regex.split (t);

		var tt = parts[1]; // type
		var s = parts[3];  // precision

		coltype = Benchwell.SQL.ColType.String;

		if ( parts.length >= 5 ) {
			unsigned = parts[4] == "unsigned"; // unsigned
		}

		switch (tt) {
			case "enum":
				coltype = Benchwell.SQL.ColType.List;
				options = s.split (",");
				unsigned = false;
				break;
			case "text", "mediumtext", "longtext", "blob", "mediumblob", "longblob":
				coltype = Benchwell.SQL.ColType.String;
				unsigned = false;
				break;
			case "varchar", "tinytext":
				coltype = Benchwell.SQL.ColType.String;
				precision = int.parse(s);
				break;
			case "int", "smallint", "mediumint", "bigint":
				coltype = Benchwell.SQL.ColType.Int;
				precision = int.parse(s);
				break;
			case "tinyint":
				if ( s == "1" ) {
					coltype = Benchwell.SQL.ColType.Boolean;
					break;
				}

				precision = int.parse(s);
				coltype = Benchwell.SQL.ColType.Int;
				break;
			case "double precision", "double", "float", "decimal":
				coltype = Benchwell.SQL.ColType.Float;
				break;
			case "time", "datetime":
				coltype = Benchwell.SQL.ColType.Date;
				break;
		}
	}
}
