command line
=========================

#### Query base
Starts from `cli q did`

| **Command**             | **Arguments** | **Description**                |
|-------------------------|---------------|--------------------------------|
| get-address-from-did    | [did]         |  Query for an account address by DID           |
| get-did-doc             | [did]         | Query DidDoc for a DID        |
| get-all-dids            | N/A           | Query all DIDs              |
| get-all-did-docs        | N/A         | Query all DID documents                |



#### Tx base
Starts from `cli tx did`


| **Command**             | **Arguments** | **Description**                |
|-------------------------|---------------|--------------------------------|
| add-did-doc    | [sovrin-did]         |  Add a new SovrinDid          |
| add-kyc-credential             | [did] [signer-did-doc]        | Add a new KYC Credential for a Did by the signer        |
