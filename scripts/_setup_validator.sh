

echo "====================================="
echo "gentx for validator for horizon"
echo "====================================="

yes $PASSWORD | $DEMON gentx \
  --amount 10000$TOKEN_SYM \
  --pubkey $($DEMON tendermint show-validator) \
  --name horizon \
  --commission-rate "0.10" \
  --commission-max-rate "0.10" \
  --commission-max-change-rate "0.10"

echo "====================================="
echo "add a validator"
echo "====================================="

yes $PASSWORD | $DCLI tx staking create-validator \
  --amount=10000000$TOKEN_SYM \
  --pubkey=$($DEMON tendermint show-validator) \
  --moniker=$MONIKER \
  --chain-id=$CHAIN_ID \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-prices="0.025$TOKEN_SYM" \
  --from=horizon
  #\  --generate-only

