#!/usr/bin/env bash


export GOBIN=$HOME/go/bin
PASSWORD="123123123"
CHAIN_ID="testnetdp"
MONIKER=dptest
#NODE="127.0.0.1:26657"
NODE="0.0.0.0:26657"
LOCAL=8.210.117.181
STAKING_TOKEN="dap"
TOKEN_SYM="dap"
GPG_KEY_ID="SD2CF0F5"
DAEMON_NAME="dpd"
#DAEMON_NAME="appd"
CLI_NAME="dpcli"
#CLI_NAME="appcli"
DEMON=$GOBIN/$DAEMON_NAME
DCLI=$GOBIN/$CLI_NAME
ENVCLIFOLDER="$HOME/.dpcli"
#ENVCLIFOLDER="$HOME/.dmcli"
ENVDFOLDER="$HOME/.dpd"
#ENVDFOLDER="$HOME/.dmd"
JSONFile=$ENVDFOLDER/config/genesis.json

dclis(){
  echo "====================================="
  yes $PASSWORD | $DCLI $1
  echo "====================================<."
}

dcli_show_key(){
  echo "====================================="
  yes $PASSWORD | $DCLI keys show $1
  echo "====================================<."
}

user_oper() {
	echo "====================================="
	echo "user account operation"
	echo "====================================="
	yes $PASSWORD | $DCLI keys delete $1 --force
	yes $PASSWORD | $DCLI keys add $1
}

tx() {
  cmd=$1
  shift
  $DCLI tx bonds "$cmd" --broadcast-mode block "$@"
}

