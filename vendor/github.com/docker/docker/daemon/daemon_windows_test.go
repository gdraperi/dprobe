// +build windows

package daemon

import (
	"strings"
	"testing"

	"golang.org/x/sys/windows/svc/mgr"
)

const existingService = "Power"

func TestEnsureServicesExist(t *testing.T) ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		t.Fatal("failed to connect to service manager, this test needs admin")
	***REMOVED***
	defer m.Disconnect()
	s, err := m.OpenService(existingService)
	if err != nil ***REMOVED***
		t.Fatalf("expected to find known inbox service %q, this test needs a known inbox service to run correctly", existingService)
	***REMOVED***
	defer s.Close()

	input := []string***REMOVED***existingService***REMOVED***
	err = ensureServicesInstalled(input)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error for input %q: %q", input, err)
	***REMOVED***
***REMOVED***

func TestEnsureServicesExistErrors(t *testing.T) ***REMOVED***
	m, err := mgr.Connect()
	if err != nil ***REMOVED***
		t.Fatal("failed to connect to service manager, this test needs admin")
	***REMOVED***
	defer m.Disconnect()
	s, err := m.OpenService(existingService)
	if err != nil ***REMOVED***
		t.Fatalf("expected to find known inbox service %q, this test needs a known inbox service to run correctly", existingService)
	***REMOVED***
	defer s.Close()

	for _, testcase := range []struct ***REMOVED***
		input         []string
		expectedError string
	***REMOVED******REMOVED***
		***REMOVED***
			input:         []string***REMOVED***"daemon_windows_test_fakeservice"***REMOVED***,
			expectedError: "failed to open service daemon_windows_test_fakeservice",
		***REMOVED***,
		***REMOVED***
			input:         []string***REMOVED***"daemon_windows_test_fakeservice1", "daemon_windows_test_fakeservice2"***REMOVED***,
			expectedError: "failed to open service daemon_windows_test_fakeservice1",
		***REMOVED***,
		***REMOVED***
			input:         []string***REMOVED***existingService, "daemon_windows_test_fakeservice"***REMOVED***,
			expectedError: "failed to open service daemon_windows_test_fakeservice",
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		t.Run(strings.Join(testcase.input, ";"), func(t *testing.T) ***REMOVED***
			err := ensureServicesInstalled(testcase.input)
			if err == nil ***REMOVED***
				t.Fatalf("expected error for input %v", testcase.input)
			***REMOVED***
			if !strings.Contains(err.Error(), testcase.expectedError) ***REMOVED***
				t.Fatalf("expected error %q to contain %q", err.Error(), testcase.expectedError)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
