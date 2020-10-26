Generation of DID key
======================

## Step to generation ED25519 Base Key

Step 1 -
Generate Mnemonic String for 24 words = `mnemonic`

Step 2 -
Generate Seed by the given `mnemonic` and cut the length to 32 byte and using the hex function encode. Hex.Encode(Mnemonic[0:32]) = `seed`

Step 3 -
Generate private key and public key from the seed using the ed25519 function. Here we got `publicKey` and `privatekey`.

Step 4 -
Generate `signKey` by the using base58 encode function with length [0:32]. 

Step 5 -
Generate Second key pairs by using the `privatekey` got from the step 3. Take function curve25519 key generator to make public key `pubkey` and `prikey`. Using base58 encode function with fell length for both `pubkey` and `prikey` we got `encryptionPublicKey` (5.1) and `encryptionPrivateKey` (5.2).

Step 6 -
Taking the value `publicKey` from the step 3 and using base58 encode function with full length we get `verifyKey`.

Step 7 -
Taking the value `publicKey` from the step 3 and using base58 encode function with [0:16] length we get `did` key = "did:dxp:" + value

Step 8 -
Taking the value `publicKey` from the step 3 and use SHA256 sum function and cutting = result1. Take the hash from result1 and cut [0:20] = result2. Using bech32 convertbits function from the result2. 

`result3 = bech32.converbits(data_byte=result2, from_bits=8, to_bits=5, pad=true)`

dpaddress = bech32.encode (prefix="dx0", data_byte=result3)


Step 9 -
Taking the value `publicKey` from the step 3. Using bech32 convertbits function from the result2. 

`result4 = bech32.converbits(data_byte=publicKey, from_bits=8, to_bits=5, pad=true)`

`dpaddress = bech32.encode (prefix="dx0pub", data_byte=result4)`

Step 10 -
Collect the customer username by their email


The result concept should able to place the results into the below demonstration.
```json
{
             "did": "Step7",
             "verifyKey": "Step6",
             "encryptionPublicKey": "Step5.1",
             "secret": {
               "seed": "Step2",
               "signKey": "Step4",
               "encryptionPrivateKey": "Step5.2"
             },
             "dp": {
               "address": "Step8",
               "pubkey": "Step9",
               "name": "Step10",
               "algo": "ed25519"
             }
}
```

## Seed Generation by the program

Using the code base we can get from this 
```go
builder := exported.NewDidGeneratorBuilder().Pre()
seed:= builder.GetMnemonicString()
mnemonic:= builder.GetSeedString()
did:= builder.Finalize()
```



## Recovery of Key by Mnemonic by the program
Using the code base we can get from this 
```go
my_recovery_mnemonic="word1 word2 word3 ..."
did := exported.NewDidGeneratorBuilder().Recover(my_recovery_mnemonic)
```


Using the code base we can get from this 
```go
my_recovery_seed="74fd93fdd7508e6b2fc9f4e1ac8cef727003f2c36e3f1acf1fcb104658da8f42"
did := exported.NewDidGeneratorBuilder().RecoverBySeed(my_recovery_seed)
```
## Test Data

In total of 24 words that consist of the recovery mnemonic
`sample_did_01_mem = "better swap climb night chronic border process gift drastic cabin jazz find pupil twin breeze lawn peanut banana tail empower civil borrow edit dentist"`

Using the above mnemonic should able to generate the below key. The key to be generated with the public key and private key.
```json
           {
             "did": "did:dxp:VrsU9cUAcYgF7f397xtjsX",
             "verifyKey": "GjKLRmDSCLALj28519q8XwKTmJTfFpobEsWCCKWHhzut",
             "encryptionPublicKey": "2Pb4bkbk1oXTpypzuZFABUtDgDP8VCZEiBVuXJsVbgYb",
             "secret": {
               "seed": "74fd93fdd7508e6b2fc9f4e1ac8cef727003f2c36e3f1acf1fcb104658da8f42",
               "signKey": "8sgZQSCiu8GHTveWb1mfusT1KbaFCksriHoXhyDzwahF",
               "encryptionPrivateKey": "8sgZQSCiu8GHTveWb1mfusT1KbaFCksriHoXhyDzwahF"
             },
             "dp": {
               "address": "dx01nyx8wn3qelmdpykjcjqnw22zdmu9pjt9us0y73",
               "pubkey": "dx0pub1zcjduepqaxmxmerk2lw76qxlsf8cc0rzst5hfgy4a3xmvnaxgmkczc9hv30sh4yj2y",
               "name": "cosmos",
               "algo": "ed25519"
             }
           }
```

# Key generator source code also see the [original code](https://github.com/tokenchain/dp-hub/blob/1.3.8x-dex/x/did/exported/generator.go)
```
package exported

import (
        "bytes"
        "crypto/sha256"
        "encoding/hex"
        "fmt"
        "github.com/btcsuite/btcutil/base58"
        "github.com/cosmos/cosmos-sdk/crypto/keys"
        sdk "github.com/cosmos/cosmos-sdk/types"
        "github.com/cosmos/go-bip39"
        ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
        edgen "github.com/tokenchain/dp-block/x/did/ed25519"
        naclBox "golang.org/x/crypto/nacl/box"
)

func NewDidGeneratorBuilder() KeyGenerator {
        return KeyGenerator{
                name: "cosmos",
                mem:  "",
        }
}
func (s KeyGenerator) GetMnemonicString() string {
        return s.mem
}
func (s KeyGenerator) GetSeedString() string {
        return hex.EncodeToString(s.seed[0:32])
}
func (s KeyGenerator) WithName(n string) KeyGenerator {
        s.name = n
        return s
}

func (s KeyGenerator) WithPubKey(n []byte) KeyGenerator {
        s.pubkey = n
        return s
}

func (s KeyGenerator) WithPrivKey(n []byte) KeyGenerator {
        s.privkey = n
        return s
}

func (s KeyGenerator) WithMem(n string) KeyGenerator {
        s.mem = n
        return s
}
func (s KeyGenerator) WithSeed(seed32 [32]byte) KeyGenerator {
        s.seed = seed32
        return s
}
func (s KeyGenerator) generateSeed() KeyGenerator {
        seed := sha256.New()
        seed.Write([]byte(s.mem))
        var seed32 [32]byte
        copy(seed32[:], seed.Sum(nil)[:32])
        s.seed = seed32
        return s
}
func (s KeyGenerator) generateMnemonic() KeyGenerator {
        entropy, _ := bip39.NewEntropy(12)
        mnemonicWords, _ := bip39.NewMnemonic(entropy)
        s.mem = mnemonicWords
        return s
}
func (s KeyGenerator) generateFinal() IxoDid {
        publicKeyBytes, privateKeyBytes, err := edgen.GenerateKey(bytes.NewReader(s.seed[0:32]))
        if err != nil {
                panic(err)
        }
        //head part
        signKey := base58.Encode(privateKeyBytes[:32])
        //keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))
        keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKeyBytes[:]))

        var pubKey ed25519tm.PubKeyEd25519
        copy(pubKey[:], publicKeyBytes[:])

        sovDid := IxoDid{
                Did:                 dxpDidAddress(base58.Encode(publicKeyBytes[:16])),
                VerifyKey:           base58.Encode(publicKeyBytes[:]),
                EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),

                Secret: Secret{
                        Seed:                 hex.EncodeToString(s.seed[0:32]),
                        SignKey:              signKey,
                        EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
                },

                Dpinfo: DpInfo{
                        DpAddress: sdk.AccAddress(pubKey.Address()).String(),
                        PubKey:    sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey),
                        Name:      s.name,
                        Algo:      keys.Ed25519,
                },
        }
        return sovDid
}

func (s KeyGenerator) Build() IxoDid {
        fmt.Println(s.mem)
        if s.mem == "" {
                return s.generateMnemonic().generateSeed().generateFinal()
        } else {
                return s.generateSeed().generateFinal()
        }
}

func (s KeyGenerator) BuildWithCustomSeed(seed32 [32]byte) IxoDid {
        return s.WithSeed(seed32).generateFinal()
}

func (s KeyGenerator) Recover(mem string) IxoDid {
        return s.WithMem(mem).generateSeed().generateFinal()
}

```


#### API Doc
Please find the api document to be located at `*:1317/swagger-ui/` at the LCD.
