echo "====================================="
echo "Adding genesis oracles"
echo "====================================="
# Add genesis oracles
yes $PASSWORD | $DEMON add-genesis-oracle "$HORIZON_DID" "dap:mint"
yes $PASSWORD | $DEMON add-genesis-oracle "$BLACKHOLE_DID" "dap:mint/burn/transfer"
yes $PASSWORD | $DEMON add-genesis-oracle "$STARDUST_DID" "dollar:transfer,dap:transfer"
yes $PASSWORD | $DEMON add-genesis-oracle "$LIGHT_DID" "dollar:transfer/burn/mint"