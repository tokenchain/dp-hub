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
	tmcrypto "github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	edgen "github.com/tokenchain/ixo-blockchain/x/did/ed25519"
	naclBox "golang.org/x/crypto/nacl/box"
)

const (
	mnemonicEntropySize = 256
)

func NewDidGeneratorBuilder() KeyGenerator {
	return KeyGenerator{
		name:  "cosmos",
		mem:   "",
		debug: false,
	}
}
func (s KeyGenerator) GetMnemonicString() string {
	return s.mem
}
func (s KeyGenerator) GetSeedString() string {
	return hex.EncodeToString(s.seed[0:32])
}
func (s KeyGenerator) Debug() KeyGenerator {
	s.debug = true
	return s
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
func (s KeyGenerator) generateSeedBIP44(account, index uint32, bip39Passphrase string) KeyGenerator {
	hdPath := keys.CreateHDPath(account, index).String()
	seed_key, _ := keys.SecpDeriveKey(s.mem, bip39Passphrase, hdPath)
	var seed32 [32]byte
	copy(seed32[:], seed_key)
	s.seed = seed32
	return s
}
func (s KeyGenerator) generateHDSeed(hdPath string, bip39Passphrase string) KeyGenerator {
	seed_key, _ := keys.SecpDeriveKey(s.mem, bip39Passphrase, hdPath)
	var seed32 [32]byte
	copy(seed32[:], seed_key)
	s.seed = seed32
	return s
}
func (s KeyGenerator) generateMnemonic() KeyGenerator {
	entropy, _ := bip39.NewEntropy(mnemonicEntropySize)
	mnemonicWords, _ := bip39.NewMnemonic(entropy)
	s.mem = mnemonicWords
	return s
}
func (s KeyGenerator) makePubKey(bt *[32]byte) (pubKey tmcrypto.PubKey) {
	var pubKeyRaw ed25519tm.PubKeyEd25519
	copy(pubKeyRaw[:], bt[:])
	return pubKeyRaw
}
func (s KeyGenerator) generateFinal() IxoDid {
	publicKeyBytes, privateKeyBytes, err := edgen.GenerateKey(bytes.NewReader(s.seed[:32]))
	if err != nil {
		panic(err)
	}
	//head part
	signKey := base58.Encode(privateKeyBytes[:32])
	//keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))
	keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKeyBytes[:]))
	pub_key := base58.Encode(publicKeyBytes[:])
	sovDid := IxoDid{
		Did:                 dxpDidAddress(base58.Encode(publicKeyBytes[:16])),
		VerifyKey:           pub_key,
		EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),

		Secret: Secret{
			Seed:                 hex.EncodeToString(s.seed[:32]),
			SignKey:              signKey,
			EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
		},

		Dpinfo: DpInfo{
			DpAddress: VerifyKeyToAddrEd25519(pub_key).String(),
			PubKey:    sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, s.makePubKey(publicKeyBytes)),
			Name:      s.name,
			Algo:      keys.Ed25519,
		},
	}
	if s.debug {
		fmt.Println("========pair1 key  =========")
		fmt.Println("pub", len(publicKeyBytes), publicKeyBytes)
		fmt.Println("pri", len(privateKeyBytes), privateKeyBytes)
		fmt.Println("========pair2 key  =========")
		fmt.Println("pub", len(keyPairPublicKey), keyPairPublicKey)
		fmt.Println("pri", len(keyPairPrivateKey), keyPairPrivateKey)
		fmt.Println("========Derived private key  =========")
		fmt.Println(len(s.seed), s.seed)
		fmt.Println("================================")
	}
	return sovDid
}

//to get some parameters from the builder
func (s KeyGenerator) Pre() KeyGenerator {
	return s.generateMnemonic().generateSeed()
}
func (s KeyGenerator) Finalize() IxoDid {
	return s.generateFinal()
}

//one step finish
func (s KeyGenerator) Build() IxoDid {
	fmt.Println(s.mem)
	if s.mem == "" {
		return s.generateMnemonic().generateSeed().generateFinal()
	} else {
		return s.generateSeed().generateFinal()
	}
}
func (s KeyGenerator) RecoverBySeed(seed32 [32]byte) IxoDid {
	return s.WithSeed(seed32).generateFinal()
}
func (s KeyGenerator) BuildDocBIP44(account, index uint32, bip39Passphrase string) IxoDid {
	return s.generateSeedBIP44(account, index, bip39Passphrase).generateFinal()
}
func (s KeyGenerator) BuildDocHD(path, bip39Passphrase string) IxoDid {
	return s.generateHDSeed(path, bip39Passphrase).generateFinal()
}
func (s KeyGenerator) Recover(mnemonic string) IxoDid {
	return s.WithMem(mnemonic).generateSeed().generateFinal()
}
func (s KeyGenerator) RecoverBIP44(mnemonic, bip39Passphrase string, account, index uint32) IxoDid {
	return s.WithMem(mnemonic).generateSeedBIP44(account, index, bip39Passphrase).generateFinal()
}
func (s KeyGenerator) RecoverHD(mnemonic, bip39Passphrase, path string) IxoDid {
	return s.WithMem(mnemonic).generateHDSeed(path, bip39Passphrase).generateFinal()
}
