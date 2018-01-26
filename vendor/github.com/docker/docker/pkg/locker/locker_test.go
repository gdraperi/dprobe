package locker

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLockCounter(t *testing.T) ***REMOVED***
	l := &lockCtr***REMOVED******REMOVED***
	l.inc()

	if l.waiters != 1 ***REMOVED***
		t.Fatal("counter inc failed")
	***REMOVED***

	l.dec()
	if l.waiters != 0 ***REMOVED***
		t.Fatal("counter dec failed")
	***REMOVED***
***REMOVED***

func TestLockerLock(t *testing.T) ***REMOVED***
	l := New()
	l.Lock("test")
	ctr := l.locks["test"]

	if ctr.count() != 0 ***REMOVED***
		t.Fatalf("expected waiters to be 0, got :%d", ctr.waiters)
	***REMOVED***

	chDone := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		l.Lock("test")
		close(chDone)
	***REMOVED***()

	chWaiting := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		for range time.Tick(1 * time.Millisecond) ***REMOVED***
			if ctr.count() == 1 ***REMOVED***
				close(chWaiting)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case <-chWaiting:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for lock waiters to be incremented")
	***REMOVED***

	select ***REMOVED***
	case <-chDone:
		t.Fatal("lock should not have returned while it was still held")
	default:
	***REMOVED***

	if err := l.Unlock("test"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	select ***REMOVED***
	case <-chDone:
	case <-time.After(3 * time.Second):
		t.Fatalf("lock should have completed")
	***REMOVED***

	if ctr.count() != 0 ***REMOVED***
		t.Fatalf("expected waiters to be 0, got: %d", ctr.count())
	***REMOVED***
***REMOVED***

func TestLockerUnlock(t *testing.T) ***REMOVED***
	l := New()

	l.Lock("test")
	l.Unlock("test")

	chDone := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		l.Lock("test")
		close(chDone)
	***REMOVED***()

	select ***REMOVED***
	case <-chDone:
	case <-time.After(3 * time.Second):
		t.Fatalf("lock should not be blocked")
	***REMOVED***
***REMOVED***

func TestLockerConcurrency(t *testing.T) ***REMOVED***
	l := New()

	var wg sync.WaitGroup
	for i := 0; i <= 10000; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			l.Lock("test")
			// if there is a concurrency issue, will very likely panic here
			l.Unlock("test")
			wg.Done()
		***REMOVED***()
	***REMOVED***

	chDone := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		wg.Wait()
		close(chDone)
	***REMOVED***()

	select ***REMOVED***
	case <-chDone:
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for locks to complete")
	***REMOVED***

	// Since everything has unlocked this should not exist anymore
	if ctr, exists := l.locks["test"]; exists ***REMOVED***
		t.Fatalf("lock should not exist: %v", ctr)
	***REMOVED***
***REMOVED***

func BenchmarkLocker(b *testing.B) ***REMOVED***
	l := New()
	for i := 0; i < b.N; i++ ***REMOVED***
		l.Lock("test")
		l.Unlock("test")
	***REMOVED***
***REMOVED***

func BenchmarkLockerParallel(b *testing.B) ***REMOVED***
	l := New()
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		for pb.Next() ***REMOVED***
			l.Lock("test")
			l.Unlock("test")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkLockerMoreKeys(b *testing.B) ***REMOVED***
	l := New()
	var keys []string
	for i := 0; i < 64; i++ ***REMOVED***
		keys = append(keys, strconv.Itoa(i))
	***REMOVED***
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		for pb.Next() ***REMOVED***
			k := keys[rand.Intn(len(keys))]
			l.Lock(k)
			l.Unlock(k)
		***REMOVED***
	***REMOVED***)
***REMOVED***
