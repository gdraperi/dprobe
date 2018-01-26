#!/usr/bin/env bash

#
# This script generates a gitdm compatible email aliases file from a git
# formatted .mailmap file.
#
# Usage:
#  $> ./generate_aliases <mailmap_file> > aliases
#

cat $1 | \
    grep -v '^#' | \
    sed 's/^[^<]*<\([^>]*\)>/\1/' | \
    grep '<.*>' | sed -e 's/[<>]/ /g' | \
    awk '***REMOVED***if ($3 != "") ***REMOVED*** print $3" "$1 ***REMOVED*** else ***REMOVED***print $2" "$1***REMOVED******REMOVED***' | \
    sort | uniq
