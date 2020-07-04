#!/usr/bin/env bash

HORIZON_ADDR=$(yes $PASSWORD | $DCLI keys show horizon -a)
STARDUST_ADDR=$(yes $PASSWORD | $DCLI keys show stardust -a)
LIGHT_ADDR=$(yes $PASSWORD | $DCLI keys show light -a)
DARKNESS_ADDR=$(yes $PASSWORD | $DCLI keys show darkness -a)
MININGPOOL_ADDR=$(yes $PASSWORD | $DCLI keys show miningpool -a)
BLACKHOLE_ADDR=$(yes $PASSWORD | $DCLI keys show darkhole -a)
CCP_ADDR=$(yes $PASSWORD | $DCLI keys show ccp -a)
BOND_DID="did:dxp:PGtRP4o5cCYhe8B5Fu7xj6"
BOND_DID_FULL=$(jq -c . did_full_bond.json)

STARDUST_DID="did:dxp:QmuNtGJCJZuU8wm8e9PSWk"
STARDUST_DID_FULL=$(jq -c . did_full_bond.json)
HORIZON_DID="did:dxp:XX1rB7o9BeEmNuvB5PeAn2"
HORIZON_DID_FULL=$(jq -c . did_full_bond.json)
BLACKHOLE_DID="did:dxp:Xc7PR8S7A6AzreYR5th2h4"
BLACKHOLE_DID_FULL=$(jq -c . did_full_blackhole.json)
LIGHT_DID="did:dxp:PQQMd4iP8XUpctiETKRhW2"
LIGHT_DID_FULL=$(jq -c . did_full_light.json)

FEE1=$(yes $PASSWORD | $DCLI keys show fee1 -a)
FEE2=$(yes $PASSWORD | $DCLI keys show fee2 -a)
FEE3=$(yes $PASSWORD | $DCLI keys show fee3 -a)
FEE4=$(yes $PASSWORD | $DCLI keys show fee4 -a)

BONDDOC1_DID="did:dxp:VkydFvdB7YMSkwSARTPR5e"
BONDDOC1_DID_FULL=$(jq -c . didfullbonddoc.json)


#yes C2VS8888EX | $DCLI keys show horizon -a

#$DCLI query account dx01qxpej02f506dep649cvd9x24a98pasgkqr3vu9