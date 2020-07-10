package ed25519

import (
	"bytes"
	"crypto"
	"crypto/sha512"
	"crypto/subtle"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"io"
	"strconv"
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey []byte

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey []byte

// Public returns the PublicKey corresponding to priv.
func (priv PrivateKey) Public() crypto.PublicKey {
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, priv[32:])
	return PublicKey(publicKey)
}

// Seed returns the private key seed corresponding to priv. It is provided for
// interoperability with RFC 8032. RFC 8032's private keys correspond to seeds
// in this package.
func (priv PrivateKey) Seed() []byte {
	seed := make([]byte, SeedSize)
	copy(seed, priv[:32])
	return seed
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey tmCrypto.PubKey, err error) {
	err = ModuleCdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}
func (privKey PrivateKey) Bytes() []byte {
	return ModuleCdc.MustMarshalBinaryBare(privKey)
	//return []byte(privKey[:])
}
func (privKey PrivateKey) Sign(msg []byte) ([]byte, error) {
	signatureBytes := Sign(privKey, msg)
	return signatureBytes, nil
	//return nil, nil
}
func (privKey PrivateKey) Equals(other tmCrypto.PrivKey) bool {
	if otherEd, ok := other.(PrivateKey); ok {
		return subtle.ConstantTimeCompare(privKey[:], otherEd[:]) == 1
	}
	return false
}

func (privKey PrivateKey) PrivKey() tmCrypto.PrivKey {
	key, err := PrivKeyFromBytes(privKey.Bytes())
	if err != nil {
		panic("cannot decode binary bare to crypto")
	}
	return key
}

func (privKey PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey{%s}", base58.Encode(privKey))
}

// PubKey gets the corresponding public key from the private key.
func (privKey PrivateKey) PubKey() tmCrypto.PubKey {
	privKeyBytes := make([]byte, PrivateKeySize)

	//var privKeyBytes [PrivateKeySize]byte
	copy(privKeyBytes, privKey)
	//privKeyBytes := [PrivateKeySize]byte(privKey)
	initialized := false
	// If the latter 32 bytes of the privkey are all zero, compute the pubkey
	// otherwise privkey is initialized and we can use the cached value inside
	// of the private key.
	println(fmt.Sprintf("===== the cap is now: %d", cap(privKeyBytes)))
	for i, v := range privKeyBytes[32:] {
		fmt.Printf("Value: %d  Value-Addr: %X  ElemAddr: %X\n", v, &v, &privKeyBytes[i])
		if v != 0 {
			initialized = true
			break
		}
	}

	// println("=============")
	// println(privKey.String())

	if !initialized {
		panic("Expected PrivKeyEd25519dp to include concatenated pubkey bytes")
	}

	var pubkeyBytes [PubKeyEd25519Size]byte
	copy(pubkeyBytes[:], privKey[32:])
	return PubKeyEd25519dp(pubkeyBytes)
}

// Sign signs the message with privateKey and returns a signature. It will
// panic if len(privateKey) is not PrivateKeySize.
func Sign(privateKey PrivateKey, message []byte) []byte {
	if l := len(privateKey); l != PrivateKeySize {
		panic("ed25519: bad private key length: " + strconv.Itoa(l))
	}

	h := sha512.New()
	h.Write(privateKey[:32])

	var digest1, messageDigest, hramDigest [64]byte
	var expandedSecretKey [32]byte
	h.Sum(digest1[:0])
	copy(expandedSecretKey[:], digest1[:])
	expandedSecretKey[0] &= 248
	expandedSecretKey[31] &= 63
	expandedSecretKey[31] |= 64

	h.Reset()
	h.Write(digest1[32:])
	h.Write(message)
	h.Sum(messageDigest[:0])

	var messageDigestReduced [32]byte
	ScReduce(&messageDigestReduced, &messageDigest)
	var R ExtendedGroupElement
	GeScalarMultBase(&R, &messageDigestReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	h.Reset()
	h.Write(encodedR[:])
	h.Write(privateKey[32:])
	h.Write(message)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	ScReduce(&hramDigestReduced, &hramDigest)

	var s [32]byte
	ScMulAdd(&s, &hramDigestReduced, &expandedSecretKey, &messageDigestReduced)

	signature := make([]byte, SignatureSize)
	copy(signature[:], encodedR[:])
	copy(signature[32:], s[:])

	return signature
}

// Sign signs the message with privateKey and returns a signature.
func SignClassic(privateKey *[PrivateKeySize]byte, message []byte) *[SignatureSize]byte {
	h := sha512.New()
	h.Write(privateKey[:32])

	var digest1, messageDigest, hramDigest [64]byte
	var expandedSecretKey [32]byte
	h.Sum(digest1[:0])
	copy(expandedSecretKey[:], digest1[:])
	expandedSecretKey[0] &= 248
	expandedSecretKey[31] &= 63
	expandedSecretKey[31] |= 64

	h.Reset()
	h.Write(digest1[32:])
	h.Write(message)
	h.Sum(messageDigest[:0])

	var messageDigestReduced [32]byte
	ScReduce(&messageDigestReduced, &messageDigest)
	var R ExtendedGroupElement
	GeScalarMultBase(&R, &messageDigestReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	h.Reset()
	h.Write(encodedR[:])
	h.Write(privateKey[32:])
	h.Write(message)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	ScReduce(&hramDigestReduced, &hramDigest)

	var s [32]byte
	ScMulAdd(&s, &hramDigestReduced, &expandedSecretKey, &messageDigestReduced)

	signature := new([64]byte)
	copy(signature[:], encodedR[:])
	copy(signature[32:], s[:])
	return signature
}

func Verify(publicKey PublicKey, message, sig []byte) bool {
	if l := len(publicKey); l != PublicKeySize {
		panic("ed25519: bad public key length: " + strconv.Itoa(l))
	}

	if len(sig) != SignatureSize || sig[63]&224 != 0 {
		return false
	}

	var A ExtendedGroupElement
	var publicKeyBytes [32]byte
	copy(publicKeyBytes[:], publicKey)
	if !A.FromBytes(&publicKeyBytes) {
		return false
	}
	FeNeg(&A.X, &A.X)
	FeNeg(&A.T, &A.T)

	h := sha512.New()
	h.Write(sig[:32])
	h.Write(publicKey[:])
	h.Write(message)
	var digest [64]byte
	h.Sum(digest[:0])

	var hReduced [32]byte
	ScReduce(&hReduced, &digest)

	var R ProjectiveGroupElement
	var s [32]byte
	copy(s[:], sig[32:])

	// https://tools.ietf.org/html/rfc8032#section-5.1.7 requires that s be in
	// the range [0, order) in order to prevent signature malleability.
	if !ScMinimal(&s) {
		return false
	}

	GeDoubleScalarMultVartime(&R, &hReduced, &A, &s)

	var checkR [32]byte
	R.ToBytes(&checkR)
	return bytes.Equal(sig[:32], checkR[:])
}

// Verify returns true iff sig is a valid signature of message by publicKey.
func VerifyClassic(publicKey *[PublicKeySize]byte, message []byte, sig *[SignatureSize]byte) bool {
	if sig[63]&224 != 0 {
		return false
	}

	var A ExtendedGroupElement
	if !A.FromBytes(publicKey) {
		return false
	}
	FeNeg(&A.X, &A.X)
	FeNeg(&A.T, &A.T)

	h := sha512.New()
	h.Write(sig[:32])
	h.Write(publicKey[:])
	h.Write(message)
	var digest [64]byte
	h.Sum(digest[:0])

	var hReduced [32]byte
	ScReduce(&hReduced, &digest)

	var R ProjectiveGroupElement
	var b [32]byte
	copy(b[:], sig[32:])
	GeDoubleScalarMultVartime(&R, &hReduced, &A, &b)

	var checkR [32]byte
	R.ToBytes(&checkR)
	return subtle.ConstantTimeCompare(sig[:32], checkR[:]) == 1
}

/*
// GenerateKey generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rand io.Reader) (PublicKey, PrivateKey, error) {
	if rand == nil {
		rand = cryptorand.Reader
	}

	seed := make([]byte, SeedSize)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, nil, err
	}

	privateKey := NewKeyFromSeed(seed)
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, privateKey[32:])

	return publicKey, privateKey, nil
}*/

// GenerateKey generates a public/private key pair using randomness from rand.
func GenerateKey(rand io.Reader) (publicKey *[PublicKeySize]byte, privateKey *[PrivateKeySize]byte, err error) {
	privateKey = new([64]byte)
	_, err = io.ReadFull(rand, privateKey[:32])
	if err != nil {
		return nil, nil, err
	}

	publicKey = MakePublicKey(privateKey)
	return
}
// NewKeyFromSeed calculates a private key from a seed. It will panic if
// len(seed) is not SeedSize. This function is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this
// package.
func NewKeyFromSeed(seed []byte) PrivateKey {
	// Outline the function body so that the returned key can be stack-allocated.
	privateKey := make([]byte, PrivateKeySize)
	newKeyFromSeed(privateKey, seed)
	return privateKey
}
func NewKeyFromSeedToEdPrivateKey(seed []byte) (PrivateKey, tmCrypto.PrivKey) {
	private_key := make([]byte, PrivateKeySize)
	newKeyFromSeed(private_key, seed)
	var privKeyEd PrivateKey

	signKey := base58.Encode(seed)
	println(signKey)

	copy(privKeyEd, private_key[:PrivateKeySize])
	println(privKeyEd.String())
	copy(privKeyEd, private_key[0:PrivateKeySize])
	println(privKeyEd.String())
	copy(privKeyEd, private_key)
	println(privKeyEd.String())
	copy(privKeyEd[:], private_key[:])
	println(privKeyEd.String())
	return privKeyEd, tmCrypto.PrivKey(privKeyEd)
}

// MakePublicKey makes a publicKey from the first half of privateKey.
func MakePublicKey(privateKey *[PrivateKeySize]byte) (publicKey *[PublicKeySize]byte) {
	publicKey = new([32]byte)

	h := sha512.New()
	h.Write(privateKey[:32])
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest)
	GeScalarMultBase(&A, &hBytes)
	A.ToBytes(publicKey)

	copy(privateKey[32:], publicKey[:])
	return
}
func newKeyFromSeed(privateKey, seed []byte) {
	if l := len(seed); l != SeedSize {
		panic("ed25519: bad seed length: " + strconv.Itoa(l))
	}

	digest := sha512.Sum512(seed)
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest[:])
	GeScalarMultBase(&A, &hBytes)
	var publicKeyBytes [32]byte
	A.ToBytes(&publicKeyBytes)

	copy(privateKey, seed)
	copy(privateKey[32:], publicKeyBytes[:])
}
