# Add DID-based genesis accounts
MIGUEL_ADDR="$(ixocli keys show miguel -a)"       # address from did:dxp:DJXazTE9Se8Kzkknn8xyAe
FRANCESCO_ADDR="$(ixocli keys show francesco -a)" # address from did:dxp:UKzkhVSHc3qEFva5EY2XHt
SHAUN_ADDR="$(ixocli keys show shaun -a)"         # address from did:dxp:U4tSpzzv91HHqWW1YmFkHJ
yes "$PASSWORD" | ixod add-genesis-account miguel 100000dpl,100000000STAKE,10000usdt
yes "$PASSWORD" | ixod add-genesis-account francesco 100000000STAKE,100000dpl,10000usdt
yes "$PASSWORD" | ixod add-genesis-account shaun 100000000STAKE,100000dpl,10000usdt


yes "12348888EX"| ixod add-genesis-account francesco 100000dpl,10000usdt,100000000stake
yes "12348888EX"| ixod add-genesis-account miguel 100000dpl
yes "12348888EX"| ixod add-genesis-account shaun 100000dpl,10000usdt,100000000stake

MIGUEL_DID="did:dxp:DJXazTE9Se8Kzkknn8xyAe"
FRANCESCO_DID="did:dxp:UKzkhVSHc3qEFva5EY2XHt"
SHAUN_DID="did:dxp:U4tSpzzv91HHqWW1YmFkHJ"
SEXY_DID="did:dxp:VQJotb6QJq2VRhxBtKMiQ"

ixocli tx did addKycCredential "did:dxp:DJXazTE9Se8Kzkknn8xyAe" "did:dxp:UKzkhVSHc3qEFva5EY2XHt"
yes "12348888EX" | ixod gentx --name miguel
echo "show the validator key"
ixod tendermint show-validator
VALIDATOR_KEY="$(ixod tendermint show-validator)"

ixocli tx staking create-validator \
  --amount=1000000dpl \
  --pubkey=$VALIDATOR_KEY \
  --moniker="hongmin" \
  --chain-id=darkpool-1x \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-prices="0.025usdt" \
  --from=miguel
