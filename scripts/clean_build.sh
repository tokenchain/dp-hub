#!/usr/bin/env bash

rm -rf "$HOME"/.dxod
rm -rf "$HOME"/.dxocli
rm -rf "$HOME"/.dpd
rm -rf "$HOME"/.dpcli

cd "$HOME"/go/src/github.com/tokenchain/dp-block/ || exit
make install
