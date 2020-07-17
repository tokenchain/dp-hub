package rest

import (
	"fmt"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	didkeep "github.com/tokenchain/ixo-blockchain/x/did/internal/keeper"
	didtypes "github.com/tokenchain/ixo-blockchain/x/did/internal/types"
	"net/http"
)

func writeHeadf(w http.ResponseWriter, code int, format string, i ...interface{}) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(fmt.Sprintf(format, i...)))
}
func writeHead(w http.ResponseWriter, code int, txt string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(txt))
}

func getQuery(did_key exported.Did) string {
	return fmt.Sprintf("custom/%s/%s/%s", didtypes.QuerierRoute, didkeep.QueryDidDoc, did_key)
}
func getQuerySr() string {
	return fmt.Sprintf("custom/%s/%s/", didtypes.QuerierRoute, didkeep.QueryDidDoc)
}
func getQueryAll() string {
	return fmt.Sprintf("custom/%s/%s/", didtypes.QuerierRoute, didkeep.QueryAllDidDocs)
}
