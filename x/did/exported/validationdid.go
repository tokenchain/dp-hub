package exported

import "regexp"

var(
	ValidDid    = regexp.MustCompile(`^did:(dxp:|sov:)([a-zA-Z0-9]){21,22}([/][a-zA-Z0-9:]+|)$`)
	IsValidDid  = ValidDid.MatchString
	// https://sovrin-foundation.github.io/sovrin/spec/did-method-spec-template.html
	// IsValidDid adapted from the above link but assumes no sub-namespaces
	// TODO: ValidDid needs to be updated once we no longer want to be able
	//   to consider project accounts as DIDs (especially in treasury module),
	//   possibly should just be `^did:(dxp:|sov:)([a-zA-Z0-9]){21,22}$`.
)