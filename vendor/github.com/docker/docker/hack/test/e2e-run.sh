#!/usr/bin/env bash
set -e

TESTFLAGS=$***REMOVED***TESTFLAGS:-""***REMOVED***
# Currently only DockerSuite and DockerNetworkSuite have been adapted for E2E testing
TESTFLAGS_LEGACY=$***REMOVED***TESTFLAGS_LEGACY:-""***REMOVED***
TIMEOUT=$***REMOVED***TIMEOUT:-60m***REMOVED***

SCRIPTDIR="$(dirname $***REMOVED***BASH_SOURCE[0]***REMOVED***)"

export DOCKER_ENGINE_GOARCH=$***REMOVED***DOCKER_ENGINE_GOARCH:-amd64***REMOVED***

run_test_integration() ***REMOVED***
  run_test_integration_suites
  run_test_integration_legacy_suites
***REMOVED***

run_test_integration_suites() ***REMOVED***
  local flags="-test.timeout=$***REMOVED***TIMEOUT***REMOVED*** $TESTFLAGS"
  for dir in /tests/integration/*; do
    if ! (
      cd $dir
      echo "Running $PWD"
      ./test.main $flags
    ); then exit 1; fi
  done
***REMOVED***

run_test_integration_legacy_suites() ***REMOVED***
  (
    flags="-check.timeout=$***REMOVED***TIMEOUT***REMOVED*** -test.timeout=360m $TESTFLAGS_LEGACY"
    cd /tests/integration-cli
    echo "Running $PWD"
    ./test.main $flags
  )
***REMOVED***

bash $SCRIPTDIR/ensure-emptyfs.sh

echo "Run integration tests"
run_test_integration
