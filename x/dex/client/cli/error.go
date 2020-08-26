package cli

import "github.com/cosmos/cosmos-sdk/types/errors"


var(
	Err9 = errors.Register("dex", 591, "cli errors")
)


func ErrCli(msg string) error {
	return errors.Wrapf(Err9, "dex: %s", msg)
}
