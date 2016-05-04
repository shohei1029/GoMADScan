#!/bin/sh
set  -eu

brew install git
brew install go
brew install pkgconfig
export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
brew install gtk+
brew install gtksourceview

mkdir ~/go
export GOPATH=~/go
go get github.com/carushi/GoMADScan
