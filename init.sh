#sudo chown -R $USER: $HOME
#sudo rm -rf $GOPATH/pkg/mod/
echo "Removing all packages..."
sudo go clean -modcache
#ls $GOPATH/pkg/mod
mkdir $HOME/go/pkg/mod
sudo chmod -R 755 $HOME/go/pkg/mod
sleep 3
echo "Donwloading packages..."
go mod download
sleep 3
echo "list all..."
go list -m -json all