public class Benchwell.MysqlDB : Benchwell.Driver {
	public Benchwell.Connection connect (Benchwell.ConnectionInfo c) throws Benchwell.Error {
		return new MysqlConnection (c);
	}

	public static bool validate_connection (Benchwell.ConnectionInfo c) {
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

public class Benchwell.MysqlConnection : Benchwell.Connection, Object {
	private Mysql.Database db;
	public Benchwell.ConnectionInfo info;

	public MysqlConnection (Benchwell.ConnectionInfo c) throws Benchwell.Error {
		info = c;
		db = new Mysql.Database ();

		Mysql.ClientFlag cflag = Mysql.ClientFlag.MULTI_STATEMENTS;
		var isConnected = db.real_connect (info.host, info.user, info.password, info.database, info.port, null, cflag);
		if ( ! isConnected ) {
			throw new Benchwell.Error.CONNECTION (@"$(db.errno()): $(db.error())");
		}
		db.options (Mysql.Option.OPT_RECONNECT, "1");
	}

	public List<string> databases () throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION (@"connection lost");
		}

		var result = db.list_dbs ();
		var databases = new List<string> ();

		string[] row;
		while ( ( row = result.fetch_row () ) != null ) {
			databases.append ( row[0] );
		}

		return databases;
	}

	public void use_database (string name) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		var rc = db.select_db (name);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY("failed query");
		}
	}

	public new bool disconnect () {
		return true;
	}

	public void reconnect () throws Benchwell.Error {
	}

	public TableDef[] tables () throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.QUERY(@"connection lost");
		}

		var result = db.list_tables ();

		string[] row;
		TableDef[] tables = {};
		while ( ( row = result.fetch_row () ) != null ) {
			var def = new TableDef () {
				name = row[0]
			};

			if ( row[0] == "VIEW" ) {
				def.ttype = Benchwell.TableType.View;
			}

			if ( row[0] == "BASE TABLE" ) {
				def.ttype = Benchwell.TableType.Regular;
			}

			tables += def;
		}

		return tables;
	}

	public ColDef[] table_definition(string name) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		var query = @"DESCRIBE $name";

		string[] row;

		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY("failed query");
		}

		Benchwell.ColDef[] cols = {};
		var result = db.use_result ();
		while ( ( row = result.fetch_row () ) != null ) {
			Benchwell.ColDef col = new Benchwell.ColDef ();
			col.name = row[0];
			col.nullable = row[2] == "YES";
			col.pk = row[3] == "PRI";

			var coltype = Benchwell.ColType.String;
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

	public void delete_table(Benchwell.TableDef def) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		switch ( def.ttype ) {
			case Benchwell.TableType.Dummy:
				// TODO: delete from config;
				break;
			default:
				var query = @"DROP TABLE $(def.name)";
				var rc = db.query(query);
				if ( rc != 0 ) {
					throw new Benchwell.Error.QUERY("failed to drop table");
				}
				break;
		}
	}

	public void truncate_table(Benchwell.TableDef def) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		switch ( def.ttype ) {
			case Benchwell.TableType.Dummy:
				// TODO: delete from config;
				break;
			default:
				var query = @"TRUNCATE TABLE `$(def.name)`";
				var rc = db.query(query);
				if ( rc != 0 ) {
					throw new Benchwell.Error.QUERY("failed to truncate table");
				}
				break;
		}
	}

	public List<List<string?>> fetch_table (
		string name,
		Benchwell.CondStmt[]? conditions,
		Benchwell.SortOption[]? sorts,
		int limit,
		int offset
	) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		string[] wheres = {};
		int i = 0;

		foreach (Benchwell.CondStmt cond in conditions) {
			if ( cond.field.name == "" ) {
				continue;
			}

			switch ( cond.op ){
				case Benchwell.Operator.IsNull:
					wheres += @"`$(cond.field.name)` IS NULL";
					break;
				case Benchwell.Operator.IsNotNull:
					wheres += @"`$(cond.field.name)` IS NOT NULL";
					break;
				case Benchwell.Operator.Nin:
					var val = sanitize_string_array (cond.val);
					wheres += @"`$(cond.field.name)` NOT IN ($val)";
					break;
				case Benchwell.Operator.In:
					var val = sanitize_string_array (cond.val);
					wheres += @"`$(cond.field.name)` IN ($val)";
					break;
				default:
					var val = sanitize_string (cond.val);
					wheres += @"`$(cond.field.name)` $(cond.op) $val";
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
		foreach (Benchwell.SortOption sort in sorts) {
			var dir = "ASC";
			if ( sort.dir == Benchwell.SortType.Desc ) {
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
			throw new Benchwell.Error.QUERY (db.error());
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

	public void update_field (string table, ColDef[] columns, string[] row) throws Benchwell.Error
		requires(columns.length == row.length)
		requires(columns.length > 1)
	{
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		string[] wheres = {};
		for (var i = 0; i < columns.length - 1; i++) {
			var val = sanitize_string (row[i]);
			wheres += @"`$(columns[i].name)` = $val";
		}

		var new_value = sanitize_string (row[row.length - 1]);
		var query = @"UPDATE `$table` SET `$(columns[columns.length -1].name)` = $new_value WHERE $(string.joinv (" AND ", wheres))";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}
	}

	public string[]? insert_record(string name, ColDef[] columns, string[] row) throws Benchwell.Error
		requires (name != "")
		requires (row.length > 0)
		requires (columns.length == row.length)
	{
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		var builder = new StringBuilder ();
		builder.append ("INSERT INTO `");
		builder.append (name);
		builder.append ("`(");
		for (var i = 0; i < row.length; i++) {
			builder.append (columns[i].name);
			if (i != row.length -1) {
				builder.append (",");
			}
		}

		builder.append (") VALUES (");


		for (var i=0;i<row.length;i++) {
			builder.append(sanitize_string (row[i]));
			if (i != row.length -1) {
				builder.append (",");
			}
		}
		builder.append (")");

		var query = builder.str;

		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
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
			throw new Benchwell.Error.QUERY (db.error());
		}

		var result = db.use_result ();
		while ((row = result.fetch_row () ) != null) {
			return row;
		}

		return null;
	}

	public void delete_record(string name, ColDef[] columns, string[] row) throws Benchwell.Error
		requires (name != "")
		requires (row.length > 0)
		requires (columns.length == row.length)
	{
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION (@"connection lost");
		}

		string[] wheres = {};
		// delete using PK
		for (var i = 0; i < columns.length; i++) {
			if (!columns[i].pk) {
				continue;
			}

			var val = "";
			if (row[i] == Benchwell.null_string) {
				val = "IS NULL";
			} else {
				val = @"= \"$(row[i])\"";
			}

				wheres += @"`$(columns[i].name)` $val";
		}

		if (wheres.length == 0) {
			for (var i = 0; i < columns.length; i++) {
				string val = "";
				if (row[i] == Benchwell.null_string) {
					val = "IS NULL";
				} else {
					val = @"= $(row[i])";
				}

				wheres += @"";
				wheres += @"`$(columns[i].name)` $val";
			}
		}

		var query = @"DELETE FROM `$name` WHERE $(string.joinv (" AND ", wheres))";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}

		return;
	}

	public string get_create_table(string name) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		var query = @"SHOW CREATE TABLE `$name`";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}

		string[] row;
		var result = db.use_result ();
		while ((row = result.fetch_row () ) != null) {
			return row[1];
		}

		return "";
	}

	public void query(string query, out string[] columns, out List<List<string?>> rows) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION (@"connection lost");
		}

		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (@"$(db.errno()): $(db.error())");
		}

		var result = db.use_result ();
		Mysql.Field* field;
		rows = new List<List<string?>> ();

		// SELECT
		if (result != null) {
			string[] names = {};
			while ((field = result.fetch_field ()) != null) {
				names += field.name;
			}
			columns = names;

			string[] row;
			while ((row = result.fetch_row () ) != null) {
				List<string> rowl = null;
				foreach (string s in row) {
					rowl.append (s);
				}

				rows.append ((owned) rowl);
			}
		} else {
			// DML

			columns = {"affected rows", "inserted id"};
			var id = db.insert_id ();
			var count = db.affected_rows ();

			var rowl = new List<string> ();
			rowl.append (count.to_string ());
			rowl.append (id.to_string ());

			rows.append ((owned) rowl);
		}
	}

	private void parse_type(
		string t,
		ref Benchwell.ColType coltype,
		ref int precision,
		ref string[] options,
		ref bool unsigned
	) throws GLib.RegexError {
		var regex = new Regex ("([a-z ]+)(\\((.+)\\))?\\s?(unsigned)?");
		var parts = regex.split (t);

		var tt = parts[1]; // type
		var s = parts[3];  // precision

		coltype = Benchwell.ColType.String;

		if ( parts.length >= 5 ) {
			unsigned = parts[4] == "unsigned"; // unsigned
		}

		switch (tt) {
			case "enum":
				coltype = Benchwell.ColType.List;
				options = s.split (",");
				unsigned = false;
				break;
			case "text", "mediumtext", "longtext", "blob", "mediumblob", "longblob":
				coltype = Benchwell.ColType.String;
				unsigned = false;
				break;
			case "varchar", "tinytext":
				coltype = Benchwell.ColType.String;
				precision = int.parse(s);
				break;
			case "int", "smallint", "mediumint", "bigint":
				coltype = Benchwell.ColType.Int;
				precision = int.parse(s);
				break;
			case "tinyint":
				if ( s == "1" ) {
					coltype = Benchwell.ColType.Boolean;
					break;
				}

				precision = int.parse(s);
				coltype = Benchwell.ColType.Int;
				break;
			case "double precision", "double", "float", "decimal":
				coltype = Benchwell.ColType.Float;
				break;
			case "time", "datetime":
				coltype = Benchwell.ColType.Date;
				break;
		}
	}

	private string sanitize_string (string? dirty) {
		if (dirty == Benchwell.null_string || dirty == null) {
			return "NULL";
		}
		string chunk = dirty;
		db.real_escape_string (chunk, dirty, dirty.length);

		return @"\"$chunk\"";
	}

	private string sanitize_string_array (string dirty) {
		var parts = dirty.split (",");
		string[] clean = {};
		foreach (var part in parts) {
			if (part == null) {
				clean += "NULL";
				continue;
			}
			string chunk = "";
			db.real_escape_string (chunk, part, part.length);
			clean += @"\"$chunk\"";
		}

		return string.joinv (",", clean);
	}

	public string get_insert_statement(string name, Benchwell.ColDef[] columns, string[] row)
		requires(columns.length == row.length)
		requires(columns.length > 1)
	{
		var builder = new StringBuilder ();
		builder.append ("INSERT INTO `")
			.append (name)
			.append ("`(");

		for (var i = 0; i < columns.length; i++){
			builder.append ("`")
				.append (columns[i].name)
				.append ("`");
			if (i < columns.length -1) {
				builder.append (", ");
			}
		}

		builder.append(") VALUES(");
		for (var i = 0; i < row.length; i++){
			builder.append (sanitize_string (row[i]));
			if (i < row.length - 1) {
				builder.append (", ");
			}
		}
		builder.append (");");

		var s = builder.str;
		return s;
	}
}
