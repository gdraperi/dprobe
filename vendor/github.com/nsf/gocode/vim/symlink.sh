#!/bin/sh
cd "$***REMOVED***0%/****REMOVED***"
ROOTDIR=`pwd`
mkdir -p "$HOME/.vim/autoload"
mkdir -p "$HOME/.vim/ftplugin/go"
ln -fs "$ROOTDIR/autoload/gocomplete.vim" "$HOME/.vim/autoload/"
ln -fs "$ROOTDIR/ftplugin/go/gocomplete.vim" "$HOME/.vim/ftplugin/go/"
