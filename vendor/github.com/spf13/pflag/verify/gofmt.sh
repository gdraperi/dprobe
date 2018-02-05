#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname "$***REMOVED***BASH_SOURCE***REMOVED***")/..

pushd "$***REMOVED***ROOT***REMOVED***" > /dev/null

GOFMT=$***REMOVED***GOFMT:-"gofmt"***REMOVED***
bad_files=$(find . -name '*.go' | xargs $GOFMT -s -l)
if [[ -n "$***REMOVED***bad_files***REMOVED***" ]]; then
  echo "!!! '$GOFMT' needs to be run on the following files: "
  echo "$***REMOVED***bad_files***REMOVED***"
  exit 1
fi

# ex: ts=2 sw=2 et filetype=sh
