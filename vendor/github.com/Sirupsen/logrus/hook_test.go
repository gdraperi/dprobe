package logrus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestHook struct ***REMOVED***
	Fired bool
***REMOVED***

func (hook *TestHook) Fire(entry *Entry) error ***REMOVED***
	hook.Fired = true
	return nil
***REMOVED***

func (hook *TestHook) Levels() []Level ***REMOVED***
	return []Level***REMOVED***
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
	***REMOVED***
***REMOVED***

func TestHookFires(t *testing.T) ***REMOVED***
	hook := new(TestHook)

	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Hooks.Add(hook)
		assert.Equal(t, hook.Fired, false)

		log.Print("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, hook.Fired, true)
	***REMOVED***)
***REMOVED***

type ModifyHook struct ***REMOVED***
***REMOVED***

func (hook *ModifyHook) Fire(entry *Entry) error ***REMOVED***
	entry.Data["wow"] = "whale"
	return nil
***REMOVED***

func (hook *ModifyHook) Levels() []Level ***REMOVED***
	return []Level***REMOVED***
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
	***REMOVED***
***REMOVED***

func TestHookCanModifyEntry(t *testing.T) ***REMOVED***
	hook := new(ModifyHook)

	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Hooks.Add(hook)
		log.WithField("wow", "elephant").Print("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["wow"], "whale")
	***REMOVED***)
***REMOVED***

func TestCanFireMultipleHooks(t *testing.T) ***REMOVED***
	hook1 := new(ModifyHook)
	hook2 := new(TestHook)

	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Hooks.Add(hook1)
		log.Hooks.Add(hook2)

		log.WithField("wow", "elephant").Print("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["wow"], "whale")
		assert.Equal(t, hook2.Fired, true)
	***REMOVED***)
***REMOVED***

type ErrorHook struct ***REMOVED***
	Fired bool
***REMOVED***

func (hook *ErrorHook) Fire(entry *Entry) error ***REMOVED***
	hook.Fired = true
	return nil
***REMOVED***

func (hook *ErrorHook) Levels() []Level ***REMOVED***
	return []Level***REMOVED***
		ErrorLevel,
	***REMOVED***
***REMOVED***

func TestErrorHookShouldntFireOnInfo(t *testing.T) ***REMOVED***
	hook := new(ErrorHook)

	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Hooks.Add(hook)
		log.Info("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, hook.Fired, false)
	***REMOVED***)
***REMOVED***

func TestErrorHookShouldFireOnError(t *testing.T) ***REMOVED***
	hook := new(ErrorHook)

	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Hooks.Add(hook)
		log.Error("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, hook.Fired, true)
	***REMOVED***)
***REMOVED***

func TestAddHookRace(t *testing.T) ***REMOVED***
	var wg sync.WaitGroup
	wg.Add(2)
	hook := new(ErrorHook)
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		go func() ***REMOVED***
			defer wg.Done()
			log.AddHook(hook)
		***REMOVED***()
		go func() ***REMOVED***
			defer wg.Done()
			log.Error("test")
		***REMOVED***()
		wg.Wait()
	***REMOVED***, func(fields Fields) ***REMOVED***
		// the line may have been logged
		// before the hook was added, so we can't
		// actually assert on the hook
	***REMOVED***)
***REMOVED***
