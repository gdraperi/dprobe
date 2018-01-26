// Package units provides helper function to parse and print size and time units
// in human-readable format.
package units

import (
	"fmt"
	"time"
)

// HumanDuration returns a human-readable approximation of a duration
// (eg. "About a minute", "4 hours ago", etc.).
func HumanDuration(d time.Duration) string ***REMOVED***
	if seconds := int(d.Seconds()); seconds < 1 ***REMOVED***
		return "Less than a second"
	***REMOVED*** else if seconds == 1 ***REMOVED***
		return "1 second"
	***REMOVED*** else if seconds < 60 ***REMOVED***
		return fmt.Sprintf("%d seconds", seconds)
	***REMOVED*** else if minutes := int(d.Minutes()); minutes == 1 ***REMOVED***
		return "About a minute"
	***REMOVED*** else if minutes < 46 ***REMOVED***
		return fmt.Sprintf("%d minutes", minutes)
	***REMOVED*** else if hours := int(d.Hours() + 0.5); hours == 1 ***REMOVED***
		return "About an hour"
	***REMOVED*** else if hours < 48 ***REMOVED***
		return fmt.Sprintf("%d hours", hours)
	***REMOVED*** else if hours < 24*7*2 ***REMOVED***
		return fmt.Sprintf("%d days", hours/24)
	***REMOVED*** else if hours < 24*30*2 ***REMOVED***
		return fmt.Sprintf("%d weeks", hours/24/7)
	***REMOVED*** else if hours < 24*365*2 ***REMOVED***
		return fmt.Sprintf("%d months", hours/24/30)
	***REMOVED***
	return fmt.Sprintf("%d years", int(d.Hours())/24/365)
***REMOVED***
