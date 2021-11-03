public class Benchwell.Http.CollectionsComboBox : Gtk.ComboBoxText {
	public Benchwell.ApplicationWindow    window { get; construct; }

	public CollectionsComboBox (Benchwell.ApplicationWindow window) {
		Object (
			window: window,
			name: "HttpCollectionSelect"
		);

		append ("new", _("Add collection"));
		foreach (var collection in Config.http_collections) {
			append (collection.id.to_string (), collection.name);
		}

		changed.connect (on_collection_selected);
		Config.http_collection_added.connect ((collection) => {
			append (collection.id.to_string (), collection.name);
		});
		set_active_id (Config.settings.http_collection_id.to_string ());
	}

	private void on_collection_selected () {
		if (get_active_id () == "new") {
			on_add_new_collection ();
			return;
		}

		var collection_id = int64.parse (get_active_id ());

		foreach(var collection in Config.http_collections) {
			if (collection.id != collection_id){
				continue;
			}

			// TODO: make this eager loading optional. have no use case yet
			try {
				Config.load_http_items (collection);
			} catch (Benchwell.ConfigError err) {
				Config.show_alert (this, err.message);
				return;
			}

			break;
		}
	}

	private void on_add_new_collection () {
		var dialog = new Gtk.Dialog.with_buttons (_("Add collection"), window,
									Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
									_("Save"), Gtk.ResponseType.OK,
									_("Cancel"), Gtk.ResponseType.CANCEL);
		dialog.set_default_size (250, 130);

		var label = new Gtk.Label (_("Enter collection name"));
		label.show ();

		var entry = new Gtk.Entry ();
		entry.show ();

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 10);
		box.show ();

		box.pack_start (label, true, true, 0);
		box.pack_start (entry, true, true, 0);

		dialog.get_content_area ().add (box);

		var resp = (Gtk.ResponseType) dialog.run ();
		var name = entry.get_text ();
		dialog.destroy ();

		if (resp != Gtk.ResponseType.OK) {
			return;
		}

		try {
			var collection = new Benchwell.HttpCollection ();
			collection.name = name;
			Config.add_http_collection (collection);
			//append (collection.id.to_string (), collection.name);
			set_active_id (collection.id.to_string ());
		} catch (Benchwell.ConfigError err) {
			Config.show_alert (this, err.message);
		}
	}
}
