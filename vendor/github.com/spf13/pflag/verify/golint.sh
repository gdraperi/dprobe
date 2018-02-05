#!/bin/bash

ROOT=$(dirname "$***REMOVED***BASH_SOURCE***REMOVED***")/..
GOLINT=$***REMOVED***GOLINT:-"golint"***REMOVED***

pushd "$***REMOVED***ROOT***REMOVED***" > /dev/null
  bad_files=$($GOLINT -min_confidence=0.9 ./...)
  if [[ -n "$***REMOVED***bad_files***REMOVED***" ]]; then
    echo "!!! '$GOLINT' problems: "
    echo "$***REMOVED***bad_files***REMOVED***"
    exit 1
  fi
popd > /dev/null

# ex: ts=2 sw=2 et filetype=sh
