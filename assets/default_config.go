package assets

const DEFAULT_CONFIG = `
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS "settings" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name NVARCHAR(300) UNIQUE NOT NULL,
    value NVARCHAR(300) NOT NULL
);

CREATE TABLE IF NOT EXISTS "connections" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name text NOT NULL DEFAULT "",
    adapter text NOT NULL DEFAULT "",
    type text NOT NULL DEFAULT "",
    database text NULL DEFAULT "",
    host text NULL DEFAULT "",
    options text NULL DEFAULT "",
    user text NULL DEFAULT "",
    password text NULL DEFAULT "",
    port integer NULL DEFAULT "",
    encrypted BOOLEAN NOT NULL DEFAULT 0 CHECK (encrypted IN (0,1)),

    socket    text NOT NULL DEFAULT "",
    file      text NOT NULL DEFAULT "",
    sshHost   text NOT NULL DEFAULT "",
    sshAgent  text NOT NULL DEFAULT ""
);

CREATE TABLE IF NOT EXISTS "queries" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name text NOT NULL DEFAULT "",
    query text NOT NULL DEFAULT "",
    connections_id INTEGER NOT NULL,
    FOREIGN KEY(connections_id) REFERENCES connections(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO settings(name, value) VALUES
    ("gui.editor.word_wrap", "word"),
    ("gui.page_size", 100),
    ("gui.table_tab_position", "top"),
    ("gui.connection_tab_position", "top"),
    ("gui.initial_cell_width", 50),
    ("gui.dark_mode", 1),
    ("encryption_mode", "DBUS")`
