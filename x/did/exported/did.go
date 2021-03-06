package exported

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
)

type (
	Did    = string
	DidDoc interface {
		SetDid(did Did) error
		GetDid() Did
		SetPubKey(pubkey string) error
		GetPubKey() string
		Address() sdk.AccAddress
		AddressUnverified() sdk.AccAddress
	}
	Claim struct {
		Id           Did  `json:"id" yaml:"id"`
		KYCValidated bool `json:"KYCValidated" yaml:"KYCValidated"`
	}
	DidCredential struct {
		CredType []string `json:"type" yaml:"type"`
		Issuer   Did      `json:"issuer" yaml:"issuer"`
		Issued   string   `json:"issued" yaml:"issued"`
		Claim    Claim    `json:"claim" yaml:"claim"`
	}
	Secret struct {
		Seed                 string `json:"seed" yaml:"seed"`
		SignKey              string `json:"signKey" yaml:"signKey"`
		EncryptionPrivateKey string `json:"encryptionPrivateKey" yaml:"encryptionPrivateKey"`
	}
	IxoDid struct {
		Did                 string `json:"did" yaml:"did"`
		VerifyKey           string `json:"verifyKey" yaml:"verifyKey"`
		EncryptionPublicKey string `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
		Secret              Secret `json:"secret" yaml:"secret"`
		Dpinfo              DpInfo `json:"dp" yaml:"dp"`
	}
	DpInfo struct {
		DpAddress string           `json:"address" yaml:"address"`
		PubKey    string           `json:"pubkey" yaml:"pubkey"`
		Name      string           `json:"name" yaml:"name"`
		Algo      keys.SigningAlgo `json:"algo" yaml:"algo"`
	}
	KeyGenerator struct {
		mem     string
		pubkey  []byte
		privkey []byte
		name    string
		seed    [32]byte
		debug   bool
	}
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*DidDoc)(nil), nil)
	cdc.RegisterConcrete(&IxoDid{}, "darkpool/IxoDid", nil)
	cdc.RegisterConcrete(&DpInfo{}, "darkpool/DpInfo", nil)
	cdc.RegisterConcrete(&Secret{}, "darkpool/Secret", nil)
	cdc.RegisterConcrete(&DidCredential{}, "darkpool/DidCredential", nil)
	cdc.RegisterConcrete(&Claim{}, "darkpool/Claim", nil)
}

func (id IxoDid) Equals(other IxoDid) bool {
	return id.Did == other.Did &&
		id.VerifyKey == other.VerifyKey &&
		id.EncryptionPublicKey == other.EncryptionPublicKey &&
		id.Secret.Equals(other.Secret)
}

func (id IxoDid) GetPubKeyByte() [32]byte {
	return RecoverDidToEd25519PubKey(id)
}
func (id IxoDid) GetPriKeyByte() [64]byte {
	return RecoverDidToEd25519PrivateKey(id)
}
func (id IxoDid) FromAddressDx0() sdk.AccAddress {
	address, _ := sdk.AccAddressFromBech32(id.Dpinfo.DpAddress)
	return address
}
func (id IxoDid) FromPubKeyDx0() tmcrypto.PubKey {
	address, e := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, id.Dpinfo.PubKey)
	//address, _ := sdk.AccAddressFromBech32( id.VerifyKey)
	if e != nil {
		fmt.Println("cannot get the pubkey, ", e, id.Dpinfo.PubKey)
	}
	return address
}
func (id IxoDid) Address() sdk.AccAddress {
	return VerifyKeyToAddrEd25519(id.VerifyKey)
}
func (id IxoDid) AddressEd() sdk.AccAddress {
	return UnverifiedToAddr(id.VerifyKey)
}
func (id IxoDid) DidAddress() string {
	return id.Did
}
func (id IxoDid) String() string {
	output, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v", string(output))
}
func (id IxoDid) MarshaDid() ([]byte, error) {
	t, err := json.Marshal(id)
	if err != nil {
		return nil, err
	}
	return t, nil
}
func (id IxoDid) GetPubKey() string {
	return id.VerifyKey
}

func (id IxoDid) SignMessage(msg []byte) ([]byte, error) {
	var privateKey ed25519tm.PrivKeyEd25519
	copy(privateKey[:], base58.Decode(id.Secret.SignKey))
	copy(privateKey[32:], base58.Decode(id.VerifyKey))
	return privateKey.Sign(msg)
}

func (id IxoDid) VerifySignedMessage(msg []byte, sig []byte) bool {
	var publicKey ed25519tm.PubKeyEd25519
	copy(publicKey[:], base58.Decode(id.VerifyKey))
	return publicKey.VerifyBytes(msg, sig)
}

func fromJsonStringDp(jsonSovrinDid string) (IxoDid, error) {
	var did IxoDid
	err := json.Unmarshal([]byte(jsonSovrinDid), &did)
	if err != nil {
		err := fmt.Errorf("Could not unmarshal did into struct. Dxp Error: %s! ", err.Error())
		return IxoDid{}, err
	}
	return did, nil
}

/*
func VerifyKeyToAddr(verifyKey string) sdk.AccAddress {
	var pubKey ed25519.PubKeyEd25519
	copy(pubKey[:], base58.Decode(verifyKey))
	return sdk.AccAddress(pubKey.Address())
}
func UnverifiedToAddr(ver string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(ver)))
}

func UnmarshalDxpDid(jsonSovrinDid string) (IxoDid, error) {
	return fromJsonStringDp(jsonSovrinDid)
}
*/

// Above IxoDid modelled after Sovrin documents
// Ref: https://www.npmjs.com/package/sovrin-did
// {
//    did: "<base58 did>",
//    verifyKey: "<base58 publicKey>",
//    publicKey: "<base58 publicKey>",
//
//    secret: {
//        seed: "<hex encoded 32-byte seed>",
//        signKey: "<base58 secretKey>",
//        privateKey: "<base58 privateKey>"
//    }
// }

func NewSecret(seed, signKey, encryptionPrivateKey string) Secret {
	return Secret{
		Seed:                 seed,
		SignKey:              signKey,
		EncryptionPrivateKey: encryptionPrivateKey,
	}
}

func (s Secret) Equals(other Secret) bool {
	return s.Seed == other.Seed &&
		s.SignKey == other.SignKey &&
		s.EncryptionPrivateKey == other.EncryptionPrivateKey
}

func (s Secret) String() string {
	output, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v", string(output))
}
