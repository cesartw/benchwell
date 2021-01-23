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
			// keep active field before replace columns
			if (field_combo.sensitive) {
				Gtk.TreeIter iter;
				GLib.Value? val;
				field_combo.get_active_iter (out iter);
				if (store.iter_is_valid (iter)) {
					store.get_value (iter, 0, out val);
					if (val != null) {
						_active_field = val.get_string ();
					} else {
						_active_field = null;
					}
				}
			}

			_columns = value;
			_update_fields ();
		}
	}

	private string? _active_field;

	public signal void search();

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
		completion.inline_completion = true;
		completion.inline_selection = true;
		completion.minimum_key_length = 1;
		completion.set_model (store);

		var entry = field_combo.get_child () as Gtk.Entry;
		entry.set_completion (completion);
		entry.focus_out_event.connect (on_field_focus_out);

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
	}

	public Benchwell.CondStmt? get_condition () {
		if (!active_switch.get_active () || !active_switch.sensitive) {
			return null;
		}

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

		return stmt;
	}

	private string? selected_field () {
		Gtk.TreeIter? iter;
		field_combo.get_active_iter (out iter);
		if (!store.iter_is_valid(iter)) {
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

	private void _update_fields () {
		store.clear ();

		bool enable = _active_field == "" || _active_field == null;
		foreach (var column in _columns) {
			Gtk.TreeIter iter;
			store.append (out iter);
			store.set_value (iter, 0, column.name);
			if (column.name == _active_field) {
				enable = true;
				field_combo.set_active_iter (iter);
			}
		}

		active_switch.sensitive = enable;
		field_combo.sensitive = enable;
		operator_combo.sensitive = enable;
		value_entry.sensitive = enable;
	}

	private bool on_field_focus_out () {
		var entry = field_combo.get_child () as Gtk.Entry;
		var written_field_name = entry.get_text ();
		var selected_field_name = selected_field ();

		if (written_field_name == selected_field_name) {
			return false;
		}

		var i = -1;
		foreach (var col in _columns) {
			i++;

			if (col.name != written_field_name) {
				continue;
			}

			field_combo.set_active(i);
		}
		return false;
	}
}

public class Benchwell.Database.Conditions : Gtk.Grid {
	private List<Benchwell.Database.Condition> conditions;
	private Benchwell.ColDef[] _columns;
	public Benchwell.ColDef[] columns {
		get { return _columns; }
		set {
			_columns = value;
			conditions.foreach ((condition) => {
				condition.columns = _columns;
			});
		}
	}

	public signal void search();

	public Conditions () {
		Object (
			row_spacing: 5,
			column_spacing: 5
		);
		set_name ("conditions");
		add_condition ();
	}

	public void add_condition () {
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
}
