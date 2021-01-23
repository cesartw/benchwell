CREATE TABLE "config" (
    id    INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name  NVARCHAR(300) NOT NULL,
    value NVARCHAR(300) NOT NULL
);

CREATE TABLE "environments" (
    id   INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name text NOT NULL
);

CREATE TABLE "environment_variables" (
    id             INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    key            text NOT NULL,
    value          text NOT NULL,
    enabled        BOOLEAN NOT NULL DEFAULT 1 CHECK (enabled IN (0,1)),
    environment_id integer NOT NULL
);

CREATE TABLE "db_connections" (
    id        INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name      text NOT NULL,
    adapter   text NOT NULL,
    type      text NOT NULL,
    database  text NULL,
    host      text NULL,
    options   text NULL,
    user      text NULL,
    port      integer NULL,
    encrypted BOOLEAN NOT NULL DEFAULT 0 CHECK (encrypted IN (0,1)),

    Socket    text NULL,
    File      text NULL,
    SshHost   text NULL,
    SshAgent  text NULL
);

CREATE TABLE "db_queries" (
    id             INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name           text NOT NULL,
    query          text NOT NULL,
	query_type     text NOT NULL DEFAULT "fav", -- history
    connections_id integer NOT NULL
);

CREATE TABLE "http_collections" (
    id    INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    count integer default 0,
    name  text NOT NULL
);

CREATE TABLE "http_items" (
	id                  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name                text NOT NULL,
	description         text NOT NULL DEFAULT "",
	parent_id           integer,
	is_folder           integer,
	count               integer default 0,
	sort                integer NOT NULL,
	http_collections_id integer NOT NULL,
	external_data       text NOT NULL DEFAULT "",

	method text DEFAULT "",
	url    text DEFAULT "",
	body   text DEFAULT "",
	mime   text DEFAULT "json"
);

CREATE TRIGGER increment_http_collections_count AFTER INSERT ON http_items
    BEGIN
        UPDATE http_collections SET count = count + 1 WHERE http_collections.id = NEW.http_collections_id;
    END;

CREATE TRIGGER decrement_http_collections_count AFTER DELETE ON http_items
    BEGIN
        UPDATE http_collections SET count = count - 1 WHERE http_collections.id = OLD.http_collections_id;
    END;

CREATE TABLE "http_kvs" (
    id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    key           text NOT NULL,
    value         text NOT NULL,
    type          text NOT NULL,
    http_items_id integer NOT NULL,
    sort          integer NOT NULL,
    enabled       BOOLEAN NOT NULL DEFAULT 1 CHECK (enabled IN (0,1))
);

INSERT INTO db_connections(name, adapter, type, database, host, options, user, port, encrypted)
      VALUES("localhost", "mysql", "tcp", "", "localhost", "", "", 3306, 0);

INSERT INTO config(name, value) VALUES("gui.editor.word_wrap", "word");
INSERT INTO config(name, value) VALUES("gui.page_size", 100);
INSERT INTO config(name, value) VALUES("gui.tab_position", "top");
INSERT INTO config(name, value) VALUES("gui.cell_width", 45);
INSERT INTO config(name, value) VALUES("gui.dark_mode", 1);
INSERT INTO config(name, value) VALUES("encryption_mode", "DBUS");
