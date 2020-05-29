#! /bin/bash

if [ ! "$(which sassc 2> /dev/null)" ]; then
   echo sassc needs to be installed to generate the css.
   exit 1
fi

SASSC_OPT="-M -t compact"

echo Generating the css...

sassc $SASSC_OPT Adwaita/gtk-contained.scss Adwaita/gtk-contained.css
sassc $SASSC_OPT Adwaita/gtk-contained-dark.scss Adwaita/gtk-contained-dark.css
