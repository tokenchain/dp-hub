#cd static
#grunt dev
#cd ../
#go run starbubbles.go
PATH_NOW=$HOME/Documents/ixo/dp-hub
cd $PATH_NOW
echo "========================================="
echo "==GOTO>=$PATH_NOW"
echo "========================================="
#GOPATH="/Users/hesk/go/src/"
godepgraph -novendor -s ./cmd/dpd | dot -Tpng -o $PATH_NOW/godepgraph_ixod.png
godepgraph -novendor -s ./cmd/dpcli | dot -Tpng -o $PATH_NOW/godepgraph_ixocli.png
open $PATH_NOW/godepgraph_ixod.png
open $PATH_NOW/godepgraph_ixocli.png
#go list -f '{{join .DepsErrors "\n"}}' <import-path>