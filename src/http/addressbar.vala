public class Benchwell.Http.AddressBar : Gtk.Box {
	public Gtk.ComboBoxText method_combo;
	public Gtk.Entry address;
	public Gtk.Label address_label;
	public Benchwell.OptButton send_btn;
	public Benchwell.HttpItem? item;

	public signal void changed ();

	internal bool disable_changed = false;

	public AddressBar () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5,
			name: "HttpAddressBar"
		);

		method_combo = new Gtk.ComboBoxText ();
		foreach (Benchwell.Http.Method method in Benchwell.Http.Method.all())
			method_combo.append (method.to_string (), method.to_string ());
		method_combo.set_active (0);
		method_combo.vexpand = false;
		method_combo.set_valign (Gtk.Align.START);
		method_combo.show ();

		address = new Gtk.Entry ();
		address.placeholder_text = "http://localhost/path.json";
		address.show ();

		address_label = new Gtk.Label ("");
		address_label.set_halign (Gtk.Align.START);
		address_label.ellipsize = Pango.EllipsizeMode.END;
		address_label.name = "address-tag-line";
		address_label.sensitive = false;
		address_label.show ();

		var address_label_eventbox = new Gtk.EventBox ();
		address_label_eventbox.add (address_label);
		address_label_eventbox.show ();

		// TODO: add to window
		send_btn = new Benchwell.OptButton(_("SEND"),
										   _("Save as"), "win.saveas",
										   _("Copy curl"), "win.copycurl");
		send_btn.btn.get_style_context ().add_class ("suggested-action");
		send_btn.menu_btn.get_style_context ().add_class ("suggested-action");
		send_btn.sensitive = false;
		send_btn.show ();

		var address_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		address_box.pack_start (address, false, false, 0);
		address_box.pack_start (address_label_eventbox, false, false, 0);
		address_box.show ();

		pack_start(method_combo, false, false, 0);
		pack_start(address_box, true, true, 0);
		pack_end(send_btn, false, false, 0);

		Config.environments.selected_changed.connect (on_environment_selected_change);
		method_combo.changed.connect (on_method_change);
		address.changed.connect (on_address_change);

		address_label_eventbox.button_press_event.connect (on_copy_url);
	}

	private bool on_copy_url (Gtk.Widget w, Gdk.EventButton event) {
		if (event.button != Gdk.BUTTON_PRIMARY) {
			return false;
		}

		var cb = Gtk.Clipboard.get_default (Gdk.Display.get_default ());
		var st = address_label.get_text ();
		cb.set_text (st, st.length);

		return true;
	}

	private void on_address_change () {
		if (item != null) {
			item.url = address.text;
			changed();
		}

		send_btn.sensitive = address.text.strip () != "";
	}

	private void on_method_change () {
		if (disable_changed || item == null) { return; }

		item.method = method_combo.get_active_text ();

		changed();
	}

	private void on_environment_selected_change () {
		if (Config.environments.selected == null) {
			return;
		}
		address.tooltip_text = Config.environments.selected.interpolate (address.text);
		address_label.set_text (address.tooltip_text);
	}

	public void set_request (Benchwell.HttpItem item) {
		disable_changed = true;
		this.item = item;

		address.set_text (item.url);
		address_label.set_text (item.url);
		method_combo.set_active_id (item.method);
		disable_changed = false;
	}
}
