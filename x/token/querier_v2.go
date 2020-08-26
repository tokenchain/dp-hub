package token

import (
	"fmt"
	did "github.com/tokenchain/ixo-blockchain/x/did/exported"
	"github.com/tokenchain/ixo-blockchain/x/token/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tokenchain/ixo-blockchain/x/common"
	"github.com/tokenchain/ixo-blockchain/x/token/types"
)

func queryAccountV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, did.InvalidAddress(fmt.Sprintf("invalid address：%s", path[0]))
	}

	//var queryPage QueryPage
	var accountParam types.AccountParamV2
	//var symbol string
	err = codec.Cdc.UnmarshalJSON(req.Data, &accountParam)
	if err != nil {
		return nil, did.UnknownRequest(err.Error())
	}

	coinsInfo := keeper.GetCoinsInfo(ctx, addr)
	coinsInfoChosen := make([]CoinInfo, 0)
	if accountParam.Currency == "" {
		coinsInfoChosen = coinsInfo

		// hide_zero yes or no
		if accountParam.HideZero == "no" {
			tokens := keeper.GetTokensInfo(ctx)

			for _, token := range tokens {
				found := false
				for _, coinInfo := range coinsInfo {
					if coinInfo.Symbol == token.Symbol {
						found = true
						break
					}
				}
				// not found
				if !found {
					ci := types.NewCoinInfo(token.Symbol, "0", "0")
					coinsInfoChosen = append(coinsInfoChosen, *ci)
				}
			}
		}
	} else {
		for _, coinInfo := range coinsInfo {
			if coinInfo.Symbol == accountParam.Currency {
				coinsInfoChosen = append(coinsInfoChosen, coinInfo)
			}
		}
	}

	res, err := common.JSONMarshalV2(coinsInfoChosen)
	if err != nil {
		return nil, did.ErrJsonMars(err.Error())
	}
	return res, nil
}

func queryTokensV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	tokens := keeper.GetTokensInfo(ctx)

	var tokensResp types.Tokens
	for _, token := range tokens {
		tokenResp := types.GenTokenResp(token)
		tokenResp.TotalSupply = keeper.GetTokenTotalSupply(ctx, token.Symbol)
		tokensResp = append(tokensResp, tokenResp)
	}
	res, err := common.JSONMarshalV2(tokensResp)
	if err != nil {
		return nil, did.ErrJsonMars(err.Error())
	}
	return res, nil
}

func queryTokenV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {
	name := path[0]

	token := keeper.GetTokenInfo(ctx, name)
	if token.Symbol == "" {
		return nil, did.ErrInvalidCoins("unknown token")
	}

	tokenResp := types.GenTokenResp(token)
	tokenResp.TotalSupply = keeper.GetTokenTotalSupply(ctx, name)
	res, err := common.JSONMarshalV2(tokenResp)
	if err != nil {
		return nil, did.ErrJsonMars(err.Error())
	}
	return res, nil
}
