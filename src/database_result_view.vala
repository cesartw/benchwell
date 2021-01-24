public class Benchwell.Database.ResultView : Gtk.Paned {
	public Benchwell.ApplicationWindow window { get; construct; }
	public Benchwell.DatabaseService service { get; construct; }
	public Benchwell.SourceView editor;
	public Benchwell.Database.Table table;
	public Gtk.Button btn_load_query;
	public Gtk.MenuButton save_menu;
	public Gtk.InfoBar infobar;
	public Gtk.Label infobar_label;

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
		editor.show_line_numbers = true;
		editor.show_line_marks = true;
		editor.margin_top = 10;
		editor.show ();

		var editor_sw = new Gtk.ScrolledWindow (null, null);
		editor_sw.show ();
		editor_sw.add (editor);

		// table controls
		btn_load_query = new Benchwell.Button ("open", Gtk.IconSize.BUTTON);
		btn_load_query.show ();

		var img = new Benchwell.Image("save");
		save_menu = new Gtk.MenuButton ();
		save_menu.show ();
		save_menu.set_image (img);

		var save_menu_model = new GLib.Menu ();
		save_menu_model.append (_("Save As"), "win.save.file");
		save_menu_model.append (_("Save fav"), "win.save.fav");

		save_menu.set_menu_model (save_menu_model);
		var action_save_file = new GLib.SimpleAction ("save.file", null);
		var action_save_fav = new GLib.SimpleAction ("save.fav", null);
		window.add_action (action_save_file);
		window.add_action (action_save_fav);

		infobar = new Gtk.InfoBar ();
		var infobar_label_sw = new Gtk.ScrolledWindow (null, null);
		infobar_label_sw.show ();
		//infobar_label_sw.vexpand = true;
		infobar_label_sw.hexpand = true;
		infobar.no_show_all = true;
		infobar.add_button (_("Ok"), Gtk.ResponseType.OK);

		infobar_label = new Gtk.Label("");
		infobar_label.wrap = true;
		infobar_label.show ();
		infobar_label_sw.add (infobar_label);
		infobar.get_content_area ().add (infobar_label_sw);

		var editor_actionbar = new Gtk.ActionBar ();
		editor_actionbar.show ();
		editor_actionbar.pack_end (save_menu);
		editor_actionbar.pack_end (btn_load_query);
		editor_actionbar.set_name ("queryactionbar");
		/////////////////

		var editor_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		editor_box.pack_start (editor_sw, true, true, 0);
		editor_box.pack_end (editor_actionbar, false, false, 0);
		editor_box.show ();

		var table_box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		table_box.vexpand = true;
		table_box.hexpand = true;
		table_box.pack_start (table, true, true);
		table_box.pack_end (infobar, false, false);
		table_box.show ();

		pack1 (editor_box, false, false);
		pack2 (table_box, true, true);

		action_save_file.activate.connect (on_save_file);
		action_save_fav.activate.connect (on_save_fav);
		btn_load_query.clicked.connect (on_open_file);

		editor.key_press_event.connect (on_editor_key_press);
		infobar.response.connect (infobar.hide);
	}

	public void on_open_file () {
		var dialog = new Gtk.FileChooserDialog (_("Select file"), window,
											 Gtk.FileChooserAction.OPEN,
											_("Open"), Gtk.ResponseType.OK,
											_("Cancel"), Gtk.ResponseType.CANCEL);
		var resp = (Gtk.ResponseType) dialog.run ();

		if (resp == Gtk.ResponseType.CANCEL) {
			dialog.destroy ();
			return;
		}

		var filename = dialog.get_filename ();
		dialog.destroy ();

		string text;
		var ok = GLib.FileUtils.get_contents (filename, out text, null);
		if (!ok) {
			return;
		}

		editor.get_buffer ().set_text (text);
	}

	public void on_save_file () {
		var dialog = new Gtk.FileChooserDialog (_("Save file"), window,
											 Gtk.FileChooserAction.SAVE,
											_("Ok"), Gtk.ResponseType.OK,
											_("Cancel"), Gtk.ResponseType.CANCEL);
		var resp = (Gtk.ResponseType) dialog.run ();

		if (resp == Gtk.ResponseType.CANCEL) {
			dialog.destroy ();
			return;
		}

		var filename = dialog.get_filename ();
		dialog.destroy ();

		var buffer = editor.get_buffer ();
		Gtk.TextIter start, end;
		buffer.get_start_iter (out start);
		buffer.get_end_iter (out end);
		var txt = buffer.get_text (start, end, false);

		var fs = FileStream.open (filename, "w");
		fs.puts (txt);
	}

	public void on_save_fav () {
		var query_name = ask_fav_name ();

		if (query_name == "") {
			return;
		}


		var buffer = editor.get_buffer ();
		Gtk.TextIter start, end;
		buffer.get_start_iter (out start);
		buffer.get_end_iter (out end);

		var query_text = buffer.get_text (start, end, false);
		var query = service.info.add_query ();

		query.touch_without_save (() => {
			query.query = query_text;
			query.query_type = "fav";
			query.name = query_name;
		});
		query.save ();
		fav_saved ();
	}

	public void show_alert (string message, Gtk.MessageType type = Gtk.MessageType.ERROR, bool autohide = false, int timeout = 0) {
		infobar_label.set_text (message);
		infobar.message_type = type;
		infobar.show ();

		if (autohide) {
			if (timeout == 0) {
				switch (message.split (" ").length) {
					case 1:
						timeout = 1000;
						break;
					case 2, 3:
						timeout = 1500;
						break;
					default:
						timeout = 3000;
						break;
				}
			}

			Timeout.add (timeout, () => {
				infobar.hide ();
				return false;
			});
		}
	}

	public void hide_alert () {
		infobar.hide ();
	}

	private string ask_fav_name () {
		var dialog = new Gtk.Dialog.with_buttons (_("Choose"), window,
									Gtk.DialogFlags.DESTROY_WITH_PARENT|Gtk.DialogFlags.MODAL,
									_("Ok"), Gtk.ResponseType.OK,
									_("Cancel"), Gtk.ResponseType.CANCEL);
		dialog.set_default_size (250, 130);

		var label = new Gtk.Label (_("Enter favorite name"));
		label.show ();

		var entry = new Gtk.Entry ();
		entry.show ();

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 10);
		box.show ();

		box.pack_start (label, true, true, 0);
		box.pack_start (entry, true, true, 0);

		dialog.get_content_area ().add (box);

		var resp = (Gtk.ResponseType) dialog.run ();
		var filename = entry.get_text ();
		dialog.destroy ();

		if (resp != Gtk.ResponseType.OK) {
			return "";
		}

		return filename;
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
		var query = editor.get_buffer ().text;
		exec_query (query);
	}
}

