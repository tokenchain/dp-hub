# Darkpool 客户端 The second part

RPC remote port by cosmos

https://cosmos.network/rpc/v0.38.3

#### API Doc
Please find the api document to be located at `*:1317/swagger-ui/` at the LCD.


#### Other source for development guide
The core blockchain
https://github.com/tokenchain/dp-hub

The explorer
https://github.com/tokenchain/dpexplorer

The js SDK
https://github.com/tokenchain/ledger-darkpool-js

# Darkpool HDWallet Generation Online for dedicated IP hosting

mainnet: 8.210.227.164
port: 1315

1. 生成地址的 接口

GET `/hdwallet/create/mnemonic`
POST `/hdwallet/recovery`
POST `/hdwallet/create/{at_index}/`

Please read the postman documentation: https://documenter.getpostman.com/view/2597586/T1Ds9Fs5

port: 1312
3. 签名 发送交易接口

GET `/txs/{hash}`
GET `/txs/{hash}`
POST `/txs`
POST `/txs/decode`

port: 1312
2. 查询余额 接口

GET `/auth/accounts/{address}`