#!/usr/bin/env bash

rm -rf "$HOME"/.dxod
rm -rf "$HOME"/.dxocli

cd "$HOME"/go/src/github.com/tokenchain/ixo-blockchain/ || exit
make install
