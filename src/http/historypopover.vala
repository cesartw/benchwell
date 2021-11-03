public class Benchwell.Http.HistoryPopover : Gtk.Popover {
	private Gtk.Grid grid;
	public weak Benchwell.HttpItem? item { owned get; set; }

	public signal void result_activated (Result result);

	public HistoryPopover (Gtk.Widget relative_to, Benchwell.HttpItem? item) {
		Object (
			relative_to: relative_to,
			item: item
		);

		build_results ();
		notify["item"].connect ((sender, property) => {
			if (this.item == null)
				return;
			this.item.response_added.connect (build_results);
			build_results ();
		});
	}

	public void disconnect_from_item () {
		if (item == null)
			return;
		item.response_added.disconnect (build_results);
	}

	private void build_results () {
		if (item == null) {
			return;
		}

		if (grid != null) {
			remove (grid);
			grid.destroy ();
		}

		grid = new Gtk.Grid ();
		grid.row_spacing = 5;
		grid.column_spacing = 5;
		add (grid);

		var i = 0;
		foreach (var response in item.responses) {
			add_response (i, response);
			i++;
		}

		grid.show_all ();
	}

	private void add_response (int at, Result response) {
		var status_label = new Benchwell.Label ("0");
		status_label.set_text (response.status.to_string ());
		status_label.get_style_context ().add_class ("tag");

		if (response.status >= 200 && response.status < 300) {
			status_label.get_style_context ().add_class("ok");
		}
		if (response.status >= 300 && response.status < 500) {
			status_label.get_style_context ().add_class("warning");
		}
		if (response.status >= 500) {
			status_label.get_style_context ().add_class("bad");
		}

		var duration_label  = new Benchwell.Label (@"$(response.duration/1000)ms");
		duration_label.get_style_context ().add_class ("tag");

		var size_label  = new Benchwell.Label ("0KB");
		size_label.get_style_context ().add_class ("tag");
		size_label.set_text (@"0B");
		if (response.body != null) {
			if (response.body.length < 1024) {
				size_label.set_text (@"$(response.body.length)B");
			} else {
				size_label.set_text (@"$(response.body.length/1024)kB");
			}
		}

		var url_label = new Benchwell.Label (@"$(response.method) $(response.url)");
		url_label.max_width_chars = 30;
		url_label.ellipsize = Pango.EllipsizeMode.END;
		url_label.tooltip_text = response.url;

		var time_fmt = "%Y-%m-%d %H:%M:%S";
		var now = new DateTime.now_local ();
		if (now.get_day_of_year () == response.created_at.get_day_of_year ())
			time_fmt = "%H:%M:%S";

		var time_label = new Benchwell.Label (response.created_at.format (time_fmt));

		grid.attach (status_label, 0, at, 1, 1);
		grid.attach (duration_label, 1, at, 1, 1);
		grid.attach (size_label, 2, at, 1, 1);
		grid.attach (url_label, 3, at, 1, 1);
		grid.attach (time_label, 4, at, 1, 1);

		status_label.clicked.connect ( () => { result_activated (response); });
		duration_label.clicked.connect ( () => { result_activated (response); });
		size_label.clicked.connect ( () => { result_activated (response); });
		url_label.clicked.connect ( () => { result_activated (response); });
		time_label.clicked.connect ( () => { result_activated (response); });
	}
}
