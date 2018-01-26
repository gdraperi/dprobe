package ansiterm

import (
	"errors"
	"log"
	"os"
)

type AnsiParser struct ***REMOVED***
	currState          state
	eventHandler       AnsiEventHandler
	context            *ansiContext
	csiEntry           state
	csiParam           state
	dcsEntry           state
	escape             state
	escapeIntermediate state
	error              state
	ground             state
	oscString          state
	stateMap           []state

	logf func(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type Option func(*AnsiParser)

func WithLogf(f func(string, ...interface***REMOVED******REMOVED***)) Option ***REMOVED***
	return func(ap *AnsiParser) ***REMOVED***
		ap.logf = f
	***REMOVED***
***REMOVED***

func CreateParser(initialState string, evtHandler AnsiEventHandler, opts ...Option) *AnsiParser ***REMOVED***
	ap := &AnsiParser***REMOVED***
		eventHandler: evtHandler,
		context:      &ansiContext***REMOVED******REMOVED***,
	***REMOVED***
	for _, o := range opts ***REMOVED***
		o(ap)
	***REMOVED***

	if isDebugEnv := os.Getenv(LogEnv); isDebugEnv == "1" ***REMOVED***
		logFile, _ := os.Create("ansiParser.log")
		logger := log.New(logFile, "", log.LstdFlags)
		if ap.logf != nil ***REMOVED***
			l := ap.logf
			ap.logf = func(s string, v ...interface***REMOVED******REMOVED***) ***REMOVED***
				l(s, v...)
				logger.Printf(s, v...)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			ap.logf = logger.Printf
		***REMOVED***
	***REMOVED***

	if ap.logf == nil ***REMOVED***
		ap.logf = func(string, ...interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
	***REMOVED***

	ap.csiEntry = csiEntryState***REMOVED***baseState***REMOVED***name: "CsiEntry", parser: ap***REMOVED******REMOVED***
	ap.csiParam = csiParamState***REMOVED***baseState***REMOVED***name: "CsiParam", parser: ap***REMOVED******REMOVED***
	ap.dcsEntry = dcsEntryState***REMOVED***baseState***REMOVED***name: "DcsEntry", parser: ap***REMOVED******REMOVED***
	ap.escape = escapeState***REMOVED***baseState***REMOVED***name: "Escape", parser: ap***REMOVED******REMOVED***
	ap.escapeIntermediate = escapeIntermediateState***REMOVED***baseState***REMOVED***name: "EscapeIntermediate", parser: ap***REMOVED******REMOVED***
	ap.error = errorState***REMOVED***baseState***REMOVED***name: "Error", parser: ap***REMOVED******REMOVED***
	ap.ground = groundState***REMOVED***baseState***REMOVED***name: "Ground", parser: ap***REMOVED******REMOVED***
	ap.oscString = oscStringState***REMOVED***baseState***REMOVED***name: "OscString", parser: ap***REMOVED******REMOVED***

	ap.stateMap = []state***REMOVED***
		ap.csiEntry,
		ap.csiParam,
		ap.dcsEntry,
		ap.escape,
		ap.escapeIntermediate,
		ap.error,
		ap.ground,
		ap.oscString,
	***REMOVED***

	ap.currState = getState(initialState, ap.stateMap)

	ap.logf("CreateParser: parser %p", ap)
	return ap
***REMOVED***

func getState(name string, states []state) state ***REMOVED***
	for _, el := range states ***REMOVED***
		if el.Name() == name ***REMOVED***
			return el
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (ap *AnsiParser) Parse(bytes []byte) (int, error) ***REMOVED***
	for i, b := range bytes ***REMOVED***
		if err := ap.handle(b); err != nil ***REMOVED***
			return i, err
		***REMOVED***
	***REMOVED***

	return len(bytes), ap.eventHandler.Flush()
***REMOVED***

func (ap *AnsiParser) handle(b byte) error ***REMOVED***
	ap.context.currentChar = b
	newState, err := ap.currState.Handle(b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if newState == nil ***REMOVED***
		ap.logf("WARNING: newState is nil")
		return errors.New("New state of 'nil' is invalid.")
	***REMOVED***

	if newState != ap.currState ***REMOVED***
		if err := ap.changeState(newState); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (ap *AnsiParser) changeState(newState state) error ***REMOVED***
	ap.logf("ChangeState %s --> %s", ap.currState.Name(), newState.Name())

	// Exit old state
	if err := ap.currState.Exit(); err != nil ***REMOVED***
		ap.logf("Exit state '%s' failed with : '%v'", ap.currState.Name(), err)
		return err
	***REMOVED***

	// Perform transition action
	if err := ap.currState.Transition(newState); err != nil ***REMOVED***
		ap.logf("Transition from '%s' to '%s' failed with: '%v'", ap.currState.Name(), newState.Name, err)
		return err
	***REMOVED***

	// Enter new state
	if err := newState.Enter(); err != nil ***REMOVED***
		ap.logf("Enter state '%s' failed with: '%v'", newState.Name(), err)
		return err
	***REMOVED***

	ap.currState = newState
	return nil
***REMOVED***
