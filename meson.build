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

	'src/sql_driver.vala',
	'src/sql_mysql.vala',
	'src/sql_engine.vala',

	'src/widget_image.vala',
	'src/widget_menu_item.vala',
	'src/widget_button.vala',
	'src/widget_tab.vala',
	'src/widget_source_view.vala',
	'src/widget_button_opt.vala',
	'src/widget_keyvalues.vala',

	'src/database_mysql_form.vala',
	'src/database_sqlite_form.vala',
	'src/database.vala',
	'src/database_connect.vala',
	'src/database_data.vala',
	'src/database_condition.vala',
	'src/database_result_view.vala',
	'src/database_table.vala',
	'src/database_tables.vala',

	'src/http.vala',
	'src/http_addressbar.vala',
	'src/http_sidebar.vala',

	'src/database_service.vala',

	'src/config/config.vala',
	'src/config/plugin.vala',
	'src/config/http.vala',
	'src/config/environment.vala',
	'src/config/keyvalue.vala',
	'src/config/connection.vala',

	'src/constants.vala',
	'src/application_window.vala',
	'src/environment.vala',
	'src/settings_panel.vala',

	'src/application.vala',
	#'src/http_importer.vala',
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
		dependency('javascriptcoregtk-4.0'),
		dependency('gio-2.0'),
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
