project('io.benchwell', 'vala', 'c')

i18n = import('i18n')
add_global_arguments('-DGETTEXT_PACKAGE="@0@"'.format (meson.project_name()), language: 'c')
vala_args = ['--pkg', 'posix']
if get_option('profile') == 'development'
	application_id = '@0@Devel'.format(meson.project_name())
	vala_args += ['-D', '-g', '--save-temps', 'IS_DEVEL']
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

vapi_dir = meson.current_source_dir() / 'vapi'
add_project_arguments(['--vapidir', vapi_dir], language: 'vala')

executable(
	meson.project_name(),

	'src/backend/sql_driver.vala',
	'src/backend/sql_mysql.vala',
	'src/backend/sql_engine.vala',

	'src/widgets/image.vala',
	'src/widgets/menu_item.vala',
	'src/widgets/button.vala',
	'src/widgets/tab.vala',
	'src/widgets/source_view.vala',
	'src/widgets/button_opt.vala',
	'src/widgets/keyvalues.vala',

	'src/database/mysql_form.vala',
	'src/database/sqlite_form.vala',
	'src/database/database.vala',
	'src/database/connect.vala',
	'src/database/data.vala',
	'src/database/condition.vala',
	'src/database/result_view.vala',
	'src/database/table.vala',
	'src/database/tables.vala',

	'src/http/http.vala',
	'src/http/addressbar.vala',
	'src/http/sidebar.vala',

	'src/services/database.vala',

	'src/config/config.vala',
	'src/config/http.vala',
	'src/config/environment.vala',
	'src/config/keyvalue.vala',

	'src/constants.vala',
	'src/application.vala',
	'src/application_window.vala',
	'src/environment.vala',

	#'src/main.vala',

	config_headers,
	asresources,
	dependencies: [
		dependency('gobject-2.0'),
		dependency('glib-2.0'),
		dependency('gtk+-3.0'),
		dependency('gtksourceview-3.0'),
		dependency('sqlite3'),
		dependency('libsecret-1'),
		dependency('json-glib-1.0'),
		dependency('libcurl'),
		meson.get_compiler('c').find_library('mysqlclient'),
		meson.get_compiler('c').find_library('mariadb'),
		meson.get_compiler('vala').find_library('mysql')
		#meson.get_compiler('c').find_library('libcurl')
	],
	install: true,
	vala_args: vala_args,
)

subdir('data')
subdir('po')

meson.add_install_script('build-aux/post_install.py')
