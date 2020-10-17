public class Benchwell.Database.Table : Gtk.Box {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.Services.Database service { get; construct; }
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
	public Benchwell.Database.Conditions conditions;

	public signal void field_change (Benchwell.Backend.Sql.ColDef[] column, string[] row);

	public Table (Benchwell.ApplicationWindow window, Benchwell.Services.Database service) {
		Object (
			window: window,
			service: service,
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

		conditions = new Benchwell.Database.Conditions ();

		pack_start (table_actionbar, false, false, 0);
		pack_start (conditions, false,false, 5);
		pack_start (table_sw, true, true, 0);

		// signals
		table.button_press_event.connect (on_button_press);

		btn_show_filters.toggled.connect (() => {
			if (btn_show_filters.active) {
				conditions.show ();
			} else {
				conditions.hide ();
			}
		});

		//column
		search.search_changed.connect (on_column_search);

		btn_add_row.clicked.connect (add_empty_row);

		btn_delete_row.clicked.connect (on_delete_row);
	}

	public Benchwell.Backend.Sql.SortOption[] get_sort_options () {
		Benchwell.Backend.Sql.SortOption[] sorts = {};

		for (var i = 0; i < table.get_n_columns (); i++) {
			var col = table.get_column (i);

			if (!col.get_sort_indicator ()) {
				continue;
			}

			switch (col.get_sort_order ()) {
			case Gtk.SortType.DESCENDING:
				sorts += new Benchwell.Backend.Sql.SortOption(service.columns[i], Benchwell.Backend.Sql.SortType.Asc);
				break;
			case Gtk.SortType.ASCENDING:
				sorts += new Benchwell.Backend.Sql.SortOption(service.columns[i], Benchwell.Backend.Sql.SortType.Desc);
				break;
			}
		}

		return sorts;
	}

	public Benchwell.Backend.Sql.CondStmt[] get_conditions () {
		return conditions.get_conditions ();
	}

	public void load_table () {
		if (store != null) {
			store.clear ();
		}

		while ( table.get_column(0) != null ) {
			table.remove_column (table.get_column(0));
		}

		GLib.Type[] column_types = new GLib.Type[service.columns.length + 1];
		var i = 0;
		foreach (var column in service.columns) {
			table.insert_column (build_column (column, i ), i);
			column_types[i] = GLib.Type.STRING;
			i++;
		};

		column_types[service.columns.length] = GLib.Type.INT;

		store = new Gtk.ListStore.newv (column_types);
		table.model = store;

		service.data.foreach ( (row) => {
			add_row (row);
		});
	}

	public void refresh_data () {
		if (store != null) {
			store.clear ();
		}

		service.data.foreach ( (row) => {
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

	public void clear () {
		if (store != null) {
			store.clear ();
		}
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
		for (i = 0; i < service.columns.length; i++) {
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

		var values = new string[service.columns.length];
		for (var i = 0; i < service.columns.length; i++) {
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

	private Gtk.TreeViewColumn build_column(Benchwell.Backend.Sql.ColDef column, int column_index) {
		var renderer = new Gtk.CellRendererText ();
		renderer.editable = true;
		renderer.xpad = 10;
		renderer.height = 23;
		renderer.cell_background = Benchwell.Colors.PKHL.to_string ();
		renderer.cell_background_set = column.pk;
		renderer.ellipsize = Pango.EllipsizeMode.END;

		var _column = new Gtk.TreeViewColumn.with_attributes (column.name.replace ("_", "__"), renderer, "text", column_index);
		_column.resizable = true;
		_column.clickable = true;
		_column.min_width = 85;
		//_column.max_width = 250;

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
		_column.set_cell_data_func (renderer, (cell_layout, cell, tree_model, iter) => {
			GLib.Value val;
			var index = column_index;
			tree_model.get_value (iter, index, out val);
			var path = tree_model.get_path (iter);

			if ( val.holds (GLib.Type.STRING) ) {
				if ( val.get_string () == Benchwell.null_string ){
					cell.set_property ("markup", @"<span foreground=\"$(Benchwell.Colors.NULLHL.to_string ())\">&lt;NULL&gt;</span>");
				}
			}
		});

		renderer.edited.connect ((cell, path, new_value) => {
			on_edited (cell, path, new_value, column_index);
		});

		return _column;
	}

	private void on_edited (Gtk.CellRendererText cell, string path, string new_value, int column_index) {
		// update cell
		Gtk.TreeIter iter ;
		store.get_iter (out iter, new Gtk.TreePath.from_string(path));
		//////////////

		// new record
		GLib.Value val;
		store.get_value (iter, (int) service.columns.length, out val);
		if (val.get_int () == 1) {
			store.set_value (iter, column_index, new_value);
			return;
		}
		/////////////


		// update record
		Benchwell.Backend.Sql.ColDef[] pks = {};
		string[] values = {};
		var col_id = 0;

		foreach (var column in service.columns) {
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
			foreach (var column in service.columns) {
				pks += column;
				values += val.get_string ();

				col_id++;
			};
		}

		// // append changing value
		pks += service.columns[column_index];
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
		store.get_value (iter, (int) service.columns.length, out val);
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
		for (var i = 0; i < service.columns.length; i++) {
			if (service.columns[i].pk) {
				continue;
			}
			store.set_value (iter, i, data[i]);
		}
		store.set_value (iter, service.columns.length, 1);
		selection.select_iter (iter);
	}

	private void on_column_search () {
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
	}
}
