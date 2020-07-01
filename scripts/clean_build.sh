#!/usr/bin/env bash

rm -rf $ENVDFOLDER
rm -rf $ENVCLIFOLDER

cd "$HOME"/go/src/github.com/tokenchain/ixo-blockchain/ || exit
make install