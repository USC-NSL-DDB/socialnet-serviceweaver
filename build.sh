#!/bin/bash

script_dir="$(dirname "$0")"

./$script_dir/prepare.sh

pushd $script_dir/src/server
go build -o server.out
popd

pushd $script_dir/src/client
go build -o client.out
popd