package assets

const BRAND = `* {
	font-size: 1em;
	font-family: 'JetBrainsMono Nerd Font';
	-GtkScrollbar-has-backward-stepper: true;
	-GtkScrollbar-has-forward-stepper: true;
}

notebook > header > tabs > tab {
	border-left: 1px solid #dddddd;
	border-right: 1px solid #dddddd;
}

notebook, notebook header  {
	border: none;
}

notebook > header.top > tabs > tab {
	border-radius: 0.5em 0.5em 0em 0.0em;
	border-top: 1px solid #dddddd;
	border-bottom: none;
	margin-top: 5px;
}

notebook > header.bottom > tabs > tab {
	border-radius: 0em 0.0em 0.5em 0.5em;
	border-bottom: 1px solid #dddddd;
	border-top: none;
	margin-bottom: 5px;
}

notebook tabs tab label {
	padding: 0.2em;
	font-weight: 700;
}

notebook#MainNotebook header {
	border-bottom-width: 1px;
	border-bottom-style: solid;
}

textview {
	font-size: 1em;
}

statusbar * {
	font-size: 0.9em;
}

/*header*/
treeview.view button box label {
	padding: 0.5em;
}

list row {
	padding: 0.5em 0em;
}

list row image {
	margin-right: 1px;
}

#conditions {
	margin-left: 5px;
}

#logger {
	font-size: 0.9em;
}

#queryactionbar box {
	border-width: 0;
}

overlay > box {
	background-color: rgba(0, 0, 0, 0.4);
}
`
