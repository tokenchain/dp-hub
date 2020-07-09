package ed25519

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"io"
)

const (
	// PublicKeySize is the size, in bytes, of public keys as used in this package.
	PublicKeySize = 32
	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
	PrivateKeySize = 64
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = 64
	// SeedSize is the size, in bytes, of private key seeds. These are the private key representations used by RFC 8032.
	SeedSize          = 32
	PubKeyEd25519Size = 32
)

// PrivKeyEd25519 implements crypto.PrivKey.
type (
	PrivKeyEd25519 [PrivateKeySize]byte
	PubKeyEd25519  [PubKeyEd25519Size]byte
)

var _ crypto.PubKey = PubKeyEd25519{}
var _ crypto.PrivKey = PrivKeyEd25519{}

var cdc = amino.NewCodec()

// Address is the SHA256-20 of the raw pubkey bytes.
func (pubKey PubKeyEd25519) Address() crypto.Address {
	return crypto.Address(tmhash.SumTruncated(pubKey[:]))
}

// Bytes marshals the PubKey using amino encoding.
func (pubKey PubKeyEd25519) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(pubKey)
}

func (pubKey PubKeyEd25519) VerifyBytes(msg []byte, sig []byte) bool {
	// make sure we use the same algorithm to sign
	if len(sig) != SignatureSize {
		return false
	}
	return false

	//return ed25519.Verify(pubKey[:], msg, sig)
}

func (pubKey PubKeyEd25519) String() string {
	return fmt.Sprintf("PubKeyEd25519{%X}", pubKey[:])
}

// nolint: golint
func (pubKey PubKeyEd25519) Equals(other crypto.PubKey) bool {
	if otherEd, ok := other.(PubKeyEd25519); ok {
		return bytes.Equal(pubKey[:], otherEd[:])
	}
	return false
}

// PrivKeyFromBytes unmarshals private key bytes and returns a PrivKey
func PrivKeyFromBytes(privKeyBytes []byte) (privKey crypto.PrivKey, err error) {
	err = cdc.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}
func (privKey PrivKeyEd25519) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(privKey)
}
func (privKey PrivKeyEd25519) Sign(msg []byte) ([]byte, error) {
	//signatureBytes := ed25519.Sign(privKey[:], msg)
	//return signatureBytes, nil
	return nil, nil
}
func (privKey PrivKeyEd25519) Equals(other crypto.PrivKey) bool {
	if otherEd, ok := other.(PrivKeyEd25519); ok {
		return subtle.ConstantTimeCompare(privKey[:], otherEd[:]) == 1
	}
	return false
}
func (privKey PrivKeyEd25519) PrivKey() crypto.PrivKey {
	key, err := PrivKeyFromBytes(privKey.Bytes())
	if err != nil {
		panic("cannot decode binary bare to crypto")
	}
	return key
}
func (privKey PrivKeyEd25519) String() string {
	return fmt.Sprintf("PrivKeyEd25519{%X}", privKey[:])
}
func GenPrivKey() PrivKeyEd25519 {
	return genPrivKey(crypto.CReader())
}

// genPrivKey generates a new ed25519 private key using the provided reader.
func genPrivKey(rand io.Reader) PrivKeyEd25519 {
	seed := make([]byte, 32)
	_, err := io.ReadFull(rand, seed)
	if err != nil {
		panic(err)
	}

	privKey := NewKeyFromSeed(seed)
	var privKeyEd PrivKeyEd25519
	copy(privKeyEd[:], privKey)
	return privKeyEd
}

// PubKey gets the corresponding public key from the private key.
func (privKey PrivKeyEd25519) PubKey() crypto.PubKey {
	privKeyBytes := [PrivateKeySize]byte(privKey)
	initialized := false
	// If the latter 32 bytes of the privkey are all zero, compute the pubkey
	// otherwise privkey is initialized and we can use the cached value inside
	// of the private key.
	for _, v := range privKeyBytes[32:] {
		if v != 0 {
			initialized = true
			break
		}
	}

	if !initialized {
		panic("Expected PrivKeyEd25519 to include concatenated pubkey bytes")
	}

	var pubkeyBytes [PubKeyEd25519Size]byte
	copy(pubkeyBytes[:], privKeyBytes[32:])
	return PubKeyEd25519(pubkeyBytes)
}


// genPrivKey generates a new ed25519 private key using the provided reader.
func PrivKeyToEdPrivateKey(privKey PrivateKey) PrivKeyEd25519 {
	var privKeyEd PrivKeyEd25519
	copy(privKeyEd[:], privKey)
	return privKeyEd
}

// GenPrivKeyFromSecret hashes the secret with SHA2, and uses
// that 32 byte output to create the private key.
// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyFromSecret(secret []byte) PrivKeyEd25519 {
	seed := crypto.Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	privKey := NewKeyFromSeed(seed)
	var privKeyEd PrivKeyEd25519
	copy(privKeyEd[:], privKey)
	return privKeyEd
}
