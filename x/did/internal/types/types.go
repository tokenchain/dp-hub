package types

import (
	"errors"
	"github.com/tokenchain/ixo-blockchain/x/ixo/types"
)

var _ types.DidDoc = (*BaseDidDoc)(nil)

type BaseDidDoc struct {
	Did         types.Did       `json:"did" yaml:"did"`
	PubKey      string          `json:"pubKey" yaml:"pubKey"`
	Credentials []DidCredential `json:"credentials" yaml:"credentials"`
}

func NewBaseDidDoc(did types.Did, pubKey string) BaseDidDoc {
	return BaseDidDoc{
		Did:         did,
		PubKey:      pubKey,
		Credentials: []DidCredential{},
	}
}

func (dd BaseDidDoc) GetDid() types.Did { return dd.Did }
func (dd BaseDidDoc) GetPubKey() string { return dd.PubKey }
func (dd BaseDidDoc) GetCredentials() []DidCredential { return dd.Credentials }

func (dd BaseDidDoc) SetDid(did types.Did) error {
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

type DidCredential struct {
	CredType []string  `json:"type" yaml:"type"`
	Issuer   types.Did `json:"issuer" yaml:"issuer"`
	Issued   string    `json:"issued" yaml:"issued"`
	Claim    Claim     `json:"claim" yaml:"claim"`
}

type Claim struct {
	Id           types.Did `json:"id" yaml:"id"`
	KYCValidated bool      `json:"KYCValidated" yaml:"KYCValidated"`
}

type Credential struct{}
