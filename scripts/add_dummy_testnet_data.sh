#!/usr/bin/env bash

wait() {
  echo "Waiting for chain to start..."
  while :; do
    RET=$(ixocli status 2>&1)
    if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
      sleep 1
    else
      echo "A few more seconds..."
      sleep 6
      break
    fi
  done
}

RET=$(ixocli status 2>&1)
if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
  wait
fi

PASSWORD="12348888EX"

FEE1=$(yes $PASSWORD | ixocli keys show fee -a)
FEE2=$(yes $PASSWORD | ixocli keys show fee2 -a)
FEE3=$(yes $PASSWORD | ixocli keys show fee3 -a)
FEE4=$(yes $PASSWORD | ixocli keys show fee4 -a)

BOND1_DID="did:dxp:U7GK8p8rVhJMKhBVRCJJ8c"
BOND2_DID="did:dxp:JHcN95bkS4aAWk3TKXapA2"
BOND3_DID="did:dxp:48PVm1uyF6QVDSPdGRWw4T"
BOND4_DID="did:dxp:RYLHkfNpbA8Losy68jt4yF"
BOND1_DID_FULL="{\"did\":\"did:dxp:U7GK8p8rVhJMKhBVRCJJ8c\",\"verifyKey\":\"FmwNAfvV2xEqHwszrVJVBR3JgQ8AFCQEVzo1p6x4L8VW\",\"encryptionPublicKey\":\"domKpTpjrHQtKUnaFLjCuDLe2oHeS4b1sKt7yU9cq7m\",\"secret\":{\"seed\":\"933e454dbcfc1437f3afc10a0cd512cf0339787b6595819849f53707c268b053\",\"signKey\":\"Aun1EpjR1HQu1idBsPQ4u4C4dMwtbYPe1SdSC5bUerFC\",\"encryptionPrivateKey\":\"Aun1EpjR1HQu1idBsPQ4u4C4dMwtbYPe1SdSC5bUerFC\"}}"
BOND2_DID_FULL="{\"did\":\"did:dxp:JHcN95bkS4aAWk3TKXapA2\",\"verifyKey\":\"ARTUGePyi4rm3ogq3kjp8dEAVq1RR7Z3HWwLze6Ey4qg\",\"encryptionPublicKey\":\"FjLXxW1N68XgBKekfB2isCLHPwbqhLiQLQFhLiiivXqP\",\"secret\":{\"seed\":\"9fc38edde7b9d7097b0aeefcb22a8fcccb6c0748fc3034eec5abdac0740339f7\",\"signKey\":\"Bken38cPuz2Mosb3poKQHd81Q2BiTWFznG1Wa7ZZkrni\",\"encryptionPrivateKey\":\"Bken38cPuz2Mosb3poKQHd81Q2BiTWFznG1Wa7ZZkrni\"}}"
BOND3_DID_FULL="{\"did\":\"did:dxp:48PVm1uyF6QVDSPdGRWw4T\",\"verifyKey\":\"2hs2cb232Ev97aSQLvrfK4q8ZceBR8cf33UTstWpKU9M\",\"encryptionPublicKey\":\"9k2THnNbTziXGRjn77tvWujffgigRPqPyKZUSdwjmfh2\",\"secret\":{\"seed\":\"82949a422215a5999846beaadf398659157c345564787993f92e91d192f2a9c5\",\"signKey\":\"9njRge76sTYdfcpFfBG5p2NwbDXownFzUyTeN3iDQdjz\",\"encryptionPrivateKey\":\"9njRge76sTYdfcpFfBG5p2NwbDXownFzUyTeN3iDQdjz\"}}"
BOND4_DID_FULL="{\"did\":\"did:dxp:RYLHkfNpbA8Losy68jt4yF\",\"verifyKey\":\"ENmMCsfNmjYoTRhNgnwXbQAw6p8JKH9DCJfGTPXNfsxW\",\"encryptionPublicKey\":\"5unQBt6JPW1pq9AqoRNhFJmibv8JqeoyyNvN3gF24EaU\",\"secret\":{\"seed\":\"d2c05b107acc2dfe3e9d67e98c993a9c03d227ed8f4505c43997cf4e7819bee2\",\"signKey\":\"FBgjjwoVPfd8ZUVs4nXVKtf4iV6xwnKhaMKBBEAHvGtH\",\"encryptionPrivateKey\":\"FBgjjwoVPfd8ZUVs4nXVKtf4iV6xwnKhaMKBBEAHvGtH\"}}"

PROJECT1_DID="did:dxp:U7GK8p8rVhJMKhBVRCJJ8c"
PROJECT2_DID="did:dxp:JHcN95bkS4aAWk3TKXapA2"
PROJECT1_DID_FULL="{\"did\":\"did:dxp:U7GK8p8rVhJMKhBVRCJJ8c\",\"verifyKey\":\"FmwNAfvV2xEqHwszrVJVBR3JgQ8AFCQEVzo1p6x4L8VW\",\"encryptionPublicKey\":\"domKpTpjrHQtKUnaFLjCuDLe2oHeS4b1sKt7yU9cq7m\",\"secret\":{\"seed\":\"933e454dbcfc1437f3afc10a0cd512cf0339787b6595819849f53707c268b053\",\"signKey\":\"Aun1EpjR1HQu1idBsPQ4u4C4dMwtbYPe1SdSC5bUerFC\",\"encryptionPrivateKey\":\"Aun1EpjR1HQu1idBsPQ4u4C4dMwtbYPe1SdSC5bUerFC\"}}"
PROJECT2_DID_FULL="{\"did\":\"did:dxp:JHcN95bkS4aAWk3TKXapA2\",\"verifyKey\":\"ARTUGePyi4rm3ogq3kjp8dEAVq1RR7Z3HWwLze6Ey4qg\",\"encryptionPublicKey\":\"FjLXxW1N68XgBKekfB2isCLHPwbqhLiQLQFhLiiivXqP\",\"secret\":{\"seed\":\"9fc38edde7b9d7097b0aeefcb22a8fcccb6c0748fc3034eec5abdac0740339f7\",\"signKey\":\"Bken38cPuz2Mosb3poKQHd81Q2BiTWFznG1Wa7ZZkrni\",\"encryptionPrivateKey\":\"Bken38cPuz2Mosb3poKQHd81Q2BiTWFznG1Wa7ZZkrni\"}}"
PROJECT1_INFO="{\"nodeDid\":\"nodeDid\",\"requiredClaims\":\"500\",\"evaluatorPayPerClaim\":\"50\",\"serviceEndpoint\":\"serviceEndpoint\",\"createdOn\":\"2020-01-01T01:01:01.000Z\",\"createdBy\":\"Miguel\",\"status\":\"\"}"
PROJECT2_INFO="{\"nodeDid\":\"nodeDid\",\"requiredClaims\":\"100\",\"evaluatorPayPerClaim\":\"10\",\"serviceEndpoint\":\"serviceEndpoint\",\"createdOn\":\"2020-02-02T02:02:02.000Z\",\"createdBy\":\"Francesco\",\"status\":\"\"}"

BONDDOC1_DID="did:dxp:48PVm1uyF6QVDSPdGRWw4T"
BONDDOC2_DID="did:dxp:RYLHkfNpbA8Losy68jt4yF"
BONDDOC1_DID_FULL="{\"did\":\"did:dxp:48PVm1uyF6QVDSPdGRWw4T\",\"verifyKey\":\"2hs2cb232Ev97aSQLvrfK4q8ZceBR8cf33UTstWpKU9M\",\"encryptionPublicKey\":\"9k2THnNbTziXGRjn77tvWujffgigRPqPyKZUSdwjmfh2\",\"secret\":{\"seed\":\"82949a422215a5999846beaadf398659157c345564787993f92e91d192f2a9c5\",\"signKey\":\"9njRge76sTYdfcpFfBG5p2NwbDXownFzUyTeN3iDQdjz\",\"encryptionPrivateKey\":\"9njRge76sTYdfcpFfBG5p2NwbDXownFzUyTeN3iDQdjz\"}}"
BONDDOC2_DID_FULL="{\"did\":\"did:dxp:RYLHkfNpbA8Losy68jt4yF\",\"verifyKey\":\"ENmMCsfNmjYoTRhNgnwXbQAw6p8JKH9DCJfGTPXNfsxW\",\"encryptionPublicKey\":\"5unQBt6JPW1pq9AqoRNhFJmibv8JqeoyyNvN3gF24EaU\",\"secret\":{\"seed\":\"d2c05b107acc2dfe3e9d67e98c993a9c03d227ed8f4505c43997cf4e7819bee2\",\"signKey\":\"FBgjjwoVPfd8ZUVs4nXVKtf4iV6xwnKhaMKBBEAHvGtH\",\"encryptionPrivateKey\":\"FBgjjwoVPfd8ZUVs4nXVKtf4iV6xwnKhaMKBBEAHvGtH\"}}"
BONDDOC1_INFO="{\"createdOn\":\"createdOn\",\"createdBy\":\"createdBy\"}"
BONDDOC2_INFO="{\"createdOn\":\"createdOn\",\"createdBy\":\"createdBy\"}"

MIGUEL_DID="did:dxp:4XJLBfGtWSGKSz4BeRxdun"
FRANCESCO_DID="did:dxp:UKzkhVSHc3qEFva5EY2XHt"
SHAUN_DID="did:dxp:U4tSpzzv91HHqWW1YmFkHJ"
SEXY_DID="did:dxp:U4tSexy1HHqFW1YmFkHJ"

MIGUEL_DID_FULL="{\"did\":\"did:dxp:4XJLBfGtWSGKSz4BeRxdun\",\"verifyKey\":\"2vMHhssdhrBCRFiq9vj7TxGYDybW4yYdrYh9JG56RaAt\",\"encryptionPublicKey\":\"6GBp8qYgjE3ducksUa9Ar26ganhDFcmYfbZE9ezFx5xS\",\"secret\":{\"seed\":\"38734eeb53b5d69177da1fa9a093f10d218b3e0f81087226be6ce0cdce478180\",\"signKey\":\"4oMozrMR6BXRN93MDk6UYoqBVBLiPn9RnZhR3wQd6tBh\",\"encryptionPrivateKey\":\"4oMozrMR6BXRN93MDk6UYoqBVBLiPn9RnZhR3wQd6tBh\"}}"
FRANCESCO_DID_FULL="{\"did\":\"did:dxp:UKzkhVSHc3qEFva5EY2XHt\",\"verifyKey\":\"Ftsqjc2pEvGLqBtgvVx69VXLe1dj2mFzoi4kqQNGo3Ej\",\"encryptionPublicKey\":\"8YScf3mY4eeHoxDT9MRxiuGX5Fw7edWFnwHpgWYSn1si\",\"secret\":{\"seed\":\"94f3c48a9b19b4881e582ba80f5767cd3f3c5d7b7103cb9a50fa018f108d89de\",\"signKey\":\"B2Svs8GoQnUJHg8W2Ch7J53Goq36AaF6C6W4PD2MCPrM\",\"encryptionPrivateKey\":\"B2Svs8GoQnUJHg8W2Ch7J53Goq36AaF6C6W4PD2MCPrM\"}}"
SHAUN_DID_FULL="{\"did\":\"did:dxp:U4tSpzzv91HHqWW1YmFkHJ\",\"verifyKey\":\"FkeDue5it82taeheMprdaPrctfK3DeVV9NnEPYDgwwRG\",\"encryptionPublicKey\":\"DtdGbZB2nSQvwhs6QoN5Cd8JTxWgfVRAGVKfxj8LA15i\",\"secret\":{\"seed\":\"6ef0002659d260a0bbad194d1aa28650ccea6c6862f994dfdbd48648e1a05c5e\",\"signKey\":\"8U474VrG2QiUFKfeNnS84CAsqHdmVRjEx4vQje122ycR\",\"encryptionPrivateKey\":\"8U474VrG2QiUFKfeNnS84CAsqHdmVRjEx4vQje122ycR\"}}"

# ----------------------------------------------------------------------------------------- dids
# Ledger DIDs
echo "Ledgering DID 1/3..."
ixocli tx did addDidDoc "$MIGUEL_DID_FULL"
echo "Ledgering DID 2/3..."
ixocli tx did addDidDoc "$FRANCESCO_DID_FULL"
echo "Ledgering DID 3/3..."
ixocli tx did addDidDoc "$SHAUN_DID_FULL"

echo "Sleeping for a bit..."
sleep 6 # to make sure DIDs were ledgered before proceeding

# Adding KYC credentials
echo "Adding KYC credential 1/1..."
ixocli tx did addKycCredential "$MIGUEL_DID" "$FRANCESCO_DID_FULL"

# ----------------------------------------------------------------------------------------- mints/burns
# Mint and burn dxp tokens
echo "Minting 1000ixo tokens to Miguel using Miguel oracle..."
ixocli tx treasury oracle-mint "$MIGUEL_DID" 1000dpl "$MIGUEL_DID_FULL" "dummy proof"
echo "Burning 1000ixo tokens from Francesco using Francesco oracle..."
ixocli tx treasury oracle-burn "$FRANCESCO_DID" 1000dpl "$FRANCESCO_DID_FULL" "dummy proof"

# ----------------------------------------------------------------------------------------- bonds
# Power function with m:12,n:2,c:100, rez reserve, non-zero fees, and batch_blocks=1
echo "Creating bond 1/4..."
ixocli tx bonds create-bond \
  --token=token1 \
  --name="Test Token 1" \
  --description="Power function with non-zero fees and batch_blocks=1" \
  --function-type=power_function \
  --function-parameters="m:12,n:2,c:100" \
  --reserve-tokens=res \
  --tx-fee-percentage=0.5 \
  --exit-fee-percentage=0.1 \
  --fee-address="$FEE1" \
  --max-supply=1000000token1 \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells=true \
  --batch-blocks=1 \
  --bond-did="$BOND1_DID_FULL" \
  --creator-did="$MIGUEL_DID"

# Power function with m:10,n:3,c:0, res reserve, zero fees, and batch_blocks=3
echo "Creating bond 2/4..."
ixocli tx bonds create-bond \
  --token=token2 \
  --name="Test Token 2" \
  --description="Power function with zero fees and batch_blocks=4" \
  --function-type=power_function \
  --function-parameters="m:10,n:3,c:1" \
  --reserve-tokens=res \
  --tx-fee-percentage=0 \
  --exit-fee-percentage=0 \
  --fee-address="$FEE2" \
  --max-supply=1000000token2 \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells=true \
  --batch-blocks=3 \
  --bond-did="$BOND2_DID_FULL" \
  --creator-did="$MIGUEL_DID"

# Swapper function between res and rez with zero fees, and batch_blocks=2
echo "Creating bond 3/4..."
ixocli tx bonds create-bond \
  --token=token3 \
  --name="Test Token 3" \
  --description="Swapper function between res and rez" \
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens="res,rez" \
  --tx-fee-percentage=0 \
  --exit-fee-percentage=0 \
  --fee-address="$FEE3" \
  --max-supply=1000000token3 \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells=true \
  --batch-blocks=2 \
  --bond-did="$BOND3_DID_FULL" \
  --creator-did="$MIGUEL_DID"

# Swapper function between token1 and token2 with non-zero fees, and batch_blocks=1
echo "Creating bond 4/4..."
ixocli tx bonds create-bond \
  --token=token4 \
  --name="Test Token 4" \
  --description="Swapper function between res and rez" \
  --function-type=swapper_function \
  --function-parameters="" \
  --reserve-tokens="token1,token2" \
  --tx-fee-percentage=2.5 \
  --exit-fee-percentage=5 \
  --fee-address="$FEE4" \
  --max-supply=1000000token4 \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells=true \
  --batch-blocks=1 \
  --bond-did="$BOND4_DID_FULL" \
  --creator-did="$MIGUEL_DID"

echo "Sleeping for a bit..."
sleep 6 # to make sure bonds were ledgered before proceeding

# Buy 5token1, 5token2 from Miguel
echo "Buying 5token1 from Miguel..."
ixocli tx bonds buy 5token1 "100000res" "$BOND1_DID" "$MIGUEL_DID_FULL"
echo "Buying 5token2 from Miguel..."
ixocli tx bonds buy 5token2 "100000res" "$BOND2_DID" "$MIGUEL_DID_FULL"

echo "Sleeping for a bit..."
sleep 6 # to make sure buys were processed before proceeding

# Buy token2 and token3 from Francesco and Shaun
echo "Buying 5token2 from Francesco..."
ixocli tx bonds buy 5token2 "100000res" "$BOND2_DID" "$FRANCESCO_DID_FULL"
echo "Buying 5token3 from Shaun..."
ixocli tx bonds buy 5token3 "100res,100rez" "$BOND3_DID" "$SHAUN_DID_FULL"

echo "Sleeping for a bit..."
sleep 6 # to make sure buys were processed before proceeding

# Buy 5token4 from Miguel (using token1 and token2)
echo "Buying 5token4 from Miguel..."
ixocli tx bonds buy 5token4 "2token1,2token2" "$BOND4_DID" "$MIGUEL_DID_FULL"

# ----------------------------------------------------------------------------------------- projects
# Create projects (this creates a project doc for the respective project)
SENDER_DID="$SHAUN_DID"
echo "Creating project 1/2..."
ixocli tx project createProject "$SENDER_DID" "$PROJECT1_INFO" "$PROJECT1_DID_FULL"
echo "Creating project 2/2..."
ixocli tx project createProject "$SENDER_DID" "$PROJECT2_INFO" "$PROJECT2_DID_FULL"

echo "Sleeping for a bit..."
sleep 6 # to make sure projects were ledgered before proceeding

# Update project status (this updates the status in the project doc for the respective project)
SENDER_DID="$SHAUN_DID"
echo "Updating project 1 to CREATED..."
ixocli tx project updateProjectStatus "$SENDER_DID" CREATED "$PROJECT1_DID_FULL"
echo "Updating project 2 to CREATED..."
ixocli tx project updateProjectStatus "$SENDER_DID" CREATED "$PROJECT2_DID_FULL" --broadcast-mode block
echo "Updating project 2 to PENDING..."
ixocli tx project updateProjectStatus "$SENDER_DID" PENDING "$PROJECT2_DID_FULL" --broadcast-mode block

# Fund project (using treasury 'send' and 'oracle-transfer')
echo "Funding project 2 (using treasury 'send' from Miguel)..."
ixocli tx treasury send "$PROJECT2_DID/$PROJECT2_DID" 5000000000ixo "$MIGUEL_DID_FULL" --broadcast-mode block
echo "Funding project 2 (using treasury 'oracle-transfer' from Miguel using Francesco oracle)..."
ixocli tx treasury oracle-transfer "$MIGUEL_DID" "$PROJECT2_DID/$PROJECT2_DID" 5000000000ixo "$FRANCESCO_DID_FULL" "dummy proof" --broadcast-mode block
# The address behind "$PROJECT2_DID/$PROJECT2_DID" can also be obtained from (ixocli q project getProjectAccounts $PROJECT2_DID)
# Note that we're actually sending just 100ixo, since ixoDecimals is 1e8 and we're sending 100e8ixo
echo "Updating project 2 to FUNDED..."
SENDER_DID="$SHAUN_DID"
ixocli tx project updateProjectStatus "$SENDER_DID" FUNDED "$PROJECT2_DID_FULL" --broadcast-mode block

# Adding a claim and evaluation
echo "Creating a claim in project 2..."
SENDER_DID="$SHAUN_DID"
ixocli tx project createClaim "tx_hash" "$SENDER_DID" "claim_id" "$PROJECT2_DID_FULL" --broadcast-mode block
echo "Creating an evaluation in project 2..."
SENDER_DID="$MIGUEL_DID"
STATUS="1" # createEvaluation updates status of claim from 0 to 1 implicitly (explicitly in blocksync)
ixocli tx project createEvaluation "tx_hash" "$SENDER_DID" "claim_id" $STATUS "$PROJECT2_DID_FULL" --broadcast-mode block

# Adding agents (this creates a project account for the agent in the respective project)
echo "Adding agent to project 1..."
SENDER_DID="did:dxp:48PVm1uyF6QVDSPdGRWw4T"
AGENT_DID="did:dxp:RYLHkfNpbA8Losy68jt4yF"
ROLE="SA"
ixocli tx project createAgent "tx_hash" "$SENDER_DID" "$AGENT_DID" "$ROLE" "$PROJECT1_DID_FULL"

# ----------------------------------------------------------------------------------------- bonddocs
# Creating bonddoc
SENDER_DID="$SHAUN_DID"
echo "Creating bonddoc 1/2..."
ixocli tx bonddoc createBond "$SENDER_DID" "$BONDDOC1_INFO" "$BONDDOC1_DID_FULL"
echo "Creating bonddoc 1/2..."
ixocli tx bonddoc createBond "$SENDER_DID" "$BONDDOC2_INFO" "$BONDDOC2_DID_FULL"

echo "Sleeping for a bit..."
sleep 6 # to make sure bonddocs were ledgered before proceeding

# Updating bonddoc status
SENDER_DID="$SHAUN_DID"
echo "Updating bonddoc 1 to PREISSUANCE..."
ixocli tx bonddoc updateBondStatus "$SENDER_DID" PREISSUANCE "$BONDDOC1_DID_FULL"
echo "Updating bonddoc 2 to PREISSUANCE..."
ixocli tx bonddoc updateBondStatus "$SENDER_DID" PREISSUANCE "$BONDDOC2_DID_FULL" --broadcast-mode block
echo "Updating bonddoc 2 to OPEN..."
ixocli tx bonddoc updateBondStatus "$SENDER_DID" OPEN "$BONDDOC2_DID_FULL"
