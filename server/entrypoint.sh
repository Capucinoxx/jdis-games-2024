#!/bin/sh

if [ -z "$DOMAIN" ]; then
  echo "ERROR: The environment variable DOMAIN is not set."
  exit 1
fi

CERT_DIR="/app/certs"
CERT_KEY="$CERT_DIR/server.key"
CERT_CRT="$CERT_DIR/server.crt"

mkdir -p $CERT_DIR

if [ ! -f "$CERT_KEY" ] || [ ! -f "$CERT_CRT" ]; then
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$CERT_KEY" \
    -out "$CERT_CRT" \
    -subj "/C=US/ST=State/L=City/O=Organization/OU=Department/CN=$DOMAIN"
  echo "Self-signed certificate generated."
else
  echo "Using existing self-signed certificate."
fi

exec ./server