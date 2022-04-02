#!/bin/bash +x

VERSION="0.0.0"
if [ -n "$1" ]; then
    VERSION=$1
fi

echo "Building version: $VERSION"


env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./dist/windows_amd64/rendr.exe main.go
if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
fi

env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./dist/darwin_amd64/rendr main.go
if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
fi

env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./dist/linux_amd64/rendr main.go
if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
fi

echo "Done building version: $VERSION"