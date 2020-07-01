package types

import (
	"encoding/json"
	"fmt"
	"github.com/tokenchain/dp-hub/x"
	"github.com/tokenchain/dp-hub/x/ixo/types"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	er "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgAddDid        = "add-did"
	TypeMsgAddCredential = "add-credential"
)

var (
	_ types.IxoMsg = MsgAddDid{}
	_ types.IxoMsg = MsgAddCredential{}
)

type MsgAddDid struct {
	DidDoc BaseDidDoc `json:"didDoc" yaml:"didDoc"`
}

func (msg *MsgAddDid) UnmarshalJSON(bytes []byte) error {
	var msg2 struct {
		DidDoc BaseDidDoc `json:"didDoc" yaml:"didDoc"`
	}
	err := json.Unmarshal(bytes, &msg2)
	if err != nil {
		return err
	}

	if msg2.DidDoc.Credentials == nil {
		msg2.DidDoc.Credentials = []DidCredential{}
	}

	*msg = msg2
	return nil
}

func NewMsgAddDid(did string, publicKey string) MsgAddDid {
	return MsgAddDid{
		DidDoc: NewBaseDidDoc(did, publicKey),
	}
}

func (msg MsgAddDid) Type() string            { return TypeMsgAddDid }
func (msg MsgAddDid) Route() string           { return RouterKey }
func (msg MsgAddDid) GetSignerDid() types.Did { return msg.DidDoc.Did }
func (msg MsgAddDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}
func (msg MsgAddDid) ValidateBasic() error {
	// Check that not empty
	if strings.TrimSpace(msg.DidDoc.Did) == "" {
		return er.Wrap(x.ErrorInvalidDidE, "did should not be empty")
	} else if strings.TrimSpace(msg.DidDoc.PubKey) == "" {
		return er.Wrap(x.ErrorInvalidPubKey, "pubKey should not be empty")
	}

	// Check DidDoc credentials for empty fields
	for _, cred := range msg.DidDoc.Credentials {
		if strings.TrimSpace(cred.Issuer) == "" {
			return er.Wrap(x.ErrorInvalidIssuer, "issuer should not be empty")
		} else if strings.TrimSpace(cred.Claim.Id) == "" {
			return er.Wrap(x.ErrorInvalidDidE, "claim id should not be empty")
		}
	}

	// Check that DID valid
	if !types.IsValidDid(msg.DidDoc.Did) {
		return er.Wrap(x.ErrorInvalidDidE, "did is invalid")
	}

	return nil
}

func (msg MsgAddDid) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func (msg MsgAddDid) String() string {
	return fmt.Sprintf("MsgAddDid{Did: %v, publicKey: %v}", string(msg.DidDoc.GetDid()), msg.DidDoc.GetPubKey())
}

type MsgAddCredential struct {
	DidCredential DidCredential `json:"credential" yaml:"credential"`
}

func NewMsgAddCredential(did string, credType []string, issuer string, issued string) MsgAddCredential {
	didCredential := DidCredential{
		CredType: credType,
		Issuer:   issuer,
		Issued:   issued,
		Claim: Claim{
			Id:           did,
			KYCValidated: true,
		},
	}

	return MsgAddCredential{
		DidCredential: didCredential,
	}
}
func (msg MsgAddCredential) Type() string            { return TypeMsgAddCredential }
func (msg MsgAddCredential) Route() string           { return RouterKey }
func (msg MsgAddCredential) GetSignerDid() types.Did { return msg.DidCredential.Issuer }
func (msg MsgAddCredential) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.DidToAddr(msg.GetSignerDid())}
}
func (msg MsgAddCredential) String() string {
	return fmt.Sprintf("MsgAddCredential{Did: %v, Type: %v, Signer: %v}",
		string(msg.DidCredential.Claim.Id), msg.DidCredential.CredType, string(msg.DidCredential.Issuer))
}
func (msg MsgAddCredential) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.DidCredential.Claim.Id) == "" {
		return er.Wrap(x.ErrorInvalidDidE, "claim id should not be empty")
	} else if strings.TrimSpace(msg.DidCredential.Issuer) == "" {
		return er.Wrap(x.ErrorInvalidIssuer, "issuer should not be empty")
	}
	// Check that DID valid
	if !types.IsValidDid(msg.DidCredential.Issuer) {
		return er.Wrap(x.ErrorInvalidDidE, "issuer id is invalid")
	}
	return nil
}
func (msg MsgAddCredential) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}
