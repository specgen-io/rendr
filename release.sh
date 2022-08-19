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
    echo 'GITHUB_TOKEN is not set'
    exit 1
fi

if [ -n "$3" ]; then
    JFROG_USER=$3
else
    echo 'JFROG_USER is not set'
    exit 1
fi

if [ -n "$4" ]; then
    JFROG_PASS=$4
else
    echo 'JFROG_PASS is not set'
    exit 1
fi


zip "./rendr_darwin_amd64.zip" "./dist/darwin_amd64/rendr" -q -j
zip "./rendr_darwin_arm64.zip" "./dist/darwin_arm64/rendr" -q -j
zip "./rendr_linux_amd64.zip" "./dist/linux_amd64/rendr" -q -j
zip "./rendr_windows_amd64.zip" "./dist/windows_amd64/rendr.exe" -q -j

RELEASE_NAME=v$VERSION

go install github.com/aktau/github-release@latest

GOPATH=$(go env GOPATH)

echo "Releasing to Github: $RELEASE_NAME"
set +e
$GOPATH/bin/github-release release --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --target main
set -e

sleep 10

echo "Releasing rendr_darwin_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --name rendr_darwin_amd64.zip  --file rendr_darwin_amd64.zip
echo "Releasing rendr_darwin_arm64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --name rendr_darwin_arm64.zip  --file rendr_darwin_arm64.zip
echo "Releasing rendr_linux_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --name rendr_linux_amd64.zip   --file rendr_linux_amd64.zip
echo "Releasing rendr_windows_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo rendr --tag $RELEASE_NAME --name rendr_windows_amd64.zip --file rendr_windows_amd64.zip

echo "Done releasing to Github"

ARTFACTORY_URL="https://specgen.jfrog.io/artifactory/binaries/rendr"

echo "Releasing to Artifactory: $ARTFACTORY_URL/latest"

curl -u$JFROG_USER:$JFROG_PASS -T rendr_darwin_amd64.zip "$ARTFACTORY_URL/latest/rendr_darwin_amd64.zip"
curl -u$JFROG_USER:$JFROG_PASS -T rendr_darwin_arm64.zip "$ARTFACTORY_URL/latest/rendr_darwin_arm64.zip"
curl -u$JFROG_USER:$JFROG_PASS -T rendr_linux_amd64.zip "$ARTFACTORY_URL/latest/rendr_linux_amd64.zip"
curl -u$JFROG_USER:$JFROG_PASS -T rendr_windows_amd64.zip "$ARTFACTORY_URL/latest/rendr_windows_amd64.zip"

echo "Releasing to Artifactory: $ARTFACTORY_URL/$RELEASE_NAME"

curl -u$JFROG_USER:$JFROG_PASS -T rendr_darwin_amd64.zip "$ARTFACTORY_URL/$RELEASE_NAME/rendr_darwin_amd64.zip"
curl -u$JFROG_USER:$JFROG_PASS -T rendr_darwin_arm64.zip "$ARTFACTORY_URL/$RELEASE_NAME/rendr_darwin_arm64.zip"
curl -u$JFROG_USER:$JFROG_PASS -T rendr_linux_amd64.zip "$ARTFACTORY_URL/$RELEASE_NAME/rendr_linux_amd64.zip"
curl -u$JFROG_USER:$JFROG_PASS -T rendr_windows_amd64.zip "$ARTFACTORY_URL/$RELEASE_NAME/rendr_windows_amd64.zip"

echo "Done releasing to Artifactory"
