public class Benchwell.SQLiteForm : Gtk.Grid {
	public Benchwell.Backend.Sql.ConnectionInfo? connection;
	public string filename;

	public Gtk.Label tab_label;

	public Gtk.Entry name_entry;
	public Gtk.FileChooserButton file_btn;
	public Gtk.FileChooserDialog file_dialog;

	public Gtk.Label name_lbl;
	public Gtk.Label file_lbl;

	public SQLiteForm(){
		Object();
		set_name ("form");
		set_column_homogeneous (true);
		set_row_spacing (5);

		tab_label = new Gtk.Label ("SQLite");
		tab_label.show ();

		name_entry = new Gtk.Entry();
		name_entry.show ();

		name_lbl = new Gtk.Label("Name");
		name_lbl.set_halign (Gtk.Align.START);
		name_lbl.show ();

		file_lbl = new Gtk.Label("File");
		file_lbl.set_halign (Gtk.Align.START);
		file_lbl.show ();

		file_dialog = new Gtk.FileChooserDialog ("Select", null, Gtk.FileChooserAction.OPEN,
											"Ok", Gtk.ResponseType.OK,
											"Cancel", Gtk.ResponseType.CANCEL);
		file_btn = new Gtk.FileChooserButton.with_dialog (file_dialog);
		file_btn.show ();

		attach(name_lbl, 0, 1, 1, 1);
		attach(name_entry, 1, 1, 2, 1);

		attach(file_lbl, 0, 2, 1, 1);
		attach(file_btn, 1, 2, 2, 1);
	}

	public void clear() {
		connection = null;
		name_entry.set_text ("");
		filename = "";
	}

	public void set_connection (Benchwell.Backend.Sql.ConnectionInfo conn) {
		connection = conn;

		name_entry.set_text (connection.name);
		file_btn.select_filename (conn.file);
		filename = conn.file;
	}

	public Benchwell.Backend.Sql.ConnectionInfo get_connection () {
		var conn = new Benchwell.Backend.Sql.ConnectionInfo ();
		conn.adapter = "mysql";
		conn.ttype = "socket";
		conn.name = name_entry.get_text ();
		conn.file = filename;

		if ( connection != null ) {
			conn.id = connection.id;
		}

		return conn;
	}
}
