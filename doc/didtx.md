DID base DXP transaction self sign and board cast design
=========================================================

## Step to generation ED25519 Base Key


Please visit this documentation from: https://documenter.getpostman.com/view/2597586/T1Ds9Fs5

GET `http://localhost:1320/`
```json
{
    "endpoints": [
        "/hdwallet/create/mnemonic",
        "/",
        "/hdwallet/create/{at_index}/",
        "/hdwallet/recovery"
    ]
}
```

GET `http://localhost:1320/hdwallet/create/mnemonic`
```json
{"words":"annual job denial sleep misery guess apple april message jacket require afford swamp ticket erode stumble involve skate minute satoshi trick kit virtual boat"}
```

POST `http://localhost:1320/hdwallet/recovery`

REQUEST BODY:
```json
{
	"keywords": "annual job denial sleep misery guess apple april message jacket require afford swamp ticket erode stumble involve skate minute satoshi trick kit virtual boat",
	"names":["ann","peter","lason","ming","ming","mymy@gmail.com","lalalmoon"],
	"from_index":0
}
```


POST ` "/hdwallet/create/{at_index}/"`

REQUEST BODY:
```json
{
	"keywords": "annual job denial sleep misery guess apple april message jacket require afford swamp ticket erode stumble involve skate minute satoshi trick kit virtual boat"
}
```

RESPONSE SAMPLE:

