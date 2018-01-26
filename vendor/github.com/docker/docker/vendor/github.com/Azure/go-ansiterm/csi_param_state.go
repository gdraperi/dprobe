package ansiterm

type csiParamState struct ***REMOVED***
	baseState
***REMOVED***

func (csiState csiParamState) Handle(b byte) (s state, e error) ***REMOVED***
	csiState.parser.logf("CsiParam::Handle %#x", b)

	nextState, err := csiState.baseState.Handle(b)
	if nextState != nil || err != nil ***REMOVED***
		return nextState, err
	***REMOVED***

	switch ***REMOVED***
	case sliceContains(alphabetics, b):
		return csiState.parser.ground, nil
	case sliceContains(csiCollectables, b):
		csiState.parser.collectParam()
		return csiState, nil
	case sliceContains(executors, b):
		return csiState, csiState.parser.execute()
	***REMOVED***

	return csiState, nil
***REMOVED***

func (csiState csiParamState) Transition(s state) error ***REMOVED***
	csiState.parser.logf("CsiParam::Transition %s --> %s", csiState.Name(), s.Name())
	csiState.baseState.Transition(s)

	switch s ***REMOVED***
	case csiState.parser.ground:
		return csiState.parser.csiDispatch()
	***REMOVED***

	return nil
***REMOVED***
