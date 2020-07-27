#!/usr/bin/env bash
. ./_auth.sh

$DEMON init $MONIKER --chain-id $CHAIN_ID

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  # Mac OSX
  pass init $GPG_KEY_ID
elif [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OSX
  $DCLI config keyring-backend file
elif [[ "$OSTYPE" == "cygwin" ]]; then
  # POSIX compatibility layer and Linux environment emulation for Windows
  echo "..."
elif [[ "$OSTYPE" == "msys" ]]; then
  # Windows
  echo "..."
elif [[ "$OSTYPE" == "freebsd"* ]]; then
  # ...
  echo "..."
fi

user_oper horizon
user_oper stardust
user_oper light
user_oper stellar
user_oper darkness
user_oper darkpool
user_oper meteor
user_oper southpole
user_oper darkpole
user_oper northpole
user_oper miner
user_oper mining
user_oper darkhole
user_oper singularity
user_oper miningpool
user_oper developer
user_oper worker
user_oper vollar
user_oper dollar
user_oper btc
user_oper eth
user_oper ccp
user_oper fee1
user_oper fee2
user_oper fee3
user_oper fee4

. ./_setup_cert.sh

echo "====================================="
echo "Adding dxp did, if do not have sponge please install. MacOS please run [brew install moreutil]. "
echo "====================================="
# Add ixo did
# if do not have sponge please install 
# brew install moreutil

cat $JSONFile|jq '.app_state.project.params.dp_did = "'$BOND_DID'"' -c $JSONFile | sponge $JSONFile

echo "====================================="
echo "config chain-id"
echo "====================================="
$DCLI config chain-id $CHAIN_ID
$DCLI config output json
$DCLI config indent true
$DCLI config trust-node true
$DCLI config node locahost:26657
echo "====================================="
echo "config daemon toml file wait 2 seconds"
echo "====================================="
echo "validator"
$DEMON tendermint show-validator
echo "validator address"
$DEMON tendermint show-address
$DEMON tendermint version


sleep 2

CONFIG_APPTOML=$ENVDFOLDER/config/app.toml
LINE_KEYRING="keyring-backend=\"file\""

if grep -q $LINE_KEYRING $CONFIG_APPTOML;
 then
   echo "the line is all set now.. "
 else
    echo "$LINE_KEYRING\n$(cat $CONFIG_APPTOML)" > $CONFIG_APPTOML
fi


sleep 2


echo "====================================="
echo "Adding genesis accounts"
echo "====================================="
# Note: important to add 'horizon' as a genesis-account since this is the chain's validator
yes $PASSWORD | $DEMON add-genesis-account $HORIZON_ADDR 10000000000$TOKEN_SYM,10000000dollar,100000000stake
yes $PASSWORD | $DEMON add-genesis-account $LIGHT_ADDR 100$TOKEN_SYM,10000dollar,100000000stake
yes $PASSWORD | $DEMON add-genesis-account $CCP_ADDR 100$TOKEN_SYM,10000dollar,100000000stake
yes $PASSWORD | $DEMON add-genesis-account $STARDUST_ADDR 10000000$TOKEN_SYM,10000dollar,100000000stake


echo "====================================="
echo "Get registered accounts"
echo "====================================="
# Add DID-based genesis accounts
#. ./_setup_oracle.sh
. ./_setup_validator.sh

echo "====================================="
echo "collect-gentxs"
echo "====================================="
$DEMON collect-gentxs
$DEMON validate-genesis
#$DEMON start --pruning "syncable" &
#$DCLI rest-server --chain-id darkpool-1x --trust-node && fg
#Great, now that we’ve initialized the chains, we can start both nodes in the background:
#gaiad start --home=$HOME/.gaiad1  &> gaia1.log & NODE1_PID=$!
#gaia start --home=$HOME/.gaiad2  &> gaia2.log & NODE2_PID=$!
#Note that we save the PID so we can later kill the processes. You can peak at your logs with tail gaia1.log, or follow them for a bit with tail -f gaia1.log.


sleep 2


echo "====================================="
echo "Fix denom staking token"
echo "====================================="
# Set staking token
cat $JSONFile|jq '.app_state.staking.params.bond_denom = "'$STAKING_TOKEN'"' -c $JSONFile | sponge $JSONFile
cat $JSONFile|jq '.app_state.crisis.constant_fee.denom = "'$STAKING_TOKEN'"' -c $JSONFile | sponge $JSONFile
cat $JSONFile|jq '.app_state.mint.params.mint_denom = "'$STAKING_TOKEN'"' -c $JSONFile | sponge $JSONFile
#cat $JSONFile|jq '.app_state.gov.deposit_params.min_deposit[0].denom = "'$STAKING_TOKEN'"' -c $JSONFile | sponge $JSONFile


sleep 2

echo "====================================="
echo "all done and its now ready to start"
echo "====================================="
#if [ $USER  == 'hesk' ] 
#then
#fi

#if [ $USER  == 'root' ] 
#then
#fi

  cp $JSONFile ./genesis.json
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  cp $JSONFile ./genesis.json
elif [[ "$OSTYPE" == "darwin"* ]]; then
# Mac OSX
  cp $JSONFile ./genesis.json
 # cp $JSONFile $HOME/Documents/ixo/b-explorer-settings/v1/genesis.json
 # cd $HOME/Documents/ixo/b-explorer-settings
 # sh pushcommit.sh
elif [[ "$OSTYPE" == "cygwin" ]]; then
  # POSIX compatibility layer and Linux environment emulation for Windows
  echo "..."
elif [[ "$OSTYPE" == "msys" ]]; then
  # Windows
  echo "..."
elif [[ "$OSTYPE" == "freebsd"* ]]; then
  # ...
  echo "..."
fi


