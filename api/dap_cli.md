

# Darkpool 客户端 The second part

## Darkpool CLI


### 配置 dpcli

设置`dpcli`的主要命令如下：

Add a new SovrinDid from the full json document
```shell script
dpcli tx did add-did-doc [sovrin-did]
```

Add a new KYC Credential for a Did by the signer
```shell script
dpcli tx did add-kyc-credential [did] [signer-did-doc]
```

Generate did document offline
```shell script
dpcli tx did generate-offline [name]
```

Query for an account address by DID
```shell script
dpcli q did get-address-from-did [did]
```

Query DidDoc for a DID
```shell script
dpcli q did get-did-doc [did]
```

Query all DIDs
```shell script
dpcli q did get-all-dids
```

Query all DID documents
```shell script
dpcli q did get-all-did-docs
```



### Darkpool bond handling

Creating bond with the bond name
```shell script
dpcli tx bonds create-bond
  --token=chkd \
  --name="DPHKD" \
  --description="The New HKD" \
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens=stake,$TOKEN_SYM \
  --tx-fee-percentage=0.015 \
  --exit-fee-percentage=0.02 \
  --fee-address="$FEE1" \
  --max-supply=10000000000dollar \
  --order-quantity-limits="5000dollar,5000dap" \
  --sanity-rate="0.5" \
  --sanity-margin-percentage="20" \
  --allow-sells=true \
  --batch-blocks=1 \
  --bond-did="$DID_NOVA" \
  --creator-did="$DIDSOVRIN_SINGULARITY" \
  --broadcast-mode block

```

Editing bonds
```shell script
dpcli tx bonds edit-bond
  --token=chkd \
  --name="DPHKD" \
  --description="The New HKD" \
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens=stake,$TOKEN_SYM \
  --tx-fee-percentage=0.015 \
  --exit-fee-percentage=0.02 \
  --fee-address="$FEE1" \
  --order-quantity-limits="5000dollar,5000dap" \
  --sanity-rate="0.5" \
  --sanity-margin-percentage="20" \
  --allow-sells=true \
  --batch-blocks=1 \
  --broadcast-mode block

```

Buy Bond
```shell script
dpcli q bonds buy [bond-token-with-amount] [max-prices] [bond-did] [buyer-did]
```

Sell from a bond
```shell script
dpcli q bonds sell [bond-token-with-amount] [bond-did] [seller-did]
```

Swap bonds
```shell script
dpcli q bonds swap [from-amount] [from-token] [to-token] [bond-did] [swapper-did]
```


List all bonds
```shell script
dpcli q bonds bonds-list
```

Query info of a bond
```shell script
dpcli q bonds bond [bond-did
```

Query info of a bond's current batch
```shell script
dpcli q bonds batch [bond-did]
```

Query info of a bond's last batch
```shell script
dpcli q bonds last-batch [bond-did]
```

Query current price(s) of the bond
```shell script
dpcli q bonds last-batch [bond-did]
```

Query current balance(s) of the reserve pool
```shell script
dpcli q bonds current-reserve [bond-did]
```

Query price(s) of the bond at a specific supply
```shell script
dpcli q bonds price [bond-token-with-amount] [bond-did]
```

Query price(s) of buying an amount of tokens of the bond
```shell script
dpcli q bonds buy-price [bond-token-with-amount] [bond-did]
```

Query return(s) on selling an amount of tokens of the bond
```shell script
dpcli q bonds sell-return [bond-token-with-amount] [bond-did]
```

Query return(s) on swapping an amount of tokens to another token
```shell script
dpcli q bonds swap-return [bond-did] [from-token-with-amount] [to-token]
```

