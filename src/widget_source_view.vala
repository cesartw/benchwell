public class Benchwell.SourceView : Gtk.SourceView {
	public SourceView (string lang = "auto") {
		Object (
			show_right_margin: false,
			hexpand: true,
			vexpand: true,
			auto_indent: true,
			accepts_tab: true,
			background_pattern: Gtk.SourceBackgroundPatternType.GRID,
			highlight_current_line: Config.settings.editor_highlight_line,
			tab_width: (uint)Config.settings.editor_tab_width,
			show_line_numbers: Config.settings.editor_line_number,
			insert_spaces_instead_of_tabs: Config.settings.editor_no_tabs
		);

		set_language (lang);

		// PRETTY
		//get_space_drawer ().set_types_for_locations (Gtk.SourceSpaceLocationFlags.LEADING|Gtk.SourceSpaceLocationFlags.TRAILING, Gtk.SourceSpaceTypeFlags.ALL);
		//get_space_drawer ().enable_matrix = true;


		var buffer = (Gtk.SourceBuffer) get_buffer ();
		var sm = Gtk.SourceStyleSchemeManager.get_default ();

		if (Config.settings.editor_theme in sm.scheme_ids) {
			buffer.set_style_scheme (sm.get_scheme (Config.settings.editor_theme));
		}

		if (Config.settings.editor_font != "") {
			Config.set_font (this, Pango.FontDescription.from_string (Config.settings.editor_font));
		}

		Config.settings.changed["editor-theme"].connect (() => {
			if (Config.settings.editor_theme in sm.scheme_ids) {
				buffer.set_style_scheme (sm.get_scheme (Config.settings.editor_theme));
			}
		});

		Config.settings.changed["editor-font"].connect (() => {
			Config.set_font (this, Pango.FontDescription.from_string (Config.settings.editor_font));
		});

		Config.settings.changed["editor-tab-width"].connect (() => {
			tab_width = (uint)Config.settings.editor_tab_width;
		});

		Config.settings.changed["editor-line-number"].connect (() => {
			show_line_numbers = Config.settings.editor_line_number;
		});

		Config.settings.changed["editor-highlight-line"].connect (() => {
			highlight_current_line = Config.settings.editor_highlight_line;
		});

		Config.settings.changed["editor-no-tabs"].connect (() => {
			insert_spaces_instead_of_tabs = Config.settings.editor_no_tabs;
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
