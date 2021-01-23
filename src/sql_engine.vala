public class Benchwell.Engine {
	public Benchwell.Connection? connect (Benchwell.ConnectionInfo info) throws Error {
		switch (info.adapter) {
			case "mysql":
				var driver = new MysqlDB ();
				return driver.connect (info);
		}

		return null;
	}
}
