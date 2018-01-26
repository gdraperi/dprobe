package syslog

import (
	"reflect"
	"testing"

	syslog "github.com/RackSec/srslog"
)

func functionMatches(expectedFun interface***REMOVED******REMOVED***, actualFun interface***REMOVED******REMOVED***) bool ***REMOVED***
	return reflect.ValueOf(expectedFun).Pointer() == reflect.ValueOf(actualFun).Pointer()
***REMOVED***

func TestParseLogFormat(t *testing.T) ***REMOVED***
	formatter, framer, err := parseLogFormat("rfc5424", "udp")
	if err != nil || !functionMatches(rfc5424formatterWithAppNameAsTag, formatter) ||
		!functionMatches(syslog.DefaultFramer, framer) ***REMOVED***
		t.Fatal("Failed to parse rfc5424 format", err, formatter, framer)
	***REMOVED***

	formatter, framer, err = parseLogFormat("rfc5424", "tcp+tls")
	if err != nil || !functionMatches(rfc5424formatterWithAppNameAsTag, formatter) ||
		!functionMatches(syslog.RFC5425MessageLengthFramer, framer) ***REMOVED***
		t.Fatal("Failed to parse rfc5424 format", err, formatter, framer)
	***REMOVED***

	formatter, framer, err = parseLogFormat("rfc5424micro", "udp")
	if err != nil || !functionMatches(rfc5424microformatterWithAppNameAsTag, formatter) ||
		!functionMatches(syslog.DefaultFramer, framer) ***REMOVED***
		t.Fatal("Failed to parse rfc5424 (microsecond) format", err, formatter, framer)
	***REMOVED***

	formatter, framer, err = parseLogFormat("rfc5424micro", "tcp+tls")
	if err != nil || !functionMatches(rfc5424microformatterWithAppNameAsTag, formatter) ||
		!functionMatches(syslog.RFC5425MessageLengthFramer, framer) ***REMOVED***
		t.Fatal("Failed to parse rfc5424 (microsecond) format", err, formatter, framer)
	***REMOVED***

	formatter, framer, err = parseLogFormat("rfc3164", "")
	if err != nil || !functionMatches(syslog.RFC3164Formatter, formatter) ||
		!functionMatches(syslog.DefaultFramer, framer) ***REMOVED***
		t.Fatal("Failed to parse rfc3164 format", err, formatter, framer)
	***REMOVED***

	formatter, framer, err = parseLogFormat("", "")
	if err != nil || !functionMatches(syslog.UnixFormatter, formatter) ||
		!functionMatches(syslog.DefaultFramer, framer) ***REMOVED***
		t.Fatal("Failed to parse empty format", err, formatter, framer)
	***REMOVED***

	formatter, framer, err = parseLogFormat("invalid", "")
	if err == nil ***REMOVED***
		t.Fatal("Failed to parse invalid format", err, formatter, framer)
	***REMOVED***
***REMOVED***

func TestValidateLogOptEmpty(t *testing.T) ***REMOVED***
	emptyConfig := make(map[string]string)
	if err := ValidateLogOpt(emptyConfig); err != nil ***REMOVED***
		t.Fatal("Failed to parse empty config", err)
	***REMOVED***
***REMOVED***
