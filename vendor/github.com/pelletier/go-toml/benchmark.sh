#!/bin/bash

set -e

reference_ref=$***REMOVED***1:-master***REMOVED***
reference_git=$***REMOVED***2:-.***REMOVED***

if ! `hash benchstat 2>/dev/null`; then
    echo "Installing benchstat"
    go get golang.org/x/perf/cmd/benchstat
    go install golang.org/x/perf/cmd/benchstat
fi

tempdir=`mktemp -d /tmp/go-toml-benchmark-XXXXXX`
ref_tempdir="$***REMOVED***tempdir***REMOVED***/ref"
ref_benchmark="$***REMOVED***ref_tempdir***REMOVED***/benchmark-`echo -n $***REMOVED***reference_ref***REMOVED***|tr -s '/' '-'`.txt"
local_benchmark="`pwd`/benchmark-local.txt"

echo "=== $***REMOVED***reference_ref***REMOVED*** ($***REMOVED***ref_tempdir***REMOVED***)"
git clone $***REMOVED***reference_git***REMOVED*** $***REMOVED***ref_tempdir***REMOVED*** >/dev/null 2>/dev/null
pushd $***REMOVED***ref_tempdir***REMOVED*** >/dev/null
git checkout $***REMOVED***reference_ref***REMOVED*** >/dev/null 2>/dev/null
go test -bench=. -benchmem | tee $***REMOVED***ref_benchmark***REMOVED***
popd >/dev/null

echo ""
echo "=== local"
go test -bench=. -benchmem  | tee $***REMOVED***local_benchmark***REMOVED***

echo ""
echo "=== diff"
benchstat -delta-test=none $***REMOVED***ref_benchmark***REMOVED*** $***REMOVED***local_benchmark***REMOVED***