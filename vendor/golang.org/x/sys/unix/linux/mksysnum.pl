#!/usr/bin/env perl
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

use strict;

if($ENV***REMOVED***'GOARCH'***REMOVED*** eq "" || $ENV***REMOVED***'GOOS'***REMOVED*** eq "") ***REMOVED***
	print STDERR "GOARCH or GOOS not defined in environment\n";
	exit 1;
***REMOVED***

# Check that we are using the new build system if we should
if($ENV***REMOVED***'GOLANG_SYS_BUILD'***REMOVED*** ne "docker") ***REMOVED***
	print STDERR "In the new build system, mksysnum should not be called directly.\n";
	print STDERR "See README.md\n";
	exit 1;
***REMOVED***

my $command = "$0 ". join(' ', @ARGV);

print <<EOF;
// $command
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build $ENV***REMOVED***'GOARCH'***REMOVED***,$ENV***REMOVED***'GOOS'***REMOVED***

package unix

const(
EOF

my $offset = 0;

sub fmt ***REMOVED***
	my ($name, $num) = @_;
	if($num > 999)***REMOVED***
		# ignore deprecated syscalls that are no longer implemented
		# https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/include/uapi/asm-generic/unistd.h?id=refs/heads/master#n716
		return;
	***REMOVED***
	$name =~ y/a-z/A-Z/;
	$num = $num + $offset;
	print "	SYS_$name = $num;\n";
***REMOVED***

my $prev;
open(CC, "$ENV***REMOVED***'CC'***REMOVED*** -E -dD @ARGV |") || die "can't run $ENV***REMOVED***'CC'***REMOVED***";
while(<CC>)***REMOVED***
	if(/^#define __NR_Linux\s+([0-9]+)/)***REMOVED***
		# mips/mips64: extract offset
		$offset = $1;
	***REMOVED***
	elsif(/^#define __NR(\w*)_SYSCALL_BASE\s+([0-9]+)/)***REMOVED***
		# arm: extract offset
		$offset = $1;
	***REMOVED***
	elsif(/^#define __NR_syscalls\s+/) ***REMOVED***
		# ignore redefinitions of __NR_syscalls
	***REMOVED***
	elsif(/^#define __NR_(\w*)Linux_syscalls\s+/) ***REMOVED***
		# mips/mips64: ignore definitions about the number of syscalls
	***REMOVED***
	elsif(/^#define __NR_(\w+)\s+([0-9]+)/)***REMOVED***
		$prev = $2;
		fmt($1, $2);
	***REMOVED***
	elsif(/^#define __NR3264_(\w+)\s+([0-9]+)/)***REMOVED***
		$prev = $2;
		fmt($1, $2);
	***REMOVED***
	elsif(/^#define __NR_(\w+)\s+\(\w+\+\s*([0-9]+)\)/)***REMOVED***
		fmt($1, $prev+$2)
	***REMOVED***
	elsif(/^#define __NR_(\w+)\s+\(__NR_Linux \+ ([0-9]+)/)***REMOVED***
		fmt($1, $2);
	***REMOVED***
	elsif(/^#define __NR_(\w+)\s+\(__NR_SYSCALL_BASE \+ ([0-9]+)/)***REMOVED***
		fmt($1, $2);
	***REMOVED***
***REMOVED***

print <<EOF;
)
EOF