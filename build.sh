#!/bin/bash
# build.sh builds a go executable as a static build and tests it and the api.
# Redirect all out put to the current log
export DIST=public
mkdir -p $DIST
export LOG_NAME=$DIST/build-output.log
exec 1>>$LOG_NAME
exec 2>&1
yum --quiet upgrade
yum install --assumeyes --quiet wget
wget -q https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
tar -xzf go1.23.0.linux-amd64.tar.gz
mv go /usr/local
export GOROOT=/usr/local/go
export GOPATH=$PWD
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
echo $PWD
mv README.html $DIST
go version
env
go env
go mod init
go build -o test-atomic-writes
go test -cpu 4 -parallel 20 -timeout 5m -v > $DIST/test-output.log

# Would like to be able to install this on my lamda, but haven't learned how to deploy more code to the api dectory.
# For now you have to redeploy to test again.
# mv test-atomic-writes api
ls -lR /var/task > $DIST/ls.txt

# Always exit successfully because when the build or test fails the dist directory is not served 
# by Vercel so you can't see the test output for debugging.
exit 0
