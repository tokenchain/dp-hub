#!/usr/bin/env bash

echo "Exporting app state to genesis file..."
ixod export >genesis.json

echo "Collecting genesis txs..."
ixod collect-gentxs

echo "Validating genesis file..."
ixod validate-genesis

ixod start --pruning "syncable" &
ixocli rest-server --chain-id darkpool-1 --trust-node && fg
