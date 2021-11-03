public class Benchwell.Database.Condition {
	public Gtk.Switch active_switch;
	public Gtk.ListStore store;
	public Gtk.ComboBox field_combo;
	public Gtk.ComboBoxText operator_combo;
	public Gtk.Entry value_entry;
	public Benchwell.Button remove_btn;
	private Benchwell.ColDef[] _columns;
	public Benchwell.ColDef[] columns {
		get { return _columns; }
		set {
			_columns = value;
			store.clear ();

			Gtk.TreeIter iter;
			store.append (out iter);
			store.set_value (iter, 0 ,"");

			foreach (var column in _columns) {
				store.append (out iter);
				store.set_value (iter, 0, column.name);
			}
		}
	}

	public signal void search();
	public signal void ready (Benchwell.CondStmt values);
	public bool no_ready = false;

	public Condition () {
		store = new Gtk.ListStore (1, GLib.Type.STRING);
		Gtk.TreeIter iter;
		store.append (out iter);
		store.set_value (iter, 0 ,"");

		active_switch = new Gtk.Switch ();
		active_switch.active = true;
		active_switch.valign = Gtk.Align.CENTER;
		active_switch.vexpand = false;
		active_switch.show ();

		field_combo = new Gtk.ComboBox.with_model_and_entry (store);
		field_combo.id_column = 0;
		field_combo.set_entry_text_column (0);
		field_combo.show ();

		var completion = new Gtk.EntryCompletion ();
		completion.text_column = 0;
		completion.inline_completion = false;
		completion.inline_selection = true;
		completion.minimum_key_length = 1;
		completion.set_model (store);
		completion.match_selected.connect ((model, iter) => {
			field_combo.set_active_iter (iter);
			return true;
		});

		var entry = field_combo.get_child () as Gtk.Entry;
		entry.set_completion (completion);
		//entry.focus_out_event.connect (on_field_focus_out);

		operator_combo = new Gtk.ComboBoxText ();
		foreach (var op in Benchwell.Operator.all ()) {
			operator_combo.append (op.to_string(), op.to_string ());
		}
		operator_combo.set_active (0);
		operator_combo.show ();

		value_entry = new Gtk.Entry ();
		value_entry.show ();

		remove_btn = new Benchwell.Button ("close", Gtk.IconSize.MENU);
		remove_btn.show ();

		operator_combo.changed.connect (() => {
			var op = Benchwell.Operator.parse (operator_combo.get_active_text ());
			value_entry.sensitive = true;

			if (op == Benchwell.Operator.IsNull) {
				value_entry.sensitive = false;
			}
			if (op == Benchwell.Operator.IsNotNull) {
				value_entry.sensitive = false;
			}
		});
		entry.activate.connect( () => { search (); });
		value_entry.activate.connect( () => { search (); });

		operator_combo.changed.connect (on_change);
		entry.changed.connect (on_change);
		value_entry.changed.connect (on_change);
		active_switch.state_set.connect ((s) => {
			on_change();
			return false;
		});
	}

	public Benchwell.CondStmt? get_condition () {
		var column_name = selected_field ();
		if (column_name == "" || column_name == null) {
			return null;
		}

		Benchwell.ColDef? column = null;
		foreach (var c in columns) {
			if (c.name == column_name) {
				column = c;
				break;
			}
		}
		if (column == null) {
			return null;
		}

		var op = operator_combo.get_active_text ();
		if (op == null || op == "") {
			return null;
		}
		var operator = Benchwell.Operator.parse (op);
		if (operator == null) {
			return null;
		}

		var cvalue = value_entry.get_text ();
		if ((cvalue == null || cvalue == "") && operator != Benchwell.Operator.IsNull && operator != Benchwell.Operator.IsNotNull) {
			return null;
		}

		var stmt = new Benchwell.CondStmt ();
		stmt.field = column;
		stmt.op = operator;
		stmt.val = cvalue;
		stmt.enabled = active_switch.get_active () && active_switch.sensitive;

		return stmt;
	}

	private string? selected_field () {
		Gtk.TreeIter? iter;
		var ok = field_combo.get_active_iter (out iter);
		if (!ok) {
			return null;
		}
		GLib.Value val;
		store.get_value (iter, 0, out val);
		var column_name = val.get_string ();
		if (column_name == "" || column_name == null) {
			return null;
		}

		return column_name;
	}

	private void on_change () {
		if (no_ready) {
			return;
		}

		var stmt = get_condition ();
		if (stmt != null) {
			ready(stmt);
		}
	}
}

public class Benchwell.Database.Conditions : Gtk.Grid {
	private List<Condition> conditions;
	private Benchwell.ColDef[] _columns;

	public signal void ready (Benchwell.CondStmt stmt);
	public signal void search();

	public Conditions () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);
		set_name ("conditions");
		add_condition ();
	}

	public void set_columns (Benchwell.ColDef[] cols) {
		_columns = cols;
		conditions.foreach ((condition) => {
			condition.columns = _columns;
		});
	}

	public unowned Benchwell.ColDef[] get_columns () {
		return _columns;
	}

	public Benchwell.Database.Condition add_condition () {
		var cond = new Benchwell.Database.Condition ();
		cond.columns = _columns;

		var y = (int) conditions.length ();
		attach (cond.field_combo, 0, y, 2, 1);
		attach (cond.operator_combo, 2, y, 1, 1);
		attach (cond.value_entry, 3, y, 2, 1);
		attach (cond.active_switch, 5, y, 2, 1);
		attach (cond.remove_btn, 7, y, 1, 1);

		conditions.append (cond);

		cond.remove_btn.clicked.connect ( () => {
			var index = conditions.index (cond);
			remove_row (index);
			conditions.remove (cond);

			if (conditions.length () == 0) {
				add_condition ();
			}
		});

		cond.active_switch.grab_focus.connect (() => {
			if (conditions.index (cond) != conditions.length () - 1) {
				return;
			}

			add_condition ();
		});

		var entry = cond.field_combo.get_child () as Gtk.Entry;
		entry.grab_focus.connect (() => {
			if (conditions.index (cond) != conditions.length () - 1) {
				return;
			}

			add_condition ();
		});

		cond.value_entry.grab_focus.connect (() => {
			if (conditions.index (cond) != conditions.length () - 1) {
				return;
			}

			add_condition ();
		});

		cond.search.connect ( () => {
			search ();
		});

		cond.ready.connect( (stmt) => {
			ready(stmt);
		});

		return cond;
	}

	public Benchwell.CondStmt[] get_conditions () {
		Benchwell.CondStmt[] stmts = {};

		conditions.foreach( (condition) => {
			var c = condition.get_condition ();
			if (c != null) {
				stmts += c;
			}
		});

		return stmts;
	}

	public void clear (bool add_empty = true) {
		while (!conditions.is_empty ()) {
			conditions.remove (conditions.nth_data (0));
			remove_row (0);
		}

		if (add_empty)
			add_condition ();
	}

	public void rebuild (string[] filters) {
		clear (false);

		for (var i = 0; i < filters.length; i+=4) {
			var key = filters[i];
			var op = filters[i+1];
			var val = filters[i+2];

			var cond = add_condition ();
			cond.no_ready = true;
			cond.value_entry.set_text (val);
			cond.operator_combo.set_active_id (op);
			cond.field_combo.set_active_id (key);
			cond.active_switch.set_active (filters[i+3] == "true");
			cond.no_ready = false;
		}

		add_condition ();
	}
}
