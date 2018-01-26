package shakers

import (
	"fmt"
	"time"

	"github.com/go-check/check"
)

// Default format when parsing (in addition to RFC and default time formats..)
const shortForm = "2006-01-02"

// IsBefore checker verifies the specified value is before the specified time.
// It is exclusive.
//
//    c.Assert(myTime, IsBefore, theTime, check.Commentf("bouuuhhh"))
//
var IsBefore check.Checker = &isBeforeChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "IsBefore",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type isBeforeChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *isBeforeChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return isBefore(params[0], params[1])
***REMOVED***

func isBefore(value, t interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	tTime, ok := parseTime(t)
	if !ok ***REMOVED***
		return false, "expected must be a Time struct, or parseable."
	***REMOVED***
	valueTime, valueIsTime := parseTime(value)
	if valueIsTime ***REMOVED***
		return valueTime.Before(tTime), ""
	***REMOVED***
	return false, "obtained value is not a time.Time struct or parseable as a time."
***REMOVED***

// IsAfter checker verifies the specified value is before the specified time.
// It is exclusive.
//
//    c.Assert(myTime, IsAfter, theTime, check.Commentf("bouuuhhh"))
//
var IsAfter check.Checker = &isAfterChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "IsAfter",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type isAfterChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *isAfterChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return isAfter(params[0], params[1])
***REMOVED***

func isAfter(value, t interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	tTime, ok := parseTime(t)
	if !ok ***REMOVED***
		return false, "expected must be a Time struct, or parseable."
	***REMOVED***
	valueTime, valueIsTime := parseTime(value)
	if valueIsTime ***REMOVED***
		return valueTime.After(tTime), ""
	***REMOVED***
	return false, "obtained value is not a time.Time struct or parseable as a time."
***REMOVED***

// IsBetween checker verifies the specified time is between the specified start
// and end. It's exclusive so if the specified time is at the tip of the interval.
//
//    c.Assert(myTime, IsBetween, startTime, endTime, check.Commentf("bouuuhhh"))
//
var IsBetween check.Checker = &isBetweenChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "IsBetween",
		Params: []string***REMOVED***"obtained", "start", "end"***REMOVED***,
	***REMOVED***,
***REMOVED***

type isBetweenChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *isBetweenChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return isBetween(params[0], params[1], params[2])
***REMOVED***

func isBetween(value, start, end interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	startTime, ok := parseTime(start)
	if !ok ***REMOVED***
		return false, "start must be a Time struct, or parseable."
	***REMOVED***
	endTime, ok := parseTime(end)
	if !ok ***REMOVED***
		return false, "end must be a Time struct, or parseable."
	***REMOVED***
	valueTime, valueIsTime := parseTime(value)
	if valueIsTime ***REMOVED***
		return valueTime.After(startTime) && valueTime.Before(endTime), ""
	***REMOVED***
	return false, "obtained value is not a time.Time struct or parseable as a time."
***REMOVED***

// TimeEquals checker verifies the specified time is the equal to the expected
// time.
//
//    c.Assert(myTime, TimeEquals, expected, check.Commentf("bouhhh"))
//
// It's possible to ignore some part of the time (like hours, minutes, etc..) using
// the TimeIgnore checker with it.
//
//    c.Assert(myTime, TimeIgnore(TimeEquals, time.Hour), expected, check.Commentf("... bouh.."))
//
var TimeEquals check.Checker = &timeEqualsChecker***REMOVED***
	&check.CheckerInfo***REMOVED***
		Name:   "TimeEquals",
		Params: []string***REMOVED***"obtained", "expected"***REMOVED***,
	***REMOVED***,
***REMOVED***

type timeEqualsChecker struct ***REMOVED***
	*check.CheckerInfo
***REMOVED***

func (checker *timeEqualsChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	return timeEquals(params[0], params[1])
***REMOVED***

func timeEquals(obtained, expected interface***REMOVED******REMOVED***) (bool, string) ***REMOVED***
	expectedTime, ok := parseTime(expected)
	if !ok ***REMOVED***
		return false, "expected must be a Time struct, or parseable."
	***REMOVED***
	valueTime, valueIsTime := parseTime(obtained)
	if valueIsTime ***REMOVED***
		return valueTime.Equal(expectedTime), ""
	***REMOVED***
	return false, "obtained value is not a time.Time struct or parseable as a time."
***REMOVED***

// TimeIgnore checker will ignore some part of the time on the encapsulated checker.
//
//    c.Assert(myTime, TimeIgnore(IsBetween, time.Second), start, end)
//
// FIXME use interface***REMOVED******REMOVED*** for ignore (to enable "Month", ..
func TimeIgnore(checker check.Checker, ignore time.Duration) check.Checker ***REMOVED***
	return &timeIgnoreChecker***REMOVED***
		sub:    checker,
		ignore: ignore,
	***REMOVED***
***REMOVED***

type timeIgnoreChecker struct ***REMOVED***
	sub    check.Checker
	ignore time.Duration
***REMOVED***

func (checker *timeIgnoreChecker) Info() *check.CheckerInfo ***REMOVED***
	info := *checker.sub.Info()
	info.Name = fmt.Sprintf("TimeIgnore(%s, %v)", info.Name, checker.ignore)
	return &info
***REMOVED***

func (checker *timeIgnoreChecker) Check(params []interface***REMOVED******REMOVED***, names []string) (bool, string) ***REMOVED***
	// Naive implementation : all params are supposed to be date
	mParams := make([]interface***REMOVED******REMOVED***, len(params))
	for index, param := range params ***REMOVED***
		paramTime, ok := parseTime(param)
		if !ok ***REMOVED***
			return false, fmt.Sprintf("%s must be a Time struct, or parseable.", names[index])
		***REMOVED***
		year := paramTime.Year()
		month := paramTime.Month()
		day := paramTime.Day()
		hour := paramTime.Hour()
		min := paramTime.Minute()
		sec := paramTime.Second()
		nsec := paramTime.Nanosecond()
		location := paramTime.Location()
		switch checker.ignore ***REMOVED***
		case time.Hour:
			hour = 0
			fallthrough
		case time.Minute:
			min = 0
			fallthrough
		case time.Second:
			sec = 0
			fallthrough
		case time.Millisecond:
			fallthrough
		case time.Microsecond:
			fallthrough
		case time.Nanosecond:
			nsec = 0
		***REMOVED***
		mParams[index] = time.Date(year, month, day, hour, min, sec, nsec, location)
	***REMOVED***
	return checker.sub.Check(mParams, names)
***REMOVED***

func parseTime(datetime interface***REMOVED******REMOVED***) (time.Time, bool) ***REMOVED***
	switch datetime.(type) ***REMOVED***
	case time.Time:
		return datetime.(time.Time), true
	case string:
		return parseTimeAsString(datetime.(string))
	default:
		if datetimeWithStr, ok := datetime.(fmt.Stringer); ok ***REMOVED***
			return parseTimeAsString(datetimeWithStr.String())
		***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
***REMOVED***

func parseTimeAsString(timeAsStr string) (time.Time, bool) ***REMOVED***
	forms := []string***REMOVED***shortForm, time.RFC3339, time.RFC3339Nano, time.RFC822, time.RFC822Z***REMOVED***
	for _, form := range forms ***REMOVED***
		datetime, err := time.Parse(form, timeAsStr)
		if err == nil ***REMOVED***
			return datetime, true
		***REMOVED***
	***REMOVED***
	return time.Time***REMOVED******REMOVED***, false
***REMOVED***
