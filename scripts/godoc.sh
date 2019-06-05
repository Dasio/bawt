#! /usr/bin/env bash

mkdir -p /tmp/tmpgoroot/doc
rm -rf /tmp/tmpgopath/src/github.com/gopherworks/bawt
mkdir -p /tmp/tmpgopath/src/github.com/gopherworks/bawt
tar -c --exclude='.git' --exclude='tmp' . | tar -x -C /tmp/tmpgopath/src/github.com/gopherworks/bawt
echo -e "open http://localhost:6060/pkg/github.com/gopherworks/bawt\n"
GOROOT=/tmp/tmpgoroot/ GOPATH=/tmp/tmpgopath/ godoc -http=localhost:6060