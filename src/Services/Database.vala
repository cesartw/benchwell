public class Benchwell.Service.Database {
	public weak Benchwell.Window window { get; construct; }
	private Benchwell.SQL.Engine engine { get; }

	public Database (Benchwell.Window widow) {
		Object(
			window: window
		);

		engine = new Benchwell.SQL.Engine ();
	}
}
