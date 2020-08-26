package handle

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/tokenchain/ixo-blockchain/x/dex/client/cli"
	"github.com/tokenchain/ixo-blockchain/x/dex/client/rest"
)

// param change proposal handler
var (
	// DelistProposalHandler alias gov NewProposalHandler
	DelistProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitDelistProposal, rest.DelistProposalRESTHandler)
)
