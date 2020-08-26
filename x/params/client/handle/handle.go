package handle

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/tokenchain/ixo-blockchain/x/params/client/cli"
	"github.com/tokenchain/ixo-blockchain/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler in cmsdk
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
