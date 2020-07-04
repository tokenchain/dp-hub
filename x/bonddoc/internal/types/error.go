package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tokenchain/ixo-blockchain/x/did"
)
import 	"github.com/tokenchain/ixo-blockchain/x"


func ErrFromAndToCannotBeTheSameToken() error {
	return errors.Wrap(x.ErrFromAndToCannotBeTheSameToken_E, "From and To tokens cannot be the same token.")
}

func ErrBondDoesNotExist(bondDid did.Did) error {
	return errors.Wrapf(x.ErrCodeBondDoesNotExist, "Bond '%s' does not exist", bondDid)
}
func ErrBondAlreadyExists(bonddid did.Did) error {
	return errors.Wrapf(x.ErrCodeBondAlreadyExists, "Bond '%s' already exists", bonddid)
}
func ErrBondTokenIsTaken(bondToken string) error {
	return errors.Wrapf(x.ErrCodeBondAlreadyExists, "Bond token '%s' is taken", bondToken)
}
func ErrBondDoesNotAllowSelling() error {
	return errors.Wrap(x.ErrCodeBondDoesNotAllowSelling, "Bond does not allow selling.")
}
func ErrDidNotEditAnything() error {
	return errors.Wrap(x.ErrCodeDidNotEditAnything, "Did not edit anything from the bond.")
}
