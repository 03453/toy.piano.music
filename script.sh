#!/usr/bin/env zsh

sed -i .bak -e "s/<chord>HERE<\/chord>/<chord\/>/g" -e "s/<rest>HERE<\/rest>/<rest\/>/g" -e "s/<dot>HERE<\/dot>/<dot\/>/g" score.xml
rm score.xml.bak
