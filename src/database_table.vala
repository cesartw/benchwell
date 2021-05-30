public class Benchwell.Database.Table : Gtk.Box {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.DatabaseService service { get; construct; }
	private bool _raw_mode = false;
	public bool raw_mode {
		get { return _raw_mode; }
		set {
			_raw_mode = value;
			btn_prev.sensitive = !_raw_mode;
			btn_next.sensitive = !_raw_mode;
			btn_show_filters.sensitive = !_raw_mode;
			btn_add_row.sensitive = !_raw_mode;
			btn_refresh.sensitive = !_raw_mode;
			btn_show_filters.active = !_raw_mode;
		}
	}
	public Gtk.TreeView table;
	public Benchwell.Database.EditPanel edit_panel;
	public Gtk.Menu menu;
	public Gtk.MenuItem clone_menu;
	public Gtk.MenuItem copy_insert_menu;
	public Gtk.MenuItem copy_menu;
	public Gtk.ListStore store;

	public Gtk.Button btn_load_query;
	public Gtk.MenuButton save_menu;
	public Gtk.Button btn_prev;
	public Gtk.Button btn_next;
	public Gtk.Button btn_refresh;
	public Gtk.ToggleButton btn_show_filters;
	public Gtk.Button btn_add_row;
	public Gtk.Button btn_delete_row;
	public Gtk.Button btn_save_row;
	public Gtk.SearchEntry search;

	public Benchwell.Database.Conditions conditions;

	public signal void field_changed (Benchwell.Column[] columns);
	public signal bool delete_record (Benchwell.Column[]? data);
	public signal void file_opened (string query);
	public signal void file_saved (string filename);
	public signal void fav_saved (string query_name);

	private bool disable_selection;

	public Table (Benchwell.ApplicationWindow window, Benchwell.DatabaseService service) {
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
		table.get_selection ().set_mode (Gtk.SelectionMode.MULTIPLE);
		table.show ();

		var table_sw = new Gtk.ScrolledWindow (null, null);
		table_sw.add (table);
		table_sw.show ();

		var table_edit_paned = new Gtk.Paned (Gtk.Orientation.HORIZONTAL);
		table_edit_paned.wide_handle = true;
		table_edit_paned.show ();

		if (Config.settings.db_edit_panel) {
			table_edit_paned.pack2 (build_edit_panel (), true, true);
		}

		table_edit_paned.pack1 (table_sw, true, false);

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
		btn_add_row = new Benchwell.Button ("add-record", Gtk.IconSize.BUTTON) {
			sensitive = false
		};
		btn_add_row.show ();

		btn_delete_row = new Benchwell.Button ("delete-record", Gtk.IconSize.BUTTON) {
			sensitive = false
		};
		btn_delete_row.show ();

		btn_save_row = new Benchwell.Button ("save-record", Gtk.IconSize.BUTTON) {
			sensitive = false
		};
		btn_save_row.sensitive = false;
		btn_save_row.show ();

		btn_show_filters = new Benchwell.ToggleButton ("filter", Gtk.IconSize.BUTTON) {
			sensitive = false
		};
		btn_show_filters.show ();

		btn_prev = new Benchwell.Button ("back", Gtk.IconSize.BUTTON) {
			sensitive = false
		};
		btn_prev.show ();

		btn_next = new Benchwell.Button ("next", Gtk.IconSize.BUTTON) {
			sensitive = false
		};
		btn_next.show ();

		btn_refresh = new Benchwell.Button ("refresh", Gtk.IconSize.BUTTON) {
			sensitive = false
		};
		btn_refresh.show ();

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

		action_save_file.activate.connect (on_save_file);
		action_save_fav.activate.connect (on_save_fav);
		btn_load_query.clicked.connect (on_open_file);

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
		table_actionbar.pack_end (btn_load_query);
		table_actionbar.pack_end (save_menu);
		table_actionbar.pack_end (btn_refresh);
		table_actionbar.pack_end (btn_next);
		table_actionbar.pack_end (btn_prev);
		////////////

		conditions = new Benchwell.Database.Conditions ();

		pack_start (table_actionbar, false, false, 0);
		pack_start (conditions, false,false, 5);
		pack_start (table_edit_paned, true, true, 0);

		// signals
		table.button_press_event.connect (on_button_press);
		table.key_press_event.connect (on_table_key_press);
		table.get_selection ().changed.connect (on_selection_changed);

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

		clone_menu.activate.connect (on_clone);
	}

	private void on_selection_changed () {
		if (disable_selection)
			return;

		btn_delete_row.sensitive = btn_save_row.sensitive = table.get_selection ().count_selected_rows () != 0;

		if (table.get_selection ().count_selected_rows () == 0)
			return;

		var data = get_selected_data ();
		if (data == null)
			return;

		if (table.get_selection ().count_selected_rows () > 1) {
			edit_panel.get_parent ().hide ();
			return;
		}

		edit_panel.set_record (data);
		edit_panel.get_parent ().show ();
		edit_panel.show ();
	}

	private bool on_table_key_press (Gtk.Widget widget, Gdk.EventKey event) {
		if (event.keyval == Gdk.Key.c && event.state == Gdk.ModifierType.CONTROL_MASK) {
			var selection = table.get_selection ();

			var builder = new StringBuilder ();
			Gtk.TreeIter? iter = null;
			selection.get_selected_rows (null).foreach ((path) => {
				var ok = store.get_iter (out iter, path);
				if (!ok) {
					return;
				}
				for (var i = 0; i < service.columns.length; i++) {
					GLib.Value val;
					store.get_value (iter, i, out val);
					var str = val.get_string ();

					str = str.replace ("\\t", "\\\\t").replace ("\t", "\\t");
					builder.append (str);

					if (i == service.columns.length - 1)
						builder.append ("\n");
					else
						builder.append ("\t");
				}
			});

			var cb = Gtk.Clipboard.get_default (Gdk.Display.get_default ());
			cb.set_text (builder.str, builder.str.length);
		}

		if (event.keyval == Gdk.Key.Escape) {
			table.get_selection ().unselect_all ();
			return true;
		}

		if (event.keyval == Gdk.Key.Delete) {
			delete_rows ();
			return true;
		}

		return false;
	}

	public Gtk.Widget build_edit_panel () {
		edit_panel = new Benchwell.Database.EditPanel ();
		edit_panel.show ();

		var edit_panel_sw = new Gtk.ScrolledWindow (null, null);
		edit_panel_sw.add (edit_panel);
		edit_panel_sw.show ();
		return edit_panel_sw;
	}

	public void on_open_file () {
		var dialog = new Gtk.FileChooserDialog (_("Select file"), window,
											 Gtk.FileChooserAction.OPEN,
											_("Open"), Gtk.ResponseType.OK,
											_("Cancel"), Gtk.ResponseType.CANCEL);
		var resp = (Gtk.ResponseType) dialog.run ();

		if (resp == Gtk.ResponseType.CANCEL) {
			dialog.destroy ();
			return;
		}

		var filename = dialog.get_filename ();
		dialog.destroy ();

		try {
			string text;
			var ok = GLib.FileUtils.get_contents (filename, out text, null);
			if (!ok) {
				return;
			}

			file_opened (text);
		} catch( GLib.FileError err) {
			Config.show_alert (this, err.message);
		}
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

		file_saved (dialog.get_filename ());
		dialog.destroy ();
	}

	public void on_save_fav () {
		var query_name = ask_fav_name ();
		if (query_name == "") {
			return;
		}

		fav_saved (query_name);
	}

	private string ask_fav_name () {
		var dialog = new Gtk.Dialog.with_buttons (_("Choose"), window,
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

	public Benchwell.SortOption[] get_sort_options () {
		Benchwell.SortOption[] sorts = {};

		for (var i = 0; i < table.get_n_columns (); i++) {
			var col = table.get_column (i);

			if (!col.get_sort_indicator ()) {
				continue;
			}

			switch (col.get_sort_order ()) {
			case Gtk.SortType.DESCENDING:
				sorts += new Benchwell.SortOption(service.columns[i], Benchwell.SortType.Asc);
				break;
			case Gtk.SortType.ASCENDING:
				sorts += new Benchwell.SortOption(service.columns[i], Benchwell.SortType.Desc);
				break;
			}
		}

		return sorts;
	}

	public Benchwell.CondStmt[] get_conditions () {
		return conditions.get_conditions ();
	}

	public void load_table () {
		disable_selection = true;
		edit_panel.clear ();
		table.get_selection ().unselect_all ();
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
		disable_selection = false;
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
		if (store == null) {
			return;
		}
		store.clear ();
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
		table.row_activated (path, table.get_column(0));
		table.scroll_to_cell (path, null, true, (float) 0.5, (float) 0);
	}

	public Benchwell.Column[]? get_selected_data () {
		Gtk.TreeIter? selected = null;
		var selection = table.get_selection ();
		if (selection.count_selected_rows () == 0)
			return null;

		selection.selected_foreach ( (model, path, iter) => {
			if (selected != null){
				return;
			}
			store.get_iter (out selected, path);
		});

		if (selected == null) {
			return null;
		}

		//var values = new string[service.columns.length];
		var values = new Benchwell.Column[service.columns.length];
		for (var i = 0; i < service.columns.length; i++) {
			GLib.Value val;
			store.get_value (selected, i, out val);

			values[i] = new Benchwell.Column ();
			values[i].val = val.get_string ();
			values[i].coldef = service.columns[i];
		}

		return values;
	}

	public Benchwell.Column[]? get_data_at (Gtk.TreeIter iter) {
		var values = new Benchwell.Column[service.columns.length];
		for (var i = 0; i < service.columns.length; i++) {
			GLib.Value val;
			store.get_value (iter, i, out val);

			values[i] = new Benchwell.Column ();
			values[i].coldef = service.columns[i];
			values[i].val = val.get_string ();
		}

		return values;
	}

	public void delete_selected_row () {
		var selection = table.get_selection ();
		var tm = store as Gtk.TreeModel;
		var paths = selection.get_selected_rows (out tm);

		paths.reverse ();
		paths.foreach ((path) => {
			Gtk.TreeIter iter;
			store.get_iter (out iter, path);
			store.remove (ref iter);
		});
	}

	public void update_selected_row (Benchwell.Column[] data) {
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
			store.set_value (selected, i, data[i].val);
		}
	}

	private Gtk.TreeViewColumn build_column(Benchwell.ColDef column, int column_index) {
		var renderer = new Gtk.CellRendererText ();
		renderer.editable = true;
		renderer.xpad = 10;
		renderer.height = 23;
		renderer.cell_background = Benchwell.HighlightColors.PKHL.to_string ();
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

			if ( val.holds (GLib.Type.STRING) ) {
				if ( val.get_string () == Benchwell.null_string ){
					cell.set_property ("markup", @"<span foreground=\"$(Benchwell.HighlightColors.NULLHL.to_string ())\">&lt;NULL&gt;</span>");
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
		var columns = new Benchwell.Column[service.columns.length + 1];
		for (var i = 0; i < service.columns.length; i++) {
			store.get_value (iter, i, out val);
			columns[i] = new Benchwell.Column ();
			columns[i].coldef = service.columns[i];
			columns[i].val = val.get_string ();

			if (column_index == i && val.get_string () == new_value) {
				return;
			}
		};

		columns[service.columns.length] = new Benchwell.Column ();
		columns[service.columns.length].coldef = service.columns[column_index];
		columns[service.columns.length].val = new_value;

		store.set_value (iter, column_index, new_value);
		field_changed (columns);
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
		delete_rows ();
	}

	private void delete_rows () {
		GLib.List<Gtk.TreeRowReference> references = null;
		table.get_selection ().get_selected_rows (null).foreach((path) => {;
			references.append(new Gtk.TreeRowReference (store, path));
		});

		references.foreach ((tref) => {
			var path = tref.get_path ();
			if (path == null)
				return;

			Gtk.TreeIter iter;
			var ok = store.get_iter (out iter, path);
			if (!ok)
				return;

			GLib.Value val;
			store.get_value (iter, (int) service.columns.length, out val);

			// new unsaved record
			if (val.get_int () == 1) {
				store.remove (ref iter);
				return;
			}

			var data = get_data_at (iter);
			if (delete_record (data)) {
				store.remove (ref iter);
			}
		});
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
		Regex regex = null;
		try {
			regex = new Regex (expr, RegexCompileFlags.CASELESS);
		} catch (RegexError e) {
			try { regex = new Regex (".*"); } catch (RegexError err) {};
		}


		table.get_columns ().foreach ((column) => {
			column.visible = regex.match (column.get_title ().replace ("__", "_"));
		});
	}
}

public class Benchwell.Database.EditPanel : Gtk.Box {
	public Gtk.Button btn_save;

	public EditPanel () {
		Object (
			orientation: Gtk.Orientation.VERTICAL
		);

		btn_save = new Gtk.Button.with_label (_("Save"));
		btn_save.get_style_context ().add_class ("suggested-action");
	}

	public void clear () {
		get_children ().foreach ((w) => {
			w.destroy ();
		});
	}

	public void set_record (Benchwell.Column[] columns) {
		clear ();

		for (var i = 0; i < columns.length; i++) {
			var label = new Gtk.Label (@"<b>$(columns[i].coldef.name)</b>") {
				halign = Gtk.Align.START,
				use_markup = true
			};
			pack_start (label, false, false, 5);

			Gtk.Widget input;
			columns[i].coldef.ttype == Benchwell.ColType.LongString;
			if (columns[i].coldef.ttype == Benchwell.ColType.LongString || (columns[i].val != null && columns[i].val.index_of ("\n") == -1)) {
				input = new Gtk.Entry () {
					text = columns[i].val
				};
			} else {
				var tx = new Gtk.TextView ();
				input = tx;
				if (columns[i].val != null)
					tx.get_buffer ().set_text (columns[i].val);
			}

			label.show ();
			input.show ();

			pack_start (input, false, false, 5);
		}
	}
}
