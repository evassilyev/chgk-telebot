#!/usr/bin/env bash

# Usage: build.sh {VERSION}

if [ $# -ne 1 ]; then
    echo "Enter the building version"
    echo "Usage: build.sh {VERSION}"
else
    macversionpath="Builds/$1/MacOS"
    winversionpath="Builds/$1/Windows"
    nixversionpath="Builds/$1/Linux"

    echo -n "Building MacOS version... "
    GOOS=darwin GOARCH=amd64 go build
    mkdir -p $macversionpath
    mv chgk-telebot $macversionpath
    echo "Done"

    echo -n "Building Windows version... "
    GOOS=windows GOARCH=amd64 go build
    mkdir -p $winversionpath
    mv chgk-telebot.exe $winversionpath
    echo "Done"

    echo -n "Building Linux version... "
    GOOS=linux GOARCH=amd64 go build
    mkdir -p $nixversionpath
    mv chgk-telebot $nixversionpath
    echo "Done"
fi