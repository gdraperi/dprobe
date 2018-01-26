package ansiterm

type csiEntryState struct ***REMOVED***
	baseState
***REMOVED***

func (csiState csiEntryState) Handle(b byte) (s state, e error) ***REMOVED***
	csiState.parser.logf("CsiEntry::Handle %#x", b)

	nextState, err := csiState.baseState.Handle(b)
	if nextState != nil || err != nil ***REMOVED***
		return nextState, err
	***REMOVED***

	switch ***REMOVED***
	case sliceContains(alphabetics, b):
		return csiState.parser.ground, nil
	case sliceContains(csiCollectables, b):
		return csiState.parser.csiParam, nil
	case sliceContains(executors, b):
		return csiState, csiState.parser.execute()
	***REMOVED***

	return csiState, nil
***REMOVED***

func (csiState csiEntryState) Transition(s state) error ***REMOVED***
	csiState.parser.logf("CsiEntry::Transition %s --> %s", csiState.Name(), s.Name())
	csiState.baseState.Transition(s)

	switch s ***REMOVED***
	case csiState.parser.ground:
		return csiState.parser.csiDispatch()
	case csiState.parser.csiParam:
		switch ***REMOVED***
		case sliceContains(csiParams, csiState.parser.context.currentChar):
			csiState.parser.collectParam()
		case sliceContains(intermeds, csiState.parser.context.currentChar):
			csiState.parser.collectInter()
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (csiState csiEntryState) Enter() error ***REMOVED***
	csiState.parser.clear()
	return nil
***REMOVED***
