#!/usr/bin/env bash

# TODO: add parameters for "short" (linter and test) and "cover & race" (test)

# Execute linter
echo
echo -en "\E[40;35m\033[1mExecute linters\033[0m"
echo
gometalinter --config=gometalinter.json
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
go test -v -cover -race -coverprofile=./reports/coverage.txt -covermode=atomic ./...
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
