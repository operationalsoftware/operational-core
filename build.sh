#!/usr/bin/env bash

# Build assets that can be served from memory
packages=("static" "components" "layout" "views")

for package in "${packages[@]}"
do
    go-assets-builder --package="$package" --strip-prefix=/"$package" --output="$package/assets.go" $(find ./"$package" -type f ! -name "*.go")
done

# error if executable name not provided as first argument
if [ -z "$1" ]; then
    echo "Error: executable name not provided"
    exit 1
fi

# Check if OS/architecture argument is provided as second argument
if [ -n "$2" ]; then
    # If argument is provided, use the specified OS/architecture
    target_os_arch="$2"
    IFS="/" read -ra os_arch <<< "$target_os_arch"
    build_os="${os_arch[0]}"
    build_arch="${os_arch[1]}"

    GOOS="$build_os" GOARCH="$build_arch" go build -o "$1"
else
    # If no argument is provided, use the original behavior
    go build -o "$1"
fi
