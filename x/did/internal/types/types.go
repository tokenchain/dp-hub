package types

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tokenchain/dp-block/x/did/exported"
)

var _ exported.DidDoc = (*BaseDidDoc)(nil)

type BaseDidDoc struct {
	Did         exported.Did             `json:"did" yaml:"did"`
	PubKey      string                   `json:"pubKey" yaml:"pubKey"` //that also is the verify key
	Credentials []exported.DidCredential `json:"credentials,omitempty" yaml:"credentials"`
}

func NewBaseDidDoc(did exported.Did, pubKey string) BaseDidDoc {
	return BaseDidDoc{
		Did:         did,
		PubKey:      pubKey,
		Credentials: []exported.DidCredential{},
	}
}
func (dd BaseDidDoc) GetDid() exported.Did                     { return dd.Did }
func (dd BaseDidDoc) GetPubKey() string                        { return dd.PubKey }
func (dd BaseDidDoc) GetCredentials() []exported.DidCredential { return dd.Credentials }
func (dd BaseDidDoc) SetDid(did exported.Did) error {
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
func (dd *BaseDidDoc) AddCredential(cred exported.DidCredential) {
	if dd.Credentials == nil {
		dd.Credentials = make([]exported.DidCredential, 0)
	}
	dd.Credentials = append(dd.Credentials, cred)
}
func (dd BaseDidDoc) Address() sdk.AccAddress {
	return exported.VerifyKeyToAddrEd25519(dd.GetPubKey())
}
func (dd BaseDidDoc) AddressUnverified() sdk.AccAddress {
	return exported.UnverifiedToAddr(dd.GetPubKey())
}
type Credential struct{}
func fromJsonString(jsonIxoDid string) (exported.IxoDid, error) {
	var did exported.IxoDid
	err := json.Unmarshal([]byte(jsonIxoDid), &did)
	if err != nil {
		err := fmt.Errorf("Could not unmarshal did into struct. Error: %s .", err.Error())
		return exported.IxoDid{}, err
	}

	return did, nil
}
func UnmarshalIxoDid(jsonIxoDid string) (exported.IxoDid, error) {
	return fromJsonString(jsonIxoDid)
}
