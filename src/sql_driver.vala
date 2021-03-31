public interface Benchwell.Driver {
	public abstract Benchwell.Connection connect(Benchwell.ConnectionInfo c) throws Benchwell.Error;
}

public interface Benchwell.Connection : Object {
	public abstract List<string> databases () throws Benchwell.Error;
	public abstract void use_database (string name) throws Benchwell.Error;
	public abstract bool disconnect ();
	public abstract void reconnect () throws Benchwell.Error;

	public abstract Benchwell.TableDef[] tables () throws Benchwell.Error;
	public abstract Benchwell.ColDef[] table_definition (string name) throws Benchwell.Error;
	public abstract void delete_table(TableDef def) throws Benchwell.Error;
	public abstract void truncate_table(TableDef def) throws Benchwell.Error;
	public abstract List<List<string?>> fetch_table(
		string name,
		Benchwell.CondStmt[]? conditions,
		Benchwell.SortOption[]? opts,
		int limit,
		int offset
		) throws Benchwell.Error;
	public abstract void update_field (string name, ColDef[] defs, string[] row) throws Error;
	public abstract string[]? insert_record(string name, ColDef[] defs, string[] row) throws Error;
	public abstract void delete_record(string name, ColDef[] defs, string[] row) throws Error;
	public abstract string get_create_table(string name) throws Error;
	public abstract void query(string query, out string[] columns, out List<List<string?>> rows) throws Error;
	public abstract string get_insert_statement(string name, ColDef[] columns, string[] row);

	//public abstract string get_select_statement(TableDef def) throws Error;
	//public abstract string update_record(string name, ColDef[] defs, string[] newrow, string[] oldrow) throws Error; // new, oldvalues;
	//public abstract string update_fields(string name, ColDef[] defs, string[] row, int keys) throws Error;
	// NOTE: everything is an string... so ? public abstract string ParseValue(def ColDef, value string) interface{}
	//public abstract void execute(string query, ref string lastId, ref int64 count) throws Error;
	//public abstract string name();
	// DDL
}

public errordomain Benchwell.Error {
	CONNECTION,
	QUERY
}

public enum Benchwell.ColType {
	Boolean,
	String,
	LongString,
	Float,
	Int,
	Date,
	List
}

public enum Benchwell.TableType {
	Regular,
	View,
	Dummy
}

public enum Benchwell.SortType {
	Asc,
	Desc
}

public enum Benchwell.Operator {
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

	public static Benchwell.Operator[] all () {
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

	public static Benchwell.Operator? parse (string s) {
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

public class Benchwell.TableDef : Object {
	public string name     { get; set; }
	public Benchwell.TableType ttype { get; set; }
	public string query    { get; set; }
	public Object source { get; set; }

	public TableDef.with_name(string n) {
		Object();
		name = n;
	}

	public string to_string() {
		return name;
	}
}

// ColDef describe a column
public class Benchwell.ColDef : Object {
	public string name     { get; set; }
	public bool pk         { get; set; }
	public bool fk         { get; set; }
	public int precision   { get; set; }
	public bool unsigned   { get; set; }
	public bool nullable   { get; set; }
	public Benchwell.ColType ttype   { get; set; }
	public string[] values { get; set; }

	public ColDef.with_name (string n) {
		name = n;
	}
}

public class Benchwell.CondStmt : Object {
	public Benchwell.ColDef field { get; set; }
	public Benchwell.Operator op  { get; set; }
	public string val   { get; set; }
}

public class Benchwell.SortOption: Object {
	public Benchwell.ColDef column { get; construct; }
	public Benchwell.SortType dir      { get; construct; }

	public SortOption(Benchwell.ColDef column, Benchwell.SortType dir){
		Object(
			column: column,
			dir: dir
		);
	}
}
