package ansiterm

type escapeIntermediateState struct ***REMOVED***
	baseState
***REMOVED***

func (escState escapeIntermediateState) Handle(b byte) (s state, e error) ***REMOVED***
	escState.parser.logf("escapeIntermediateState::Handle %#x", b)
	nextState, err := escState.baseState.Handle(b)
	if nextState != nil || err != nil ***REMOVED***
		return nextState, err
	***REMOVED***

	switch ***REMOVED***
	case sliceContains(intermeds, b):
		return escState, escState.parser.collectInter()
	case sliceContains(executors, b):
		return escState, escState.parser.execute()
	case sliceContains(escapeIntermediateToGroundBytes, b):
		return escState.parser.ground, nil
	***REMOVED***

	return escState, nil
***REMOVED***

func (escState escapeIntermediateState) Transition(s state) error ***REMOVED***
	escState.parser.logf("escapeIntermediateState::Transition %s --> %s", escState.Name(), s.Name())
	escState.baseState.Transition(s)

	switch s ***REMOVED***
	case escState.parser.ground:
		return escState.parser.escDispatch()
	***REMOVED***

	return nil
***REMOVED***
