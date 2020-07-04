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
#tendermint node --rpc.laddr=tcp://0.0.0.0:26657
#When I try to do `tx staking edit-validator` because the demon cannot start with out a validator in the genesis file for the first start privat full node. It gots me this error from creation of validator ABCIQuery: Post failed: Post "http://localhost:26657": dial tcp [::1]:26657: connect: connection refused. What is the correct way to set the configuration?