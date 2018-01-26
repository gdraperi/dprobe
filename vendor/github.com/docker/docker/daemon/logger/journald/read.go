// +build linux,cgo,!static_build,journald

package journald

// #include <sys/types.h>
// #include <sys/poll.h>
// #include <systemd/sd-journal.h>
// #include <errno.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include <time.h>
// #include <unistd.h>
//
//static int get_message(sd_journal *j, const char **msg, size_t *length, int *partial)
//***REMOVED***
//	int rc;
//	size_t plength;
//	*msg = NULL;
//	*length = 0;
//	plength = strlen("CONTAINER_PARTIAL_MESSAGE=true");
//	rc = sd_journal_get_data(j, "CONTAINER_PARTIAL_MESSAGE", (const void **) msg, length);
//	*partial = ((rc == 0) && (*length == plength) && (memcmp(*msg, "CONTAINER_PARTIAL_MESSAGE=true", plength) == 0));
//	rc = sd_journal_get_data(j, "MESSAGE", (const void **) msg, length);
//	if (rc == 0) ***REMOVED***
//		if (*length > 8) ***REMOVED***
//			(*msg) += 8;
//			*length -= 8;
//		***REMOVED*** else ***REMOVED***
//			*msg = NULL;
//			*length = 0;
//			rc = -ENOENT;
//		***REMOVED***
//	***REMOVED***
//	return rc;
//***REMOVED***
//static int get_priority(sd_journal *j, int *priority)
//***REMOVED***
//	const void *data;
//	size_t i, length;
//	int rc;
//	*priority = -1;
//	rc = sd_journal_get_data(j, "PRIORITY", &data, &length);
//	if (rc == 0) ***REMOVED***
//		if ((length > 9) && (strncmp(data, "PRIORITY=", 9) == 0)) ***REMOVED***
//			*priority = 0;
//			for (i = 9; i < length; i++) ***REMOVED***
//				*priority = *priority * 10 + ((const char *)data)[i] - '0';
//			***REMOVED***
//			if (length > 9) ***REMOVED***
//				rc = 0;
//			***REMOVED***
//		***REMOVED***
//	***REMOVED***
//	return rc;
//***REMOVED***
//static int is_attribute_field(const char *msg, size_t length)
//***REMOVED***
//	static const struct known_field ***REMOVED***
//		const char *name;
//		size_t length;
//	***REMOVED*** fields[] = ***REMOVED***
//		***REMOVED***"MESSAGE", sizeof("MESSAGE") - 1***REMOVED***,
//		***REMOVED***"MESSAGE_ID", sizeof("MESSAGE_ID") - 1***REMOVED***,
//		***REMOVED***"PRIORITY", sizeof("PRIORITY") - 1***REMOVED***,
//		***REMOVED***"CODE_FILE", sizeof("CODE_FILE") - 1***REMOVED***,
//		***REMOVED***"CODE_LINE", sizeof("CODE_LINE") - 1***REMOVED***,
//		***REMOVED***"CODE_FUNC", sizeof("CODE_FUNC") - 1***REMOVED***,
//		***REMOVED***"ERRNO", sizeof("ERRNO") - 1***REMOVED***,
//		***REMOVED***"SYSLOG_FACILITY", sizeof("SYSLOG_FACILITY") - 1***REMOVED***,
//		***REMOVED***"SYSLOG_IDENTIFIER", sizeof("SYSLOG_IDENTIFIER") - 1***REMOVED***,
//		***REMOVED***"SYSLOG_PID", sizeof("SYSLOG_PID") - 1***REMOVED***,
//		***REMOVED***"CONTAINER_NAME", sizeof("CONTAINER_NAME") - 1***REMOVED***,
//		***REMOVED***"CONTAINER_ID", sizeof("CONTAINER_ID") - 1***REMOVED***,
//		***REMOVED***"CONTAINER_ID_FULL", sizeof("CONTAINER_ID_FULL") - 1***REMOVED***,
//		***REMOVED***"CONTAINER_TAG", sizeof("CONTAINER_TAG") - 1***REMOVED***,
//	***REMOVED***;
//	unsigned int i;
//	void *p;
//	if ((length < 1) || (msg[0] == '_') || ((p = memchr(msg, '=', length)) == NULL)) ***REMOVED***
//		return -1;
//	***REMOVED***
//	length = ((const char *) p) - msg;
//	for (i = 0; i < sizeof(fields) / sizeof(fields[0]); i++) ***REMOVED***
//		if ((fields[i].length == length) && (memcmp(fields[i].name, msg, length) == 0)) ***REMOVED***
//			return -1;
//		***REMOVED***
//	***REMOVED***
//	return 0;
//***REMOVED***
//static int get_attribute_field(sd_journal *j, const char **msg, size_t *length)
//***REMOVED***
//	int rc;
//	*msg = NULL;
//	*length = 0;
//	while ((rc = sd_journal_enumerate_data(j, (const void **) msg, length)) > 0) ***REMOVED***
//		if (is_attribute_field(*msg, *length) == 0) ***REMOVED***
//			break;
//		***REMOVED***
//		rc = -ENOENT;
//	***REMOVED***
//	return rc;
//***REMOVED***
//static int wait_for_data_cancelable(sd_journal *j, int pipefd)
//***REMOVED***
//	struct pollfd fds[2];
//	uint64_t when = 0;
//	int timeout, jevents, i;
//	struct timespec ts;
//	uint64_t now;
//
//	memset(&fds, 0, sizeof(fds));
//	fds[0].fd = pipefd;
//	fds[0].events = POLLHUP;
//	fds[1].fd = sd_journal_get_fd(j);
//	if (fds[1].fd < 0) ***REMOVED***
//		return fds[1].fd;
//	***REMOVED***
//
//	do ***REMOVED***
//		jevents = sd_journal_get_events(j);
//		if (jevents < 0) ***REMOVED***
//			return jevents;
//		***REMOVED***
//		fds[1].events = jevents;
//		sd_journal_get_timeout(j, &when);
//		if (when == -1) ***REMOVED***
//			timeout = -1;
//		***REMOVED*** else ***REMOVED***
//			clock_gettime(CLOCK_MONOTONIC, &ts);
//			now = (uint64_t) ts.tv_sec * 1000000 + ts.tv_nsec / 1000;
//			timeout = when > now ? (int) ((when - now + 999) / 1000) : 0;
//		***REMOVED***
//		i = poll(fds, 2, timeout);
//		if ((i == -1) && (errno != EINTR)) ***REMOVED***
//			/* An unexpected error. */
//			return (errno != 0) ? -errno : -EINTR;
//		***REMOVED***
//		if (fds[0].revents & POLLHUP) ***REMOVED***
//			/* The close notification pipe was closed. */
//			return 0;
//		***REMOVED***
//		if (sd_journal_process(j) == SD_JOURNAL_APPEND) ***REMOVED***
//			/* Data, which we might care about, was appended. */
//			return 1;
//		***REMOVED***
//	***REMOVED*** while ((fds[0].revents & POLLHUP) == 0);
//	return 0;
//***REMOVED***
import "C"

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/coreos/go-systemd/journal"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/daemon/logger"
	"github.com/sirupsen/logrus"
)

func (s *journald) Close() error ***REMOVED***
	s.mu.Lock()
	s.closed = true
	for reader := range s.readers.readers ***REMOVED***
		reader.Close()
	***REMOVED***
	s.mu.Unlock()
	return nil
***REMOVED***

func (s *journald) drainJournal(logWatcher *logger.LogWatcher, j *C.sd_journal, oldCursor *C.char, untilUnixMicro uint64) (*C.char, bool) ***REMOVED***
	var msg, data, cursor *C.char
	var length C.size_t
	var stamp C.uint64_t
	var priority, partial C.int
	var done bool

	// Walk the journal from here forward until we run out of new entries
	// or we reach the until value (if provided).
drain:
	for ***REMOVED***
		// Try not to send a given entry twice.
		if oldCursor != nil ***REMOVED***
			for C.sd_journal_test_cursor(j, oldCursor) > 0 ***REMOVED***
				if C.sd_journal_next(j) <= 0 ***REMOVED***
					break drain
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// Read and send the logged message, if there is one to read.
		i := C.get_message(j, &msg, &length, &partial)
		if i != -C.ENOENT && i != -C.EADDRNOTAVAIL ***REMOVED***
			// Read the entry's timestamp.
			if C.sd_journal_get_realtime_usec(j, &stamp) != 0 ***REMOVED***
				break
			***REMOVED***
			// Break if the timestamp exceeds any provided until flag.
			if untilUnixMicro != 0 && untilUnixMicro < uint64(stamp) ***REMOVED***
				done = true
				break
			***REMOVED***

			// Set up the time and text of the entry.
			timestamp := time.Unix(int64(stamp)/1000000, (int64(stamp)%1000000)*1000)
			line := C.GoBytes(unsafe.Pointer(msg), C.int(length))
			if partial == 0 ***REMOVED***
				line = append(line, "\n"...)
			***REMOVED***
			// Recover the stream name by mapping
			// from the journal priority back to
			// the stream that we would have
			// assigned that value.
			source := ""
			if C.get_priority(j, &priority) != 0 ***REMOVED***
				source = ""
			***REMOVED*** else if priority == C.int(journal.PriErr) ***REMOVED***
				source = "stderr"
			***REMOVED*** else if priority == C.int(journal.PriInfo) ***REMOVED***
				source = "stdout"
			***REMOVED***
			// Retrieve the values of any variables we're adding to the journal.
			var attrs []backend.LogAttr
			C.sd_journal_restart_data(j)
			for C.get_attribute_field(j, &data, &length) > C.int(0) ***REMOVED***
				kv := strings.SplitN(C.GoStringN(data, C.int(length)), "=", 2)
				attrs = append(attrs, backend.LogAttr***REMOVED***Key: kv[0], Value: kv[1]***REMOVED***)
			***REMOVED***
			// Send the log message.
			logWatcher.Msg <- &logger.Message***REMOVED***
				Line:      line,
				Source:    source,
				Timestamp: timestamp.In(time.UTC),
				Attrs:     attrs,
			***REMOVED***
		***REMOVED***
		// If we're at the end of the journal, we're done (for now).
		if C.sd_journal_next(j) <= 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	// free(NULL) is safe
	C.free(unsafe.Pointer(oldCursor))
	if C.sd_journal_get_cursor(j, &cursor) != 0 ***REMOVED***
		// ensure that we won't be freeing an address that's invalid
		cursor = nil
	***REMOVED***
	return cursor, done
***REMOVED***

func (s *journald) followJournal(logWatcher *logger.LogWatcher, j *C.sd_journal, pfd [2]C.int, cursor *C.char, untilUnixMicro uint64) *C.char ***REMOVED***
	s.mu.Lock()
	s.readers.readers[logWatcher] = logWatcher
	if s.closed ***REMOVED***
		// the journald Logger is closed, presumably because the container has been
		// reset.  So we shouldn't follow, because we'll never be woken up.  But we
		// should make one more drainJournal call to be sure we've got all the logs.
		// Close pfd[1] so that one drainJournal happens, then cleanup, then return.
		C.close(pfd[1])
	***REMOVED***
	s.mu.Unlock()

	newCursor := make(chan *C.char)

	go func() ***REMOVED***
		for ***REMOVED***
			// Keep copying journal data out until we're notified to stop
			// or we hit an error.
			status := C.wait_for_data_cancelable(j, pfd[0])
			if status < 0 ***REMOVED***
				cerrstr := C.strerror(C.int(-status))
				errstr := C.GoString(cerrstr)
				fmtstr := "error %q while attempting to follow journal for container %q"
				logrus.Errorf(fmtstr, errstr, s.vars["CONTAINER_ID_FULL"])
				break
			***REMOVED***

			var done bool
			cursor, done = s.drainJournal(logWatcher, j, cursor, untilUnixMicro)

			if status != 1 || done ***REMOVED***
				// We were notified to stop
				break
			***REMOVED***
		***REMOVED***

		// Clean up.
		C.close(pfd[0])
		s.mu.Lock()
		delete(s.readers.readers, logWatcher)
		s.mu.Unlock()
		close(logWatcher.Msg)
		newCursor <- cursor
	***REMOVED***()

	// Wait until we're told to stop.
	select ***REMOVED***
	case cursor = <-newCursor:
	case <-logWatcher.WatchClose():
		// Notify the other goroutine that its work is done.
		C.close(pfd[1])
		cursor = <-newCursor
	***REMOVED***

	return cursor
***REMOVED***

func (s *journald) readLogs(logWatcher *logger.LogWatcher, config logger.ReadConfig) ***REMOVED***
	var j *C.sd_journal
	var cmatch, cursor *C.char
	var stamp C.uint64_t
	var sinceUnixMicro uint64
	var untilUnixMicro uint64
	var pipes [2]C.int

	// Get a handle to the journal.
	rc := C.sd_journal_open(&j, C.int(0))
	if rc != 0 ***REMOVED***
		logWatcher.Err <- fmt.Errorf("error opening journal")
		close(logWatcher.Msg)
		return
	***REMOVED***
	// If we end up following the log, we can set the journal context
	// pointer and the channel pointer to nil so that we won't close them
	// here, potentially while the goroutine that uses them is still
	// running.  Otherwise, close them when we return from this function.
	following := false
	defer func(pfollowing *bool) ***REMOVED***
		if !*pfollowing ***REMOVED***
			close(logWatcher.Msg)
		***REMOVED***
		C.sd_journal_close(j)
	***REMOVED***(&following)
	// Remove limits on the size of data items that we'll retrieve.
	rc = C.sd_journal_set_data_threshold(j, C.size_t(0))
	if rc != 0 ***REMOVED***
		logWatcher.Err <- fmt.Errorf("error setting journal data threshold")
		return
	***REMOVED***
	// Add a match to have the library do the searching for us.
	cmatch = C.CString("CONTAINER_ID_FULL=" + s.vars["CONTAINER_ID_FULL"])
	defer C.free(unsafe.Pointer(cmatch))
	rc = C.sd_journal_add_match(j, unsafe.Pointer(cmatch), C.strlen(cmatch))
	if rc != 0 ***REMOVED***
		logWatcher.Err <- fmt.Errorf("error setting journal match")
		return
	***REMOVED***
	// If we have a cutoff time, convert it to Unix time once.
	if !config.Since.IsZero() ***REMOVED***
		nano := config.Since.UnixNano()
		sinceUnixMicro = uint64(nano / 1000)
	***REMOVED***
	// If we have an until value, convert it too
	if !config.Until.IsZero() ***REMOVED***
		nano := config.Until.UnixNano()
		untilUnixMicro = uint64(nano / 1000)
	***REMOVED***
	if config.Tail > 0 ***REMOVED***
		lines := config.Tail
		// If until time provided, start from there.
		// Otherwise start at the end of the journal.
		if untilUnixMicro != 0 && C.sd_journal_seek_realtime_usec(j, C.uint64_t(untilUnixMicro)) < 0 ***REMOVED***
			logWatcher.Err <- fmt.Errorf("error seeking provided until value")
			return
		***REMOVED*** else if C.sd_journal_seek_tail(j) < 0 ***REMOVED***
			logWatcher.Err <- fmt.Errorf("error seeking to end of journal")
			return
		***REMOVED***
		if C.sd_journal_previous(j) < 0 ***REMOVED***
			logWatcher.Err <- fmt.Errorf("error backtracking to previous journal entry")
			return
		***REMOVED***
		// Walk backward.
		for lines > 0 ***REMOVED***
			// Stop if the entry time is before our cutoff.
			// We'll need the entry time if it isn't, so go
			// ahead and parse it now.
			if C.sd_journal_get_realtime_usec(j, &stamp) != 0 ***REMOVED***
				break
			***REMOVED*** else ***REMOVED***
				// Compare the timestamp on the entry to our threshold value.
				if sinceUnixMicro != 0 && sinceUnixMicro > uint64(stamp) ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			lines--
			// If we're at the start of the journal, or
			// don't need to back up past any more entries,
			// stop.
			if lines == 0 || C.sd_journal_previous(j) <= 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Start at the beginning of the journal.
		if C.sd_journal_seek_head(j) < 0 ***REMOVED***
			logWatcher.Err <- fmt.Errorf("error seeking to start of journal")
			return
		***REMOVED***
		// If we have a cutoff date, fast-forward to it.
		if sinceUnixMicro != 0 && C.sd_journal_seek_realtime_usec(j, C.uint64_t(sinceUnixMicro)) != 0 ***REMOVED***
			logWatcher.Err <- fmt.Errorf("error seeking to start time in journal")
			return
		***REMOVED***
		if C.sd_journal_next(j) < 0 ***REMOVED***
			logWatcher.Err <- fmt.Errorf("error skipping to next journal entry")
			return
		***REMOVED***
	***REMOVED***
	cursor, _ = s.drainJournal(logWatcher, j, nil, untilUnixMicro)
	if config.Follow ***REMOVED***
		// Allocate a descriptor for following the journal, if we'll
		// need one.  Do it here so that we can report if it fails.
		if fd := C.sd_journal_get_fd(j); fd < C.int(0) ***REMOVED***
			logWatcher.Err <- fmt.Errorf("error opening journald follow descriptor: %q", C.GoString(C.strerror(-fd)))
		***REMOVED*** else ***REMOVED***
			// Create a pipe that we can poll at the same time as
			// the journald descriptor.
			if C.pipe(&pipes[0]) == C.int(-1) ***REMOVED***
				logWatcher.Err <- fmt.Errorf("error opening journald close notification pipe")
			***REMOVED*** else ***REMOVED***
				cursor = s.followJournal(logWatcher, j, pipes, cursor, untilUnixMicro)
				// Let followJournal handle freeing the journal context
				// object and closing the channel.
				following = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	C.free(unsafe.Pointer(cursor))
	return
***REMOVED***

func (s *journald) ReadLogs(config logger.ReadConfig) *logger.LogWatcher ***REMOVED***
	logWatcher := logger.NewLogWatcher()
	go s.readLogs(logWatcher, config)
	return logWatcher
***REMOVED***
