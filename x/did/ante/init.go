package ante

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
	ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	types2 "github.com/tokenchain/ixo-blockchain/x/did/internal/types"
)

const (
	Ed25519SignatureLen = 64
)

var (
	// simulation signature values used to estimate gas consumption
	simEd25519Pubkey   ed25519tm.PubKeyEd25519
	simEd25519Sig      [Ed25519SignatureLen]byte
	simSecp256k1Pubkey secp256k1.PubKeySecp256k1
	simSecp256k1Sig    [Ed25519SignatureLen]byte
)

func init() {
	// This decodes a valid hex string into a ed25519Pubkey for use in transaction simulation
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E14")
	copy(simEd25519Pubkey[:], bz)
	//copy(simSecp256k1Pubkey[:], bz)
}

type PubKeyGetter func(ctx sdk.Context, msg types2.IxoMsg) (crypto.PubKey, error)
type SigVerificationGasConsumer func(meter sdk.GasMeter, sig []byte, pubkey crypto.PubKey, params types.Params) error
