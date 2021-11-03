public class Benchwell.Http.SaveAsSelector : Gtk.Grid {
	public Benchwell.ApplicationWindow   window { get; construct; }
	private CollectionsComboBox collections_combo;
	private Gtk.ComboBoxText folder_combo;
	private Gtk.Entry name_entry;

	public signal void changed ();

	public SaveAsSelector (Benchwell.ApplicationWindow window) {
		Object(
			window: window,
			row_spacing: 5,
			column_spacing: 5
		);

		collections_combo = new CollectionsComboBox (window);
		collections_combo.show ();

		folder_combo = new Gtk.ComboBoxText ();
		folder_combo.show ();
		folder_combo.sensitive = Config.settings.http_collection_id != 0;
		folder_combo.append (null, "");

		name_entry = new Gtk.Entry ();
		name_entry.show ();

		attach (collections_combo, 0, 0, 1, 1);
		attach (folder_combo, 1, 0, 1, 1);
		attach (name_entry, 0, 1, 2, 1);

		collections_combo.changed.connect (on_collection_selected);
		load_folder ();

		collections_combo.changed.connect (()=>{changed ();});
		folder_combo.changed.connect (()=>{changed ();});
		name_entry.changed.connect (()=>{changed ();});
	}

	public int64? get_selected_collection_id () {
		var string_id = collections_combo.get_active_id ();
		if (string_id == null)
			return null;

		return int64.parse (string_id);
	}

	public int64? get_selected_item_id () {
		var string_id = folder_combo.get_active_id ();
		if (string_id == null)
			return null;

		return int64.parse (string_id);
	}

	public string get_name () {
		return name_entry.text;
	}

	private void on_collection_selected () {
		folder_combo.remove_all ();
	}

	private void load_folder ()
		requires (Config.settings.http_collection_id != 0) {

		var collection = Config.get_selected_http_collection ();
		append_items_name_with_parent (null, collection.items);
	}

	private void append_items_name_with_parent (string? parent, HttpItem[] items) {
		foreach(var item in items) {
			if (!item.is_folder)
				continue;

			var name = item.name;
			if (parent != null)
				name = @"$parent -> $name";

			folder_combo.append (item.id.to_string (), name);
			append_items_name_with_parent (name, item.items);
		}
	}
}

