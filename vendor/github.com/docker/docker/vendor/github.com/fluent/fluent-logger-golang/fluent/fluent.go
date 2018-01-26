package fluent

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/tinylib/msgp/msgp"
)

const (
	defaultHost                   = "127.0.0.1"
	defaultNetwork                = "tcp"
	defaultSocketPath             = ""
	defaultPort                   = 24224
	defaultTimeout                = 3 * time.Second
	defaultWriteTimeout           = time.Duration(0) // Write() will not time out
	defaultBufferLimit            = 8 * 1024 * 1024
	defaultRetryWait              = 500
	defaultMaxRetry               = 13
	defaultReconnectWaitIncreRate = 1.5
	// Default sub-second precision value to false since it is only compatible
	// with fluentd versions v0.14 and above.
	defaultSubSecondPrecision = false
)

type Config struct ***REMOVED***
	FluentPort       int           `json:"fluent_port"`
	FluentHost       string        `json:"fluent_host"`
	FluentNetwork    string        `json:"fluent_network"`
	FluentSocketPath string        `json:"fluent_socket_path"`
	Timeout          time.Duration `json:"timeout"`
	WriteTimeout     time.Duration `json:"write_timeout"`
	BufferLimit      int           `json:"buffer_limit"`
	RetryWait        int           `json:"retry_wait"`
	MaxRetry         int           `json:"max_retry"`
	TagPrefix        string        `json:"tag_prefix"`
	AsyncConnect     bool          `json:"async_connect"`
	MarshalAsJSON    bool          `json:"marshal_as_json"`

	// Sub-second precision timestamps are only possible for those using fluentd
	// v0.14+ and serializing their messages with msgpack.
	SubSecondPrecision bool `json:"sub_second_precision"`
***REMOVED***

type Fluent struct ***REMOVED***
	Config

	mubuff  sync.Mutex
	pending []byte

	muconn       sync.Mutex
	conn         net.Conn
	reconnecting bool
***REMOVED***

// New creates a new Logger.
func New(config Config) (f *Fluent, err error) ***REMOVED***
	if config.FluentNetwork == "" ***REMOVED***
		config.FluentNetwork = defaultNetwork
	***REMOVED***
	if config.FluentHost == "" ***REMOVED***
		config.FluentHost = defaultHost
	***REMOVED***
	if config.FluentPort == 0 ***REMOVED***
		config.FluentPort = defaultPort
	***REMOVED***
	if config.FluentSocketPath == "" ***REMOVED***
		config.FluentSocketPath = defaultSocketPath
	***REMOVED***
	if config.Timeout == 0 ***REMOVED***
		config.Timeout = defaultTimeout
	***REMOVED***
	if config.WriteTimeout == 0 ***REMOVED***
		config.WriteTimeout = defaultWriteTimeout
	***REMOVED***
	if config.BufferLimit == 0 ***REMOVED***
		config.BufferLimit = defaultBufferLimit
	***REMOVED***
	if config.RetryWait == 0 ***REMOVED***
		config.RetryWait = defaultRetryWait
	***REMOVED***
	if config.MaxRetry == 0 ***REMOVED***
		config.MaxRetry = defaultMaxRetry
	***REMOVED***
	if config.AsyncConnect ***REMOVED***
		f = &Fluent***REMOVED***Config: config, reconnecting: true***REMOVED***
		go f.reconnect()
	***REMOVED*** else ***REMOVED***
		f = &Fluent***REMOVED***Config: config, reconnecting: false***REMOVED***
		err = f.connect()
	***REMOVED***
	return
***REMOVED***

// Post writes the output for a logging event.
//
// Examples:
//
//  // send map[string]
//  mapStringData := map[string]string***REMOVED***
//  	"foo":  "bar",
//  ***REMOVED***
//  f.Post("tag_name", mapStringData)
//
//  // send message with specified time
//  mapStringData := map[string]string***REMOVED***
//  	"foo":  "bar",
//  ***REMOVED***
//  tm := time.Now()
//  f.PostWithTime("tag_name", tm, mapStringData)
//
//  // send struct
//  structData := struct ***REMOVED***
//  		Name string `msg:"name"`
//  ***REMOVED*** ***REMOVED***
//  		"john smith",
//  ***REMOVED***
//  f.Post("tag_name", structData)
//
func (f *Fluent) Post(tag string, message interface***REMOVED******REMOVED***) error ***REMOVED***
	timeNow := time.Now()
	return f.PostWithTime(tag, timeNow, message)
***REMOVED***

func (f *Fluent) PostWithTime(tag string, tm time.Time, message interface***REMOVED******REMOVED***) error ***REMOVED***
	if len(f.TagPrefix) > 0 ***REMOVED***
		tag = f.TagPrefix + "." + tag
	***REMOVED***

	if m, ok := message.(msgp.Marshaler); ok ***REMOVED***
		return f.EncodeAndPostData(tag, tm, m)
	***REMOVED***

	msg := reflect.ValueOf(message)
	msgtype := msg.Type()

	if msgtype.Kind() == reflect.Struct ***REMOVED***
		// message should be tagged by "codec" or "msg"
		kv := make(map[string]interface***REMOVED******REMOVED***)
		fields := msgtype.NumField()
		for i := 0; i < fields; i++ ***REMOVED***
			field := msgtype.Field(i)
			name := field.Name
			if n1 := field.Tag.Get("msg"); n1 != "" ***REMOVED***
				name = n1
			***REMOVED*** else if n2 := field.Tag.Get("codec"); n2 != "" ***REMOVED***
				name = n2
			***REMOVED***
			kv[name] = msg.FieldByIndex(field.Index).Interface()
		***REMOVED***
		return f.EncodeAndPostData(tag, tm, kv)
	***REMOVED***

	if msgtype.Kind() != reflect.Map ***REMOVED***
		return errors.New("fluent#PostWithTime: message must be a map")
	***REMOVED*** else if msgtype.Key().Kind() != reflect.String ***REMOVED***
		return errors.New("fluent#PostWithTime: map keys must be strings")
	***REMOVED***

	kv := make(map[string]interface***REMOVED******REMOVED***)
	for _, k := range msg.MapKeys() ***REMOVED***
		kv[k.String()] = msg.MapIndex(k).Interface()
	***REMOVED***

	return f.EncodeAndPostData(tag, tm, kv)
***REMOVED***

func (f *Fluent) EncodeAndPostData(tag string, tm time.Time, message interface***REMOVED******REMOVED***) error ***REMOVED***
	var data []byte
	var err error
	if data, err = f.EncodeData(tag, tm, message); err != nil ***REMOVED***
		return fmt.Errorf("fluent#EncodeAndPostData: can't convert '%#v' to msgpack:%v", message, err)
	***REMOVED***
	return f.postRawData(data)
***REMOVED***

// Deprecated: Use EncodeAndPostData instead
func (f *Fluent) PostRawData(data []byte) ***REMOVED***
	f.postRawData(data)
***REMOVED***

func (f *Fluent) postRawData(data []byte) error ***REMOVED***
	if err := f.appendBuffer(data); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := f.send(); err != nil ***REMOVED***
		f.close()
		return err
	***REMOVED***
	return nil
***REMOVED***

// For sending forward protocol adopted JSON
type MessageChunk struct ***REMOVED***
	message Message
***REMOVED***

// Golang default marshaler does not support
// ["value", "value2", ***REMOVED***"key":"value"***REMOVED***] style marshaling.
// So, it should write JSON marshaler by hand.
func (chunk *MessageChunk) MarshalJSON() ([]byte, error) ***REMOVED***
	data, err := json.Marshal(chunk.message.Record)
	return []byte(fmt.Sprintf("[\"%s\",%d,%s,null]", chunk.message.Tag,
		chunk.message.Time, data)), err
***REMOVED***

func (f *Fluent) EncodeData(tag string, tm time.Time, message interface***REMOVED******REMOVED***) (data []byte, err error) ***REMOVED***
	timeUnix := tm.Unix()
	if f.Config.MarshalAsJSON ***REMOVED***
		msg := Message***REMOVED***Tag: tag, Time: timeUnix, Record: message***REMOVED***
		chunk := &MessageChunk***REMOVED***message: msg***REMOVED***
		data, err = json.Marshal(chunk)
	***REMOVED*** else if f.Config.SubSecondPrecision ***REMOVED***
		msg := &MessageExt***REMOVED***Tag: tag, Time: EventTime(tm), Record: message***REMOVED***
		data, err = msg.MarshalMsg(nil)
	***REMOVED*** else ***REMOVED***
		msg := &Message***REMOVED***Tag: tag, Time: timeUnix, Record: message***REMOVED***
		data, err = msg.MarshalMsg(nil)
	***REMOVED***
	return
***REMOVED***

// Close closes the connection.
func (f *Fluent) Close() (err error) ***REMOVED***
	if len(f.pending) > 0 ***REMOVED***
		err = f.send()
	***REMOVED***
	f.close()
	return
***REMOVED***

// appendBuffer appends data to buffer with lock.
func (f *Fluent) appendBuffer(data []byte) error ***REMOVED***
	f.mubuff.Lock()
	defer f.mubuff.Unlock()
	if len(f.pending)+len(data) > f.Config.BufferLimit ***REMOVED***
		return errors.New(fmt.Sprintf("fluent#appendBuffer: Buffer full, limit %v", f.Config.BufferLimit))
	***REMOVED***
	f.pending = append(f.pending, data...)
	return nil
***REMOVED***

// close closes the connection.
func (f *Fluent) close() ***REMOVED***
	f.muconn.Lock()
	if f.conn != nil ***REMOVED***
		f.conn.Close()
		f.conn = nil
	***REMOVED***
	f.muconn.Unlock()
***REMOVED***

// connect establishes a new connection using the specified transport.
func (f *Fluent) connect() (err error) ***REMOVED***
	f.muconn.Lock()
	defer f.muconn.Unlock()

	switch f.Config.FluentNetwork ***REMOVED***
	case "tcp":
		f.conn, err = net.DialTimeout(f.Config.FluentNetwork, f.Config.FluentHost+":"+strconv.Itoa(f.Config.FluentPort), f.Config.Timeout)
	case "unix":
		f.conn, err = net.DialTimeout(f.Config.FluentNetwork, f.Config.FluentSocketPath, f.Config.Timeout)
	default:
		err = net.UnknownNetworkError(f.Config.FluentNetwork)
	***REMOVED***

	if err == nil ***REMOVED***
		f.reconnecting = false
	***REMOVED***
	return
***REMOVED***

func e(x, y float64) int ***REMOVED***
	return int(math.Pow(x, y))
***REMOVED***

func (f *Fluent) reconnect() ***REMOVED***
	for i := 0; ; i++ ***REMOVED***
		err := f.connect()
		if err == nil ***REMOVED***
			f.send()
			return
		***REMOVED***
		if i == f.Config.MaxRetry ***REMOVED***
			// TODO: What we can do when connection failed MaxRetry times?
			panic("fluent#reconnect: failed to reconnect!")
		***REMOVED***
		waitTime := f.Config.RetryWait * e(defaultReconnectWaitIncreRate, float64(i-1))
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	***REMOVED***
***REMOVED***

func (f *Fluent) send() error ***REMOVED***
	f.muconn.Lock()
	defer f.muconn.Unlock()

	if f.conn == nil ***REMOVED***
		if f.reconnecting == false ***REMOVED***
			f.reconnecting = true
			go f.reconnect()
		***REMOVED***
		return errors.New("fluent#send: can't send logs, client is reconnecting")
	***REMOVED***

	f.mubuff.Lock()
	defer f.mubuff.Unlock()

	var err error
	if len(f.pending) > 0 ***REMOVED***
		t := f.Config.WriteTimeout
		if time.Duration(0) < t ***REMOVED***
			f.conn.SetWriteDeadline(time.Now().Add(t))
		***REMOVED*** else ***REMOVED***
			f.conn.SetWriteDeadline(time.Time***REMOVED******REMOVED***)
		***REMOVED***
		_, err = f.conn.Write(f.pending)
		if err != nil ***REMOVED***
			f.conn.Close()
			f.conn = nil
		***REMOVED*** else ***REMOVED***
			f.pending = f.pending[:0]
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***
