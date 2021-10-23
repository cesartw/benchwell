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
		var r = db.options (Mysql.Option.OPT_RECONNECT, "1");
		if (r != 0) {
			throw new Benchwell.Error.CONNECTION (@"unable to set RECONNECT");
		}
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

	public ColDef[] table_definition (string name) throws Benchwell.Error {
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

			var coltype = Benchwell.ColType.STRING;
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

	public void delete_table (Benchwell.TableDef def) throws Benchwell.Error {
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

	public void truncate_table (Benchwell.TableDef def) throws Benchwell.Error {
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
		int offset,
		out QueryInfo? query_info
	) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		string[] wheres = {};
		int i = 0;

		foreach (Benchwell.CondStmt cond in conditions) {
			if (cond.field.name == "" || !cond.enabled) {
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

		var then = get_real_time ();
		int64 row_count = 0;
		var query = @"SELECT * FROM $name $whereStmt $sortStmt LIMIT $(limit) OFFSET $(offset)";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}
		var now = get_real_time ();

		var rows = new List<List<string?>> ();
		var result = db.use_result ();

		string[] row;
		while ((row = result.fetch_row () ) != null) {
			row_count++;
			List<string> rowl = null;
			foreach (string s in row) {
				rowl.append (s);
			}

			rows.append ((owned) rowl);
		}

		query_info = new QueryInfo (now - then, row_count);
		return rows;
	}

	public void update_field (string table, Column[] columns) throws Benchwell.Error
		requires(columns.length > 1)
	{
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		string[] wheres = {};
		bool hasPk = false;
		for (var i = 0; i < columns.length - 1; i++) {
			if (!columns[i].coldef.pk) {
				continue;
			}
			hasPk = true;
			var val = sanitize_string (columns[i].val);
			wheres += @"`$(columns[i].coldef.name)` = $val";
		}

		if (!hasPk) {
			for (var i = 0; i < columns.length - 1; i++) {
				var val = sanitize_string (columns[i].val);
				if (val == "NULL") {
					wheres += @"`$(columns[i].coldef.name)` IS NULL";
				} else {
					wheres += @"`$(columns[i].coldef.name)` = $val";
				}
			}
		}

		var new_value = sanitize_string (columns[columns.length - 1].val);
		var query = @"UPDATE `$table` SET `$(columns[columns.length -1].coldef.name)` = $new_value WHERE $(string.joinv (" AND ", wheres))";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}
	}

	public void update_fields (string table, Column[] columns) throws Benchwell.Error
		requires(columns.length > 1)
		requires(columns.length % 2 == 0)
	{
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		string[] wheres = {};
		bool hasPk = false;
		for (var i = 0; i < columns.length / 2; i++) {
			if (!columns[i].coldef.pk) {
				continue;
			}
			hasPk = true;
			var val = sanitize_string (columns[i].val);
			wheres += @"`$(columns[i].coldef.name)` = $val";
		}

		if (!hasPk) {
			for (var i = 0; i < columns.length - 1; i++) {
				var val = sanitize_string (columns[i].val);
				if (val == "NULL") {
					wheres += @"`$(columns[i].coldef.name)` IS NULL";
				} else {
					wheres += @"`$(columns[i].coldef.name)` = $val";
				}
			}
		}

		string[] sets = {};
		for (var i = columns.length / 2; i < columns.length; i++) {
			var new_value = sanitize_string (columns[i].val);
			sets += @"`$(columns[i].coldef.name)` = $new_value";
		}
		var query = @"UPDATE `$table` SET $(string.joinv (",", sets)) WHERE $(string.joinv (" AND ", wheres))";
		return;
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}
	}

	public Benchwell.Column[]? insert_record (string name, Column[] columns) throws Benchwell.Error
		requires (name != "")
		requires (columns.length > 0)
	{
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION(@"connection lost");
		}

		var builder = new StringBuilder ();
		builder.append ("INSERT INTO `");
		builder.append (name);
		builder.append ("`(");
		for (var i = 0; i < columns.length; i++) {
			builder.append (columns[i].coldef.name);
			if (i != columns.length -1) {
				builder.append (",");
			}
		}

		builder.append (") VALUES (");


		for (var i=0;i<columns.length;i++) {
			builder.append(sanitize_string (columns[i].val));
			if (i != columns.length -1) {
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
			if (column.coldef.pk) {
				pk = column.coldef.name;
				break;
			}
		}

		if (pk == null || pk == "") {
			return columns;
		}

		query = @"SELECT * FROM `$name` WHERE `$pk` = $id";

		rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}

		var result = db.use_result ();
		string[] data;
		while ((data = result.fetch_row () ) != null) {
			for (var i = 0; i < data.length; i++) {
				columns[i].val = data[i];
			}
			return columns;
		}

		return null;
	}

	public void delete_record (string name, Column[] columns) throws Benchwell.Error
		requires (name != "")
		requires (columns.length > 0)
	{
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION (@"connection lost");
		}

		string[] wheres = {};
		// delete using PK
		for (var i = 0; i < columns.length; i++) {
			if (!columns[i].coldef.pk) {
				continue;
			}

			var val = "";
			if (columns[i].val == Benchwell.null_string) {
				val = "IS NULL";
			} else {
				val = @"= \"$(columns[i].val)\"";
			}

				wheres += @"`$(columns[i].coldef.name)` $val";
		}

		if (wheres.length == 0) {
			for (var i = 0; i < columns.length; i++) {
				string val = "";
				if (columns[i].val == Benchwell.null_string) {
					val = "IS NULL";
				} else {
					val = @"= $(columns[i].val)";
				}

				wheres += @"";
				wheres += @"`$(columns[i].coldef.name)` $val";
			}
		}

		var query = @"DELETE FROM `$name` WHERE $(string.joinv (" AND ", wheres))";
		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (db.error());
		}

		return;
	}

	public string get_create_table (string name) throws Benchwell.Error {
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

	public void query (string query, out string[] columns, out List<List<string?>> rows, out Benchwell.QueryInfo? query_info) throws Benchwell.Error {
		if (db.ping () != 0) {
			throw new Benchwell.Error.CONNECTION (@"connection lost");
		}

		var then = get_real_time ();
		int64 row_count = 0;

		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new Benchwell.Error.QUERY (@"$(db.errno()): $(db.error())");
		}
		var now = get_real_time ();

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
				row_count++;
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

		query_info = new QueryInfo (now - then, row_count);
	}

	private void parse_type (
		string t,
		ref Benchwell.ColType coltype,
		ref int precision,
		ref string[] options,
		ref bool unsigned
	) {
		Regex regex = null;
		try {
			regex = new Regex ("([a-z ]+)(\\((.+)\\))?\\s?(unsigned)?");
		} catch (RegexError err) {
		}

		var parts = regex.split (t);

		var tt = parts[1]; // type
		var s = parts[3];  // precision

		coltype = Benchwell.ColType.STRING;

		if ( parts.length >= 5 ) {
			unsigned = parts[4] == "unsigned"; // unsigned
		}

		switch (tt) {
			case "enum":
				coltype = Benchwell.ColType.LIST;
				options = s.split (",");
				unsigned = false;
				break;
			case "text", "mediumtext":
				coltype = Benchwell.ColType.STRING;
				unsigned = false;
				break;
			case "blob", "mediumblob", "longblob":
				coltype = Benchwell.ColType.BLOB;
				break;
			case "longtext":
				coltype = Benchwell.ColType.LONG_STRING;
				break;
			case "varchar", "tinytext":
				coltype = Benchwell.ColType.STRING;
				precision = int.parse(s);
				break;
			case "int", "smallint", "mediumint", "bigint":
				coltype = Benchwell.ColType.INT;
				precision = int.parse(s);
				break;
			case "tinyint":
				if ( s == "1" ) {
					coltype = Benchwell.ColType.BOOLEAN;
					break;
				}

				precision = int.parse(s);
				coltype = Benchwell.ColType.INT;
				break;
			case "double precision", "double", "float", "decimal":
				coltype = Benchwell.ColType.FLOAT;
				break;
			case "time", "datetime":
				coltype = Benchwell.ColType.DATE;
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

	public string get_insert_statement(string name, Benchwell.Column[] columns)
		requires(columns.length > 1)
	{
		var builder = new StringBuilder ();
		builder.append ("INSERT INTO `")
			.append (name)
			.append ("`(");

		for (var i = 0; i < columns.length; i++){
			builder.append ("`")
				.append (columns[i].coldef.name)
				.append ("`");
			if (i < columns.length -1) {
				builder.append (", ");
			}
		}

		builder.append(") VALUES(");
		for (var i = 0; i < columns[i].val.length; i++){
			builder.append (sanitize_string (columns[i].val));
			if (i < columns[i].val.length - 1) {
				builder.append (", ");
			}
		}
		builder.append (");");

		var s = builder.str;
		return s;
	}
}
