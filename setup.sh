#!/bin/bash

echo "[CHECKING] Go Installation"
if go version &> /dev/null; then
  echo "[PASSED] Go installed."
else
  echo "[FAILED] Install Go manually first"
  exit 1
fi

echo "[INSTALLING] Service Weaver"
go install github.com/ServiceWeaver/weaver/cmd/weaver@latest

PATH=$PATH:$HOME/go/bin
EXPORT_CMD="export PATH=$PATH"

if [[ $SHELL == */bash ]]; then
  echo "$EXPORT_CMD" >> $HOME/.bashrc
  echo "source .bashrc manually"
elif [[ $SHELL == */zsh ]]; then
  echo "$EXPORT_CMD" >> $HOME/.zshrc
  echo "source .zshrc manually"
else
  echo "Current shell is not supported"
fi

echo "[DONE] Service Weaver Installation"

