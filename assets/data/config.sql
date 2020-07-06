CREATE TABLE "config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name NVARCHAR(300) NOT NULL,
    value NVARCHAR(300) NOT NULL
);

CREATE TABLE "connections" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name text NOT NULL,
    adapter text NOT NULL,
    type text NOT NULL,
    database text NULL,
    host text NULL,
    options text NULL,
    user text NULL,
    password text NULL,
    port integer NULL,
    encrypted BOOLEAN NOT NULL DEFAULT 0 CHECK (encrypted IN (0,1))

	Socket    text NULL,
	File      text NULL,
	SshHost   text NULL,
	SshAgent  text NULL
);

INSERT INTO connections(name, adapter, type, database, host, options, user, password, port, encrypted)
      VALUES("localhost", "mysql", "tcp", "", "localhost", "", "", "", 3306, 0);

INSERT INTO config(name, value) VALUES("gui.editor.word_wrap", "word");
INSERT INTO config(name, value) VALUES("gui.page_size", 100);
INSERT INTO config(name, value) VALUES("gui.table_tab_position", "top");
INSERT INTO config(name, value) VALUES("gui.connection_tab_position", "top");
INSERT INTO config(name, value) VALUES("gui.initial_cell_width", 50);
INSERT INTO config(name, value) VALUES("gui.dark_mode", 1);
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.background", " bg,#3d3d3d");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.comment", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.commenthashbang", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.commentmultiline", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.commentpreproc", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.commentsingle", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.commentspecial", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.error", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.generic", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericdeleted", "#8b080b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericemph", "#f8f8f2 underline");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericerror", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericheading", "#f8f8f2 bold");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericinserted", "#f8f8f2 bold");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericoutput", "#44475a");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericprompt", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericstrong", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericsubheading", "#f8f8f2 bold");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.generictraceback", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.genericunderline", "underline");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.keyword", "#ff5f5f");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.keywordconstant", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.keyworddeclaration", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.keywordnamespace", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.keywordpseudo", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.keywordreserved", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.keywordtype", "#8be9fd");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literal", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literaldate", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalnumber", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalnumberbin", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalnumberfloat", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalnumberhex", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalnumberinteger", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalnumberintegerlong", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalnumberoct", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstring", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringbacktick", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringchar", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringdoc", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringdouble", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringescape", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringheredoc", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringinterpol", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringother", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringregex", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringsingle", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.literalstringsymbol", "#f1fa8c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.name", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.nameattribute", "#50fa7b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namebuiltin", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namebuiltinpseudo", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.nameclass", "#50fa7b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.nameconstant", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namedecorator", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.nameentity", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.nameexception", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namefunction", "#50fa7b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namelabel", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namenamespace", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.nameother", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.nametag", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namevariable", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namevariableclass", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namevariableglobal", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.namevariableinstance", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.operator", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.operatorword", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.other", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.punctuation", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.text", "#f8f8f2");
INSERT INTO config(name, value) VALUES("gui.editor.theme.dark.textwhitespace", "#f8f8f2";
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.background", " bg,#ffffff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.comment", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.commenthashbang", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.commentmultiline", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.commentpreproc", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.commentsingle", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.commentspecial", "#00f000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.error", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.generic", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericdeleted", "#8b080b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericemph", "#000000 underline");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericerror", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericheading", "#000000 bold");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericinserted", "#000000 bold");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericoutput", "#44475a");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericprompt", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericstrong", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericsubheading", "#000000 bold");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.generictraceback", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.genericunderline", "underline");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.keyword", "#ff5f5f");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.keywordconstant", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.keyworddeclaration", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.keywordnamespace", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.keywordpseudo", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.keywordreserved", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.keywordtype", "#8be9fd");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literal", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literaldate", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalnumber", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalnumberbin", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalnumberfloat", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalnumberhex", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalnumberinteger", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalnumberintegerlong", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalnumberoct", "#b1b1ff");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstring", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringbacktick", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringchar", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringdoc", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringdouble", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringescape", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringheredoc", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringinterpol", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringother", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringregex", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringsingle", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.literalstringsymbol", "#919a9c");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.name", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.nameattribute", "#50fa7b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namebuiltin", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namebuiltinpseudo", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.nameclass", "#50fa7b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.nameconstant", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namedecorator", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.nameentity", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.nameexception", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namefunction", "#50fa7b");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namelabel", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namenamespace", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.nameother", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.nametag", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namevariable", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namevariableclass", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namevariableglobal", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.namevariableinstance", "#8be9fd italic");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.operator", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.operatorword", "#ff79c6");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.other", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.punctuation", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.text", "#000000");
INSERT INTO config(name, value) VALUES("gui.editor.theme.light.textwhitespace", "#000000");
INSERT INTO config(name, value) VALUES("encryption_mode", "DBUS");
