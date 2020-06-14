#!/usr/bin/env bash

echo "Exporting app state to genesis file..."
ixod export >genesis.json

echo "Fixing genesis file..."
sed -i 's/"genutil":null/"genutil":{"gentxs":null}/g' genesis.json
# https://github.com/cosmos/cosmos-sdk/issues/5086

echo "Backing up existing genesis file..."
cp "$HOME"/.ixod/config/genesis.json "$HOME"/.ixod/config/genesis.json.backup

echo "Moving new genesis file to $HOME/.ixod/config/genesis.json..."
mv genesis.json "$HOME"/.ixod/config/genesis.json

ixod unsafe-reset-all
ixod validate-genesis

ixocli init test --chain-id=darkpool-1

ixocli config output json
ixocli config indent true
ixocli config trust-node true
ixocli config chain-id namechain
ixocli config keyring-backend test

ixocli keys add jack
ixocli keys add alice

ixod add-genesis-account $(ixocli keys show jack -a) 1000dp,100000000stake
ixod add-genesis-account $(ixocli keys show alice -a) 1000dp,100000000stake

ixod gentx --name jack --keyring-backend test

echo "Collecting genesis txs..."
ixod collect-gentxs

echo "Validating genesis file..."
ixod validate-genesis

ixod start --pruning "syncable" &
ixocli rest-server --chain-id darkpool-1 --trust-node && fg
