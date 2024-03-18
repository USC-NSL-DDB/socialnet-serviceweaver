#!/bin/bash

script_dir="$(dirname "$0")"

pushd $script_dir/src/server
go build -o server.out
popd