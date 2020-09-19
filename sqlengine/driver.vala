namespace SQLEngine {
	errordomain ErrorConnection {
		CODE_1
	}

	errordomain ErrorQuery {
		CODE_1
	}

	enum ColType {
		Boolean,
		String,
		LongString,
		Float,
		Int,
		Date,
		List
	}

	enum TableType {
		Regular,
		View,
		Dummy
	}

	enum Sort {
		Asc,
		Desc
	}

	enum Operator {
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
		IsNotNull   // = "NOT NULL";
	}

	class ConnectionInfo {
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
	}

	class Query {
		public int64 id           { get; set; }
		public string name        { get; set; }
		public string query       { get; set; }
		public int64 connectionId { get; set; }
	}

	class TableDef {
		public string name     { get; set; }
		public TableType ttype { get; set; }
		public string query    { get; set; }
	}

	// ColDef describe a column
	class ColDef {
		public string name     { get; set; }
		public bool pk         { get; set; }
		public bool fk         { get; set; }
		public int precision   { get; set; }
		public bool unsigned   { get; set; }
		public bool nullable   { get; set; }
		public ColType ttype   { get; set; }
		public string[] values { get; set; }
	}

	class CondStmt {
		public ColDef field { get; set; }
		public Operator op  { get; set; }
		public string val   { get; set; }
	}

	class SortOption {
		public ColDef col   { get; set; }
		public Sort sortdir { get; set; }
	}

	class FetchTableOptions  {
		public int64 offset          { set; get; }
		public int64 limit           { set; get; }
		public SortOption[] sort     { set; get; }
		public CondStmt[] conditions { set; get; }
	}

	interface Driver {
		public abstract Connection Connect(ConnectionInfo c) throws ErrorConnection;
		public abstract bool ValidateConnection(ConnectionInfo c);
	}

	interface Connection {
		public abstract List<string> Databases() throws ErrorQuery;
		public abstract void UseDatabase(string name) throws ErrorQuery;
		public abstract bool Disconnect();
		public abstract void Reconnect() throws ErrorConnection;

		public abstract List<TableDef> Tables() throws ErrorQuery;
		public abstract List<ColDef> TableDefinition(string name) throws ErrorQuery;
		public abstract void DeleteTable(TableDef def) throws ErrorQuery;
		public abstract void TruncateTable(TableDef def) throws ErrorQuery;
		public abstract void FetchTable(string name, FetchTableOptions opts, ref List<ColDef> def, ref List<List<string>> rows) throws ErrorQuery;
		//public abstract void DeleteRecord(string name, ColDef[] defs, string[] row) throws ErrorQuery;
		//public abstract string UpdateRecord(string name, ColDef[] defs, string[] newrow, string[] oldrow) throws ErrorQuery; // new, oldvalues;
		//public abstract string UpdateField(string name, ColDef[] defs, string[] row ) throws ErrorQuery;
		//public abstract string UpdateFields(string name, ColDef[] defs, string[] row, int keys) throws ErrorQuery;
		//public abstract string[] InsertRecord(string name, ColDef[] defs, string[] row) throws ErrorQuery;
		// NOTE: everything is an string... so ? public abstract string ParseValue(def ColDef, value string) interface{}
		//public abstract void Query(string query, ref string[] colnames, ref string[,] rows) throws ErrorQuery;
		//public abstract void Execute(string query, ref string lastId, ref int64 count) throws ErrorQuery;
		//public abstract string Name();
		// DDL
		//public abstract string GetCreateTable(string name) throws ErrorQuery;
		//public abstract string GetInsertStatement(string name, ColDef[] def, string[] row) throws ErrorQuery;
		//public abstract string GetSelectStatement(TableDef def) throws ErrorQuery;
	}
}
