package ed25519

import (
	"bytes"
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"io"
)

const (
	PrivKeyAminoName = "darkpool/PrivKeyEd25519dp"
	PubKeyAminoName  = "darkpool/PubKeyEd25519dp"
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

var _ crypto.PrivKey = PrivateKey{}

// PrivKeyEd25519dp implements crypto.PrivKey.

type (
	PubKeyEd25519dp [PubKeyEd25519Size]byte
)

// Address is the SHA256-20 of the raw pubkey bytes.
func (pubKey PubKeyEd25519dp) Address() crypto.Address {
	return crypto.Address(tmhash.SumTruncated(pubKey[:]))
}

// Bytes marshals the PubKey using amino encoding.
func (pubKey PubKeyEd25519dp) Bytes() []byte {
	return ModuleCdc.MustMarshalBinaryBare(pubKey)
}

func (pubKey PubKeyEd25519dp) VerifyBytes(msg []byte, sig []byte) bool {
	// make sure we use the same algorithm to sign
	if len(sig) != SignatureSize {
		return false
	}
	//return false
	return Verify(pubKey[:], msg, sig)
}

func (pubKey PubKeyEd25519dp) String() string {
	return fmt.Sprintf("PubKeyEd25519dp{%X}", pubKey[:])
}

// nolint: golint
func (pubKey PubKeyEd25519dp) Equals(other crypto.PubKey) bool {
	if otherEd, ok := other.(PubKeyEd25519dp); ok {
		return bytes.Equal(pubKey[:], otherEd[:])
	}
	return false
}

// PrivKeyFromBytes unmarshals private key bytes and returns a PrivKey
func PrivKeyFromBytes(privKeyBytes []byte) (privKey crypto.PrivKey, err error) {
	err = ModuleCdc.UnmarshalBinaryBare(privKeyBytes, &privKey)
	return
}

func GenPrivKey() PrivateKey {
	return genPrivKey(crypto.CReader())
}

// genPrivKey generates a new ed25519 private key using the provided reader.
func genPrivKey(rand io.Reader) PrivateKey {
	seed := make([]byte, 32)
	_, err := io.ReadFull(rand, seed)
	if err != nil {
		panic(err)
	}

	privKey := NewKeyFromSeed(seed)
	var privKeyEd PrivateKey
	copy(privKeyEd[:], privKey)
	return privKeyEd
}

// genPrivKey generates a new ed25519 private key using the provided reader.
func PrivKeyToEdPrivateKeyccccc(privKey PrivateKey) PrivateKey {
	var privKeyEd PrivateKey
	copy(privKeyEd[:], privKey)
	return privKeyEd
}

// GenPrivKeyFromSecret hashes the secret with SHA2, and uses
// that 32 byte output to create the private key.
// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyFromSecret(secret []byte) PrivateKey {
	seed := crypto.Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	privKey := NewKeyFromSeed(seed)
	var privKeyEd PrivateKey
	copy(privKeyEd[:], privKey)
	return privKeyEd
}
