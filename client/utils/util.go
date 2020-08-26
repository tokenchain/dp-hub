package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/console"
	"github.com/spf13/viper"
	core "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tokenchain/ixo-blockchain/x/did/ante"
	"io"
	"io/ioutil"
	"time"
)

/*custom coin transaction layer*/
func GetKeybase(transient bool, buf io.Reader) (keys.Keybase, error) {
	if transient {
		return keys.NewInMemory(), nil
	}
	return keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf)
}
func QueryWithData(cliCtx context.CLIContext, format string, arg ...interface{}) ([]byte, int64, error) {
	return cliCtx.QueryWithData(fmt.Sprintf(format, arg...), nil)
}

func QueryWithDataPost(cliCtx context.CLIContext, data []byte, format string, arg ...interface{}) ([]byte, int64, error) {
	return cliCtx.QueryWithData(fmt.Sprintf(format, arg...), data)
}

func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
	return ante.DefaultTxDecoder(cdc)(txBytes)
}

func formatTxResult(cdc *codec.Codec, resTx *core.ResultTx, resBlock *core.ResultBlock) (sdk.TxResponse, error) {
	tx, err := parseTx(cdc, resTx.Tx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return sdk.NewResponseResultTx(resTx, tx, resBlock.Block.Time.Format(time.RFC3339)), nil
}

func getBlocksForTxResults(cliCtx context.CLIContext, resTxs []*core.ResultTx) (map[int64]*core.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resBlocks := make(map[int64]*core.ResultBlock)
	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}

func ValidateTxResult(cliCtx context.CLIContext, resTx *core.ResultTx) error {
	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(resTx.Height)
		if err != nil {
			return err
		}

		err = resTx.Proof.Validate(check.Header.DataHash)
		if err != nil {
			return err
		}
	}

	return nil
}

func QueryTx(cliCtx context.CLIContext, hashHexStr string) (sdk.TxResponse, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	resTx, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	if !cliCtx.TrustNode {
		if err = ValidateTxResult(cliCtx, resTx); err != nil {
			return sdk.TxResponse{}, err
		}
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, []*core.ResultTx{resTx})
	if err != nil {
		return sdk.TxResponse{}, err
	}

	out, err := formatTxResult(cliCtx.Codec, resTx, resBlocks[resTx.Height])
	if err != nil {
		return out, err
	}

	return out, nil
}

// getPassPhrase retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func GetPassPhrase(prompt string, confirmation bool, i int, passwords []string) string {
	// If a list of passwords was supplied, retrieve from them
	if len(passwords) > 0 {
		if i < len(passwords) {
			return passwords[i]
		}
		return passwords[len(passwords)-1]
	}
	// Otherwise prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}
	password, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		utils.Fatalf("Failed to read passphrase: %v", err)
	}
	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			utils.Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if password != confirm {
			utils.Fatalf("Passphrases do not match")
		}
	}
	return password
}

type DelistProposalJSON struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	BaseAsset   string       `json:"base_asset" yaml:"base_asset"`
	QuoteAsset  string       `json:"quote_asset" yaml:"quote_asset"`
	Deposit     sdk.DecCoins `json:"deposit" yaml:"deposit"`
}

// ParseDelistProposalJSON parse json from proposal file to DelistProposalJSON struct
func ParseDelistProposalJSON(cdc *codec.Codec, proposalFilePath string) (proposal DelistProposalJSON, err error) {
	contents, err := ioutil.ReadFile(proposalFilePath)
	if err != nil {
		return proposal, err
	}
	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}
	return proposal, nil
}
func ParseDecCoinRounded(coins sdk.DecCoins) sdk.Coins {
	coinName := coins.GetDenomByIndex(0)
	outAmt := coins.AmountOf(coinName)
	return sdk.Coins{sdk.NewCoin(coinName, outAmt.RoundInt())}
}
func ParseDecCoinSingleRounded(coin sdk.DecCoin) sdk.Coin {
	return sdk.NewCoin(coin.Denom, coin.Amount.RoundInt())
}
func ListParsedCoinRounded(coin sdk.DecCoin) sdk.Coins {
	return sdk.Coins{ParseDecCoinSingleRounded(coin)}
}

func MustParseCoins(symbol, amount string) sdk.Coins {
	coinAmt, _ := sdk.NewIntFromString(amount)
	coins := sdk.Coins{sdk.NewCoin(symbol, coinAmt)}
	return coins
}
