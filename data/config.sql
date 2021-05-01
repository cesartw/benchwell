CREATE TABLE "environments" (
    name TEXT NOT NULL
);

CREATE TABLE "environment_variables" (
	key            TEXT NOT NULL,
	value          TEXT NOT NULL,
	enabled        BOOLEAN NOT NULL DEFAULT 1 CHECK (enabled IN (0,1)),
	environment_id INTEGER NOT NULL,
	kvtype      INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE "db_connections" (
	name      TEXT NOT NULL,
	adapter   TEXT NOT NULL,
	type      TEXT NOT NULL,
	database  TEXT NULL,
	host      TEXT NULL,
	options   TEXT NULL,
	user      TEXT NULL,
	port      INTEGER NULL,
	encrypted BOOLEAN NOT NULL DEFAULT 0 CHECK (encrypted IN (0,1)),

	Socket    TEXT NULL,
	File      TEXT NULL,
	SshHost   TEXT NULL,
	SshAgent  TEXT NULL
);

CREATE TABLE "db_queries" (
	name           TEXT NOT NULL,
	query          TEXT NOT NULL,
	query_type     TEXT NOT NULL DEFAULT "fav", -- history
	connections_id INTEGER NOT NULL,
	created_at     INTEGER NOT NULL
);

CREATE TABLE "http_collections" (
	count integer default 0,
	name  TEXT NOT NULL
);

CREATE TABLE "http_items" (
	name                TEXT NOT NULL,
	description         TEXT NOT NULL DEFAULT "",
	parent_id           INTEGER,
	is_folder           INTEGER,
	count               INTEGER default 0,
	sort                INTEGER NOT NULL,
	http_collections_id INTEGER NOT NULL,
	external_data       TEXT NOT NULL DEFAULT "",

	method TEXT DEFAULT "",
	url    TEXT DEFAULT "",
	body   TEXT DEFAULT "",
	mime   TEXT DEFAULT "json",
);

CREATE TABLE "http_responses" (
	http_items_id INTEGER NOT NULL,
	url           TEXT DEFAULT "",
	method        TEXT DEFAULT "",
	headers	      TEXT DEFAULT "",
	body          TEXT DEFAULT "",
	content_type  TEXT DEFAULT "",
	duration      INTEGER DEFAULT 0,
	code          INTEGER DEFAULT 0,
	created_at    INTEGER NOT NULL
);

CREATE TRIGGER increment_http_collections_count AFTER INSERT ON http_items
    BEGIN
        UPDATE http_collections SET count = count + 1 WHERE http_collections.rowid = NEW.http_collections_id;
    END;

CREATE TRIGGER decrement_http_collections_count AFTER DELETE ON http_items
    BEGIN
        UPDATE http_collections SET count = count - 1 WHERE http_collections.rowid = OLD.http_collections_id;
    END;

CREATE TABLE "http_kvs" (
	key           TEXT NOT NULL,
	value         TEXT NOT NULL,
	type          TEXT NOT NULL,
	kvtype     INTEGER NOT NULL DEFAULT 1,
	http_items_id INTEGER NOT NULL,
	sort          INTEGER NOT NULL,
	enabled       BOOLEAN NOT NULL DEFAULT 1 CHECK (enabled IN (0,1))
);

INSERT INTO db_connections(name, adapter, type, database, host, options, user, port, encrypted)
      VALUES("localhost", "mysql", "tcp", "", "localhost", "", "", 3306, 0);
