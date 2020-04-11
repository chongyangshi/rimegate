#!/bin/sh

set -e

# Use this script to generate a self-signed CA certificate for signing client certificates
# for authentication with a reverse proxy in front of Rimegate, allowing Rimegate to be used
# in a mode relying on static Grafana API tokens.

openssl genrsa -out rimegate-self-signed-root-key.pem 4096
openssl req -x509 -new -nodes -subj "/C=GB/ST=London/O=Rimegate/CN=Rimegate" -key rimegate-self-signed-root-key.pem -sha256 -days 3560 -out rimegate-self-signed-root.pem

echo "Self-signed Root CA key":
echo "-----------------------------------------"
cat rimegate-self-signed-root-key.pem
echo "-----------------------------------------"
openssl rsa -in rimegate-self-signed-root-key.pem -noout -text
echo "-----------------------------------------"
echo "Self-signed Root CA":
echo "-----------------------------------------"
cat rimegate-self-signed-root.pem
echo "-----------------------------------------"
openssl x509 -in rimegate-self-signed-root.pem -noout -text
echo "-----------------------------------------"