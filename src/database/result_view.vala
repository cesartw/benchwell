public class Benchwell.Database.ResultView : Gtk.Paned {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.DatabaseService service { get; construct; }
	public Benchwell.SourceView editor;
	public Benchwell.Database.Table table;

	public signal void exec_query (string query);
	public signal void fav_saved ();

	public ResultView (Benchwell.ApplicationWindow window, Benchwell.DatabaseService service) {
		Object (
			window: window,
			service: service,
			orientation: Gtk.Orientation.VERTICAL
		);

		table = new Benchwell.Database.Table (window, service);
		table.show ();

		// editor
		editor = new Benchwell.SourceView ("sql");
		editor.show ();

		editor.statement_selector = new SqlStatementSelector ();

		var editor_sw = new Gtk.ScrolledWindow (null, null);
		editor_sw.show ();
		editor_sw.add (editor);

		var table_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		table_box.vexpand = true;
		table_box.hexpand = true;
		table_box.pack_start (table, true, true);
		table_box.show ();

		pack1 (editor_sw, false, false);
		pack2 (table_box, true, true);

		editor.key_press_event.connect (on_editor_key_press);
		table.file_opened.connect ((query) => {
			editor.get_buffer ().set_text (query);
		});

		table.file_saved.connect ((filename) => {
			var buffer = editor.get_buffer ();
			Gtk.TextIter start, end;
			buffer.get_start_iter (out start);
			buffer.get_end_iter (out end);
			var txt = buffer.get_text (start, end, false);

			var fs = FileStream.open (filename, "w");
			fs.puts (txt);
		});

		table.fav_saved.connect( (query_name) => {
			var buffer = editor.get_buffer ();
			Gtk.TextIter start, end;
			buffer.get_start_iter (out start);
			buffer.get_end_iter (out end);

			var query_text = buffer.get_text (start, end, false);
			try {
				var query = service.info.add_query ();

				query.touch_without_save (() => {
					query.query = query_text;
					query.query_type = "fav";
					query.name = query_name;
				});
				query.save ();
			} catch (Benchwell.ConfigError err) {
				Config.show_alert (this, err.message);
			}
		});

		editor.completion.add_provider (new Benchwell.SQL.TableCompletion (service));
	}

	private bool on_editor_key_press (Gtk.Widget widget, Gdk.EventKey event) {
		if (event.keyval != Gdk.Key.Return) {
			return false;
		}

		if (event.state != Gdk.ModifierType.CONTROL_MASK) {
			return false;
		}

		_exec_query ();
		return true;
	}

	public void _exec_query () {
		var buffer = editor.get_buffer ();
		Gtk.TextIter start, end;

		var query = buffer.text;
		if (buffer.get_selection_bounds (out start, out end)) {
			query = buffer.get_text (start, end, false);
		}

		exec_query (query);
	}
}

public class Benchwell.Database.SqlStatementSelector : Benchwell.SourceViewStatementSelector, Object {
	public bool select (Gtk.SourceBuffer buffer) {
		Gtk.TextIter backward_iter;
		buffer.get_iter_at_offset (out backward_iter, buffer.cursor_position);

		Gtk.TextIter forward_iter;
		buffer.get_iter_at_offset (out forward_iter, buffer.cursor_position);

		do {
			Gtk.TextIter end;
			backward_iter.set_line_offset (0);

			buffer.get_iter_at_line (out end, backward_iter.get_line ());
			if (!end.ends_line ())
				end.forward_to_line_end ();
			var txt = buffer.get_text (backward_iter, end, false);
			if (txt.strip () == "" || txt == "\n" || txt == null) {
				break;
			}
		} while (backward_iter.backward_line());

		if (backward_iter.get_line () > 0) {
			backward_iter.forward_line ();
		}

		do {
			forward_iter.set_line_offset (0);
			Gtk.TextIter end;
			buffer.get_iter_at_line (out end, forward_iter.get_line ());
			end.forward_to_line_end ();
			var txt = buffer.get_text (forward_iter, end, false);
			if (txt == "" || txt == "\n" || txt == null) {
				break;
			}
		} while (forward_iter.forward_line());

		buffer.select_range (backward_iter, forward_iter);


		return true;
	}
}
