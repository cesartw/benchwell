#!/usr/bin/env python3

import os
import platform
import subprocess

install_prefix = os.environ['MESON_INSTALL_PREFIX']
schemadir = os.path.join(install_prefix, 'share/glib-2.0/schemas')

if not os.environ.get('DESTDIR'):
    print(f'Compiling gsettings schemas {schemadir}...')
    subprocess.call(['glib-compile-schemas', schemadir])

    if platform.system() == "Linux":
        icon_cache_dir = os.path.join(install_prefix, 'share/icons/hicolor')
        print(f'Updating icon cache...{icon_cache_dir}')
        subprocess.call(['gtk-update-icon-cache', '-qtf', icon_cache_dir])

        desktop_database_dir = os.path.join(install_prefix, 'share/applications')
        print(f'Updating desktop database {desktop_database_dir}...')
        subprocess.call(['update-desktop-database', '-q', desktop_database_dir])
