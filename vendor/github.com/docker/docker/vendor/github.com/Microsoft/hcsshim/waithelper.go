package hcsshim

import (
	"time"

	"github.com/sirupsen/logrus"
)

func processAsyncHcsResult(err error, resultp *uint16, callbackNumber uintptr, expectedNotification hcsNotification, timeout *time.Duration) error ***REMOVED***
	err = processHcsResult(err, resultp)
	if IsPending(err) ***REMOVED***
		return waitForNotification(callbackNumber, expectedNotification, timeout)
	***REMOVED***

	return err
***REMOVED***

func waitForNotification(callbackNumber uintptr, expectedNotification hcsNotification, timeout *time.Duration) error ***REMOVED***
	callbackMapLock.RLock()
	channels := callbackMap[callbackNumber].channels
	callbackMapLock.RUnlock()

	expectedChannel := channels[expectedNotification]
	if expectedChannel == nil ***REMOVED***
		logrus.Errorf("unknown notification type in waitForNotification %x", expectedNotification)
		return ErrInvalidNotificationType
	***REMOVED***

	var c <-chan time.Time
	if timeout != nil ***REMOVED***
		timer := time.NewTimer(*timeout)
		c = timer.C
		defer timer.Stop()
	***REMOVED***

	select ***REMOVED***
	case err, ok := <-expectedChannel:
		if !ok ***REMOVED***
			return ErrHandleClose
		***REMOVED***
		return err
	case err, ok := <-channels[hcsNotificationSystemExited]:
		if !ok ***REMOVED***
			return ErrHandleClose
		***REMOVED***
		// If the expected notification is hcsNotificationSystemExited which of the two selects
		// chosen is random. Return the raw error if hcsNotificationSystemExited is expected
		if channels[hcsNotificationSystemExited] == expectedChannel ***REMOVED***
			return err
		***REMOVED***
		return ErrUnexpectedContainerExit
	case _, ok := <-channels[hcsNotificationServiceDisconnect]:
		if !ok ***REMOVED***
			return ErrHandleClose
		***REMOVED***
		// hcsNotificationServiceDisconnect should never be an expected notification
		// it does not need the same handling as hcsNotificationSystemExited
		return ErrUnexpectedProcessAbort
	case <-c:
		return ErrTimeout
	***REMOVED***
	return nil
***REMOVED***
