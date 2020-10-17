namespace Benchwell {
	class Lexer {
		static string[] whitespaces = {"\n", " ", "\t", "\r"};

		protected string _input;
		public string input {
			get {
				return _input;
			}
		}

		protected int _position;
		public int position {
			get {
				return _position;
			}
		}

		public Lexer (string input) {
			_input = input;
		}

		public bool is_end_reached () {
			return _position == _input.length;
		}

		public string peek (int length = 1, int offset = 0) {
			var total_offset = offset + _position;

			if (!is_end_reached () && total_offset + length <= _input.length) {
				return _input.substring(total_offset, length);
			} else {
				return "";
			}
		}

		public string next (int length = 1, int offset = 0) {
			var str = peek (length, offset);

			if (str != "") {
				_position = _position + length + offset;
			}

			return str;
		}

		public void consume_whitespaces () {
			while(true) {
				if (is_end_reached() || !(peek (1) in Lexer.whitespaces)) {
					break;
				} else {
					next (1);
				}
			}
		}

		public bool expect (string wish, int offset = 0) {
			if (peek (wish.length, offset) == wish) {
				next (wish.length, offset);
				return true;
			}

			return false;
		}

		public bool expect_multiple (string[] wish) {
			foreach(string single in wish) {
				var matches = expect (single);

				if (matches) {
					return true;
				}
			}

			return false;
		}
	}

	enum NodeType {
		O_BRACKET,
		E_BRACKET,
		VAR,
		STRING;
	}

	class Node {
		public int start_at;
		public int end_at;
		public NodeType type;
	}

	class Parse {
		public Lexer lexer;
		public List<Node> tree;

		public void parse (string input) {
			lexer = new Lexer (input);

			while (!lexer.is_end_reached ()){
				var node = new Node ();
				var last_node = tree.last ();
				node.start_at = lexer.position;

				if (is_o_bracket ()) {
					node.type = NodeType.O_BRACKET;
					node.end_at = node.start_at + 2;
					tree.append (node);
				}

				if (is_e_bracket ()) {
					node.type = NodeType.E_BRACKET;
					node.end_at = node.start_at + 2;
					tree.append (node);
				}

				lexer.next ();
			}
		}

		private bool is_o_bracket() {
			return lexer.peek (2) == "{{";
		}

		private bool is_e_bracket() {
			return lexer.peek (2) == "}}";
		}
	}

	public static int main(string[] args) {
		var t = new Parse ();
		t.parse ("http://{{token}}/desk/api/v1/tickets.json");

		t.tree.foreach ((node) => {
			print (@"====$(node.start_at):$(node.end_at):$(node.type)\n");
		});

		return 0;
	}
}
