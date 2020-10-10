namespace Benchwell.SQL {
	public class Engine {
		public Connection? connect (ConnectionInfo info) throws Error {
			switch (info.adapter) {
				case "mysql":
					var driver = new MysqlDB ();
					return driver.connect (info);
			}

			return null;
		}
	}
}
