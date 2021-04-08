install_data(
	'gschema.xml',
	install_dir: join_paths(get_option ('datadir'), 'glib-2.0', 'schemas'),
	rename: meson.project_name() + '.gschema.xml'
)

icons_dir = join_paths(get_option('datadir'), 'icons', 'hicolor')
install_data(
		join_paths('assets', 'icons', 'scalable', '@0@.svg'.format(meson.project_name())),
		install_dir:  join_paths(icons_dir, 'scalable', 'apps'),
		rename: '@0@.svg'.format(application_id)
)

install_data(
		join_paths('assets', 'icons', 'status', 'table.svg'),
		install_dir: join_paths(icons_dir , 'scalable', 'status'),
		rename: 'bw-table.svg'
)

install_data(
		join_paths('assets', 'icons', 'status', 'table-v.svg'),
		install_dir: join_paths(icons_dir , 'scalable', 'status'),
		rename: 'bw-table-v.svg'
)

install_data(
		join_paths('assets', 'icons', 'third-party', 'mariadb.svg'),
		install_dir: join_paths(icons_dir , 'scalable', 'status'),
		rename: 'bw-mariadb.svg'
)

install_data(
		join_paths('assets', 'icons', 'third-party', 'sqlite.svg'),
		install_dir: join_paths(icons_dir , 'scalable', 'status'),
		rename: 'bw-sqlite.svg'
)

actions = [
	'add-record.svg',
	'add-tab.svg',
	'add.svg',
	'back.svg',
	'close.svg',
	'config.svg',
	'connection.svg',
	'copy.svg',
	'cowboy.svg',
	'db.svg',
	'delete-record.svg',
	'delete-table.svg',
	'directory.svg',
	'edit-table.svg',
	'filter.svg',
	'next.svg',
	'open.svg',
	'refresh.svg',
	'save-record.svg',
	'save.svg',
	'truncate.svg'
]

foreach action : actions
	install_data(
			join_paths('assets', 'icons', 'actions', 'orange', action),
			install_dir: join_paths(icons_dir , 'scalable', 'actions'),
			rename: 'bw-@0@'.format(action)
	)
	install_data(
			join_paths('assets', 'icons', 'actions', 'white', action),
			install_dir: join_paths(icons_dir , 'scalable', 'actions'),
			rename: 'bw-white-@0@'.format(action)
	)
endforeach

# Install the Desktop file
desktop_conf = configuration_data()
desktop_conf.set('icon', application_id)
desktop_conf.set('binary', application_id)

desktop_file = i18n.merge_file(
		input: configure_file(
			input: '@0@.desktop.in.in'.format(meson.project_name()),
			output: '@0@.desktop.in'.format(application_id),
			configuration: desktop_conf
		),
		output:'@0@.desktop'.format(application_id),
		po_dir: join_paths(meson.source_root(), 'po', 'extra'),
		type: 'desktop',
		install: true,
		install_dir: join_paths(get_option('datadir'), 'applications')
)

# Validate desktop file
desktop_file_validate = find_program('desktop-file-validate', required: false)
if desktop_file_validate.found()
	test(
		'validate-desktop',
		desktop_file_validate,
		args: [
			desktop_file.full_path()
		]
	)
endif

# Install the AppData file
appdata_conf = configuration_data()
appdata_conf.set('appid', application_id)

appdata_file = i18n.merge_file(
		input: configure_file(
			input: '@0@.appdata.xml.in.in'.format(meson.project_name()),
			output: '@0@.appdata.xml.in'.format(application_id),
			configuration: appdata_conf
		),
		output: '@0@.appdata.xml'.format(application_id),
		po_dir: join_paths(meson.source_root(), 'po'),
		install: true,
		install_dir: join_paths(get_option('datadir'), 'metainfo')
)

# Validate AppData file
appstream_util = find_program('appstream-util', required: false)
if appstream_util.found()
	test(
		'validate-appdata', appstream_util,
		args: [
			'validate-relax', appdata_file.full_path()
		]
	)
endif
