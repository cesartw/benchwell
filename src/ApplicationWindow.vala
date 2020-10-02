public class Benchwell.ApplicationWindow : Gtk.ApplicationWindow {
	public Gtk.Notebook notebook;
	public Gtk.Statusbar statusbar;
	public Gtk.ComboBox env_combo;
	public Gtk.ListStore env_store;
	public Gtk.Button btn_env;

	public SimpleAction new_connection_menu;
	public SimpleAction new_database_tab_menu;
	public SimpleAction new_http_tab_menu;
	public SimpleAction close_menu;

	public class ApplicationWindow(Gtk.Application app) {
		Object (
			application: app
		);

		set_title ("Benchwell");

		new_database_tab_menu = new SimpleAction("new.db", null);
		new_http_tab_menu = new SimpleAction("new.http", null);
		close_menu = new SimpleAction("close", null);

		add_action(new_database_tab_menu);
		add_action(new_http_tab_menu);
		add_action(close_menu);

		notebook = new Gtk.Notebook ();
		notebook.set_name ("MainNotebook");
		notebook.set_property ("scrollable", true);
		notebook.set_group_name ("mainwindow");
		notebook.popup_enable ();

		notebook.set_property ("tab-pos", Config.tab_position());
		notebook.show ();

		statusbar = new Gtk.Statusbar ();
		statusbar.show ();

		var box = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		box.show ();

		// header bar
		var header = new Gtk.HeaderBar ();
		header.title ="Benchwell";
		header.subtitle ="version";
		header.show_close_button =true;
		header.show ();

		var windowBtnMenu = new Gtk.MenuButton();
		windowBtnMenu.show ();

		var addImg = new Benchwell.Image ("add-tab", "orange", 16);
		addImg.show ();
		windowBtnMenu.set_image (addImg);

		var windowMenu = new GLib.Menu ();
		windowBtnMenu.set_menu_model (windowMenu);

		windowMenu.append ("Window", "app.new");
		windowMenu.append ("Database", "win.new.db");
		windowMenu.append ("HTTP", "win.new.http");

		var appBtnMenu = new Gtk.MenuButton ();
		appBtnMenu.show ();
		/////////////

		set_titlebar (header);

		header.pack_start (windowBtnMenu);
		header.pack_end (appBtnMenu);

		box.pack_start (notebook, true, true, 0);
		box.pack_end (statusbar, false, false, 0);

		add(box);

		set_default_size (Config.window_width(), Config.window_height());
		move (Config.window_position_x(), Config.window_position_y());

		new_database_tab_menu.activate.connect ( () => {
			add_database_tab ();
		});

		close_menu.activate.connect ( () => {
			notebook.remove_page (notebook.get_current_page ());
		});
	}

	public void add_database_tab() {
		var tab  = new Benchwell.Tab ();
		tab.show ();

		var database = new Benchwell.Views.DBDatabase (this);
		database.notify["title"].connect ((s, p) => {
			tab.label.set_text (database.title);
		});
		database.show ();

		tab.label.set_text ( database.title );
		tab.pack_start (database, true, true, 0);

		notebook.append_page (tab, tab.header);
		tab.btn.clicked.connect( () => {
			notebook.remove_page (notebook.page_num (tab));
		});
	}
}

