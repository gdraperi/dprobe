#!/bin/sh

set -e

[ $(id -u) -eq 0 ] || ***REMOVED***
	printf >&2 '%s requires root\n' "$0"
	exit 1
***REMOVED***

usage() ***REMOVED***
	printf >&2 '%s: [-r release] [-m mirror] [-s] [-c additional repository] [-a arch]\n' "$0"
	exit 1
***REMOVED***

tmp() ***REMOVED***
	TMP=$(mktemp -d $***REMOVED***TMPDIR:-/var/tmp***REMOVED***/alpine-docker-XXXXXXXXXX)
	ROOTFS=$(mktemp -d $***REMOVED***TMPDIR:-/var/tmp***REMOVED***/alpine-docker-rootfs-XXXXXXXXXX)
	trap "rm -rf $TMP $ROOTFS" EXIT TERM INT
***REMOVED***

apkv() ***REMOVED***
	curl -sSL $MAINREPO/$ARCH/APKINDEX.tar.gz | tar -Oxz |
		grep --text '^P:apk-tools-static$' -A1 | tail -n1 | cut -d: -f2
***REMOVED***

getapk() ***REMOVED***
	curl -sSL $MAINREPO/$ARCH/apk-tools-static-$(apkv).apk |
		tar -xz -C $TMP sbin/apk.static
***REMOVED***

mkbase() ***REMOVED***
	$TMP/sbin/apk.static --repository $MAINREPO --update-cache --allow-untrusted \
		--root $ROOTFS --initdb add alpine-base
***REMOVED***

conf() ***REMOVED***
	printf '%s\n' $MAINREPO > $ROOTFS/etc/apk/repositories
	printf '%s\n' $ADDITIONALREPO >> $ROOTFS/etc/apk/repositories
***REMOVED***

pack() ***REMOVED***
	local id
	id=$(tar --numeric-owner -C $ROOTFS -c . | docker import - alpine:$REL)

	docker tag $id alpine:latest
	docker run -i -t --rm alpine printf 'alpine:%s with id=%s created!\n' $REL $id
***REMOVED***

save() ***REMOVED***
	[ $SAVE -eq 1 ] || return 0

	tar --numeric-owner -C $ROOTFS -c . | xz > rootfs.tar.xz
***REMOVED***

while getopts "hr:m:sc:a:" opt; do
	case $opt in
		r)
			REL=$OPTARG
			;;
		m)
			MIRROR=$OPTARG
			;;
		s)
			SAVE=1
			;;
		c)
			ADDITIONALREPO=$OPTARG
			;;
		a)
			ARCH=$OPTARG
			;;
		*)
			usage
			;;
	esac
done

REL=$***REMOVED***REL:-edge***REMOVED***
MIRROR=$***REMOVED***MIRROR:-http://nl.alpinelinux.org/alpine***REMOVED***
SAVE=$***REMOVED***SAVE:-0***REMOVED***
MAINREPO=$MIRROR/$REL/main
ADDITIONALREPO=$MIRROR/$REL/$***REMOVED***ADDITIONALREPO:-community***REMOVED***
ARCH=$***REMOVED***ARCH:-$(uname -m)***REMOVED***

tmp
getapk
mkbase
conf
pack
save
