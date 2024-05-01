#!/bin/bash

go install github.com/ServiceWeaver/weaver/cmd/weaver@latest
go install github.com/ServiceWeaver/weaver-kube/cmd/weaver-kube@latest

source_cmd="export PATH=$PATH:$HOME/go/bin"
echo $source_cmd >> $HOME/.bashrc
