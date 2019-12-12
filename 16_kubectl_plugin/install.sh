for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o bin/kubectl-snapshot-$GOOS-$GOARCH
  done
done

cd bin
rm -rf kubectl-snapshot.tar.gz
tar -cvzf kubectl-snapshot.tar.gz *
cd ..

go install
cp bin/kubectl-snapshot-linux-amd64 ~/repos/github/fbrubbo/kubectl-plugins/kubectl-snapshot
chmod +x ./sh/*
cp ./sh/* ~/repos/github/fbrubbo/kubectl-plugins/ -f