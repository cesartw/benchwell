using SQLEngine;

class MysqlDB : Driver {
	public Connection Connect (ConnectionInfo c) throws ErrorConnection {
		return new MysqlConnection (c);
	}

	public bool ValidateConnection (ConnectionInfo c) {
		return true;
	}
}

class MysqlConnection : Connection, Object {
	private Mysql.Database db;
	public ConnectionInfo info;

	public MysqlConnection (ConnectionInfo c) {
		info = c;
		db = new Mysql.Database ();

		Mysql.ClientFlag cflag = 0;
		var isConnected = db.real_connect (info.host, info.user, info.password, info.database, info.port, null, cflag);
		if ( ! isConnected ) {
			throw new ErrorConnection.CODE_1("not connected");
		}
	}

	public List<string> Databases () throws ErrorQuery {
		var result = db.list_dbs ();
		var databases = new List<string> ();

		string[] row;
		while ( ( row = result.fetch_row () ) != null ) {
			databases.append ( row[0] );
		}

		return databases;
	}

	public void UseDatabase (string name) throws ErrorQuery {
		var rc = db.select_db (name);
		if ( rc != 0 ) {
			throw new ErrorQuery.CODE_1("failed query");
		}
	}

	public bool Disconnect () {
		return true;
	}

	public void Reconnect () throws ErrorConnection {
	}

	public List<TableDef> Tables() throws ErrorQuery {
		var result = db.list_tables ();

		string[] row;
		var tables = new List<TableDef> ();
		while ( ( row = result.fetch_row () ) != null ) {
			var def = new TableDef () {
				name = row[0]
			};

			if ( row[0] == "VIEW" ) {
				def.ttype = TableType.View;
			}

			if ( row[0] == "BASE TABLE" ) {
				def.ttype = TableType.Regular;
			}

			tables.append ( def );
		}

		return tables;
	}

	public List<ColDef> TableDefinition(string name) throws ErrorQuery {
		var query = @"DESCRIBE $name";

		string[] row;

		var rc = db.query (query);
		if ( rc != 0 ) {
			throw new ErrorQuery.CODE_1("failed query");
		}

		var cols = new List<ColDef> ();
		var result = db.use_result ();
		while ( ( row = result.fetch_row () ) != null ) {
			var col = new ColDef () {
				name = row[0],
				nullable = row[2] == "YES",
				pk = row[3] == "PRI"
			};

			ColType coltype = ColType.String;
			int precision = 0;
			bool unsigned = false;
			string[] options = null;
			parseType(row[1],  ref coltype, ref precision, ref options, ref unsigned);

			col.precision = precision;
			col.unsigned = unsigned;
			col.values = options;
			col.ttype = coltype;

			cols.append (col);
		}

		return cols;
	}

	public void DeleteTable(TableDef def) throws ErrorQuery {
		switch ( def.ttype ) {
			case TableType.Dummy:
				// TODO: delete from config;
				break;
			default:
				var query = @"DROP TABLE $(def.name)";
				var rc = db.query(query);
				if ( rc != 0 ) {
					throw new ErrorQuery.CODE_1("failed to drop table");
				}
				break;
		}
	}

	public void TruncateTable(TableDef def) throws ErrorQuery {
		switch ( def.ttype ) {
			case TableType.Dummy:
				// TODO: delete from config;
				break;
			default:
				var query = @"TRUNCATE TABLE $(def.name)";
				var rc = db.query(query);
				if ( rc != 0 ) {
					throw new ErrorQuery.CODE_1("failed to truncate table");
				}
				break;
		}
	}

	public void FetchTable(
		string name, FetchTableOptions opts,
		ref List<ColDef> def,
		ref List<List<string>> rows
		) throws ErrorQuery {
	}

	private void parseType(string t, ref ColType coltype, ref int precision, ref string[] options, ref bool unsigned) {
		var regex = new Regex ("([a-z ]+)(\\((.+)\\))?\\s?(unsigned)?");
		var parts = regex.split (t);

		var tt = parts[1]; // type
		var s = parts[3];  // precision

		coltype = ColType.String;

		if ( parts.length >= 5 ) {
			unsigned = parts[4] == "unsigned"; // unsigned
		}

		switch (tt) {
			case "enum":
				coltype = ColType.List;
				options = s.split (",");
				unsigned = false;
				break;
			case "text", "mediumtext", "longtext", "blob", "mediumblob", "longblob":
				coltype = ColType.String;
				unsigned = false;
				break;
			case "varchar", "tinytext":
				coltype = ColType.String;
				precision = int.parse(s);
				break;
			case "int", "smallint", "mediumint", "bigint":
				coltype = ColType.Int;
				precision = int.parse(s);
				break;
			case "tinyint":
				if ( s == "1" ) {
					coltype = ColType.Boolean;
					break;
				}

				precision = int.parse(s);
				coltype = ColType.Int;
				break;
			case "double precision", "double", "float", "decimal":
				coltype = ColType.Float;
				break;
			case "time", "datetime":
				coltype = ColType.Date;
				break;
		}
	}
}

int main(string[] args) {
	var info = new ConnectionInfo ();
	info.host = "db";
	info.user = "root";
	info.password = "dev";
	info.database = "teamworkdesk_shard6";

	var c = new MysqlConnection (info);
	var dbs = c.Databases ();
	dbs.foreach ((entry) => {
		//print (@"$entry\n");
	});

	c.UseDatabase ("teamworkdesk_shard6");
	var tables = c.Tables();
	tables.foreach ((table) => {
		//var name = table.name;
		//print (@"$name\n");
	});

	var fields = c.TableDefinition ("threads");
	fields.foreach ((field) => {
		//print (@"$(field.name) ");
		//print (@"$(field.pk) ");
		//print (@"$(field.fk) ");
		//print (@"$(field.precision) ");
		//print (@"$(field.nullable) ");
		//print (@"$(field.ttype) ");
		//for ( int i = 0; i < field.values.length; i++ ) {
			//print (@"$(field.values[i])");
		//}
		//print("\n");
	});

	return 0;
}
