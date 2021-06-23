public class Benchwell.SecretEntry : Gtk.Entry {
	private bool _open;
	public bool open {
		get { return _open; }
		set {
			_open = value;
			if ( _open ) {
				placeholder_text = "";
			} else {
				set_text ("");
				placeholder_text = _("Stored");
			}
		}
	}

	public SecretEntry (bool open = false) {
		Object(
			caps_lock_warning: true,
			placeholder_text: open ? "" : _("Stored"),
			input_purpose: Gtk.InputPurpose.PASSWORD,
			visibility: false
		);
		this.open = open;
	}
}
