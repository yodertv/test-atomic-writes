#!/bin/bash
# build.sh builds a go executable as a static build and tests it and the api.
# Redirect all output to the current log
set -x
export DIST=public
export BIN=bin
mkdir -p ${DIST}
mkdir -p ${BIN}
export LOG_NAME=${DIST}/build-output.log
echo "Building $LOG_NAME at"
date
exec 1>${LOG_NAME}
exec 2>&1
echo $PATH
env | grep URL
date
pwd

# Deploy static content
cp -p ./README.html $DIST/index.html
cp -p ./LICENSE.html $DIST/LICENSE.html

# Source target Go version from go.mod
export `sed -En 's/^go (.*)$/GO_VERSION=\1/p' go.mod`

echo ${GO_VERSION}
echo ${GOARCH}
echo ${GOOS}

# Check if 'go' exists in the PATH
if ! command -v go &> /dev/null; then
    echo "Go command not found. Initializing Go environment..."

    # 1. Download and extract portable Go into the current working directory
    curl -sL https://go.dev/dl/go${GO_VERSION}.${GOOS}-${GOARCH}.tar.gz | tar -xzf -

    # 2. Add the freshly extracted binary directory to the PATH
    export PATH="$PATH:$(pwd)/go/bin"
    echo "Go environment initialized successfully."
else
    echo "Go is already available in the system PATH."
fi

# Print the active version to confirm everything works
go version
go mod tidy

# Would like to be able to install this on my lamda, but haven't learned how to deploy more code to the api dectory.
# Now attempting to control the test via a monolith go host.

go build -o ./${BIN}/test-atomic-writes ./cmd/test-atomic-writes

# Truncate and date the new test-output.log file.
date > ${DIST}/test-output.log
go test -cpu 4 -parallel 20 -timeout 5m -v ./... >> ${DIST}/test-output.log

# List the filesystem for curiosity.
ls -lR .

# Build main last because it embeds the build and test output as static data to serve.
go build main.go

# Always exit successfully because when the build or test fails the dist directory is not served 
# by Vercel so you can't see the test output for debugging.
exit 0
