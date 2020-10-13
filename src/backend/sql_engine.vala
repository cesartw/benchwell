namespace Benchwell.Backend.Sql {
	public class Engine {
		public Benchwell.Backend.Sql.Connection? connect (Benchwell.Backend.Sql.ConnectionInfo info) throws Error {
			switch (info.adapter) {
				case "mysql":
					var driver = new MysqlDB ();
					return driver.connect (info);
			}

			return null;
		}
	}
}
