#!/usr/bin/env bash
CompileDaemon \
	-pattern="(.+\.go|.+\.c|.+\.png|.+\.svg|.+\.css|.+\.js|.+\.env)$" \
	-exclude-dir=".git" \
	-exclude="assets.go" \
	-build="./build.sh app-dev" \
	-command="./app-dev"
