public class Benchwell.Http.HttpAddressBar : Gtk.Box {
	public Gtk.ComboBoxText method_combo;
	public Gtk.Entry address;
	public Gtk.Label address_label;
	public Benchwell.OptButton send_btn;
	public Benchwell.HttpItem? item;

	public signal void changed ();

	internal bool disable_changed = false;

	enum Method {
		GET,
		POST,
		PUT,
		PATCH,
		DELETE,
		OPTIONS,
		HEAD;

		public string to_string () {
			switch (this) {
				case GET:
					return "GET";
				case POST:
					return  "POST";
				case PUT:
					return "PUT";
				case PATCH:
					return "PATCH";
				case DELETE:
					return "DELETE";
				case OPTIONS:
					return "OPTIONS";
				case HEAD:
					return "HEAD";
				default:
					assert_not_reached ();
			}
		}

		public static Method[] all () {
			return {
				GET,
				POST,
				PUT,
				PATCH,
				DELETE,
				OPTIONS,
				HEAD
			};
		}
	}

	public HttpAddressBar () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5,
			name: "HttpAddressBar"
		);

		method_combo = new Gtk.ComboBoxText ();
		foreach (Method method in Method.all())
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
		send_btn = new Benchwell.OptButton(_("SEND"), _("Save as"), "win.saveas");
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

		Config.environment_changed.connect (() => {
			if (Config.environment != null) {
				address.tooltip_text = Config.environment.interpolate (address.text);
				address_label.set_text (address.tooltip_text);
			}
		});

		method_combo.changed.connect (() => {
			if (disable_changed || item == null) { return; }

			item.method = method_combo.get_active_text ();

			changed();
		});

		address.changed.connect (() => {
			if (item != null) {
				item.url = address.text;
				changed();
			}

			send_btn.sensitive = address.text.strip () != "";
		});

		address_label_eventbox.button_press_event.connect ((w, event) => {
			if (event.button != Gdk.BUTTON_PRIMARY) {
				return false;
			}

			var cb = Gtk.Clipboard.get_default (Gdk.Display.get_default ());
			var st = address_label.get_text ();
			cb.set_text (st, st.length);

			return true;
		});
	}

	public void set_request (Benchwell.HttpItem item) {
		disable_changed = true;
		this.item = item;

		address.set_text (item.url);
		address_label.set_text (item.url);
		method_combo.set_active_id (item.method);
		update_url ();
		disable_changed = false;
	}

	public void update_url () {
		// TODO: add tagline with interpolated url
		//if ( item.query_params.length ==  0 ){
			//return;
		//}

		//if (item.url.index_of ("?", 0) == -1)
			//address.text = "%s?".printf (item.url);

		//for (var i = 0; i < item.query_params.length; i++) {
			//if ( item.query_params[i].key == "") {
				//continue;
			//}
			//if (i > 0) {
				//address.text += "&";
			//}

			//address.text += "%s=%s".printf (item.query_params[i].key, item.query_params[i].val);
		//}
	}
}

