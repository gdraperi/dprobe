#!/usr/bin/env tclsh

set red "\033\[0;31m"
set grn "\033\[0;32m"
set nc  "\033\[0m"

set pass "$***REMOVED***grn***REMOVED***PASS!$***REMOVED***nc***REMOVED***"
set fail "$***REMOVED***red***REMOVED***FAIL!$***REMOVED***nc***REMOVED***"

set stats.total 0
set stats.ok 0
set stats.fail 0

proc print_fail_report ***REMOVED***t out expected***REMOVED*** ***REMOVED***
	global fail

	set hr [join [lrepeat 65 "-"] ""]
	puts "$***REMOVED***t***REMOVED***: $***REMOVED***fail***REMOVED***"
	puts $hr
	puts "Got:\n$***REMOVED***out***REMOVED***"
	puts $hr
	puts "Expected:\n$***REMOVED***expected***REMOVED***"
	puts $hr
***REMOVED***

proc print_pass_report ***REMOVED***t***REMOVED*** ***REMOVED***
	global pass

	puts "$***REMOVED***t***REMOVED***: $***REMOVED***pass***REMOVED***"
***REMOVED***

proc print_stats ***REMOVED******REMOVED*** ***REMOVED***
	global red grn nc stats.total stats.ok stats.fail

	set hr [join [lrepeat 72 "█"] ""]
	set hrcol [expr ***REMOVED***$***REMOVED***stats.fail***REMOVED*** ? $red : $grn***REMOVED***]
	puts "\nSummary (total: $***REMOVED***stats.total***REMOVED***)"
	puts "$***REMOVED***grn***REMOVED***  PASS$***REMOVED***nc***REMOVED***: $***REMOVED***stats.ok***REMOVED***"
	puts "$***REMOVED***red***REMOVED***  FAIL$***REMOVED***nc***REMOVED***: $***REMOVED***stats.fail***REMOVED***"
	puts "$***REMOVED***hrcol***REMOVED***$***REMOVED***hr***REMOVED***$***REMOVED***nc***REMOVED***"
***REMOVED***

proc read_file ***REMOVED***filename***REMOVED*** ***REMOVED***
	set f [open $filename r]
	set data [read $f]
	close $f
	return $data
***REMOVED***

proc run_test ***REMOVED***t***REMOVED*** ***REMOVED***
	global stats.total stats.ok stats.fail

	incr stats.total
	set cursorpos [string range [file extension [glob "$***REMOVED***t***REMOVED***/cursor.*"]] 1 end]
	set expected [read_file "$***REMOVED***t***REMOVED***/out.expected"]
	set filename "$***REMOVED***t***REMOVED***/test.go.in"

	set out [read_file "| gocode -in $***REMOVED***filename***REMOVED*** autocomplete $***REMOVED***filename***REMOVED*** $***REMOVED***cursorpos***REMOVED***"]
	if ***REMOVED***$out eq $expected***REMOVED*** ***REMOVED***
		print_pass_report $t
		incr stats.ok
	***REMOVED*** else ***REMOVED***
		print_fail_report $t $out $expected
		incr stats.fail
	***REMOVED***
***REMOVED***

if ***REMOVED***$argc == 1***REMOVED*** ***REMOVED***
	run_test $argv
***REMOVED*** else ***REMOVED***
	foreach t [lsort [glob test.*]] ***REMOVED***
		run_test $t
	***REMOVED***
***REMOVED***

print_stats


