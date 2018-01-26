package container

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
)

func TestIsValidHealthString(t *testing.T) ***REMOVED***
	contexts := []struct ***REMOVED***
		Health   string
		Expected bool
	***REMOVED******REMOVED***
		***REMOVED***types.Healthy, true***REMOVED***,
		***REMOVED***types.Unhealthy, true***REMOVED***,
		***REMOVED***types.Starting, true***REMOVED***,
		***REMOVED***types.NoHealthcheck, true***REMOVED***,
		***REMOVED***"fail", false***REMOVED***,
	***REMOVED***

	for _, c := range contexts ***REMOVED***
		v := IsValidHealthString(c.Health)
		if v != c.Expected ***REMOVED***
			t.Fatalf("Expected %t, but got %t", c.Expected, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestStateRunStop(t *testing.T) ***REMOVED***
	s := NewState()

	// Begin another wait with WaitConditionRemoved. It should complete
	// within 200 milliseconds.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	removalWait := s.Wait(ctx, WaitConditionRemoved)

	// Full lifecycle two times.
	for i := 1; i <= 2; i++ ***REMOVED***
		// A wait with WaitConditionNotRunning should return
		// immediately since the state is now either "created" (on the
		// first iteration) or "exited" (on the second iteration). It
		// shouldn't take more than 50 milliseconds.
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		// Expectx exit code to be i-1 since it should be the exit
		// code from the previous loop or 0 for the created state.
		if status := <-s.Wait(ctx, WaitConditionNotRunning); status.ExitCode() != i-1 ***REMOVED***
			t.Fatalf("ExitCode %v, expected %v, err %q", status.ExitCode(), i-1, status.Err())
		***REMOVED***

		// A wait with WaitConditionNextExit should block until the
		// container has started and exited. It shouldn't take more
		// than 100 milliseconds.
		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		initialWait := s.Wait(ctx, WaitConditionNextExit)

		// Set the state to "Running".
		s.Lock()
		s.SetRunning(i, true)
		s.Unlock()

		// Assert desired state.
		if !s.IsRunning() ***REMOVED***
			t.Fatal("State not running")
		***REMOVED***
		if s.Pid != i ***REMOVED***
			t.Fatalf("Pid %v, expected %v", s.Pid, i)
		***REMOVED***
		if s.ExitCode() != 0 ***REMOVED***
			t.Fatalf("ExitCode %v, expected 0", s.ExitCode())
		***REMOVED***

		// Now that it's running, a wait with WaitConditionNotRunning
		// should block until we stop the container. It shouldn't take
		// more than 100 milliseconds.
		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		exitWait := s.Wait(ctx, WaitConditionNotRunning)

		// Set the state to "Exited".
		s.Lock()
		s.SetStopped(&ExitStatus***REMOVED***ExitCode: i***REMOVED***)
		s.Unlock()

		// Assert desired state.
		if s.IsRunning() ***REMOVED***
			t.Fatal("State is running")
		***REMOVED***
		if s.ExitCode() != i ***REMOVED***
			t.Fatalf("ExitCode %v, expected %v", s.ExitCode(), i)
		***REMOVED***
		if s.Pid != 0 ***REMOVED***
			t.Fatalf("Pid %v, expected 0", s.Pid)
		***REMOVED***

		// Receive the initialWait result.
		if status := <-initialWait; status.ExitCode() != i ***REMOVED***
			t.Fatalf("ExitCode %v, expected %v, err %q", status.ExitCode(), i, status.Err())
		***REMOVED***

		// Receive the exitWait result.
		if status := <-exitWait; status.ExitCode() != i ***REMOVED***
			t.Fatalf("ExitCode %v, expected %v, err %q", status.ExitCode(), i, status.Err())
		***REMOVED***
	***REMOVED***

	// Set the state to dead and removed.
	s.SetDead()
	s.SetRemoved()

	// Wait for removed status or timeout.
	if status := <-removalWait; status.ExitCode() != 2 ***REMOVED***
		// Should have the final exit code from the loop.
		t.Fatalf("Removal wait exitCode %v, expected %v, err %q", status.ExitCode(), 2, status.Err())
	***REMOVED***
***REMOVED***

func TestStateTimeoutWait(t *testing.T) ***REMOVED***
	s := NewState()

	s.Lock()
	s.SetRunning(0, true)
	s.Unlock()

	// Start a wait with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	waitC := s.Wait(ctx, WaitConditionNotRunning)

	// It should timeout *before* this 200ms timer does.
	select ***REMOVED***
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Stop callback doesn't fire in 200 milliseconds")
	case status := <-waitC:
		t.Log("Stop callback fired")
		// Should be a timeout error.
		if status.Err() == nil ***REMOVED***
			t.Fatal("expected timeout error, got nil")
		***REMOVED***
		if status.ExitCode() != -1 ***REMOVED***
			t.Fatalf("expected exit code %v, got %v", -1, status.ExitCode())
		***REMOVED***
	***REMOVED***

	s.Lock()
	s.SetStopped(&ExitStatus***REMOVED***ExitCode: 0***REMOVED***)
	s.Unlock()

	// Start another wait with a timeout. This one should return
	// immediately.
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	waitC = s.Wait(ctx, WaitConditionNotRunning)

	select ***REMOVED***
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Stop callback doesn't fire in 200 milliseconds")
	case status := <-waitC:
		t.Log("Stop callback fired")
		if status.ExitCode() != 0 ***REMOVED***
			t.Fatalf("expected exit code %v, got %v, err %q", 0, status.ExitCode(), status.Err())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsValidStateString(t *testing.T) ***REMOVED***
	states := []struct ***REMOVED***
		state    string
		expected bool
	***REMOVED******REMOVED***
		***REMOVED***"paused", true***REMOVED***,
		***REMOVED***"restarting", true***REMOVED***,
		***REMOVED***"running", true***REMOVED***,
		***REMOVED***"dead", true***REMOVED***,
		***REMOVED***"start", false***REMOVED***,
		***REMOVED***"created", true***REMOVED***,
		***REMOVED***"exited", true***REMOVED***,
		***REMOVED***"removing", true***REMOVED***,
		***REMOVED***"stop", false***REMOVED***,
	***REMOVED***

	for _, s := range states ***REMOVED***
		v := IsValidStateString(s.state)
		if v != s.expected ***REMOVED***
			t.Fatalf("Expected %t, but got %t", s.expected, v)
		***REMOVED***
	***REMOVED***
***REMOVED***
