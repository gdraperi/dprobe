// Package gelf provides the log driver for forwarding server logs to
// endpoints that support the Graylog Extended Log Format.
package gelf

import (
	"compress/flate"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/sirupsen/logrus"
)

const name = "gelf"

type gelfLogger struct ***REMOVED***
	writer   gelf.Writer
	info     logger.Info
	hostname string
	rawExtra json.RawMessage
***REMOVED***

func init() ***REMOVED***
	if err := logger.RegisterLogDriver(name, New); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil ***REMOVED***
		logrus.Fatal(err)
	***REMOVED***
***REMOVED***

// New creates a gelf logger using the configuration passed in on the
// context. The supported context configuration variable is gelf-address.
func New(info logger.Info) (logger.Logger, error) ***REMOVED***
	// parse gelf address
	address, err := parseAddress(info.Config["gelf-address"])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// collect extra data for GELF message
	hostname, err := info.Hostname()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("gelf: cannot access hostname to set source field")
	***REMOVED***

	// parse log tag
	tag, err := loggerutils.ParseLogTag(info, loggerutils.DefaultTemplate)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	extra := map[string]interface***REMOVED******REMOVED******REMOVED***
		"_container_id":   info.ContainerID,
		"_container_name": info.Name(),
		"_image_id":       info.ContainerImageID,
		"_image_name":     info.ContainerImageName,
		"_command":        info.Command(),
		"_tag":            tag,
		"_created":        info.ContainerCreated,
	***REMOVED***

	extraAttrs, err := info.ExtraAttributes(func(key string) string ***REMOVED***
		if key[0] == '_' ***REMOVED***
			return key
		***REMOVED***
		return "_" + key
	***REMOVED***)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for k, v := range extraAttrs ***REMOVED***
		extra[k] = v
	***REMOVED***

	rawExtra, err := json.Marshal(extra)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var gelfWriter gelf.Writer
	if address.Scheme == "udp" ***REMOVED***
		gelfWriter, err = newGELFUDPWriter(address.Host, info)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else if address.Scheme == "tcp" ***REMOVED***
		gelfWriter, err = newGELFTCPWriter(address.Host, info)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return &gelfLogger***REMOVED***
		writer:   gelfWriter,
		info:     info,
		hostname: hostname,
		rawExtra: rawExtra,
	***REMOVED***, nil
***REMOVED***

// create new TCP gelfWriter
func newGELFTCPWriter(address string, info logger.Info) (gelf.Writer, error) ***REMOVED***
	gelfWriter, err := gelf.NewTCPWriter(address)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("gelf: cannot connect to GELF endpoint: %s %v", address, err)
	***REMOVED***

	if v, ok := info.Config["gelf-tcp-max-reconnect"]; ok ***REMOVED***
		i, err := strconv.Atoi(v)
		if err != nil || i < 0 ***REMOVED***
			return nil, fmt.Errorf("gelf-tcp-max-reconnect must be a positive integer")
		***REMOVED***
		gelfWriter.MaxReconnect = i
	***REMOVED***

	if v, ok := info.Config["gelf-tcp-reconnect-delay"]; ok ***REMOVED***
		i, err := strconv.Atoi(v)
		if err != nil || i < 0 ***REMOVED***
			return nil, fmt.Errorf("gelf-tcp-reconnect-delay must be a positive integer")
		***REMOVED***
		gelfWriter.ReconnectDelay = time.Duration(i)
	***REMOVED***

	return gelfWriter, nil
***REMOVED***

// create new UDP gelfWriter
func newGELFUDPWriter(address string, info logger.Info) (gelf.Writer, error) ***REMOVED***
	gelfWriter, err := gelf.NewUDPWriter(address)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("gelf: cannot connect to GELF endpoint: %s %v", address, err)
	***REMOVED***

	if v, ok := info.Config["gelf-compression-type"]; ok ***REMOVED***
		switch v ***REMOVED***
		case "gzip":
			gelfWriter.CompressionType = gelf.CompressGzip
		case "zlib":
			gelfWriter.CompressionType = gelf.CompressZlib
		case "none":
			gelfWriter.CompressionType = gelf.CompressNone
		default:
			return nil, fmt.Errorf("gelf: invalid compression type %q", v)
		***REMOVED***
	***REMOVED***

	if v, ok := info.Config["gelf-compression-level"]; ok ***REMOVED***
		val, err := strconv.Atoi(v)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("gelf: invalid compression level %s, err %v", v, err)
		***REMOVED***
		gelfWriter.CompressionLevel = val
	***REMOVED***

	return gelfWriter, nil
***REMOVED***

func (s *gelfLogger) Log(msg *logger.Message) error ***REMOVED***
	level := gelf.LOG_INFO
	if msg.Source == "stderr" ***REMOVED***
		level = gelf.LOG_ERR
	***REMOVED***

	m := gelf.Message***REMOVED***
		Version:  "1.1",
		Host:     s.hostname,
		Short:    string(msg.Line),
		TimeUnix: float64(msg.Timestamp.UnixNano()/int64(time.Millisecond)) / 1000.0,
		Level:    int32(level),
		RawExtra: s.rawExtra,
	***REMOVED***
	logger.PutMessage(msg)

	if err := s.writer.WriteMessage(&m); err != nil ***REMOVED***
		return fmt.Errorf("gelf: cannot send GELF message: %v", err)
	***REMOVED***
	return nil
***REMOVED***

func (s *gelfLogger) Close() error ***REMOVED***
	return s.writer.Close()
***REMOVED***

func (s *gelfLogger) Name() string ***REMOVED***
	return name
***REMOVED***

// ValidateLogOpt looks for gelf specific log option gelf-address.
func ValidateLogOpt(cfg map[string]string) error ***REMOVED***
	address, err := parseAddress(cfg["gelf-address"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for key, val := range cfg ***REMOVED***
		switch key ***REMOVED***
		case "gelf-address":
		case "tag":
		case "labels":
		case "env":
		case "env-regex":
		case "gelf-compression-level":
			if address.Scheme != "udp" ***REMOVED***
				return fmt.Errorf("compression is only supported on UDP")
			***REMOVED***
			i, err := strconv.Atoi(val)
			if err != nil || i < flate.DefaultCompression || i > flate.BestCompression ***REMOVED***
				return fmt.Errorf("unknown value %q for log opt %q for gelf log driver", val, key)
			***REMOVED***
		case "gelf-compression-type":
			if address.Scheme != "udp" ***REMOVED***
				return fmt.Errorf("compression is only supported on UDP")
			***REMOVED***
			switch val ***REMOVED***
			case "gzip", "zlib", "none":
			default:
				return fmt.Errorf("unknown value %q for log opt %q for gelf log driver", val, key)
			***REMOVED***
		case "gelf-tcp-max-reconnect", "gelf-tcp-reconnect-delay":
			if address.Scheme != "tcp" ***REMOVED***
				return fmt.Errorf("%q is only valid for TCP", key)
			***REMOVED***
			i, err := strconv.Atoi(val)
			if err != nil || i < 0 ***REMOVED***
				return fmt.Errorf("%q must be a positive integer", key)
			***REMOVED***
		default:
			return fmt.Errorf("unknown log opt %q for gelf log driver", key)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func parseAddress(address string) (*url.URL, error) ***REMOVED***
	if address == "" ***REMOVED***
		return nil, fmt.Errorf("gelf-address is a required parameter")
	***REMOVED***
	if !urlutil.IsTransportURL(address) ***REMOVED***
		return nil, fmt.Errorf("gelf-address should be in form proto://address, got %v", address)
	***REMOVED***
	url, err := url.Parse(address)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// we support only udp
	if url.Scheme != "udp" && url.Scheme != "tcp" ***REMOVED***
		return nil, fmt.Errorf("gelf: endpoint needs to be TCP or UDP")
	***REMOVED***

	// get host and port
	if _, _, err = net.SplitHostPort(url.Host); err != nil ***REMOVED***
		return nil, fmt.Errorf("gelf: please provide gelf-address as proto://host:port")
	***REMOVED***

	return url, nil
***REMOVED***
