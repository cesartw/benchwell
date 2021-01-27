public class Benchwell.DatabaseService : Object {
	public Benchwell.Engine engine;
	public Benchwell.Connection? connection;
	public Benchwell.ConnectionInfo? info;
	public string dbname;
	public Benchwell.TableDef? table_def;
	public Benchwell.TableDef[]? tables;
	public Benchwell.ColDef[]? columns;
	public List<List<string?>> data;

	public DatabaseService () {
		engine = new Benchwell.Engine ();
	}

	public async void connect (Benchwell.ConnectionInfo _info) throws Benchwell.Error {
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

	public void use_database (string _dbname) throws Benchwell.Error {
		dbname = _dbname;

		connection.use_database (dbname);
		var _tables = connection.tables ();
		foreach (var q in info.queries) {
			var t = new Benchwell.TableDef.with_name (q.name);
			t.ttype = Benchwell.TableType.Dummy;
			t.source = q;
			_tables += t;
		}
		tables = _tables;
	}

	public void delete_table (Benchwell.TableDef tabledef) throws Benchwell.Error {
		connection.delete_table (tabledef);

		Benchwell.TableDef[] new_tables = {};

		foreach (var t in tables) {
			if (t.name != tabledef.name) {
				new_tables += t;
			}
		}

		tables = new_tables;
	}

	public void load_table (Benchwell.CondStmt[]? conditions,
							Benchwell.SortOption[]? sorts,
							int page, int page_size) throws Benchwell.Error {
		columns = connection.table_definition (table_def.name);
		data = connection.fetch_table (table_def.name,
										 conditions,
										 sorts,
										 page_size, page*page_size);
	}
}
