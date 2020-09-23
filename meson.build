project('io.benchwell', 'vala', 'c')

executable(
	meson.project_name(),

	'src/SQL/Driver.vala',
	'src/SQL/Mysql.vala',
	'src/SQL/Engine.vala',

	'src/Widgets/Image.vala',
	'src/Widgets/MenuItem.vala',
	'src/Widgets/Button.vala',
	'src/Widgets/Tab.vala',
	'src/Widgets/MysqlForm.vala',
	'src/Widgets/SqliteForm.vala',

	'src/Views/DBDatabase.vala',
	'src/Views/DBConnect.vala',
	'src/Views/DBData.vala',

	'src/Config.vala',
	'src/Constants.vala',
	'src/Application.vala',
	'src/ApplicationWindow.vala',

	dependencies: [
		dependency('gobject-2.0'),
		dependency('glib-2.0'),
		dependency('gtk+-3.0'),
		dependency('gtksourceview-3.0'),
		dependency('sqlite3'),
		dependency('libsecret-1'),
		meson.get_compiler('c').find_library('mysqlclient'),
		meson.get_compiler('c').find_library('mariadb'),
		meson.get_compiler('vala').find_library('mysql')
	],
	install: true,
)

subdir('data')

meson.add_install_script('meson/post_install.py')
