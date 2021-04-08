public class Benchwell.SQLiteForm : Gtk.Grid {
	public weak Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.ConnectionInfo? connection;
	public string filename;

	public Gtk.Label tab_label;

	public Gtk.Entry name_entry;
	public Gtk.FileChooserButton file_btn;
	public Gtk.FileChooserDialog file_dialog;

	public Gtk.Label name_lbl;
	public Gtk.Label file_lbl;

	public SQLiteForm (Benchwell.ApplicationWindow w, Gtk.SizeGroup size_group){
		Object(
			window: w,
			name: "form",
			column_homogeneous: true,
			row_spacing: 5
		);

		tab_label = new Gtk.Label ("SQLite");
		tab_label.show ();

		name_entry = new Gtk.Entry();
		name_entry.show ();

		name_lbl = new Gtk.Label("Name") {
			xalign = 0
		};
		name_lbl.show ();

		file_lbl = new Gtk.Label("File") {
			xalign = 0
		};
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

		size_group.add_widget (name_lbl);
		size_group.add_widget (name_entry);

		size_group.add_widget (file_lbl);
		size_group.add_widget (file_btn);
	}

	public void clear() {
		connection = null;
		name_entry.set_text ("");
		filename = "";
	}

	public void set_connection (Benchwell.ConnectionInfo conn) {
		connection = conn;

		name_entry.set_text (connection.name);
		file_btn.select_filename (conn.file);
		filename = conn.file;
	}

	public Benchwell.ConnectionInfo get_connection () {
		var conn = new Benchwell.ConnectionInfo ();
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
