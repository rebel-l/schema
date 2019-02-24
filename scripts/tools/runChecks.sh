#!/usr/bin/env bash

COVER="-cover -coverprofile=./reports/coverage.txt -covermode=atomic"
FAST=""
RACE="-race"
SHORT=""

showHelp(){
    echo
    echo "Script executes linters and tests."
    echo "usage: ./scripts/tools/runChecks.sh [options]"
    echo
    echo "Options:"
    echo "-c        excludes coverage report"
    echo "-h, -?    shows this help"
    echo "-r        exclude checks for race conditions"
    echo "-s        executes only fast linters & tests"
    echo
}

while getopts "h?crs" opt; do
    case "$opt" in
        c)
            COVER=""
            ;;
        h|\?)
            showHelp
            exit 0
            ;;
        r)
            RACE=""
            ;;
        s)
            FAST="--fast"
            SHORT="-short"
            ;;
    esac
done

# Execute linter
echo
echo -en "\E[40;35m\033[1mExecute linters\033[0m"
echo
gometalinter ${FAST} --config=gometalinter.json
EXIT_CODE=$?
if [[ ${EXIT_CODE} != 0 ]]
then
	echo
	echo -en "\E[40;31m\033[1mLinting failed with exit code: \033[0m" ${EXIT_CODE}
	echo
	echo
	exit ${EXIT_CODE}
fi

# Execute tests
echo -en "\E[40;35m\033[1mExecute tests\033[0m"
echo
go test -v ${SHORT} ${RACE} ${COVER} ./...
EXIT_CODE=$?
if [[ ${EXIT_CODE} != 0 ]]
then
	echo
	echo -en "\E[40;31m\033[1mTests failed with exit code: \033[0m" ${EXIT_CODE}
	echo
	echo
	exit ${EXIT_CODE}
fi

# Success
echo
echo -en "\E[40;32m\033[1mChecks run successful :-)\033[0m"
echo
echo
