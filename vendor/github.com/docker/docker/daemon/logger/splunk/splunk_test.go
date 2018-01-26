package splunk

import (
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/require"
)

// Validate options
func TestValidateLogOpt(t *testing.T) ***REMOVED***
	err := ValidateLogOpt(map[string]string***REMOVED***
		splunkURLKey:                  "http://127.0.0.1",
		splunkTokenKey:                "2160C7EF-2CE9-4307-A180-F852B99CF417",
		splunkSourceKey:               "mysource",
		splunkSourceTypeKey:           "mysourcetype",
		splunkIndexKey:                "myindex",
		splunkCAPathKey:               "/usr/cert.pem",
		splunkCANameKey:               "ca_name",
		splunkInsecureSkipVerifyKey:   "true",
		splunkFormatKey:               "json",
		splunkVerifyConnectionKey:     "true",
		splunkGzipCompressionKey:      "true",
		splunkGzipCompressionLevelKey: "1",
		envKey:      "a",
		envRegexKey: "^foo",
		labelsKey:   "b",
		tagKey:      "c",
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = ValidateLogOpt(map[string]string***REMOVED***
		"not-supported-option": "a",
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expecting error on unsupported options")
	***REMOVED***
***REMOVED***

// Driver require user to specify required options
func TestNewMissedConfig(t *testing.T) ***REMOVED***
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED******REMOVED***,
	***REMOVED***
	_, err := New(info)
	if err == nil ***REMOVED***
		t.Fatal("Logger driver should fail when no required parameters specified")
	***REMOVED***
***REMOVED***

// Driver require user to specify splunk-url
func TestNewMissedUrl(t *testing.T) ***REMOVED***
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkTokenKey: "4642492F-D8BD-47F1-A005-0C08AE4657DF",
		***REMOVED***,
	***REMOVED***
	_, err := New(info)
	if err.Error() != "splunk: splunk-url is expected" ***REMOVED***
		t.Fatal("Logger driver should fail when no required parameters specified")
	***REMOVED***
***REMOVED***

// Driver require user to specify splunk-token
func TestNewMissedToken(t *testing.T) ***REMOVED***
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey: "http://127.0.0.1:8088",
		***REMOVED***,
	***REMOVED***
	_, err := New(info)
	if err.Error() != "splunk: splunk-token is expected" ***REMOVED***
		t.Fatal("Logger driver should fail when no required parameters specified")
	***REMOVED***
***REMOVED***

// Test default settings
func TestDefault(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:   hec.URL(),
			splunkTokenKey: hec.token,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	hostname, err := info.Hostname()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if loggerDriver.Name() != driverName ***REMOVED***
		t.Fatal("Unexpected logger driver name")
	***REMOVED***

	if !hec.connectionVerified ***REMOVED***
		t.Fatal("By default connection should be verified")
	***REMOVED***

	splunkLoggerDriver, ok := loggerDriver.(*splunkLoggerInline)
	if !ok ***REMOVED***
		t.Fatal("Unexpected Splunk Logging Driver type")
	***REMOVED***

	if splunkLoggerDriver.url != hec.URL()+"/services/collector/event/1.0" ||
		splunkLoggerDriver.auth != "Splunk "+hec.token ||
		splunkLoggerDriver.nullMessage.Host != hostname ||
		splunkLoggerDriver.nullMessage.Source != "" ||
		splunkLoggerDriver.nullMessage.SourceType != "" ||
		splunkLoggerDriver.nullMessage.Index != "" ||
		splunkLoggerDriver.gzipCompression ||
		splunkLoggerDriver.postMessagesFrequency != defaultPostMessagesFrequency ||
		splunkLoggerDriver.postMessagesBatchSize != defaultPostMessagesBatchSize ||
		splunkLoggerDriver.bufferMaximum != defaultBufferMaximum ||
		cap(splunkLoggerDriver.stream) != defaultStreamChannelSize ***REMOVED***
		t.Fatal("Found not default values setup in Splunk Logging Driver.")
	***REMOVED***

	message1Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("***REMOVED***\"a\":\"b\"***REMOVED***"), Source: "stdout", Timestamp: message1Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	message2Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("notajson"), Source: "stdout", Timestamp: message2Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 2 ***REMOVED***
		t.Fatal("Expected two messages")
	***REMOVED***

	if *hec.gzipEnabled ***REMOVED***
		t.Fatal("Gzip should not be used")
	***REMOVED***

	message1 := hec.messages[0]
	if message1.Time != fmt.Sprintf("%f", float64(message1Time.UnixNano())/float64(time.Second)) ||
		message1.Host != hostname ||
		message1.Source != "" ||
		message1.SourceType != "" ||
		message1.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 1 %v", message1)
	***REMOVED***

	if event, err := message1.EventAsMap(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event["line"] != "***REMOVED***\"a\":\"b\"***REMOVED***" ||
			event["source"] != "stdout" ||
			event["tag"] != "containeriid" ||
			len(event) != 3 ***REMOVED***
			t.Fatalf("Unexpected event in message %v", event)
		***REMOVED***
	***REMOVED***

	message2 := hec.messages[1]
	if message2.Time != fmt.Sprintf("%f", float64(message2Time.UnixNano())/float64(time.Second)) ||
		message2.Host != hostname ||
		message2.Source != "" ||
		message2.SourceType != "" ||
		message2.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 1 %v", message2)
	***REMOVED***

	if event, err := message2.EventAsMap(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event["line"] != "notajson" ||
			event["source"] != "stdout" ||
			event["tag"] != "containeriid" ||
			len(event) != 3 ***REMOVED***
			t.Fatalf("Unexpected event in message %v", event)
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify inline format with a not default settings for most of options
func TestInlineFormatWithNonDefaultOptions(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:             hec.URL(),
			splunkTokenKey:           hec.token,
			splunkSourceKey:          "mysource",
			splunkSourceTypeKey:      "mysourcetype",
			splunkIndexKey:           "myindex",
			splunkFormatKey:          splunkFormatInline,
			splunkGzipCompressionKey: "true",
			tagKey:      "***REMOVED******REMOVED***.ImageName***REMOVED******REMOVED***/***REMOVED******REMOVED***.Name***REMOVED******REMOVED***",
			labelsKey:   "a",
			envRegexKey: "^foo",
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
		ContainerLabels: map[string]string***REMOVED***
			"a": "b",
		***REMOVED***,
		ContainerEnv: []string***REMOVED***"foo_finder=bar"***REMOVED***,
	***REMOVED***

	hostname, err := info.Hostname()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if !hec.connectionVerified ***REMOVED***
		t.Fatal("By default connection should be verified")
	***REMOVED***

	splunkLoggerDriver, ok := loggerDriver.(*splunkLoggerInline)
	if !ok ***REMOVED***
		t.Fatal("Unexpected Splunk Logging Driver type")
	***REMOVED***

	if splunkLoggerDriver.url != hec.URL()+"/services/collector/event/1.0" ||
		splunkLoggerDriver.auth != "Splunk "+hec.token ||
		splunkLoggerDriver.nullMessage.Host != hostname ||
		splunkLoggerDriver.nullMessage.Source != "mysource" ||
		splunkLoggerDriver.nullMessage.SourceType != "mysourcetype" ||
		splunkLoggerDriver.nullMessage.Index != "myindex" ||
		!splunkLoggerDriver.gzipCompression ||
		splunkLoggerDriver.gzipCompressionLevel != gzip.DefaultCompression ||
		splunkLoggerDriver.postMessagesFrequency != defaultPostMessagesFrequency ||
		splunkLoggerDriver.postMessagesBatchSize != defaultPostMessagesBatchSize ||
		splunkLoggerDriver.bufferMaximum != defaultBufferMaximum ||
		cap(splunkLoggerDriver.stream) != defaultStreamChannelSize ***REMOVED***
		t.Fatal("Values do not match configuration.")
	***REMOVED***

	messageTime := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("1"), Source: "stdout", Timestamp: messageTime***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 1 ***REMOVED***
		t.Fatal("Expected one message")
	***REMOVED***

	if !*hec.gzipEnabled ***REMOVED***
		t.Fatal("Gzip should be used")
	***REMOVED***

	message := hec.messages[0]
	if message.Time != fmt.Sprintf("%f", float64(messageTime.UnixNano())/float64(time.Second)) ||
		message.Host != hostname ||
		message.Source != "mysource" ||
		message.SourceType != "mysourcetype" ||
		message.Index != "myindex" ***REMOVED***
		t.Fatalf("Unexpected values of message %v", message)
	***REMOVED***

	if event, err := message.EventAsMap(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event["line"] != "1" ||
			event["source"] != "stdout" ||
			event["tag"] != "container_image_name/container_name" ||
			event["attrs"].(map[string]interface***REMOVED******REMOVED***)["a"] != "b" ||
			event["attrs"].(map[string]interface***REMOVED******REMOVED***)["foo_finder"] != "bar" ||
			len(event) != 4 ***REMOVED***
			t.Fatalf("Unexpected event in message %v", event)
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify JSON format
func TestJsonFormat(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:                  hec.URL(),
			splunkTokenKey:                hec.token,
			splunkFormatKey:               splunkFormatJSON,
			splunkGzipCompressionKey:      "true",
			splunkGzipCompressionLevelKey: "1",
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	hostname, err := info.Hostname()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if !hec.connectionVerified ***REMOVED***
		t.Fatal("By default connection should be verified")
	***REMOVED***

	splunkLoggerDriver, ok := loggerDriver.(*splunkLoggerJSON)
	if !ok ***REMOVED***
		t.Fatal("Unexpected Splunk Logging Driver type")
	***REMOVED***

	if splunkLoggerDriver.url != hec.URL()+"/services/collector/event/1.0" ||
		splunkLoggerDriver.auth != "Splunk "+hec.token ||
		splunkLoggerDriver.nullMessage.Host != hostname ||
		splunkLoggerDriver.nullMessage.Source != "" ||
		splunkLoggerDriver.nullMessage.SourceType != "" ||
		splunkLoggerDriver.nullMessage.Index != "" ||
		!splunkLoggerDriver.gzipCompression ||
		splunkLoggerDriver.gzipCompressionLevel != gzip.BestSpeed ||
		splunkLoggerDriver.postMessagesFrequency != defaultPostMessagesFrequency ||
		splunkLoggerDriver.postMessagesBatchSize != defaultPostMessagesBatchSize ||
		splunkLoggerDriver.bufferMaximum != defaultBufferMaximum ||
		cap(splunkLoggerDriver.stream) != defaultStreamChannelSize ***REMOVED***
		t.Fatal("Values do not match configuration.")
	***REMOVED***

	message1Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("***REMOVED***\"a\":\"b\"***REMOVED***"), Source: "stdout", Timestamp: message1Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	message2Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("notjson"), Source: "stdout", Timestamp: message2Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 2 ***REMOVED***
		t.Fatal("Expected two messages")
	***REMOVED***

	message1 := hec.messages[0]
	if message1.Time != fmt.Sprintf("%f", float64(message1Time.UnixNano())/float64(time.Second)) ||
		message1.Host != hostname ||
		message1.Source != "" ||
		message1.SourceType != "" ||
		message1.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 1 %v", message1)
	***REMOVED***

	if event, err := message1.EventAsMap(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event["line"].(map[string]interface***REMOVED******REMOVED***)["a"] != "b" ||
			event["source"] != "stdout" ||
			event["tag"] != "containeriid" ||
			len(event) != 3 ***REMOVED***
			t.Fatalf("Unexpected event in message 1 %v", event)
		***REMOVED***
	***REMOVED***

	message2 := hec.messages[1]
	if message2.Time != fmt.Sprintf("%f", float64(message2Time.UnixNano())/float64(time.Second)) ||
		message2.Host != hostname ||
		message2.Source != "" ||
		message2.SourceType != "" ||
		message2.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 2 %v", message2)
	***REMOVED***

	// If message cannot be parsed as JSON - it should be sent as a line
	if event, err := message2.EventAsMap(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event["line"] != "notjson" ||
			event["source"] != "stdout" ||
			event["tag"] != "containeriid" ||
			len(event) != 3 ***REMOVED***
			t.Fatalf("Unexpected event in message 2 %v", event)
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify raw format
func TestRawFormat(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:    hec.URL(),
			splunkTokenKey:  hec.token,
			splunkFormatKey: splunkFormatRaw,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	hostname, err := info.Hostname()
	require.NoError(t, err)

	loggerDriver, err := New(info)
	require.NoError(t, err)

	if !hec.connectionVerified ***REMOVED***
		t.Fatal("By default connection should be verified")
	***REMOVED***

	splunkLoggerDriver, ok := loggerDriver.(*splunkLoggerRaw)
	if !ok ***REMOVED***
		t.Fatal("Unexpected Splunk Logging Driver type")
	***REMOVED***

	if splunkLoggerDriver.url != hec.URL()+"/services/collector/event/1.0" ||
		splunkLoggerDriver.auth != "Splunk "+hec.token ||
		splunkLoggerDriver.nullMessage.Host != hostname ||
		splunkLoggerDriver.nullMessage.Source != "" ||
		splunkLoggerDriver.nullMessage.SourceType != "" ||
		splunkLoggerDriver.nullMessage.Index != "" ||
		splunkLoggerDriver.gzipCompression ||
		splunkLoggerDriver.postMessagesFrequency != defaultPostMessagesFrequency ||
		splunkLoggerDriver.postMessagesBatchSize != defaultPostMessagesBatchSize ||
		splunkLoggerDriver.bufferMaximum != defaultBufferMaximum ||
		cap(splunkLoggerDriver.stream) != defaultStreamChannelSize ||
		string(splunkLoggerDriver.prefix) != "containeriid " ***REMOVED***
		t.Fatal("Values do not match configuration.")
	***REMOVED***

	message1Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("***REMOVED***\"a\":\"b\"***REMOVED***"), Source: "stdout", Timestamp: message1Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	message2Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("notjson"), Source: "stdout", Timestamp: message2Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 2 ***REMOVED***
		t.Fatal("Expected two messages")
	***REMOVED***

	message1 := hec.messages[0]
	if message1.Time != fmt.Sprintf("%f", float64(message1Time.UnixNano())/float64(time.Second)) ||
		message1.Host != hostname ||
		message1.Source != "" ||
		message1.SourceType != "" ||
		message1.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 1 %v", message1)
	***REMOVED***

	if event, err := message1.EventAsString(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event != "containeriid ***REMOVED***\"a\":\"b\"***REMOVED***" ***REMOVED***
			t.Fatalf("Unexpected event in message 1 %v", event)
		***REMOVED***
	***REMOVED***

	message2 := hec.messages[1]
	if message2.Time != fmt.Sprintf("%f", float64(message2Time.UnixNano())/float64(time.Second)) ||
		message2.Host != hostname ||
		message2.Source != "" ||
		message2.SourceType != "" ||
		message2.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 2 %v", message2)
	***REMOVED***

	if event, err := message2.EventAsString(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event != "containeriid notjson" ***REMOVED***
			t.Fatalf("Unexpected event in message 1 %v", event)
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify raw format with labels
func TestRawFormatWithLabels(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:    hec.URL(),
			splunkTokenKey:  hec.token,
			splunkFormatKey: splunkFormatRaw,
			labelsKey:       "a",
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
		ContainerLabels: map[string]string***REMOVED***
			"a": "b",
		***REMOVED***,
	***REMOVED***

	hostname, err := info.Hostname()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if !hec.connectionVerified ***REMOVED***
		t.Fatal("By default connection should be verified")
	***REMOVED***

	splunkLoggerDriver, ok := loggerDriver.(*splunkLoggerRaw)
	if !ok ***REMOVED***
		t.Fatal("Unexpected Splunk Logging Driver type")
	***REMOVED***

	if splunkLoggerDriver.url != hec.URL()+"/services/collector/event/1.0" ||
		splunkLoggerDriver.auth != "Splunk "+hec.token ||
		splunkLoggerDriver.nullMessage.Host != hostname ||
		splunkLoggerDriver.nullMessage.Source != "" ||
		splunkLoggerDriver.nullMessage.SourceType != "" ||
		splunkLoggerDriver.nullMessage.Index != "" ||
		splunkLoggerDriver.gzipCompression ||
		splunkLoggerDriver.postMessagesFrequency != defaultPostMessagesFrequency ||
		splunkLoggerDriver.postMessagesBatchSize != defaultPostMessagesBatchSize ||
		splunkLoggerDriver.bufferMaximum != defaultBufferMaximum ||
		cap(splunkLoggerDriver.stream) != defaultStreamChannelSize ||
		string(splunkLoggerDriver.prefix) != "containeriid a=b " ***REMOVED***
		t.Fatal("Values do not match configuration.")
	***REMOVED***

	message1Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("***REMOVED***\"a\":\"b\"***REMOVED***"), Source: "stdout", Timestamp: message1Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	message2Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("notjson"), Source: "stdout", Timestamp: message2Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 2 ***REMOVED***
		t.Fatal("Expected two messages")
	***REMOVED***

	message1 := hec.messages[0]
	if message1.Time != fmt.Sprintf("%f", float64(message1Time.UnixNano())/float64(time.Second)) ||
		message1.Host != hostname ||
		message1.Source != "" ||
		message1.SourceType != "" ||
		message1.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 1 %v", message1)
	***REMOVED***

	if event, err := message1.EventAsString(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event != "containeriid a=b ***REMOVED***\"a\":\"b\"***REMOVED***" ***REMOVED***
			t.Fatalf("Unexpected event in message 1 %v", event)
		***REMOVED***
	***REMOVED***

	message2 := hec.messages[1]
	if message2.Time != fmt.Sprintf("%f", float64(message2Time.UnixNano())/float64(time.Second)) ||
		message2.Host != hostname ||
		message2.Source != "" ||
		message2.SourceType != "" ||
		message2.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 2 %v", message2)
	***REMOVED***

	if event, err := message2.EventAsString(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event != "containeriid a=b notjson" ***REMOVED***
			t.Fatalf("Unexpected event in message 2 %v", event)
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify that Splunk Logging Driver can accept tag="" which will allow to send raw messages
// in the same way we get them in stdout/stderr
func TestRawFormatWithoutTag(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:    hec.URL(),
			splunkTokenKey:  hec.token,
			splunkFormatKey: splunkFormatRaw,
			tagKey:          "",
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	hostname, err := info.Hostname()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if !hec.connectionVerified ***REMOVED***
		t.Fatal("By default connection should be verified")
	***REMOVED***

	splunkLoggerDriver, ok := loggerDriver.(*splunkLoggerRaw)
	if !ok ***REMOVED***
		t.Fatal("Unexpected Splunk Logging Driver type")
	***REMOVED***

	if splunkLoggerDriver.url != hec.URL()+"/services/collector/event/1.0" ||
		splunkLoggerDriver.auth != "Splunk "+hec.token ||
		splunkLoggerDriver.nullMessage.Host != hostname ||
		splunkLoggerDriver.nullMessage.Source != "" ||
		splunkLoggerDriver.nullMessage.SourceType != "" ||
		splunkLoggerDriver.nullMessage.Index != "" ||
		splunkLoggerDriver.gzipCompression ||
		splunkLoggerDriver.postMessagesFrequency != defaultPostMessagesFrequency ||
		splunkLoggerDriver.postMessagesBatchSize != defaultPostMessagesBatchSize ||
		splunkLoggerDriver.bufferMaximum != defaultBufferMaximum ||
		cap(splunkLoggerDriver.stream) != defaultStreamChannelSize ||
		string(splunkLoggerDriver.prefix) != "" ***REMOVED***
		t.Log(string(splunkLoggerDriver.prefix) + "a")
		t.Fatal("Values do not match configuration.")
	***REMOVED***

	message1Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("***REMOVED***\"a\":\"b\"***REMOVED***"), Source: "stdout", Timestamp: message1Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	message2Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("notjson"), Source: "stdout", Timestamp: message2Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	message3Time := time.Now()
	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(" "), Source: "stdout", Timestamp: message3Time***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// message3 would have an empty or whitespace only string in the "event" field
	// both of which are not acceptable to HEC
	// thus here we must expect 2 messages, not 3
	if len(hec.messages) != 2 ***REMOVED***
		t.Fatal("Expected two messages")
	***REMOVED***

	message1 := hec.messages[0]
	if message1.Time != fmt.Sprintf("%f", float64(message1Time.UnixNano())/float64(time.Second)) ||
		message1.Host != hostname ||
		message1.Source != "" ||
		message1.SourceType != "" ||
		message1.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 1 %v", message1)
	***REMOVED***

	if event, err := message1.EventAsString(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event != "***REMOVED***\"a\":\"b\"***REMOVED***" ***REMOVED***
			t.Fatalf("Unexpected event in message 1 %v", event)
		***REMOVED***
	***REMOVED***

	message2 := hec.messages[1]
	if message2.Time != fmt.Sprintf("%f", float64(message2Time.UnixNano())/float64(time.Second)) ||
		message2.Host != hostname ||
		message2.Source != "" ||
		message2.SourceType != "" ||
		message2.Index != "" ***REMOVED***
		t.Fatalf("Unexpected values of message 2 %v", message2)
	***REMOVED***

	if event, err := message2.EventAsString(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event != "notjson" ***REMOVED***
			t.Fatalf("Unexpected event in message 2 %v", event)
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify that we will send messages in batches with default batching parameters,
// but change frequency to be sure that numOfRequests will match expected 17 requests
func TestBatching(t *testing.T) ***REMOVED***
	if err := os.Setenv(envVarPostMessagesFrequency, "10h"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:   hec.URL(),
			splunkTokenKey: hec.token,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for i := 0; i < defaultStreamChannelSize*4; i++ ***REMOVED***
		if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(fmt.Sprintf("%d", i)), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != defaultStreamChannelSize*4 ***REMOVED***
		t.Fatal("Not all messages delivered")
	***REMOVED***

	for i, message := range hec.messages ***REMOVED***
		if event, err := message.EventAsMap(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			if event["line"] != fmt.Sprintf("%d", i) ***REMOVED***
				t.Fatalf("Unexpected event in message %v", event)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// 1 to verify connection and 16 batches
	if hec.numOfRequests != 17 ***REMOVED***
		t.Fatalf("Unexpected number of requests %d", hec.numOfRequests)
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarPostMessagesFrequency, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify that test is using time to fire events not rare than specified frequency
func TestFrequency(t *testing.T) ***REMOVED***
	if err := os.Setenv(envVarPostMessagesFrequency, "5ms"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:   hec.URL(),
			splunkTokenKey: hec.token,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for i := 0; i < 10; i++ ***REMOVED***
		if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(fmt.Sprintf("%d", i)), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		time.Sleep(15 * time.Millisecond)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 10 ***REMOVED***
		t.Fatal("Not all messages delivered")
	***REMOVED***

	for i, message := range hec.messages ***REMOVED***
		if event, err := message.EventAsMap(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			if event["line"] != fmt.Sprintf("%d", i) ***REMOVED***
				t.Fatalf("Unexpected event in message %v", event)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// 1 to verify connection and 10 to verify that we have sent messages with required frequency,
	// but because frequency is too small (to keep test quick), instead of 11, use 9 if context switches will be slow
	if hec.numOfRequests < 9 ***REMOVED***
		t.Fatalf("Unexpected number of requests %d", hec.numOfRequests)
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarPostMessagesFrequency, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Simulate behavior similar to first version of Splunk Logging Driver, when we were sending one message
// per request
func TestOneMessagePerRequest(t *testing.T) ***REMOVED***
	if err := os.Setenv(envVarPostMessagesFrequency, "10h"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarPostMessagesBatchSize, "1"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarBufferMaximum, "1"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarStreamChannelSize, "0"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	hec := NewHTTPEventCollectorMock(t)

	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:   hec.URL(),
			splunkTokenKey: hec.token,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for i := 0; i < 10; i++ ***REMOVED***
		if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(fmt.Sprintf("%d", i)), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 10 ***REMOVED***
		t.Fatal("Not all messages delivered")
	***REMOVED***

	for i, message := range hec.messages ***REMOVED***
		if event, err := message.EventAsMap(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			if event["line"] != fmt.Sprintf("%d", i) ***REMOVED***
				t.Fatalf("Unexpected event in message %v", event)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// 1 to verify connection and 10 messages
	if hec.numOfRequests != 11 ***REMOVED***
		t.Fatalf("Unexpected number of requests %d", hec.numOfRequests)
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarPostMessagesFrequency, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarPostMessagesBatchSize, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarBufferMaximum, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarStreamChannelSize, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Driver should not be created when HEC is unresponsive
func TestVerify(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)
	hec.simulateServerError = true
	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:   hec.URL(),
			splunkTokenKey: hec.token,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	_, err := New(info)
	if err == nil ***REMOVED***
		t.Fatal("Expecting driver to fail, when server is unresponsive")
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify that user can specify to skip verification that Splunk HEC is working.
// Also in this test we verify retry logic.
func TestSkipVerify(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)
	hec.simulateServerError = true
	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:              hec.URL(),
			splunkTokenKey:            hec.token,
			splunkVerifyConnectionKey: "false",
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if hec.connectionVerified ***REMOVED***
		t.Fatal("Connection should not be verified")
	***REMOVED***

	for i := 0; i < defaultStreamChannelSize*2; i++ ***REMOVED***
		if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(fmt.Sprintf("%d", i)), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	if len(hec.messages) != 0 ***REMOVED***
		t.Fatal("No messages should be accepted at this point")
	***REMOVED***

	hec.simulateErr(false)

	for i := defaultStreamChannelSize * 2; i < defaultStreamChannelSize*4; i++ ***REMOVED***
		if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(fmt.Sprintf("%d", i)), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != defaultStreamChannelSize*4 ***REMOVED***
		t.Fatal("Not all messages delivered")
	***REMOVED***

	for i, message := range hec.messages ***REMOVED***
		if event, err := message.EventAsMap(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			if event["line"] != fmt.Sprintf("%d", i) ***REMOVED***
				t.Fatalf("Unexpected event in message %v", event)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify logic for when we filled whole buffer
func TestBufferMaximum(t *testing.T) ***REMOVED***
	if err := os.Setenv(envVarPostMessagesBatchSize, "2"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarBufferMaximum, "10"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarStreamChannelSize, "0"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	hec := NewHTTPEventCollectorMock(t)
	hec.simulateErr(true)
	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:              hec.URL(),
			splunkTokenKey:            hec.token,
			splunkVerifyConnectionKey: "false",
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if hec.connectionVerified ***REMOVED***
		t.Fatal("Connection should not be verified")
	***REMOVED***

	for i := 0; i < 11; i++ ***REMOVED***
		if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(fmt.Sprintf("%d", i)), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	if len(hec.messages) != 0 ***REMOVED***
		t.Fatal("No messages should be accepted at this point")
	***REMOVED***

	hec.simulateServerError = false

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 9 ***REMOVED***
		t.Fatalf("Expected # of messages %d, got %d", 9, len(hec.messages))
	***REMOVED***

	// First 1000 messages are written to daemon log when buffer was full
	for i, message := range hec.messages ***REMOVED***
		if event, err := message.EventAsMap(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED*** else ***REMOVED***
			if event["line"] != fmt.Sprintf("%d", i+2) ***REMOVED***
				t.Fatalf("Unexpected event in message %v", event)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarPostMessagesBatchSize, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarBufferMaximum, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarStreamChannelSize, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Verify that we are not blocking close when HEC is down for the whole time
func TestServerAlwaysDown(t *testing.T) ***REMOVED***
	if err := os.Setenv(envVarPostMessagesBatchSize, "2"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarBufferMaximum, "4"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarStreamChannelSize, "0"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	hec := NewHTTPEventCollectorMock(t)
	hec.simulateServerError = true
	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:              hec.URL(),
			splunkTokenKey:            hec.token,
			splunkVerifyConnectionKey: "false",
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if hec.connectionVerified ***REMOVED***
		t.Fatal("Connection should not be verified")
	***REMOVED***

	for i := 0; i < 5; i++ ***REMOVED***
		if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte(fmt.Sprintf("%d", i)), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(hec.messages) != 0 ***REMOVED***
		t.Fatal("No messages should be sent")
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarPostMessagesBatchSize, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarBufferMaximum, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := os.Setenv(envVarStreamChannelSize, ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Cannot send messages after we close driver
func TestCannotSendAfterClose(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)
	go hec.Serve()

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:   hec.URL(),
			splunkTokenKey: hec.token,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	loggerDriver, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("message1"), Source: "stdout", Timestamp: time.Now()***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	err = loggerDriver.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := loggerDriver.Log(&logger.Message***REMOVED***Line: []byte("message2"), Source: "stdout", Timestamp: time.Now()***REMOVED***); err == nil ***REMOVED***
		t.Fatal("Driver should not allow to send messages after close")
	***REMOVED***

	if len(hec.messages) != 1 ***REMOVED***
		t.Fatal("Only one message should be sent")
	***REMOVED***

	message := hec.messages[0]
	if event, err := message.EventAsMap(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else ***REMOVED***
		if event["line"] != "message1" ***REMOVED***
			t.Fatalf("Unexpected event in message %v", event)
		***REMOVED***
	***REMOVED***

	err = hec.Close()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestDeadlockOnBlockedEndpoint(t *testing.T) ***REMOVED***
	hec := NewHTTPEventCollectorMock(t)
	go hec.Serve()
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			splunkURLKey:   hec.URL(),
			splunkTokenKey: hec.token,
		***REMOVED***,
		ContainerID:        "containeriid",
		ContainerName:      "/container_name",
		ContainerImageID:   "contaimageid",
		ContainerImageName: "container_image_name",
	***REMOVED***

	l, err := New(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ctx, unblock := context.WithCancel(context.Background())
	hec.withBlock(ctx)
	defer unblock()

	batchSendTimeout = 1 * time.Second

	if err := l.Log(&logger.Message***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		l.Close()
		close(done)
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(60 * time.Second):
		buf := make([]byte, 1e6)
		buf = buf[:runtime.Stack(buf, true)]
		t.Logf("STACK DUMP: \n\n%s\n\n", string(buf))
		t.Fatal("timeout waiting for close to finish")
	case <-done:
	***REMOVED***
***REMOVED***
