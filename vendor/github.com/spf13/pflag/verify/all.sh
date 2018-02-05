#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname "$***REMOVED***BASH_SOURCE***REMOVED***")/..

# Some useful colors.
if [[ -z "$***REMOVED***color_start-***REMOVED***" ]]; then
  declare -r color_start="\033["
  declare -r color_red="$***REMOVED***color_start***REMOVED***0;31m"
  declare -r color_yellow="$***REMOVED***color_start***REMOVED***0;33m"
  declare -r color_green="$***REMOVED***color_start***REMOVED***0;32m"
  declare -r color_norm="$***REMOVED***color_start***REMOVED***0m"
fi

SILENT=true

function is-excluded ***REMOVED***
  for e in $EXCLUDE; do
    if [[ $1 -ef $***REMOVED***BASH_SOURCE***REMOVED*** ]]; then
      return
    fi
    if [[ $1 -ef "$ROOT/hack/$e" ]]; then
      return
    fi
  done
  return 1
***REMOVED***

while getopts ":v" opt; do
  case $opt in
    v)
      SILENT=false
      ;;
    \?)
      echo "Invalid flag: -$OPTARG" >&2
      exit 1
      ;;
  esac
done

if $SILENT ; then
  echo "Running in the silent mode, run with -v if you want to see script logs."
fi

EXCLUDE="all.sh"

ret=0
for t in `ls $ROOT/verify/*.sh`
do
  if is-excluded $t ; then
    echo "Skipping $t"
    continue
  fi
  if $SILENT ; then
    echo -e "Verifying $t"
    if bash "$t" &> /dev/null; then
      echo -e "$***REMOVED***color_green***REMOVED***SUCCESS$***REMOVED***color_norm***REMOVED***"
    else
      echo -e "$***REMOVED***color_red***REMOVED***FAILED$***REMOVED***color_norm***REMOVED***"
      ret=1
    fi
  else
    bash "$t" || ret=1
  fi
done
exit $ret
