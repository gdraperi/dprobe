package loggerutils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils/multireader"
	"github.com/docker/docker/pkg/filenotify"
	"github.com/docker/docker/pkg/pubsub"
	"github.com/docker/docker/pkg/tailfile"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LogFile is Logger implementation for default Docker logging.
type LogFile struct ***REMOVED***
	f             *os.File // store for closing
	closed        bool
	mu            sync.RWMutex
	capacity      int64 //maximum size of each file
	currentSize   int64 // current size of the latest file
	maxFiles      int   //maximum number of files
	notifyRotate  *pubsub.Publisher
	marshal       logger.MarshalFunc
	createDecoder makeDecoderFunc
***REMOVED***

type makeDecoderFunc func(rdr io.Reader) func() (*logger.Message, error)

//NewLogFile creates new LogFile
func NewLogFile(logPath string, capacity int64, maxFiles int, marshaller logger.MarshalFunc, decodeFunc makeDecoderFunc) (*LogFile, error) ***REMOVED***
	log, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	size, err := log.Seek(0, os.SEEK_END)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &LogFile***REMOVED***
		f:             log,
		capacity:      capacity,
		currentSize:   size,
		maxFiles:      maxFiles,
		notifyRotate:  pubsub.NewPublisher(0, 1),
		marshal:       marshaller,
		createDecoder: decodeFunc,
	***REMOVED***, nil
***REMOVED***

// WriteLogEntry writes the provided log message to the current log file.
// This may trigger a rotation event if the max file/capacity limits are hit.
func (w *LogFile) WriteLogEntry(msg *logger.Message) error ***REMOVED***
	b, err := w.marshal(msg)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error marshalling log message")
	***REMOVED***

	logger.PutMessage(msg)

	w.mu.Lock()
	if w.closed ***REMOVED***
		w.mu.Unlock()
		return errors.New("cannot write because the output file was closed")
	***REMOVED***

	if err := w.checkCapacityAndRotate(); err != nil ***REMOVED***
		w.mu.Unlock()
		return err
	***REMOVED***

	n, err := w.f.Write(b)
	if err == nil ***REMOVED***
		w.currentSize += int64(n)
	***REMOVED***
	w.mu.Unlock()
	return err
***REMOVED***

func (w *LogFile) checkCapacityAndRotate() error ***REMOVED***
	if w.capacity == -1 ***REMOVED***
		return nil
	***REMOVED***

	if w.currentSize >= w.capacity ***REMOVED***
		name := w.f.Name()
		if err := w.f.Close(); err != nil ***REMOVED***
			return errors.Wrap(err, "error closing file")
		***REMOVED***
		if err := rotate(name, w.maxFiles); err != nil ***REMOVED***
			return err
		***REMOVED***
		file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0640)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		w.f = file
		w.currentSize = 0
		w.notifyRotate.Publish(struct***REMOVED******REMOVED******REMOVED******REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

func rotate(name string, maxFiles int) error ***REMOVED***
	if maxFiles < 2 ***REMOVED***
		return nil
	***REMOVED***
	for i := maxFiles - 1; i > 1; i-- ***REMOVED***
		toPath := name + "." + strconv.Itoa(i)
		fromPath := name + "." + strconv.Itoa(i-1)
		if err := os.Rename(fromPath, toPath); err != nil && !os.IsNotExist(err) ***REMOVED***
			return errors.Wrap(err, "error rotating old log entries")
		***REMOVED***
	***REMOVED***

	if err := os.Rename(name, name+".1"); err != nil && !os.IsNotExist(err) ***REMOVED***
		return errors.Wrap(err, "error rotating current log")
	***REMOVED***
	return nil
***REMOVED***

// LogPath returns the location the given writer logs to.
func (w *LogFile) LogPath() string ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.f.Name()
***REMOVED***

// MaxFiles return maximum number of files
func (w *LogFile) MaxFiles() int ***REMOVED***
	return w.maxFiles
***REMOVED***

// Close closes underlying file and signals all readers to stop.
func (w *LogFile) Close() error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed ***REMOVED***
		return nil
	***REMOVED***
	if err := w.f.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.closed = true
	return nil
***REMOVED***

// ReadLogs decodes entries from log files and sends them the passed in watcher
func (w *LogFile) ReadLogs(config logger.ReadConfig, watcher *logger.LogWatcher) ***REMOVED***
	w.mu.RLock()
	files, err := w.openRotatedFiles()
	if err != nil ***REMOVED***
		w.mu.RUnlock()
		watcher.Err <- err
		return
	***REMOVED***
	defer func() ***REMOVED***
		for _, f := range files ***REMOVED***
			f.Close()
		***REMOVED***
	***REMOVED***()

	currentFile, err := os.Open(w.f.Name())
	if err != nil ***REMOVED***
		w.mu.RUnlock()
		watcher.Err <- err
		return
	***REMOVED***
	defer currentFile.Close()

	currentChunk, err := newSectionReader(currentFile)
	w.mu.RUnlock()

	if err != nil ***REMOVED***
		watcher.Err <- err
		return
	***REMOVED***

	if config.Tail != 0 ***REMOVED***
		seekers := make([]io.ReadSeeker, 0, len(files)+1)
		for _, f := range files ***REMOVED***
			seekers = append(seekers, f)
		***REMOVED***
		seekers = append(seekers, currentChunk)
		tailFile(multireader.MultiReadSeeker(seekers...), watcher, w.createDecoder, config)
	***REMOVED***

	w.mu.RLock()
	if !config.Follow || w.closed ***REMOVED***
		w.mu.RUnlock()
		return
	***REMOVED***
	w.mu.RUnlock()

	notifyRotate := w.notifyRotate.Subscribe()
	defer w.notifyRotate.Evict(notifyRotate)
	followLogs(currentFile, watcher, notifyRotate, w.createDecoder, config.Since, config.Until)
***REMOVED***

func (w *LogFile) openRotatedFiles() (files []*os.File, err error) ***REMOVED***
	defer func() ***REMOVED***
		if err == nil ***REMOVED***
			return
		***REMOVED***
		for _, f := range files ***REMOVED***
			f.Close()
		***REMOVED***
	***REMOVED***()

	for i := w.maxFiles; i > 1; i-- ***REMOVED***
		f, err := os.Open(fmt.Sprintf("%s.%d", w.f.Name(), i-1))
		if err != nil ***REMOVED***
			if !os.IsNotExist(err) ***REMOVED***
				return nil, err
			***REMOVED***
			continue
		***REMOVED***
		files = append(files, f)
	***REMOVED***

	return files, nil
***REMOVED***

func newSectionReader(f *os.File) (*io.SectionReader, error) ***REMOVED***
	// seek to the end to get the size
	// we'll leave this at the end of the file since section reader does not advance the reader
	size, err := f.Seek(0, os.SEEK_END)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error getting current file size")
	***REMOVED***
	return io.NewSectionReader(f, 0, size), nil
***REMOVED***

type decodeFunc func() (*logger.Message, error)

func tailFile(f io.ReadSeeker, watcher *logger.LogWatcher, createDecoder makeDecoderFunc, config logger.ReadConfig) ***REMOVED***
	var rdr io.Reader = f
	if config.Tail > 0 ***REMOVED***
		ls, err := tailfile.TailFile(f, config.Tail)
		if err != nil ***REMOVED***
			watcher.Err <- err
			return
		***REMOVED***
		rdr = bytes.NewBuffer(bytes.Join(ls, []byte("\n")))
	***REMOVED***

	decodeLogLine := createDecoder(rdr)
	for ***REMOVED***
		msg, err := decodeLogLine()
		if err != nil ***REMOVED***
			if err != io.EOF ***REMOVED***
				watcher.Err <- err
			***REMOVED***
			return
		***REMOVED***
		if !config.Since.IsZero() && msg.Timestamp.Before(config.Since) ***REMOVED***
			continue
		***REMOVED***
		if !config.Until.IsZero() && msg.Timestamp.After(config.Until) ***REMOVED***
			return
		***REMOVED***
		select ***REMOVED***
		case <-watcher.WatchClose():
			return
		case watcher.Msg <- msg:
		***REMOVED***
	***REMOVED***
***REMOVED***

func followLogs(f *os.File, logWatcher *logger.LogWatcher, notifyRotate chan interface***REMOVED******REMOVED***, createDecoder makeDecoderFunc, since, until time.Time) ***REMOVED***
	decodeLogLine := createDecoder(f)

	name := f.Name()
	fileWatcher, err := watchFile(name)
	if err != nil ***REMOVED***
		logWatcher.Err <- err
		return
	***REMOVED***
	defer func() ***REMOVED***
		f.Close()
		fileWatcher.Remove(name)
		fileWatcher.Close()
	***REMOVED***()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() ***REMOVED***
		select ***REMOVED***
		case <-logWatcher.WatchClose():
			fileWatcher.Remove(name)
			cancel()
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***()

	var retries int
	handleRotate := func() error ***REMOVED***
		f.Close()
		fileWatcher.Remove(name)

		// retry when the file doesn't exist
		for retries := 0; retries <= 5; retries++ ***REMOVED***
			f, err = os.Open(name)
			if err == nil || !os.IsNotExist(err) ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := fileWatcher.Add(name); err != nil ***REMOVED***
			return err
		***REMOVED***
		decodeLogLine = createDecoder(f)
		return nil
	***REMOVED***

	errRetry := errors.New("retry")
	errDone := errors.New("done")
	waitRead := func() error ***REMOVED***
		select ***REMOVED***
		case e := <-fileWatcher.Events():
			switch e.Op ***REMOVED***
			case fsnotify.Write:
				decodeLogLine = createDecoder(f)
				return nil
			case fsnotify.Rename, fsnotify.Remove:
				select ***REMOVED***
				case <-notifyRotate:
				case <-ctx.Done():
					return errDone
				***REMOVED***
				if err := handleRotate(); err != nil ***REMOVED***
					return err
				***REMOVED***
				return nil
			***REMOVED***
			return errRetry
		case err := <-fileWatcher.Errors():
			logrus.Debug("logger got error watching file: %v", err)
			// Something happened, let's try and stay alive and create a new watcher
			if retries <= 5 ***REMOVED***
				fileWatcher.Close()
				fileWatcher, err = watchFile(name)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				retries++
				return errRetry
			***REMOVED***
			return err
		case <-ctx.Done():
			return errDone
		***REMOVED***
	***REMOVED***

	handleDecodeErr := func(err error) error ***REMOVED***
		if err != io.EOF ***REMOVED***
			return err
		***REMOVED***

		for ***REMOVED***
			err := waitRead()
			if err == nil ***REMOVED***
				break
			***REMOVED***
			if err == errRetry ***REMOVED***
				continue
			***REMOVED***
			return err
		***REMOVED***
		return nil
	***REMOVED***

	// main loop
	for ***REMOVED***
		msg, err := decodeLogLine()
		if err != nil ***REMOVED***
			if err := handleDecodeErr(err); err != nil ***REMOVED***
				if err == errDone ***REMOVED***
					return
				***REMOVED***
				// we got an unrecoverable error, so return
				logWatcher.Err <- err
				return
			***REMOVED***
			// ready to try again
			continue
		***REMOVED***

		retries = 0 // reset retries since we've succeeded
		if !since.IsZero() && msg.Timestamp.Before(since) ***REMOVED***
			continue
		***REMOVED***
		if !until.IsZero() && msg.Timestamp.After(until) ***REMOVED***
			return
		***REMOVED***
		select ***REMOVED***
		case logWatcher.Msg <- msg:
		case <-ctx.Done():
			logWatcher.Msg <- msg
			for ***REMOVED***
				msg, err := decodeLogLine()
				if err != nil ***REMOVED***
					return
				***REMOVED***
				if !since.IsZero() && msg.Timestamp.Before(since) ***REMOVED***
					continue
				***REMOVED***
				if !until.IsZero() && msg.Timestamp.After(until) ***REMOVED***
					return
				***REMOVED***
				logWatcher.Msg <- msg
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func watchFile(name string) (filenotify.FileWatcher, error) ***REMOVED***
	fileWatcher, err := filenotify.New()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	logger := logrus.WithFields(logrus.Fields***REMOVED***
		"module": "logger",
		"fille":  name,
	***REMOVED***)

	if err := fileWatcher.Add(name); err != nil ***REMOVED***
		logger.WithError(err).Warnf("falling back to file poller")
		fileWatcher.Close()
		fileWatcher = filenotify.NewPollingWatcher()

		if err := fileWatcher.Add(name); err != nil ***REMOVED***
			fileWatcher.Close()
			logger.WithError(err).Debugf("error watching log file for modifications")
			return nil, err
		***REMOVED***
	***REMOVED***
	return fileWatcher, nil
***REMOVED***
