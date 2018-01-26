package awslogs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type mockcwlogsclient struct ***REMOVED***
	createLogGroupArgument  chan *cloudwatchlogs.CreateLogGroupInput
	createLogGroupResult    chan *createLogGroupResult
	createLogStreamArgument chan *cloudwatchlogs.CreateLogStreamInput
	createLogStreamResult   chan *createLogStreamResult
	putLogEventsArgument    chan *cloudwatchlogs.PutLogEventsInput
	putLogEventsResult      chan *putLogEventsResult
***REMOVED***

type createLogGroupResult struct ***REMOVED***
	successResult *cloudwatchlogs.CreateLogGroupOutput
	errorResult   error
***REMOVED***

type createLogStreamResult struct ***REMOVED***
	successResult *cloudwatchlogs.CreateLogStreamOutput
	errorResult   error
***REMOVED***

type putLogEventsResult struct ***REMOVED***
	successResult *cloudwatchlogs.PutLogEventsOutput
	errorResult   error
***REMOVED***

func newMockClient() *mockcwlogsclient ***REMOVED***
	return &mockcwlogsclient***REMOVED***
		createLogGroupArgument:  make(chan *cloudwatchlogs.CreateLogGroupInput, 1),
		createLogGroupResult:    make(chan *createLogGroupResult, 1),
		createLogStreamArgument: make(chan *cloudwatchlogs.CreateLogStreamInput, 1),
		createLogStreamResult:   make(chan *createLogStreamResult, 1),
		putLogEventsArgument:    make(chan *cloudwatchlogs.PutLogEventsInput, 1),
		putLogEventsResult:      make(chan *putLogEventsResult, 1),
	***REMOVED***
***REMOVED***

func newMockClientBuffered(buflen int) *mockcwlogsclient ***REMOVED***
	return &mockcwlogsclient***REMOVED***
		createLogStreamArgument: make(chan *cloudwatchlogs.CreateLogStreamInput, buflen),
		createLogStreamResult:   make(chan *createLogStreamResult, buflen),
		putLogEventsArgument:    make(chan *cloudwatchlogs.PutLogEventsInput, buflen),
		putLogEventsResult:      make(chan *putLogEventsResult, buflen),
	***REMOVED***
***REMOVED***

func (m *mockcwlogsclient) CreateLogGroup(input *cloudwatchlogs.CreateLogGroupInput) (*cloudwatchlogs.CreateLogGroupOutput, error) ***REMOVED***
	m.createLogGroupArgument <- input
	output := <-m.createLogGroupResult
	return output.successResult, output.errorResult
***REMOVED***

func (m *mockcwlogsclient) CreateLogStream(input *cloudwatchlogs.CreateLogStreamInput) (*cloudwatchlogs.CreateLogStreamOutput, error) ***REMOVED***
	m.createLogStreamArgument <- input
	output := <-m.createLogStreamResult
	return output.successResult, output.errorResult
***REMOVED***

func (m *mockcwlogsclient) PutLogEvents(input *cloudwatchlogs.PutLogEventsInput) (*cloudwatchlogs.PutLogEventsOutput, error) ***REMOVED***
	events := make([]*cloudwatchlogs.InputLogEvent, len(input.LogEvents))
	copy(events, input.LogEvents)
	m.putLogEventsArgument <- &cloudwatchlogs.PutLogEventsInput***REMOVED***
		LogEvents:     events,
		SequenceToken: input.SequenceToken,
		LogGroupName:  input.LogGroupName,
		LogStreamName: input.LogStreamName,
	***REMOVED***

	// Intended mock output
	output := <-m.putLogEventsResult

	// Checked enforced limits in mock
	totalBytes := 0
	for _, evt := range events ***REMOVED***
		if evt.Message == nil ***REMOVED***
			continue
		***REMOVED***
		eventBytes := len([]byte(*evt.Message))
		if eventBytes > maximumBytesPerEvent ***REMOVED***
			// exceeded per event message size limits
			return nil, fmt.Errorf("maximum bytes per event exceeded: Event too large %d, max allowed: %d", eventBytes, maximumBytesPerEvent)
		***REMOVED***
		// total event bytes including overhead
		totalBytes += eventBytes + perEventBytes
	***REMOVED***

	if totalBytes > maximumBytesPerPut ***REMOVED***
		// exceeded per put maximum size limit
		return nil, fmt.Errorf("maximum bytes per put exceeded: Upload too large %d, max allowed: %d", totalBytes, maximumBytesPerPut)
	***REMOVED***

	return output.successResult, output.errorResult
***REMOVED***

type mockmetadataclient struct ***REMOVED***
	regionResult chan *regionResult
***REMOVED***

type regionResult struct ***REMOVED***
	successResult string
	errorResult   error
***REMOVED***

func newMockMetadataClient() *mockmetadataclient ***REMOVED***
	return &mockmetadataclient***REMOVED***
		regionResult: make(chan *regionResult, 1),
	***REMOVED***
***REMOVED***

func (m *mockmetadataclient) Region() (string, error) ***REMOVED***
	output := <-m.regionResult
	return output.successResult, output.errorResult
***REMOVED***
