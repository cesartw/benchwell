public class Benchwell.SettingsPanel : Gtk.Box {
	private Gtk.Notebook notebook;
	private Benchwell.EnvironmentEditor env_editor;
	private Benchwell.EditorSettings editor_settings;
	private Benchwell.HttpSettings http_settings;
	private Benchwell.PomodoroSettings pomodoro_settings;
	private Gtk.Switch dark_switch;

	public SettingsPanel () {
		Object (
			orientation: Gtk.Orientation.VERTICAL,
			vexpand: true,
			hexpand: true
		);

		var header_bar = new Gtk.HeaderBar ();
		header_bar.title = _("Settings");
		header_bar.show ();

		dark_switch = new Gtk.Switch ();
		dark_switch.show();
		var dark_icon = new Gtk.Image.from_icon_name ("weather-clear-night", Gtk.IconSize.SMALL_TOOLBAR);
		dark_icon.show ();


		notebook = new Gtk.Notebook ();
		notebook.tab_pos = Gtk.PositionType.TOP;
		notebook.show ();

		env_editor = new Benchwell.EnvironmentEditor ();
		env_editor.show ();

		editor_settings = new Benchwell.EditorSettings ();
		editor_settings.show ();

		http_settings = new Benchwell.HttpSettings ();
		http_settings.show ();

		pomodoro_settings = new Benchwell.PomodoroSettings ();
		pomodoro_settings.show ();

		notebook.append_page (env_editor, new Gtk.Label (_("Environments")));
		notebook.append_page (editor_settings, new Gtk.Label (_("Editor")));
		notebook.append_page (http_settings, new Gtk.Label (_("HTTP")));
		notebook.append_page (pomodoro_settings, new Gtk.Label (_("Pomodoro")));

		header_bar.pack_end (dark_switch);
		header_bar.pack_end (dark_icon);
		pack_start (header_bar, false, false, 0);
		pack_start (notebook, true, true, 0);

		dark_switch.state = Config.settings.get_boolean ("dark-mode");
		dark_switch.state_set.connect ((state) => {
			Config.settings.set_boolean ("dark-mode", state);
			Gtk.Settings.get_default ().gtk_application_prefer_dark_theme = state;
			return false;
		});
	}

	public void select_env (Benchwell.Environment env) {
		env_editor.select_env (env);
	}
}

public class Benchwell.EditorSettings : Gtk.Grid {
	public EditorSettings () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);

		var label_alignment = Gtk.Align.START;

		// LOOK & FEEL
		var laf_label = new Gtk.Label (_("Look & Feel")) {
			valign = label_alignment,
			halign = label_alignment
		};
		laf_label.show ();

		// EDITOR THEME
		var laf_theme_combo = new Gtk.ComboBoxText ();
		foreach (var id in Gtk.SourceStyleSchemeManager.get_default ().scheme_ids) {
			laf_theme_combo.append (id, id);
		}
		laf_theme_combo.set_active_id (Config.settings.get_string ("editor-theme"));
		laf_theme_combo.show ();

		var laf_theme_label = new Gtk.Label (_("Theme")) {
			valign = label_alignment,
			halign = label_alignment
		};
		laf_theme_label.show ();

		///////////////

		// FONT
		var laf_font_label = new Gtk.Label (_("Font")) {
			valign = label_alignment,
			halign = label_alignment
		};
		laf_font_label.show ();

		var laf_font_btn = new Gtk.FontButton ();
		if (Config.settings.get_string ("editor-font") != "") {
			laf_font_btn.font = Config.settings.get_string ("editor-font");
		}
		laf_font_btn.show ();

		///////

		// EDITING
		var editing_label = new Gtk.Label (_("Editing")) {
			valign = label_alignment,
			halign = label_alignment
		};
		editing_label.show ();

		// TAB WIDTH
		var editing_tabwidth_label = new Gtk.Label (_("Tab width")) {
			valign = label_alignment,
			halign = label_alignment
		};
		editing_tabwidth_label.show ();

		var editing_tabwidth_spin = new Gtk.SpinButton.with_range (2, 8, 2);
		editing_tabwidth_spin.value = (double) Config.settings.get_int64 ("editor-tab-width");
		editing_tabwidth_spin.show ();
		////////////

		// SHOW LINE NUMBER
		var editing_ln_label = new Gtk.Label (_("Show line number")) {
			valign = label_alignment,
			halign = label_alignment
		};
		editing_ln_label.show ();

		var editing_ln_sw = new Gtk.Switch ();
		editing_ln_sw.state = Config.settings.get_boolean ("editor-line-number");
		editing_ln_sw.show ();
		///////////////////

		// HIGHLIGHT CURRENT LINE
		var editing_hl_label = new Gtk.Label (_("Highlight current line")) {
			valign = label_alignment,
			halign = label_alignment
		};
		editing_hl_label.show ();

		var editing_hl_sw = new Gtk.Switch ();
		editing_hl_sw.state = Config.settings.get_boolean ("editor-highlight-line");
		editing_hl_sw.show ();
		/////////////////////////

		// SPACES FOR TABS
		var editing_notabs_label = new Gtk.Label (_("Insert spaces instead of tabs")) {
			valign = label_alignment,
			halign = label_alignment
		};
		editing_notabs_label.show ();

		var editing_notabs_sw = new Gtk.Switch ();
		editing_notabs_sw.state = Config.settings.get_boolean ("editor-no-tabs");
		editing_notabs_sw.show ();
		//////////////////

		attach (laf_label, 0, 0, 2, 1);

		attach (laf_theme_label, 1, 1, 1, 1);
		attach (laf_theme_combo, 2, 1, 1, 1);

		attach (laf_font_label, 1, 2, 1, 1);
		attach (laf_font_btn, 2, 2, 1, 1);

		attach (editing_label, 3, 0, 2, 1);

		attach (editing_tabwidth_label, 4, 1, 1, 1);
		attach (editing_tabwidth_spin, 5, 1, 1, 1);

		attach (editing_ln_label, 4, 2, 1, 1);
		attach (editing_ln_sw, 5, 2, 1, 1);

		attach (editing_hl_label, 4, 3, 1, 1);
		attach (editing_hl_sw, 5, 3, 1, 1);

		attach (editing_notabs_label, 4, 4, 1, 1);
		attach (editing_notabs_sw, 5, 4, 1, 1);

		laf_theme_combo.changed.connect (() => {
			Config.settings.set_string ("editor-theme", laf_theme_combo.get_active_id ());
		});

		laf_font_btn.font_set.connect (() => {
			Config.settings.set_string ("editor-font", laf_font_btn.font);
		});

		editing_tabwidth_spin.changed.connect (() => {
			Config.settings.set_int64 ("editor-tab-width", (int64)editing_tabwidth_spin.value);
		});

		editing_ln_sw.state_set.connect ((state) => {
			Config.settings.set_boolean ("editor-line-number", state);
			return false;
		});

		editing_hl_sw.state_set.connect ((state) => {
			Config.settings.set_boolean ("editor-highlight-line", state);
			return false;
		});

		editing_notabs_sw.state_set.connect ((state) => {
			Config.settings.set_boolean ("editor-no-tabs", state);
			return false;
		});
	}
}

public class Benchwell.HttpSettings : Gtk.Grid {
	public HttpSettings () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);

		var label_alignment = Gtk.Align.START;

		// LOOK & FEEL
		var laf_label = new Gtk.Label (_("Look & Feel")) {
			valign = label_alignment,
			halign = label_alignment
		};
		laf_label.show ();

		// FONT
		var laf_font_label = new Gtk.Label (_("Font")) {
			valign = label_alignment,
			halign = label_alignment
		};
		laf_font_label.show ();

		var laf_font_btn = new Gtk.FontButton ();
		if (Config.settings.get_string ("http-font") != "") {
			laf_font_btn.font = Config.settings.get_string ("http-font");
		}
		laf_font_btn.show ();

		///////

		// OTHER
		var requests_label = new Gtk.Label (_("Requests")) {
			valign = label_alignment,
			halign = label_alignment
		};
		requests_label.show ();
		////////

		// ACTIVATE ON SINGLE CLICK
		var requests_single_click_label = new Gtk.Label (_("Open on single click")) {
			valign = label_alignment,
			halign = label_alignment
		};
		requests_single_click_label.show ();

		var requests_single_click_sw = new Gtk.Switch ();
		requests_single_click_sw.state = Config.settings.get_boolean ("http-single-click-activate");
		requests_single_click_sw.show ();
		//////////////////


		attach (laf_label, 0, 0, 2, 1);

		attach (laf_font_label, 1, 1, 1, 1);
		attach (laf_font_btn, 2, 1, 1, 1);

		attach (requests_label, 3, 0, 2, 1);
		attach (requests_single_click_label, 4, 1, 1, 1);
		attach (requests_single_click_sw, 5, 1, 1, 1);

		laf_font_btn.font_set.connect (() => {
			Config.settings.set_string ("http-font", laf_font_btn.font);
		});

		requests_single_click_sw.state_set.connect ((state) => {
			Config.settings.set_boolean ("http-single-click-activate", state);
			return false;
		});
	}
}

public class Benchwell.PomodoroSettings : Gtk.Grid {
	public PomodoroSettings () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);

		var label_alignment = Gtk.Align.START;

		var timings_label = new Gtk.Label (_("Timings")) {
			valign = label_alignment,
			halign = label_alignment
		};
		timings_label.show ();

		// WORK DURATION
		var duration_label = new Gtk.Label (_("Work session duration")) {
			valign = label_alignment,
			halign = label_alignment
		};
		duration_label.show ();

		var duration_spin = new Gtk.SpinButton.with_range (20, 60, 5);
		duration_spin.value = (double)Config.settings.get_int64 ("pomodoro-duration") / 60;
		duration_spin.show ();
		////////////////

		// BREAK DURATION
		var break_duration_label = new Gtk.Label (_("Short break duration")) {
			valign = label_alignment,
			halign = label_alignment
		};
		break_duration_label.show ();

		var break_duration_spin = new Gtk.SpinButton.with_range (5, 10, 1);
		break_duration_spin.value = (double)Config.settings.get_int64 ("pomodoro-break-duration") / 60;
		break_duration_spin.show ();
		/////////////////

		// LONG BREAK DURATION
		var long_break_duration_label = new Gtk.Label (_("Long break duration")) {
			valign = label_alignment,
			halign = label_alignment
		};
		long_break_duration_label.show ();

		var long_break_duration_spin = new Gtk.SpinButton.with_range (10, 20, 1);
		long_break_duration_spin.value = (double)Config.settings.get_int64 ("pomodoro-long-break-duration") / 60;
		long_break_duration_spin.show ();
		/////////////////

		// SET COUNT
		var set_count_label = new Gtk.Label (_("Sets before long break")) {
			valign = label_alignment,
			halign = label_alignment
		};
		set_count_label.show ();

		var set_count_spin = new Gtk.SpinButton.with_range (10, 20, 1);
		set_count_spin.value = (double)Config.settings.get_int64 ("pomodoro-set-count");
		set_count_spin.show ();
		/////////////////

		attach (timings_label, 0, 0, 2, 1);

		attach (duration_label, 1, 1, 1, 1);
		attach (duration_spin, 2, 1, 1, 1);

		attach (break_duration_label, 1, 2, 1, 1);
		attach (break_duration_spin, 2, 2, 1, 1);

		attach (long_break_duration_label, 1, 3, 1, 1);
		attach (long_break_duration_spin, 2, 3, 1, 1);

		attach (set_count_label, 1, 4, 1, 1);
		attach (set_count_spin, 2, 4, 1, 1);


		duration_spin.changed.connect (() => {
			Config.settings.set_int64 ("pomodoro-duration", (int64)duration_spin.value*60);
		});

		break_duration_spin.changed.connect (() => {
			Config.settings.set_int64 ("pomodoro-break-duration", (int64)break_duration_spin.value*60);
		});

		long_break_duration_spin.changed.connect (() => {
			Config.settings.set_int64 ("pomodoro-long-break-duration", (int64)long_break_duration_spin.value*60);
		});

		set_count_spin.changed.connect (() => {
			Config.settings.set_int64 ("pomodoro-set-count", (int64)set_count_spin.value);
		});
	}
}

