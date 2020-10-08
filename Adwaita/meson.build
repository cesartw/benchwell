sassc = find_program('sassc', required: false)
if not sassc.found()
  subproject('sassc')
  sassc = find_program('sassc', required: true)
endif

sassc_opts = ['-M', '-t', 'compact' ]

adwaita_scss_files = files([
  '_colors-public.scss',
  '_colors.scss',
  '_common.scss',
  '_drawing.scss',
])

adwaita_theme_variants = [
  'dark',
]

theme = custom_target('theme-light',
    input: 'gtk-container.scss',
    output: 'benchwell.css',
    command: [
      sassc, sassc_opts, '@INPUT@', '@OUTPUT@',
    ],
    depend_files: adwaita_scss_files,
  )

theme_light = custom_target('theme-dark',
	input: 'gtk-container-dark.scss',
	output: 'benchwell-dark.css',
	command: [
	  sassc, sassc_opts, '@INPUT@', '@OUTPUT@',
	],
	depend_files: adwaita_scss_files,
)
