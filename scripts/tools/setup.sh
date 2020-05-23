#!/usr/bin/env bash

function showHelp() {
    echo
    echo "This script installs the basic tools you need (go is required) to start development of this project."
    echo "It doesn't expect any parameter and installs gometalinters, goconvey and git hooks."
    echo
    echo "usage: ./scripts/tools/setup.sh [options]"
    echo
    echo "Options"
    echo "-h, -?    shows this help"
    echo
}

while getopts "h?" opt; do
    case "$opt" in
        h|\?)
            showHelp
            exit 0
            ;;
    esac
done

echo
echo -en "\E[40;34m\033[1mSetup local environment\033[0m"
echo
echo

# check if go is installed
go version
EXIT_CODE=$?
if [[ ${EXIT_CODE} != 0 ]]
then
    echo
    echo -en "\E[40;31m\033[1mGo is not installed! Please download and install Go from https://golang.org/dl/ before executing this script\033[0m"
    echo
    echo
    exit ${EXIT_CODE}
fi

# install golangci-lint
echo
echo -en "\E[40;34m\033[1mInstall: golangci-lint\033[0m"
echo
go get github.com/golangci/golangci-lint/cmd/golangci-lint

# install goconvey
echo
echo -en "\E[40;34m\033[1mInstall: goconvey\033[0m"
echo
go get github.com/smartystreets/goconvey

# install gomock
echo
echo -en "\E[40;34m\033[1mInstall: gomock\033[0m"
echo
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen

# tidy
echo
echo -en "\E[40;34m\033[1mTidy: go\033[0m"
echo
go mod tidy

# hooks
echo
echo -en "\E[40;34m\033[1mSetup: hooks\033[0m"
echo
cp ./scripts/hooks/* ./.git/hooks

echo
echo -en "\E[40;34m\033[1mSetup local environment done!\033[0m"
echo
