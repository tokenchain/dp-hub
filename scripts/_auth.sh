#!/usr/bin/env bash

PASSWORD="xxxxxxxx"
CHAIN_ID="xxxxx"
MONIKER=xxxxx
NODE="127.0.0.1:26657"
LOCAL=8.210.117.181
STAKING_TOKEN="xxx"
TOKEN_SYM="xxx"
GPG_KEY_ID="XD2CF0F5"
DEMON=appd
DCLI=appcli
ENVCLIFOLDER="$HOME/.dmcli"
ENVDFOLDER="$HOME/.dmd"


user_oper() {
	echo "====================================="
	echo "user account operation"
	echo "====================================="
	yes $PASSWORD | $DCLI keys delete $1 --force
	yes $PASSWORD | $DCLI keys add $1
}