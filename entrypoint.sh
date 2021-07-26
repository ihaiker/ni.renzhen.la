#!/usr/bin/env sh

set -e

echo "copy /build to /source"
cp -r /build /source

echo "cd /source"
cd /source

echo "generator document."
go run scripts/main.go

echo "mkdocs $@"
mkdocs $@

cp -r /source/site /build/site
