public class Benchwell.Services.Database : Object {
	public Benchwell.Backend.Sql.Engine engine;
	public Benchwell.Backend.Sql.Connection? connection;
	public Benchwell.Backend.Sql.ConnectionInfo? info;
	public string dbname;
	public Benchwell.Backend.Sql.TableDef? table_def;
	public Benchwell.Backend.Sql.TableDef[]? tables;
	public Benchwell.Backend.Sql.ColDef[]? columns;
	public List<List<string?>> data;

	public Database () {
		engine = new Benchwell.Backend.Sql.Engine ();
	}

	public async void connect (Benchwell.Backend.Sql.ConnectionInfo _info) throws Benchwell.Backend.Sql.Error {
		info = _info;

		yield Config.ping_dbus ();

		if (info.password == "") {
			var loop = new MainLoop ();
			var password = "";
			Config.decrypt.begin (info, (obj, res) => {
				password = Config.decrypt.end (res);
				loop.quit ();
			});
			loop.run ();

			info.password = password;
		}

		connection = engine.connect (info);
	}

	public void use_database (string _dbname) throws Benchwell.Backend.Sql.Error {
		dbname = _dbname;

		connection.use_database (dbname);
		var _tables = connection.tables ();
		foreach (var q in info.queries) {
			var t = new Benchwell.Backend.Sql.TableDef.with_name (q.name);
			t.ttype = Benchwell.Backend.Sql.TableType.Dummy;
			_tables += t;
		}
		tables = _tables;
	}

	public void delete_table (Benchwell.Backend.Sql.TableDef tabledef) throws Benchwell.Backend.Sql.Error {
		connection.delete_table (tabledef);

		Benchwell.Backend.Sql.TableDef[] new_tables = {};

		foreach (var t in tables) {
			if (t.name != tabledef.name) {
				new_tables += t;
			}
		}

		tables = new_tables;
	}

	public void load_table (Benchwell.Backend.Sql.CondStmt[] conditions,
							Benchwell.Backend.Sql.SortOption[] sorts,
							int page, int page_size) throws Benchwell.Backend.Sql.Error {
		columns = connection.table_definition (table_def.name);
		data = connection.fetch_table (table_def.name,
										 conditions,
										 sorts,
										 page_size, page*page_size);
	}
}
