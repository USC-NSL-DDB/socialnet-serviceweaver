#!/bin/bash

echo "[CHECKING] Go Installation"
if go version &> /dev/null; then
  echo "[PASSED] Go already installed."
else
  mkdir -p /tmp/go_install
  mkdir -p $HOME/.local/bin
  pushd /tmp/go_install
  wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
  rm -rf /usr/local/go
  rm -rf $HOME/.local/go
  tar -C $HOME/.local -xzf go1.22.2.linux-amd64.tar.gz
fi

source_cmd="export PATH=$PATH:$HOME/.local/go/bin"
echo $source_cmd >> ~/.bashrc