#!/usr/bin/env bash

: ${1?'missing key directory'}
: ${2?'Common name subject. eg: webhook.namespace.svc.local'}

key_dir="$1"
common_name="$2"

chmod 0700 "$key_dir"
cd "$key_dir"

# Generate the CA cert and private key
openssl req -nodes -new -x509 -days 3000 -keyout ca.key -out ca.crt -subj "/CN=Admission Controller Webhook CA"
# Generate the private key for the webhook server
openssl genrsa -out webhook-server-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key webhook-server-tls.key -days 3000 -subj "/CN=$common_name" \
| openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook-server-tls.crt -days 3000
