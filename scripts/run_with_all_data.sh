#!/usr/bin/env bash

PASSWORD="12348888EX"

ixod init local --chain-id darkpool-1x

yes $PASSWORD | ixocli keys delete miguel --force
yes $PASSWORD | ixocli keys delete francesco --force
yes $PASSWORD | ixocli keys delete shaun --force
yes $PASSWORD | ixocli keys delete sexy --force
yes $PASSWORD | ixocli keys delete fee1 --force
yes $PASSWORD | ixocli keys delete fee2 --force
yes $PASSWORD | ixocli keys delete fee3 --force
yes $PASSWORD | ixocli keys delete fee4 --force
yes $PASSWORD | ixocli keys delete fee5 --force

yes $PASSWORD | ixocli keys add miguel
yes $PASSWORD | ixocli keys add francesco
yes $PASSWORD | ixocli keys add shaun
yes $PASSWORD | ixocli keys add sexy
yes $PASSWORD | ixocli keys add fee1
yes $PASSWORD | ixocli keys add fee2
yes $PASSWORD | ixocli keys add fee3
yes $PASSWORD | ixocli keys add fee4
yes $PASSWORD | ixocli keys add fee5
echo "====================================="
echo "Adding genesis accounts"
echo "====================================="
# Note: important to add 'miguel' as a genesis-account since this is the chain's validator
yes $PASSWORD | ixod add-genesis-account miguel 1000000dpl,100000000stake,10000usdt
yes $PASSWORD | ixod add-genesis-account "$(ixocli keys show sexy -a)" 100000dpl,1000000stake,1000usdt
yes $PASSWORD | ixod add-genesis-account francesco 1000000stake,100000dpl,10000usdt
yes $PASSWORD | ixod add-genesis-account shaun 1000000stake,100000dpl,10000usdt
echo "====================================="
echo "Get registered accounts"
echo "====================================="
# Add DID-based genesis accounts
MIGUEL_ADDR="$(ixocli keys show miguel -a)"       # address from did:dxp:DJXazTE9Se8Kzkknn8xyAe
FRANCESCO_ADDR="$(ixocli keys show francesco -a)" # address from did:dxp:UKzkhVSHc3qEFva5EY2XHt
SHAUN_ADDR="$(ixocli keys show shaun -a)"         # address from did:dxp:U4tSpzzv91HHqWW1YmFkHJ
#yes $PASSWORD | ixod add-genesis-account "$MIGUEL_ADDR" 100000000stake,100000dpl,10000usdt
#yes $PASSWORD | ixod add-genesis-account "$FRANCESCO_ADDR" 100000000stake,100000dpl,10000usdt
#yes $PASSWORD | ixod add-genesis-account "$SHAUN_ADDR" 100000000stake,100000dpl,10000usdt
echo "====================================="
echo "Adding genesis oracles"
echo "====================================="
# Add genesis oracles
MIGUEL_DID="did:dxp:DJXazTE9Se8Kzkknn8xyAe"
FRANCESCO_DID="did:dxp:UKzkhVSHc3qEFva5EY2XHt"
SHAUN_DID="did:dxp:U4tSpzzv91HHqWW1YmFkHJ"
SEXY_DID="did:dxp:VQJotb6QJq2VRhxBtKMiQ"
yes $PASSWORD | ixod add-genesis-oracle "$MIGUEL_DID" "dpl:mint"
yes $PASSWORD | ixod add-genesis-oracle "$FRANCESCO_DID" "dpl:mint/burn/transfer"
yes $PASSWORD | ixod add-genesis-oracle "$SHAUN_DID" "usdt:transfer,dpl:transfer"
yes $PASSWORD | ixod add-genesis-oracle "$SEXY_DID" "usdt:transfer/burn/mint"
echo "====================================="
echo "Adding dxp did, if do not have sponge please install. MacOS please run [brew install moreutil]. "
echo "====================================="
# Add ixo did
DXP_DID=$MIGUEL_DID
# if do not have sponge please install 
# brew install moreutil
JSONFile="$HOME"/.dxod/config/genesis.json
cat $JSONFile|jq '.app_state.project.params.dp_did = "'$DXP_DID'"' -c $JSONFile | sponge $JSONFile

echo "====================================="
echo "config chain-id"
echo "====================================="
ixocli config chain-id darkpool-1x
echo "====================================="
echo "config output json"
echo "====================================="
ixocli config output json
ixocli config indent true
ixocli config trust-node true
echo "====================================="
echo "gentx for validator"
echo "====================================="
yes $PASSWORD | ixod gentx --name miguel

ixod collect-gentxs
ixod validate-genesis

ixod start --pruning "syncable" &
ixocli rest-server --chain-id darkpool-1x --trust-node && fg