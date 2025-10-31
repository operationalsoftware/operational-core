#!/usr/bin/env bash
CompileDaemon \
	-pattern="(.+\.go|.+\.c|.+\.png|.+\.svg|.+\.css|.+\.js|.+\.sql)$" \
	-exclude-dir=".git" \
	-exclude="assets.go" \
	-build="./build.sh app-dev" \
	-command="./app-dev"
