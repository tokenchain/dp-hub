package exported

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func IntErr(m string) error {
	return errors.Wrap(errors.ErrPanic, m)
}
func UnknownRequest(m string) error {
	return errors.Wrap(errors.ErrUnknownRequest, m)
}
func Invalid(r string) error {
	return errors.Wrap(ErrorInvalid, r)
}
func Invalidf(r string, a ...interface{}) error {
	return errors.Wrap(ErrorInvalid, fmt.Sprintf(r, a...))
}
func InvalidCredentials(ar string) error {
	return errors.Wrap(ErrorInvalidCredentials, ar)
}
func InvalidCoinsf(r string, a ...interface{}) error {
	return errors.Wrap(EInvalidCoin, fmt.Sprintf(r, a...))
}

func InvalidCoins(args string) error {
	return errors.Wrap(EInvalidCoin, args)
}
func Unauthorized(m string) error {
	return errors.Wrap(errors.ErrUnauthorized, m)
}
func InvalidTxDecode() error {
	return errors.Wrap(errors.ErrTxDecode, "invalid tx type")
}
func InvalidTxDecodeMsg(ar string) error {
	return errors.Wrap(errors.ErrTxDecode, ar)
}
func InvalidAddress(r string) error {
	return errors.Wrapf(errors.ErrInvalidAddress, r)
}
func InvalidIssuer(ar string) error {
	return errors.Wrap(ErrInvalidIssuer, ar)
}

func InvalidPubKey(ar string) error {
	return errors.Wrap(ErrInvalidPubKey, ar)
}

func InvalidDidMsg(r string) error {
	return errors.Wrap(EInvalidId, r)
}

func ErrJsonMars(m string) error {
	return errors.Wrapf(errors.ErrJSONMarshal, "Json marshall error %s", m)
}
