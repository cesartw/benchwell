public interface Benchwell.SQL.Driver {
	public abstract Benchwell.SQL.Connection connect(Benchwell.SQL.ConnectionInfo c) throws Benchwell.SQL.ErrorConnection;
}

public interface Benchwell.SQL.Connection : Object {
	public abstract List<string> databases () throws Benchwell.SQL.ErrorQuery;
	public abstract void use_database (string name) throws Benchwell.SQL.ErrorQuery;
	public abstract bool disconnect ();
	public abstract void reconnect () throws Benchwell.SQL.ErrorConnection;

	public abstract Benchwell.SQL.TableDef[] tables () throws Benchwell.SQL.ErrorQuery;
	public abstract Benchwell.SQL.ColDef[] table_definition (string name) throws Benchwell.SQL.ErrorQuery;
	public abstract void delete_table(TableDef def) throws Benchwell.SQL.ErrorQuery;
	public abstract void truncate_table(TableDef def) throws Benchwell.SQL.ErrorQuery;
	public abstract List<List<string?>> fetch_table(
		string name,
		Benchwell.SQL.CondStmt[]? conditions,
		Benchwell.SQL.SortOption[]? opts,
		int limit,
		int offset
		) throws Benchwell.SQL.ErrorQuery;
	public abstract void update_field (string name, ColDef[] defs, string[] row) throws ErrorQuery;
	public abstract string[] insert_record(string name, ColDef[] defs, string[] row) throws ErrorQuery;
	//public abstract void delete_record(string name, ColDef[] defs, string[] row) throws ErrorQuery;
	//public abstract string update_record(string name, ColDef[] defs, string[] newrow, string[] oldrow) throws ErrorQuery; // new, oldvalues;
	//public abstract string update_fields(string name, ColDef[] defs, string[] row, int keys) throws ErrorQuery;
	// NOTE: everything is an string... so ? public abstract string ParseValue(def ColDef, value string) interface{}
	//public abstract void query(string query, ref string[] colnames, ref string[,] rows) throws ErrorQuery;
	//public abstract void execute(string query, ref string lastId, ref int64 count) throws ErrorQuery;
	//public abstract string name();
	// DDL
	//public abstract string get_create_table(string name) throws ErrorQuery;
	//public abstract string get_insert_statement(string name, ColDef[] def, string[] row) throws ErrorQuery;
	//public abstract string get_select_statement(TableDef def) throws ErrorQuery;
}

public errordomain Benchwell.SQL.ErrorConnection {
	CODE_1
}

public errordomain Benchwell.SQL.ErrorQuery {
	CODE_1
}

public enum Benchwell.SQL.ColType {
	Boolean,
	String,
	LongString,
	Float,
	Int,
	Date,
	List
}

public enum Benchwell.SQL.TableType {
	Regular,
	View,
	Dummy
}

public enum Benchwell.SQL.SortType {
	Asc,
	Desc
}

public enum Benchwell.SQL.Operator {
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

	public static Benchwell.SQL.Operator[] all () {
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
}

public class Benchwell.SQL.ConnectionInfo {
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
		adapter  = "";
		ttype    = "";
		name     = "";
		socket   = "";
		file     = "";
		host     = "";
		port     = 0;
		user     = "";
		password = "";
		database = "";
		sshhost  = "";
		sshagent = "";
		options  = "";
		encrypted  = false;
	}
	public string to_string() {
		return name;
	}
}

public class Benchwell.SQL.Query : Object {
	public int64 id           { get; set; }
	public string name        { get; set; }
	public string query       { get; set; }
	public int64 connectionId { get; set; }
}

public class Benchwell.SQL.TableDef : Object {
	public string name     { get; set; }
	public Benchwell.SQL.TableType ttype { get; set; }
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
public class Benchwell.SQL.ColDef : Object {
	public string name     { get; set; }
	public bool pk         { get; set; }
	public bool fk         { get; set; }
	public int precision   { get; set; }
	public bool unsigned   { get; set; }
	public bool nullable   { get; set; }
	public Benchwell.SQL.ColType ttype   { get; set; }
	public string[] values { get; set; }
}

public class Benchwell.SQL.CondStmt : Object {
	public Benchwell.SQL.ColDef field { get; set; }
	public Benchwell.SQL.Operator op  { get; set; }
	public string val   { get; set; }
}

public class Benchwell.SQL.SortOption: Object {
	public Benchwell.SQL.ColDef column { get; construct; }
	public Benchwell.SQL.SortType dir      { get; construct; }

	public SortOption(Benchwell.SQL.ColDef column, Benchwell.SQL.SortType dir){
		Object(
			column: column,
			dir: dir
		);
	}
}