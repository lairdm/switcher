#!/bin/env bash

declare screen

ARGONE_REGEX="[0-9]+"
SCREEN_REGEX="\:([0-9])+"

if [ ! -z ${2+x} ]; then
    if [[ ! $2 =~ $ARGONE_REGEX ]]; then
	echo "Error, \"$2\" is not a valid display label"
	exit 1
    fi

    echo "Trying to change ($1) display $2"
    screen=$2
fi

w -sh | while read -r user tty from rest; do
    [[ ! $from =~ $SCREEN_REGEX ]] && continue

    if [ ! -z ${screen+x} ] && [ ${BASH_REMATCH[1]} != $screen ]; then
	continue
    fi	
    
    echo "Reseting blanking for display $from for user $user"
    sudo -u $user xset -display $from dpms $1 $1 $1
    sudo -u $user xset -display $from s $1 $1
done
