package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/ixo-blockchain/x"
)

func InvalidTxDecodePubkeyNotFound(e error) error {
	return errors.Wrapf(errors.ErrTxDecode, "retrieve Darkpool pubkey failure %s! ", e.Error())
}
func InvalidTxDecode() error {
	return errors.Wrap(errors.ErrTxDecode, "invalid tx type")
}
func UnknownRequest(m string) error {
	return errors.Wrap(errors.ErrUnknownRequest, m)
}
func Unauthorized(m string) error {
	return errors.Wrap(errors.ErrUnauthorized, m)
}
func Unauthorizedf(format string, a ...interface{}) error {
	return errors.Wrap(errors.ErrUnauthorized, fmt.Sprintf(format, a...))
}
func IntErr(m string) error {
	return errors.Wrap(errors.ErrPanic, m)
}
func ErrJsonMars(m string) error {
	return errors.Wrapf(errors.ErrJSONMarshal, "Json marshall error %s", m)
}

func InvalidPubKey(m string) error {
	return errors.Wrapf(errors.ErrInvalidPubKey, "PubKey error %s", m)
}

func ErrInvalidBasicMsg(msg string) error {
	errMsg := fmt.Sprintf("invalid basic message %s", msg)
	return errors.Wrap(x.ErrInvalidBasicMsg, errMsg)
}

func ErrBadDataValue(msg string) error {
	errMsg := fmt.Sprintf("bad data value %s", msg)
	return errors.Wrap(x.ErrBadDataValue, errMsg)
}

func ErrUnauthorizedPermission(msg string) error {
	errMsg := fmt.Sprintf("permission deny %s", msg)

	return errors.Wrap(x.ErrUnauthorizedPermission, errMsg)
}

func ErrItemDuplication(msg string) error {
	errMsg := fmt.Sprintf("item duplicated %s", msg)

	return errors.Wrap(x.ErrItemDuplication, errMsg)
}

func ErrItemNotFound(msg string) error {
	errMsg := fmt.Sprintf("item is not found %s", msg)

	return errors.Wrap(x.ErrItemNotFound, errMsg)
}

func ErrInvalidState(msg string) error {
	errMsg := fmt.Sprintf("invalid state %s", msg)

	return errors.Wrap(x.ErrInvalidState, errMsg)
}

func ErrBadWasmExecution(msg string) error {
	errMsg := fmt.Sprintf("bad wasm execution %s", msg)

	return errors.Wrap(x.ErrBadWasmExecution, errMsg)
}

func ErrOnlyOneDenomAllowed(msg string) error {
	errMsg := fmt.Sprintf("only one denom allowed %s", msg)

	return errors.Wrap(x.ErrOnlyOneDenomAllowed, errMsg)
}

func ErrInvalidDenom(msg string) error {
	errMsg := fmt.Sprintf("invalid denom %s", msg)

	return errors.Wrap(x.ErrInvalidDenom, errMsg)
}

func ErrUnknownClientID() error {
	errMsg := fmt.Sprintf("unknown client ID")

	return errors.Wrap(x.ErrUnknownClientID, errMsg)
}
