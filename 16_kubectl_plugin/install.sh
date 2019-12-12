for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    output_name="bin/kubectl-snapshot-$GOOS-$GOARCH"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    echo "Building $output_name"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
    if [ $? -ne 0 ]; then
      echo 'An error has occurred! Aborting the script execution...'
      exit 1
    fi
  done
done

chmod +x ./sh/*
cp ./sh/* ./bin/

mkdir -p ~/repos/github/fbrubbo/kubectl-plugins/bin
cp bin/* ~/repos/github/fbrubbo/kubectl-plugins/bin


