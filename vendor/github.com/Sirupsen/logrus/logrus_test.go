package logrus

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LogAndAssertJSON(t *testing.T, log func(*Logger), assertions func(fields Fields)) ***REMOVED***
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	log(logger)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assertions(fields)
***REMOVED***

func LogAndAssertText(t *testing.T, log func(*Logger), assertions func(fields map[string]string)) ***REMOVED***
	var buffer bytes.Buffer

	logger := New()
	logger.Out = &buffer
	logger.Formatter = &TextFormatter***REMOVED***
		DisableColors: true,
	***REMOVED***

	log(logger)

	fields := make(map[string]string)
	for _, kv := range strings.Split(buffer.String(), " ") ***REMOVED***
		if !strings.Contains(kv, "=") ***REMOVED***
			continue
		***REMOVED***
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' ***REMOVED***
			var err error
			val, err = strconv.Unquote(val)
			assert.NoError(t, err)
		***REMOVED***
		fields[key] = val
	***REMOVED***
	assertions(fields)
***REMOVED***

func TestPrint(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Print("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	***REMOVED***)
***REMOVED***

func TestInfo(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Info("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "info")
	***REMOVED***)
***REMOVED***

func TestWarn(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Warn("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["level"], "warning")
	***REMOVED***)
***REMOVED***

func TestInfolnShouldAddSpacesBetweenStrings(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Infoln("test", "test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test test")
	***REMOVED***)
***REMOVED***

func TestInfolnShouldAddSpacesBetweenStringAndNonstring(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Infoln("test", 10)
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test 10")
	***REMOVED***)
***REMOVED***

func TestInfolnShouldAddSpacesBetweenTwoNonStrings(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Infoln(10, 10)
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "10 10")
	***REMOVED***)
***REMOVED***

func TestInfoShouldAddSpacesBetweenTwoNonStrings(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Infoln(10, 10)
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "10 10")
	***REMOVED***)
***REMOVED***

func TestInfoShouldNotAddSpacesBetweenStringAndNonstring(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Info("test", 10)
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test10")
	***REMOVED***)
***REMOVED***

func TestInfoShouldNotAddSpacesBetweenStrings(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.Info("test", "test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "testtest")
	***REMOVED***)
***REMOVED***

func TestWithFieldsShouldAllowAssignments(t *testing.T) ***REMOVED***
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	localLog := logger.WithFields(Fields***REMOVED***
		"key1": "value1",
	***REMOVED***)

	localLog.WithField("key2", "value2").Info("test")
	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assert.Equal(t, "value2", fields["key2"])
	assert.Equal(t, "value1", fields["key1"])

	buffer = bytes.Buffer***REMOVED******REMOVED***
	fields = Fields***REMOVED******REMOVED***
	localLog.Info("test")
	err = json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	_, ok := fields["key2"]
	assert.Equal(t, false, ok)
	assert.Equal(t, "value1", fields["key1"])
***REMOVED***

func TestUserSuppliedFieldDoesNotOverwriteDefaults(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.WithField("msg", "hello").Info("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test")
	***REMOVED***)
***REMOVED***

func TestUserSuppliedMsgFieldHasPrefix(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.WithField("msg", "hello").Info("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["msg"], "test")
		assert.Equal(t, fields["fields.msg"], "hello")
	***REMOVED***)
***REMOVED***

func TestUserSuppliedTimeFieldHasPrefix(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.WithField("time", "hello").Info("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["fields.time"], "hello")
	***REMOVED***)
***REMOVED***

func TestUserSuppliedLevelFieldHasPrefix(t *testing.T) ***REMOVED***
	LogAndAssertJSON(t, func(log *Logger) ***REMOVED***
		log.WithField("level", 1).Info("test")
	***REMOVED***, func(fields Fields) ***REMOVED***
		assert.Equal(t, fields["level"], "info")
		assert.Equal(t, fields["fields.level"], 1.0) // JSON has floats only
	***REMOVED***)
***REMOVED***

func TestDefaultFieldsAreNotPrefixed(t *testing.T) ***REMOVED***
	LogAndAssertText(t, func(log *Logger) ***REMOVED***
		ll := log.WithField("herp", "derp")
		ll.Info("hello")
		ll.Info("bye")
	***REMOVED***, func(fields map[string]string) ***REMOVED***
		for _, fieldName := range []string***REMOVED***"fields.level", "fields.time", "fields.msg"***REMOVED*** ***REMOVED***
			if _, ok := fields[fieldName]; ok ***REMOVED***
				t.Fatalf("should not have prefixed %q: %v", fieldName, fields)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestDoubleLoggingDoesntPrefixPreviousFields(t *testing.T) ***REMOVED***

	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	llog := logger.WithField("context", "eating raw fish")

	llog.Info("looks delicious")

	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.NoError(t, err, "should have decoded first message")
	assert.Equal(t, len(fields), 4, "should only have msg/time/level/context fields")
	assert.Equal(t, fields["msg"], "looks delicious")
	assert.Equal(t, fields["context"], "eating raw fish")

	buffer.Reset()

	llog.Warn("omg it is!")

	err = json.Unmarshal(buffer.Bytes(), &fields)
	assert.NoError(t, err, "should have decoded second message")
	assert.Equal(t, len(fields), 4, "should only have msg/time/level/context fields")
	assert.Equal(t, fields["msg"], "omg it is!")
	assert.Equal(t, fields["context"], "eating raw fish")
	assert.Nil(t, fields["fields.msg"], "should not have prefixed previous `msg` entry")

***REMOVED***

func TestConvertLevelToString(t *testing.T) ***REMOVED***
	assert.Equal(t, "debug", DebugLevel.String())
	assert.Equal(t, "info", InfoLevel.String())
	assert.Equal(t, "warning", WarnLevel.String())
	assert.Equal(t, "error", ErrorLevel.String())
	assert.Equal(t, "fatal", FatalLevel.String())
	assert.Equal(t, "panic", PanicLevel.String())
***REMOVED***

func TestParseLevel(t *testing.T) ***REMOVED***
	l, err := ParseLevel("panic")
	assert.Nil(t, err)
	assert.Equal(t, PanicLevel, l)

	l, err = ParseLevel("PANIC")
	assert.Nil(t, err)
	assert.Equal(t, PanicLevel, l)

	l, err = ParseLevel("fatal")
	assert.Nil(t, err)
	assert.Equal(t, FatalLevel, l)

	l, err = ParseLevel("FATAL")
	assert.Nil(t, err)
	assert.Equal(t, FatalLevel, l)

	l, err = ParseLevel("error")
	assert.Nil(t, err)
	assert.Equal(t, ErrorLevel, l)

	l, err = ParseLevel("ERROR")
	assert.Nil(t, err)
	assert.Equal(t, ErrorLevel, l)

	l, err = ParseLevel("warn")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("WARN")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("warning")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("WARNING")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("info")
	assert.Nil(t, err)
	assert.Equal(t, InfoLevel, l)

	l, err = ParseLevel("INFO")
	assert.Nil(t, err)
	assert.Equal(t, InfoLevel, l)

	l, err = ParseLevel("debug")
	assert.Nil(t, err)
	assert.Equal(t, DebugLevel, l)

	l, err = ParseLevel("DEBUG")
	assert.Nil(t, err)
	assert.Equal(t, DebugLevel, l)

	l, err = ParseLevel("invalid")
	assert.Equal(t, "not a valid logrus Level: \"invalid\"", err.Error())
***REMOVED***

func TestGetSetLevelRace(t *testing.T) ***REMOVED***
	wg := sync.WaitGroup***REMOVED******REMOVED***
	for i := 0; i < 100; i++ ***REMOVED***
		wg.Add(1)
		go func(i int) ***REMOVED***
			defer wg.Done()
			if i%2 == 0 ***REMOVED***
				SetLevel(InfoLevel)
			***REMOVED*** else ***REMOVED***
				GetLevel()
			***REMOVED***
		***REMOVED***(i)

	***REMOVED***
	wg.Wait()
***REMOVED***

func TestLoggingRace(t *testing.T) ***REMOVED***
	logger := New()

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ ***REMOVED***
		go func() ***REMOVED***
			logger.Info("info")
			wg.Done()
		***REMOVED***()
	***REMOVED***
	wg.Wait()
***REMOVED***

// Compile test
func TestLogrusInterface(t *testing.T) ***REMOVED***
	var buffer bytes.Buffer
	fn := func(l FieldLogger) ***REMOVED***
		b := l.WithField("key", "value")
		b.Debug("Test")
	***REMOVED***
	// test logger
	logger := New()
	logger.Out = &buffer
	fn(logger)

	// test Entry
	e := logger.WithField("another", "value")
	fn(e)
***REMOVED***

// Implements io.Writer using channels for synchronization, so we can wait on
// the Entry.Writer goroutine to write in a non-racey way. This does assume that
// there is a single call to Logger.Out for each message.
type channelWriter chan []byte

func (cw channelWriter) Write(p []byte) (int, error) ***REMOVED***
	cw <- p
	return len(p), nil
***REMOVED***

func TestEntryWriter(t *testing.T) ***REMOVED***
	cw := channelWriter(make(chan []byte, 1))
	log := New()
	log.Out = cw
	log.Formatter = new(JSONFormatter)
	log.WithField("foo", "bar").WriterLevel(WarnLevel).Write([]byte("hello\n"))

	bs := <-cw
	var fields Fields
	err := json.Unmarshal(bs, &fields)
	assert.Nil(t, err)
	assert.Equal(t, fields["foo"], "bar")
	assert.Equal(t, fields["level"], "warning")
***REMOVED***
