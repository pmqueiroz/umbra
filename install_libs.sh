#!/bin/bash

if [ -z "${UMBRA_PATH}" ]; then
  UMBRA_PATH="$HOME/.umbra"
  echo "export UMBRA_PATH=\"$UMBRA_PATH\"" >> ~/.bashrc
fi

mkdir -p "$UMBRA_PATH/lib"

cp -r lib/* "$UMBRA_PATH/lib"
