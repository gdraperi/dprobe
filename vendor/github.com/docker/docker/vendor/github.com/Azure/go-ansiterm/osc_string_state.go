package ansiterm

type oscStringState struct ***REMOVED***
	baseState
***REMOVED***

func (oscState oscStringState) Handle(b byte) (s state, e error) ***REMOVED***
	oscState.parser.logf("OscString::Handle %#x", b)
	nextState, err := oscState.baseState.Handle(b)
	if nextState != nil || err != nil ***REMOVED***
		return nextState, err
	***REMOVED***

	switch ***REMOVED***
	case isOscStringTerminator(b):
		return oscState.parser.ground, nil
	***REMOVED***

	return oscState, nil
***REMOVED***

// See below for OSC string terminators for linux
// http://man7.org/linux/man-pages/man4/console_codes.4.html
func isOscStringTerminator(b byte) bool ***REMOVED***

	if b == ANSI_BEL || b == 0x5C ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***
