#!/bin/bash
set -e

if [ -n "$1" ]; then
    VERSION=$1
else
    echo 'Version not set'
    exit 1
fi

echo "Releasing $VERSION"

if [ -n "$2" ]; then
    GITHUB_TOKEN=$2
else
    echo 'Github token not set'
    exit 1
fi

zip "./rendr_darwin_amd64.zip" "./dist/darwin_amd64/rendr" -q -j
zip "./rendr_linux_amd64.zip" "./dist/linux_amd64/rendr" -q -j
zip "./rendr_windows_amd64.zip" "./dist/windows_amd64/rendr.exe" -q -j

RELEASE_NAME=v$VERSION

go install github.com/aktau/github-release@latest

GOPATH=$(go env GOPATH)

echo "Creating release in Github: $RELEASE_NAME"
set +e
$GOPATH/bin/github-release release --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --target main
set -e

sleep 5

echo "Releasing rendr_darwin_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --name rendr_darwin_amd64.zip  --file rendr_darwin_amd64.zip
echo "Releasing rendr_linux_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --name rendr_linux_amd64.zip   --file rendr_linux_amd64.zip
echo "Releasing rendr_windows_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --name rendr_windows_amd64.zip --file rendr_windows_amd64.zip

echo "Done releasing $RELEASE_NAME"
