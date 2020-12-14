public class Benchwell.Http.HttpAddressBar : Gtk.Box {
	public Gtk.ComboBoxText method_combo;
	public Gtk.Entry address;
	public Benchwell.OptButton send_btn;
	public Benchwell.HttpItem? item;

	public signal void changed ();

	internal bool disable_changed = false;

	public HttpAddressBar () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		method_combo = new Gtk.ComboBoxText ();
		for (var i = 0; i < Benchwell.Methods.length; i++) {
			method_combo.append (Benchwell.Methods[i], Benchwell.Methods[i]);
		}
		method_combo.set_active (0);
		method_combo.show ();

		address = new Gtk.Entry ();
		address.placeholder_text = "http://localhost/path.json";
		address.show ();

		// TODO: add to window
		send_btn = new Benchwell.OptButton(_("SEND"), _("Save as"), "win.saveas");
		send_btn.btn.get_style_context ().add_class ("suggested-action");
		send_btn.menu_btn.get_style_context ().add_class ("suggested-action");
		send_btn.show ();

		pack_start(method_combo, false, false, 0);
		pack_start(address, true, true, 0);
		pack_end(send_btn, false, false, 0);

		Config.environment_changed.connect (() => {
			if (item != null && Config.environment != null) {
				address.tooltip_text = Config.environment.interpolate (item.url);
			}
		});

		method_combo.changed.connect (() => {
			if (disable_changed || item == null) { return; }

			item.method = method_combo.get_active_text ();
			changed();
		});

		address.changed.connect (() => {
			if (disable_changed || item == null) { return; }

			item.url = address.text;
			changed();
		});
	}

	public void set_request (Benchwell.HttpItem item) {
		disable_changed = true;
		this.item = item;
		address.set_text (item.url);
		if (Config.environment != null) {
			address.tooltip_text = Config.environment.interpolate (item.url);
		}
		method_combo.set_active_id (item.method);
		disable_changed = false;
	}
}

