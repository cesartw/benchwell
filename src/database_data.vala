public class Benchwell.Database.Data : Gtk.Paned {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.DatabaseService service { get; construct; }

	public Gtk.SearchEntry table_search;
	public Gtk.ComboBoxText database_combo;
	public Benchwell.Database.Tables tables;
	public Benchwell.Database.ResultView result_view;
	public Gtk.ListBox history;

	private Benchwell.Views.CancelOverlay overlay;
	private List<string> databases;

	private int current_page = 0;
	private int page_size = 100;

	public signal void database_selected(string dbname);

	public Data (Benchwell.ApplicationWindow window, Benchwell.DatabaseService service) {
		Object(
			window: window,
			service: service,
			orientation: Gtk.Orientation.VERTICAL,
			wide_handle: true,
			vexpand: true,
			hexpand: true
		);

		build ();

		try {
			service.info.load_history ();
		} catch (Benchwell.ConfigError err) {
			stderr.printf ("loading history: %s", err.message);
		}

		foreach (Benchwell.Query query in service.info.history) {
			add_history_row (query);
		}

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

		tables.table_selected.connect (on_load_table);

		result_view.table.field_changed.connect (on_field_changed);
		result_view.table.btn_refresh.clicked.connect (on_refresh_table);
		result_view.table.conditions.search.connect (on_refresh_table);
		result_view.table.conditions.ready.connect (on_condition_ready);
		result_view.fav_saved.connect (on_refresh_tables);

		tables.schema_menu.activate.connect (on_show_schema);
		tables.refresh_menu.activate.connect (on_refresh_tables);
		tables.truncate_menu.activate.connect (on_truncate_table);
		tables.delete_menu.activate.connect (on_delete_table);

		result_view.table.btn_prev.clicked.connect (on_prev_page);
		result_view.table.btn_next.clicked.connect (on_next_page);

		result_view.table.copy_insert_menu.activate.connect (on_copy_insert);

		fill ();
	}

	private void on_condition_ready (Benchwell.CondStmt stmt) {
		Config.save_filters (service.info, service.table_def.name, result_view.table.conditions.get_conditions ());
	}

	private void on_copy_insert () {
		var data = result_view.table.get_selected_data ();
		if (data == null) {
			return;
		}

		var st = service.connection.get_insert_statement (service.table_def.name, service.columns, data);
		var cb = Gtk.Clipboard.get_default (Gdk.Display.get_default ());
		cb.set_text (st, st.length);
	}

	private void on_show_schema () {
		var tabledef = tables.selected_tabledef;
		if (tabledef == null) {
			return;
		}

		string sql = "";
		try {
			sql = service.connection.get_create_table (tabledef.name);
		} catch (Benchwell.Error err) {
			Config.show_alert (this, err.message);
			return;
		}
		var dialog = new Gtk.Dialog.with_buttons (@"$(tabledef.name) schema", window,
								Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
								_("Ok"), Gtk.ResponseType.OK);
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

		window.infobar.hide ();
		try {
			if (tabledef.ttype == Benchwell.TableType.Dummy) {
				var query = (Benchwell.Query) tabledef.source;
				service.info.remove_query (query);
			} else {
				service.delete_table (tabledef);
			}
			tables.remove_selected ();
		} catch (Benchwell.Error err) {
			window.show_alert (err.message);
		} catch (Benchwell.ConfigError err) {
			window.show_alert (err.message);
		}
	}

	private void on_truncate_table () {
		window.infobar.hide ();

		var tabledef = tables.selected_tabledef;
		if (tabledef == null) {
			return;
		}

		try {
			service.connection.truncate_table (tabledef);
		} catch (Benchwell.Error err) {
			window.show_alert (err.message);
			return;
		}
		window.show_alert (_("Done"), Gtk.MessageType.INFO, true);
	}

	private void build () {
		var hpaned = new Gtk.Paned (Gtk.Orientation.HORIZONTAL);
		hpaned.set_wide_handle (true);
		hpaned.set_hexpand (true);
		hpaned.set_vexpand (true);
		hpaned.show ();

		var sidebar = new Gtk.Box (Gtk.Orientation.VERTICAL, 5);
		sidebar.get_style_context ().add_class ("bw-spacing");
		sidebar.show ();

		table_search = new Gtk.SearchEntry ();
		table_search.set_placeholder_text (_("Search table: .*"));
		table_search.show ();

		tables = new Benchwell.Database.Tables (service);
		tables.activate_on_single_click = false;
		tables.show ();

		var tables_sw = new Gtk.ScrolledWindow (null, null);
		tables_sw.show ();
		tables_sw.add (tables);


		history = new Gtk.ListBox ();
		history.activate_on_single_click = false;
		history.show ();

		var history_sw = new Gtk.ScrolledWindow (null, null);
		history_sw.add (history);
		history_sw.show ();

		var tables_and_history = new Gtk.Paned (Gtk.Orientation.VERTICAL);
		tables_and_history.pack1 (tables_sw, true, true);
		tables_and_history.pack2 (history_sw, false, false);
		tables_and_history.show ();

		database_combo = new Gtk.ComboBoxText ();
		database_combo.set_id_column (0);
		database_combo.show ();

		sidebar.pack_start (database_combo, false, true, 0);
		sidebar.pack_start (table_search, false, true, 0);
		sidebar.pack_start (tables_and_history, true, true, 0);

		// main section
		var main_section = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		main_section.show ();

		result_view = new Benchwell.Database.ResultView (window, service);
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
			window.infobar.hide ();

			var data = result_view.table.get_selected_data ();
			if (data == null) {
				return;
			}

			try {
				data = service.connection.insert_record (service.table_def.name, service.columns, data);
				result_view.table.update_selected_row (data);
			} catch (Benchwell.Error err) {
				window.show_alert (err.message);
				return;
			}
		});

		result_view.table.delete_record.connect ((data) => {
			window.infobar.hide ();

			try {
				service.connection.delete_record (service.table_def.name, service.columns, data);
				//result_view.table.delete_selected_row ();
			} catch (Benchwell.Error err) {
				window.show_alert (err.message);
				return false;
			}

			return true;
		});

		//result_view.table.btn_delete_row.clicked.connect (() => {
			//result_view.infobar.hide ();

			//var data = result_view.table.get_selected_data ();
			//if (data == null) {
				//return;
			//}
			//try {
				//service.connection.delete_record (service.table_def.name, service.columns, data);
				//result_view.table.delete_selected_row ();
			//} catch (Benchwell.Error err) {
				//window.show_alert (err.message);
				//return;
			//}
		//});

		result_view.exec_query.connect (on_exec_query);

		history.row_activated.connect (on_history_activated);
	}

	private void on_history_activated () {
		var row = history.get_selected_row ();
		var query = service.info.history[service.info.history.length - row.get_index () - 1];
		result_view.editor.get_buffer ().set_text (query.query);
	}

	private void on_exec_query (string raw_query) {
		window.infobar.hide ();

		var interpolated = raw_query;
		interpolated = Config.environment.interpolate_variables (interpolated);
		interpolated = Config.environment.interpolate_functions (interpolated);

		try {
			var query = service.info.save_history (interpolated);
			add_history_row (query);
		} catch (Benchwell.ConfigError err) {
			stderr.printf ("saving history: %s", err.message);
		}

		try {
			string[] columns;
			List<List<string?>> data;
			service.connection.query(interpolated, out columns, out data);

			Benchwell.ColDef[] cols = {};
			foreach (var column in columns) {
				cols += new Benchwell.ColDef.with_name (column);
			}

			service.columns = cols;
			service.data = (owned) data;
			result_view.table.load_table ();
			result_view.table.raw_mode = true;
		} catch (Benchwell.Error err) {
			window.show_alert (err.message);
			return;
		}
	}

	private void add_history_row (owned Benchwell.Query query) {
		var row = new Gtk.ListBoxRow ();
		row.show ();

		var box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		box.show ();

		var time_fmt = "%Y-%m-%d %H:%M:%S";
		var now = new DateTime.now_local ();
		if (now.get_day_of_year () == query.created_at.get_day_of_year ())
			time_fmt = "%H:%M:%S";

		var time_lbl = new Gtk.Label (query.created_at.format (time_fmt));
		time_lbl.set_halign (Gtk.Align.START);
		time_lbl.show ();

		var q = query.query.replace("\n", " ");
		var length = q.length;
		if (length > 25)
			length = 25;
		var query_lbl = new Gtk.Label (q.substring (0, length));
		query_lbl.set_halign (Gtk.Align.START);
		query_lbl.show ();

		box.pack_start (time_lbl, false, false, 5);
		box.pack_start (query_lbl, false, false, 0);

		row.add (box);

		history.prepend (row);
	}

	private void fill () {
		databases = service.connection.databases ();
		databases.foreach ( db => {
			database_combo.append (db, db);
		});

		if ( service.info.database != "" ) {
			database_combo.set_active_id (service.info.database);
		}
	}

	private void on_database_selected () {
		window.infobar.hide ();
		var dbname = database_combo.get_active_text ();
		try {
			service.use_database (dbname);
			tables.update_items ();
			database_selected (dbname);
		} catch (Benchwell.Error err) {
			window.show_alert (err.message);
			return;
		}
		result_view.table.clear ();
		window.show_alert (_("Using %s").printf (dbname), Gtk.MessageType.INFO, true);
	}

	private void on_field_changed(Benchwell.ColDef[] columns, string[] row) {
		window.infobar.hide ();
		if (service.table_def == null) {
			window.show_alert (_("No table selected"), Gtk.MessageType.ERROR);
			return;
		}
		try {
			service.connection.update_field (service.table_def.name, columns, row);
		} catch (Benchwell.Error err) {
			window.show_alert (err.message);
			return;
		}
		window.show_alert (_("Updated"), Gtk.MessageType.INFO, true);
	}

	private void on_load_table (Benchwell.TableDef _table_def) {
		window.infobar.hide ();
		service.table_def = _table_def;

		current_page = 0;
		try {
			service.load_table(null, null, current_page, page_size);
			result_view.table.raw_mode = false;
			result_view.table.conditions.columns = service.columns;
			result_view.table.load_table ();

			var filters = Config.get_table_filters (service.info, service.table_def.name);
			result_view.table.conditions.rebuild (filters);
		} catch (Benchwell.Error err) {
			window.show_alert (err.message);
			return;
		}

		window.show_alert (_("Loaded"), Gtk.MessageType.INFO, true);
	}

	private void on_refresh_table () {
		window.infobar.hide ();

		if (result_view.table.raw_mode) {
			result_view._exec_query ();
		} else {
			if (service.table_def == null ) {
				return;
			}

			try {
				service.load_table(result_view.table.get_conditions (),
								result_view.table.get_sort_options (),
								current_page, page_size);
				result_view.table.raw_mode = false;
				result_view.table.refresh_data ();
			} catch (Benchwell.Error err) {
				window.show_alert (err.message);
				return;
			}

			window.show_alert (_("Refresh"), Gtk.MessageType.INFO, true);
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
