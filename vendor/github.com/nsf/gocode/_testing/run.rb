#!/usr/bin/env ruby
# encoding: utf-8

RED = "\033[0;31m"
GRN = "\033[0;32m"
NC  = "\033[0m"

PASS = "#***REMOVED***GRN***REMOVED***PASS!#***REMOVED***NC***REMOVED***"
FAIL = "#***REMOVED***RED***REMOVED***FAIL!#***REMOVED***NC***REMOVED***"

Stats = Struct.new :total, :ok, :fail
$stats = Stats.new 0, 0, 0

def print_fail_report(t, out, outexpected)
	puts "#***REMOVED***t***REMOVED***: #***REMOVED***FAIL***REMOVED***"
	puts "-"*65
	puts "Got:\n#***REMOVED***out***REMOVED***"
	puts "-"*65
	puts "Expected:\n#***REMOVED***outexpected***REMOVED***"
	puts "-"*65
end

def print_pass_report(t)
	puts "#***REMOVED***t***REMOVED***: #***REMOVED***PASS***REMOVED***"
end

def print_stats
	puts "\nSummary (total: #***REMOVED***$stats.total***REMOVED***)"
	puts "#***REMOVED***GRN***REMOVED***  PASS#***REMOVED***NC***REMOVED***: #***REMOVED***$stats.ok***REMOVED***"
	puts "#***REMOVED***RED***REMOVED***  FAIL#***REMOVED***NC***REMOVED***: #***REMOVED***$stats.fail***REMOVED***"
	puts "#***REMOVED***$stats.fail == 0 ? GRN : RED***REMOVED***#***REMOVED***"█"*72***REMOVED***#***REMOVED***NC***REMOVED***"
end

def run_test(t)
	$stats.total += 1

	cursorpos = Dir["#***REMOVED***t***REMOVED***/cursor.*"].map***REMOVED***|d| File.extname(d)[1..-1]***REMOVED***.first
	outexpected = IO.read("#***REMOVED***t***REMOVED***/out.expected") rescue "To be determined"
	filename = "#***REMOVED***t***REMOVED***/test.go.in"

	out = %x[gocode -in #***REMOVED***filename***REMOVED*** autocomplete #***REMOVED***filename***REMOVED*** #***REMOVED***cursorpos***REMOVED***]

	if out != outexpected then
		print_fail_report(t, out, outexpected)
		$stats.fail += 1
	else
		print_pass_report(t)
		$stats.ok += 1
	end
end

if ARGV.one?
	run_test ARGV[0]
else
	Dir["test.*"].sort.each do |t| 
		run_test t
	end
end

print_stats
