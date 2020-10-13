namespace Benchwell {
	public static int main(string[] args) {
		var info = new Benchwell.SQLEngine.ConnectionInfo ();
		info.host = "db";
		info.user = "root";
		info.password = "dev";
		info.database = "teamworkdesk_shard6";

		var c = new Benchwell.SQLEngine.MysqlConnection (info);
		var dbs = c.Databases ();
		dbs.foreach ((entry) => {
			print (@"Database: $entry\n");
		});

		c.UseDatabase ("teamworkdesk_shard6");
		var tables = c.Tables ();
		tables.foreach ((table) => {
			var name = table.name;
			print (@"	Table: $name\n");
		});

		var fields = c.TableDefinition ("customers");
		fields.foreach ((field) => {
			print ("		field:");
			print (@"name: $(field.name) ");
			print (@"pk: $(field.pk) ");
			print (@"fk: $(field.fk) ");
			print (@"precision: $(field.precision) ");
			print (@"nullable: $(field.nullable) ");
			print (@"type: $(field.ttype) ");
			for ( int i = 0; i < field.values.length; i++ ) {
				print (@"opts: $(field.values[i])");
			}
			print("\n");
		});


		var rows = new List<List<string?>>();
		var opts = new Benchwell.SQLEngine.FetchTableOptions();
		c.FetchTable("customers", opts, ref rows);

		rows.foreach((row) => {
			print("==========\n");
			var i = 0;
			row.foreach((s) => {
				if ( s == null || s == "" ) {
					s = "NULL";
				}
				print(@"$(fields.nth_data(i).name) = $s\n");
				i++;
			});
			print("==========\n");
		});

		return 0;
	}
}
