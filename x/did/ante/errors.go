package ante

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/dp-hub/x/did/exported"
)

func InvalidTxDecodePubkeyNotFound(e error) error {
	return errors.Wrapf(errors.ErrTxDecode, "retrieve Darkpool pubkey failure %s! ", e.Error())
}
func InvalidTxDecode() error {
	return errors.Wrap(errors.ErrTxDecode, "invalid tx type")
}
func InvalidTxDecodeMsg(ar string) error {
	return errors.Wrap(errors.ErrTxDecode, ar)
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

func UnknownAddress(m string) error {
	return errors.Wrap(errors.ErrUnknownAddress, m)
}
func UnknownAddressf(format string, a ...interface{}) error {
	return errors.Wrap(errors.ErrUnknownAddress, fmt.Sprintf(format, a...))
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
func InvalidPubKeyf(format string, a ...interface{}) error {
	return errors.Wrapf(errors.ErrInvalidPubKey, format, a)
}
func ErrInvalidBasicMsg(msg string) error {
	errMsg := fmt.Sprintf("invalid basic message %s", msg)
	return errors.Wrap(exported.ErrInvalidBasicMsg, errMsg)
}

func ErrBadDataValue(msg string) error {
	errMsg := fmt.Sprintf("bad data value %s", msg)
	return errors.Wrap(exported.ErrBadDataValue, errMsg)
}

func ErrUnauthorizedPermission(msg string) error {
	errMsg := fmt.Sprintf("permission deny %s", msg)

	return errors.Wrap(exported.ErrUnauthorizedPermission, errMsg)
}

func ErrItemDuplication(msg string) error {
	errMsg := fmt.Sprintf("item duplicated %s", msg)

	return errors.Wrap(exported.ErrItemDuplication, errMsg)
}

func ItemSigNotFound(msg string) error {
	errMsg := fmt.Sprintf("No multi signatures found %s", msg)
	return errors.Wrap(exported.ErrItemNotFound, errMsg)
}
func ErrItemNotFound(msg string) error {
	errMsg := fmt.Sprintf("item is not found %s", msg)
	return errors.Wrap(exported.ErrItemNotFound, errMsg)
}

func ErrInvalidState(msg string) error {
	errMsg := fmt.Sprintf("invalid state %s", msg)

	return errors.Wrap(exported.ErrInvalidState, errMsg)
}

func ErrBadWasmExecution(msg string) error {
	errMsg := fmt.Sprintf("bad wasm execution %s", msg)

	return errors.Wrap(exported.ErrBadWasmExecution, errMsg)
}

func ErrOnlyOneDenomAllowed(msg string) error {
	errMsg := fmt.Sprintf("only one denom allowed %s", msg)

	return errors.Wrap(exported.ErrOnlyOneDenomAllowed, errMsg)
}

func ErrInvalidDenom(msg string) error {
	errMsg := fmt.Sprintf("invalid denom %s", msg)

	return errors.Wrap(exported.ErrInvalidDenom, errMsg)
}

func ErrUnknownClientID() error {
	errMsg := fmt.Sprintf("unknown client ID")

	return errors.Wrap(exported.ErrUnknownClientID, errMsg)
}
