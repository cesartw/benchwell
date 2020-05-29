package assets

const THEME_LIGHT = `/*************************** Check and Radio buttons * */
* { padding: 0; -GtkToolButton-icon-spacing: 4; -GtkTextView-error-underline-color: #cc0000; -GtkScrolledWindow-scrollbar-spacing: 0; -GtkToolItemGroup-expander-size: 11; -GtkWidget-text-handle-width: 20; -GtkWidget-text-handle-height: 24; -GtkDialog-button-spacing: 4; -GtkDialog-action-area-border: 0; outline-color: alpha(currentColor,0.3); outline-style: dashed; outline-offset: -3px; outline-width: 1px; -gtk-outline-radius: 3px; -gtk-secondary-caret-color: #3584e4; }

/*************** Base States * */
.background { color: #2e3436; background-color: #f6f5f4; }

.background:backdrop { color: #929595; background-color: #f6f5f4; text-shadow: none; -gtk-icon-shadow: none; }

/* These wildcard seems unavoidable, need to investigate. Wildcards are bad and troublesome, use them with care, or better, just don't. Everytime a wildcard is used a kitten dies, painfully. */
*:disabled { -gtk-icon-effect: dim; }

.gtkstyle-fallback { color: #2e3436; background-color: #f6f5f4; }

.gtkstyle-fallback:hover { color: #2e3436; background-color: white; }

.gtkstyle-fallback:active { color: #2e3436; background-color: #dfdcd8; }

.gtkstyle-fallback:disabled { color: #929595; background-color: #faf9f8; }

.gtkstyle-fallback:selected { color: #ffffff; background-color: #3584e4; }

.view, iconview, .view text, iconview text, textview text { color: black; background-color: #ffffff; }

.view:backdrop, iconview:backdrop, .view text:backdrop, iconview text:backdrop, textview text:backdrop { color: #323232; background-color: #fcfcfc; }

.view:backdrop:disabled, iconview:backdrop:disabled, .view text:backdrop:disabled, iconview text:backdrop:disabled, textview text:backdrop:disabled { color: #d4cfca; }

.view:disabled, iconview:disabled, .view text:disabled, iconview text:disabled, textview text:disabled { color: #929595; background-color: #faf9f8; }

.view:selected:focus, iconview:selected:focus, .view:selected, iconview:selected, .view text:selected:focus, iconview text:selected:focus, textview text:selected:focus, .view text:selected, iconview text:selected, textview text:selected { border-radius: 3px; }

textview border { background-color: #fbfafa; }

.rubberband, rubberband, .content-view rubberband, .content-view .rubberband, treeview.view rubberband, flowbox rubberband { border: 1px solid #1b6acb; background-color: rgba(27, 106, 203, 0.2); }

flowbox flowboxchild { padding: 3px; }

flowbox flowboxchild:selected { outline-offset: -2px; }

.content-view .tile { margin: 2px; background-color: transparent; border-radius: 0; padding: 0; }

.content-view .tile:backdrop { background-color: transparent; }

.content-view .tile:active, .content-view .tile:selected { background-color: transparent; }

.content-view .tile:disabled { background-color: transparent; }

label { caret-color: currentColor; }

label selection { background-color: #3584e4; color: #ffffff; }

label:disabled { color: #929595; }

button label:disabled { color: inherit; }

label:disabled:backdrop { color: #d4cfca; }

button label:disabled:backdrop { color: inherit; }

.dim-label, .titlebar:not(headerbar) .subtitle, headerbar .subtitle, label.separator { opacity: 0.55; text-shadow: none; }

assistant .sidebar { background-color: #ffffff; border-top: 1px solid #cdc7c2; }

assistant .sidebar:backdrop { background-color: #fcfcfc; border-color: #d5d0cc; }

assistant.csd .sidebar { border-top-style: none; }

assistant .sidebar label { padding: 6px 12px; }

assistant .sidebar label.highlight { background-color: #cecece; }

.osd .scale-popup, .app-notification, .app-notification.frame, .csd popover.background.osd, popover.background.osd, .csd popover.background.touch-selection, .csd popover.background.magnifier, popover.background.touch-selection, popover.background.magnifier, .osd { color: #eeeeec; border: none; background-color: rgba(53, 53, 53, 0.9); background-clip: padding-box; text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; }

.osd .scale-popup:backdrop, .app-notification:backdrop, popover.background.osd:backdrop, popover.background.touch-selection:backdrop, popover.background.magnifier:backdrop, .osd:backdrop { text-shadow: none; -gtk-icon-shadow: none; }

/********************* Spinner Animation * */
@keyframes spin { to { -gtk-icon-transform: rotate(1turn); } }

spinner { background: none; opacity: 0; -gtk-icon-source: -gtk-icontheme("process-working-symbolic"); }

spinner:backdrop { color: #929595; }

spinner:checked { opacity: 1; animation: spin 1s linear infinite; }

spinner:checked:disabled { opacity: 0.5; }

/********************** General Typography * */
.large-title { font-weight: 300; font-size: 24pt; letter-spacing: 0.2rem; }

.title-1 { font-weight: 800; font-size: 20pt; }

.title-2 { font-weight: 800; font-size: 15pt; }

.title-3 { font-weight: 700; font-size: 15pt; }

.title-4 { font-weight: 700; font-size: 13pt; }

.heading { font-weight: 700; font-size: 11pt; }

.body { font-weight: 400; font-size: 11pt; }

.caption-heading { font-weight: 700; font-size: 9pt; }

.caption { font-weight: 400; font-size: 9pt; }

/**************** Text Entries * */
spinbutton:not(.vertical), entry { min-height: 32px; padding-left: 8px; padding-right: 8px; border: 1px solid; border-radius: 5px; transition: all 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); color: black; border-color: #cdc7c2; background-color: #ffffff; box-shadow: inset 0 0 0 1px rgba(53, 132, 228, 0); }

spinbutton:not(.vertical) image.left, entry image.left { margin-right: 6px; }

spinbutton:not(.vertical) image.right, entry image.right { margin-left: 6px; }

spinbutton.flat:not(.vertical), entry.flat:focus, entry.flat:backdrop, entry.flat:disabled, entry.flat { min-height: 0; padding: 2px; background-color: transparent; border-color: transparent; border-radius: 0; }

spinbutton:focus:not(.vertical), entry:focus { box-shadow: inset 0 0 0 1px #3584e4; border-color: #3584e4; }

spinbutton:disabled:not(.vertical), entry:disabled { color: #929595; border-color: #cdc7c2; background-color: #faf9f8; box-shadow: none; }

spinbutton:backdrop:not(.vertical), entry:backdrop { color: #323232; border-color: #d5d0cc; background-color: #fcfcfc; box-shadow: none; transition: 200ms ease-out; }

spinbutton:backdrop:disabled:not(.vertical), entry:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-color: #faf9f8; box-shadow: none; }

spinbutton.error:not(.vertical), entry.error { color: #cc0000; border-color: #cc0000; }

spinbutton.error:focus:not(.vertical), entry.error:focus { box-shadow: inset 0 0 0 1px #cc0000; border-color: #cc0000; }

spinbutton.error:not(.vertical) selection, entry.error selection { background-color: #cc0000; }

spinbutton.warning:not(.vertical), entry.warning { color: #f57900; border-color: #f57900; }

spinbutton.warning:focus:not(.vertical), entry.warning:focus { box-shadow: inset 0 0 0 1px #f57900; border-color: #f57900; }

spinbutton.warning:not(.vertical) selection, entry.warning selection { background-color: #f57900; }

spinbutton:not(.vertical) image, entry image { color: #585d5e; }

spinbutton:not(.vertical) image:hover, entry image:hover { color: #2e3436; }

spinbutton:not(.vertical) image:active, entry image:active { color: #3584e4; }

spinbutton:not(.vertical) image:backdrop, entry image:backdrop { color: #a7aaaa; }

spinbutton:drop(active):not(.vertical), entry:drop(active):focus, entry:drop(active) { border-color: #4e9a06; box-shadow: inset 0 0 0 1px #4e9a06; }

.osd spinbutton:not(.vertical), .osd entry { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: rgba(0, 0, 0, 0.5); background-clip: padding-box; box-shadow: none; text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; }

.osd spinbutton:focus:not(.vertical), .osd entry:focus { color: white; border-color: #3584e4; background-color: rgba(0, 0, 0, 0.5); background-clip: padding-box; box-shadow: inset 0 0 0 1px #3584e4; text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; }

.osd spinbutton:backdrop:not(.vertical), .osd entry:backdrop { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: rgba(0, 0, 0, 0.5); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.osd spinbutton:disabled:not(.vertical), .osd entry:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: rgba(71, 71, 71, 0.5); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

spinbutton:not(.vertical) progress, entry progress { margin: 2px -6px; background-color: transparent; background-image: none; border-radius: 0; border-width: 0 0 2px; border-color: #3584e4; border-style: solid; box-shadow: none; }

spinbutton:not(.vertical) progress:backdrop, entry progress:backdrop { background-color: transparent; }

.linked:not(.vertical) > spinbutton:focus:not(.vertical) + spinbutton:not(.vertical), .linked:not(.vertical) > spinbutton:focus:not(.vertical) + button, .linked:not(.vertical) > spinbutton:focus:not(.vertical) + combobox > box > button.combo, .linked:not(.vertical) > spinbutton:focus:not(.vertical) + entry, .linked:not(.vertical) > entry:focus + button, .linked:not(.vertical) > entry:focus + combobox > box > button.combo, .linked:not(.vertical) > entry:focus + spinbutton:not(.vertical), .linked:not(.vertical) > entry:focus + entry { border-left-color: #3584e4; }

.linked:not(.vertical) > spinbutton.error:focus:not(.vertical) + spinbutton:not(.vertical), .linked:not(.vertical) > spinbutton.error:focus:not(.vertical) + button, .linked:not(.vertical) > spinbutton.error:focus:not(.vertical) + combobox > box > button.combo, .linked:not(.vertical) > spinbutton.error:focus:not(.vertical) + entry, .linked:not(.vertical) > entry.error:focus + button, .linked:not(.vertical) > entry.error:focus + combobox > box > button.combo, .linked:not(.vertical) > entry.error:focus + spinbutton:not(.vertical), .linked:not(.vertical) > entry.error:focus + entry { border-left-color: #cc0000; }

.linked:not(.vertical) > spinbutton:drop(active):not(.vertical) + spinbutton:not(.vertical), .linked:not(.vertical) > spinbutton:drop(active):not(.vertical) + button, .linked:not(.vertical) > spinbutton:drop(active):not(.vertical) + combobox > box > button.combo, .linked:not(.vertical) > spinbutton:drop(active):not(.vertical) + entry, .linked:not(.vertical) > entry:drop(active) + button, .linked:not(.vertical) > entry:drop(active) + combobox > box > button.combo, .linked:not(.vertical) > entry:drop(active) + spinbutton:not(.vertical), .linked:not(.vertical) > entry:drop(active) + entry { border-left-color: #4e9a06; }

.linked.vertical > spinbutton:not(:disabled):not(.vertical) + entry:not(:disabled), .linked.vertical > spinbutton:not(:disabled):not(.vertical) + spinbutton:not(:disabled):not(.vertical), .linked.vertical > entry:not(:disabled) + entry:not(:disabled), .linked.vertical > entry:not(:disabled) + spinbutton:not(:disabled):not(.vertical) { border-top-color: #f0eeed; }

.linked.vertical > spinbutton:not(:disabled):not(.vertical) + entry:not(:disabled):backdrop, .linked.vertical > spinbutton:not(:disabled):not(.vertical) + spinbutton:not(:disabled):backdrop:not(.vertical), .linked.vertical > entry:not(:disabled) + entry:not(:disabled):backdrop, .linked.vertical > entry:not(:disabled) + spinbutton:not(:disabled):backdrop:not(.vertical) { border-top-color: #f1efee; }

.linked.vertical > spinbutton:disabled:not(.vertical) + spinbutton:disabled:not(.vertical), .linked.vertical > spinbutton:disabled:not(.vertical) + entry:disabled, .linked.vertical > entry:disabled + spinbutton:disabled:not(.vertical), .linked.vertical > entry:disabled + entry:disabled { border-top-color: #f0eeed; }

.linked.vertical > spinbutton:not(.vertical) + spinbutton:focus:not(:only-child):not(.vertical), .linked.vertical > spinbutton:not(.vertical) + entry:focus:not(:only-child), .linked.vertical > entry + spinbutton:focus:not(:only-child):not(.vertical), .linked.vertical > entry + entry:focus:not(:only-child) { border-top-color: #3584e4; }

.linked.vertical > spinbutton:not(.vertical) + spinbutton.error:focus:not(:only-child):not(.vertical), .linked.vertical > spinbutton:not(.vertical) + entry.error:focus:not(:only-child), .linked.vertical > entry + spinbutton.error:focus:not(:only-child):not(.vertical), .linked.vertical > entry + entry.error:focus:not(:only-child) { border-top-color: #cc0000; }

.linked.vertical > spinbutton:not(.vertical) + spinbutton:drop(active):not(:only-child):not(.vertical), .linked.vertical > spinbutton:not(.vertical) + entry:drop(active):not(:only-child), .linked.vertical > entry + spinbutton:drop(active):not(:only-child):not(.vertical), .linked.vertical > entry + entry:drop(active):not(:only-child) { border-top-color: #4e9a06; }

.linked.vertical > spinbutton:focus:not(:only-child):not(.vertical) + spinbutton:not(.vertical), .linked.vertical > spinbutton:focus:not(:only-child):not(.vertical) + entry, .linked.vertical > spinbutton:focus:not(:only-child):not(.vertical) + button, .linked.vertical > spinbutton:focus:not(:only-child):not(.vertical) + combobox > box > button.combo, .linked.vertical > entry:focus:not(:only-child) + spinbutton:not(.vertical), .linked.vertical > entry:focus:not(:only-child) + entry, .linked.vertical > entry:focus:not(:only-child) + button, .linked.vertical > entry:focus:not(:only-child) + combobox > box > button.combo { border-top-color: #3584e4; }

.linked.vertical > spinbutton.error:focus:not(:only-child):not(.vertical) + spinbutton:not(.vertical), .linked.vertical > spinbutton.error:focus:not(:only-child):not(.vertical) + entry, .linked.vertical > spinbutton.error:focus:not(:only-child):not(.vertical) + button, .linked.vertical > spinbutton.error:focus:not(:only-child):not(.vertical) + combobox > box > button.combo, .linked.vertical > entry.error:focus:not(:only-child) + spinbutton:not(.vertical), .linked.vertical > entry.error:focus:not(:only-child) + entry, .linked.vertical > entry.error:focus:not(:only-child) + button, .linked.vertical > entry.error:focus:not(:only-child) + combobox > box > button.combo { border-top-color: #cc0000; }

.linked.vertical > spinbutton:drop(active):not(:only-child):not(.vertical) + spinbutton:not(.vertical), .linked.vertical > spinbutton:drop(active):not(:only-child):not(.vertical) + entry, .linked.vertical > spinbutton:drop(active):not(:only-child):not(.vertical) + button, .linked.vertical > spinbutton:drop(active):not(:only-child):not(.vertical) + combobox > box > button.combo, .linked.vertical > entry:drop(active):not(:only-child) + spinbutton:not(.vertical), .linked.vertical > entry:drop(active):not(:only-child) + entry, .linked.vertical > entry:drop(active):not(:only-child) + button, .linked.vertical > entry:drop(active):not(:only-child) + combobox > box > button.combo { border-top-color: #4e9a06; }

treeview entry:focus:dir(rtl), treeview entry:focus:dir(ltr) { background-color: #ffffff; transition-property: color, background; }

treeview entry.flat, treeview entry { border-radius: 0; background-image: none; background-color: #ffffff; }

treeview entry.flat:focus, treeview entry:focus { border-color: #3584e4; }

.entry-tag { padding: 5px; margin-top: 2px; margin-bottom: 2px; border-style: none; color: #ffffff; background-color: #3584e4; }

:dir(ltr) .entry-tag { margin-left: 8px; margin-right: -5px; }

:dir(rtl) .entry-tag { margin-left: -5px; margin-right: 8px; }

.entry-tag:hover { background-color: #629fea; }

:backdrop .entry-tag { color: #fcfcfc; background-color: #3584e4; }

.entry-tag.button { background-color: transparent; color: rgba(255, 255, 255, 0.7); }

:not(:backdrop) .entry-tag.button:hover { border: 1px solid #3584e4; color: #ffffff; }

:not(:backdrop) .entry-tag.button:active { background-color: #3584e4; color: rgba(255, 255, 255, 0.7); }

/*********** Buttons * */
@keyframes needs_attention { from { background-image: -gtk-gradient(radial, center center, 0, center center, 0.01, to(#3584e4), to(transparent)); }
  to { background-image: -gtk-gradient(radial, center center, 0, center center, 0.5, to(#3584e4), to(transparent)); } }

button.titlebutton, notebook > header > tabs > arrow, button { min-height: 24px; min-width: 16px; padding: 4px 9px; border: 1px solid; border-radius: 5px; transition: all 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; background-image: linear-gradient(to top, #edebe9 2px, #f6f5f4); text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); }

button.titlebutton, button.sidebar-button, notebook > header > tabs > arrow, notebook > header > tabs > arrow.flat, button.flat { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; transition: none; }

button.titlebutton:hover, button.sidebar-button:hover, notebook > header > tabs > arrow:hover, button.flat:hover { transition: all 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); transition-duration: 500ms; }

button.titlebutton:hover:active, button.sidebar-button:hover:active, notebook > header > tabs > arrow:hover:active, button.flat:hover:active { transition: all 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); }

notebook > header > tabs > arrow:hover, button:hover { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px); -gtk-icon-effect: highlight; }

notebook > header > tabs > arrow:active, notebook > header > tabs > arrow:checked, button:active, button:checked { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; background-image: image(#d6d1cd); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; transition-duration: 50ms; }

notebook > header > tabs > arrow:backdrop, button.flat:backdrop, button:backdrop { color: #929595; border-color: #d5d0cc; background-image: image(#f6f5f4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); transition: 200ms ease-out; -gtk-icon-effect: none; }

notebook > header > tabs > arrow:backdrop:active, notebook > header > tabs > arrow:backdrop:checked, button.flat:backdrop:active, button.flat:backdrop:checked, button:backdrop:active, button:backdrop:checked { color: #929595; border-color: #d5d0cc; background-image: image(#e4e4e0); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

notebook > header > tabs > arrow:backdrop:disabled, button.flat:backdrop:disabled, button:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

notebook > header > tabs > arrow:backdrop:disabled:active, notebook > header > tabs > arrow:backdrop:disabled:checked, button.flat:backdrop:disabled:active, button.flat:backdrop:disabled:checked, button:backdrop:disabled:active, button:backdrop:disabled:checked { color: #d4cfca; border-color: #d5d0cc; background-image: image(#e4e4e0); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.titlebutton:backdrop, button.sidebar-button:backdrop, notebook > header > tabs > arrow:backdrop, button.titlebutton:disabled, button.sidebar-button:disabled, notebook > header > tabs > arrow:disabled, button.flat:backdrop, button.flat:disabled, button.flat:backdrop:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

notebook > header > tabs > arrow:disabled, button:disabled { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

notebook > header > tabs > arrow:disabled:active, notebook > header > tabs > arrow:disabled:checked, button:disabled:active, button:disabled:checked { color: #929595; border-color: #cdc7c2; background-image: image(#e4e4e0); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

notebook > header > tabs > arrow.image-button, button.image-button { min-width: 24px; padding-left: 5px; padding-right: 5px; }

notebook > header > tabs > arrow.text-button, button.text-button { padding-left: 16px; padding-right: 16px; }

notebook > header > tabs > arrow.text-button.image-button, button.text-button.image-button { padding-left: 8px; padding-right: 8px; }

notebook > header > tabs > arrow.text-button.image-button label, button.text-button.image-button label { padding-left: 8px; padding-right: 8px; }

combobox:drop(active) button.combo, notebook > header > tabs > arrow:drop(active), button:drop(active) { color: #4e9a06; border-color: #4e9a06; box-shadow: inset 0 0 0 1px #4e9a06; }

row:selected button { border-color: #185fb4; }

row:selected button.sidebar-button:not(:active):not(:checked):not(:hover):not(disabled), row:selected button.flat:not(:active):not(:checked):not(:hover):not(disabled) { color: #ffffff; border-color: transparent; }

row:selected button.sidebar-button:not(:active):not(:checked):not(:hover):not(disabled):backdrop, row:selected button.flat:not(:active):not(:checked):not(:hover):not(disabled):backdrop { color: #fcfcfc; }

button.osd { min-width: 26px; min-height: 32px; color: #eeeeec; border-radius: 5px; color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); border: none; box-shadow: none; }

button.osd.image-button { min-width: 34px; }

button.osd:hover { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(83, 83, 83, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); border: none; box-shadow: none; }

button.osd:active, button.osd:checked { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(0, 0, 0, 0.7)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); border: none; box-shadow: none; }

button.osd:disabled:backdrop, button.osd:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; border: none; }

button.osd:backdrop { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; border: none; }

.app-notification button, .app-notification.frame button, .csd popover.background.touch-selection button, .csd popover.background.magnifier button, popover.background.touch-selection button, popover.background.magnifier button, .osd button { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.app-notification button:hover, popover.background.touch-selection button:hover, popover.background.magnifier button:hover, .osd button:hover { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(83, 83, 83, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.app-notification button:active, popover.background.touch-selection button:active, popover.background.magnifier button:active, .app-notification button:checked, popover.background.touch-selection button:checked, popover.background.magnifier button:checked, .osd button:active:backdrop, .osd button:active, .osd button:checked:backdrop, .osd button:checked { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(0, 0, 0, 0.7)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); }

.app-notification button:disabled, popover.background.touch-selection button:disabled, popover.background.magnifier button:disabled, .osd button:disabled:backdrop, .osd button:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.app-notification button:backdrop, popover.background.touch-selection button:backdrop, popover.background.magnifier button:backdrop, .osd button:backdrop { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.app-notification button.flat, popover.background.touch-selection button.flat, popover.background.magnifier button.flat, .osd button.flat { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; box-shadow: none; text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; }

.app-notification button.flat:hover, popover.background.touch-selection button.flat:hover, popover.background.magnifier button.flat:hover, .osd button.flat:hover { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(83, 83, 83, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.app-notification button.flat:disabled, popover.background.touch-selection button.flat:disabled, popover.background.magnifier button.flat:disabled, .osd button.flat:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; background-image: none; border-color: transparent; box-shadow: none; }

.app-notification button.flat:backdrop, popover.background.touch-selection button.flat:backdrop, popover.background.magnifier button.flat:backdrop, .osd button.flat:backdrop { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

.app-notification button.flat:active, popover.background.touch-selection button.flat:active, popover.background.magnifier button.flat:active, .app-notification button.flat:checked, popover.background.touch-selection button.flat:checked, popover.background.magnifier button.flat:checked, .osd button.flat:active, .osd button.flat:checked { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(0, 0, 0, 0.7)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); }

button.suggested-action { color: white; outline-color: rgba(255, 255, 255, 0.3); border-color: #1b6acb; border-bottom-color: #15539e; background-image: linear-gradient(to top, #2379e2 2px, #3584e4); text-shadow: 0 -1px rgba(0, 0, 0, 0.559216); -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.559216); box-shadow: inset 0 1px rgba(255, 255, 255, 0.2), 0 1px 2px rgba(0, 0, 0, 0.07); }

button.suggested-action.flat { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #3584e4; }

button.suggested-action:hover { color: white; outline-color: rgba(255, 255, 255, 0.3); border-color: #1b6acb; border-bottom-color: #15539e; text-shadow: 0 -1px rgba(0, 0, 0, 0.511216); -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.511216); box-shadow: inset 0 1px rgba(255, 255, 255, 0.2), 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #3584e4, #3987e5 1px); }

button.suggested-action:active, button.suggested-action:checked { color: white; outline-color: rgba(255, 255, 255, 0.3); border-color: #1b6acb; background-image: image(#1961b9); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

button.suggested-action:backdrop, button.suggested-action.flat:backdrop { color: #d7e6fa; border-color: #3584e4; background-image: image(#3584e4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.suggested-action:backdrop:active, button.suggested-action:backdrop:checked, button.suggested-action.flat:backdrop:active, button.suggested-action.flat:backdrop:checked { color: #d5e6f9; border-color: #2f80e3; background-image: image(#2f80e3); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.suggested-action:backdrop:disabled, button.suggested-action.flat:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.suggested-action:backdrop:disabled:active, button.suggested-action:backdrop:disabled:checked, button.suggested-action.flat:backdrop:disabled:active, button.suggested-action.flat:backdrop:disabled:checked { color: #78aced; border-color: #2f80e3; background-image: image(#2f80e3); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.suggested-action.flat:backdrop, button.suggested-action.flat:disabled, button.suggested-action.flat:backdrop:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: rgba(53, 132, 228, 0.8); }

button.suggested-action:disabled { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.suggested-action:disabled:active, button.suggested-action:disabled:checked { color: #acccf4; border-color: #1b6acb; background-image: image(#2f80e3); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.osd button.suggested-action { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 132, 228, 0.5)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd button.suggested-action:hover { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 132, 228, 0.7)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd button.suggested-action:active:backdrop, .osd button.suggested-action:active, .osd button.suggested-action:checked:backdrop, .osd button.suggested-action:checked { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(#3584e4); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); }

.osd button.suggested-action:disabled:backdrop, .osd button.suggested-action:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.osd button.suggested-action:backdrop { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 132, 228, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

button.destructive-action { color: white; outline-color: rgba(255, 255, 255, 0.3); border-color: #b2161d; border-bottom-color: #851015; background-image: linear-gradient(to top, #ce1921 2px, #e01b24); text-shadow: 0 -1px rgba(0, 0, 0, 0.606275); -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.606275); box-shadow: inset 0 1px rgba(255, 255, 255, 0.1), 0 1px 2px rgba(0, 0, 0, 0.07); }

button.destructive-action.flat { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #e01b24; }

button.destructive-action:hover { color: white; outline-color: rgba(255, 255, 255, 0.3); border-color: #b2161d; border-bottom-color: #851015; text-shadow: 0 -1px rgba(0, 0, 0, 0.558275); -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.558275); box-shadow: inset 0 1px rgba(255, 255, 255, 0.2), 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #e01b24, #e41c26 1px); }

button.destructive-action:active, button.destructive-action:checked { color: white; outline-color: rgba(255, 255, 255, 0.3); border-color: #b2161d; background-image: image(#a0131a); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

button.destructive-action:backdrop, button.destructive-action.flat:backdrop { color: #f9d1d3; border-color: #e01b24; background-image: image(#e01b24); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.destructive-action:backdrop:active, button.destructive-action:backdrop:checked, button.destructive-action.flat:backdrop:active, button.destructive-action.flat:backdrop:checked { color: #f8d2d4; border-color: #dc1d27; background-image: image(#dc1d27); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.destructive-action:backdrop:disabled, button.destructive-action.flat:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.destructive-action:backdrop:disabled:active, button.destructive-action:backdrop:disabled:checked, button.destructive-action.flat:backdrop:disabled:active, button.destructive-action.flat:backdrop:disabled:checked { color: #e86c72; border-color: #dc1d27; background-image: image(#dc1d27); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.destructive-action.flat:backdrop, button.destructive-action.flat:disabled, button.destructive-action.flat:backdrop:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: rgba(224, 27, 36, 0.8); }

button.destructive-action:disabled { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

button.destructive-action:disabled:active, button.destructive-action:disabled:checked { color: #f1a5a8; border-color: #b2161d; background-image: image(#dc1d27); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.osd button.destructive-action { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(224, 27, 36, 0.5)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd button.destructive-action:hover { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(224, 27, 36, 0.7)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd button.destructive-action:active:backdrop, .osd button.destructive-action:active, .osd button.destructive-action:checked:backdrop, .osd button.destructive-action:checked { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(#e01b24); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); }

.osd button.destructive-action:disabled:backdrop, .osd button.destructive-action:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.osd button.destructive-action:backdrop { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(224, 27, 36, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.stack-switcher > button { outline-offset: -3px; }

.stack-switcher > button > label { padding-left: 6px; padding-right: 6px; }

.stack-switcher > button > image { padding-left: 6px; padding-right: 6px; padding-top: 3px; padding-bottom: 3px; }

.stack-switcher > button.text-button { padding-left: 10px; padding-right: 10px; }

.stack-switcher > button.image-button { padding-left: 2px; padding-right: 2px; }

.stack-switcher > button.needs-attention:active > label, .stack-switcher > button.needs-attention:active > image, .stack-switcher > button.needs-attention:checked > label, .stack-switcher > button.needs-attention:checked > image { animation: none; background-image: none; }

button.font separator, button.file separator { background-color: transparent; }

button.font > box > box > label { font-weight: bold; }

.inline-toolbar button, .inline-toolbar button:backdrop { border-radius: 2px; border-width: 1px; }

.primary-toolbar button { -gtk-icon-shadow: none; }

button.circular { border-radius: 9999px; -gtk-outline-radius: 9999px; padding: 4px; /* circles instead of ellipses */ background-origin: padding-box, border-box; background-clip: padding-box, border-box; }

button.circular label { padding: 0; }

button.circular:not(.flat):not(.osd):not(:checked):not(:active):not(:disabled):not(:backdrop) { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; background-image: linear-gradient(to top, #edebe9 2px, #f6f5f4); text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #edebe9 2px, #f6f5f4), linear-gradient(to top, #bfb8b1 25%, #cdc7c2 50%); border-color: transparent; }

button.circular:hover:not(.osd):not(:checked):not(:active):not(:disabled):not(:backdrop) { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px), linear-gradient(to top, #bfb8b1 25%, #cdc7c2 50%); border-color: transparent; }

stacksidebar row.needs-attention > label, .stack-switcher > button.needs-attention > label, .stack-switcher > button.needs-attention > image { animation: needs_attention 150ms ease-in; background-image: -gtk-gradient(radial, center center, 0, center center, 0.5, to(#3584e4), to(transparent)), -gtk-gradient(radial, center center, 0, center center, 0.5, to(rgba(255, 255, 255, 0.769231)), to(transparent)); background-size: 6px 6px, 6px 6px; background-repeat: no-repeat; background-position: right 3px, right 4px; }

stacksidebar row.needs-attention > label:backdrop, .stack-switcher > button.needs-attention > label:backdrop, .stack-switcher > button.needs-attention > image:backdrop { background-size: 6px 6px, 0 0; }

stacksidebar row.needs-attention > label:dir(rtl), .stack-switcher > button.needs-attention > label:dir(rtl), .stack-switcher > button.needs-attention > image:dir(rtl) { background-position: left 3px, left 4px; }

.inline-toolbar toolbutton > button { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; background-image: linear-gradient(to top, #edebe9 2px, #f6f5f4); text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); }

.inline-toolbar toolbutton > button:hover { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px); }

.inline-toolbar toolbutton > button:active, .inline-toolbar toolbutton > button:checked { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; background-image: image(#d6d1cd); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

.inline-toolbar toolbutton > button:disabled { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.inline-toolbar toolbutton > button:disabled:active, .inline-toolbar toolbutton > button:disabled:checked { color: #929595; border-color: #cdc7c2; background-image: image(#e4e4e0); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.inline-toolbar toolbutton > button:backdrop { color: #929595; border-color: #d5d0cc; background-image: image(#f6f5f4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.inline-toolbar toolbutton > button:backdrop:active, .inline-toolbar toolbutton > button:backdrop:checked { color: #929595; border-color: #d5d0cc; background-image: image(#e4e4e0); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.inline-toolbar toolbutton > button:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.inline-toolbar toolbutton > button:backdrop:disabled:active, .inline-toolbar toolbutton > button:backdrop:disabled:checked { color: #d4cfca; border-color: #d5d0cc; background-image: image(#e4e4e0); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.linked > combobox > box > button.combo:dir(ltr), .linked > combobox > box > button.combo:dir(rtl), filechooser .path-bar.linked > button, .linked:not(.vertical) > spinbutton:not(.vertical), .linked:not(.vertical) > entry, .inline-toolbar button, .inline-toolbar button:backdrop, .linked > button, .linked > button:hover, .linked > button:active, .linked > button:checked, .linked > button:backdrop, toolbar.inline-toolbar toolbutton > button.flat, toolbar.inline-toolbar toolbutton:backdrop > button.flat { border-radius: 0; border-right-style: none; -gtk-outline-radius: 0; }

.linked:not(.vertical) > combobox:first-child > box > button.combo, combobox.linked button:nth-child(2):dir(rtl), filechooser .path-bar.linked > button:dir(rtl):last-child, filechooser .path-bar.linked > button:dir(ltr):first-child, .linked:not(.vertical) > spinbutton:first-child:not(.vertical), .linked:not(.vertical) > entry:first-child, .inline-toolbar button:first-child, .inline-toolbar button:first-child:backdrop, .linked > button:first-child, .linked > button:first-child:hover, .linked > button:first-child:active, .linked > button:first-child:checked, .linked > button:first-child:backdrop, toolbar.inline-toolbar toolbutton:first-child > button.flat { border-top-left-radius: 5px; border-bottom-left-radius: 5px; border-top-right-radius: 0; border-bottom-right-radius: 0; border-right-style: none; -gtk-outline-bottom-left-radius: 5px; -gtk-outline-top-left-radius: 5px; -gtk-outline-top-right-radius: 0; -gtk-outline-bottom-right-radius: 0; }

.linked:not(.vertical) > combobox:last-child > box > button.combo, combobox.linked button:nth-child(2):dir(ltr), filechooser .path-bar.linked > button:dir(rtl):first-child, filechooser .path-bar.linked > button:dir(ltr):last-child, .linked:not(.vertical) > spinbutton:last-child:not(.vertical), .linked:not(.vertical) > entry:last-child, .inline-toolbar button:last-child, .inline-toolbar button:last-child:backdrop, .linked > button:last-child, .linked > button:last-child:hover, .linked > button:last-child:active, .linked > button:last-child:checked, .linked > button:last-child:backdrop, toolbar.inline-toolbar toolbutton:last-child > button.flat { border-top-left-radius: 0; border-bottom-left-radius: 0; border-top-right-radius: 5px; border-bottom-right-radius: 5px; border-right-style: solid; -gtk-outline-bottom-right-radius: 5px; -gtk-outline-top-right-radius: 5px; -gtk-outline-bottom-left-radius: 0; -gtk-outline-top-left-radius: 0; }

.linked:not(.vertical) > combobox:only-child > box > button.combo, filechooser .path-bar.linked > button:only-child, .linked:not(.vertical) > spinbutton:only-child:not(.vertical), .linked:not(.vertical) > entry:only-child, .inline-toolbar button:only-child, .inline-toolbar button:only-child:backdrop, .linked > button:only-child, .linked > button:only-child:hover, .linked > button:only-child:active, .linked > button:only-child:checked, .linked > button:only-child:backdrop, toolbar.inline-toolbar toolbutton:only-child > button.flat { border-radius: 5px; border-style: solid; -gtk-outline-radius: 5px; }

.linked.vertical > combobox > box > button.combo, .linked.vertical > spinbutton:not(.vertical), .linked.vertical > entry, .linked.vertical > button, .linked.vertical > button:hover, .linked.vertical > button:active, .linked.vertical > button:checked, .linked.vertical > button:backdrop { border-style: solid solid none solid; border-radius: 0; }

.linked.vertical > combobox:first-child > box > button.combo, .linked.vertical > spinbutton:first-child:not(.vertical), .linked.vertical > entry:first-child, .linked.vertical > button:first-child, .linked.vertical > button:first-child:hover, .linked.vertical > button:first-child:active, .linked.vertical > button:first-child:checked, .linked.vertical > button:first-child:backdrop { border-top-left-radius: 5px; border-top-right-radius: 5px; }

.linked.vertical > combobox:last-child > box > button.combo, .linked.vertical > spinbutton:last-child:not(.vertical), .linked.vertical > entry:last-child, .linked.vertical > button:last-child, .linked.vertical > button:last-child:hover, .linked.vertical > button:last-child:active, .linked.vertical > button:last-child:checked, .linked.vertical > button:last-child:backdrop { border-bottom-left-radius: 5px; border-bottom-right-radius: 5px; border-style: solid; }

.linked.vertical > combobox:only-child > box > button.combo, .linked.vertical > spinbutton:only-child:not(.vertical), .linked.vertical > entry:only-child, .linked.vertical > button:only-child, .linked.vertical > button:only-child:hover, .linked.vertical > button:only-child:active, .linked.vertical > button:only-child:checked, .linked.vertical > button:only-child:backdrop { border-radius: 5px; border-style: solid; }

.scale-popup button:backdrop:hover, .scale-popup button:backdrop:disabled, .scale-popup button:backdrop, .scale-popup button:hover, calendar.button, button:link:hover, button:link:active, button:link:checked, button:visited:hover, button:visited:active, button:visited:checked, button:link, button:visited, list row button.image-button:not(.flat), modelbutton.flat:backdrop, modelbutton.flat:backdrop:hover, modelbutton.flat, .menuitem.button.flat { background-color: transparent; background-image: none; border-color: transparent; box-shadow: inset 0 1px rgba(255, 255, 255, 0), 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

/* menu buttons */
modelbutton.flat, .menuitem.button.flat { min-height: 26px; padding-left: 5px; padding-right: 5px; border-radius: 5px; outline-offset: -2px; }

modelbutton.flat:hover, .menuitem.button.flat:hover { background-color: white; }

modelbutton.flat arrow { background: none; }

modelbutton.flat arrow:hover { background: none; }

modelbutton.flat arrow.left { -gtk-icon-source: -gtk-icontheme("pan-start-symbolic"); }

modelbutton.flat arrow.right { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic"); }

button.color { padding: 4px; }

button.color colorswatch:only-child { box-shadow: 0 1px rgba(255, 255, 255, 0.769231); }

button.color colorswatch:only-child, button.color colorswatch:only-child overlay { border-radius: 0; }

.osd button.color colorswatch:only-child { box-shadow: none; }

.osd button.color:disabled colorswatch:only-child, .osd button.color:backdrop colorswatch:only-child, .osd button.color:active colorswatch:only-child, .osd button.color:checked colorswatch:only-child, button.color:disabled colorswatch:only-child, button.color:backdrop colorswatch:only-child, button.color:active colorswatch:only-child, button.color:checked colorswatch:only-child { box-shadow: none; }

/* list buttons */
/* tone down as per new designs, see issue #1473, #1748 */
list row button.image-button:not(.flat) { border: 1px solid rgba(205, 199, 194, 0.5); }

list row button.image-button:not(.flat):hover { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px); }

list row button.image-button:not(.flat):active, list row button.image-button:not(.flat):checked { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; background-image: image(#d6d1cd); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

/********* Links * */
button:link > label, button:visited > label, button:link, button:visited, *:link { color: #1b6acb; }

button:link > label:visited, button:visited > label:visited, button:visited, *:link:visited { color: #15539e; }

*:selected button:link > label:visited, *:selected button:visited > label:visited, *:selected button:visited, *:selected *:link:visited { color: #a1bad8; }

button:link > label:hover, button:visited > label:hover, button:hover:link, button:hover:visited, *:link:hover { color: #3584e4; }

*:selected button:link > label:hover, *:selected button:visited > label:hover, *:selected button:hover:link, *:selected button:hover:visited, *:selected *:link:hover { color: #ebf3fc; }

button:link > label:active, button:visited > label:active, button:active:link, button:active:visited, *:link:active { color: #1b6acb; }

*:selected button:link > label:active, *:selected button:visited > label:active, *:selected button:active:link, *:selected button:active:visited, *:selected *:link:active { color: #d1e1f5; }

button:link > label:disabled, button:visited > label:disabled, button:disabled:link, button:disabled:visited, *:link:disabled, *:link:disabled:backdrop { color: rgba(115, 115, 115, 0.8); }

button:link > label:backdrop, button:visited > label:backdrop, button:backdrop:link, button:backdrop:visited, *:link:backdrop:backdrop:hover, *:link:backdrop:backdrop:hover:selected, *:link:backdrop { color: rgba(27, 106, 203, 0.9); }

.selection-mode .titlebar:not(headerbar) .subtitle:link, .selection-mode.titlebar:not(headerbar) .subtitle:link, .selection-mode headerbar .subtitle:link, headerbar.selection-mode .subtitle:link, button:link > label:selected, button:visited > label:selected, button:selected:link, button:selected:visited, *:selected button:link > label, *:selected button:visited > label, *:selected button:link, *:selected button:visited, *:link:selected, *:selected *:link { color: #d1e1f5; }

button:link, button:visited { text-shadow: none; }

button:link:hover, button:link:active, button:link:checked, button:visited:hover, button:visited:active, button:visited:checked { text-shadow: none; }

button:link > label, button:visited > label { text-decoration-line: underline; }

/***************** GtkSpinButton * */
spinbutton { font-feature-settings: "tnum"; }

spinbutton:not(.vertical) { padding: 0; }

.osd spinbutton:not(.vertical) entry, spinbutton:not(.vertical) entry { min-width: 28px; margin: 0; background: none; background-color: transparent; border: none; border-radius: 0; box-shadow: none; }

spinbutton:not(.vertical) entry:backdrop:disabled { background-color: transparent; }

spinbutton:not(.vertical) button { min-height: 16px; margin: 0; padding-bottom: 0; padding-top: 0; color: #43484a; background-image: none; border-style: none none none solid; border-color: rgba(205, 199, 194, 0.3); border-radius: 0; box-shadow: none; }

spinbutton:not(.vertical) button:dir(rtl) { border-style: none solid none none; }

spinbutton:not(.vertical) button:hover { color: #2e3436; background-color: rgba(46, 52, 54, 0.05); }

spinbutton:not(.vertical) button:disabled { color: rgba(146, 149, 149, 0.3); background-color: transparent; }

spinbutton:not(.vertical) button:active { background-color: rgba(0, 0, 0, 0.1); box-shadow: inset 0 2px 3px -1px rgba(0, 0, 0, 0.2); }

spinbutton:not(.vertical) button:backdrop { color: #9d9f9f; background-color: transparent; border-color: rgba(213, 208, 204, 0.3); transition: 200ms ease-out; }

spinbutton:not(.vertical) button:backdrop:disabled { color: rgba(212, 207, 202, 0.3); background-color: transparent; background-image: none; border-style: none none none solid; }

spinbutton:not(.vertical) button:backdrop:disabled:dir(rtl) { border-style: none solid none none; }

spinbutton:not(.vertical) button:dir(ltr):last-child { border-radius: 0 5px 5px 0; }

spinbutton:not(.vertical) button:dir(rtl):first-child { border-radius: 5px 0 0 5px; }

.osd spinbutton:not(.vertical) button { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #eeeeec; border-style: none none none solid; border-color: rgba(0, 0, 0, 0.4); border-radius: 0; box-shadow: none; -gtk-icon-shadow: 0 1px black; }

.osd spinbutton:not(.vertical) button:dir(rtl) { border-style: none solid none none; }

.osd spinbutton:not(.vertical) button:hover { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #eeeeec; border-color: rgba(0, 0, 0, 0.5); background-color: rgba(238, 238, 236, 0.1); -gtk-icon-shadow: 0 1px black; box-shadow: none; }

.osd spinbutton:not(.vertical) button:backdrop { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #eeeeec; border-color: rgba(0, 0, 0, 0.5); -gtk-icon-shadow: none; box-shadow: none; }

.osd spinbutton:not(.vertical) button:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #919190; border-color: rgba(0, 0, 0, 0.5); -gtk-icon-shadow: none; box-shadow: none; }

.osd spinbutton:not(.vertical) button:dir(ltr):last-child { border-radius: 0 5px 5px 0; }

.osd spinbutton:not(.vertical) button:dir(rtl):first-child { border-radius: 5px 0 0 5px; }

spinbutton.vertical:disabled { color: #929595; }

spinbutton.vertical:backdrop:disabled { color: #d4cfca; }

spinbutton.vertical:drop(active) { border-color: transparent; box-shadow: none; }

spinbutton.vertical entry { min-height: 32px; min-width: 32px; padding: 0; border-radius: 0; }

spinbutton.vertical button { min-height: 32px; min-width: 32px; padding: 0; }

spinbutton.vertical button.up { border-radius: 5px 5px 0 0; border-style: solid solid none solid; }

spinbutton.vertical button.down { border-radius: 0 0 5px 5px; border-style: none solid solid solid; }

.osd spinbutton.vertical button:first-child { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd spinbutton.vertical button:first-child:hover { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(83, 83, 83, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd spinbutton.vertical button:first-child:active { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(0, 0, 0, 0.7)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); }

.osd spinbutton.vertical button:first-child:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.osd spinbutton.vertical button:first-child:backdrop { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

treeview spinbutton:not(.vertical) { min-height: 0; border-style: none; border-radius: 0; }

treeview spinbutton:not(.vertical) entry { min-height: 0; padding: 1px 2px; }

/************** ComboBoxes * */
combobox arrow { -gtk-icon-source: -gtk-icontheme("pan-down-symbolic"); min-height: 16px; min-width: 16px; }

combobox:drop(active) { box-shadow: none; }

/************ Toolbars * */
searchbar > revealer > box, .location-bar, .inline-toolbar, toolbar { -GtkWidget-window-dragging: true; padding: 4px; background-color: #f6f5f4; }

toolbar { padding: 4px 3px 3px 4px; }

.osd toolbar { background-color: transparent; }

toolbar.osd { padding: 13px; border: none; border-radius: 5px; background-color: rgba(53, 53, 53, 0.9); }

toolbar.osd.left, toolbar.osd.right, toolbar.osd.top, toolbar.osd.bottom { border-radius: 0; }

toolbar.horizontal separator { margin: 0 7px 1px 6px; }

toolbar.vertical separator { margin: 6px 1px 7px 0; }

toolbar:not(.inline-toolbar):not(.osd) > *:not(.toggle):not(.popup) > * { margin-right: 1px; margin-bottom: 1px; }

.inline-toolbar { padding: 3px; border-width: 0 1px 1px; border-radius: 0  0 5px 5px; }

searchbar > revealer > box, .location-bar { border-width: 0 0 1px; padding: 3px; }

searchbar > revealer > box { margin: -6px; padding: 6px; }

.inline-toolbar, searchbar > revealer > box, .location-bar { border-style: solid; border-color: #cdc7c2; background-color: #eae7e5; }

.inline-toolbar:backdrop, searchbar > revealer > box:backdrop, .location-bar:backdrop { border-color: #d5d0cc; background-color: #eae8e6; box-shadow: none; transition: 200ms ease-out; }

/*************** Header bars * */
.titlebar:not(headerbar), headerbar { padding: 0 6px; min-height: 46px; border-width: 0 0 1px; border-style: solid; border-color: #bfb8b1; border-radius: 0; background: #dfdcd8 linear-gradient(to top, #dad6d2, #e1dedb); box-shadow: inset 0 1px rgba(255, 255, 255, 0.8); /* Darken switchbuttons for headerbars. issue #1588 */ /* hide the close button separator */ }

.titlebar:backdrop:not(headerbar), headerbar:backdrop { border-color: #d5d0cc; background-color: #f6f5f4; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0.8); transition: 200ms ease-out; }

.titlebar:not(headerbar) .title, headerbar .title { padding-left: 12px; padding-right: 12px; font-weight: bold; }

.titlebar:not(headerbar) .subtitle, headerbar .subtitle { font-size: smaller; padding-left: 12px; padding-right: 12px; }

.titlebar:not(headerbar) stackswitcher button:checked, .titlebar:not(headerbar) button.toggle:checked, headerbar stackswitcher button:checked, headerbar button.toggle:checked { background: image(#cfcac4); border-color: #c6bfb9; border-top-color: #bab3ab; }

.titlebar:not(headerbar) stackswitcher button:checked:backdrop, .titlebar:not(headerbar) button.toggle:checked:backdrop, headerbar stackswitcher button:checked:backdrop, headerbar button.toggle:checked:backdrop { color: #929595; border-color: #d5d0cc; background-image: image(#e4e4e0); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.selection-mode .titlebar:not(headerbar), .selection-mode.titlebar:not(headerbar), .selection-mode headerbar, headerbar.selection-mode { color: #ffffff; border-color: #185fb4; text-shadow: 0 -1px rgba(0, 0, 0, 0.5); background: #3584e4 linear-gradient(to top, #2c7fe3, #3987e5); box-shadow: inset 0 1px rgba(134, 181, 239, 0.9); }

.selection-mode .titlebar:backdrop:not(headerbar), .selection-mode.titlebar:backdrop:not(headerbar), .selection-mode headerbar:backdrop, headerbar.selection-mode:backdrop { background-color: #3584e4; background-image: none; box-shadow: inset 0 1px rgba(154, 194, 242, 0.88); }

.selection-mode .titlebar:backdrop:not(headerbar) label, .selection-mode.titlebar:backdrop:not(headerbar) label, .selection-mode headerbar:backdrop label, headerbar.selection-mode:backdrop label { text-shadow: none; color: #ffffff; }

.selection-mode .titlebar:not(headerbar) button, .selection-mode.titlebar:not(headerbar) button, .selection-mode headerbar button, headerbar.selection-mode button { color: #ffffff; outline-color: rgba(255, 255, 255, 0.3); border-color: #1b6acb; border-bottom-color: #15539e; background-image: linear-gradient(to top, #2379e2 2px, #3584e4); text-shadow: 0 -1px rgba(0, 0, 0, 0.559216); -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.559216); box-shadow: inset 0 1px rgba(255, 255, 255, 0.2), 0 1px 2px rgba(0, 0, 0, 0.07); }

.selection-mode button.titlebutton, .selection-mode .titlebar:not(headerbar) button.flat, .selection-mode.titlebar:not(headerbar) button.flat, .selection-mode headerbar button.flat, headerbar.selection-mode button.flat { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

.selection-mode .titlebar:not(headerbar) button:hover, .selection-mode.titlebar:not(headerbar) button:hover, .selection-mode headerbar button:hover, headerbar.selection-mode button:hover { color: #ffffff; outline-color: rgba(255, 255, 255, 0.3); border-color: #1b6acb; border-bottom-color: #15539e; text-shadow: 0 -1px rgba(0, 0, 0, 0.511216); -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.511216); box-shadow: inset 0 1px rgba(255, 255, 255, 0.2), 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #3584e4, #3987e5 1px); }

.selection-mode .titlebar:not(headerbar) button:active, .selection-mode .titlebar:not(headerbar) button:checked, .selection-mode.titlebar:not(headerbar) button:active, .selection-mode.titlebar:not(headerbar) button:checked, .selection-mode headerbar button:active, .selection-mode headerbar button:checked, .selection-mode headerbar button.toggle:checked, .selection-mode headerbar button.toggle:active, headerbar.selection-mode button:active, headerbar.selection-mode button:checked, headerbar.selection-mode button.toggle:checked, headerbar.selection-mode button.toggle:active { color: #ffffff; outline-color: rgba(255, 255, 255, 0.3); border-color: #1b6acb; background-image: image(#1961b9); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

.selection-mode .titlebar:not(headerbar) button:backdrop, .selection-mode.titlebar:not(headerbar) button:backdrop, .selection-mode headerbar button.flat:backdrop, .selection-mode headerbar button:backdrop, headerbar.selection-mode button.flat:backdrop, headerbar.selection-mode button:backdrop { color: #d7e6fa; border-color: #3584e4; background-image: image(#3584e4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); -gtk-icon-effect: none; border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button:backdrop:active, .selection-mode .titlebar:not(headerbar) button:backdrop:checked, .selection-mode.titlebar:not(headerbar) button:backdrop:active, .selection-mode.titlebar:not(headerbar) button:backdrop:checked, .selection-mode headerbar button.flat:backdrop:active, .selection-mode headerbar button.flat:backdrop:checked, .selection-mode headerbar button:backdrop:active, .selection-mode headerbar button:backdrop:checked, headerbar.selection-mode button.flat:backdrop:active, headerbar.selection-mode button.flat:backdrop:checked, headerbar.selection-mode button:backdrop:active, headerbar.selection-mode button:backdrop:checked { color: #d5e6f9; border-color: #2f80e3; background-image: image(#2f80e3); box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button:backdrop:disabled, .selection-mode.titlebar:not(headerbar) button:backdrop:disabled, .selection-mode headerbar button.flat:backdrop:disabled, .selection-mode headerbar button:backdrop:disabled, headerbar.selection-mode button.flat:backdrop:disabled, headerbar.selection-mode button:backdrop:disabled { color: #8fbbf0; border-color: #5396e8; background-image: image(#5396e8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button:backdrop:disabled:active, .selection-mode .titlebar:not(headerbar) button:backdrop:disabled:checked, .selection-mode.titlebar:not(headerbar) button:backdrop:disabled:active, .selection-mode.titlebar:not(headerbar) button:backdrop:disabled:checked, .selection-mode headerbar button:backdrop:disabled:active, .selection-mode headerbar button:backdrop:disabled:checked, headerbar.selection-mode button:backdrop:disabled:active, headerbar.selection-mode button:backdrop:disabled:checked { color: #78aced; border-color: #2f80e3; background-image: image(#2f80e3); box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #185fb4; }

.selection-mode button.titlebutton:backdrop, .selection-mode button.titlebutton:disabled, .selection-mode .titlebar:not(headerbar) button.flat:backdrop, .selection-mode .titlebar:not(headerbar) button.flat:disabled, .selection-mode.titlebar:not(headerbar) button.flat:backdrop, .selection-mode.titlebar:not(headerbar) button.flat:disabled, .selection-mode headerbar button.flat:backdrop, .selection-mode headerbar button.flat:disabled, .selection-mode headerbar button.flat:backdrop:disabled, headerbar.selection-mode button.flat:backdrop, headerbar.selection-mode button.flat:disabled, headerbar.selection-mode button.flat:backdrop:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

.selection-mode .titlebar:not(headerbar) button:disabled, .selection-mode.titlebar:not(headerbar) button:disabled, .selection-mode headerbar button:disabled, headerbar.selection-mode button:disabled { color: #a9cbf4; border-color: #1b6acb; background-image: image(#5396e8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.selection-mode .titlebar:not(headerbar) button:disabled:active, .selection-mode .titlebar:not(headerbar) button:disabled:checked, .selection-mode.titlebar:not(headerbar) button:disabled:active, .selection-mode.titlebar:not(headerbar) button:disabled:checked, .selection-mode headerbar button:disabled:active, .selection-mode headerbar button:disabled:checked, headerbar.selection-mode button:disabled:active, headerbar.selection-mode button:disabled:checked { color: #acccf4; border-color: #1b6acb; background-image: image(#2f80e3); box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

.selection-mode .titlebar:not(headerbar) button.suggested-action, .selection-mode.titlebar:not(headerbar) button.suggested-action, .selection-mode headerbar button.suggested-action, headerbar.selection-mode button.suggested-action { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; background-image: linear-gradient(to top, #edebe9 2px, #f6f5f4); text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button.suggested-action:hover, .selection-mode.titlebar:not(headerbar) button.suggested-action:hover, .selection-mode headerbar button.suggested-action:hover, headerbar.selection-mode button.suggested-action:hover { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px); border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button.suggested-action:active, .selection-mode.titlebar:not(headerbar) button.suggested-action:active, .selection-mode headerbar button.suggested-action:active, headerbar.selection-mode button.suggested-action:active { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; background-image: image(#d6d1cd); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button.suggested-action:disabled, .selection-mode.titlebar:not(headerbar) button.suggested-action:disabled, .selection-mode headerbar button.suggested-action:disabled, headerbar.selection-mode button.suggested-action:disabled { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button.suggested-action:backdrop, .selection-mode.titlebar:not(headerbar) button.suggested-action:backdrop, .selection-mode headerbar button.suggested-action:backdrop, headerbar.selection-mode button.suggested-action:backdrop { color: #929595; border-color: #d5d0cc; background-image: image(#f6f5f4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) button.suggested-action:backdrop:disabled, .selection-mode.titlebar:not(headerbar) button.suggested-action:backdrop:disabled, .selection-mode headerbar button.suggested-action:backdrop:disabled, headerbar.selection-mode button.suggested-action:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #185fb4; }

.selection-mode .titlebar:not(headerbar) .selection-menu, .selection-mode.titlebar:not(headerbar) .selection-menu, .selection-mode headerbar .selection-menu:backdrop, .selection-mode headerbar .selection-menu, headerbar.selection-mode .selection-menu:backdrop, headerbar.selection-mode .selection-menu { border-color: rgba(53, 132, 228, 0); background-color: rgba(53, 132, 228, 0); background-image: none; box-shadow: none; min-height: 20px; padding: 6px 10px; }

.selection-mode .titlebar:not(headerbar) .selection-menu arrow, .selection-mode.titlebar:not(headerbar) .selection-menu arrow, .selection-mode headerbar .selection-menu:backdrop arrow, .selection-mode headerbar .selection-menu arrow, headerbar.selection-mode .selection-menu:backdrop arrow, headerbar.selection-mode .selection-menu arrow { -GtkArrow-arrow-scaling: 1; }

.selection-mode .titlebar:not(headerbar) .selection-menu .arrow, .selection-mode.titlebar:not(headerbar) .selection-menu .arrow, .selection-mode headerbar .selection-menu:backdrop .arrow, .selection-mode headerbar .selection-menu .arrow, headerbar.selection-mode .selection-menu:backdrop .arrow, headerbar.selection-mode .selection-menu .arrow { -gtk-icon-source: -gtk-icontheme("pan-down-symbolic"); color: rgba(255, 255, 255, 0.5); -gtk-icon-shadow: none; }

.tiled .titlebar:not(headerbar), .tiled-top .titlebar:not(headerbar), .tiled-right .titlebar:not(headerbar), .tiled-bottom .titlebar:not(headerbar), .tiled-left .titlebar:not(headerbar), .maximized .titlebar:not(headerbar), .fullscreen .titlebar:not(headerbar), .tiled headerbar, .tiled-top headerbar, .tiled-right headerbar, .tiled-bottom headerbar, .tiled-left headerbar, .maximized headerbar, .fullscreen headerbar { border-radius: 0; }

.default-decoration.titlebar:not(headerbar), headerbar.default-decoration { min-height: 28px; padding: 4px; }

.default-decoration.titlebar:not(headerbar) button.titlebutton, headerbar.default-decoration button.titlebutton { min-height: 26px; min-width: 26px; margin: 0; padding: 0; }

.titlebar:not(headerbar) separator.titlebutton, headerbar separator.titlebutton { opacity: 0; }

.solid-csd .titlebar:dir(rtl):not(headerbar), .solid-csd .titlebar:dir(ltr):not(headerbar), .solid-csd headerbar:backdrop:dir(rtl), .solid-csd headerbar:backdrop:dir(ltr), .solid-csd headerbar:dir(rtl), .solid-csd headerbar:dir(ltr) { margin-left: -1px; margin-right: -1px; margin-top: -1px; border-radius: 0; box-shadow: none; }

headerbar entry, headerbar spinbutton, headerbar separator:not(.sidebar), headerbar button { margin-top: 6px; margin-bottom: 6px; }

headerbar switch { margin-top: 10px; margin-bottom: 10px; }

headerbar.titlebar headerbar:not(.titlebar) { background: none; box-shadow: none; }

.background .titlebar:backdrop, .background .titlebar { border-top-left-radius: 8px; border-top-right-radius: 8px; }

.background.tiled .titlebar:backdrop, .background.tiled .titlebar, .background.tiled-top .titlebar:backdrop, .background.tiled-top .titlebar, .background.tiled-right .titlebar:backdrop, .background.tiled-right .titlebar, .background.tiled-bottom .titlebar:backdrop, .background.tiled-bottom .titlebar, .background.tiled-left .titlebar:backdrop, .background.tiled-left .titlebar, .background.maximized .titlebar:backdrop, .background.maximized .titlebar, .background.solid-csd .titlebar:backdrop, .background.solid-csd .titlebar { border-top-left-radius: 0; border-top-right-radius: 0; }

window separator:first-child + headerbar:backdrop, window separator:first-child + headerbar, window headerbar:first-child:backdrop, window headerbar:first-child { border-top-left-radius: 7px; }

window headerbar:last-child:backdrop, window headerbar:last-child { border-top-right-radius: 7px; }

window stack headerbar:first-child:backdrop, window stack headerbar:first-child, window stack headerbar:last-child:backdrop, window stack headerbar:last-child { border-top-left-radius: 7px; border-top-right-radius: 7px; }

window.tiled headerbar, window.tiled headerbar:first-child, window.tiled headerbar:last-child, window.tiled headerbar:only-child, window.tiled headerbar:backdrop, window.tiled headerbar:backdrop:first-child, window.tiled headerbar:backdrop:last-child, window.tiled headerbar:backdrop:only-child, window.tiled-top headerbar, window.tiled-top headerbar:first-child, window.tiled-top headerbar:last-child, window.tiled-top headerbar:only-child, window.tiled-top headerbar:backdrop, window.tiled-top headerbar:backdrop:first-child, window.tiled-top headerbar:backdrop:last-child, window.tiled-top headerbar:backdrop:only-child, window.tiled-right headerbar, window.tiled-right headerbar:first-child, window.tiled-right headerbar:last-child, window.tiled-right headerbar:only-child, window.tiled-right headerbar:backdrop, window.tiled-right headerbar:backdrop:first-child, window.tiled-right headerbar:backdrop:last-child, window.tiled-right headerbar:backdrop:only-child, window.tiled-bottom headerbar, window.tiled-bottom headerbar:first-child, window.tiled-bottom headerbar:last-child, window.tiled-bottom headerbar:only-child, window.tiled-bottom headerbar:backdrop, window.tiled-bottom headerbar:backdrop:first-child, window.tiled-bottom headerbar:backdrop:last-child, window.tiled-bottom headerbar:backdrop:only-child, window.tiled-left headerbar, window.tiled-left headerbar:first-child, window.tiled-left headerbar:last-child, window.tiled-left headerbar:only-child, window.tiled-left headerbar:backdrop, window.tiled-left headerbar:backdrop:first-child, window.tiled-left headerbar:backdrop:last-child, window.tiled-left headerbar:backdrop:only-child, window.maximized headerbar, window.maximized headerbar:first-child, window.maximized headerbar:last-child, window.maximized headerbar:only-child, window.maximized headerbar:backdrop, window.maximized headerbar:backdrop:first-child, window.maximized headerbar:backdrop:last-child, window.maximized headerbar:backdrop:only-child, window.fullscreen headerbar, window.fullscreen headerbar:first-child, window.fullscreen headerbar:last-child, window.fullscreen headerbar:only-child, window.fullscreen headerbar:backdrop, window.fullscreen headerbar:backdrop:first-child, window.fullscreen headerbar:backdrop:last-child, window.fullscreen headerbar:backdrop:only-child, window.solid-csd headerbar, window.solid-csd headerbar:first-child, window.solid-csd headerbar:last-child, window.solid-csd headerbar:only-child, window.solid-csd headerbar:backdrop, window.solid-csd headerbar:backdrop:first-child, window.solid-csd headerbar:backdrop:last-child, window.solid-csd headerbar:backdrop:only-child { border-top-left-radius: 0; border-top-right-radius: 0; }

window.csd > .titlebar:not(headerbar) { padding: 0; background-color: transparent; background-image: none; border-style: none; border-color: transparent; box-shadow: none; }

.titlebar:not(headerbar) separator { background-color: #cdc7c2; }

window.devel headerbar.titlebar:not(.selection-mode) { background: #f6f5f4 cross-fade(10% -gtk-icontheme("system-run-symbolic"), image(transparent)) 90% 0/256px 256px no-repeat, linear-gradient(to right, transparent 65%, rgba(53, 132, 228, 0.2)), linear-gradient(to top, #d8d4d0, #dfdcd8 3px, #edebe9); }

window.devel headerbar.titlebar:not(.selection-mode):backdrop { background: #f6f5f4 cross-fade(10% -gtk-icontheme("system-run-symbolic"), image(transparent)) 90% 0/256px 256px no-repeat, image(#f6f5f4); /* background-color would flash */ }

/************ Pathbars * */
.path-bar button.text-button, .path-bar button.image-button, .path-bar button { padding-left: 4px; padding-right: 4px; }

.path-bar button.text-button.image-button label { padding-left: 0; padding-right: 0; }

.path-bar button.text-button.image-button label:last-child, .path-bar button label:last-child { padding-right: 8px; }

.path-bar button.text-button.image-button label:first-child, .path-bar button label:first-child { padding-left: 8px; }

.path-bar button image { padding-left: 4px; padding-right: 4px; }

.path-bar button.slider-button { padding-left: 0; padding-right: 0; }

/************** Tree Views * */
treeview.view { border-left-color: #979a9b; border-top-color: #f6f5f4; }

* { -GtkTreeView-horizontal-separator: 4; -GtkTreeView-grid-line-width: 1; -GtkTreeView-grid-line-pattern: ''; -GtkTreeView-tree-line-width: 1; -GtkTreeView-tree-line-pattern: ''; -GtkTreeView-expander-size: 16; }

treeview.view:selected:focus, treeview.view:selected { border-radius: 0; }

treeview.view:selected:backdrop, treeview.view:selected { border-left-color: #9ac2f2; border-top-color: rgba(46, 52, 54, 0.1); }

treeview.view:disabled { color: #929595; }

treeview.view:disabled:selected { color: #86b5ef; }

treeview.view:disabled:selected:backdrop { color: #71a8eb; }

treeview.view:disabled:backdrop { color: #d4cfca; }

treeview.view.separator { min-height: 2px; color: #f6f5f4; }

treeview.view.separator:backdrop { color: rgba(246, 245, 244, 0.1); }

treeview.view:backdrop { border-left-color: #c4c5c5; border-top: #f6f5f4; }

treeview.view:drop(active) { border-style: solid none; border-width: 1px; border-color: #185fb4; }

treeview.view.after:drop(active) { border-top-style: none; }

treeview.view.before:drop(active) { border-bottom-style: none; }

treeview.view.expander { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic"); color: #4d4d4d; }

treeview.view.expander:dir(rtl) { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic-rtl"); }

treeview.view.expander:hover { color: black; }

treeview.view.expander:selected { color: #c2daf7; }

treeview.view.expander:selected:hover { color: #ffffff; }

treeview.view.expander:selected:backdrop { color: #c1d8f5; }

treeview.view.expander:checked { -gtk-icon-source: -gtk-icontheme("pan-down-symbolic"); }

treeview.view.expander:backdrop { color: #b2b4b4; }

treeview.view.progressbar { color: #ffffff; background-color: #3584e4; background-image: image(#3584e4); box-shadow: none; }

treeview.view.progressbar:selected:focus, treeview.view.progressbar:selected { color: #3584e4; background-image: image(#ffffff); }

treeview.view.progressbar:selected:focus:backdrop, treeview.view.progressbar:selected:backdrop { color: #3584e4; background-color: #fcfcfc; }

treeview.view.progressbar:backdrop { color: #fcfcfc; background-image: none; box-shadow: none; }

treeview.view.trough { background-color: rgba(46, 52, 54, 0.1); }

treeview.view.trough:selected:focus, treeview.view.trough:selected { background-color: rgba(255, 255, 255, 0.3); }

treeview.view header button { color: #979a9b; background-color: #ffffff; font-weight: bold; text-shadow: none; box-shadow: none; }

treeview.view header button:hover { color: #636769; box-shadow: none; transition: none; }

treeview.view header button:active { color: #2e3436; transition: none; }

treeview.view button.dnd:active, treeview.view button.dnd:selected, treeview.view button.dnd:hover, treeview.view button.dnd, treeview.view header.button.dnd:active, treeview.view header.button.dnd:selected, treeview.view header.button.dnd:hover, treeview.view header.button.dnd { padding: 0 6px; color: #ffffff; background-image: none; background-color: #3584e4; border-style: none; border-radius: 0; box-shadow: inset 0 0 0 1px #ffffff; text-shadow: none; transition: none; }

treeview.view acceleditor > label { background-color: #3584e4; }

treeview.view header button, treeview.view header button:hover, treeview.view header button:active { padding: 0 6px; background-image: none; border-style: none solid solid none; border-color: #f6f5f4; border-radius: 0; text-shadow: none; }

treeview.view header button:disabled { border-color: #f6f5f4; background-image: none; }

treeview.view header button:backdrop { color: #c4c5c5; border-color: #f6f5f4; border-style: none solid solid none; background-image: none; background-color: #fcfcfc; }

treeview.view header button:backdrop:disabled { border-color: #f6f5f4; background-image: none; }

treeview.view header button:last-child { border-right-style: none; }

/********* Menus * */
menubar, .menubar { -GtkWidget-window-dragging: true; padding: 0px; box-shadow: inset 0 -1px rgba(0, 0, 0, 0.1); }

menubar:backdrop, .menubar:backdrop { background-color: #f6f5f4; }

menubar > menuitem, .menubar > menuitem { min-height: 16px; padding: 4px 8px; }

menubar > menuitem menu:dir(rtl), menubar > menuitem menu:dir(ltr), .menubar > menuitem menu:dir(rtl), .menubar > menuitem menu:dir(ltr) { border-radius: 0; padding: 0; }

menubar > menuitem:hover, .menubar > menuitem:hover { box-shadow: inset 0 -3px #3584e4; color: #1b6acb; }

menubar > menuitem:disabled, .menubar > menuitem:disabled { color: #929595; box-shadow: none; }

menubar .csd.popup decoration, .menubar .csd.popup decoration { border-radius: 0; }

.background.popup { background-color: transparent; }

menu, .menu, .context-menu { margin: 4px; padding: 4px 0px; background-color: #ffffff; border: 1px solid #cdc7c2; }

.csd menu, .csd .menu, .csd .context-menu { border: none; border-radius: 5px; }

menu:backdrop, .menu:backdrop, .context-menu:backdrop { background-color: #fcfcfc; }

menu menuitem, .menu menuitem, .context-menu menuitem { min-height: 16px; min-width: 40px; padding: 4px 6px; text-shadow: none; }

menu menuitem:hover, .menu menuitem:hover, .context-menu menuitem:hover { color: #ffffff; background-color: #3584e4; }

menu menuitem:disabled, .menu menuitem:disabled, .context-menu menuitem:disabled { color: #929595; }

menu menuitem:disabled:backdrop, .menu menuitem:disabled:backdrop, .context-menu menuitem:disabled:backdrop { color: #d4cfca; }

menu menuitem:backdrop, menu menuitem:backdrop:hover, .menu menuitem:backdrop, .menu menuitem:backdrop:hover, .context-menu menuitem:backdrop, .context-menu menuitem:backdrop:hover { color: #929595; background-color: transparent; }

menu menuitem arrow, .menu menuitem arrow, .context-menu menuitem arrow { min-height: 16px; min-width: 16px; }

menu menuitem arrow:dir(ltr), .menu menuitem arrow:dir(ltr), .context-menu menuitem arrow:dir(ltr) { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic"); margin-left: 10px; }

menu menuitem arrow:dir(rtl), .menu menuitem arrow:dir(rtl), .context-menu menuitem arrow:dir(rtl) { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic-rtl"); margin-right: 10px; }

menu menuitem label:dir(rtl), menu menuitem label:dir(ltr), .menu menuitem label:dir(rtl), .menu menuitem label:dir(ltr), .context-menu menuitem label:dir(rtl), .context-menu menuitem label:dir(ltr) { color: inherit; }

menu > arrow, .menu > arrow, .context-menu > arrow { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; min-height: 16px; min-width: 16px; padding: 4px; background-color: #ffffff; border-radius: 0; }

menu > arrow.top, .menu > arrow.top, .context-menu > arrow.top { margin-top: -4px; border-bottom: 1px solid #eaebeb; border-top-right-radius: 5px; border-top-left-radius: 5px; -gtk-icon-source: -gtk-icontheme("pan-up-symbolic"); }

menu > arrow.bottom, .menu > arrow.bottom, .context-menu > arrow.bottom { margin-top: 8px; margin-bottom: -12px; border-top: 1px solid #eaebeb; border-bottom-right-radius: 5px; border-bottom-left-radius: 5px; -gtk-icon-source: -gtk-icontheme("pan-down-symbolic"); }

menu > arrow:hover, .menu > arrow:hover, .context-menu > arrow:hover { background-color: #eaebeb; }

menu > arrow:backdrop, .menu > arrow:backdrop, .context-menu > arrow:backdrop { background-color: #fcfcfc; }

menu > arrow:disabled, .menu > arrow:disabled, .context-menu > arrow:disabled { color: transparent; background-color: transparent; border-color: transparent; }

menuitem accelerator { color: alpha(currentColor,0.55); }

menuitem check, menuitem radio { min-height: 16px; min-width: 16px; }

menuitem check:dir(ltr), menuitem radio:dir(ltr) { margin-right: 7px; }

menuitem check:dir(rtl), menuitem radio:dir(rtl) { margin-left: 7px; }

/*************** Popovers   * */
popover.background { padding: 2px; background-color: #f6f5f4; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.3); }

.csd popover.background, popover.background { border: 1px solid #cdc7c2; border-radius: 9px; }

.csd popover.background { background-clip: padding-box; border-color: rgba(0, 0, 0, 0.23); }

popover.background:backdrop { background-color: #f6f5f4; box-shadow: none; }

popover.background > list, popover.background > .view, popover.background > iconview, popover.background > toolbar { border-style: none; background-color: transparent; }

.csd popover.background.touch-selection, .csd popover.background.magnifier, popover.background.touch-selection, popover.background.magnifier { border: 1px solid rgba(255, 255, 255, 0.1); }

popover.background separator { margin: 3px; }

popover.background list separator { margin: 0px; }

/************* Notebooks * */
notebook > header { padding: 1px; border-color: #cdc7c2; border-width: 1px; background-color: #e1dedb; }

notebook > header:backdrop { border-color: #d5d0cc; background-color: #eae8e6; }

notebook > header tabs { margin: -1px; }

notebook > header.top { border-bottom-style: solid; }

notebook > header.top > tabs { margin-bottom: -2px; }

notebook > header.top > tabs > tab:hover { box-shadow: inset 0 -3px #cdc7c2; }

notebook > header.top > tabs > tab:backdrop { box-shadow: none; }

notebook > header.top > tabs > tab:checked { box-shadow: inset 0 -3px #3584e4; }

notebook > header.bottom { border-top-style: solid; }

notebook > header.bottom > tabs { margin-top: -2px; }

notebook > header.bottom > tabs > tab:hover { box-shadow: inset 0 3px #cdc7c2; }

notebook > header.bottom > tabs > tab:backdrop { box-shadow: none; }

notebook > header.bottom > tabs > tab:checked { box-shadow: inset 0 3px #3584e4; }

notebook > header.left { border-right-style: solid; }

notebook > header.left > tabs { margin-right: -2px; }

notebook > header.left > tabs > tab:hover { box-shadow: inset -3px 0 #cdc7c2; }

notebook > header.left > tabs > tab:backdrop { box-shadow: none; }

notebook > header.left > tabs > tab:checked { box-shadow: inset -3px 0 #3584e4; }

notebook > header.right { border-left-style: solid; }

notebook > header.right > tabs { margin-left: -2px; }

notebook > header.right > tabs > tab:hover { box-shadow: inset 3px 0 #cdc7c2; }

notebook > header.right > tabs > tab:backdrop { box-shadow: none; }

notebook > header.right > tabs > tab:checked { box-shadow: inset 3px 0 #3584e4; }

notebook > header.top > tabs > arrow { border-top-style: none; }

notebook > header.bottom > tabs > arrow { border-bottom-style: none; }

notebook > header.top > tabs > arrow, notebook > header.bottom > tabs > arrow { margin-left: -5px; margin-right: -5px; padding-left: 4px; padding-right: 4px; }

notebook > header.top > tabs > arrow.down, notebook > header.bottom > tabs > arrow.down { -gtk-icon-source: -gtk-icontheme("pan-start-symbolic"); }

notebook > header.top > tabs > arrow.up, notebook > header.bottom > tabs > arrow.up { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic"); }

notebook > header.left > tabs > arrow { border-left-style: none; }

notebook > header.right > tabs > arrow { border-right-style: none; }

notebook > header.left > tabs > arrow, notebook > header.right > tabs > arrow { margin-top: -5px; margin-bottom: -5px; padding-top: 4px; padding-bottom: 4px; }

notebook > header.left > tabs > arrow.down, notebook > header.right > tabs > arrow.down { -gtk-icon-source: -gtk-icontheme("pan-up-symbolic"); }

notebook > header.left > tabs > arrow.up, notebook > header.right > tabs > arrow.up { -gtk-icon-source: -gtk-icontheme("pan-down-symbolic"); }

notebook > header > tabs > arrow { min-height: 16px; min-width: 16px; border-radius: 0; }

notebook > header > tabs > arrow:hover:not(:active):not(:backdrop) { background-clip: padding-box; background-image: none; background-color: rgba(255, 255, 255, 0.3); border-color: transparent; box-shadow: none; }

notebook > header > tabs > arrow:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

notebook > header tab { min-height: 30px; min-width: 30px; padding: 3px 12px; outline-offset: -5px; color: #929595; font-weight: bold; border-width: 1px; border-color: transparent; }

notebook > header tab:hover { color: #606566; }

notebook > header tab.reorderable-page:hover { border-color: rgba(205, 199, 194, 0.3); background-color: rgba(246, 245, 244, 0.2); }

notebook > header tab:backdrop { color: #babbbb; }

notebook > header tab.reorderable-page:backdrop { border-color: transparent; background-color: transparent; }

notebook > header tab:checked { color: #2e3436; }

notebook > header tab.reorderable-page:checked { border-color: rgba(205, 199, 194, 0.5); background-color: rgba(246, 245, 244, 0.5); }

notebook > header tab.reorderable-page:checked:hover { background-color: rgba(246, 245, 244, 0.7); }

notebook > header tab:backdrop:checked { color: #929595; }

notebook > header tab.reorderable-page:backdrop:checked { border-color: #d5d0cc; background-color: #f6f5f4; }

notebook > header tab button.flat { padding: 0; margin-top: 4px; margin-bottom: 4px; min-width: 20px; min-height: 20px; }

notebook > header tab button.flat:hover { color: currentColor; }

notebook > header tab button.flat, notebook > header tab button.flat:backdrop { color: alpha(currentColor,0.3); }

notebook > header tab button.flat:last-child { margin-left: 4px; margin-right: -4px; }

notebook > header tab button.flat:first-child { margin-left: -4px; margin-right: 4px; }

notebook > header.top tabs, notebook > header.bottom tabs { padding-left: 4px; padding-right: 4px; }

notebook > header.top tabs:not(:only-child), notebook > header.bottom tabs:not(:only-child) { margin-left: 3px; margin-right: 3px; }

notebook > header.top tabs:not(:only-child):first-child, notebook > header.bottom tabs:not(:only-child):first-child { margin-left: -1px; }

notebook > header.top tabs:not(:only-child):last-child, notebook > header.bottom tabs:not(:only-child):last-child { margin-right: -1px; }

notebook > header.top tabs tab, notebook > header.bottom tabs tab { margin-left: 4px; margin-right: 4px; }

notebook > header.top tabs tab.reorderable-page, notebook > header.bottom tabs tab.reorderable-page { border-style: none solid; }

notebook > header.left tabs, notebook > header.right tabs { padding-top: 4px; padding-bottom: 4px; }

notebook > header.left tabs:not(:only-child), notebook > header.right tabs:not(:only-child) { margin-top: 3px; margin-bottom: 3px; }

notebook > header.left tabs:not(:only-child):first-child, notebook > header.right tabs:not(:only-child):first-child { margin-top: -1px; }

notebook > header.left tabs:not(:only-child):last-child, notebook > header.right tabs:not(:only-child):last-child { margin-bottom: -1px; }

notebook > header.left tabs tab, notebook > header.right tabs tab { margin-top: 4px; margin-bottom: 4px; }

notebook > header.left tabs tab.reorderable-page, notebook > header.right tabs tab.reorderable-page { border-style: solid none; }

notebook > header.top tab { padding-bottom: 4px; }

notebook > header.bottom tab { padding-top: 4px; }

notebook > stack:not(:only-child) { background-color: #ffffff; }

notebook > stack:not(:only-child):backdrop { background-color: #fcfcfc; }

/************** Scrollbars * */
scrollbar { background-color: #cecece; transition: 300ms cubic-bezier(0.25, 0.46, 0.45, 0.94); }

* { -GtkScrollbar-has-backward-stepper: false; -GtkScrollbar-has-forward-stepper: false; }

scrollbar.top { border-bottom: 1px solid #cdc7c2; }

scrollbar.bottom { border-top: 1px solid #cdc7c2; }

scrollbar.left { border-right: 1px solid #cdc7c2; }

scrollbar.right { border-left: 1px solid #cdc7c2; }

scrollbar:backdrop { background-color: #efedec; border-color: #d5d0cc; transition: 200ms ease-out; }

scrollbar slider { min-width: 6px; min-height: 6px; margin: -1px; border: 4px solid transparent; border-radius: 8px; background-clip: padding-box; background-color: #7e8182; }

scrollbar slider:hover { background-color: #565b5c; }

scrollbar slider:hover:active { background-color: #1b6acb; }

scrollbar slider:backdrop { background-color: #cecfce; }

scrollbar slider:disabled { background-color: transparent; }

scrollbar.fine-tune slider { min-width: 4px; min-height: 4px; }

scrollbar.fine-tune.horizontal slider { border-width: 5px 4px; }

scrollbar.fine-tune.vertical slider { border-width: 4px 5px; }

scrollbar.overlay-indicator:not(.dragging):not(.hovering) { border-color: transparent; opacity: 0.4; background-color: transparent; }

scrollbar.overlay-indicator:not(.dragging):not(.hovering) slider { margin: 0; min-width: 3px; min-height: 3px; background-color: #2e3436; border: 1px solid white; }

scrollbar.overlay-indicator:not(.dragging):not(.hovering) button { min-width: 5px; min-height: 5px; background-color: #2e3436; background-clip: padding-box; border-radius: 100%; border: 1px solid white; -gtk-icon-source: none; }

scrollbar.overlay-indicator.horizontal:not(.dragging):not(.hovering) slider { margin: 0 2px; min-width: 40px; }

scrollbar.overlay-indicator.horizontal:not(.dragging):not(.hovering) button { margin: 1px 2px; min-width: 5px; }

scrollbar.overlay-indicator.vertical:not(.dragging):not(.hovering) slider { margin: 2px 0; min-height: 40px; }

scrollbar.overlay-indicator.vertical:not(.dragging):not(.hovering) button { margin: 2px 1px; min-height: 5px; }

scrollbar.overlay-indicator.dragging, scrollbar.overlay-indicator.hovering { opacity: 0.8; }

scrollbar.horizontal slider { min-width: 40px; }

scrollbar.vertical slider { min-height: 40px; }

scrollbar button { padding: 0; min-width: 12px; min-height: 12px; border-style: none; border-radius: 0; transition-property: min-height, min-width, color; border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #7e8182; }

scrollbar button:hover { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #565b5c; }

scrollbar button:active, scrollbar button:checked { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #1b6acb; }

scrollbar button:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: rgba(126, 129, 130, 0.2); }

scrollbar button:backdrop { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: #cecfce; }

scrollbar button:backdrop:disabled { border-color: transparent; background-color: transparent; background-image: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; color: rgba(206, 207, 206, 0.2); }

scrollbar.vertical button.down { -gtk-icon-source: -gtk-icontheme("pan-down-symbolic"); }

scrollbar.vertical button.up { -gtk-icon-source: -gtk-icontheme("pan-up-symbolic"); }

scrollbar.horizontal button.down { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic"); }

scrollbar.horizontal button.up { -gtk-icon-source: -gtk-icontheme("pan-start-symbolic"); }

treeview ~ scrollbar.vertical { border-top: 1px solid #cdc7c2; margin-top: -1px; }

/********** Switch * */
switch { outline-offset: -4px; border: 1px solid #cdc7c2; border-radius: 14px; color: #2e3436; background-color: #e1dedb; text-shadow: 0 1px rgba(0, 0, 0, 0.1); /* only show i / o for the accessible theme */ }

switch:checked { color: #ffffff; border-color: #15539e; background-color: #3584e4; text-shadow: 0 1px rgba(24, 95, 180, 0.5), 0 0 2px rgba(255, 255, 255, 0.6); }

switch:disabled { color: #929595; border-color: #cdc7c2; background-color: #faf9f8; text-shadow: none; }

switch:backdrop { color: #929595; border-color: #d5d0cc; background-color: #eae8e6; text-shadow: none; transition: 200ms ease-out; }

switch:backdrop:checked { color: #f6f5f4; border-color: #15539e; background-color: #3584e4; }

switch:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-color: #faf9f8; }

switch slider { margin: -1px; min-width: 24px; min-height: 24px; border: 1px solid; border-radius: 50%; transition: all 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); -gtk-outline-radius: 20px; color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); background-image: linear-gradient(to bottom, white 20%, #f6f5f4 90%); box-shadow: inset 0 1px white, 0 1px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.07); }

switch image { color: transparent; }

switch:hover slider { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #bfb8b1; box-shadow: inset 0 1px white, 0 1px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to bottom, white 10%, white 90%); }

switch:checked > slider { border: 1px solid #15539e; }

switch:disabled slider { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

switch:backdrop slider { transition: 200ms ease-out; color: #929595; border-color: #d5d0cc; background-image: image(#f6f5f4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

switch:backdrop:checked > slider { border-color: #15539e; }

switch:backdrop:disabled slider { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

row:selected switch { box-shadow: none; border-color: #15539e; }

row:selected switch:backdrop { border-color: #15539e; }

row:selected switch > slider:checked, row:selected switch > slider { border-color: #15539e; }

/************************* Check and Radio items * */
.view.content-view.check:not(list), iconview.content-view.check:not(list), .content-view .tile check:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: transparent; background-color: #3584e4; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: none; -gtk-icon-shadow: none; }

.view.content-view.check:hover:not(list), iconview.content-view.check:hover:not(list), .content-view .tile check:hover:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: transparent; background-color: #3584e4; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: none; -gtk-icon-shadow: none; }

.view.content-view.check:active:not(list), iconview.content-view.check:active:not(list), .content-view .tile check:active:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: transparent; background-color: #3584e4; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: none; -gtk-icon-shadow: none; }

.view.content-view.check:backdrop:not(list), iconview.content-view.check:backdrop:not(list), .content-view .tile check:backdrop:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: transparent; background-color: #8d8d8d; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: none; -gtk-icon-shadow: none; }

.view.content-view.check:checked:not(list), iconview.content-view.check:checked:not(list), .content-view .tile check:checked:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: #eeeeec; background-color: #3584e4; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: -gtk-icontheme('object-select-symbolic'); -gtk-icon-shadow: none; }

.view.content-view.check:checked:hover:not(list), iconview.content-view.check:checked:hover:not(list), .content-view .tile check:checked:hover:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: #eeeeec; background-color: #3584e4; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: -gtk-icontheme('object-select-symbolic'); -gtk-icon-shadow: none; }

.view.content-view.check:checked:active:not(list), iconview.content-view.check:checked:active:not(list), .content-view .tile check:checked:active:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: #eeeeec; background-color: #3584e4; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: -gtk-icontheme('object-select-symbolic'); -gtk-icon-shadow: none; }

.view.content-view.check:backdrop:checked:not(list), iconview.content-view.check:backdrop:checked:not(list), .content-view .tile check:backdrop:checked:not(list) { margin: 4px; min-width: 32px; min-height: 32px; color: rgba(238, 238, 236, 0.8); background-color: #8d8d8d; border-radius: 5px; background-image: none; transition: 200ms; box-shadow: none; border-width: 0; -gtk-icon-source: -gtk-icontheme('object-select-symbolic'); -gtk-icon-shadow: none; }

checkbutton.text-button, radiobutton.text-button { padding: 2px 0; outline-offset: 0; }

checkbutton.text-button label:not(:only-child):first-child, radiobutton.text-button label:not(:only-child):first-child { margin-left: 4px; }

checkbutton.text-button label:not(:only-child):last-child, radiobutton.text-button label:not(:only-child):last-child { margin-right: 4px; }

check, radio { margin: 0 4px; min-height: 14px; min-width: 14px; border: 1px solid; -gtk-icon-source: none; }

check:only-child, radio:only-child { margin: 0; }

popover check.left:dir(rtl), popover radio.left:dir(rtl) { margin-left: 0; margin-right: 12px; }

popover check.right:dir(ltr), popover radio.right:dir(ltr) { margin-left: 12px; margin-right: 0; }

check, radio { background-clip: padding-box; background-image: linear-gradient(to bottom, white 20%, white 90%); border-color: #bfb8b1; box-shadow: 0 1px rgba(0, 0, 0, 0.05); color: #2e3436; }

check:hover, radio:hover { background-image: image(#f2f2f2); }

check:active, radio:active { box-shadow: inset 0 1px 1px 0px rgba(0, 0, 0, 0.2); }

check:disabled, radio:disabled { box-shadow: none; color: rgba(46, 52, 54, 0.7); }

check:backdrop, radio:backdrop { background-image: image(white); box-shadow: none; color: #2e3436; }

check:backdrop:disabled, radio:backdrop:disabled { box-shadow: none; color: rgba(46, 52, 54, 0.7); }

check:checked, radio:checked { background-clip: border-box; background-image: linear-gradient(to bottom, #4b92e7 20%, #3584e4 90%); border-color: #3584e4; box-shadow: 0 1px rgba(0, 0, 0, 0.05); color: #ffffff; }

check:checked:hover, radio:checked:hover { background-image: linear-gradient(to bottom, #5d9de9 10%, #478fe6 90%); }

check:checked:active, radio:checked:active { box-shadow: inset 0 1px 1px 0px rgba(0, 0, 0, 0.2); }

check:checked:disabled, radio:checked:disabled { box-shadow: none; color: rgba(255, 255, 255, 0.7); }

check:checked:backdrop, radio:checked:backdrop { background-image: image(#3584e4); box-shadow: none; color: #ffffff; }

check:checked:backdrop:disabled, radio:checked:backdrop:disabled { box-shadow: none; color: rgba(255, 255, 255, 0.7); }

check:indeterminate, radio:indeterminate { background-clip: border-box; background-image: linear-gradient(to bottom, #4b92e7 20%, #3584e4 90%); border-color: #3584e4; box-shadow: 0 1px rgba(0, 0, 0, 0.05); color: #ffffff; }

check:indeterminate:hover, radio:indeterminate:hover { background-image: linear-gradient(to bottom, #5d9de9 10%, #478fe6 90%); }

check:indeterminate:active, radio:indeterminate:active { box-shadow: inset 0 1px 1px 0px rgba(0, 0, 0, 0.2); }

check:indeterminate:disabled, radio:indeterminate:disabled { box-shadow: none; color: rgba(255, 255, 255, 0.7); }

check:indeterminate:backdrop, radio:indeterminate:backdrop { background-image: image(#3584e4); box-shadow: none; color: #ffffff; }

check:indeterminate:backdrop:disabled, radio:indeterminate:backdrop:disabled { box-shadow: none; color: rgba(255, 255, 255, 0.7); }

check:backdrop, radio:backdrop { transition: 200ms ease-out; }

row:selected check, row:selected radio { border-color: #15539e; }

.osd check, .osd radio { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd check:hover, .osd radio:hover { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); }

.osd check:active, .osd radio:active { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(0, 0, 0, 0.7)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); }

.osd check:backdrop, .osd radio:backdrop { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

.osd check:disabled, .osd radio:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; }

menu menuitem check, menu menuitem radio { margin: 0; }

menu menuitem check, menu menuitem check:hover, menu menuitem check:disabled, menu menuitem check:checked, menu menuitem check:checked:hover, menu menuitem check:checked:disabled, menu menuitem check:indeterminate, menu menuitem check:indeterminate:hover, menu menuitem check:indeterminate:disabled, menu menuitem radio, menu menuitem radio:hover, menu menuitem radio:disabled, menu menuitem radio:checked, menu menuitem radio:checked:hover, menu menuitem radio:checked:disabled, menu menuitem radio:indeterminate, menu menuitem radio:indeterminate:hover, menu menuitem radio:indeterminate:disabled { min-height: 14px; min-width: 14px; background-image: none; background-color: transparent; box-shadow: none; -gtk-icon-shadow: none; color: inherit; border-color: currentColor; }

check { border-radius: 3px; }

check:checked { -gtk-icon-source: image(-gtk-recolor(url("assets/check-symbolic.svg")), -gtk-recolor(url("assets/check-symbolic.symbolic.png"))); }

check:indeterminate { -gtk-icon-source: image(-gtk-recolor(url("assets/dash-symbolic.svg")), -gtk-recolor(url("assets/dash-symbolic.symbolic.png"))); }

treeview.view radio:selected:focus, treeview.view radio:selected, radio { border-radius: 100%; }

treeview.view radio:checked:selected, radio:checked { -gtk-icon-source: image(-gtk-recolor(url("assets/bullet-symbolic.svg")), -gtk-recolor(url("assets/bullet-symbolic.symbolic.png"))); }

treeview.view radio:indeterminate:selected, radio:indeterminate { -gtk-icon-source: image(-gtk-recolor(url("assets/dash-symbolic.svg")), -gtk-recolor(url("assets/dash-symbolic.symbolic.png"))); }

radio:not(:indeterminate):not(:checked):active:not(:backdrop) { -gtk-icon-transform: scale(0); }

check:not(:indeterminate):not(:checked):active:not(:backdrop) { -gtk-icon-transform: translate(6px, -3px) rotate(-45deg) scaleY(0.2) rotate(45deg) scaleX(0); }

radio:active, check:active { -gtk-icon-transform: scale(0, 1); }

radio:checked:not(:backdrop), radio:indeterminate:not(:backdrop), check:checked:not(:backdrop), check:indeterminate:not(:backdrop) { -gtk-icon-transform: unset; transition: 400ms; }

menu menuitem radio:checked:not(:backdrop), menu menuitem radio:indeterminate:not(:backdrop), menu menuitem check:checked:not(:backdrop), menu menuitem check:indeterminate:not(:backdrop) { transition: none; }

treeview.view check:selected:focus, treeview.view check:selected, treeview.view radio:selected:focus, treeview.view radio:selected { color: #ffffff; border-color: #185fb4; }

/************ GtkScale * */
progressbar trough, scale fill, scale trough { border: 1px solid #cdc7c2; border-radius: 3px; background-color: #e1dedb; }

progressbar trough:disabled, scale fill:disabled, scale trough:disabled { background-color: #faf9f8; }

progressbar trough:backdrop, scale fill:backdrop, scale trough:backdrop { background-color: #eae8e6; border-color: #d5d0cc; transition: 200ms ease-out; }

progressbar trough:backdrop:disabled, scale fill:backdrop:disabled, scale trough:backdrop:disabled { background-color: #faf9f8; }

row:selected progressbar trough, progressbar row:selected trough, row:selected scale fill, scale row:selected fill, row:selected scale trough, scale row:selected trough { border-color: #185fb4; }

.osd progressbar trough, progressbar .osd trough, .osd scale fill, scale .osd fill, .osd scale trough, scale .osd trough { border-color: rgba(0, 0, 0, 0.7); background-color: rgba(0, 0, 0, 0.5); }

.osd progressbar trough:disabled, progressbar .osd trough:disabled, .osd scale fill:disabled, scale .osd fill:disabled, .osd scale trough:disabled, scale .osd trough:disabled { background-color: rgba(71, 71, 71, 0.5); }

progressbar progress, scale highlight { border: 1px solid #185fb4; border-radius: 3px; background-color: #3584e4; }

progressbar progress:disabled, scale highlight:disabled { background-color: transparent; border-color: transparent; }

progressbar progress:backdrop, scale highlight:backdrop { border-color: #3584e4; }

progressbar progress:backdrop:disabled, scale highlight:backdrop:disabled { background-color: transparent; border-color: transparent; }

row:selected progressbar progress, progressbar row:selected progress, row:selected scale highlight, scale row:selected highlight { border-color: #185fb4; }

.osd progressbar progress, progressbar .osd progress, .osd scale highlight, scale .osd highlight { border-color: rgba(0, 0, 0, 0.7); }

.osd progressbar progress:disabled, progressbar .osd progress:disabled, .osd scale highlight:disabled, scale .osd highlight:disabled { border-color: transparent; }

scale { min-height: 10px; min-width: 10px; padding: 12px; }

scale fill, scale highlight { margin: -1px; }

scale slider { min-height: 18px; min-width: 18px; margin: -9px; }

scale.fine-tune.horizontal { padding-top: 9px; padding-bottom: 9px; min-height: 16px; }

scale.fine-tune.vertical { padding-left: 9px; padding-right: 9px; min-width: 16px; }

scale.fine-tune slider { margin: -6px; }

scale.fine-tune fill, scale.fine-tune highlight, scale.fine-tune trough { border-radius: 5px; -gtk-outline-radius: 7px; }

scale trough { outline-offset: 2px; -gtk-outline-radius: 5px; }

scale fill:backdrop, scale fill { background-color: #cdc7c2; }

scale fill:disabled:backdrop, scale fill:disabled { border-color: transparent; background-color: transparent; }

.osd scale fill { background-color: rgba(91, 91, 90, 0.775); }

.osd scale fill:disabled:backdrop, .osd scale fill:disabled { border-color: transparent; background-color: transparent; }

scale slider { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); background-image: linear-gradient(to bottom, white 20%, #f6f5f4 90%); box-shadow: inset 0 1px white, 0 1px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.07); border: 1px solid #b8b0a8; border-radius: 100%; transition: all 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); transition-property: background, border, box-shadow; }

scale slider:hover { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #bfb8b1; box-shadow: inset 0 1px white, 0 1px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to bottom, white 10%, white 90%); }

scale slider:active { border-color: #185fb4; }

scale slider:disabled { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

scale slider:backdrop { transition: 200ms ease-out; color: #929595; border-color: #d5d0cc; background-image: image(#f6f5f4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

scale slider:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

row:selected scale slider:disabled, row:selected scale slider { border-color: #185fb4; }

.osd scale slider { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); border-color: rgba(0, 0, 0, 0.7); background-color: #353535; }

.osd scale slider:hover { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(83, 83, 83, 0.9)); background-clip: padding-box; box-shadow: inset 0 1px rgba(255, 255, 255, 0.1); text-shadow: 0 1px black; -gtk-icon-shadow: 0 1px black; outline-color: rgba(238, 238, 236, 0.3); background-color: #353535; }

.osd scale slider:active { color: white; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(0, 0, 0, 0.7)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; outline-color: rgba(238, 238, 236, 0.3); background-color: #353535; }

.osd scale slider:disabled { color: #919190; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(71, 71, 71, 0.5)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; background-color: #353535; }

.osd scale slider:backdrop { color: #eeeeec; border-color: rgba(0, 0, 0, 0.7); background-color: transparent; background-image: image(rgba(53, 53, 53, 0.9)); background-clip: padding-box; box-shadow: none; text-shadow: none; -gtk-icon-shadow: none; background-color: #353535; }

.osd scale slider:backdrop:disabled { background-color: #353535; }

scale marks, scale value { color: alpha(currentColor,0.55); font-feature-settings: "tnum"; }

scale.horizontal marks.top { margin-bottom: 6px; margin-top: -12px; }

scale.horizontal.fine-tune marks.top { margin-bottom: 6px; margin-top: -9px; }

scale.horizontal marks.bottom { margin-top: 6px; margin-bottom: -12px; }

scale.horizontal.fine-tune marks.bottom { margin-top: 6px; margin-bottom: -9px; }

scale.vertical marks.top { margin-right: 6px; margin-left: -12px; }

scale.vertical.fine-tune marks.top { margin-right: 6px; margin-left: -9px; }

scale.vertical marks.bottom { margin-left: 6px; margin-right: -12px; }

scale.vertical.fine-tune marks.bottom { margin-left: 6px; margin-right: -9px; }

scale.horizontal indicator { min-height: 6px; min-width: 1px; }

scale.horizontal.fine-tune indicator { min-height: 3px; }

scale.vertical indicator { min-height: 1px; min-width: 6px; }

scale.vertical.fine-tune indicator { min-width: 3px; }

scale.horizontal.marks-before:not(.marks-after) slider { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above.png"), url("assets/slider-horz-scale-has-marks-above@2.png")); min-height: 26px; min-width: 22px; margin-top: -14px; background-position: top; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-top: -11px; }

scale.horizontal.marks-before:not(.marks-after) slider:hover { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above-hover.png"), url("assets/slider-horz-scale-has-marks-above-hover@2.png")); min-height: 26px; min-width: 22px; margin-top: -14px; background-position: top; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-top: -11px; }

scale.horizontal.marks-before:not(.marks-after) slider:active { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above-active.png"), url("assets/slider-horz-scale-has-marks-above-active@2.png")); min-height: 26px; min-width: 22px; margin-top: -14px; background-position: top; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-top: -11px; }

scale.horizontal.marks-before:not(.marks-after) slider:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above-insensitive.png"), url("assets/slider-horz-scale-has-marks-above-insensitive@2.png")); min-height: 26px; min-width: 22px; margin-top: -14px; background-position: top; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-top: -11px; }

scale.horizontal.marks-before:not(.marks-after) slider:backdrop { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above-backdrop.png"), url("assets/slider-horz-scale-has-marks-above-backdrop@2.png")); min-height: 26px; min-width: 22px; margin-top: -14px; background-position: top; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-top: -11px; }

scale.horizontal.marks-before:not(.marks-after) slider:backdrop:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above-backdrop-insensitive.png"), url("assets/slider-horz-scale-has-marks-above-backdrop-insensitive@2.png")); min-height: 26px; min-width: 22px; margin-top: -14px; background-position: top; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-top: -11px; }

scale.horizontal.marks-after:not(.marks-before) slider { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-below.png"), url("assets/slider-horz-scale-has-marks-below@2.png")); min-height: 26px; min-width: 22px; margin-bottom: -14px; background-position: bottom; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-bottom: -11px; }

scale.horizontal.marks-after:not(.marks-before) slider:hover { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-below-hover.png"), url("assets/slider-horz-scale-has-marks-below-hover@2.png")); min-height: 26px; min-width: 22px; margin-bottom: -14px; background-position: bottom; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-bottom: -11px; }

scale.horizontal.marks-after:not(.marks-before) slider:active { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-below-active.png"), url("assets/slider-horz-scale-has-marks-below-active@2.png")); min-height: 26px; min-width: 22px; margin-bottom: -14px; background-position: bottom; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-bottom: -11px; }

scale.horizontal.marks-after:not(.marks-before) slider:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-below-insensitive.png"), url("assets/slider-horz-scale-has-marks-below-insensitive@2.png")); min-height: 26px; min-width: 22px; margin-bottom: -14px; background-position: bottom; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-bottom: -11px; }

scale.horizontal.marks-after:not(.marks-before) slider:backdrop { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-below-backdrop.png"), url("assets/slider-horz-scale-has-marks-below-backdrop@2.png")); min-height: 26px; min-width: 22px; margin-bottom: -14px; background-position: bottom; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-bottom: -11px; }

scale.horizontal.marks-after:not(.marks-before) slider:backdrop:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-horz-scale-has-marks-below-backdrop-insensitive.png"), url("assets/slider-horz-scale-has-marks-below-backdrop-insensitive@2.png")); min-height: 26px; min-width: 22px; margin-bottom: -14px; background-position: bottom; background-repeat: no-repeat; box-shadow: none; }

scale.horizontal.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-bottom: -11px; }

scale.vertical.marks-before:not(.marks-after) slider { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-above.png"), url("assets/slider-vert-scale-has-marks-above@2.png")); min-height: 22px; min-width: 26px; margin-left: -14px; background-position: left bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-left: -11px; }

scale.vertical.marks-before:not(.marks-after) slider:hover { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-above-hover.png"), url("assets/slider-vert-scale-has-marks-above-hover@2.png")); min-height: 22px; min-width: 26px; margin-left: -14px; background-position: left bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-left: -11px; }

scale.vertical.marks-before:not(.marks-after) slider:active { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-above-active.png"), url("assets/slider-vert-scale-has-marks-above-active@2.png")); min-height: 22px; min-width: 26px; margin-left: -14px; background-position: left bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-left: -11px; }

scale.vertical.marks-before:not(.marks-after) slider:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-above-insensitive.png"), url("assets/slider-vert-scale-has-marks-above-insensitive@2.png")); min-height: 22px; min-width: 26px; margin-left: -14px; background-position: left bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-left: -11px; }

scale.vertical.marks-before:not(.marks-after) slider:backdrop { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-above-backdrop.png"), url("assets/slider-vert-scale-has-marks-above-backdrop@2.png")); min-height: 22px; min-width: 26px; margin-left: -14px; background-position: left bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-left: -11px; }

scale.vertical.marks-before:not(.marks-after) slider:backdrop:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-above-backdrop-insensitive.png"), url("assets/slider-vert-scale-has-marks-above-backdrop-insensitive@2.png")); min-height: 22px; min-width: 26px; margin-left: -14px; background-position: left bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-before.fine-tune:not(.marks-after) slider { margin: -7px; margin-left: -11px; }

scale.vertical.marks-after:not(.marks-before) slider { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-below.png"), url("assets/slider-vert-scale-has-marks-below@2.png")); min-height: 22px; min-width: 26px; margin-right: -14px; background-position: right bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-right: -11px; }

scale.vertical.marks-after:not(.marks-before) slider:hover { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-below-hover.png"), url("assets/slider-vert-scale-has-marks-below-hover@2.png")); min-height: 22px; min-width: 26px; margin-right: -14px; background-position: right bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-right: -11px; }

scale.vertical.marks-after:not(.marks-before) slider:active { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-below-active.png"), url("assets/slider-vert-scale-has-marks-below-active@2.png")); min-height: 22px; min-width: 26px; margin-right: -14px; background-position: right bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-right: -11px; }

scale.vertical.marks-after:not(.marks-before) slider:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-below-insensitive.png"), url("assets/slider-vert-scale-has-marks-below-insensitive@2.png")); min-height: 22px; min-width: 26px; margin-right: -14px; background-position: right bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-right: -11px; }

scale.vertical.marks-after:not(.marks-before) slider:backdrop { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-below-backdrop.png"), url("assets/slider-vert-scale-has-marks-below-backdrop@2.png")); min-height: 22px; min-width: 26px; margin-right: -14px; background-position: right bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-right: -11px; }

scale.vertical.marks-after:not(.marks-before) slider:backdrop:disabled { margin: -10px; border-style: none; border-radius: 0; background-color: transparent; background-image: -gtk-scaled(url("assets/slider-vert-scale-has-marks-below-backdrop-insensitive.png"), url("assets/slider-vert-scale-has-marks-below-backdrop-insensitive@2.png")); min-height: 22px; min-width: 26px; margin-right: -14px; background-position: right bottom; background-repeat: no-repeat; box-shadow: none; }

scale.vertical.marks-after.fine-tune:not(.marks-before) slider { margin: -7px; margin-right: -11px; }

scale.color { min-height: 0; min-width: 0; }

scale.color trough { background-image: image(#cdc7c2); background-repeat: no-repeat; }

scale.color.horizontal { padding: 0 0 15px 0; }

scale.color.horizontal trough { padding-bottom: 4px; background-position: 0 -3px; border-top-left-radius: 0; border-top-right-radius: 0; }

scale.color.horizontal slider:dir(ltr):hover, scale.color.horizontal slider:dir(ltr):backdrop, scale.color.horizontal slider:dir(ltr):disabled, scale.color.horizontal slider:dir(ltr):backdrop:disabled, scale.color.horizontal slider:dir(ltr), scale.color.horizontal slider:dir(rtl):hover, scale.color.horizontal slider:dir(rtl):backdrop, scale.color.horizontal slider:dir(rtl):disabled, scale.color.horizontal slider:dir(rtl):backdrop:disabled, scale.color.horizontal slider:dir(rtl) { margin-bottom: -15px; margin-top: 6px; }

scale.color.vertical:dir(ltr) { padding: 0 0 0 15px; }

scale.color.vertical:dir(ltr) trough { padding-left: 4px; background-position: 3px 0; border-bottom-right-radius: 0; border-top-right-radius: 0; }

scale.color.vertical:dir(ltr) slider:hover, scale.color.vertical:dir(ltr) slider:backdrop, scale.color.vertical:dir(ltr) slider:disabled, scale.color.vertical:dir(ltr) slider:backdrop:disabled, scale.color.vertical:dir(ltr) slider { margin-left: -15px; margin-right: 6px; }

scale.color.vertical:dir(rtl) { padding: 0 15px 0 0; }

scale.color.vertical:dir(rtl) trough { padding-right: 4px; background-position: -3px 0; border-bottom-left-radius: 0; border-top-left-radius: 0; }

scale.color.vertical:dir(rtl) slider:hover, scale.color.vertical:dir(rtl) slider:backdrop, scale.color.vertical:dir(rtl) slider:disabled, scale.color.vertical:dir(rtl) slider:backdrop:disabled, scale.color.vertical:dir(rtl) slider { margin-right: -15px; margin-left: 6px; }

scale.color.fine-tune.horizontal:dir(ltr), scale.color.fine-tune.horizontal:dir(rtl) { padding: 0 0 12px 0; }

scale.color.fine-tune.horizontal:dir(ltr) trough, scale.color.fine-tune.horizontal:dir(rtl) trough { padding-bottom: 7px; background-position: 0 -6px; }

scale.color.fine-tune.horizontal:dir(ltr) slider, scale.color.fine-tune.horizontal:dir(rtl) slider { margin-bottom: -15px; margin-top: 6px; }

scale.color.fine-tune.vertical:dir(ltr) { padding: 0 0 0 12px; }

scale.color.fine-tune.vertical:dir(ltr) trough { padding-left: 7px; background-position: 6px 0; }

scale.color.fine-tune.vertical:dir(ltr) slider { margin-left: -15px; margin-right: 6px; }

scale.color.fine-tune.vertical:dir(rtl) { padding: 0 12px 0 0; }

scale.color.fine-tune.vertical:dir(rtl) trough { padding-right: 7px; background-position: -6px 0; }

scale.color.fine-tune.vertical:dir(rtl) slider { margin-right: -15px; margin-left: 6px; }

/***************** Progress bars * */
progressbar { font-size: smaller; color: rgba(46, 52, 54, 0.4); font-feature-settings: "tnum"; }

progressbar.horizontal trough, progressbar.horizontal progress { min-height: 2px; }

progressbar.vertical trough, progressbar.vertical progress { min-width: 2px; }

progressbar.horizontal progress { margin: 0 -1px; }

progressbar.vertical progress { margin: -1px 0; }

progressbar:backdrop { box-shadow: none; transition: 200ms ease-out; }

progressbar progress { border-radius: 1.5px; }

progressbar progress.left { border-top-left-radius: 2px; border-bottom-left-radius: 2px; }

progressbar progress.right { border-top-right-radius: 2px; border-bottom-right-radius: 2px; }

progressbar progress.top { border-top-right-radius: 2px; border-top-left-radius: 2px; }

progressbar progress.bottom { border-bottom-right-radius: 2px; border-bottom-left-radius: 2px; }

progressbar.osd { min-width: 3px; min-height: 3px; background-color: transparent; }

progressbar.osd trough { border-style: none; border-radius: 0; background-color: transparent; box-shadow: none; }

progressbar.osd progress { border-style: none; border-radius: 0; }

progressbar trough.empty progress { all: unset; }

/************* Level Bar * */
levelbar.horizontal block { min-height: 1px; }

levelbar.horizontal.discrete block { margin: 0 1px; min-width: 32px; }

levelbar.vertical block { min-width: 1px; }

levelbar.vertical.discrete block { margin: 1px 0; min-height: 32px; }

levelbar:backdrop { transition: 200ms ease-out; }

levelbar trough { border: 1px solid; padding: 2px; border-radius: 3px; color: black; border-color: #cdc7c2; background-color: #ffffff; box-shadow: inset 0 0 0 1px rgba(53, 132, 228, 0); }

levelbar trough:backdrop { color: #323232; border-color: #d5d0cc; background-color: #fcfcfc; box-shadow: none; }

levelbar block { border: 1px solid; border-radius: 1px; }

levelbar block.low { border-color: #8f4700; background-color: #f57900; }

levelbar block.low:backdrop { border-color: #f57900; }

levelbar block.high, levelbar block:not(.empty) { border-color: #15539e; background-color: #3584e4; }

levelbar block.high:backdrop, levelbar block:not(.empty):backdrop { border-color: #3584e4; }

levelbar block.full { border-color: #1d814a; background-color: #33d17a; }

levelbar block.full:backdrop { border-color: #33d17a; }

levelbar block.empty { background-color: transparent; border-color: rgba(46, 52, 54, 0.2); }

levelbar block.empty:backdrop { border-color: rgba(146, 149, 149, 0.15); }

/**************** Print dialog * */
printdialog paper { color: #2e3436; border: 1px solid #cdc7c2; background: white; padding: 0; }

printdialog paper:backdrop { color: #929595; border-color: #d5d0cc; }

printdialog .dialog-action-box { margin: 12px; }

/********** Frames * */
frame > border, .frame { box-shadow: none; margin: 0; padding: 0; border-radius: 0; border: 1px solid #cdc7c2; }

frame > border.flat, .frame.flat { border-style: none; }

frame > border:backdrop, .frame:backdrop { border-color: #d5d0cc; }

actionbar > revealer > box { padding: 6px; border-top: 1px solid #cdc7c2; }

actionbar > revealer > box:backdrop { border-color: #d5d0cc; }

scrolledwindow viewport.frame { border-style: none; }

scrolledwindow overshoot.top { background-image: -gtk-gradient(radial, center top, 0, center top, 0.5, to(#b6aea5), to(rgba(182, 174, 165, 0))), -gtk-gradient(radial, center top, 0, center top, 0.6, from(rgba(46, 52, 54, 0.07)), to(rgba(46, 52, 54, 0))); background-size: 100% 5%, 100% 100%; background-repeat: no-repeat; background-position: center top; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow overshoot.top:backdrop { background-image: -gtk-gradient(radial, center top, 0, center top, 0.5, to(#d5d0cc), to(rgba(213, 208, 204, 0))); background-size: 100% 5%; background-repeat: no-repeat; background-position: center top; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow overshoot.bottom { background-image: -gtk-gradient(radial, center bottom, 0, center bottom, 0.5, to(#b6aea5), to(rgba(182, 174, 165, 0))), -gtk-gradient(radial, center bottom, 0, center bottom, 0.6, from(rgba(46, 52, 54, 0.07)), to(rgba(46, 52, 54, 0))); background-size: 100% 5%, 100% 100%; background-repeat: no-repeat; background-position: center bottom; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow overshoot.bottom:backdrop { background-image: -gtk-gradient(radial, center bottom, 0, center bottom, 0.5, to(#d5d0cc), to(rgba(213, 208, 204, 0))); background-size: 100% 5%; background-repeat: no-repeat; background-position: center bottom; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow overshoot.left { background-image: -gtk-gradient(radial, left center, 0, left center, 0.5, to(#b6aea5), to(rgba(182, 174, 165, 0))), -gtk-gradient(radial, left center, 0, left center, 0.6, from(rgba(46, 52, 54, 0.07)), to(rgba(46, 52, 54, 0))); background-size: 5% 100%, 100% 100%; background-repeat: no-repeat; background-position: left center; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow overshoot.left:backdrop { background-image: -gtk-gradient(radial, left center, 0, left center, 0.5, to(#d5d0cc), to(rgba(213, 208, 204, 0))); background-size: 5% 100%; background-repeat: no-repeat; background-position: left center; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow overshoot.right { background-image: -gtk-gradient(radial, right center, 0, right center, 0.5, to(#b6aea5), to(rgba(182, 174, 165, 0))), -gtk-gradient(radial, right center, 0, right center, 0.6, from(rgba(46, 52, 54, 0.07)), to(rgba(46, 52, 54, 0))); background-size: 5% 100%, 100% 100%; background-repeat: no-repeat; background-position: right center; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow overshoot.right:backdrop { background-image: -gtk-gradient(radial, right center, 0, right center, 0.5, to(#d5d0cc), to(rgba(213, 208, 204, 0))); background-size: 5% 100%; background-repeat: no-repeat; background-position: right center; background-color: transparent; border: none; box-shadow: none; }

scrolledwindow junction { border-color: transparent; border-image: linear-gradient(to bottom, #cdc7c2 1px, transparent 1px) 0 0 0 1/0 1px stretch; background-color: #cecece; }

scrolledwindow junction:dir(rtl) { border-image-slice: 0 1 0 0; }

scrolledwindow junction:backdrop { border-image-source: linear-gradient(to bottom, #d5d0cc 1px, transparent 1px); background-color: #efedec; transition: 200ms ease-out; }

separator { background: rgba(0, 0, 0, 0.1); min-width: 1px; min-height: 1px; }

/********* Lists * */
list { color: black; background-color: #ffffff; border-color: #cdc7c2; }

list:backdrop { color: #323232; background-color: #fcfcfc; border-color: #d5d0cc; }

list row { padding: 2px; }

row { transition: all 150ms cubic-bezier(0.25, 0.46, 0.45, 0.94); }

row:hover { transition: none; }

row:backdrop { transition: 200ms ease-out; }

row.activatable.has-open-popup, row.activatable:hover { background-color: rgba(46, 52, 54, 0.05); }

row.activatable:active { box-shadow: inset 0 2px 2px -2px rgba(0, 0, 0, 0.2); }

row.activatable:backdrop:hover { background-color: transparent; }

row.activatable:selected:active { box-shadow: inset 0 2px 3px -1px rgba(0, 0, 0, 0.5); }

row.activatable.has-open-popup:selected, row.activatable:selected:hover { background-color: #347cd3; }

row.activatable:selected:backdrop { background-color: #3584e4; }

/********************* App Notifications * */
.app-notification, .app-notification.frame { padding: 10px; border-radius: 0 0 5px 5px; background-color: rgba(53, 53, 53, 0.9); background-image: linear-gradient(to bottom, rgba(0, 0, 0, 0.2), transparent 2px); background-clip: padding-box; }

.app-notification:backdrop, .app-notification.frame:backdrop { background-image: none; transition: 200ms ease-out; }

.app-notification border, .app-notification.frame border { border: none; }

/************* Expanders * */
expander title > arrow { min-width: 16px; min-height: 16px; -gtk-icon-source: -gtk-icontheme("pan-end-symbolic"); }

expander title > arrow:dir(rtl) { -gtk-icon-source: -gtk-icontheme("pan-end-symbolic-rtl"); }

expander title > arrow:hover { color: #748489; }

expander title > arrow:disabled { color: #929595; }

expander title > arrow:disabled:backdrop { color: #d4cfca; }

expander title > arrow:checked { -gtk-icon-source: -gtk-icontheme("pan-down-symbolic"); }

/************ Calendar * */
calendar { color: black; border: 1px solid #cdc7c2; }

calendar:selected { border-radius: 3px; }

calendar.header { border-bottom-color: rgba(0, 0, 0, 0.1); }

calendar.header:backdrop { border-bottom-color: rgba(0, 0, 0, 0.1); }

calendar.button { color: rgba(46, 52, 54, 0.45); }

calendar.button:hover { color: #2e3436; }

calendar.button:backdrop { color: rgba(146, 149, 149, 0.45); }

calendar.button:disabled { color: rgba(146, 149, 149, 0.45); }

calendar.highlight { color: #929595; }

calendar.highlight:backdrop { color: #d4cfca; }

calendar:backdrop { color: #323232; border-color: #d5d0cc; }

calendar:indeterminate { color: alpha(currentColor,0.1); }

/*********** Dialogs * */
messagedialog .titlebar { min-height: 20px; background-image: none; background-color: #f6f5f4; border-style: none; border-top-left-radius: 7px; border-top-right-radius: 7px; }

messagedialog.csd.background { border-bottom-left-radius: 9px; border-bottom-right-radius: 9px; }

messagedialog.csd .dialog-action-area button { padding: 10px 14px; border-radius: 0; border-left-style: solid; border-right-style: none; border-bottom-style: none; }

messagedialog.csd .dialog-action-area button:first-child:not(:only-child) { border-left-style: none; border-bottom-left-radius: 7px; -gtk-outline-bottom-left-radius: 7px; -gtk-outline-top-left-radius: 0px; -gtk-outline-top-right-radius: 0px; -gtk-outline-bottom-right-radius: 0px; }

messagedialog.csd .dialog-action-area button:last-child:not(:only-child) { border-bottom-right-radius: 7px; -gtk-outline-bottom-right-radius: 7px; -gtk-outline-top-right-radius: 0px; -gtk-outline-bottom-left-radius: 0px; -gtk-outline-top-left-radius: 0px; }

messagedialog.csd .dialog-action-area button:only-child { border-top-right-radius: 0; border-top-left-radius: 0; border-bottom-left-radius: 7px; border-bottom-right-radius: 7px; -gtk-outline-top-right-radius: 0px; -gtk-outline-top-left-radius: 0px; -gtk-outline-bottom-left-radius: 7px; -gtk-outline-bottom-right-radius: 7px; }

filechooser .dialog-action-box { border-top: 1px solid #cdc7c2; }

filechooser .dialog-action-box:backdrop { border-top-color: #d5d0cc; }

filechooser #pathbarbox { border-bottom: 1px solid #f6f5f4; }

filechooserbutton:drop(active) { box-shadow: none; border-color: transparent; }

/*********** Sidebar * */
.sidebar { border-style: none; background-color: #fbfafa; }

stacksidebar.sidebar:dir(ltr) list, stacksidebar.sidebar.left list, stacksidebar.sidebar.left:dir(rtl) list, .sidebar:not(separator):dir(ltr), .sidebar.left:not(separator) { border-right: 1px solid #cdc7c2; border-left-style: none; }

stacksidebar.sidebar:dir(rtl) list, stacksidebar.sidebar.right list, .sidebar:not(separator):dir(rtl), .sidebar.right:not(separator) { border-left: 1px solid #cdc7c2; border-right-style: none; }

.sidebar:backdrop { background-color: #f9f9f8; border-color: #d5d0cc; transition: 200ms ease-out; }

.sidebar list { background-color: transparent; }

paned .sidebar.left, paned .sidebar.right, paned .sidebar.left:dir(rtl), paned .sidebar:dir(rtl), paned .sidebar:dir(ltr), paned .sidebar { border-style: none; }

stacksidebar row { padding: 10px 4px; }

stacksidebar row > label { padding-left: 6px; padding-right: 6px; }

stacksidebar row.needs-attention > label { background-size: 6px 6px, 0 0; }

separator.sidebar { background-color: #cdc7c2; }

separator.sidebar:backdrop { background-color: #d5d0cc; }

separator.sidebar.selection-mode, .selection-mode separator.sidebar { background-color: #15539e; }

/**************** File chooser * */
row image.sidebar-icon { opacity: 0.7; }

placessidebar > viewport.frame { border-style: none; }

placessidebar row { min-height: 36px; padding: 0px; }

placessidebar row > revealer { padding: 0 14px; }

placessidebar row:selected { color: #ffffff; }

placessidebar row:disabled { color: #929595; }

placessidebar row:backdrop { color: #929595; }

placessidebar row:backdrop:selected { color: #fcfcfc; }

placessidebar row:backdrop:disabled { color: #d4cfca; }

placessidebar row image.sidebar-icon:dir(ltr) { padding-right: 8px; }

placessidebar row image.sidebar-icon:dir(rtl) { padding-left: 8px; }

placessidebar row label.sidebar-label:dir(ltr) { padding-right: 2px; }

placessidebar row label.sidebar-label:dir(rtl) { padding-left: 2px; }

button.sidebar-button { min-height: 26px; min-width: 26px; margin-top: 3px; margin-bottom: 3px; padding: 0; border-radius: 100%; -gtk-outline-radius: 100%; }

button.sidebar-button:not(:hover):not(:active) > image, button.sidebar-button:backdrop > image { opacity: 0.7; }

placessidebar row:selected:active { box-shadow: none; }

placessidebar row.sidebar-placeholder-row { padding: 0 8px; min-height: 2px; background-image: image(#4e9a06); background-clip: content-box; }

placessidebar row.sidebar-new-bookmark-row { color: #3584e4; }

placessidebar row:drop(active):not(:disabled) { color: #4e9a06; box-shadow: inset 0 1px #4e9a06, inset 0 -1px #4e9a06; }

placessidebar row:drop(active):not(:disabled):selected { color: #ffffff; background-color: #4e9a06; }

placesview .server-list-button > image { transition: 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); -gtk-icon-transform: rotate(0turn); }

placesview .server-list-button:checked > image { transition: 200ms cubic-bezier(0.25, 0.46, 0.45, 0.94); -gtk-icon-transform: rotate(-0.5turn); }

placesview row.activatable:hover { background-color: transparent; }

placesview > actionbar > revealer > box > label { padding-left: 8px; padding-right: 8px; }

/********* Paned * */
paned > separator { min-width: 1px; min-height: 1px; -gtk-icon-source: none; border-style: none; background-color: transparent; background-image: image(#cdc7c2); background-size: 1px 1px; }

paned > separator:selected { background-image: image(#3584e4); }

paned > separator:backdrop { background-image: image(#d5d0cc); }

paned > separator.wide { min-width: 5px; min-height: 5px; background-color: #f6f5f4; background-image: image(#cdc7c2), image(#cdc7c2); background-size: 1px 1px, 1px 1px; }

paned > separator.wide:backdrop { background-color: #f6f5f4; background-image: image(#d5d0cc), image(#d5d0cc); }

paned.horizontal > separator { background-repeat: repeat-y; }

paned.horizontal > separator:dir(ltr) { margin: 0 -8px 0 0; padding: 0 8px 0 0; background-position: left; }

paned.horizontal > separator:dir(rtl) { margin: 0 0 0 -8px; padding: 0 0 0 8px; background-position: right; }

paned.horizontal > separator.wide { margin: 0; padding: 0; background-repeat: repeat-y, repeat-y; background-position: left, right; }

paned.vertical > separator { margin: 0 0 -8px 0; padding: 0 0 8px 0; background-repeat: repeat-x; background-position: top; }

paned.vertical > separator.wide { margin: 0; padding: 0; background-repeat: repeat-x, repeat-x; background-position: bottom, top; }

/************** GtkInfoBar * */
infobar { border-style: none; }

infobar.action:hover > revealer > box { background-color: #f4ebe1; border-bottom: 1px solid #d8d4d0; }

infobar.info, infobar.question, infobar.warning, infobar.error { text-shadow: none; }

infobar.info:backdrop > revealer > box, infobar.info > revealer > box, infobar.question:backdrop > revealer > box, infobar.question > revealer > box, infobar.warning:backdrop > revealer > box, infobar.warning > revealer > box, infobar.error:backdrop > revealer > box, infobar.error > revealer > box { background-color: #f1e6d9; border-bottom: 1px solid #d8d4d0; }

infobar.info:backdrop > revealer > box label, infobar.info:backdrop > revealer > box, infobar.info > revealer > box label, infobar.info > revealer > box, infobar.question:backdrop > revealer > box label, infobar.question:backdrop > revealer > box, infobar.question > revealer > box label, infobar.question > revealer > box, infobar.warning:backdrop > revealer > box label, infobar.warning:backdrop > revealer > box, infobar.warning > revealer > box label, infobar.warning > revealer > box, infobar.error:backdrop > revealer > box label, infobar.error:backdrop > revealer > box, infobar.error > revealer > box label, infobar.error > revealer > box { color: #2e3436; }

infobar.info:backdrop, infobar.question:backdrop, infobar.warning:backdrop, infobar.error:backdrop { text-shadow: none; }

infobar.info button, infobar.question button, infobar.warning button, infobar.error button { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; background-image: linear-gradient(to top, #edebe9 2px, #f6f5f4); text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); }

infobar.info button:hover, infobar.question button:hover, infobar.warning button:hover, infobar.error button:hover { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px); }

infobar.info button:active, infobar.info button:checked, infobar.question button:active, infobar.question button:checked, infobar.warning button:active, infobar.warning button:checked, infobar.error button:active, infobar.error button:checked { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; background-image: image(#d6d1cd); box-shadow: inset 0 1px rgba(255, 255, 255, 0); text-shadow: none; -gtk-icon-shadow: none; }

infobar.info button:disabled, infobar.question button:disabled, infobar.warning button:disabled, infobar.error button:disabled { color: #929595; border-color: #cdc7c2; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

infobar.info button:backdrop, infobar.question button:backdrop, infobar.warning button:backdrop, infobar.error button:backdrop { color: #929595; border-color: #d5d0cc; background-image: image(#f6f5f4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #dfdcd8; }

infobar.info button:backdrop:disabled, infobar.question button:backdrop:disabled, infobar.warning button:backdrop:disabled, infobar.error button:backdrop:disabled { color: #d4cfca; border-color: #d5d0cc; background-image: image(#faf9f8); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); border-color: #dfdcd8; }

infobar.info button:backdrop label, infobar.info button:backdrop, infobar.info button label, infobar.info button, infobar.question button:backdrop label, infobar.question button:backdrop, infobar.question button label, infobar.question button, infobar.warning button:backdrop label, infobar.warning button:backdrop, infobar.warning button label, infobar.warning button, infobar.error button:backdrop label, infobar.error button:backdrop, infobar.error button label, infobar.error button { color: #2e3436; }

infobar.info selection, infobar.question selection, infobar.warning selection, infobar.error selection { background-color: #dfdcd8; }

infobar.info *:link, infobar.question *:link, infobar.warning *:link, infobar.error *:link { color: #1b6acb; }

/************ Tooltips * */
tooltip { padding: 4px; /* not working */ border-radius: 5px; box-shadow: none; text-shadow: 0 1px black; }

tooltip.background { background-color: rgba(0, 0, 0, 0.8); background-clip: padding-box; border: 1px solid rgba(255, 255, 255, 0.1); }

tooltip decoration { background-color: transparent; }

tooltip * { padding: 4px; background-color: transparent; color: white; }

/***************** Color Chooser * */
colorswatch:drop(active), colorswatch { border-style: none; }

colorswatch.top { border-top-left-radius: 5.5px; border-top-right-radius: 5.5px; }

colorswatch.top overlay { border-top-left-radius: 5px; border-top-right-radius: 5px; }

colorswatch.bottom { border-bottom-left-radius: 5.5px; border-bottom-right-radius: 5.5px; }

colorswatch.bottom overlay { border-bottom-left-radius: 5px; border-bottom-right-radius: 5px; }

colorswatch.left, colorswatch:first-child:not(.top) { border-top-left-radius: 5.5px; border-bottom-left-radius: 5.5px; }

colorswatch.left overlay, colorswatch:first-child:not(.top) overlay { border-top-left-radius: 5px; border-bottom-left-radius: 5px; }

colorswatch.right, colorswatch:last-child:not(.bottom) { border-top-right-radius: 5.5px; border-bottom-right-radius: 5.5px; }

colorswatch.right overlay, colorswatch:last-child:not(.bottom) overlay { border-top-right-radius: 5px; border-bottom-right-radius: 5px; }

colorswatch.dark { outline-color: rgba(255, 255, 255, 0.6); }

colorswatch.dark overlay { color: white; }

colorswatch.dark overlay:hover { border-color: rgba(0, 0, 0, 0.8); }

colorswatch.dark overlay:backdrop { color: rgba(255, 255, 255, 0.5); }

colorswatch.light { outline-color: rgba(0, 0, 0, 0.6); }

colorswatch.light overlay { color: black; }

colorswatch.light overlay:hover { border-color: rgba(0, 0, 0, 0.5); }

colorswatch.light overlay:backdrop { color: rgba(0, 0, 0, 0.5); }

colorswatch:drop(active) { box-shadow: none; }

colorswatch.light:drop(active) overlay { border-color: #4e9a06; box-shadow: inset 0 0 0 2px #3d7805, inset 0 0 0 1px #4e9a06; }

colorswatch.dark:drop(active) overlay { border-color: #4e9a06; box-shadow: inset 0 0 0 2px rgba(0, 0, 0, 0.3), inset 0 0 0 1px #4e9a06; }

colorswatch overlay { border: 1px solid rgba(0, 0, 0, 0.3); }

colorswatch overlay:hover { box-shadow: inset 0 1px rgba(255, 255, 255, 0.4), inset 0 -1px rgba(0, 0, 0, 0.2); }

colorswatch overlay:backdrop, colorswatch overlay:backdrop:hover { border-color: rgba(0, 0, 0, 0.3); box-shadow: none; }

colorswatch#add-color-button { border-radius: 5px 5px 0 0; }

colorswatch#add-color-button:only-child { border-radius: 5px; }

colorswatch#add-color-button overlay { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; background-image: linear-gradient(to top, #edebe9 2px, #f6f5f4); text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); }

colorswatch#add-color-button overlay:hover { color: #2e3436; outline-color: rgba(46, 52, 54, 0.3); border-color: #cdc7c2; border-bottom-color: #bfb8b1; text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); box-shadow: inset 0 1px white, 0 1px 2px rgba(0, 0, 0, 0.07); background-image: linear-gradient(to top, #f6f5f4, #f8f8f7 1px); }

colorswatch#add-color-button overlay:backdrop { color: #929595; border-color: #d5d0cc; background-image: image(#f6f5f4); text-shadow: none; -gtk-icon-shadow: none; box-shadow: inset 0 1px rgba(255, 255, 255, 0); }

colorswatch:disabled { opacity: 0.5; }

colorswatch:disabled overlay { border-color: rgba(0, 0, 0, 0.6); box-shadow: none; }

row:selected colorswatch { box-shadow: 0 0 0 2px #ffffff; }

colorswatch#editor-color-sample { border-radius: 4px; }

colorswatch#editor-color-sample overlay { border-radius: 4.5px; }

colorchooser .popover.osd { border-radius: 5px; }

/******** Misc * */
.content-view { background-color: #e6e3e0; }

.content-view:hover { -gtk-icon-effect: highlight; }

.content-view:backdrop { background-color: #e6e3e0; }

.osd .scale-popup button.flat { border-style: none; border-radius: 5px; }

.scale-popup button:hover { background-color: rgba(46, 52, 54, 0.1); border-radius: 5px; }

/********************** Window Decorations * */
decoration { border-radius: 8px 8px 0 0; border-width: 0px; box-shadow: 0 3px 9px 1px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(0, 0, 0, 0.23); margin: 10px; }

decoration:backdrop { box-shadow: 0 3px 9px 1px transparent, 0 2px 6px 2px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(0, 0, 0, 0.18); transition: 200ms ease-out; }

.maximized decoration, .fullscreen decoration, .tiled decoration, .tiled-top decoration, .tiled-right decoration, .tiled-bottom decoration, .tiled-left decoration { border-radius: 0; }

.popup decoration { box-shadow: none; }

.csd decoration { background-color: black; }

.ssd decoration { box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.23); }

.csd.popup decoration { border-radius: 5px; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(0, 0, 0, 0.13); }

tooltip.csd decoration { border-radius: 5px; box-shadow: none; }

messagedialog.csd decoration { border-radius: 8px; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(0, 0, 0, 0.13); }

.solid-csd decoration { margin: 0; padding: 4px; background-color: #cdc7c2; border: solid 1px #cdc7c2; border-radius: 0; box-shadow: inset 0 0 0 3px white, inset 0 1px rgba(255, 255, 255, 0.8); }

.solid-csd decoration:backdrop { box-shadow: inset 0 0 0 3px #f6f5f4, inset 0 1px rgba(255, 255, 255, 0.8); }

button.titlebutton { text-shadow: 0 1px rgba(255, 255, 255, 0.769231); -gtk-icon-shadow: 0 1px rgba(255, 255, 255, 0.769231); }

button.titlebutton:not(.appmenu) { border-radius: 9999px; padding: 6px; margin: 0 2px; min-width: 0; min-height: 0; }

button.titlebutton:backdrop { -gtk-icon-shadow: none; }

.selection-mode headerbar button.titlebutton, .selection-mode .titlebar button.titlebutton, headerbar.selection-mode button.titlebutton, .titlebar.selection-mode button.titlebutton { text-shadow: 0 -1px rgba(0, 0, 0, 0.559216); -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.559216); }

.selection-mode headerbar button.titlebutton:backdrop, .selection-mode .titlebar button.titlebutton:backdrop, headerbar.selection-mode button.titlebutton:backdrop, .titlebar.selection-mode button.titlebutton:backdrop { -gtk-icon-shadow: none; }

.view:selected:focus, .view:selected, .view text:selected:focus, textview text:selected:focus, .view text:selected, textview text:selected, .view text selection:focus, .view text selection, textview text selection:focus, textview text selection, iconview:selected:focus, iconview:selected, iconview text selection:focus, .view text selection, iconview text selection, flowbox flowboxchild:selected, entry selection, modelbutton.flat:selected, .menuitem.button.flat:selected, spinbutton:not(.vertical) selection, treeview.view:selected:focus, treeview.view:selected, row:selected, calendar:selected { background-color: #3584e4; }

label:selected, .selection-mode button.titlebutton, .view:selected:focus, .view:selected, .view text:selected:focus, textview text:selected:focus, .view text:selected, textview text:selected, .view text selection:focus, .view text selection, textview text selection:focus, textview text selection, iconview:selected:focus, iconview:selected, iconview text selection:focus, .view text selection, iconview text selection, flowbox flowboxchild:selected, entry selection, modelbutton.flat:selected, .menuitem.button.flat:selected, spinbutton:not(.vertical) selection, treeview.view:selected:focus, treeview.view:selected, row:selected, calendar:selected { color: #ffffff; }

label:disabled selection, label:disabled:selected, .selection-mode button.titlebutton:disabled, .view:disabled:selected, textview text:disabled:selected:focus, .view text:disabled:selected, textview text:disabled:selected, .view text selection:disabled, textview text selection:disabled:focus, textview text selection:disabled, iconview:disabled:selected:focus, iconview:disabled:selected, iconview text selection:disabled:focus, iconview text selection:disabled, flowbox flowboxchild:disabled:selected, entry selection:disabled, modelbutton.flat:disabled:selected, .menuitem.button.flat:disabled:selected, spinbutton:not(.vertical) selection:disabled, treeview.view:disabled:selected, row:disabled:selected, calendar:disabled:selected { color: #9ac2f2; }

label:backdrop selection, label:backdrop:selected, .selection-mode button.titlebutton:backdrop, .view:backdrop:selected, textview text:backdrop:selected:focus, .view text:backdrop:selected, textview text:backdrop:selected, .view text selection:backdrop, textview text selection:backdrop:focus, textview text selection:backdrop, iconview:backdrop:selected:focus, iconview:backdrop:selected, iconview text selection:backdrop:focus, iconview text selection:backdrop, flowbox flowboxchild:backdrop:selected, entry selection:backdrop, modelbutton.flat:backdrop:selected, .menuitem.button.flat:backdrop:selected, spinbutton:not(.vertical) selection:backdrop, treeview.view:backdrop:selected, row:backdrop:selected, calendar:backdrop:selected { color: #fcfcfc; }

label:backdrop selection:disabled, label:backdrop:disabled:selected, .selection-mode button.titlebutton:backdrop:disabled, .view:backdrop:disabled:selected, .view text:backdrop:disabled:selected, textview text:backdrop:disabled:selected, .view text selection:backdrop:disabled, textview text selection:backdrop:disabled, iconview:backdrop:disabled:selected, iconview text selection:backdrop:disabled, flowbox flowboxchild:backdrop:disabled:selected, entry selection:backdrop:disabled, modelbutton.flat:backdrop:disabled:selected, .menuitem.button.flat:backdrop:disabled:selected, spinbutton:not(.vertical) selection:backdrop:disabled, row:backdrop:disabled:selected, calendar:backdrop:disabled:selected { color: #71a8eb; }

.monospace { font-family: monospace; }

/********************** Touch Copy & Paste * */
cursor-handle { background-color: transparent; background-image: none; box-shadow: none; border-style: none; }

cursor-handle.top:dir(ltr), cursor-handle.bottom:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/text-select-start.png"), url("assets/text-select-start@2.png")); padding-left: 10px; }

cursor-handle.bottom:dir(ltr), cursor-handle.top:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/text-select-end.png"), url("assets/text-select-end@2.png")); padding-right: 10px; }

cursor-handle.insertion-cursor:dir(ltr), cursor-handle.insertion-cursor:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above.png"), url("assets/slider-horz-scale-has-marks-above@2.png")); }

cursor-handle.top:hover:dir(ltr), cursor-handle.bottom:hover:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/text-select-start-hover.png"), url("assets/text-select-start-hover@2.png")); padding-left: 10px; }

cursor-handle.bottom:hover:dir(ltr), cursor-handle.top:hover:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/text-select-end-hover.png"), url("assets/text-select-end-hover@2.png")); padding-right: 10px; }

cursor-handle.insertion-cursor:hover:dir(ltr), cursor-handle.insertion-cursor:hover:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above-hover.png"), url("assets/slider-horz-scale-has-marks-above-hover@2.png")); }

cursor-handle.top:active:dir(ltr), cursor-handle.bottom:active:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/text-select-start-active.png"), url("assets/text-select-start-active@2.png")); padding-left: 10px; }

cursor-handle.bottom:active:dir(ltr), cursor-handle.top:active:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/text-select-end-active.png"), url("assets/text-select-end-active@2.png")); padding-right: 10px; }

cursor-handle.insertion-cursor:active:dir(ltr), cursor-handle.insertion-cursor:active:dir(rtl) { -gtk-icon-source: -gtk-scaled(url("assets/slider-horz-scale-has-marks-above-active.png"), url("assets/slider-horz-scale-has-marks-above-active@2.png")); }

.context-menu { font: initial; }

.keycap { min-width: 20px; min-height: 25px; margin-top: 2px; padding-bottom: 3px; padding-left: 6px; padding-right: 6px; color: #2e3436; background-color: #ffffff; border: 1px solid; border-color: #e1dedb; border-radius: 5px; box-shadow: inset 0 -3px #f8f7f6; font-size: smaller; }

.keycap:backdrop { background-color: #fcfcfc; color: #929595; transition: 200ms ease-out; }

:not(decoration):not(window):drop(active):focus, :not(decoration):not(window):drop(active) { border-color: #4e9a06; box-shadow: inset 0 0 0 1px #4e9a06; caret-color: #4e9a06; }

stackswitcher button.text-button { min-width: 100px; }

stackswitcher button.circular, stackswitcher button.text-button.circular { min-width: 32px; min-height: 32px; padding: 0; }

/************* App Icons * */
/* Outline for low res icons */
.lowres-icon { -gtk-icon-shadow: 0 -1px rgba(0, 0, 0, 0.05), 1px 0 rgba(0, 0, 0, 0.1), 0 1px rgba(0, 0, 0, 0.3), -1px 0 rgba(0, 0, 0, 0.1); }

/* Dropshadow for large icons */
.icon-dropshadow { -gtk-icon-shadow: 0 1px 12px rgba(0, 0, 0, 0.05), 0 -1px rgba(0, 0, 0, 0.05), 1px 0 rgba(0, 0, 0, 0.1), 0 1px rgba(0, 0, 0, 0.3), -1px 0 rgba(0, 0, 0, 0.1); }

/********* Emoji * */
popover.emoji-picker { padding-left: 0; padding-right: 0; }

popover.emoji-picker entry.search { margin: 3px 5px 5px 5px; }

button.emoji-section { border-color: transparent; border-width: 3px; border-style: none none solid; border-radius: 0; margin: 2px 4px 2px 4px; padding: 3px 0 0; min-width: 32px; min-height: 28px; /* reset props inherited from the button style */ background: none; box-shadow: none; text-shadow: none; outline-offset: -5px; }

button.emoji-section:first-child { margin-left: 7px; }

button.emoji-section:last-child { margin-right: 7px; }

button.emoji-section:backdrop:not(:checked) { border-color: transparent; }

button.emoji-section:hover { border-color: #cdc7c2; }

button.emoji-section:checked { border-color: #3584e4; }

button.emoji-section label { padding: 0; opacity: 0.55; }

button.emoji-section:hover label { opacity: 0.775; }

button.emoji-section:checked label { opacity: 1; }

popover.emoji-picker .emoji { font-size: x-large; padding: 6px; }

popover.emoji-picker .emoji :hover { background: #3584e4; border-radius: 6px; }

popover.emoji-completion arrow { border: none; background: none; }

popover.emoji-completion contents row box { padding: 2px 10px; }

popover.emoji-completion .emoji:hover { background: white; }

/* GTK NAMED COLORS ---------------- use responsibly! */
/*
widget text/foreground color */
@define-color theme_fg_color #2e3436;
/*
text color for entries, views and content in general */
@define-color theme_text_color black;
/*
widget base background color */
@define-color theme_bg_color #f6f5f4;
/*
text widgets and the like base background color */
@define-color theme_base_color #ffffff;
/*
base background color of selections */
@define-color theme_selected_bg_color #3584e4;
/*
text/foreground color of selections */
@define-color theme_selected_fg_color #ffffff;
/*
base background color of insensitive widgets */
@define-color insensitive_bg_color #faf9f8;
/*
text foreground color of insensitive widgets */
@define-color insensitive_fg_color #929595;
/*
insensitive text widgets and the like base background color */
@define-color insensitive_base_color #ffffff;
/*
widget text/foreground color on backdrop windows */
@define-color theme_unfocused_fg_color #929595;
/*
text color for entries, views and content in general on backdrop windows */
@define-color theme_unfocused_text_color black;
/*
widget base background color on backdrop windows */
@define-color theme_unfocused_bg_color #f6f5f4;
/*
text widgets and the like base background color on backdrop windows */
@define-color theme_unfocused_base_color #fcfcfc;
/*
base background color of selections on backdrop windows */
@define-color theme_unfocused_selected_bg_color #3584e4;
/*
text/foreground color of selections on backdrop windows */
@define-color theme_unfocused_selected_fg_color #ffffff;
/*
insensitive color on backdrop windows*/
@define-color unfocused_insensitive_color #d4cfca;
/*
widgets main borders color */
@define-color borders #cdc7c2;
/*
widgets main borders color on backdrop windows */
@define-color unfocused_borders #d5d0cc;
/*
these are pretty self explicative */
@define-color warning_color #f57900;
@define-color error_color #cc0000;
@define-color success_color #33d17a;
/*
these colors are exported for the window manager and shouldn't be used in applications,
read if you used those and something break with a version upgrade you're on your own... */
@define-color wm_title shade(#2e3436, 1.8);
@define-color wm_unfocused_title #929595;
@define-color wm_highlight rgba(255, 255, 255, 0.8);
@define-color wm_borders_edge rgba(255, 255, 255, 0.8);
@define-color wm_bg_a shade(#f6f5f4, 1.2);
@define-color wm_bg_b #f6f5f4;
@define-color wm_shadow alpha(black, 0.35);
@define-color wm_border alpha(black, 0.18);
@define-color wm_button_hover_color_a shade(#f6f5f4, 1.3);
@define-color wm_button_hover_color_b #f6f5f4;
@define-color wm_button_active_color_a shade(#f6f5f4, 0.85);
@define-color wm_button_active_color_b shade(#f6f5f4, 0.89);
@define-color wm_button_active_color_c shade(#f6f5f4, 0.9);
@define-color content_view_bg #ffffff;
`