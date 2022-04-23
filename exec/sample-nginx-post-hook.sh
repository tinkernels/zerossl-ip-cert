#!/usr/bin/env bash

echo "nginx post hook running"

echo "ZEROSSL_CERT_FPATH: $ZEROSSL_CERT_FPATH"
echo "ZEROSSL_KEY_FPATH: $ZEROSSL_KEY_FPATH"

nginx_bin=$(which nginx)
echo "nginx_bin: $nginx_bin"

"$nginx_bin" -t
"$nginx_bin" -s reload