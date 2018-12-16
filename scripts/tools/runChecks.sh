#!/usr/bin/env bash
# Execute linter
echo
echo -en "\E[40;35m\033[1mExecute linters\033[0m"
echo
gometalinter --config=gometalinter.json
EXITCODE=$?
if [ $EXITCODE != 0 ]
then
	echo
	echo -en "\E[40;31m\033[1mLinting failed with exit code: \033[0m" $EXITCODE
	echo
	echo
	exit $EXITCODE
fi

# Execute tests
echo -en "\E[40;35m\033[1mExecute tests\033[0m"
echo
go test -v github.com/rebel-l/schema github.com/rebel-l/schema/store
EXITCODE=$?
if [ $EXITCODE != 0 ]
then
	echo
	echo -en "\E[40;31m\033[1mTests failed with exit code: \033[0m" $EXITCODE
	echo
	echo
	exit $EXITCODE
fi

# Success
echo
echo -en "\E[40;32m\033[1mChecks run successful :-)\033[0m"
echo
echo
