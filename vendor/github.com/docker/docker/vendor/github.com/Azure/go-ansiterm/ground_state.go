package ansiterm

type groundState struct ***REMOVED***
	baseState
***REMOVED***

func (gs groundState) Handle(b byte) (s state, e error) ***REMOVED***
	gs.parser.context.currentChar = b

	nextState, err := gs.baseState.Handle(b)
	if nextState != nil || err != nil ***REMOVED***
		return nextState, err
	***REMOVED***

	switch ***REMOVED***
	case sliceContains(printables, b):
		return gs, gs.parser.print()

	case sliceContains(executors, b):
		return gs, gs.parser.execute()
	***REMOVED***

	return gs, nil
***REMOVED***
