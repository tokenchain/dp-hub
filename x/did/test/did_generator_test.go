package test

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/go-bip39"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tokenchain/dp-hub/x/did/exported"
	"testing"
)

func Test_recover_ke(t *testing.T) {
	setPrefix()
	name := "cosmos"
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	require.Nil(t, err, "memory phrase generated")
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	did_document := exported.NewDidGeneratorBuilder().WithMem(mnemonic).WithName(name).Build()
	fmt.Println("========🔑 Account Info Save This To a secured place===============================================")
	fmt.Println(did_document)

	fmt.Println("=======🔑  The passphrase please keep in the secured place. Its an important private key! ")
	fmt.Println(mnemonic)
}
func Test_recover_public_key(t *testing.T) {
	setPrefix()
	doc := exported.NewDapDid(
		sample_did,
		sample_verifyKey,
		sample_encryptionPublicKey,
		sample_seed,
		sample_signKey,
		sample_encryptionPrivateKey,
		sample_address,
		sample_pub,
		"cosmos",
	)

	//err := dap.SignAndBroadcastTxCli()

	/*

		var recover_pub [32]byte
		name := substring(doc.Did, 8, len(doc.Did))
		p1 := base58.Decode(name)
		p2 := base58.Decode(doc.EncryptionPublicKey)
		copy(recover_pub[:], p1)
		copy(recover_pub[16:], p2)

		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
		config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
		config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
		config.Seal()

	*/
	//both parts are filled
	r := exported.RecoverDidToEd25519PubKey(doc)
	//only the first part filled
	//f := exported.RecoverDidEd25519ToPrivateKey(doc)
	//get the private key
	fmt.Println(doc.FromAddressDx0().String())
	fmt.Println(doc.Address().String())

	fmt.Println("========key recover  pubkey from pub key=========")
	fmt.Println(doc.FromPubKeyDx0())
	fmt.Println(doc.GetPubKey())

	fmt.Println("========key recover  pubkey=========")
	fmt.Println(doc.GetPubKeyByte())

	fmt.Println("========key recover  private key=========")
	fmt.Println(doc.GetPriKeyByte())
	pri := doc.GetPriKeyByte()
	fmt.Println(len(pri), pri)
	require.Equal(t, r, ed25519.PrivKeyEd25519(doc.GetPriKeyByte()))
	fmt.Println("========it is now name  =========")
	//	fmt.Println(name, doc.Address().String())

	fmt.Println("========key recover  =========")
	fmt.Println(len(r), r)

	fmt.Println("========key recover  =========")

	fmt.Println("========key recover  =========")
	require.Nil(t, nil, "successfull signed")
}
func Test_generate_recover(t *testing.T) {
	setPrefix()
	fmt.Println("========Seed==========")
	did := exported.NewDidGeneratorBuilder().WithMem(sample_did_01_mem).Build()

	fmt.Println("========Key Recover  =========")
	fmt.Println(did)
	fmt.Println(did.Did)

	fmt.Println("========Check The Equal ===========")
	require.Equal(t, did.Did, sample_did)
}
func Test_generator(t *testing.T) {
	/*

	ADDRESS_SINGULARITY=$(jq '.dp.address' $DID_FOLDER/did_singularity.json -r)
	DIDSOVRIN_SINGULARITY=$(jq -c . $DID_FOLDER/did_singularity.json)
	DID_SINGULARITY=$(jq '.did' $DID_FOLDER/did_singularity.json -r)

	*/
	setPrefix()
	fmt.Println("========Seed==========")
	total_accounts := uint32(20)
	account_index := uint32(177)
	list := []string{
		"singularity",
		"blackhole",
		"cosmos",
		"cosmic",
		"darkness",
		"nova",
		"proton",
		"rednova",
		"bitcm",
		"dollar",
		"kbs",
	}
	var i uint32
	var name string
	for i = 1; i < total_accounts; i++ {

		if i < uint32(len(list)+1) {
			name = list[i-1]
		} else {
			name = "cosmos---"
			continue
		}

		did := exported.NewDidGeneratorBuilder().
			WithName(name).
			RecoverBIP44(sample_did_01_mem, "", account_index, i)


		jsonString, _ := json.Marshal(did)
		makeFile(name, jsonString)



		fmt.Println("========Check The DID ===========")
		fmt.Println(did.Did)
		fmt.Println("========DID Recover  =========")
		fmt.Println(did)
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
	}
}
