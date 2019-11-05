#!/bin/bash

set -e

CUR_DIR=$PWD
SRC_DIR=$PWD
cmd=$1

export GOPROXY=https://goproxy.io
export GOPATH=~/go
export GOBIN=~/go/bin

PKG_LIST=$(go list ./... | grep -v /vendor/)
LINT_VER=v0.0.0-20190409202823-959b441ac422

case $cmd in
    lint) $0 dep && $GOBIN/golint -set_exit_status ${PKG_LIST} ;;
    test) go test  -short ${PKG_LIST} ;;
    race) $0 dep && go test  -race -short ${PKG_LIST} ;;
    coverage) rm -rf coverage.* && go test ${PKG_LIST} -coverprofile=coverage.cov -covermode=count && go tool cover -func=./coverage.cov ;;
    dep) go get -v golang.org/x/lint@$LINT_VER && cd $GOPATH/pkg/mod/golang.org/x/lint@$LINT_VER/ && go install ./... && cd $CUR_DIR ;;
    build) $0 dep && go build ./... ;;
    clean) rm -rf ${PKG_LIST} && git checkout . ;;
esac

cd $CUR_DIR

