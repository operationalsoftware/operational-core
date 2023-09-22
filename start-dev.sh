#!/usr/bin/env bash
watchmedo auto-restart --directory=./ --patterns="*.py;*.scss;*.jinja2;*.json" --recursive -- hypercorn main:app 
