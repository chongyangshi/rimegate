#!/bin/sh

set -e

# Use this script to sign client certificates with a self-signed CA certificate generated with p12-ca.sh

if [ -z "$1" ]; then
    echo "Usage: sh p12-client.sh <client-name>"
    exit 1
fi

if [ ! -f "rimegate-self-signed-root-key.pem" ]; then
    echo "CA key (rimegate-self-signed-root-key.pem) not present in working directory, cannot sign client credentials. Do you need to run p12-ca.sh?"
    exit 1
fi

if [ ! -f "rimegate-self-signed-root.pem" ]; then
    echo "CA cert (rimegate-self-signed-root.pem) not present in working directory, cannot sign client credentials. Do you need to run p12-ca.sh?"
    exit 1
fi

openssl genrsa -out client-$1-key.pem 4096
openssl req -new -sha256 -key client-$1-key.pem -subj "/C=GB/ST=London/O=Rimegate/CN=$1" -out client-$1-csr.pem
openssl x509 -req -in client-$1-csr.pem -CA rimegate-self-signed-root.pem -CAkey rimegate-self-signed-root-key.pem -CAcreateserial -out client-$1.pem -days 1095 -sha256

echo "Packaging p12, you will need to enter a password to encrypt it"
openssl pkcs12 -export -out client-$1.p12 -inkey client-$1-key.pem -in client-$1.pem -certfile rimegate-self-signed-root.pem

echo "Client key":
echo "-----------------------------------------"
cat client-$1-key.pem
echo "-----------------------------------------"
openssl rsa -in client-$1-key.pem -noout -text
echo "-----------------------------------------"
echo "Client certificate":
echo "-----------------------------------------"
cat client-$1.pem
echo "-----------------------------------------"
openssl x509 -in client-$1.pem -noout -text
echo "-----------------------------------------"
echo "Encrypted client p12 saved at client-$1.p12."