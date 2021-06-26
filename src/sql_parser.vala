namespace Benchwell {
	namespace SQL {
		public enum TokenType {
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
			CLOSE_PARENTHESIS
			;

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
					CLOSE_PARENTHESIS
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
			public Token[] parse (string sql) {
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
						case ';':
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
						case ' ', '\t':
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

			private TokenType set_token_type (string token, TokenType default_type = TokenType.BAREWORD) {
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
				if (ntoken == "string") {
					return TokenType.STRING;
				}
				if (ntoken == "symbol") {
					return TokenType.SYMBOL;
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

				return default_type;
			}
		}

		public class Node {
			public Token? parent;
			public int start;
			public int end;
		}

		public class SelectNode : Node {
			public SelectFieldListNode columns;
			public SelectFromNode from;
			public WhereNode where;
		}

		public class SelectFieldListNode : Node {
			public ColumnNameNode[] fields;
		}

		public class ColumnNameNode : Node {
			public Token token;
		}

		public class SelectFromNode : Node {
			public Token table;
			public SelectNode select;
		}

		public class WhereNode : Node {
			public BinOpNode[] statements;
		}

		public class BinOpNode : Node {
			public ColumnNameNode column;
			public OperatorNode operator;
			public ValueNode? val;
			public BinOpNode? binNode;
		}

		public class OperatorNode : Node {
			public Token token;
		}

		public class ValueNode : Node {
			public Token token;
		}

		public class AST : Object {
		}
	}
}
