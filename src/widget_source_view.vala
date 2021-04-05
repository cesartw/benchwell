[Compact]
private struct line_desc
{
	uint ilevel;
	uint start_line;
	uint end_line;
}

class SimpleStack<T> {
	private T[] stack;

	public bool empty () {
		return stack == null || stack.length == 0;
	}

	public T peek () {
		assert_false (empty ());
		return stack[stack.length - 1];
	}

	public void push (T v) {
		if (stack == null)
			stack = {};
		stack += v;
	}

	public T pop () {
		T v = stack[stack.length - 1];
		stack = stack[0:stack.length -1];

		return v;
	}
}

public class Benchwell.SourceView : Gtk.SourceView {
	private uint indent_timeout = 0;

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
			show_line_marks: true,
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

		build_markers ();

        buffer.create_tag ("foldable", "invisible", true);

		buffer.changed.connect (on_buffer_changed);
		line_mark_activated.connect (on_mark_activated);
	}

	private void on_mark_activated (Gtk.TextIter iter) {
		var buffer = get_buffer () as Gtk.SourceBuffer;
		var marks = buffer.get_source_marks_at_iter (iter, null);

		foreach (var mark in marks) {
			// NOTE: I use the mark.name to transfer the line span {start-line}-{end-line}. Hacky asf
			var mark_name = mark.name;
			var start_end = mark_name.split ("-");
			Gtk.TextIter start, end, next_to_start;

			int start_line = int.parse (start_end[0]);
			int end_line = int.parse (start_end[1]);

			buffer.get_iter_at_line (out start, start_line);
			buffer.get_iter_at_line (out end, end_line);
			buffer.get_iter_at_line (out next_to_start, start_line + 1);

			switch (mark.category) {
				case "fold_collapse":
					buffer.apply_tag_by_name ("foldable", next_to_start, end);

					buffer.remove_source_marks (start, start, null);
					buffer.create_source_mark (@"$(mark_name)", "fold_expand", start);
					buffer.create_source_mark (@"$(mark_name)-more", "fold_more", end);

					break;
				case "fold_expand", "fold_more":
					buffer.remove_tag_by_name ("foldable", start , end);

					buffer.remove_source_marks (start, start, null);
					buffer.remove_source_marks (end, end, null);
					buffer.create_source_mark (@"$(mark_name)", "fold_collapse", start);

					break;
				default:
					return;
			}
		}
	}

	private void on_buffer_changed () {
		add_folding_marks ();
	}

	//    0 {
	// +  1     "inbox": {
	//    2         "name": "nightowlstud.io",
	//    3         "email": "support@nightowlstud.io",
	// +  4         "users": [
	// +  5             {
	//    6                 "id": 7,
	// +  7                 "meta": {
	//    8                     "isAdmin": true
	//    9                 }
	//   10             },
	// + 11             {
	//   12                 "id": 1,
	// + 13                 "meta": {
	//   14                     "isAdmin": true
	//   15                 }
	//   16             }
	//   17         ]
	//   18     }
	//   19 }
	private void add_folding_marks () {
		if (indent_timeout > 0) {
			Source.remove (indent_timeout);
		}

		var buffer = get_buffer () as Gtk.SourceBuffer;

		indent_timeout = Timeout.add (100, () => {
			Gtk.TextIter start, end;
			buffer.get_start_iter (out start);
			buffer.get_end_iter (out end);
			buffer.remove_source_marks (start, end, null);

			indent_timeout = 0;
			var pointer_stack = new SimpleStack<int> ();
			uint line_number = -1;
			string[] lines = buffer.text.split ("\n");
			line_desc[] lines_meta = {};

			foreach (var line in lines) {
				line_number++;
				var ilevel = indent_level (line);
				if (ilevel == 0) {
					continue; // won't add mark at 0
				}

				// there's a next line and it's ilevel is greater
				if (lines.length >= line_number && indent_level (lines[line_number+1]) > ilevel) {
					lines_meta += line_desc () {
						start_line = line_number,
						ilevel = ilevel
					};
					pointer_stack.push (lines_meta.length - 1);
					continue;
				}

				if (!pointer_stack.empty () && ilevel <= lines_meta[pointer].ilevel) {
					var pointer = pointer_stack.pop ();
					lines_meta[pointer].end_line = line_number;
					continue;
				}
			}

			foreach (var meta in lines_meta) {
				Gtk.TextIter iter;
                buffer.get_iter_at_line (out iter, (int)meta.start_line);
				buffer.create_source_mark (@"$(meta.start_line)-$(meta.end_line)", "fold_collapse", iter);
			}

			return Source.REMOVE;
		});
	}

	private void build_markers () {
		//var px = Gtk.IconTheme.get_default ().load_icon ("bw-directory", Gtk.IconSize.BUTTON, Gtk.IconLookupFlags.USE_BUILTIN);
		var mark_attr = new Gtk.SourceMarkAttributes ();
		//mark_attr.set_pixbuf (px);
		mark_attr.set_icon_name ("pan-down-symbolic");
		set_mark_attributes ("fold_collapse", mark_attr, 0);

		mark_attr = new Gtk.SourceMarkAttributes ();
		mark_attr.set_icon_name ("pan-end-symbolic");
		set_mark_attributes ("fold_expand", mark_attr, 0);

		mark_attr = new Gtk.SourceMarkAttributes ();
		mark_attr.set_icon_name ("view-more-symbolic");
		set_mark_attributes ("fold_more", mark_attr, 0);
	}

	private void dump_line_descriptions (owned line_desc[] meta) {
		foreach (var m in meta) {
			dump_line_description(m);
		}
	}

	private void dump_line_description (owned line_desc m) {
		print (@"===$(m.start_line):$(m.end_line)\n");
	}

	private uint indent_level (owned string line) {
		uint count = 0;

		foreach (var c in line.to_utf8 ()) {
			switch (c) {
				case ' ':
					count++;
					break;
				case '\t':
					count += tab_width;
					break;
				default:
					return count;
			}
		}

		return count;
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
