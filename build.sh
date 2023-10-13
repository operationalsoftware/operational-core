#!/usr/bin/env bash

# build assets that can be served from memory
packages=("static" "components" "layout" "views")

for package in "${packages[@]}"
do
	go-assets-builder --package=$package --strip-prefix=/$package --output=$package/assets.go $(find ./$package -type f ! -name "*.go")
done

go build
