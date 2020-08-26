package types

import (
	"testing"

	"github.com/tokenchain/ixo-blockchain/x/common"

	"github.com/tokenchain/ixo-blockchain/x/params"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	param := DefaultParams()
	expectedString := `Params: 
FeeIssue: 2500.00000000` + common.NativeToken + `
FeeMint: 10.00000000` + common.NativeToken + `
FeeBurn: 10.00000000` + common.NativeToken + `
FeeModify: 0.00000000` + common.NativeToken + `
FeeChown: 10.00000000` + common.NativeToken + `
`
	paramStr := param.String()
	require.EqualValues(t, expectedString, paramStr)

	psp := params.ParamSetPairs{
		{Key: KeyFeeIssue, Value: &param.FeeIssue},
		{Key: KeyFeeMint, Value: &param.FeeMint},
		{Key: KeyFeeBurn, Value: &param.FeeBurn},
		{Key: KeyFeeModify, Value: &param.FeeModify},
		{Key: KeyFeeChown, Value: &param.FeeChown},
	}

	require.EqualValues(t, psp, param.ParamSetPairs())
}
