#!/bin/sh

CERT_DIR=/etc/nginx/ssl
CERT_KEY=$CERT_DIR/server.key
CERT_CRT=$CERT_DIR/server.crt

mkdir -p $CERT_DIR
if [ ! -f "$CERT_KEY" ] || [ ! -f "$CERT_CRT" ]; then
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$CERT_KEY" \
    -out "$CERT_CRT" \
    -subj "/C=US/ST=State/L=City/O=Organization/OU=Department/CN=$DOMAIN"
else
    echo "Using existing self-signed certificate."
fi
