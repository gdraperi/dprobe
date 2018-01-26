// Package awslogs provides the logdriver for forwarding container logs to Amazon CloudWatch Logs
package awslogs

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/docker/docker/dockerversion"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	name                   = "awslogs"
	regionKey              = "awslogs-region"
	regionEnvKey           = "AWS_REGION"
	logGroupKey            = "awslogs-group"
	logStreamKey           = "awslogs-stream"
	logCreateGroupKey      = "awslogs-create-group"
	tagKey                 = "tag"
	datetimeFormatKey      = "awslogs-datetime-format"
	multilinePatternKey    = "awslogs-multiline-pattern"
	credentialsEndpointKey = "awslogs-credentials-endpoint"
	batchPublishFrequency  = 5 * time.Second

	// See: http://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_PutLogEvents.html
	perEventBytes          = 26
	maximumBytesPerPut     = 1048576
	maximumLogEventsPerPut = 10000

	// See: http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/cloudwatch_limits.html
	maximumBytesPerEvent = 262144 - perEventBytes

	resourceAlreadyExistsCode = "ResourceAlreadyExistsException"
	dataAlreadyAcceptedCode   = "DataAlreadyAcceptedException"
	invalidSequenceTokenCode  = "InvalidSequenceTokenException"
	resourceNotFoundCode      = "ResourceNotFoundException"

	credentialsEndpoint = "http://169.254.170.2"

	userAgentHeader = "User-Agent"
)

type logStream struct ***REMOVED***
	logStreamName    string
	logGroupName     string
	logCreateGroup   bool
	multilinePattern *regexp.Regexp
	client           api
	messages         chan *logger.Message
	lock             sync.RWMutex
	closed           bool
	sequenceToken    *string
***REMOVED***

type api interface ***REMOVED***
	CreateLogGroup(*cloudwatchlogs.CreateLogGroupInput) (*cloudwatchlogs.CreateLogGroupOutput, error)
	CreateLogStream(*cloudwatchlogs.CreateLogStreamInput) (*cloudwatchlogs.CreateLogStreamOutput, error)
	PutLogEvents(*cloudwatchlogs.PutLogEventsInput) (*cloudwatchlogs.PutLogEventsOutput, error)
***REMOVED***

type regionFinder interface ***REMOVED***
	Region() (string, error)
***REMOVED***

type wrappedEvent struct ***REMOVED***
	inputLogEvent *cloudwatchlogs.InputLogEvent
	insertOrder   int
***REMOVED***
type byTimestamp []wrappedEvent

// init registers the awslogs driver
func init() ***REMOVED***
	if err := logger.RegisterLogDriver(name, New); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
***REMOVED***

// eventBatch holds the events that are batched for submission and the
// associated data about it.
//
// Warning: this type is not threadsafe and must not be used
// concurrently. This type is expected to be consumed in a single go
// routine and never concurrently.
type eventBatch struct ***REMOVED***
	batch []wrappedEvent
	bytes int
***REMOVED***

// New creates an awslogs logger using the configuration passed in on the
// context.  Supported context configuration variables are awslogs-region,
// awslogs-group, awslogs-stream, awslogs-create-group, awslogs-multiline-pattern
// and awslogs-datetime-format.  When available, configuration is
// also taken from environment variables AWS_REGION, AWS_ACCESS_KEY_ID,
// AWS_SECRET_ACCESS_KEY, the shared credentials file (~/.aws/credentials), and
// the EC2 Instance Metadata Service.
func New(info logger.Info) (logger.Logger, error) ***REMOVED***
	logGroupName := info.Config[logGroupKey]
	logStreamName, err := loggerutils.ParseLogTag(info, "***REMOVED******REMOVED***.FullID***REMOVED******REMOVED***")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	logCreateGroup := false
	if info.Config[logCreateGroupKey] != "" ***REMOVED***
		logCreateGroup, err = strconv.ParseBool(info.Config[logCreateGroupKey])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if info.Config[logStreamKey] != "" ***REMOVED***
		logStreamName = info.Config[logStreamKey]
	***REMOVED***

	multilinePattern, err := parseMultilineOptions(info)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	client, err := newAWSLogsClient(info)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	containerStream := &logStream***REMOVED***
		logStreamName:    logStreamName,
		logGroupName:     logGroupName,
		logCreateGroup:   logCreateGroup,
		multilinePattern: multilinePattern,
		client:           client,
		messages:         make(chan *logger.Message, 4096),
	***REMOVED***
	err = containerStream.create()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go containerStream.collectBatch()

	return containerStream, nil
***REMOVED***

// Parses awslogs-multiline-pattern and awslogs-datetime-format options
// If awslogs-datetime-format is present, convert the format from strftime
// to regexp and return.
// If awslogs-multiline-pattern is present, compile regexp and return
func parseMultilineOptions(info logger.Info) (*regexp.Regexp, error) ***REMOVED***
	dateTimeFormat := info.Config[datetimeFormatKey]
	multilinePatternKey := info.Config[multilinePatternKey]
	// strftime input is parsed into a regular expression
	if dateTimeFormat != "" ***REMOVED***
		// %. matches each strftime format sequence and ReplaceAllStringFunc
		// looks up each format sequence in the conversion table strftimeToRegex
		// to replace with a defined regular expression
		r := regexp.MustCompile("%.")
		multilinePatternKey = r.ReplaceAllStringFunc(dateTimeFormat, func(s string) string ***REMOVED***
			return strftimeToRegex[s]
		***REMOVED***)
	***REMOVED***
	if multilinePatternKey != "" ***REMOVED***
		multilinePattern, err := regexp.Compile(multilinePatternKey)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "awslogs could not parse multiline pattern key %q", multilinePatternKey)
		***REMOVED***
		return multilinePattern, nil
	***REMOVED***
	return nil, nil
***REMOVED***

// Maps strftime format strings to regex
var strftimeToRegex = map[string]string***REMOVED***
	/*weekdayShort          */ `%a`: `(?:Mon|Tue|Wed|Thu|Fri|Sat|Sun)`,
	/*weekdayFull           */ `%A`: `(?:Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)`,
	/*weekdayZeroIndex      */ `%w`: `[0-6]`,
	/*dayZeroPadded         */ `%d`: `(?:0[1-9]|[1,2][0-9]|3[0,1])`,
	/*monthShort            */ `%b`: `(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)`,
	/*monthFull             */ `%B`: `(?:January|February|March|April|May|June|July|August|September|October|November|December)`,
	/*monthZeroPadded       */ `%m`: `(?:0[1-9]|1[0-2])`,
	/*yearCentury           */ `%Y`: `\d***REMOVED***4***REMOVED***`,
	/*yearZeroPadded        */ `%y`: `\d***REMOVED***2***REMOVED***`,
	/*hour24ZeroPadded      */ `%H`: `(?:[0,1][0-9]|2[0-3])`,
	/*hour12ZeroPadded      */ `%I`: `(?:0[0-9]|1[0-2])`,
	/*AM or PM              */ `%p`: "[A,P]M",
	/*minuteZeroPadded      */ `%M`: `[0-5][0-9]`,
	/*secondZeroPadded      */ `%S`: `[0-5][0-9]`,
	/*microsecondZeroPadded */ `%f`: `\d***REMOVED***6***REMOVED***`,
	/*utcOffset             */ `%z`: `[+-]\d***REMOVED***4***REMOVED***`,
	/*tzName                */ `%Z`: `[A-Z]***REMOVED***1,4***REMOVED***T`,
	/*dayOfYearZeroPadded   */ `%j`: `(?:0[0-9][1-9]|[1,2][0-9][0-9]|3[0-5][0-9]|36[0-6])`,
	/*milliseconds          */ `%L`: `\.\d***REMOVED***3***REMOVED***`,
***REMOVED***

// newRegionFinder is a variable such that the implementation
// can be swapped out for unit tests.
var newRegionFinder = func() regionFinder ***REMOVED***
	return ec2metadata.New(session.New())
***REMOVED***

// newSDKEndpoint is a variable such that the implementation
// can be swapped out for unit tests.
var newSDKEndpoint = credentialsEndpoint

// newAWSLogsClient creates the service client for Amazon CloudWatch Logs.
// Customizations to the default client from the SDK include a Docker-specific
// User-Agent string and automatic region detection using the EC2 Instance
// Metadata Service when region is otherwise unspecified.
func newAWSLogsClient(info logger.Info) (api, error) ***REMOVED***
	var region *string
	if os.Getenv(regionEnvKey) != "" ***REMOVED***
		region = aws.String(os.Getenv(regionEnvKey))
	***REMOVED***
	if info.Config[regionKey] != "" ***REMOVED***
		region = aws.String(info.Config[regionKey])
	***REMOVED***
	if region == nil || *region == "" ***REMOVED***
		logrus.Info("Trying to get region from EC2 Metadata")
		ec2MetadataClient := newRegionFinder()
		r, err := ec2MetadataClient.Region()
		if err != nil ***REMOVED***
			logrus.WithFields(logrus.Fields***REMOVED***
				"error": err,
			***REMOVED***).Error("Could not get region from EC2 metadata, environment, or log option")
			return nil, errors.New("Cannot determine region for awslogs driver")
		***REMOVED***
		region = &r
	***REMOVED***

	sess, err := session.NewSession()
	if err != nil ***REMOVED***
		return nil, errors.New("Failed to create a service client session for for awslogs driver")
	***REMOVED***

	// attach region to cloudwatchlogs config
	sess.Config.Region = region

	if uri, ok := info.Config[credentialsEndpointKey]; ok ***REMOVED***
		logrus.Debugf("Trying to get credentials from awslogs-credentials-endpoint")

		endpoint := fmt.Sprintf("%s%s", newSDKEndpoint, uri)
		creds := endpointcreds.NewCredentialsClient(*sess.Config, sess.Handlers, endpoint,
			func(p *endpointcreds.Provider) ***REMOVED***
				p.ExpiryWindow = 5 * time.Minute
			***REMOVED***)

		// attach credentials to cloudwatchlogs config
		sess.Config.Credentials = creds
	***REMOVED***

	logrus.WithFields(logrus.Fields***REMOVED***
		"region": *region,
	***REMOVED***).Debug("Created awslogs client")

	client := cloudwatchlogs.New(sess)

	client.Handlers.Build.PushBackNamed(request.NamedHandler***REMOVED***
		Name: "DockerUserAgentHandler",
		Fn: func(r *request.Request) ***REMOVED***
			currentAgent := r.HTTPRequest.Header.Get(userAgentHeader)
			r.HTTPRequest.Header.Set(userAgentHeader,
				fmt.Sprintf("Docker %s (%s) %s",
					dockerversion.Version, runtime.GOOS, currentAgent))
		***REMOVED***,
	***REMOVED***)
	return client, nil
***REMOVED***

// Name returns the name of the awslogs logging driver
func (l *logStream) Name() string ***REMOVED***
	return name
***REMOVED***

func (l *logStream) BufSize() int ***REMOVED***
	return maximumBytesPerEvent
***REMOVED***

// Log submits messages for logging by an instance of the awslogs logging driver
func (l *logStream) Log(msg *logger.Message) error ***REMOVED***
	l.lock.RLock()
	defer l.lock.RUnlock()
	if !l.closed ***REMOVED***
		l.messages <- msg
	***REMOVED***
	return nil
***REMOVED***

// Close closes the instance of the awslogs logging driver
func (l *logStream) Close() error ***REMOVED***
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.closed ***REMOVED***
		close(l.messages)
	***REMOVED***
	l.closed = true
	return nil
***REMOVED***

// create creates log group and log stream for the instance of the awslogs logging driver
func (l *logStream) create() error ***REMOVED***
	if err := l.createLogStream(); err != nil ***REMOVED***
		if l.logCreateGroup ***REMOVED***
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == resourceNotFoundCode ***REMOVED***
				if err := l.createLogGroup(); err != nil ***REMOVED***
					return err
				***REMOVED***
				return l.createLogStream()
			***REMOVED***
		***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// createLogGroup creates a log group for the instance of the awslogs logging driver
func (l *logStream) createLogGroup() error ***REMOVED***
	if _, err := l.client.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput***REMOVED***
		LogGroupName: aws.String(l.logGroupName),
	***REMOVED***); err != nil ***REMOVED***
		if awsErr, ok := err.(awserr.Error); ok ***REMOVED***
			fields := logrus.Fields***REMOVED***
				"errorCode":      awsErr.Code(),
				"message":        awsErr.Message(),
				"origError":      awsErr.OrigErr(),
				"logGroupName":   l.logGroupName,
				"logCreateGroup": l.logCreateGroup,
			***REMOVED***
			if awsErr.Code() == resourceAlreadyExistsCode ***REMOVED***
				// Allow creation to succeed
				logrus.WithFields(fields).Info("Log group already exists")
				return nil
			***REMOVED***
			logrus.WithFields(fields).Error("Failed to create log group")
		***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// createLogStream creates a log stream for the instance of the awslogs logging driver
func (l *logStream) createLogStream() error ***REMOVED***
	input := &cloudwatchlogs.CreateLogStreamInput***REMOVED***
		LogGroupName:  aws.String(l.logGroupName),
		LogStreamName: aws.String(l.logStreamName),
	***REMOVED***

	_, err := l.client.CreateLogStream(input)

	if err != nil ***REMOVED***
		if awsErr, ok := err.(awserr.Error); ok ***REMOVED***
			fields := logrus.Fields***REMOVED***
				"errorCode":     awsErr.Code(),
				"message":       awsErr.Message(),
				"origError":     awsErr.OrigErr(),
				"logGroupName":  l.logGroupName,
				"logStreamName": l.logStreamName,
			***REMOVED***
			if awsErr.Code() == resourceAlreadyExistsCode ***REMOVED***
				// Allow creation to succeed
				logrus.WithFields(fields).Info("Log stream already exists")
				return nil
			***REMOVED***
			logrus.WithFields(fields).Error("Failed to create log stream")
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

// newTicker is used for time-based batching.  newTicker is a variable such
// that the implementation can be swapped out for unit tests.
var newTicker = func(freq time.Duration) *time.Ticker ***REMOVED***
	return time.NewTicker(freq)
***REMOVED***

// collectBatch executes as a goroutine to perform batching of log events for
// submission to the log stream.  If the awslogs-multiline-pattern or
// awslogs-datetime-format options have been configured, multiline processing
// is enabled, where log messages are stored in an event buffer until a multiline
// pattern match is found, at which point the messages in the event buffer are
// pushed to CloudWatch logs as a single log event.  Multiline messages are processed
// according to the maximumBytesPerPut constraint, and the implementation only
// allows for messages to be buffered for a maximum of 2*batchPublishFrequency
// seconds.  When events are ready to be processed for submission to CloudWatch
// Logs, the processEvents method is called.  If a multiline pattern is not
// configured, log events are submitted to the processEvents method immediately.
func (l *logStream) collectBatch() ***REMOVED***
	ticker := newTicker(batchPublishFrequency)
	var eventBuffer []byte
	var eventBufferTimestamp int64
	var batch = newEventBatch()
	for ***REMOVED***
		select ***REMOVED***
		case t := <-ticker.C:
			// If event buffer is older than batch publish frequency flush the event buffer
			if eventBufferTimestamp > 0 && len(eventBuffer) > 0 ***REMOVED***
				eventBufferAge := t.UnixNano()/int64(time.Millisecond) - eventBufferTimestamp
				eventBufferExpired := eventBufferAge >= int64(batchPublishFrequency)/int64(time.Millisecond)
				eventBufferNegative := eventBufferAge < 0
				if eventBufferExpired || eventBufferNegative ***REMOVED***
					l.processEvent(batch, eventBuffer, eventBufferTimestamp)
					eventBuffer = eventBuffer[:0]
				***REMOVED***
			***REMOVED***
			l.publishBatch(batch)
			batch.reset()
		case msg, more := <-l.messages:
			if !more ***REMOVED***
				// Flush event buffer and release resources
				l.processEvent(batch, eventBuffer, eventBufferTimestamp)
				eventBuffer = eventBuffer[:0]
				l.publishBatch(batch)
				batch.reset()
				return
			***REMOVED***
			if eventBufferTimestamp == 0 ***REMOVED***
				eventBufferTimestamp = msg.Timestamp.UnixNano() / int64(time.Millisecond)
			***REMOVED***
			line := msg.Line
			if l.multilinePattern != nil ***REMOVED***
				if l.multilinePattern.Match(line) || len(eventBuffer)+len(line) > maximumBytesPerEvent ***REMOVED***
					// This is a new log event or we will exceed max bytes per event
					// so flush the current eventBuffer to events and reset timestamp
					l.processEvent(batch, eventBuffer, eventBufferTimestamp)
					eventBufferTimestamp = msg.Timestamp.UnixNano() / int64(time.Millisecond)
					eventBuffer = eventBuffer[:0]
				***REMOVED***
				// Append new line if event is less than max event size
				if len(line) < maximumBytesPerEvent ***REMOVED***
					line = append(line, "\n"...)
				***REMOVED***
				eventBuffer = append(eventBuffer, line...)
				logger.PutMessage(msg)
			***REMOVED*** else ***REMOVED***
				l.processEvent(batch, line, msg.Timestamp.UnixNano()/int64(time.Millisecond))
				logger.PutMessage(msg)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// processEvent processes log events that are ready for submission to CloudWatch
// logs.  Batching is performed on time- and size-bases.  Time-based batching
// occurs at a 5 second interval (defined in the batchPublishFrequency const).
// Size-based batching is performed on the maximum number of events per batch
// (defined in maximumLogEventsPerPut) and the maximum number of total bytes in a
// batch (defined in maximumBytesPerPut).  Log messages are split by the maximum
// bytes per event (defined in maximumBytesPerEvent).  There is a fixed per-event
// byte overhead (defined in perEventBytes) which is accounted for in split- and
// batch-calculations.
func (l *logStream) processEvent(batch *eventBatch, events []byte, timestamp int64) ***REMOVED***
	for len(events) > 0 ***REMOVED***
		// Split line length so it does not exceed the maximum
		lineBytes := len(events)
		if lineBytes > maximumBytesPerEvent ***REMOVED***
			lineBytes = maximumBytesPerEvent
		***REMOVED***
		line := events[:lineBytes]

		event := wrappedEvent***REMOVED***
			inputLogEvent: &cloudwatchlogs.InputLogEvent***REMOVED***
				Message:   aws.String(string(line)),
				Timestamp: aws.Int64(timestamp),
			***REMOVED***,
			insertOrder: batch.count(),
		***REMOVED***

		added := batch.add(event, lineBytes)
		if added ***REMOVED***
			events = events[lineBytes:]
		***REMOVED*** else ***REMOVED***
			l.publishBatch(batch)
			batch.reset()
		***REMOVED***
	***REMOVED***
***REMOVED***

// publishBatch calls PutLogEvents for a given set of InputLogEvents,
// accounting for sequencing requirements (each request must reference the
// sequence token returned by the previous request).
func (l *logStream) publishBatch(batch *eventBatch) ***REMOVED***
	if batch.isEmpty() ***REMOVED***
		return
	***REMOVED***
	cwEvents := unwrapEvents(batch.events())

	nextSequenceToken, err := l.putLogEvents(cwEvents, l.sequenceToken)

	if err != nil ***REMOVED***
		if awsErr, ok := err.(awserr.Error); ok ***REMOVED***
			if awsErr.Code() == dataAlreadyAcceptedCode ***REMOVED***
				// already submitted, just grab the correct sequence token
				parts := strings.Split(awsErr.Message(), " ")
				nextSequenceToken = &parts[len(parts)-1]
				logrus.WithFields(logrus.Fields***REMOVED***
					"errorCode":     awsErr.Code(),
					"message":       awsErr.Message(),
					"logGroupName":  l.logGroupName,
					"logStreamName": l.logStreamName,
				***REMOVED***).Info("Data already accepted, ignoring error")
				err = nil
			***REMOVED*** else if awsErr.Code() == invalidSequenceTokenCode ***REMOVED***
				// sequence code is bad, grab the correct one and retry
				parts := strings.Split(awsErr.Message(), " ")
				token := parts[len(parts)-1]
				nextSequenceToken, err = l.putLogEvents(cwEvents, &token)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		logrus.Error(err)
	***REMOVED*** else ***REMOVED***
		l.sequenceToken = nextSequenceToken
	***REMOVED***
***REMOVED***

// putLogEvents wraps the PutLogEvents API
func (l *logStream) putLogEvents(events []*cloudwatchlogs.InputLogEvent, sequenceToken *string) (*string, error) ***REMOVED***
	input := &cloudwatchlogs.PutLogEventsInput***REMOVED***
		LogEvents:     events,
		SequenceToken: sequenceToken,
		LogGroupName:  aws.String(l.logGroupName),
		LogStreamName: aws.String(l.logStreamName),
	***REMOVED***
	resp, err := l.client.PutLogEvents(input)
	if err != nil ***REMOVED***
		if awsErr, ok := err.(awserr.Error); ok ***REMOVED***
			logrus.WithFields(logrus.Fields***REMOVED***
				"errorCode":     awsErr.Code(),
				"message":       awsErr.Message(),
				"origError":     awsErr.OrigErr(),
				"logGroupName":  l.logGroupName,
				"logStreamName": l.logStreamName,
			***REMOVED***).Error("Failed to put log events")
		***REMOVED***
		return nil, err
	***REMOVED***
	return resp.NextSequenceToken, nil
***REMOVED***

// ValidateLogOpt looks for awslogs-specific log options awslogs-region,
// awslogs-group, awslogs-stream, awslogs-create-group, awslogs-datetime-format,
// awslogs-multiline-pattern
func ValidateLogOpt(cfg map[string]string) error ***REMOVED***
	for key := range cfg ***REMOVED***
		switch key ***REMOVED***
		case logGroupKey:
		case logStreamKey:
		case logCreateGroupKey:
		case regionKey:
		case tagKey:
		case datetimeFormatKey:
		case multilinePatternKey:
		case credentialsEndpointKey:
		default:
			return fmt.Errorf("unknown log opt '%s' for %s log driver", key, name)
		***REMOVED***
	***REMOVED***
	if cfg[logGroupKey] == "" ***REMOVED***
		return fmt.Errorf("must specify a value for log opt '%s'", logGroupKey)
	***REMOVED***
	if cfg[logCreateGroupKey] != "" ***REMOVED***
		if _, err := strconv.ParseBool(cfg[logCreateGroupKey]); err != nil ***REMOVED***
			return fmt.Errorf("must specify valid value for log opt '%s': %v", logCreateGroupKey, err)
		***REMOVED***
	***REMOVED***
	_, datetimeFormatKeyExists := cfg[datetimeFormatKey]
	_, multilinePatternKeyExists := cfg[multilinePatternKey]
	if datetimeFormatKeyExists && multilinePatternKeyExists ***REMOVED***
		return fmt.Errorf("you cannot configure log opt '%s' and '%s' at the same time", datetimeFormatKey, multilinePatternKey)
	***REMOVED***
	return nil
***REMOVED***

// Len returns the length of a byTimestamp slice.  Len is required by the
// sort.Interface interface.
func (slice byTimestamp) Len() int ***REMOVED***
	return len(slice)
***REMOVED***

// Less compares two values in a byTimestamp slice by Timestamp.  Less is
// required by the sort.Interface interface.
func (slice byTimestamp) Less(i, j int) bool ***REMOVED***
	iTimestamp, jTimestamp := int64(0), int64(0)
	if slice != nil && slice[i].inputLogEvent.Timestamp != nil ***REMOVED***
		iTimestamp = *slice[i].inputLogEvent.Timestamp
	***REMOVED***
	if slice != nil && slice[j].inputLogEvent.Timestamp != nil ***REMOVED***
		jTimestamp = *slice[j].inputLogEvent.Timestamp
	***REMOVED***
	if iTimestamp == jTimestamp ***REMOVED***
		return slice[i].insertOrder < slice[j].insertOrder
	***REMOVED***
	return iTimestamp < jTimestamp
***REMOVED***

// Swap swaps two values in a byTimestamp slice with each other.  Swap is
// required by the sort.Interface interface.
func (slice byTimestamp) Swap(i, j int) ***REMOVED***
	slice[i], slice[j] = slice[j], slice[i]
***REMOVED***

func unwrapEvents(events []wrappedEvent) []*cloudwatchlogs.InputLogEvent ***REMOVED***
	cwEvents := make([]*cloudwatchlogs.InputLogEvent, len(events))
	for i, input := range events ***REMOVED***
		cwEvents[i] = input.inputLogEvent
	***REMOVED***
	return cwEvents
***REMOVED***

func newEventBatch() *eventBatch ***REMOVED***
	return &eventBatch***REMOVED***
		batch: make([]wrappedEvent, 0),
		bytes: 0,
	***REMOVED***
***REMOVED***

// events returns a slice of wrappedEvents sorted in order of their
// timestamps and then by their insertion order (see `byTimestamp`).
//
// Warning: this method is not threadsafe and must not be used
// concurrently.
func (b *eventBatch) events() []wrappedEvent ***REMOVED***
	sort.Sort(byTimestamp(b.batch))
	return b.batch
***REMOVED***

// add adds an event to the batch of events accounting for the
// necessary overhead for an event to be logged. An error will be
// returned if the event cannot be added to the batch due to service
// limits.
//
// Warning: this method is not threadsafe and must not be used
// concurrently.
func (b *eventBatch) add(event wrappedEvent, size int) bool ***REMOVED***
	addBytes := size + perEventBytes

	// verify we are still within service limits
	switch ***REMOVED***
	case len(b.batch)+1 > maximumLogEventsPerPut:
		return false
	case b.bytes+addBytes > maximumBytesPerPut:
		return false
	***REMOVED***

	b.bytes += addBytes
	b.batch = append(b.batch, event)

	return true
***REMOVED***

// count is the number of batched events.  Warning: this method
// is not threadsafe and must not be used concurrently.
func (b *eventBatch) count() int ***REMOVED***
	return len(b.batch)
***REMOVED***

// size is the total number of bytes that the batch represents.
//
// Warning: this method is not threadsafe and must not be used
// concurrently.
func (b *eventBatch) size() int ***REMOVED***
	return b.bytes
***REMOVED***

func (b *eventBatch) isEmpty() bool ***REMOVED***
	zeroEvents := b.count() == 0
	zeroSize := b.size() == 0
	return zeroEvents && zeroSize
***REMOVED***

// reset prepares the batch for reuse.
func (b *eventBatch) reset() ***REMOVED***
	b.bytes = 0
	b.batch = b.batch[:0]
***REMOVED***
