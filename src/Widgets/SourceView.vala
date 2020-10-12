public class Benchwell.SourceView : Gtk.SourceView {
	public SourceView (string lang = "auto", string style = "oblivion") {
		Object (
			show_line_numbers: false,
			show_right_margin: true,
			hexpand: true,
			vexpand: true
		);

		set_language (lang);

		var buffer = (Gtk.SourceBuffer) get_buffer ();
		var sm = Gtk.SourceStyleSchemeManager.get_default ();
		buffer.set_style_scheme (sm.get_scheme (style));
	}

	public void set_language (string? lang) {
		var buffer = (Gtk.SourceBuffer) get_buffer ();
		if (lang == null) {
			buffer.set_language (null);
		}

		var lm = Gtk.SourceLanguageManager.get_default ();
		buffer.set_language (lm.get_language (lang));
	}

	public void set_language_by_mime_type (string mime_type) {
		switch (mime_type) {
			case "application/json":
				set_language ("json");
				break;
			case "application/html":
				set_language ("html");
				break;
			case "application/xml":
				set_language ("xml");
				break;
			case "application/yaml":
				set_language ("yaml");
				break;
			case "auto":
				var buffer = (Gtk.SourceBuffer) get_buffer ();
				buffer.highlight_syntax = true;
				set_language (null);
				break;
			default:
				var buffer = (Gtk.SourceBuffer) get_buffer ();
				buffer.highlight_syntax = false;
				set_language (null);
				break;
		}
	}
}
