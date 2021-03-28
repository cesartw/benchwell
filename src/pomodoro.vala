public enum Benchwell.PomodoroStep {
	WORK,
	BREAK,
	LONG_BREAK
}

public class Benchwell.PomodoroCycler : Object {
	public Benchwell.PomodoroStep state;
	public int pomodoro_count = 0;

	// returns whether to stop
	public void next (Benchwell.Pomodoro w, int64 elapse) {
		var seconds = elapse / 1000000;
		var pomodoro_duration = Config.settings.get_int64("pomodoro-duration");
		var break_duration = Config.settings.get_int64("pomodoro-break-duration");
		var long_break_duration = Config.settings.get_int64("pomodoro-long-break-duration");
		var pomodoro_for_long_break = Config.settings.get_int64("pomodoro-set-count");

		switch (state) {
			case Benchwell.PomodoroStep.WORK:
				if (seconds == pomodoro_duration) {
					state = Benchwell.PomodoroStep.BREAK;
					pomodoro_count++;
		 			if (pomodoro_count == pomodoro_for_long_break) {
						state = Benchwell.PomodoroStep.LONG_BREAK;
					}

					w.reset ();
				}

				break;
			case Benchwell.PomodoroStep.BREAK:
				if (seconds == break_duration) {
					state = Benchwell.PomodoroStep.WORK;
					w.reset ();
				}

				break;
			case Benchwell.PomodoroStep.LONG_BREAK:
				pomodoro_count = 0;

				if (seconds == long_break_duration) {
					state = Benchwell.PomodoroStep.WORK;
					w.reset ();
				}

				break;
		}
	}

	public void reset () {
		state = Benchwell.PomodoroStep.WORK;
		pomodoro_count = 0;
	}
}

public class Benchwell.Pomodoro : Gtk.Box {
	private Gtk.Label clock_lbl;
	private Gtk.Label counter_lbl;
	public Gtk.Button pause_btn;
	private int64 last_cycle_start_at = 0;
	private int64 since_paused = 0;
	private bool paused = false;
	private bool stopped = true;
	private Gtk.Image pause_img;
	private Gtk.Image start_img;
	private Gtk.Image stop_img;

	private Benchwell.PomodoroCycler cycler;

	public Pomodoro () {
		Object (
			orientation: Gtk.Orientation.HORIZONTAL,
			spacing: 5
		);

		cycler = new Benchwell.PomodoroCycler ();
		cycler.reset ();

		pause_img = new Gtk.Image.from_icon_name ("media-playback-pause", Gtk.IconSize.SMALL_TOOLBAR);
		pause_img.show ();

		start_img = new Gtk.Image.from_icon_name ("media-playback-start", Gtk.IconSize.SMALL_TOOLBAR);
		start_img.show ();

		stop_img = new Gtk.Image.from_icon_name ("media-playback-stop", Gtk.IconSize.SMALL_TOOLBAR);
		stop_img .show ();

		clock_lbl = new Gtk.Label ("00:00");
		clock_lbl.show ();

		counter_lbl = new Gtk.Label (_("0 of %lld").printf (Config.settings.get_int64("pomodoro-set-count")));
		counter_lbl.get_style_context ().add_class ("pomodoro-counter");
		counter_lbl.show ();

		pause_btn = new Gtk.Button ();
		pause_btn.image = start_img;
		pause_btn.show ();

		var vbox = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
		vbox.show ();

		vbox.pack_start (clock_lbl);
		vbox.pack_start (counter_lbl);

		pack_start (vbox);
		pack_end (pause_btn);

		pause_btn.clicked.connect (on_pause);
	}

	public void start () {
		stopped = false;
		paused = false;
		last_cycle_start_at = get_real_time ();
		pause_btn.image = pause_img;
		//GLib.Timeout.add_seconds (1, tick); // add_tick_callback works a lot better
		clock_lbl.add_tick_callback (tick);
	}

	public void stop () {
		last_cycle_start_at = 0;
		since_paused = 0;
		pause_btn.image = start_img;
		stopped = true;
		cycler.reset ();
	}

	public void reset () {
		last_cycle_start_at = get_real_time ();
		since_paused = 0;
	}

	private bool tick () {
		if (paused || stopped) {
			return false;
		}

		var elapse = get_real_time () - last_cycle_start_at;
		update_label (elapse);
		cycler.next (this, since_paused + elapse);

		pause_btn.tooltip_text = _("%d set").printf (cycler.pomodoro_count+1);

		if (cycler.state == Benchwell.PomodoroStep.WORK)
			counter_lbl.set_text (_("%lld of %lld").printf (cycler.pomodoro_count+1, Config.settings.get_int64("pomodoro-set-count")));
		else
			counter_lbl.set_text (_("BREAK"));

		return true;
	}

	private void update_label (int64 elapse) {
		int64 seconds = 0;
		int64 minutes = 0;

		seconds = (since_paused + elapse) / 1000000;
		minutes = seconds / 60;
		seconds = seconds % 60;

		clock_lbl.set_text ("%02lld:%02lld".printf (minutes, seconds));
	}

	private void on_pause () {
		// resume
		if (paused) {
			last_cycle_start_at = get_real_time ();
			paused = false;
			pause_btn.image = pause_img;
			GLib.Timeout.add_seconds_full (GLib.Priority.HIGH_IDLE, 1, tick);
			update_label (get_real_time () - last_cycle_start_at);
			return;
		}

		// starting
		if (last_cycle_start_at == 0) {
			start ();
			return;
		}

		since_paused = since_paused + get_real_time () - last_cycle_start_at;
		pause_btn.image = start_img;
		paused = true;
	}
}
