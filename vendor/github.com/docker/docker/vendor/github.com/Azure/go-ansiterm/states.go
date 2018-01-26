package ansiterm

type stateID int

type state interface ***REMOVED***
	Enter() error
	Exit() error
	Handle(byte) (state, error)
	Name() string
	Transition(state) error
***REMOVED***

type baseState struct ***REMOVED***
	name   string
	parser *AnsiParser
***REMOVED***

func (base baseState) Enter() error ***REMOVED***
	return nil
***REMOVED***

func (base baseState) Exit() error ***REMOVED***
	return nil
***REMOVED***

func (base baseState) Handle(b byte) (s state, e error) ***REMOVED***

	switch ***REMOVED***
	case b == CSI_ENTRY:
		return base.parser.csiEntry, nil
	case b == DCS_ENTRY:
		return base.parser.dcsEntry, nil
	case b == ANSI_ESCAPE_PRIMARY:
		return base.parser.escape, nil
	case b == OSC_STRING:
		return base.parser.oscString, nil
	case sliceContains(toGroundBytes, b):
		return base.parser.ground, nil
	***REMOVED***

	return nil, nil
***REMOVED***

func (base baseState) Name() string ***REMOVED***
	return base.name
***REMOVED***

func (base baseState) Transition(s state) error ***REMOVED***
	if s == base.parser.ground ***REMOVED***
		execBytes := []byte***REMOVED***0x18***REMOVED***
		execBytes = append(execBytes, 0x1A)
		execBytes = append(execBytes, getByteRange(0x80, 0x8F)...)
		execBytes = append(execBytes, getByteRange(0x91, 0x97)...)
		execBytes = append(execBytes, 0x99)
		execBytes = append(execBytes, 0x9A)

		if sliceContains(execBytes, base.parser.context.currentChar) ***REMOVED***
			return base.parser.execute()
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type dcsEntryState struct ***REMOVED***
	baseState
***REMOVED***

type errorState struct ***REMOVED***
	baseState
***REMOVED***
