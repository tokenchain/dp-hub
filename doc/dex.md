DEX Documentation
======================

###Enlist



###Delist

Submit a dex delist proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal delist-proposal <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:
```json
{
 "title": "delist xxx/%s",
 "description": "delist asset from dex",
 "base_asset": "xxx",
 "quote_asset": "%s",
 "deposit": [
   {
     "denom": "%s",
     "amount": "100"
   }
 ]
}
```