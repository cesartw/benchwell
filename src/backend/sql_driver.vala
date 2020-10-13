public interface Benchwell.Backend.Sql.Driver {
	public abstract Benchwell.Backend.Sql.Connection connect(Benchwell.Backend.Sql.ConnectionInfo c) throws Benchwell.Backend.Sql.Error;
}

public interface Benchwell.Backend.Sql.Connection : Object {
	public abstract List<string> databases () throws Benchwell.Backend.Sql.Error;
	public abstract void use_database (string name) throws Benchwell.Backend.Sql.Error;
	public abstract bool disconnect ();
	public abstract void reconnect () throws Benchwell.Backend.Sql.Error;

	public abstract Benchwell.Backend.Sql.TableDef[] tables () throws Benchwell.Backend.Sql.Error;
	public abstract Benchwell.Backend.Sql.ColDef[] table_definition (string name) throws Benchwell.Backend.Sql.Error;
	public abstract void delete_table(TableDef def) throws Benchwell.Backend.Sql.Error;
	public abstract void truncate_table(TableDef def) throws Benchwell.Backend.Sql.Error;
	public abstract List<List<string?>> fetch_table(
		string name,
		Benchwell.Backend.Sql.CondStmt[]? conditions,
		Benchwell.Backend.Sql.SortOption[]? opts,
		int limit,
		int offset
		) throws Benchwell.Backend.Sql.Error;
	public abstract void update_field (string name, ColDef[] defs, string[] row) throws Error;
	public abstract string[] insert_record(string name, ColDef[] defs, string[] row) throws Error;
	public abstract void delete_record(string name, ColDef[] defs, string[] row) throws Error;
	public abstract string get_create_table(string name) throws Error;
	public abstract void query(string query, out string[] columns, out List<List<string?>> rows) throws Error;
	public abstract string get_insert_statement(string name, unowned ColDef[] columns, unowned string[] row);

	//public abstract string get_select_statement(TableDef def) throws Error;
	//public abstract string update_record(string name, ColDef[] defs, string[] newrow, string[] oldrow) throws Error; // new, oldvalues;
	//public abstract string update_fields(string name, ColDef[] defs, string[] row, int keys) throws Error;
	// NOTE: everything is an string... so ? public abstract string ParseValue(def ColDef, value string) interface{}
	//public abstract void execute(string query, ref string lastId, ref int64 count) throws Error;
	//public abstract string name();
	// DDL
}

public errordomain Benchwell.Backend.Sql.Error {
	CONNECTION,
	QUERY
}

public enum Benchwell.Backend.Sql.ColType {
	Boolean,
	String,
	LongString,
	Float,
	Int,
	Date,
	List
}

public enum Benchwell.Backend.Sql.TableType {
	Regular,
	View,
	Dummy
}

public enum Benchwell.Backend.Sql.SortType {
	Asc,
	Desc
}

public enum Benchwell.Backend.Sql.Operator {
	Eq        , // = "=";
	Neq       , // = "!=";
	Gt        , // = ">";
	Lt        , // = "<";
	Gte       , // = ">=";
	Lte       , // = "<=";
	Like      , // = "LIKE";
	NotLike   , // = "NOT LIKE";
	In        , // = "IN";
	Nin       , // = "NOT IN";
	IsNull    , // = "IS NULL";
	IsNotNull;   // = "NOT NULL";

	public string to_string () {
		switch (this) {
			case Eq:
				return "=";
			case Neq:
				return "!=";
			case Gt:
				return ">";
			case Lt:
				return "<";
			case Gte:
				return ">=";
			case Lte:
				return "<=";
			case Like:
				return "LIKE";
			case NotLike:
				return "NOT LIKE";
			case In:
				return "IN";
			case Nin:
				return "NOT IN";
			case IsNull:
				return "IS NULL";
			case IsNotNull:
				return "NOT NULL";
		}

		return "";
	}

	public static Benchwell.Backend.Sql.Operator[] all () {
		return {
			Eq        ,
			Neq       ,
			Gt        ,
			Lt        ,
			Gte       ,
			Lte       ,
			Like      ,
			NotLike   ,
			In        ,
			Nin       ,
			IsNull    ,
			IsNotNull
		};
	}

	public static Benchwell.Backend.Sql.Operator? parse (string s) {
		switch (s) {
			case "=":
				return Eq;
			case "!=":
				return Neq;
			case ">":
				return Gt;
			case "<":
				return Lt;
			case ">=":
				return Gte;
			case "<=":
				return Lte;
			case "LIKE":
				return Like;
			case "NOT LIKE":
				return NotLike;
			case "IN":
				return In;
			case "NOT IN":
				return Nin;
			case "IS NULL":
				return IsNull;
			case "NOT NULL":
				return IsNotNull;
		}

		return null;
	}
}

public class Benchwell.Backend.Sql.ConnectionInfo {
	public int64  id       { get; set; }
	public string adapter  { get; set; }
	public string ttype    { get; set; }
	public string name     { get; set; }
	public string socket   { get; set; }
	public string file     { get; set; }
	public string host     { get; set; }
	public int    port     { get; set; }
	public string user     { get; set; }
	public string password { get; set; }
	public string database { get; set; }
	public string sshhost  { get; set; }
	public string sshagent { get; set; }
	public string options  { get; set; }
	public bool encrypted  { get; set; }
	public Query[] queries { get; set; }

	public ConnectionInfo () {
		adapter     = "";
		ttype       = "";
		name        = "";
		socket      = "";
		file        = "";
		host        = "";
		port        = 0;
		user        = "";
		password    = "";
		database    = "";
		sshhost     = "";
		sshagent    = "";
		options     = "";
		encrypted   = false;
	}
	public string to_string() {
		return name;
	}
}

public class Benchwell.Backend.Sql.Query : Object {
	public int64 id           { get; set; }
	public string name        { get; set; }
	public string query       { get; set; }
	public int64 connection_id { get; set; }
}

public class Benchwell.Backend.Sql.TableDef : Object {
	public string name     { get; set; }
	public Benchwell.Backend.Sql.TableType ttype { get; set; }
	public string query    { get; set; }

	public TableDef.with_name(string n) {
		Object();
		name = n;
	}

	public string to_string() {
		return name;
	}
}

// ColDef describe a column
public class Benchwell.Backend.Sql.ColDef : Object {
	public string name     { get; set; }
	public bool pk         { get; set; }
	public bool fk         { get; set; }
	public int precision   { get; set; }
	public bool unsigned   { get; set; }
	public bool nullable   { get; set; }
	public Benchwell.Backend.Sql.ColType ttype   { get; set; }
	public string[] values { get; set; }

	public ColDef.with_name (string n) {
		name = n;
	}
}

public class Benchwell.Backend.Sql.CondStmt : Object {
	public Benchwell.Backend.Sql.ColDef field { get; set; }
	public Benchwell.Backend.Sql.Operator op  { get; set; }
	public string val   { get; set; }
}

public class Benchwell.Backend.Sql.SortOption: Object {
	public Benchwell.Backend.Sql.ColDef column { get; construct; }
	public Benchwell.Backend.Sql.SortType dir      { get; construct; }

	public SortOption(Benchwell.Backend.Sql.ColDef column, Benchwell.Backend.Sql.SortType dir){
		Object(
			column: column,
			dir: dir
		);
	}
}
