#!/usr/bin/env bash
case `uname -s` in
    Linux*)     sslConfig=/etc/ssl/openssl.cnf;;
    Darwin*)    sslConfig=/System/Library/OpenSSL/openssl.cnf;;
esac
openssl req
    -newkey rsa:2048
    -x509
    -nodes
    -keyout server.key
    -new
    -out server.pem
    -subj /CN=localhost
    -reqexts SAN
    -extensions SAN
    -config <(cat $sslConfig
        <(printf '[SAN]nsubjectAltName=DNS:localhost'))
    -sha256
    -days 3650


    # openssl req -newkey rsa:2048 -x509 -nodes -keyout server.key -new -out server.pem -subj /CN=localhost -reqexts SAN -extensions SAN -config <(cat $sslConfig <(printf '[SAN]nsubjectAltName=DNS:localhost')) -sha256 -days 3650
    # openssl genrsa -out private.key 2048
    # openssl req -new -x509 -key private.key -out cacert.pem -days 1095