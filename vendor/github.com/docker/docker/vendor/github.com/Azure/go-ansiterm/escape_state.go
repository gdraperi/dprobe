package ansiterm

type escapeState struct ***REMOVED***
	baseState
***REMOVED***

func (escState escapeState) Handle(b byte) (s state, e error) ***REMOVED***
	escState.parser.logf("escapeState::Handle %#x", b)
	nextState, err := escState.baseState.Handle(b)
	if nextState != nil || err != nil ***REMOVED***
		return nextState, err
	***REMOVED***

	switch ***REMOVED***
	case b == ANSI_ESCAPE_SECONDARY:
		return escState.parser.csiEntry, nil
	case b == ANSI_OSC_STRING_ENTRY:
		return escState.parser.oscString, nil
	case sliceContains(executors, b):
		return escState, escState.parser.execute()
	case sliceContains(escapeToGroundBytes, b):
		return escState.parser.ground, nil
	case sliceContains(intermeds, b):
		return escState.parser.escapeIntermediate, nil
	***REMOVED***

	return escState, nil
***REMOVED***

func (escState escapeState) Transition(s state) error ***REMOVED***
	escState.parser.logf("Escape::Transition %s --> %s", escState.Name(), s.Name())
	escState.baseState.Transition(s)

	switch s ***REMOVED***
	case escState.parser.ground:
		return escState.parser.escDispatch()
	case escState.parser.escapeIntermediate:
		return escState.parser.collectInter()
	***REMOVED***

	return nil
***REMOVED***

func (escState escapeState) Enter() error ***REMOVED***
	escState.parser.clear()
	return nil
***REMOVED***
