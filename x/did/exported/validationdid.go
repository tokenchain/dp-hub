package exported

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"regexp"
	"strings"
)

var (
	ValidDid   = regexp.MustCompile(`^did:(dxp:|sov:)([a-zA-Z0-9]){21,22}([/][a-zA-Z0-9:]+|)$`)
	IsValidDid = ValidDid.MatchString
	// https://sovrin-foundation.github.io/sovrin/spec/did-method-spec-template.html
	// IsValidDid adapted from the above link but assumes no sub-namespaces
	// TODO: ValidDid needs to be updated once we no longer want to be able
	//   to consider project accounts as DIDs (especially in treasury module),
	//   possibly should just be `^did:(dxp:|sov:)([a-zA-Z0-9]){21,22}$`.
)
func RecoverDidToEd25519PrivateKey(did_doc IxoDid)  ed25519tm.PrivKeyEd25519 {
	var privKey  ed25519tm.PrivKeyEd25519
	copy(privKey[:], base58.Decode(did_doc.Secret.EncryptionPrivateKey)[:])
	return privKey
}
func RecoverDidToEd25519PubKey(did_doc IxoDid)  ed25519tm.PubKeyEd25519 {
	var privKey  ed25519tm.PubKeyEd25519
	copy(privKey[:], base58.Decode(did_doc.EncryptionPublicKey)[:])
	return privKey
}
func RecoverDidToCosmosPrivateKey(did_doc IxoDid) secp256k1.PrivKeySecp256k1 {
	var privKey secp256k1.PrivKeySecp256k1
	copy(privKey[:], base58.Decode(did_doc.Secret.EncryptionPrivateKey)[:])
	return privKey
}
func RecoverDidToPrivateKeyClassic(did_doc IxoDid) ed25519tm.PrivKeyEd25519 {
	var privKey ed25519tm.PrivKeyEd25519
	copy(privKey[:], base58.Decode(did_doc.Secret.SignKey))
	copy(privKey[32:], base58.Decode(did_doc.VerifyKey))
	return privKey
}
func RecoverDidEd25519ToPrivateKeyC(doc IxoDid) ed25519tm.PrivKeyEd25519 {
	var recover_priv_key_ed [64]byte
	p1, _ := hex.DecodeString(strings.ToLower(doc.Secret.EncryptionPrivateKey))
	p2, _ := hex.DecodeString(strings.ToLower(doc.Secret.SignKey))
	copy(recover_priv_key_ed[:], p1[37:])
	copy(recover_priv_key_ed[32:], p2)
	return recover_priv_key_ed
}
func RecoverDidEd25519ToPrivateKey(doc IxoDid) ed25519tm.PrivKeyEd25519 {
	var recover_priv_key_ed [64]byte
	copy(recover_priv_key_ed[:], base58.Decode(doc.Secret.EncryptionPrivateKey))
	return recover_priv_key_ed
}
func RecoverDidEd25519PublicKey(doc IxoDid) [32]byte {
	var recover_pub [32]byte
	name := substring(doc.Did, 8, len(doc.Did))
	p1 := base58.Decode(name)
	p2 := base58.Decode(doc.EncryptionPublicKey)
	copy(recover_pub[:], p1)
	copy(recover_pub[16:], p2)
	return recover_pub
}

func RecoverDidSecpK1ToPrivateKey(doc IxoDid) [32]byte {
	var recover_privKey secp256k1.PrivKeySecp256k1
	p1, _ := hex.DecodeString(strings.ToLower(doc.Secret.EncryptionPrivateKey))
	p2, _ := hex.DecodeString(strings.ToLower(doc.Secret.SignKey))
	copy(recover_privKey[:], p1)
	copy(recover_privKey[24:], p2)
	return secp256k1.PrivKeySecp256k1(recover_privKey)
}

func PrivateKeyToSecp256k1(privKey tmcrypto.PrivKey) secp256k1.PrivKeySecp256k1 {
	var privKey_orginal secp256k1.PrivKeySecp256k1
	copy(privKey_orginal[:], privKey.Bytes()[5:])
	return privKey_orginal
}
func SecpPrivKey(bz []byte) secp256k1.PrivKeySecp256k1 {
	var bzArr [32]byte
	copy(bzArr[:], bz)
	return secp256k1.PrivKeySecp256k1(bzArr)
}

func PrivateKeyToEd25519(privKey tmcrypto.PrivKey) ed25519tm.PrivKeyEd25519 {
	var privKey_orginal ed25519tm.PrivKeyEd25519
	copy(privKey_orginal[:], privKey.Bytes()[:])
	return privKey_orginal
}

func getArrayFromKey(key string) []byte {
	fmt.Println(len(base58.Decode(key)))
	return base58.Decode(key)
}
