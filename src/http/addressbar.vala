public class Benchwell.Http.HttpAddressBar : Gtk.Box {
	public Gtk.ComboBoxText method_combo;
	public Gtk.Entry address;
	public Gtk.Button send_btn;
	public Benchwell.OptButton save_btn;

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

		send_btn = new Gtk.Button.with_label (_("SEND"));
		send_btn.get_style_context ().add_class ("suggested-action");
		send_btn.show ();

		// TODO: add to window
		var save_as_action = new GLib.SimpleAction ("win.saveas", null);
		save_btn = new Benchwell.OptButton(_("SAVE"), _("Save as"), "win.saveas");
		save_btn.show ();

		pack_start(method_combo, false, false, 0);
		pack_start(address, true, true, 0);
		pack_end(save_btn, false, false, 0);
		pack_end(send_btn, false, false, 0);
	}

	public void set_request (Benchwell.HttpItem item) {
		address.set_text (item.url);
		address.tooltip_text = Config.environments.nth_data (0).interpolate (item.url);
		method_combo.set_active_id (item.method);
	}
}

