#!/bin/bash
# build.sh builds a go executable as a static build and tests it and the api.
# Redirect all out put to the current log
export DIST=public
mkdir -p $DIST
export LOG_NAME=$DIST/build-output.log
exec 1>>$LOG_NAME
exec 2>&1
echo "PWD = "$PWD
echo "cp -p README.html $DIST/index.html"
cp -p README.html $DIST/index.html
echo "go version = "
go version
echo "go mod init main"
go mod init main
echo "go mod tidy"
go mod tidy
echo "go build -o test-atomic-writes"
go build -o test-atomic-writes
echo "go test main_test.go -cpu 4 -parallel 20 -timeout 5m -v > $DIST/test-output.log"
go test main_test.go -cpu 4 -parallel 20 -timeout 5m -v > $DIST/test-output.log

# Would like to be able to install this on my lamda, but haven't learned how to deploy more code to the api dectory.
# For now you have to redeploy to test again.
# mv test-atomic-writes api
echo "ls -lR"
ls -lR .

# Always exit successfully because when the build or test fails the dist directory is not served 
# by Vercel so you can't see the test output for debugging.
exit 0
