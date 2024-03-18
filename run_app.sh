#!/bin/bash
script_dir="$(dirname "$0")"

pushd $script_dir/src
go mod tidy

go run .
popd
