# Add DID-based genesis accounts

PASSWORD="12348888EX"
CHAIN_ID="darkpool-1x"

MIGUEL_ADDR="$(ixocli keys show horizon -a)"       # address from did:dxp:DJXazTE9Se8Kzkknn8xyAe
FRANCESCO_ADDR="$(ixocli keys show stardust -a)" # address from did:dxp:UKzkhVSHc3qEFva5EY2XHt
SHAUN_ADDR="$(ixocli keys show stellar -a)"         # address from did:dxp:U4tSpzzv91HHqWW1YmFkHJ

MIGUEL_DID="did:dxp:DJXazTE9Se8Kzkknn8xyAe"
FRANCESCO_DID="did:dxp:UKzkhVSHc3qEFva5EY2XHt"
SHAUN_DID="did:dxp:U4tSpzzv91HHqWW1YmFkHJ"
SEXY_DID="did:dxp:VQJotb6QJq2VRhxBtKMiQ"

ixocli tx did addKycCredential "did:dxp:DJXazTE9Se8Kzkknn8xyAe" "did:dxp:UKzkhVSHc3qEFva5EY2XHt"
yes "12348888EX" | ixod gentx --name horizon
echo "show the validator key"
ixod tendermint show-validator
VALIDATOR_KEY="$(ixod tendermint show-validator)"

echo "====================================="
echo "add a validator"
echo "====================================="
ixocli tx staking create-validator \
  --amount=29800000000000dpl \
  --pubkey=$(ixod tendermint show-validator) \
  --moniker="local" \
  --chain-id=$CHAIN_ID \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-prices="0.025dpl" \
  --from=horizon


echo "====================================="
echo "validator add for horizon success"
echo "====================================="
