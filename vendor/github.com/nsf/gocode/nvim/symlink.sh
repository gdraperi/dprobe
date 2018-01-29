#!/bin/sh
cd "$***REMOVED***0%/****REMOVED***"
ROOTDIR=`pwd`
mkdir -p "$HOME/.config/nvim/autoload"
mkdir -p "$HOME/.config/nvim/ftplugin/go"
ln -fs "$ROOTDIR/autoload/gocomplete.vim" "$HOME/.config/nvim/autoload/"
ln -fs "$ROOTDIR/ftplugin/go/gocomplete.vim" "$HOME/.config/nvim/ftplugin/go/"
