#!/bin/sh

# This is a convenience script for reporting issues that include a base
# template of information. See https://github.com/docker/docker/pull/8845

set -e

DOCKER_ISSUE_URL=$***REMOVED***DOCKER_ISSUE_URL:-"https://github.com/docker/docker/issues/new"***REMOVED***
DOCKER_ISSUE_NAME_PREFIX=$***REMOVED***DOCKER_ISSUE_NAME_PREFIX:-"Report: "***REMOVED***
DOCKER=$***REMOVED***DOCKER:-"docker"***REMOVED***
DOCKER_COMMAND="$***REMOVED***DOCKER***REMOVED***"
export DOCKER_COMMAND

# pulled from https://gist.github.com/cdown/1163649
function urlencode() ***REMOVED***
	# urlencode <string>

	local length="$***REMOVED***#1***REMOVED***"
	for (( i = 0; i < length; i++ )); do
			local c="$***REMOVED***1:i:1***REMOVED***"
			case $c in
					[a-zA-Z0-9.~_-]) printf "$c" ;;
					*) printf '%%%02X' "'$c"
			esac
	done
***REMOVED***

function template() ***REMOVED***
# this should always match the template from CONTRIBUTING.md
	cat <<- EOM
	Description of problem:


	\`docker version\`:
	`$***REMOVED***DOCKER_COMMAND***REMOVED*** -D version`


	\`docker info\`:
	`$***REMOVED***DOCKER_COMMAND***REMOVED*** -D info`


	\`uname -a\`:
	`uname -a`


	Environment details (AWS, VirtualBox, physical, etc.):


	How reproducible:


	Steps to Reproduce:
	1.
	2.
	3.


	Actual Results:


	Expected Results:


	Additional info:


	EOM
***REMOVED***

function format_issue_url() ***REMOVED***
	if [ $***REMOVED***#@***REMOVED*** -ne 2 ] ; then
		return 1
	fi
	local issue_name=$(urlencode "$***REMOVED***DOCKER_ISSUE_NAME_PREFIX***REMOVED***$***REMOVED***1***REMOVED***")
	local issue_body=$(urlencode "$***REMOVED***2***REMOVED***")
	echo "$***REMOVED***DOCKER_ISSUE_URL***REMOVED***?title=$***REMOVED***issue_name***REMOVED***&body=$***REMOVED***issue_body***REMOVED***"
***REMOVED***


echo -ne "Do you use \`sudo\` to call docker? [y|N]: "
read -r -n 1 use_sudo
echo ""

if [ "x$***REMOVED***use_sudo***REMOVED***" = "xy" -o "x$***REMOVED***use_sudo***REMOVED***" = "xY" ]; then
	export DOCKER_COMMAND="sudo $***REMOVED***DOCKER***REMOVED***"
fi

echo -ne "Title of new issue?: "
read -r issue_title
echo ""

issue_url=$(format_issue_url "$***REMOVED***issue_title***REMOVED***" "$(template)")

if which xdg-open 2>/dev/null >/dev/null ; then
	echo -ne "Would like to launch this report in your browser? [Y|n]: "
	read -r -n 1 launch_now
	echo ""

	if [ "$***REMOVED***launch_now***REMOVED***" != "n" -a "$***REMOVED***launch_now***REMOVED***" != "N" ]; then
		xdg-open "$***REMOVED***issue_url***REMOVED***"
	fi
fi

echo "If you would like to manually open the url, you can open this link if your browser: $***REMOVED***issue_url***REMOVED***"

