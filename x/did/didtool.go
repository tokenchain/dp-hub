package did

import (
	"bytes"
	"crypto/ed25519"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
	naclBox "golang.org/x/crypto/nacl/box"
	"io"
)

type DxpDid struct {
	Did                 string       `json:"did" yaml:"did"`
	VerifyKey           string       `json:"verifyKey" yaml:"verifyKey"`
	EncryptionPublicKey string       `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
	Secret              SovrinSecret `json:"secret" yaml:"secret"`
}

func DidToAddr(did types.Did) sdk.AccAddress {
	return types.StringToAddr(did)
}

func UnmarshalDxpDid(jsonSovrinDid string) (DxpDid, error) {
	return fromJsonString(jsonSovrinDid)
}

func fromJsonString(jsonSovrinDid string) (DxpDid, error) {
	var did DxpDid
	err := json.Unmarshal([]byte(jsonSovrinDid), &did)
	if err != nil {
		err := fmt.Errorf("Could not unmarshal did into struct. Error: %s", err.Error())
		return DxpDid{}, err
	}

	return did, nil
}

func FromMnemonic(mnemonic string) DxpDid {
	seed := sha256.New()
	seed.Write([]byte(mnemonic))

	var seed32 [32]byte
	copy(seed32[:], seed.Sum(nil)[:32])

	return FromSeed(seed32)
}
func (sd DxpDid) String() string {
	output, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}
func FromSeed(seed [32]byte) DxpDid {

	publicKeyBytes, privateKeyBytes, err := ed25519.GenerateKey(bytes.NewReader(seed[0:32]))
	if err != nil {
		panic(err)
	}
	publicKey := []byte(publicKeyBytes)
	privateKey := []byte(privateKeyBytes)

	signKey := base58.Encode(privateKey[:32])
	keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))

	sovDid := DxpDid{
		Did:                 base58.Encode(publicKey[:16]),
		VerifyKey:           base58.Encode(publicKey),
		EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),

		Secret: SovrinSecret{
			Seed:                 hex.EncodeToString(seed[0:32]),
			SignKey:              signKey,
			EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
		},
	}

	return sovDid
}

func Gen() DxpDid {
	var seed [32]byte
	if _, err := io.ReadFull(cryptoRand.Reader, seed[:]); err != nil {
		panic(err)
	}
	return FromSeed(seed)
}
