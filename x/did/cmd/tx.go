package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	base58 "github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ixofoundation/ixo-cosmos/x/did"
	"github.com/ixofoundation/ixo-cosmos/x/ixo"
	"github.com/ixofoundation/ixo-cosmos/x/ixo/sovrin"
)

// Add a did doc to the ledger
func AddDidDocCmd(cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "addDidDoc sovrinDid",
		Short: "Add a new SovrinDid",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || len(args[0]) == 0 {
				return errors.New("You must provide the sovrin did document as generated by 'sovrin-did' node package")
			}

			sovrinDid := sovrin.SovrinDid{}
			err := json.Unmarshal([]byte(args[0]), &sovrinDid)
			if err != nil {
				panic(err)
			}

			ctx := context.NewCoreContextFromViper()
			// create the message
			msg := did.NewAddDidMsg(sovrinDid.Did, sovrinDid.VerifyKey, sovrinDid.Kyc)

			// Force the length to 64
			privKey := [64]byte{}
			copy(privKey[:], base58.Decode(sovrinDid.Secret.SignKey))
			copy(privKey[32:], base58.Decode(sovrinDid.VerifyKey))

			//Create the Signature
			signature := ixo.SignIxoMessage(msg, sovrinDid.Did, privKey)

			tx := ixo.NewIxoTx(msg, signature)

			fmt.Println("*******Transaction*******")
			fmt.Println(tx.String())

			bz, err := cdc.MarshalJSON(tx)
			if err != nil {
				panic(err)
			}
			// Broadcast to Tendermint
			res, err := ctx.BroadcastTx(bz)
			if err != nil {
				return err
			}

			fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			return nil
		},
	}
}

// Get a did doc to the ledger
func GetDidDocCmd(storeName string, cdc *wire.Codec, decoder did.DidDocDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "getDidDoc did",
		Short: "Get a new DidDoc for a Did",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || len(args[0]) == 0 {
				return errors.New("You must provide an did")
			}

			// find the key to look up the account
			didAddr := args[0]
			key := ixo.Did(didAddr)

			ctx := context.NewCoreContextFromViper()

			res, err := ctx.Query([]byte(key), storeName)
			if err != nil {
				return err
			}

			// decode the value
			didDoc, err := decoder(res)
			if err != nil {
				return err
			}
			// print out whole account
			output, err := json.MarshalIndent(didDoc, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))

			return nil
		},
	}
}
