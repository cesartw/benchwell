project('io.benchwell', 'vala', 'c')

i18n = import('i18n')
add_global_arguments('-DGETTEXT_PACKAGE="@0@"'.format (meson.project_name()), language: 'c')

if get_option('profile') == 'development'
	application_id = '@0@Devel'.format(meson.project_name())
	vala_args += ['-D', 'IS_DEVEL']
else
	application_id = meson.project_name()
endif

conf_data = configuration_data()
conf_data.set('PROJECT_NAME', application_id)
conf_data.set('GETTEXT_PACKAGE', meson.project_name())
conf_data.set('VERSION', meson.project_version())

config_headers = configure_file(
    input: 'config.vala.in',
    output: 'config.vala',
    configuration: conf_data
)

gnome = import('gnome')
asresources = gnome.compile_resources(
    'as-resources', 'data/assets.gresource.xml',
    source_dir: 'data',
    c_name: 'as'
)

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
	'src/Widgets/SourceView.vala',

	'src/Views/DBDatabase.vala',
	'src/Views/DBConnect.vala',
	'src/Views/DBData.vala',

	'src/Config.vala',
	'src/Constants.vala',
	'src/Application.vala',
	'src/ApplicationWindow.vala',

	config_headers,
	asresources,

	dependencies: [
		dependency('gobject-2.0'),
		dependency('glib-2.0'),
		dependency('gio-2.0'),
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
subdir('po')

meson.add_install_script('build-aux/post_install.py')
