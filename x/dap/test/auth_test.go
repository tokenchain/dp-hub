package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

func TestValidateBasic(t *testing.T) {
	// setup
	_, ctx := createTestApp(true)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()

	// msg and signatures
	msg1 := types.NewTestMsg(addr1)
	fee := types.NewTestStdFee()

	msgs := []sdk.Msg{msg1}

	privs, accNums, seqs := []crypto.PrivKey{}, []uint64{}, []uint64{}
	invalidTx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	vbd := ante.NewValidateBasicDecorator()
	antehandler := sdk.ChainAnteDecorators(vbd)
	_, err := antehandler(ctx, invalidTx, false)

	require.NotNil(t, err, "Did not error on invalid tx")

	privs, accNums, seqs = []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	validTx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	_, err = antehandler(ctx, validTx, false)
	require.Nil(t, err, "ValidateBasicDecorator returned error on valid tx. err: %v", err)

	// test decorator skips on recheck
	ctx = ctx.WithIsReCheckTx(true)

	// decorator should skip processing invalidTx on recheck and thus return nil-error
	_, err = antehandler(ctx, invalidTx, false)

	require.Nil(t, err, "ValidateBasicDecorator ran on ReCheck")
}

