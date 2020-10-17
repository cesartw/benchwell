public class Benchwell.Database.Tables : Gtk.ListBox {
	public Benchwell.Services.Database service { get; construct; }
	private Gtk.Menu menu;
	public Gtk.MenuItem edit_menu;
	public Gtk.MenuItem new_tab_menu;
	public Gtk.MenuItem schema_menu;
	public Gtk.MenuItem truncate_menu;
	public Gtk.MenuItem delete_menu;
	public Gtk.MenuItem refresh_menu;
	public Regex? filter;
	public Benchwell.Backend.Sql.TableDef? selected_tabledef {
		get {
			var row = get_selected_row ();
			if (row.get_index () < 0) {
				return null;
			}
			return service.tables[row.get_index ()];
		}
		set {
			if (value == null) {
				return;
			}

			for (var i = 0; i < service.tables.length; i++){
				if (service.tables[i].name == value.name) {
					var row = get_row_at_index (i);
					select_row (row);
					table_selected (value);
					return;
				}
			}
		}
	}
	public signal void table_selected (Benchwell.Backend.Sql.TableDef tabledef);

	public Tables (Benchwell.Services.Database service) {
		Object (
			service: service
		);

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


		row_activated.connect (on_row_activated);
		set_filter_func (search);
	}

	private void on_row_activated () {
		table_selected (get_selected_table ());
	}

	public void update_items (string name = "") {
		get_children().foreach( (row) => {
			remove (row);
		});

		foreach (var item in service.tables) {
			var row = build_row (item);
			add (row);
			if (item.name == name) {
				select_row (row);
			}
		};
	}

	private Gtk.ListBoxRow build_row (Benchwell.Backend.Sql.TableDef def) {
		var row = new Gtk.ListBoxRow ();
		row.show ();

		var box = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 0);
		box.show ();

		var label = new Gtk.Label (def.to_string());
		label.set_halign (Gtk.Align.START);
		label.show ();

		var icon_name = "table";
		if (def.ttype == Benchwell.Backend.Sql.TableType.Dummy) {
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

	public unowned Benchwell.Backend.Sql.TableDef get_selected_table () {
		var row = get_selected_row ();
		return service.tables[row.get_index ()];
	}

	public void remove_selected () {
		var row = get_selected_row ();
		var index = row.get_index ();
		if (index < 0) {
			return;
		}

		remove (row);
	}
}

