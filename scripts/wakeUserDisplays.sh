#!/bin/env bash

declare screen

ARGONE_REGEX="[0-9]+"
SCREEN_REGEX="\:([0-9])+"

if [ ! -z ${1+x} ]; then
    if [[ ! $1 =~ $ARGONE_REGEX ]]; then
	echo "Error, \"$1\" is not a valid display label"
	exit 1
    fi

    echo "Trying to wake display $1"
    screen=$1
fi

w -sh | while read -r user tty from rest; do
    [[ ! $from =~ $SCREEN_REGEX ]] && continue

    if [ ! -z ${screen+x} ] && [ ${BASH_REMATCH[1]} != $screen ]; then
	continue
    fi	
    
    echo "Waking display $from for user $user"
    sudo -u $user xset -display $from dpms force on
done
