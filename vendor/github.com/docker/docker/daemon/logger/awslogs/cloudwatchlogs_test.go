package awslogs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/docker/docker/dockerversion"
	"github.com/stretchr/testify/assert"
)

const (
	groupName         = "groupName"
	streamName        = "streamName"
	sequenceToken     = "sequenceToken"
	nextSequenceToken = "nextSequenceToken"
	logline           = "this is a log line\r"
	multilineLogline  = "2017-01-01 01:01:44 This is a multiline log entry\r"
)

// Generates i multi-line events each with j lines
func (l *logStream) logGenerator(lineCount int, multilineCount int) ***REMOVED***
	for i := 0; i < multilineCount; i++ ***REMOVED***
		l.Log(&logger.Message***REMOVED***
			Line:      []byte(multilineLogline),
			Timestamp: time.Time***REMOVED******REMOVED***,
		***REMOVED***)
		for j := 0; j < lineCount; j++ ***REMOVED***
			l.Log(&logger.Message***REMOVED***
				Line:      []byte(logline),
				Timestamp: time.Time***REMOVED******REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func testEventBatch(events []wrappedEvent) *eventBatch ***REMOVED***
	batch := newEventBatch()
	for _, event := range events ***REMOVED***
		eventlen := len([]byte(*event.inputLogEvent.Message))
		batch.add(event, eventlen)
	***REMOVED***
	return batch
***REMOVED***

func TestNewAWSLogsClientUserAgentHandler(t *testing.T) ***REMOVED***
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			regionKey: "us-east-1",
		***REMOVED***,
	***REMOVED***

	client, err := newAWSLogsClient(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	realClient, ok := client.(*cloudwatchlogs.CloudWatchLogs)
	if !ok ***REMOVED***
		t.Fatal("Could not cast client to cloudwatchlogs.CloudWatchLogs")
	***REMOVED***
	buildHandlerList := realClient.Handlers.Build
	request := &request.Request***REMOVED***
		HTTPRequest: &http.Request***REMOVED***
			Header: http.Header***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***
	buildHandlerList.Run(request)
	expectedUserAgentString := fmt.Sprintf("Docker %s (%s) %s/%s (%s; %s; %s)",
		dockerversion.Version, runtime.GOOS, aws.SDKName, aws.SDKVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	userAgent := request.HTTPRequest.Header.Get("User-Agent")
	if userAgent != expectedUserAgentString ***REMOVED***
		t.Errorf("Wrong User-Agent string, expected \"%s\" but was \"%s\"",
			expectedUserAgentString, userAgent)
	***REMOVED***
***REMOVED***

func TestNewAWSLogsClientRegionDetect(t *testing.T) ***REMOVED***
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED******REMOVED***,
	***REMOVED***

	mockMetadata := newMockMetadataClient()
	newRegionFinder = func() regionFinder ***REMOVED***
		return mockMetadata
	***REMOVED***
	mockMetadata.regionResult <- &regionResult***REMOVED***
		successResult: "us-east-1",
	***REMOVED***

	_, err := newAWSLogsClient(info)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestCreateSuccess(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
	***REMOVED***
	mockClient.createLogStreamResult <- &createLogStreamResult***REMOVED******REMOVED***

	err := stream.create()

	if err != nil ***REMOVED***
		t.Errorf("Received unexpected err: %v\n", err)
	***REMOVED***
	argument := <-mockClient.createLogStreamArgument
	if argument.LogGroupName == nil ***REMOVED***
		t.Fatal("Expected non-nil LogGroupName")
	***REMOVED***
	if *argument.LogGroupName != groupName ***REMOVED***
		t.Errorf("Expected LogGroupName to be %s", groupName)
	***REMOVED***
	if argument.LogStreamName == nil ***REMOVED***
		t.Fatal("Expected non-nil LogStreamName")
	***REMOVED***
	if *argument.LogStreamName != streamName ***REMOVED***
		t.Errorf("Expected LogStreamName to be %s", streamName)
	***REMOVED***
***REMOVED***

func TestCreateLogGroupSuccess(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:         mockClient,
		logGroupName:   groupName,
		logStreamName:  streamName,
		logCreateGroup: true,
	***REMOVED***
	mockClient.createLogGroupResult <- &createLogGroupResult***REMOVED******REMOVED***
	mockClient.createLogStreamResult <- &createLogStreamResult***REMOVED******REMOVED***

	err := stream.create()

	if err != nil ***REMOVED***
		t.Errorf("Received unexpected err: %v\n", err)
	***REMOVED***
	argument := <-mockClient.createLogStreamArgument
	if argument.LogGroupName == nil ***REMOVED***
		t.Fatal("Expected non-nil LogGroupName")
	***REMOVED***
	if *argument.LogGroupName != groupName ***REMOVED***
		t.Errorf("Expected LogGroupName to be %s", groupName)
	***REMOVED***
	if argument.LogStreamName == nil ***REMOVED***
		t.Fatal("Expected non-nil LogStreamName")
	***REMOVED***
	if *argument.LogStreamName != streamName ***REMOVED***
		t.Errorf("Expected LogStreamName to be %s", streamName)
	***REMOVED***
***REMOVED***

func TestCreateError(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client: mockClient,
	***REMOVED***
	mockClient.createLogStreamResult <- &createLogStreamResult***REMOVED***
		errorResult: errors.New("Error"),
	***REMOVED***

	err := stream.create()

	if err == nil ***REMOVED***
		t.Fatal("Expected non-nil err")
	***REMOVED***
***REMOVED***

func TestCreateAlreadyExists(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client: mockClient,
	***REMOVED***
	mockClient.createLogStreamResult <- &createLogStreamResult***REMOVED***
		errorResult: awserr.New(resourceAlreadyExistsCode, "", nil),
	***REMOVED***

	err := stream.create()

	if err != nil ***REMOVED***
		t.Fatal("Expected nil err")
	***REMOVED***
***REMOVED***

func TestPublishBatchSuccess(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	events := []wrappedEvent***REMOVED***
		***REMOVED***
			inputLogEvent: &cloudwatchlogs.InputLogEvent***REMOVED***
				Message: aws.String(logline),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	stream.publishBatch(testEventBatch(events))
	if stream.sequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil sequenceToken")
	***REMOVED***
	if *stream.sequenceToken != nextSequenceToken ***REMOVED***
		t.Errorf("Expected sequenceToken to be %s, but was %s", nextSequenceToken, *stream.sequenceToken)
	***REMOVED***
	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if argument.SequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput.SequenceToken")
	***REMOVED***
	if *argument.SequenceToken != sequenceToken ***REMOVED***
		t.Errorf("Expected PutLogEventsInput.SequenceToken to be %s, but was %s", sequenceToken, *argument.SequenceToken)
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 element, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if argument.LogEvents[0] != events[0].inputLogEvent ***REMOVED***
		t.Error("Expected event to equal input")
	***REMOVED***
***REMOVED***

func TestPublishBatchError(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		errorResult: errors.New("Error"),
	***REMOVED***

	events := []wrappedEvent***REMOVED***
		***REMOVED***
			inputLogEvent: &cloudwatchlogs.InputLogEvent***REMOVED***
				Message: aws.String(logline),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	stream.publishBatch(testEventBatch(events))
	if stream.sequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil sequenceToken")
	***REMOVED***
	if *stream.sequenceToken != sequenceToken ***REMOVED***
		t.Errorf("Expected sequenceToken to be %s, but was %s", sequenceToken, *stream.sequenceToken)
	***REMOVED***
***REMOVED***

func TestPublishBatchInvalidSeqSuccess(t *testing.T) ***REMOVED***
	mockClient := newMockClientBuffered(2)
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		errorResult: awserr.New(invalidSequenceTokenCode, "use token token", nil),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***

	events := []wrappedEvent***REMOVED***
		***REMOVED***
			inputLogEvent: &cloudwatchlogs.InputLogEvent***REMOVED***
				Message: aws.String(logline),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	stream.publishBatch(testEventBatch(events))
	if stream.sequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil sequenceToken")
	***REMOVED***
	if *stream.sequenceToken != nextSequenceToken ***REMOVED***
		t.Errorf("Expected sequenceToken to be %s, but was %s", nextSequenceToken, *stream.sequenceToken)
	***REMOVED***

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if argument.SequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput.SequenceToken")
	***REMOVED***
	if *argument.SequenceToken != sequenceToken ***REMOVED***
		t.Errorf("Expected PutLogEventsInput.SequenceToken to be %s, but was %s", sequenceToken, *argument.SequenceToken)
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 element, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if argument.LogEvents[0] != events[0].inputLogEvent ***REMOVED***
		t.Error("Expected event to equal input")
	***REMOVED***

	argument = <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if argument.SequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput.SequenceToken")
	***REMOVED***
	if *argument.SequenceToken != "token" ***REMOVED***
		t.Errorf("Expected PutLogEventsInput.SequenceToken to be %s, but was %s", "token", *argument.SequenceToken)
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 element, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if argument.LogEvents[0] != events[0].inputLogEvent ***REMOVED***
		t.Error("Expected event to equal input")
	***REMOVED***
***REMOVED***

func TestPublishBatchAlreadyAccepted(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		errorResult: awserr.New(dataAlreadyAcceptedCode, "use token token", nil),
	***REMOVED***

	events := []wrappedEvent***REMOVED***
		***REMOVED***
			inputLogEvent: &cloudwatchlogs.InputLogEvent***REMOVED***
				Message: aws.String(logline),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	stream.publishBatch(testEventBatch(events))
	if stream.sequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil sequenceToken")
	***REMOVED***
	if *stream.sequenceToken != "token" ***REMOVED***
		t.Errorf("Expected sequenceToken to be %s, but was %s", "token", *stream.sequenceToken)
	***REMOVED***

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if argument.SequenceToken == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput.SequenceToken")
	***REMOVED***
	if *argument.SequenceToken != sequenceToken ***REMOVED***
		t.Errorf("Expected PutLogEventsInput.SequenceToken to be %s, but was %s", sequenceToken, *argument.SequenceToken)
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 element, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if argument.LogEvents[0] != events[0].inputLogEvent ***REMOVED***
		t.Error("Expected event to equal input")
	***REMOVED***
***REMOVED***

func TestCollectBatchSimple(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
		messages:      make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	ticks := make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)

	ticks <- time.Time***REMOVED******REMOVED***
	stream.Close()

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 element, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if *argument.LogEvents[0].Message != logline ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", logline, *argument.LogEvents[0].Message)
	***REMOVED***
***REMOVED***

func TestCollectBatchTicker(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
		messages:      make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	ticks := make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline + " 1"),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline + " 2"),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)

	ticks <- time.Time***REMOVED******REMOVED***

	// Verify first batch
	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != 2 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 2 elements, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if *argument.LogEvents[0].Message != logline+" 1" ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", logline+" 1", *argument.LogEvents[0].Message)
	***REMOVED***
	if *argument.LogEvents[1].Message != logline+" 2" ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", logline+" 2", *argument.LogEvents[0].Message)
	***REMOVED***

	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline + " 3"),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)

	ticks <- time.Time***REMOVED******REMOVED***
	argument = <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 elements, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if *argument.LogEvents[0].Message != logline+" 3" ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", logline+" 3", *argument.LogEvents[0].Message)
	***REMOVED***

	stream.Close()

***REMOVED***

func TestCollectBatchMultilinePattern(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	multilinePattern := regexp.MustCompile("xxxx")
	stream := &logStream***REMOVED***
		client:           mockClient,
		logGroupName:     groupName,
		logStreamName:    streamName,
		multilinePattern: multilinePattern,
		sequenceToken:    aws.String(sequenceToken),
		messages:         make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	ticks := make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Now(),
	***REMOVED***)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Now(),
	***REMOVED***)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte("xxxx " + logline),
		Timestamp: time.Now(),
	***REMOVED***)

	ticks <- time.Now()

	// Verify single multiline event
	argument := <-mockClient.putLogEventsArgument
	assert.NotNil(t, argument, "Expected non-nil PutLogEventsInput")
	assert.Equal(t, 1, len(argument.LogEvents), "Expected single multiline event")
	assert.Equal(t, logline+"\n"+logline+"\n", *argument.LogEvents[0].Message, "Received incorrect multiline message")

	stream.Close()

	// Verify single event
	argument = <-mockClient.putLogEventsArgument
	assert.NotNil(t, argument, "Expected non-nil PutLogEventsInput")
	assert.Equal(t, 1, len(argument.LogEvents), "Expected single multiline event")
	assert.Equal(t, "xxxx "+logline+"\n", *argument.LogEvents[0].Message, "Received incorrect multiline message")
***REMOVED***

func BenchmarkCollectBatch(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		mockClient := newMockClient()
		stream := &logStream***REMOVED***
			client:        mockClient,
			logGroupName:  groupName,
			logStreamName: streamName,
			sequenceToken: aws.String(sequenceToken),
			messages:      make(chan *logger.Message),
		***REMOVED***
		mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
			successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
				NextSequenceToken: aws.String(nextSequenceToken),
			***REMOVED***,
		***REMOVED***
		ticks := make(chan time.Time)
		newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
			return &time.Ticker***REMOVED***
				C: ticks,
			***REMOVED***
		***REMOVED***

		go stream.collectBatch()
		stream.logGenerator(10, 100)
		ticks <- time.Time***REMOVED******REMOVED***
		stream.Close()
	***REMOVED***
***REMOVED***

func BenchmarkCollectBatchMultilinePattern(b *testing.B) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		mockClient := newMockClient()
		multilinePattern := regexp.MustCompile(`\d***REMOVED***4***REMOVED***-(?:0[1-9]|1[0-2])-(?:0[1-9]|[1,2][0-9]|3[0,1]) (?:[0,1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]`)
		stream := &logStream***REMOVED***
			client:           mockClient,
			logGroupName:     groupName,
			logStreamName:    streamName,
			multilinePattern: multilinePattern,
			sequenceToken:    aws.String(sequenceToken),
			messages:         make(chan *logger.Message),
		***REMOVED***
		mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
			successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
				NextSequenceToken: aws.String(nextSequenceToken),
			***REMOVED***,
		***REMOVED***
		ticks := make(chan time.Time)
		newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
			return &time.Ticker***REMOVED***
				C: ticks,
			***REMOVED***
		***REMOVED***
		go stream.collectBatch()
		stream.logGenerator(10, 100)
		ticks <- time.Time***REMOVED******REMOVED***
		stream.Close()
	***REMOVED***
***REMOVED***

func TestCollectBatchMultilinePatternMaxEventAge(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	multilinePattern := regexp.MustCompile("xxxx")
	stream := &logStream***REMOVED***
		client:           mockClient,
		logGroupName:     groupName,
		logStreamName:    streamName,
		multilinePattern: multilinePattern,
		sequenceToken:    aws.String(sequenceToken),
		messages:         make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	ticks := make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Now(),
	***REMOVED***)

	// Log an event 1 second later
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Now().Add(time.Second),
	***REMOVED***)

	// Fire ticker batchPublishFrequency seconds later
	ticks <- time.Now().Add(batchPublishFrequency + time.Second)

	// Verify single multiline event is flushed after maximum event buffer age (batchPublishFrequency)
	argument := <-mockClient.putLogEventsArgument
	assert.NotNil(t, argument, "Expected non-nil PutLogEventsInput")
	assert.Equal(t, 1, len(argument.LogEvents), "Expected single multiline event")
	assert.Equal(t, logline+"\n"+logline+"\n", *argument.LogEvents[0].Message, "Received incorrect multiline message")

	// Log an event 1 second later
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Now().Add(time.Second),
	***REMOVED***)

	// Fire ticker another batchPublishFrequency seconds later
	ticks <- time.Now().Add(2*batchPublishFrequency + time.Second)

	// Verify the event buffer is truly flushed - we should only receive a single event
	argument = <-mockClient.putLogEventsArgument
	assert.NotNil(t, argument, "Expected non-nil PutLogEventsInput")
	assert.Equal(t, 1, len(argument.LogEvents), "Expected single multiline event")
	assert.Equal(t, logline+"\n", *argument.LogEvents[0].Message, "Received incorrect multiline message")
	stream.Close()
***REMOVED***

func TestCollectBatchMultilinePatternNegativeEventAge(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	multilinePattern := regexp.MustCompile("xxxx")
	stream := &logStream***REMOVED***
		client:           mockClient,
		logGroupName:     groupName,
		logStreamName:    streamName,
		multilinePattern: multilinePattern,
		sequenceToken:    aws.String(sequenceToken),
		messages:         make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	ticks := make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Now(),
	***REMOVED***)

	// Log an event 1 second later
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Now().Add(time.Second),
	***REMOVED***)

	// Fire ticker in past to simulate negative event buffer age
	ticks <- time.Now().Add(-time.Second)

	// Verify single multiline event is flushed with a negative event buffer age
	argument := <-mockClient.putLogEventsArgument
	assert.NotNil(t, argument, "Expected non-nil PutLogEventsInput")
	assert.Equal(t, 1, len(argument.LogEvents), "Expected single multiline event")
	assert.Equal(t, logline+"\n"+logline+"\n", *argument.LogEvents[0].Message, "Received incorrect multiline message")

	stream.Close()
***REMOVED***

func TestCollectBatchMultilinePatternMaxEventSize(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	multilinePattern := regexp.MustCompile("xxxx")
	stream := &logStream***REMOVED***
		client:           mockClient,
		logGroupName:     groupName,
		logStreamName:    streamName,
		multilinePattern: multilinePattern,
		sequenceToken:    aws.String(sequenceToken),
		messages:         make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	ticks := make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	// Log max event size
	longline := strings.Repeat("A", maximumBytesPerEvent)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(longline),
		Timestamp: time.Now(),
	***REMOVED***)

	// Log short event
	shortline := strings.Repeat("B", 100)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(shortline),
		Timestamp: time.Now(),
	***REMOVED***)

	// Fire ticker
	ticks <- time.Now().Add(batchPublishFrequency)

	// Verify multiline events
	// We expect a maximum sized event with no new line characters and a
	// second short event with a new line character at the end
	argument := <-mockClient.putLogEventsArgument
	assert.NotNil(t, argument, "Expected non-nil PutLogEventsInput")
	assert.Equal(t, 2, len(argument.LogEvents), "Expected two events")
	assert.Equal(t, longline, *argument.LogEvents[0].Message, "Received incorrect multiline message")
	assert.Equal(t, shortline+"\n", *argument.LogEvents[1].Message, "Received incorrect multiline message")
	stream.Close()
***REMOVED***

func TestCollectBatchClose(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
		messages:      make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	var ticks = make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(logline),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)

	// no ticks
	stream.Close()

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 element, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if *argument.LogEvents[0].Message != logline ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", logline, *argument.LogEvents[0].Message)
	***REMOVED***
***REMOVED***

func TestCollectBatchLineSplit(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
		messages:      make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	var ticks = make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	longline := strings.Repeat("A", maximumBytesPerEvent)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(longline + "B"),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)

	// no ticks
	stream.Close()

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != 2 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 2 elements, but contains %d", len(argument.LogEvents))
	***REMOVED***
	if *argument.LogEvents[0].Message != longline ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", longline, *argument.LogEvents[0].Message)
	***REMOVED***
	if *argument.LogEvents[1].Message != "B" ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", "B", *argument.LogEvents[1].Message)
	***REMOVED***
***REMOVED***

func TestCollectBatchMaxEvents(t *testing.T) ***REMOVED***
	mockClient := newMockClientBuffered(1)
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
		messages:      make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	var ticks = make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	line := "A"
	for i := 0; i <= maximumLogEventsPerPut; i++ ***REMOVED***
		stream.Log(&logger.Message***REMOVED***
			Line:      []byte(line),
			Timestamp: time.Time***REMOVED******REMOVED***,
		***REMOVED***)
	***REMOVED***

	// no ticks
	stream.Close()

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != maximumLogEventsPerPut ***REMOVED***
		t.Errorf("Expected LogEvents to contain %d elements, but contains %d", maximumLogEventsPerPut, len(argument.LogEvents))
	***REMOVED***

	argument = <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain %d elements, but contains %d", 1, len(argument.LogEvents))
	***REMOVED***
***REMOVED***

func TestCollectBatchMaxTotalBytes(t *testing.T) ***REMOVED***
	expectedPuts := 2
	mockClient := newMockClientBuffered(expectedPuts)
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
		messages:      make(chan *logger.Message),
	***REMOVED***
	for i := 0; i < expectedPuts; i++ ***REMOVED***
		mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
			successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
				NextSequenceToken: aws.String(nextSequenceToken),
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	var ticks = make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	numPayloads := maximumBytesPerPut / (maximumBytesPerEvent + perEventBytes)
	// maxline is the maximum line that could be submitted after
	// accounting for its overhead.
	maxline := strings.Repeat("A", maximumBytesPerPut-(perEventBytes*numPayloads))
	// This will be split and batched up to the `maximumBytesPerPut'
	// (+/- `maximumBytesPerEvent'). This /should/ be aligned, but
	// should also tolerate an offset within that range.
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(maxline[:len(maxline)/2]),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte(maxline[len(maxline)/2:]),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)
	stream.Log(&logger.Message***REMOVED***
		Line:      []byte("B"),
		Timestamp: time.Time***REMOVED******REMOVED***,
	***REMOVED***)

	// no ticks, guarantee batch by size (and chan close)
	stream.Close()

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***

	// Should total to the maximum allowed bytes.
	eventBytes := 0
	for _, event := range argument.LogEvents ***REMOVED***
		eventBytes += len(*event.Message)
	***REMOVED***
	eventsOverhead := len(argument.LogEvents) * perEventBytes
	payloadTotal := eventBytes + eventsOverhead
	// lowestMaxBatch allows the payload to be offset if the messages
	// don't lend themselves to align with the maximum event size.
	lowestMaxBatch := maximumBytesPerPut - maximumBytesPerEvent

	if payloadTotal > maximumBytesPerPut ***REMOVED***
		t.Errorf("Expected <= %d bytes but was %d", maximumBytesPerPut, payloadTotal)
	***REMOVED***
	if payloadTotal < lowestMaxBatch ***REMOVED***
		t.Errorf("Batch to be no less than %d but was %d", lowestMaxBatch, payloadTotal)
	***REMOVED***

	argument = <-mockClient.putLogEventsArgument
	if len(argument.LogEvents) != 1 ***REMOVED***
		t.Errorf("Expected LogEvents to contain 1 elements, but contains %d", len(argument.LogEvents))
	***REMOVED***
	message := *argument.LogEvents[len(argument.LogEvents)-1].Message
	if message[len(message)-1:] != "B" ***REMOVED***
		t.Errorf("Expected message to be %s but was %s", "B", message[len(message)-1:])
	***REMOVED***
***REMOVED***

func TestCollectBatchWithDuplicateTimestamps(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: streamName,
		sequenceToken: aws.String(sequenceToken),
		messages:      make(chan *logger.Message),
	***REMOVED***
	mockClient.putLogEventsResult <- &putLogEventsResult***REMOVED***
		successResult: &cloudwatchlogs.PutLogEventsOutput***REMOVED***
			NextSequenceToken: aws.String(nextSequenceToken),
		***REMOVED***,
	***REMOVED***
	ticks := make(chan time.Time)
	newTicker = func(_ time.Duration) *time.Ticker ***REMOVED***
		return &time.Ticker***REMOVED***
			C: ticks,
		***REMOVED***
	***REMOVED***

	go stream.collectBatch()

	times := maximumLogEventsPerPut
	expectedEvents := []*cloudwatchlogs.InputLogEvent***REMOVED******REMOVED***
	timestamp := time.Now()
	for i := 0; i < times; i++ ***REMOVED***
		line := fmt.Sprintf("%d", i)
		if i%2 == 0 ***REMOVED***
			timestamp.Add(1 * time.Nanosecond)
		***REMOVED***
		stream.Log(&logger.Message***REMOVED***
			Line:      []byte(line),
			Timestamp: timestamp,
		***REMOVED***)
		expectedEvents = append(expectedEvents, &cloudwatchlogs.InputLogEvent***REMOVED***
			Message:   aws.String(line),
			Timestamp: aws.Int64(timestamp.UnixNano() / int64(time.Millisecond)),
		***REMOVED***)
	***REMOVED***

	ticks <- time.Time***REMOVED******REMOVED***
	stream.Close()

	argument := <-mockClient.putLogEventsArgument
	if argument == nil ***REMOVED***
		t.Fatal("Expected non-nil PutLogEventsInput")
	***REMOVED***
	if len(argument.LogEvents) != times ***REMOVED***
		t.Errorf("Expected LogEvents to contain %d elements, but contains %d", times, len(argument.LogEvents))
	***REMOVED***
	for i := 0; i < times; i++ ***REMOVED***
		if !reflect.DeepEqual(*argument.LogEvents[i], *expectedEvents[i]) ***REMOVED***
			t.Errorf("Expected event to be %v but was %v", *expectedEvents[i], *argument.LogEvents[i])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseLogOptionsMultilinePattern(t *testing.T) ***REMOVED***
	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			multilinePatternKey: "^xxxx",
		***REMOVED***,
	***REMOVED***

	multilinePattern, err := parseMultilineOptions(info)
	assert.Nil(t, err, "Received unexpected error")
	assert.True(t, multilinePattern.MatchString("xxxx"), "No multiline pattern match found")
***REMOVED***

func TestParseLogOptionsDatetimeFormat(t *testing.T) ***REMOVED***
	datetimeFormatTests := []struct ***REMOVED***
		format string
		match  string
	***REMOVED******REMOVED***
		***REMOVED***"%d/%m/%y %a %H:%M:%S%L %Z", "31/12/10 Mon 08:42:44.345 NZDT"***REMOVED***,
		***REMOVED***"%Y-%m-%d %A %I:%M:%S.%f%p%z", "2007-12-04 Monday 08:42:44.123456AM+1200"***REMOVED***,
		***REMOVED***"%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b|%b", "Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec"***REMOVED***,
		***REMOVED***"%B|%B|%B|%B|%B|%B|%B|%B|%B|%B|%B|%B", "January|February|March|April|May|June|July|August|September|October|November|December"***REMOVED***,
		***REMOVED***"%A|%A|%A|%A|%A|%A|%A", "Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday"***REMOVED***,
		***REMOVED***"%a|%a|%a|%a|%a|%a|%a", "Mon|Tue|Wed|Thu|Fri|Sat|Sun"***REMOVED***,
		***REMOVED***"Day of the week: %w, Day of the year: %j", "Day of the week: 4, Day of the year: 091"***REMOVED***,
	***REMOVED***
	for _, dt := range datetimeFormatTests ***REMOVED***
		t.Run(dt.match, func(t *testing.T) ***REMOVED***
			info := logger.Info***REMOVED***
				Config: map[string]string***REMOVED***
					datetimeFormatKey: dt.format,
				***REMOVED***,
			***REMOVED***
			multilinePattern, err := parseMultilineOptions(info)
			assert.Nil(t, err, "Received unexpected error")
			assert.True(t, multilinePattern.MatchString(dt.match), "No multiline pattern match found")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestValidateLogOptionsDatetimeFormatAndMultilinePattern(t *testing.T) ***REMOVED***
	cfg := map[string]string***REMOVED***
		multilinePatternKey: "^xxxx",
		datetimeFormatKey:   "%Y-%m-%d",
		logGroupKey:         groupName,
	***REMOVED***
	conflictingLogOptionsError := "you cannot configure log opt 'awslogs-datetime-format' and 'awslogs-multiline-pattern' at the same time"

	err := ValidateLogOpt(cfg)
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, err.Error(), conflictingLogOptionsError, "Received invalid error")
***REMOVED***

func TestCreateTagSuccess(t *testing.T) ***REMOVED***
	mockClient := newMockClient()
	info := logger.Info***REMOVED***
		ContainerName: "/test-container",
		ContainerID:   "container-abcdefghijklmnopqrstuvwxyz01234567890",
		Config:        map[string]string***REMOVED***"tag": "***REMOVED******REMOVED***.Name***REMOVED******REMOVED***/***REMOVED******REMOVED***.FullID***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***
	logStreamName, e := loggerutils.ParseLogTag(info, loggerutils.DefaultTemplate)
	if e != nil ***REMOVED***
		t.Errorf("Error generating tag: %q", e)
	***REMOVED***
	stream := &logStream***REMOVED***
		client:        mockClient,
		logGroupName:  groupName,
		logStreamName: logStreamName,
	***REMOVED***
	mockClient.createLogStreamResult <- &createLogStreamResult***REMOVED******REMOVED***

	err := stream.create()

	if err != nil ***REMOVED***
		t.Errorf("Received unexpected err: %v\n", err)
	***REMOVED***
	argument := <-mockClient.createLogStreamArgument

	if *argument.LogStreamName != "test-container/container-abcdefghijklmnopqrstuvwxyz01234567890" ***REMOVED***
		t.Errorf("Expected LogStreamName to be %s", "test-container/container-abcdefghijklmnopqrstuvwxyz01234567890")
	***REMOVED***
***REMOVED***

func TestIsSizedLogger(t *testing.T) ***REMOVED***
	awslogs := &logStream***REMOVED******REMOVED***
	assert.Implements(t, (*logger.SizedLogger)(nil), awslogs, "awslogs should implement SizedLogger")
***REMOVED***

func BenchmarkUnwrapEvents(b *testing.B) ***REMOVED***
	events := make([]wrappedEvent, maximumLogEventsPerPut)
	for i := 0; i < maximumLogEventsPerPut; i++ ***REMOVED***
		mes := strings.Repeat("0", maximumBytesPerEvent)
		events[i].inputLogEvent = &cloudwatchlogs.InputLogEvent***REMOVED***
			Message: &mes,
		***REMOVED***
	***REMOVED***

	as := assert.New(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		res := unwrapEvents(events)
		as.Len(res, maximumLogEventsPerPut)
	***REMOVED***
***REMOVED***

func TestNewAWSLogsClientCredentialEndpointDetect(t *testing.T) ***REMOVED***
	// required for the cloudwatchlogs client
	os.Setenv("AWS_REGION", "us-west-2")
	defer os.Unsetenv("AWS_REGION")

	credsResp := `***REMOVED***
		"AccessKeyId" :    "test-access-key-id",
		"SecretAccessKey": "test-secret-access-key"
		***REMOVED***`

	expectedAccessKeyID := "test-access-key-id"
	expectedSecretAccessKey := "test-secret-access-key"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, credsResp)
	***REMOVED***))
	defer testServer.Close()

	// set the SDKEndpoint in the driver
	newSDKEndpoint = testServer.URL

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED******REMOVED***,
	***REMOVED***

	info.Config["awslogs-credentials-endpoint"] = "/creds"

	c, err := newAWSLogsClient(info)
	assert.NoError(t, err)

	client := c.(*cloudwatchlogs.CloudWatchLogs)

	creds, err := client.Config.Credentials.Get()
	assert.NoError(t, err)

	assert.Equal(t, expectedAccessKeyID, creds.AccessKeyID)
	assert.Equal(t, expectedSecretAccessKey, creds.SecretAccessKey)
***REMOVED***

func TestNewAWSLogsClientCredentialEnvironmentVariable(t *testing.T) ***REMOVED***
	// required for the cloudwatchlogs client
	os.Setenv("AWS_REGION", "us-west-2")
	defer os.Unsetenv("AWS_REGION")

	expectedAccessKeyID := "test-access-key-id"
	expectedSecretAccessKey := "test-secret-access-key"

	os.Setenv("AWS_ACCESS_KEY_ID", expectedAccessKeyID)
	defer os.Unsetenv("AWS_ACCESS_KEY_ID")

	os.Setenv("AWS_SECRET_ACCESS_KEY", expectedSecretAccessKey)
	defer os.Unsetenv("AWS_SECRET_ACCESS_KEY")

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED******REMOVED***,
	***REMOVED***

	c, err := newAWSLogsClient(info)
	assert.NoError(t, err)

	client := c.(*cloudwatchlogs.CloudWatchLogs)

	creds, err := client.Config.Credentials.Get()
	assert.NoError(t, err)

	assert.Equal(t, expectedAccessKeyID, creds.AccessKeyID)
	assert.Equal(t, expectedSecretAccessKey, creds.SecretAccessKey)

***REMOVED***

func TestNewAWSLogsClientCredentialSharedFile(t *testing.T) ***REMOVED***
	// required for the cloudwatchlogs client
	os.Setenv("AWS_REGION", "us-west-2")
	defer os.Unsetenv("AWS_REGION")

	expectedAccessKeyID := "test-access-key-id"
	expectedSecretAccessKey := "test-secret-access-key"

	contentStr := `
	[default]
	aws_access_key_id = "test-access-key-id"
	aws_secret_access_key =  "test-secret-access-key"
	`
	content := []byte(contentStr)

	tmpfile, err := ioutil.TempFile("", "example")
	defer os.Remove(tmpfile.Name()) // clean up
	assert.NoError(t, err)

	_, err = tmpfile.Write(content)
	assert.NoError(t, err)

	err = tmpfile.Close()
	assert.NoError(t, err)

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", tmpfile.Name())
	defer os.Unsetenv("AWS_SHARED_CREDENTIALS_FILE")

	info := logger.Info***REMOVED***
		Config: map[string]string***REMOVED******REMOVED***,
	***REMOVED***

	c, err := newAWSLogsClient(info)
	assert.NoError(t, err)

	client := c.(*cloudwatchlogs.CloudWatchLogs)

	creds, err := client.Config.Credentials.Get()
	assert.NoError(t, err)

	assert.Equal(t, expectedAccessKeyID, creds.AccessKeyID)
	assert.Equal(t, expectedSecretAccessKey, creds.SecretAccessKey)
***REMOVED***
