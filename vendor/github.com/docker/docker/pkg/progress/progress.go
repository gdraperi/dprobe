package progress

import (
	"fmt"
)

// Progress represents the progress of a transfer.
type Progress struct ***REMOVED***
	ID string

	// Progress contains a Message or...
	Message string

	// ...progress of an action
	Action  string
	Current int64
	Total   int64

	// If true, don't show xB/yB
	HideCounts bool
	// If not empty, use units instead of bytes for counts
	Units string

	// Aux contains extra information not presented to the user, such as
	// digests for push signing.
	Aux interface***REMOVED******REMOVED***

	LastUpdate bool
***REMOVED***

// Output is an interface for writing progress information. It's
// like a writer for progress, but we don't call it Writer because
// that would be confusing next to ProgressReader (also, because it
// doesn't implement the io.Writer interface).
type Output interface ***REMOVED***
	WriteProgress(Progress) error
***REMOVED***

type chanOutput chan<- Progress

func (out chanOutput) WriteProgress(p Progress) error ***REMOVED***
	out <- p
	return nil
***REMOVED***

// ChanOutput returns an Output that writes progress updates to the
// supplied channel.
func ChanOutput(progressChan chan<- Progress) Output ***REMOVED***
	return chanOutput(progressChan)
***REMOVED***

type discardOutput struct***REMOVED******REMOVED***

func (discardOutput) WriteProgress(Progress) error ***REMOVED***
	return nil
***REMOVED***

// DiscardOutput returns an Output that discards progress
func DiscardOutput() Output ***REMOVED***
	return discardOutput***REMOVED******REMOVED***
***REMOVED***

// Update is a convenience function to write a progress update to the channel.
func Update(out Output, id, action string) ***REMOVED***
	out.WriteProgress(Progress***REMOVED***ID: id, Action: action***REMOVED***)
***REMOVED***

// Updatef is a convenience function to write a printf-formatted progress update
// to the channel.
func Updatef(out Output, id, format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	Update(out, id, fmt.Sprintf(format, a...))
***REMOVED***

// Message is a convenience function to write a progress message to the channel.
func Message(out Output, id, message string) ***REMOVED***
	out.WriteProgress(Progress***REMOVED***ID: id, Message: message***REMOVED***)
***REMOVED***

// Messagef is a convenience function to write a printf-formatted progress
// message to the channel.
func Messagef(out Output, id, format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	Message(out, id, fmt.Sprintf(format, a...))
***REMOVED***

// Aux sends auxiliary information over a progress interface, which will not be
// formatted for the UI. This is used for things such as push signing.
func Aux(out Output, a interface***REMOVED******REMOVED***) ***REMOVED***
	out.WriteProgress(Progress***REMOVED***Aux: a***REMOVED***)
***REMOVED***
