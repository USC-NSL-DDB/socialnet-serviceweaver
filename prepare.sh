#!/bin/bash

script_dir="$(dirname "$0")"

# Set weaver path
pushd  $script_dir/src
cp shared/copy/* client
cp shared/copy/* server 
popd

pushd $script_dir/src/client
weaver generate .
popd

pushd $script_dir/src/server
weaver generate .
popd

pushd $script_dir/src
go mod tidy
popd
