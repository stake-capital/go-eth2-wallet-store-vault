#!/bin/bash

PWD=$(pwd)
mkdir -p ./local
PLUGIN_PATH=${PWD}/local

read -r -d '' JSON_CONFIG << EOM
{
    "storage": {
        "inmem": {}
    },
    "disable_mlock": true,
    "ui": true,
    "api_addr": "http://127.0.0.1:8200"
}
EOM

echo ${JSON_CONFIG} > ./local/vault.json
vault server -dev -config=$PLUGIN_PATH/vault.json  -dev-root-token-id="golang-test"