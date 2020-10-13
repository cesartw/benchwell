public class Benchwell.Services.Database : Object {
	public Benchwell.Backend.Sql.Engine engine;
	public Benchwell.Backend.Sql.Connection? connection;
	public Benchwell.Backend.Sql.ConnectionInfo? info;
	public Benchwell.Backend.Sql.TableDef? table_def;

	public Database () {
		engine = new Benchwell.Backend.Sql.Engine ();
	}

	public void connect (Benchwell.Backend.Sql.ConnectionInfo _info) throws Benchwell.Backend.Sql.Error {
		info = _info;

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
}
