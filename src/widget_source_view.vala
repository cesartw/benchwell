public class Benchwell.SourceView : Gtk.SourceView {
	public SourceView (string lang = "auto") {
		Object (
			show_line_numbers: false,
			show_right_margin: true,
			hexpand: true,
			vexpand: true,
			auto_indent: true,
			accepts_tab: true,
			highlight_current_line: false,
			tab_width: (uint)Config.settings.get_int64("editor-tab-width")
		);

		set_language (lang);

		var buffer = (Gtk.SourceBuffer) get_buffer ();
		var sm = Gtk.SourceStyleSchemeManager.get_default ();

		if (Config.settings.get_string("editor-theme") in sm.scheme_ids) {
			buffer.set_style_scheme (sm.get_scheme (Config.settings.get_string("editor-theme")));
		}

		if (Config.settings.get_string("editor-font") != "") {
			override_font (Pango.FontDescription.from_string (Config.settings.get_string("editor-font")));
		}

		Config.settings.changed["editor-theme"].connect (() => {
			if (Config.settings.get_string("editor-theme") in sm.scheme_ids) {
				buffer.set_style_scheme (sm.get_scheme (Config.settings.get_string("editor-theme")));
			}
		});

		Config.settings.changed["editor-font"].connect (() => {
			override_font (Pango.FontDescription.from_string (Config.settings.get_string("editor-font")));
		});
		Config.settings.changed["editor-tab-width"].connect (() => {
			tab_width = (uint)Config.settings.get_int64("editor-tab-width");
		});
	}

	public void set_language (string? lang) {
		var buffer = (Gtk.SourceBuffer) get_buffer ();
		if (lang == null || lang == "") {
			buffer.set_language (null);
			return;
		}

		var lm = Gtk.SourceLanguageManager.get_default ();
		buffer.set_language (lm.get_language (lang));
	}

	public void set_language_by_mime_type (string mime_type) {
		var mime = mime_type.strip ();
		switch (mime) {
			case "application/json", "application/html", "application/xml", "application/yaml":
				var buffer = (Gtk.SourceBuffer) get_buffer ();
				set_language (mime.split("/")[1]);
				buffer.highlight_syntax = true;
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

	public string get_text () {
		Gtk.TextIter start, end;
		var buff = get_buffer ();
		buff.get_start_iter (out start);
		buff.get_end_iter (out end);
		return buff.get_text (start, end, false);
	}
}
