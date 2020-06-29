#cd static
#grunt dev
#cd ../
#go run starbubbles.go
PATH_NOW=$HOME/Documents/ixo/ixo-blockchain
cd $PATH_NOW
echo "========================================="
echo "==GOTO>=$PATH_NOW"
echo "========================================="
#GOPATH="/Users/hesk/go/src/"
godepgraph -novendor -s ./cmd/ixod | dot -Tpng -o $PATH_NOW/xxxl/godepgraph_ixod.png
godepgraph -novendor -s ./cmd/ixocli | dot -Tpng -o $PATH_NOW/xxxl/godepgraph_ixocli.png
open $PATH_NOW/xxxl/godepgraph_ixod.png
open $PATH_NOW/xxxl/godepgraph_ixocli.png
#go list -f '{{join .DepsErrors "\n"}}' <import-path>