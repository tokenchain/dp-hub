package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

var _ DidDoc = (*BaseDidDoc)(nil)

type (
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
	DidDoc interface {
		SetDid(did Did) error
		GetDid() Did
		SetPubKey(pubkey string) error
		GetPubKey() string
	}
	BaseDidDoc struct {
		Did         Did             `json:"did" yaml:"did"`
		PubKey      string          `json:"pubKey" yaml:"pubKey"`
		Credentials []DidCredential `json:"credentials" yaml:"credentials"`
	}
	DxpDid struct {
		Did                 string       `json:"did" yaml:"did"`
		VerifyKey           string       `json:"verifyKey" yaml:"verifyKey"`
		EncryptionPublicKey string       `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
		Secret              SovrinSecret `json:"secret" yaml:"secret"`
	}
	SovrinSecret struct {
		Seed                 string `json:"seed" yaml:"seed"`
		SignKey              string `json:"signKey" yaml:"signKey"`
		EncryptionPrivateKey string `json:"encryptionPrivateKey" yaml:"encryptionPrivateKey"`
	}
	Did = string
)
type Credential struct{}

func (ss SovrinSecret) String() string {
	output, err := json.MarshalIndent(ss, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

func (sd DxpDid) String() string {
	output, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

func NewBaseDidDoc(did Did, pubKey string) BaseDidDoc {
	return BaseDidDoc{
		Did:         did,
		PubKey:      pubKey,
		Credentials: []DidCredential{},
	}
}

func (dd BaseDidDoc) GetDid() Did                     { return dd.Did }
func (dd BaseDidDoc) GetPubKey() string               { return dd.PubKey }
func (dd BaseDidDoc) GetCredentials() []DidCredential { return dd.Credentials }
func (dd BaseDidDoc) SetDid(did Did) error {
	if len(dd.Did) != 0 {
		return errors.New("cannot override BaseDidDoc did")
	}

	dd.Did = did

	return nil
}
func (dd BaseDidDoc) SetPubKey(pubKey string) error {
	if len(dd.PubKey) != 0 {
		return errors.New("cannot override BaseDidDoc pubKey")
	}

	dd.PubKey = pubKey

	return nil
}
func (dd *BaseDidDoc) AddCredential(cred DidCredential) {
	if dd.Credentials == nil {
		dd.Credentials = make([]DidCredential, 0)
	}

	dd.Credentials = append(dd.Credentials, cred)
}
