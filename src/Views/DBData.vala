public class Benchwell.Views.DBData : Gtk.Paned {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.SQL.Connection connection { get; construct; }
	public Benchwell.SQL.ConnectionInfo connection_info { get; construct; }

	public Gtk.SearchEntry table_search;
	public Gtk.ComboBoxText database_combo;
	public Benchwell.Views.DBTables tables;
	public Benchwell.Views.DBResultView result_view;

	private Benchwell.Views.CancelOverlay overlay;
	private List<string> databases;
	private Benchwell.SQL.TableDef? table_def;

	private int current_page = 0;
	private int page_size = 100;

	public signal void database_selected(string dbname);

	public DBData (Benchwell.ApplicationWindow window,
				   Benchwell.SQL.Connection connection,
				   Benchwell.SQL.ConnectionInfo connection_info) {
		Object(
			window: window,
			connection: connection,
			connection_info: connection_info,
			orientation: Gtk.Orientation.VERTICAL,
			wide_handle: true,
			vexpand: true,
			hexpand: true
		);

		build ();

		fill ();

		table_search.search_changed.connect ( () => {
			var expr = table_search.get_buffer ().get_text ();
			Regex regex;
			try {
				regex = new Regex (expr, RegexCompileFlags.CASELESS);
			} catch (RegexError e) {
				regex = new Regex (".*");
			}

			tables.filter = regex;
			tables.invalidate_filter ();
		});

		tables.row_activated.connect (on_load_table);

		result_view.table.field_change.connect (on_field_change);
		result_view.table.btn_refresh.clicked.connect (on_refresh_table);
		result_view.table.conditions.search.connect (on_refresh_table);

		tables.schema_menu.activate.connect (on_show_schema);
		tables.refresh_menu.activate.connect (on_refresh_tables);

		tables.truncate_menu.activate.connect (on_truncate_table);
		tables.delete_menu.activate.connect (on_delete_table);

		result_view.table.btn_prev.clicked.connect (on_prev_page);
		result_view.table.btn_next.clicked.connect (on_next_page);

		result_view.table.copy_insert_menu.activate.connect (on_copy_insert);
	}

	private void on_copy_insert () {
		var data = result_view.table.get_selected_data ();
		if (data == null) {
			return;
		}

		var st = connection.get_insert_statement (table_def.name, result_view.table.columns, data);
		var cb = Gtk.Clipboard.get_default (Gdk.Display.get_default ());
		cb.set_text (st, st.length);
	}

	private void on_show_schema () {
		var tabledef = tables.selected_tabledef;
		if (tabledef == null) {
			return;
		}

		var sql = connection.get_create_table (tabledef.name);
		var dialog = new Gtk.Dialog.with_buttons (@"$(tabledef.name) schema", window,
								Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
								"Ok", Gtk.ResponseType.OK);
		dialog.set_default_size (400, 400);

		var sv = new Benchwell.SourceView ();
		sv.show ();
		sv.get_buffer ().set_text (sql);

		var sw = new Gtk.ScrolledWindow (null, null);
		sw.add (sv);
		sw.show ();

		dialog.get_content_area ().add (sw);

		dialog.run ();
		dialog.destroy ();
	}

	private void on_refresh_tables () {
		on_database_selected ();
	}

	private void on_delete_table () {
		var tabledef = tables.selected_tabledef;
		if (tabledef == null) {
			return;
		}

		if (tabledef.ttype == Benchwell.SQL.TableType.Dummy) {
			//Config.delete_query (tabledef);
		} else {
			connection.delete_table (tabledef);
		}
		tables.remove_selected ();
	}

	private void on_truncate_table () {
		var tabledef = tables.selected_tabledef;
		if (tabledef == null) {
			return;
		}

		try {
			connection.truncate_table (tabledef);
		} catch (Benchwell.SQL.ErrorQuery e) {
			stderr.printf (@"$(e.message)");
		}
	}

	private void build () {
		var hpaned = new Gtk.Paned (Gtk.Orientation.HORIZONTAL);
		hpaned.set_wide_handle (true);
		hpaned.set_hexpand (true);
		hpaned.set_vexpand (true);
		hpaned.show ();

		var sidebar = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		sidebar.show ();

		table_search = new Gtk.SearchEntry ();
		table_search.set_placeholder_text (_("Filter table: .*"));
		table_search.show ();

		tables = new Benchwell.Views.DBTables ();
		tables.activate_on_single_click = false;
		tables.show ();

		var tables_sw = new Gtk.ScrolledWindow (null, null);
		tables_sw.show ();
		tables_sw.add (tables);

		database_combo = new Gtk.ComboBoxText ();
		database_combo.set_id_column (0);
		database_combo.show ();

		sidebar.pack_start (database_combo, false, true, 0);
		sidebar.pack_start (table_search, false, true, 0);
		sidebar.pack_start (tables_sw, true, true, 0);

		// main section
		var main_section = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		main_section.show ();

		result_view = new Benchwell.Views.DBResultView (window);
		result_view.show ();

		main_section.add (result_view);
		main_section.set_vexpand (true);
		main_section.set_hexpand (true);
		///////////////

		overlay = new Benchwell.Views.CancelOverlay (main_section);
		overlay.show ();

		hpaned.pack1 (sidebar, false, true);
		hpaned.pack2 (overlay, true, false);

		pack1 (hpaned, false, false);

		// signals

		database_combo.changed.connect (on_database_selected);

		result_view.table.btn_save_row.clicked.connect (() => {
			var data = result_view.table.get_selected_data ();
			if (data == null) {
				return;
			}

			data = connection.insert_record (table_def.name, result_view.table.columns, data);
			result_view.table.update_selected_row (data);
		});

		result_view.table.btn_delete_row.clicked.connect (() => {
			var data = result_view.table.get_selected_data ();
			if (data == null) {
				return;
			}
			connection.delete_record (table_def.name, result_view.table.columns, data);
			result_view.table.delete_selected_row ();
		});

		result_view.exec_query.connect ((query) => {
			string[] columns;
			List<List<string?>>data;
			connection.query(query, out columns, out data);

			Benchwell.SQL.ColDef[] cols = {};
			foreach (var column in columns) {
				cols += new Benchwell.SQL.ColDef.with_name (column);
			}

			result_view.table.columns = cols;
			result_view.table.data = data;
			result_view.table.raw_mode = true;
		});
	}

	private void fill () {
		databases = connection.databases ();
		databases.foreach ( db => {
			database_combo.append (db, db);
		});

		if ( connection_info.database != "" ) {
			database_combo.set_active_id (connection_info.database);
		}
	}

	private void on_database_selected () {
		var dbname = database_combo.get_active_text ();
		try {
			connection.use_database (dbname);
			var tt = connection.tables ();

			foreach (var q in connection_info.queries) {
				var t = new Benchwell.SQL.TableDef.with_name (q.name);
				t.ttype = Benchwell.SQL.TableType.Dummy;
				tt += t;
			}


			tables.update_items (tt);

			database_selected (dbname);
		} catch (Benchwell.SQL.ErrorQuery e) {
			stderr.printf(@"error: $(e.message)");
		}
	}

	private void on_field_change(Benchwell.SQL.ColDef[] columns, string[] row) {
		connection.update_field (table_def.name, columns, row);
	}

	private void on_load_table () {
		current_page = 0;
		table_def = tables.get_selected_table ();
		try {
			result_view.table.columns = connection.table_definition (table_def.name);
			result_view.table.data = connection.fetch_table (table_def.name,
															 result_view.table.get_conditions (),
															 result_view.table.get_sort_options (),
															 page_size, current_page*page_size);
			result_view.table.raw_mode = false;
		} catch (Benchwell.SQL.ErrorQuery e) {
			stderr.printf(@"error: $(e.message)");
		}
	}

	private void on_refresh_table () {
		if (result_view.table.raw_mode) {
			result_view._exec_query ();
		} else {
			if (table_def == null ) {
				return;
			}

			result_view.table.data = connection.fetch_table (table_def.name,
															 result_view.table.get_conditions (),
															 result_view.table.get_sort_options (),
															 page_size, current_page*page_size);
		}

	}

	private void on_prev_page () {
		if (current_page == 0) {
			return;
		}
		current_page--;
		on_refresh_table ();

	}

	private void on_next_page () {
		current_page++;
		on_refresh_table ();
	}
}

public class Benchwell.Views.DBResultView : Gtk.Paned {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Gtk.SourceView editor;
	public Benchwell.Views.DBTable table;
	public Gtk.Button btn_load_query;
	public Gtk.MenuButton save_menu;

	public signal void exec_query (string query);

	public List<List<string?>> data {
		set { table.data = value; }
	}

	public DBResultView (Benchwell.ApplicationWindow window) {
		Object (
			window: window,
			orientation: Gtk.Orientation.VERTICAL
		);

		table = new Benchwell.Views.DBTable (window);
		table.show ();

		// editor
		editor = new Benchwell.SourceView ();
		editor.show ();

		var editor_sw = new Gtk.ScrolledWindow (null, null);
		editor_sw.show ();
		editor_sw.add (editor);

		// table controls
		btn_load_query = new Benchwell.Button ("open", Gtk.IconSize.BUTTON);
		btn_load_query.show ();

		var img = new Benchwell.Image("save");
		save_menu = new Gtk.MenuButton ();
		save_menu.show ();
		save_menu.set_image (img);

		var save_menu_model = new GLib.Menu ();
		save_menu_model.append (_("Save As"), "win.save.file");
		save_menu_model.append (_("Save fav"), "win.save.fav");

		save_menu.set_menu_model (save_menu_model);
		var action_save_file = new GLib.SimpleAction ("save.file", null);
		var action_save_fav = new GLib.SimpleAction ("save.fav", null);
		window.add_action (action_save_file);
		window.add_action (action_save_fav);

		var editor_actionbar = new Gtk.ActionBar ();
		editor_actionbar.show ();
		editor_actionbar.pack_end (save_menu);
		editor_actionbar.pack_end (btn_load_query);
		editor_actionbar.set_name ("queryactionbar");
		/////////////////

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		box.pack_start (editor_sw, true, true, 0);
		box.pack_end (editor_actionbar, false, false, 0);
		box.show ();

		pack1 (box, false, false);
		pack2 (table, true, false);

		action_save_file.activate.connect (on_save_file);
		action_save_fav.activate.connect (on_save_fav);
		btn_load_query.clicked.connect (on_open_file);

		editor.key_press_event.connect (on_editor_key_press);
	}

	public void on_open_file () {
		var dialog = new Gtk.FileChooserDialog (_("Select file"), window,
											 Gtk.FileChooserAction.OPEN,
											_("Ok"), Gtk.ResponseType.OK,
											_("Cancel"), Gtk.ResponseType.CANCEL);
		var resp = (Gtk.ResponseType) dialog.run ();

		if (resp == Gtk.ResponseType.CANCEL) {
			dialog.destroy ();
			return;
		}

		var filename = dialog.get_filename ();
		dialog.destroy ();

		string text;
		var ok = GLib.FileUtils.get_contents (filename, out text, null);
		if (!ok) {
			return;
		}

		editor.get_buffer ().set_text (text);
	}

	public void on_save_file () {
		var dialog = new Gtk.FileChooserDialog (_("Save file"), window,
											 Gtk.FileChooserAction.SAVE,
											_("Ok"), Gtk.ResponseType.OK,
											_("Cancel"), Gtk.ResponseType.CANCEL);
		var resp = (Gtk.ResponseType) dialog.run ();

		if (resp == Gtk.ResponseType.CANCEL) {
			dialog.destroy ();
			return;
		}

		var filename = dialog.get_filename ();
		dialog.destroy ();

		var buffer = editor.get_buffer ();
		Gtk.TextIter start, end;
		buffer.get_start_iter (out start);
		buffer.get_end_iter (out end);
		var txt = buffer.get_text (start, end, false);

		var fs = FileStream.open (filename, "w");
		fs.puts (txt);
	}

	public void on_save_fav () {
		var filename = ask_fav_name ();

		if (filename == "") {
			return;
		}

		var buffer = editor.get_buffer ();
		Gtk.TextIter start, end;
		buffer.get_start_iter (out start);
		buffer.get_end_iter (out end);
		var txt = buffer.get_text (start, end, false);
	}

	private string ask_fav_name () {
		var dialog = new Gtk.Dialog.with_buttons (_("Select"), window,
									Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
									_("Ok"), Gtk.ResponseType.OK,
									_("Cancel"), Gtk.ResponseType.CANCEL);
		dialog.set_default_size (250, 130);

		var label = new Gtk.Label (_("Enter favorite name"));
		label.show ();

		var entry = new Gtk.Entry ();
		entry.show ();

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 10);
		box.show ();

		box.pack_start (label, true, true, 0);
		box.pack_start (entry, true, true, 0);

		dialog.get_content_area ().add (box);

		var resp = (Gtk.ResponseType) dialog.run ();
		var filename = entry.get_text ();
		dialog.destroy ();

		if (resp != Gtk.ResponseType.OK) {
			return "";
		}

		return filename;
	}

	private bool on_editor_key_press (Gtk.Widget widget, Gdk.EventKey event) {
		if (event.keyval != Gdk.Key.Return) {
			return false;
		}

		if (event.state != Gdk.ModifierType.CONTROL_MASK) {
			return false;
		}

		_exec_query ();
		return true;
	}

	public void _exec_query () {
		var query = editor.get_buffer ().text;
		exec_query (query);
	}
}

public class Benchwell.Views.DBTable : Gtk.Box {
	public Benchwell.ApplicationWindow window { get; construct; }
	private bool _raw_mode = false;
	public bool raw_mode {
		get { return _raw_mode; }
		set {
			_raw_mode = value;
			btn_prev.sensitive = !_raw_mode;
			btn_next.sensitive = !_raw_mode;
			btn_show_filters.sensitive = !_raw_mode;
			btn_add_row.sensitive = !_raw_mode;
			btn_delete_row.sensitive = !_raw_mode;
			btn_save_row.sensitive = !_raw_mode;
			btn_show_filters.active = !_raw_mode;
		}
	}
	public Gtk.TreeView table;
	public Gtk.Menu menu;
	public Gtk.MenuItem clone_menu;
	public Gtk.MenuItem copy_insert_menu;
	public Gtk.MenuItem copy_menu;
	public Gtk.ListStore store;
	public Gtk.Button btn_prev;
	public Gtk.Button btn_next;
	public Gtk.Button btn_refresh;
	public Gtk.ToggleButton btn_show_filters;
	public Gtk.Button btn_add_row;
	public Gtk.Button btn_delete_row;
	public Gtk.Button btn_save_row;
	public Gtk.SearchEntry search;
	public Benchwell.Views.DBConditions conditions;
	public Benchwell.SQL.ColDef[] columns {
		get { return _columns; }
		set { _columns = value; _update_columns (); }
	}
	public List<List<string?>> data {
		set {
			_update_data (value);
		}
	}
	private Benchwell.SQL.ColDef[] _columns;

	public signal void field_change (Benchwell.SQL.ColDef[] column, string[] row);

	public DBTable (Benchwell.ApplicationWindow window) {
		Object (
			window: window,
			orientation: Gtk.Orientation.VERTICAL
		);

		table = new Gtk.TreeView ();
		table.rubber_banding = true;
		table.enable_grid_lines = Gtk.TreeViewGridLines.HORIZONTAL;
		table.activate_on_single_click = true;
		table.enable_search = true;
		table.show ();

		var table_sw = new Gtk.ScrolledWindow (null, null);
		table_sw.add (table);
		table_sw.show ();

		menu = new Gtk.Menu ();
		clone_menu = new Benchwell.MenuItem (_("Clone row"), "copy");
		clone_menu.show ();

		copy_insert_menu = new Benchwell.MenuItem (_("Copy insert"), "copy");
		copy_insert_menu.show ();

		copy_menu = new Benchwell.MenuItem (_("Copy value"), "copy");
		copy_menu.show ();

		menu.add (clone_menu);
		menu.add (copy_insert_menu);
		menu.add (copy_menu);

		// actionbar
		btn_add_row = new Benchwell.Button ("add-record", Gtk.IconSize.BUTTON);
		btn_add_row.show ();

		btn_delete_row = new Benchwell.Button ("delete-record", Gtk.IconSize.BUTTON);
		btn_delete_row.show ();

		btn_save_row = new Benchwell.Button ("save-record", Gtk.IconSize.BUTTON);
		btn_save_row.sensitive = false;
		btn_save_row.show ();

		btn_show_filters = new Benchwell.ToggleButton ("filter", Gtk.IconSize.BUTTON);
		btn_show_filters.show ();

		btn_prev = new Benchwell.Button ("back", Gtk.IconSize.BUTTON);
		btn_prev.show ();

		btn_next = new Benchwell.Button ("next", Gtk.IconSize.BUTTON);
		btn_next.show ();

		btn_refresh = new Benchwell.Button ("refresh", Gtk.IconSize.BUTTON);
		btn_refresh.show ();

		search = new Gtk.SearchEntry ();
		search.set_placeholder_text (_("Column filter: .*"));
		search.show ();

		var table_actionbar = new Gtk.ActionBar ();
		table_actionbar.show ();
		table_actionbar.add (btn_add_row);
		table_actionbar.add (btn_save_row);
		table_actionbar.add (btn_delete_row);
		table_actionbar.add (btn_show_filters);

		table_actionbar.pack_end (search);
		table_actionbar.pack_end (btn_refresh);
		table_actionbar.pack_end (btn_next);
		table_actionbar.pack_end (btn_prev);
		////////////

		conditions = new Benchwell.Views.DBConditions ();

		pack_start (table_actionbar, false, false, 0);
		pack_start (conditions, false,false, 5);
		pack_start (table_sw, true, true, 0);

		// signals
		table.button_press_event.connect (on_button_press);

		btn_show_filters.toggled.connect (() => { if (btn_show_filters.active) { conditions.show ();
			} else {
				conditions.hide ();
			}
		});

		//column
		search.search_changed.connect (() => {
			var expr = search.get_buffer ().get_text ();
			Regex regex;
			try {
				regex = new Regex (expr, RegexCompileFlags.CASELESS);
			} catch (RegexError e) {
				regex = new Regex (".*");
			}


			table.get_columns ().foreach ((column) => {
				column.visible = regex.match (column.get_title ().replace ("__", "_"));
			});
		});

		btn_add_row.clicked.connect (() => {
			add_empty_row ();
		});

		btn_delete_row.clicked.connect (on_delete_row);
	}

	private void on_delete_row () {
		var selection = table.get_selection ();
		Gtk.TreeIter? iter = null;
		selection.selected_foreach ((model, path, i) => {
			if (iter != null) {
				return;
			}
			iter = i;
		});

		GLib.Value val;
		store.get_value (iter, (int) columns.length, out val);
		if (val.get_int () == 1) {
			delete_selected_row ();
			table.unselect_all ();
		}
	}

	private void on_clone () {
		var data = get_selected_data ();
		var selection = table.get_selection ();

		Gtk.TreeIter? iter = null;
		selection.selected_foreach ((model, path, i) => {
			if (iter != null) {
				return;
			}
			iter = i;
		});

		store.insert_after (out iter, iter);
		for (var i = 0; i < _columns.length; i++) {
			if (_columns[i].pk) {
				continue;
			}
			store.set_value (iter, i, data[i]);
		}
		store.set_value (iter, _columns.length, 1);
		selection.select_iter (iter);
	}

	private void _update_columns () {
		if (store != null) {
			store.clear ();
		}

		while ( table.get_column(0) != null ) {
			table.remove_column (table.get_column(0));
		}

		GLib.Type[] column_types = new GLib.Type[columns.length + 1];
		var i = 0;
		foreach (var column in _columns) {
			table.insert_column (build_column (column, i ), i);
			column_types[i] = GLib.Type.STRING;
			i++;
		};

		column_types[_columns.length] = GLib.Type.INT;

		store = new Gtk.ListStore.newv (column_types);
		table.model = store;
		conditions.columns = _columns;
	}

	private void _update_data (List<List<string?>> data) {
		store.clear ();

		if (data == null) {
			return;
		}
		data.foreach ( (row) => {
			add_row (row);
		});
	}

	public void add_row (List<string?> data) {
		Gtk.TreeIter iter;
		store.append(out iter);

		var i = 0;
		data.foreach ( (val) => {
			if (val == null) {
				val = Benchwell.null_string;
			}

			store.set (iter, i, val);
			i++;
		});
		store.set (iter, i, 0);
	}

	public Benchwell.SQL.SortOption[] get_sort_options () {
		Benchwell.SQL.SortOption[] sorts = {};

		for (var i = 0; i < table.get_n_columns (); i++) {
			var col = table.get_column (i);

			if (!col.get_sort_indicator ()) {
				continue;
			}

			switch (col.get_sort_order ()) {
			case Gtk.SortType.DESCENDING:
				sorts += new Benchwell.SQL.SortOption(columns[i], Benchwell.SQL.SortType.Asc);
				break;
			case Gtk.SortType.ASCENDING:
				sorts += new Benchwell.SQL.SortOption(columns[i], Benchwell.SQL.SortType.Desc);
				break;
			}
		}

		return sorts;
	}

	private Gtk.TreeViewColumn build_column(SQL.ColDef column, int column_index) {
		var renderer = new Gtk.CellRendererText ();
		renderer.editable = true;
		renderer.xpad = 10;
		renderer.height = 23;
		renderer.cell_background = Benchwell.Colors.PkHL.to_string ();
		renderer.cell_background_set = column.pk;
		renderer.ellipsize = Pango.EllipsizeMode.END;

		var _column = new Gtk.TreeViewColumn.with_attributes (column.name.replace ("_", "__"), renderer, "text", column_index);
		_column.resizable = true;
		_column.clickable = true;
		_column.min_width = 85;
		_column.max_width = 250;

		_column.clicked.connect (() => {
			if ( !_column.sort_indicator ) {
				_column.sort_indicator = true;
				_column.sort_order = Gtk.SortType.ASCENDING;
				return;
			}

			if (_column.sort_order == Gtk.SortType.ASCENDING) {
				_column.sort_order = Gtk.SortType.DESCENDING;
			} else {
				_column.sort_indicator = false;
			}
		});

		// NOTE: shiiiitt. affects all cell in the column not just the one
		//_column.set_cell_data_func (renderer, (cell_layout, cell, tree_model, iter) => {
			//GLib.Value val;
			//var index = column_index;
			//tree_model.get_value (iter, index, out val);
			//var path = tree_model.get_path (iter);

			//if ( val.holds (GLib.Type.STRING) ) {
				//if ( val.get_string () == Benchwell.null_string ){
					//cell.cell_background = Benchwell.Colors.NullHL.to_string ();
					//cell.cell_background_set = true;
				//}
			//}
		//});

		renderer.edited.connect ((cell, path, new_value) => {
			on_edited (cell, path, new_value, column_index);
		});

		return _column;
	}

	public void add_empty_row () {
		Gtk.TreeIter? selected = null;
		var selection = table.get_selection ();
		selection.selected_foreach ( (model, path, iter) => {
			if (selected != null){
				return;
			}
			store.get_iter (out selected, path);
		});

		Gtk.TreeIter? insertAt = null;
		if (selected != null) {
			store.insert_after (out insertAt, selected);
		}

		if (insertAt == null) {
			store.append (out insertAt);
		}

		int i;
		for (i = 0; i < columns.length; i++) {
			store.set_value (insertAt, i, Benchwell.null_string);
		}
		store.set_value (insertAt, i, 1);

		var path = store.get_path (insertAt);
		selection.unselect_all ();
		selection.select_path (path);
		table.row_activated (path, null);
		table.scroll_to_cell (path, null, true, (float) 0.5, (float) 0);
	}

	public string[]? get_selected_data () {
		Gtk.TreeIter? selected = null;
		var selection = table.get_selection ();
		selection.selected_foreach ( (model, path, iter) => {
			if (selected != null){
				return;
			}
			store.get_iter (out selected, path);
		});

		if (selected == null) {
			return null;
		}

		var values = new string[columns.length];
		for (var i = 0; i < columns.length; i++) {
			GLib.Value val;
			store.get_value (selected, i, out val);

			values[i] = val.get_string ();
		}

		return values;
	}

	public void delete_selected_row () {
		Gtk.TreeIter? selected = null;
		var selection = table.get_selection ();
		selection.selected_foreach ( (model, path, iter) => {
			if (selected != null){
				return;
			}
			store.get_iter (out selected, path);
			store.remove (ref selected);
		});
	}

	public void update_selected_row (string[] data) {
		Gtk.TreeIter? selected = null;
		var selection = table.get_selection ();
		selection.selected_foreach ( (model, path, iter) => {
			if (selected != null){
				return;
			}
			store.get_iter (out selected, path);
		});

		if (selected == null) {
			return;
		}

		for(var i = 0; i < data.length; i++) {
			store.set_value (selected, i, data[i]);
		}
	}

	public Benchwell.SQL.CondStmt[] get_conditions () {
		return conditions.get_conditions ();
	}

	private void on_edited (Gtk.CellRendererText cell, string path, string new_value, int column_index) {
		// update cell
		Gtk.TreeIter iter ;
		store.get_iter (out iter, new Gtk.TreePath.from_string(path));
		//////////////

		// new record
		GLib.Value val;
		store.get_value (iter, (int) columns.length, out val);
		if (val.get_int () == 1) {
			store.set_value (iter, column_index, new_value);
			return;
		}
		/////////////


		// update record
		Benchwell.SQL.ColDef[] pks = {};
		string[] values = {};
		var col_id = 0;

		foreach (var column in columns) {
			if (!column.pk) {
				col_id++;
				continue;
			}
			store.get_value (iter, col_id, out val);

			pks += column;
			values += val.get_string ();

			col_id++;
		};

		if (pks.length == 0) {
			col_id = 0;
			foreach (var column in columns) {
				pks += column;
				values += val.get_string ();

				col_id++;
			};
		}

		// // append changing value
		pks += columns[column_index];
		values += new_value;

		store.set_value (iter, column_index, new_value);
		field_change (pks, values);
	}

	private bool on_button_press (Gtk.Widget w, Gdk.EventButton event) {
		if (event.button != Gdk.BUTTON_SECONDARY) {
			return false;
		}

		Gtk.TreePath path;
		table.get_path_at_pos ((int) event.x, (int) event.y , out path, null, null, null);

		table.get_selection ().select_path (path);

		menu.popup_at_pointer (event);
		return true;
	}
}

public class Benchwell.Views.DBCondition {
	public Gtk.Switch active_switch;
	public Gtk.ListStore store;
	public Gtk.ComboBox field_combo;
	public Gtk.ComboBoxText operator_combo;
	public Gtk.Entry value_entry;
	public Benchwell.Button remove_btn;
	private Benchwell.SQL.ColDef[] _columns;
	public Benchwell.SQL.ColDef[] columns {
		get { return _columns; }
		set {
			// keep active field before replace columns
			if (field_combo.sensitive) {
				Gtk.TreeIter iter;
				GLib.Value? val;
				field_combo.get_active_iter (out iter);
				if (store.iter_is_valid (iter)) {
					store.get_value (iter, 0, out val);
					if (val != null) {
						_active_field = val.get_string ();
					} else {
						_active_field = null;
					}
				}
			}

			_columns = value;
			_update_fields ();
		}
	}

	private string? _active_field;

	public signal void search();

	public DBCondition () {
		store = new Gtk.ListStore (1, GLib.Type.STRING);
		Gtk.TreeIter iter;
		store.append (out iter);
		store.set_value (iter, 0 ,"");

		active_switch = new Gtk.Switch ();
		active_switch.active = true;
		active_switch.valign = Gtk.Align.CENTER;
		active_switch.vexpand = false;
		active_switch.show ();

		field_combo = new Gtk.ComboBox.with_model_and_entry (store);
		field_combo.id_column = 0;
		field_combo.set_entry_text_column (0);
		field_combo.show ();

		var completion = new Gtk.EntryCompletion ();
		completion.text_column = 0;
		completion.inline_completion = true;
		completion.inline_selection = true;
		completion.minimum_key_length = 2;
		completion.set_model (store);

		var entry = field_combo.get_child () as Gtk.Entry;
		entry.set_completion (completion);

		operator_combo = new Gtk.ComboBoxText ();
		foreach (var op in Benchwell.SQL.Operator.all ()) {
			operator_combo.append (op.to_string(), op.to_string ());
		}
		operator_combo.set_active (0);
		operator_combo.show ();

		value_entry = new Gtk.Entry ();
		value_entry.show ();

		remove_btn = new Benchwell.Button ("close", Gtk.IconSize.MENU);
		remove_btn.show ();

		operator_combo.changed.connect (() => {
			var op = Benchwell.SQL.Operator.parse (operator_combo.get_active_text ());
			value_entry.sensitive = true;

			if (op == Benchwell.SQL.Operator.IsNull) {
				value_entry.sensitive = false;
			}
			if (op == Benchwell.SQL.Operator.IsNotNull) {
				value_entry.sensitive = false;
			}
		});

		entry.activate.connect( () => { search (); });
		value_entry.activate.connect( () => { search (); });
	}

	public Benchwell.SQL.CondStmt? get_condition () {
		if (!active_switch.get_active () || !active_switch.sensitive) {
			return null;
		}

		Gtk.TreeIter? iter;
		field_combo.get_active_iter (out iter);
		if (!store.iter_is_valid(iter)) {
			return null;
		}
		GLib.Value val;
		store.get_value (iter, 0, out val);
		var column_name = val.get_string ();
		if (column_name == "" || column_name == null) {
			return null;
		}
		Benchwell.SQL.ColDef? column = null;
		foreach (var c in columns) {
			if (c.name == column_name) {
				column = c;
				break;
			}
		}
		if (column == null) {
			return null;
		}

		var op = operator_combo.get_active_text ();
		if (op == null || op == "") {
			return null;
		}
		var operator = Benchwell.SQL.Operator.parse (op);
		if (operator == null) {
			return null;
		}

		var cvalue = value_entry.get_text ();
		if (cvalue == null && operator != Benchwell.SQL.Operator.IsNull && operator != Benchwell.SQL.Operator.IsNotNull) {
			return null;
		}

		var stmt = new Benchwell.SQL.CondStmt ();
		stmt.field = column;
		stmt.op = operator;
		stmt.val = cvalue;

		return stmt;
	}

	private void _update_fields () {
		store.clear ();

		bool enable = _active_field == "" || _active_field == null;
		foreach (var column in _columns) {
			Gtk.TreeIter iter;
			store.append (out iter);
			store.set_value (iter, 0, column.name);
			if (column.name == _active_field) {
				enable = true;
				field_combo.set_active_iter (iter);
			}
		}

		active_switch.sensitive = enable;
		field_combo.sensitive = enable;
		operator_combo.sensitive = enable;
		value_entry.sensitive = enable;
	}
}

public class Benchwell.Views.DBConditions : Gtk.Grid {
	private List<Benchwell.Views.DBCondition> conditions;
	private Benchwell.SQL.ColDef[] _columns;
	public Benchwell.SQL.ColDef[] columns {
		get { return _columns; }
		set {
			_columns = value;
			conditions.foreach ((condition) => {
				condition.columns = _columns;
			});
		}
	}

	public signal void search();

	public DBConditions () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);
		set_name ("conditions");
		add_condition ();
	}

	public void add_condition () {
		var cond = new Benchwell.Views.DBCondition ();
		cond.columns = _columns;

		var y = (int) conditions.length ();
		attach (cond.field_combo, 0, y, 2, 1);
		attach (cond.operator_combo, 2, y, 1, 1);
		attach (cond.value_entry, 3, y, 2, 1);
		attach (cond.active_switch, 5, y, 2, 1);
		attach (cond.remove_btn, 7, y, 1, 1);

		conditions.append (cond);

		cond.remove_btn.clicked.connect ( () => {
			var index = conditions.index (cond);
			remove_row (index);
			conditions.remove (cond);

			if (conditions.length () == 0) {
				add_condition ();
			}
		});

		cond.active_switch.grab_focus.connect (() => {
			if (conditions.index (cond) != conditions.length () - 1) {
				return;
			}

			add_condition ();
		});

		var entry = cond.field_combo.get_child () as Gtk.Entry;
		entry.grab_focus.connect (() => {
			if (conditions.index (cond) != conditions.length () - 1) {
				return;
			}

			add_condition ();
		});

		cond.value_entry.grab_focus.connect (() => {
			if (conditions.index (cond) != conditions.length () - 1) {
				return;
			}

			add_condition ();
		});

		cond.search.connect ( () => {
			search ();
		});
	}

	public Benchwell.SQL.CondStmt[] get_conditions () {
		Benchwell.SQL.CondStmt[] stmts = {};

		conditions.foreach( (condition) => {
			var c = condition.get_condition ();
			if (c != null) {
				stmts += c;
			}
		});

		return stmts;
	}
}

public class Benchwell.Views.DBTables : Gtk.ListBox {
	private Gtk.Menu menu;
	public Gtk.MenuItem edit_menu;
	public Gtk.MenuItem new_tab_menu;
	public Gtk.MenuItem schema_menu;
	public Gtk.MenuItem truncate_menu;
	public Gtk.MenuItem delete_menu;
	public Gtk.MenuItem refresh_menu;
	//public Gtk.MenuItem copy_select_menu;

	public Benchwell.SQL.TableDef[] _tables;

	public Regex? filter;

	public Benchwell.SQL.TableDef? selected_tabledef {
		get {
			var row = get_selected_row ();
			if (row.get_index () < 0) {
				return null;
			}
			return _tables[row.get_index ()];
		}
	}

	public DBTables () {
		Object ();

		menu = new Gtk.Menu ();
		edit_menu = new Benchwell.MenuItem (_("Edit"), "edit-table");
		edit_menu.show ();

		new_tab_menu = new Benchwell.MenuItem (_("New tab"), "add-tab");
		new_tab_menu.show ();

		schema_menu = new Benchwell.MenuItem (_("Schema"), "config");
		schema_menu.show ();

		truncate_menu = new Benchwell.MenuItem (_("Truncate"), "truncate");
		truncate_menu.show ();

		delete_menu = new Benchwell.MenuItem (_("Delete"), "delete-table");
		delete_menu.show ();

		refresh_menu = new Benchwell.MenuItem (_("Refresh"), "refresh");
		refresh_menu.show ();

		//copy_select_menu = new Benchwell.MenuItem (_("Copy SELECT"), "copy");
		//copy_select_menu.show ();

		var cowboy = new Benchwell.MenuItem (_("Cowboy"), "cowboy");
		cowboy.show ();

		menu.add (new_tab_menu);
		//menu.add (copy_select_menu);
		menu.add (schema_menu);
		menu.add (edit_menu);
		menu.add (refresh_menu);
		menu.add (cowboy);

		var cowboy_menu = new Gtk.Menu ();
		cowboy_menu.add (truncate_menu);
		cowboy_menu.add (delete_menu);
		cowboy.set_submenu (cowboy_menu);

		button_press_event.connect ( (list, event) => {
			if (event.button == Gdk.BUTTON_SECONDARY) {
				grab_focus ();
				select_row (get_row_at_y ((int)event.y));
			}

			return false;
		});

		button_press_event.connect ((list, event) => {
			if ( event.button != Gdk.BUTTON_SECONDARY){
				return false;
			}

			menu.show ();
			menu.popup_at_pointer (event);
			return true;
		});

		set_filter_func (search);
	}

	public void update_items (owned Benchwell.SQL.TableDef[] tables, string name = "") {
		_tables = (owned) tables;

		get_children().foreach( (row) => {
			remove (row);
		});

		foreach (var item in _tables) {
			var row = build_row (item);
			add (row);
			if (item.name == name) {
				select_row (row);
			}
		};
	}

	private Gtk.ListBoxRow build_row (Benchwell.SQL.TableDef def) {
		var row = new Gtk.ListBoxRow ();
		row.show ();

		var box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		box.show ();

		var label = new Gtk.Label (def.to_string());
		label.set_halign (Gtk.Align.START);
		label.show ();

		var icon_name = "table";
		if (def.ttype == Benchwell.SQL.TableType.Dummy) {
			icon_name = "table-v";
		}
		var image = new Benchwell.Image (icon_name, Gtk.IconSize.BUTTON);
		image.show ();

		box.pack_start (image, false, false, 5);
		box.pack_start (label, false, false, 0);

		row.add (box);

		return row;
	}

	public bool search (Gtk.ListBoxRow row) {
		if ( filter == null ) {
			return true;
		}

		var box = (Gtk.Box) row.get_child();

		var lbl = (Gtk.Label) box.get_children().nth_data (1);
		return filter.match (lbl.get_label ());
	}

	public unowned Benchwell.SQL.TableDef get_selected_table () {
		var row = get_selected_row ();
		return _tables[row.get_index ()];
	}

	public void remove_selected () {
		var row = get_selected_row ();
		var index = row.get_index ();
		if (index < 0) {
			return;
		}

		remove (row);

		Benchwell.SQL.TableDef[] new_tables = {};
		for (var i = 0; i < _tables.length; i++ ){
			if (i == index) {
				continue;
			}
			new_tables += _tables[i];
		}

		_tables = new_tables;
	}
}

public class Benchwell.Views.CancelOverlay : Gtk.Overlay {
	public delegate void OnCancelFunc ();

	public Gtk.Button btn_cancel;
	public Gtk.Spinner spinner;
	public Gtk.Box controls;
	private OnCancelFunc cancel;
	public Gtk.Widget overlayed { construct; }

	public CancelOverlay (Gtk.Widget overlayed) {
		Object(
			overlayed: overlayed
		);

		btn_cancel = new Gtk.Button.with_label (_("Cancel"));
		btn_cancel.set_size_request (100, 30);
		btn_cancel.show ();

		spinner = new Gtk.Spinner ();

		controls = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);

		var actions = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		actions.set_size_request (100, 150);
		actions.valign = Gtk.Align.CENTER;
		actions.halign = Gtk.Align.CENTER;
		actions.vexpand = true;
		actions.hexpand = true;
		actions.show ();

		actions.pack_start (spinner, true, true, 0);
		actions.pack_start (btn_cancel, false, false, 0);
		controls.add (actions);

		add (overlayed);

		btn_cancel.clicked.connect ( () => {
			stop ();
			cancel ();
		});
	}

	public void run (OnCancelFunc c) {
		controls.show ();
		spinner.show ();
		add_overlay (controls);
		cancel = c;
	}

	public void stop () {
		remove (controls);
		spinner.stop ();
	}
}
