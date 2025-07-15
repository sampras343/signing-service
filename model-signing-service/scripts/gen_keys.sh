#!/bin/bash

mkdir -p keys

# Generate EC private key
openssl ecparam -genkey -name prime256v1 -noout -out keys/private.pem

# Extract public key
openssl ec -in keys/private.pem -pubout -out keys/public.pem

echo "✅ Keys generated in ./keys"
