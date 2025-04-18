#!/usr/bin/env bash

openssl req -x509 -newkey rsa:4096 -sha256 -days 365 \
  -nodes -keyout key.pem -out cert.pem \
  -subj "/C=GB/ST=London/L=London/O=Dev/CN=localhost"
