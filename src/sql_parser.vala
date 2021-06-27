namespace Benchwell {
	namespace SQL {

		[Compact]
		public class GrammarNode {
			public TokenType[] token_types;
			public GrammarNode next_node;
		}

		public class Parser {
			private GrammarNode sql_grammar;

			public Parser () {
				//var select_grammar = new GrammarNode ();
				//select_grammar.token_types = { TokenType.SELECT };
				//var column_grammar = new GrammarNode ();
				//column_grammar
			}
		}

		public enum TokenType {
			ROOT,
			ADD,
			ALL,
			ALTER,
			AND,
			AS,
			ASC,
			BAREWORD,
			BETWEEN,
			BY,
			CASE,
			COLUMN,
			CONSTRAINT,
			CREATE,
			DATABASE,
			DEFAULT,
			DESC,
			DISTINCT,
			DROP,
			EXISTS,
			FOREIGN,
			FORCE,
			FROM,
			GROUP,
			HAVING,
			IN,
			INDEX,
			INNER,
			INSERT,
			INTO,
			IS,
			JOIN,
			KEY,
			LEFT,
			LIKE,
			LIMIT,
			NUMBER,
			NULL,
			NOT,
			OR,
			ORDER,
			OUTER,
			PRIMARY,
			REMOVE,
			RIGHT,
			SELECT,
			SET,
			SHOW,
			STRING,
			SYMBOL,
			TABLE,
			TRUNCATE,
			UNION,
			UPDATE,
			VALUES,
			VIEW,
			WHERE,

			OPEN_PARENTHESIS,
			CLOSE_PARENTHESIS,
			DOT,
			COMMA;

			public static TokenType[] all () {
				return {
					ADD,
					ALL,
					ALTER,
					AND,
					AS,
					ASC,
					BAREWORD,
					BETWEEN,
					BY,
					CASE,
					COLUMN,
					CONSTRAINT,
					CREATE,
					DATABASE,
					DEFAULT,
					DESC,
					DISTINCT,
					DROP,
					EXISTS,
					FOREIGN,
					FORCE,
					FROM,
					GROUP,
					HAVING,
					IN,
					INDEX,
					INNER,
					INSERT,
					INTO,
					IS,
					JOIN,
					KEY,
					LEFT,
					LIKE,
					LIMIT,
					NUMBER,
					NULL,
					NOT,
					OR,
					ORDER,
					OUTER,
					PRIMARY,
					REMOVE,
					RIGHT,
					SELECT,
					SET,
					SHOW,
					STRING,
					SYMBOL,
					TABLE,
					TRUNCATE,
					UNION,
					UPDATE,
					VALUES,
					VIEW,
					WHERE,

					OPEN_PARENTHESIS,
					CLOSE_PARENTHESIS,
					DOT,
					COMMA
				};
			}
		}

		[Compact]
		public struct Token {
			int start;
			int end;
			int len;
			char[] chars;
			TokenType type;

			public void append_char (char c) {
				if (chars.length == 0) {
					len = 0;
					chars = {0};
				}
				var lchars = this.chars;
				lchars[len] = c;
				lchars+= 0;
				len++;
				chars = lchars;
			}
		}

		[Compact]
		public struct TokenState {
			char string_wrapping_c;
			bool in_number;
			bool in_keyword;
			bool escaping;
		}

		public class Tokenizer : Object {
			public static Token[] parse (string sql) {
				Token[] tokens = {};

				TokenState state = TokenState ();
				Token ctoken = Token ();

				for (var i = 0; i < sql.length; i++) {
					var c = sql[i];

					if (state.escaping) {
						state.escaping =  false;
						ctoken.append_char (c);
						continue;
					}

					if (c == '\\') {
						// skip next character
						state.escaping = true;
						continue;
					}


					switch (c) {
						case '\'', '"', '`':
							if (state.in_number) {
								// end number
								state.in_number = false;
								ctoken.end = i;
								tokens += ctoken;

								// start string
								state.string_wrapping_c = c;
								ctoken = Token () {
									type = TokenType.STRING,
									start = i
								};
								continue;
							}

							// end string
							if (c != 0 && state.string_wrapping_c == c) {
								ctoken.append_char (c);
								ctoken.end = i;
								tokens += ctoken;
								state.string_wrapping_c = 0;
								ctoken = Token ();
								continue;
							}

							// start string
							ctoken = Token () {
								type = TokenType.STRING,
								start = i
							};
							ctoken.append_char (c);
							state.string_wrapping_c = c;
							continue;
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							if (state.in_number) {
								ctoken.append_char (c);
								continue;
							}

							if (state.in_keyword) {
								ctoken.append_char (c);
								continue;
							}

							state.in_number = true;
							ctoken = Token () {
								type = TokenType.NUMBER,
								start = i
							};
							ctoken.append_char (c);
							continue;
						case '.':
							// continue number
							if (state.in_number) {
								ctoken.append_char (c);
								continue;
							}

							// continue string
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							// end keyword
							if (state.in_keyword) {
								ctoken.end = i - 1;
								tokens += ctoken;
								state.in_keyword = false;
							}

							// raw dot
							ctoken = Token () {
								type = TokenType.SYMBOL,
								start = i,
								end = i
							};
							ctoken.append_char (c);
							tokens += ctoken;

							ctoken = Token ();
							continue;
						case ';', ',':
							// continue string
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							// end number
							if (state.in_number) {
								state.in_number = false;
								ctoken.end = i - 1;
								tokens += ctoken;
							}

							// end keyword
							if (state.in_keyword) {
								state.in_keyword = false;
								ctoken.end = i - 1;
								tokens += ctoken;
							}

							ctoken = Token () {
								type = TokenType.SYMBOL,
								start = i,
								end = i
							};
							ctoken.append_char (c);
							tokens += ctoken;

							ctoken = Token ();
							continue;
						case ' ', '\t', '\n':
							// continue string
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							// end keyword
							if (state.in_keyword) {
								state.in_keyword = false;
								ctoken.end = i - 1;
								tokens += ctoken;
								ctoken = Token ();
								continue;
							}

							// end number
							if (state.in_number) {
								state.in_number = false;
								ctoken.end = i - 1;
								tokens += ctoken;
								ctoken = Token ();
								continue;
							}

							continue;
						case '=', '<', '>', '!', '+', '-', '*':
							// continue string
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							// end keyword
							if (state.in_keyword) {
								state.in_keyword = false;
								ctoken.end = i - 1;
								tokens += ctoken;
								ctoken = Token () {
									type = TokenType.SYMBOL,
									start = i,
									end = i
								};
								ctoken.append_char (c);
								tokens += ctoken;
								ctoken = Token ();
								continue;
							}

							// end number
							if (state.in_number) {
								state.in_number = false;
								ctoken.end = i - 1;
								tokens += ctoken;
								ctoken = Token ();
								continue;
							}

							ctoken = Token () {
								type = TokenType.SYMBOL,
								start = i
							};
							ctoken.append_char (c);
							ctoken.end = i;
							tokens += ctoken;
							ctoken = Token ();

							continue;
						case '(':
							// continue string
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							// end number
							if (state.in_number) {
								state.in_number = false;
								ctoken.end = i - 1;
								tokens += ctoken;
							}

							// end keyword
							if (state.in_keyword) {
								state.in_keyword = false;
								ctoken.end = i - 1;
								tokens += ctoken;
							}

							ctoken = Token () {
								type = TokenType.SYMBOL,
								start = i,
								end = i
							};
							ctoken.append_char (c);
							tokens += ctoken;

							ctoken = Token ();
							continue;
						case ')':
							// continue string
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							// end number
							if (state.in_number) {
								state.in_number = false;
								ctoken.end = i - 1;
								tokens += ctoken;
							}

							// end keyword
							if (state.in_keyword) {
								state.in_keyword = false;
								ctoken.end = i - 1;
								tokens += ctoken;
							}

							ctoken = Token () {
								type = TokenType.SYMBOL,
								start = i,
								end = i
							};
							ctoken.append_char (c);
							tokens += ctoken;

							ctoken = Token ();
							continue;
						default:
							// continue string
							if (state.string_wrapping_c != 0) {
								ctoken.append_char (c);
								continue;
							}

							if (state.in_keyword) {
								ctoken.append_char (c);
								continue;
							}

							state.in_keyword = true;
							ctoken = Token () {
								type = TokenType.BAREWORD,
								start = i
							};
							ctoken.append_char (c);

							continue;
					}
				}

				if (ctoken.type in TokenType.all ()) {
					ctoken.end = sql.length - 1;
					tokens += ctoken;
				}

				for (var i = 0; i < tokens.length; i++) {
					if (tokens[i].chars.length == 0) {
						continue;
					}

					tokens[i].type = set_token_type ((string) tokens[i].chars, tokens[i].type);
				}

				return tokens;
			}

			private static TokenType set_token_type (string token, TokenType default_type = TokenType.BAREWORD) {
				var ntoken = token.casefold ().normalize (-1, GLib.NormalizeMode.ALL_COMPOSE);

				if (ntoken == "null") {
					return TokenType.NULL;
				}

				if (ntoken == "add") {
					return TokenType.ADD;
				}
				if (ntoken == "all") {
					return TokenType.ALL;
				}
				if (ntoken == "alter") {
					return TokenType.ALTER;
				}
				if (ntoken == "and") {
					return TokenType.AND;
				}
				if (ntoken == "as") {
					return TokenType.AS;
				}
				if (ntoken == "asc") {
					return TokenType.ASC;
				}
				if (ntoken == "bareword") {
					return TokenType.BAREWORD;
				}
				if (ntoken == "between") {
					return TokenType.BETWEEN;
				}
				if (ntoken == "by") {
					return TokenType.BY;
				}
				if (ntoken == "case") {
					return TokenType.CASE;
				}
				if (ntoken == "column") {
					return TokenType.COLUMN;
				}
				if (ntoken == "constraint") {
					return TokenType.CONSTRAINT;
				}
				if (ntoken == "create") {
					return TokenType.CREATE;
				}
				if (ntoken == "database") {
					return TokenType.DATABASE;
				}
				if (ntoken == "default") {
					return TokenType.DEFAULT;
				}
				if (ntoken == "desc") {
					return TokenType.DESC;
				}
				if (ntoken == "distinct") {
					return TokenType.DISTINCT;
				}
				if (ntoken == "drop") {
					return TokenType.DROP;
				}
				if (ntoken == "exists") {
					return TokenType.EXISTS;
				}
				if (ntoken == "foreign") {
					return TokenType.FOREIGN;
				}
				if (ntoken == "force") {
					return TokenType.FORCE;
				}
				if (ntoken == "from") {
					return TokenType.FROM;
				}
				if (ntoken == "group") {
					return TokenType.GROUP;
				}
				if (ntoken == "having") {
					return TokenType.HAVING;
				}
				if (ntoken == "in") {
					return TokenType.IN;
				}
				if (ntoken == "index") {
					return TokenType.INDEX;
				}
				if (ntoken == "inner") {
					return TokenType.INNER;
				}
				if (ntoken == "insert") {
					return TokenType.INSERT;
				}
				if (ntoken == "into") {
					return TokenType.INTO;
				}
				if (ntoken == "is") {
					return TokenType.IS;
				}
				if (ntoken == "join") {
					return TokenType.JOIN;
				}
				if (ntoken == "key") {
					return TokenType.KEY;
				}
				if (ntoken == "left") {
					return TokenType.LEFT;
				}
				if (ntoken == "like") {
					return TokenType.LIKE;
				}
				if (ntoken == "limit") {
					return TokenType.LIMIT;
				}
				if (ntoken == "number") {
					return TokenType.NUMBER;
				}
				if (ntoken == "null") {
					return TokenType.NULL;
				}
				if (ntoken == "not") {
					return TokenType.NOT;
				}
				if (ntoken == "or") {
					return TokenType.OR;
				}
				if (ntoken == "order") {
					return TokenType.ORDER;
				}
				if (ntoken == "outer") {
					return TokenType.OUTER;
				}
				if (ntoken == "primary") {
					return TokenType.PRIMARY;
				}
				if (ntoken == "remove") {
					return TokenType.REMOVE;
				}
				if (ntoken == "right") {
					return TokenType.RIGHT;
				}
				if (ntoken == "select") {
					return TokenType.SELECT;
				}
				if (ntoken == "set") {
					return TokenType.SET;
				}
				if (ntoken == "show") {
					return TokenType.SHOW;
				}
				if (ntoken == "table") {
					return TokenType.TABLE;
				}
				if (ntoken == "truncate") {
					return TokenType.TRUNCATE;
				}
				if (ntoken == "union") {
					return TokenType.UNION;
				}
				if (ntoken == "update") {
					return TokenType.UPDATE;
				}
				if (ntoken == "values") {
					return TokenType.VALUES;
				}
				if (ntoken == "view") {
					return TokenType.VIEW;
				}
				if (ntoken == "where") {
					return TokenType.WHERE;
				}
				if (ntoken == "(") {
					return TokenType.OPEN_PARENTHESIS;
				}
				if (ntoken == ")") {
					return TokenType.CLOSE_PARENTHESIS;
				}
				if (ntoken == ",") {
					return TokenType.COMMA;
				}
				if (ntoken == ".") {
					return TokenType.DOT;
				}

				return default_type;
			}
		}

		public class TableCompletion : Object, Gtk.SourceCompletionProvider {
			public weak DatabaseService db { get; construct; }
			private string prefix;
			public TableCompletion (DatabaseService db) {
				Object (
					db: db
				);
			}

			public string get_name () {
				return _("Tables");
			}

			public bool match (Gtk.SourceCompletionContext context) {
				var end = context.iter;
				Gtk.TextIter start;
				var buffer = context.completion.view.get_buffer ();
				buffer.get_start_iter (out start);

				var tokens = Benchwell.SQL.Tokenizer.parse (buffer.get_text (start, end, false));
				if (tokens.length == 0) {
					return false;
				}

				//for (var i = 0; i < tokens.length; i++) {
				//print (@"=======\"$((string) tokens[i].chars)\" $(tokens[i].type)\n");
				//}

				var ok = tokens[tokens.length - 1].type == TokenType.FROM;

				if (!ok && tokens.length > 1) {
					ok = tokens[tokens.length - 2].type == TokenType.FROM &&
						(tokens[tokens.length - 1].type == TokenType.STRING || tokens[tokens.length - 1].type == TokenType.BAREWORD);
					if (ok) {
						prefix = (string) tokens[tokens.length - 1].chars;
						if (tokens[tokens.length - 1].type == TokenType.STRING) {
							prefix = prefix.replace ("\"", "");
						}
					}
				}

				//print (@"=======ok $(ok)\n");
				return ok;
			}

			public void populate (Gtk.SourceCompletionContext context) {
				GLib.List<Gtk.SourceCompletionProposal> proposals = null;

				foreach (var table in db.tables) {
					if (table.ttype != TableType.Regular)
						continue;

					if (prefix == null || table.name.has_prefix (prefix)) {
						proposals.append (new CompletionProposal (table.name, prefix));
					}
				}
				prefix = null;
				context.add_proposals (this, proposals, true);
			}
		}

		public class CompletionProposal : Object, Gtk.SourceCompletionProposal {
			public string name { get; construct; }
			public string? prefix { get; construct; }

			public CompletionProposal (string name, string? prefix = null) {
				Object(
					name: name,
					prefix: prefix
				);
			}

			public bool equal (Gtk.SourceCompletionProposal other) {
				var other_of_this_type = other as CompletionProposal;
				if (other_of_this_type == null)
					return false;

				return this.name == other_of_this_type.name;
			}

			public string get_text () {
				return name;
			}

			// get_markup has priority
			public string get_label () {
				return "";
			}

			public unowned GLib.Icon? get_gicon () {
				return null;
			}

			public unowned Gdk.Pixbuf? get_icon () {
				return null;
			}

			public unowned string? get_icon_name () {
				return "bw-table";
			}

			public string? get_info () {
				return "";
			}

			public string get_markup () {
				if (prefix != null) {
					return name.replace (prefix, @"<b>$prefix</b>");
				}
				return name;
			}
		}

	}
}
