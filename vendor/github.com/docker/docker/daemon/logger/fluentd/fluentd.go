// Package fluentd provides the log driver for forwarding server logs
// to fluentd endpoints.
package fluentd

import (
	"fmt"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/docker/go-units"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type fluentd struct ***REMOVED***
	tag           string
	containerID   string
	containerName string
	writer        *fluent.Fluent
	extra         map[string]string
***REMOVED***

type location struct ***REMOVED***
	protocol string
	host     string
	port     int
	path     string
***REMOVED***

const (
	name = "fluentd"

	defaultProtocol    = "tcp"
	defaultHost        = "127.0.0.1"
	defaultPort        = 24224
	defaultBufferLimit = 1024 * 1024

	// logger tries to reconnect 2**32 - 1 times
	// failed (and panic) after 204 years [ 1.5 ** (2**32 - 1) - 1 seconds]
	defaultRetryWait  = 1000
	defaultMaxRetries = math.MaxInt32

	addressKey            = "fluentd-address"
	bufferLimitKey        = "fluentd-buffer-limit"
	retryWaitKey          = "fluentd-retry-wait"
	maxRetriesKey         = "fluentd-max-retries"
	asyncConnectKey       = "fluentd-async-connect"
	subSecondPrecisionKey = "fluentd-sub-second-precision"
)

func init() ***REMOVED***
	if err := logger.RegisterLogDriver(name, New); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
***REMOVED***

// New creates a fluentd logger using the configuration passed in on
// the context. The supported context configuration variable is
// fluentd-address.
func New(info logger.Info) (logger.Logger, error) ***REMOVED***
	loc, err := parseAddress(info.Config[addressKey])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	tag, err := loggerutils.ParseLogTag(info, loggerutils.DefaultTemplate)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	extra, err := info.ExtraAttributes(nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	bufferLimit := defaultBufferLimit
	if info.Config[bufferLimitKey] != "" ***REMOVED***
		bl64, err := units.RAMInBytes(info.Config[bufferLimitKey])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		bufferLimit = int(bl64)
	***REMOVED***

	retryWait := defaultRetryWait
	if info.Config[retryWaitKey] != "" ***REMOVED***
		rwd, err := time.ParseDuration(info.Config[retryWaitKey])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		retryWait = int(rwd.Seconds() * 1000)
	***REMOVED***

	maxRetries := defaultMaxRetries
	if info.Config[maxRetriesKey] != "" ***REMOVED***
		mr64, err := strconv.ParseUint(info.Config[maxRetriesKey], 10, strconv.IntSize)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		maxRetries = int(mr64)
	***REMOVED***

	asyncConnect := false
	if info.Config[asyncConnectKey] != "" ***REMOVED***
		if asyncConnect, err = strconv.ParseBool(info.Config[asyncConnectKey]); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	subSecondPrecision := false
	if info.Config[subSecondPrecisionKey] != "" ***REMOVED***
		if subSecondPrecision, err = strconv.ParseBool(info.Config[subSecondPrecisionKey]); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	fluentConfig := fluent.Config***REMOVED***
		FluentPort:         loc.port,
		FluentHost:         loc.host,
		FluentNetwork:      loc.protocol,
		FluentSocketPath:   loc.path,
		BufferLimit:        bufferLimit,
		RetryWait:          retryWait,
		MaxRetry:           maxRetries,
		AsyncConnect:       asyncConnect,
		SubSecondPrecision: subSecondPrecision,
	***REMOVED***

	logrus.WithField("container", info.ContainerID).WithField("config", fluentConfig).
		Debug("logging driver fluentd configured")

	log, err := fluent.New(fluentConfig)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &fluentd***REMOVED***
		tag:           tag,
		containerID:   info.ContainerID,
		containerName: info.ContainerName,
		writer:        log,
		extra:         extra,
	***REMOVED***, nil
***REMOVED***

func (f *fluentd) Log(msg *logger.Message) error ***REMOVED***
	data := map[string]string***REMOVED***
		"container_id":   f.containerID,
		"container_name": f.containerName,
		"source":         msg.Source,
		"log":            string(msg.Line),
	***REMOVED***
	for k, v := range f.extra ***REMOVED***
		data[k] = v
	***REMOVED***

	ts := msg.Timestamp
	logger.PutMessage(msg)
	// fluent-logger-golang buffers logs from failures and disconnections,
	// and these are transferred again automatically.
	return f.writer.PostWithTime(f.tag, ts, data)
***REMOVED***

func (f *fluentd) Close() error ***REMOVED***
	return f.writer.Close()
***REMOVED***

func (f *fluentd) Name() string ***REMOVED***
	return name
***REMOVED***

// ValidateLogOpt looks for fluentd specific log option fluentd-address.
func ValidateLogOpt(cfg map[string]string) error ***REMOVED***
	for key := range cfg ***REMOVED***
		switch key ***REMOVED***
		case "env":
		case "env-regex":
		case "labels":
		case "tag":
		case addressKey:
		case bufferLimitKey:
		case retryWaitKey:
		case maxRetriesKey:
		case asyncConnectKey:
		case subSecondPrecisionKey:
			// Accepted
		default:
			return fmt.Errorf("unknown log opt '%s' for fluentd log driver", key)
		***REMOVED***
	***REMOVED***

	_, err := parseAddress(cfg[addressKey])
	return err
***REMOVED***

func parseAddress(address string) (*location, error) ***REMOVED***
	if address == "" ***REMOVED***
		return &location***REMOVED***
			protocol: defaultProtocol,
			host:     defaultHost,
			port:     defaultPort,
			path:     "",
		***REMOVED***, nil
	***REMOVED***

	protocol := defaultProtocol
	givenAddress := address
	if urlutil.IsTransportURL(address) ***REMOVED***
		url, err := url.Parse(address)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "invalid fluentd-address %s", givenAddress)
		***REMOVED***
		// unix and unixgram socket
		if url.Scheme == "unix" || url.Scheme == "unixgram" ***REMOVED***
			return &location***REMOVED***
				protocol: url.Scheme,
				host:     "",
				port:     0,
				path:     url.Path,
			***REMOVED***, nil
		***REMOVED***
		// tcp|udp
		protocol = url.Scheme
		address = url.Host
	***REMOVED***

	host, port, err := net.SplitHostPort(address)
	if err != nil ***REMOVED***
		if !strings.Contains(err.Error(), "missing port in address") ***REMOVED***
			return nil, errors.Wrapf(err, "invalid fluentd-address %s", givenAddress)
		***REMOVED***
		return &location***REMOVED***
			protocol: protocol,
			host:     host,
			port:     defaultPort,
			path:     "",
		***REMOVED***, nil
	***REMOVED***

	portnum, err := strconv.Atoi(port)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "invalid fluentd-address %s", givenAddress)
	***REMOVED***
	return &location***REMOVED***
		protocol: protocol,
		host:     host,
		port:     portnum,
		path:     "",
	***REMOVED***, nil
***REMOVED***
