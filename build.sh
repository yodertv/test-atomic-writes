#!/bin/bash
# build.sh builds a go executable as a static build and tests it and the api.
# Redirect all output to the current log
export DIST=public
mkdir -p ${DIST}
export LOG_NAME=${DIST}/build-output.log
exec 1>${LOG_NAME}
exec 2>&1
set -x
date
pwd
# cp -p ./README.html $DIST/index.html

# Define target Go version
GO_VERSION="1.25.3"

# Check if 'go' exists in the PATH
if ! command -v go &> /dev/null; then
    echo "⚠️ Go command not found. Initializing Go environment..."

    # 1. Download and extract portable Go into the current working directory
	curl -sL https://go.dev/dl/go${GO_VERSION}.darwin-amd64.pkg | tar -xzf -

    # 2. Add the freshly extracted binary directory to the PATH
    export PATH="$PATH:$(pwd)/go/bin"

    echo "✅ Go environment initialized successfully."
else
    echo "✅ Go is already available in the system PATH."
fi

# Print the active version to confirm everything works
go version

# go mod init main
# go mod tidy
go build -o test-atomic-writes
go test ./... -cpu 4 -parallel 20 -timeout 5m -v > ${DIST}/test-output.log

# Would like to be able to install this on my lamda, but haven't learned how to deploy more code to the api dectory.
# For now you have to redeploy to test again.
cp -p ./test-atomic-writes ${DIST}
ls -lR .

# Always exit successfully because when the build or test fails the dist directory is not served 
# by Vercel so you can't see the test output for debugging.
exit 0
