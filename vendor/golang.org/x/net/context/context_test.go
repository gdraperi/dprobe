// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.7

package context

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// otherContext is a Context that's not one of the types defined in context.go.
// This lets us test code paths that differ based on the underlying type of the
// Context.
type otherContext struct ***REMOVED***
	Context
***REMOVED***

func TestBackground(t *testing.T) ***REMOVED***
	c := Background()
	if c == nil ***REMOVED***
		t.Fatalf("Background returned nil")
	***REMOVED***
	select ***REMOVED***
	case x := <-c.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	***REMOVED***
	if got, want := fmt.Sprint(c), "context.Background"; got != want ***REMOVED***
		t.Errorf("Background().String() = %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestTODO(t *testing.T) ***REMOVED***
	c := TODO()
	if c == nil ***REMOVED***
		t.Fatalf("TODO returned nil")
	***REMOVED***
	select ***REMOVED***
	case x := <-c.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	***REMOVED***
	if got, want := fmt.Sprint(c), "context.TODO"; got != want ***REMOVED***
		t.Errorf("TODO().String() = %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestWithCancel(t *testing.T) ***REMOVED***
	c1, cancel := WithCancel(Background())

	if got, want := fmt.Sprint(c1), "context.Background.WithCancel"; got != want ***REMOVED***
		t.Errorf("c1.String() = %q want %q", got, want)
	***REMOVED***

	o := otherContext***REMOVED***c1***REMOVED***
	c2, _ := WithCancel(o)
	contexts := []Context***REMOVED***c1, o, c2***REMOVED***

	for i, c := range contexts ***REMOVED***
		if d := c.Done(); d == nil ***REMOVED***
			t.Errorf("c[%d].Done() == %v want non-nil", i, d)
		***REMOVED***
		if e := c.Err(); e != nil ***REMOVED***
			t.Errorf("c[%d].Err() == %v want nil", i, e)
		***REMOVED***

		select ***REMOVED***
		case x := <-c.Done():
			t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
		default:
		***REMOVED***
	***REMOVED***

	cancel()
	time.Sleep(100 * time.Millisecond) // let cancelation propagate

	for i, c := range contexts ***REMOVED***
		select ***REMOVED***
		case <-c.Done():
		default:
			t.Errorf("<-c[%d].Done() blocked, but shouldn't have", i)
		***REMOVED***
		if e := c.Err(); e != Canceled ***REMOVED***
			t.Errorf("c[%d].Err() == %v want %v", i, e, Canceled)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParentFinishesChild(t *testing.T) ***REMOVED***
	// Context tree:
	// parent -> cancelChild
	// parent -> valueChild -> timerChild
	parent, cancel := WithCancel(Background())
	cancelChild, stop := WithCancel(parent)
	defer stop()
	valueChild := WithValue(parent, "key", "value")
	timerChild, stop := WithTimeout(valueChild, 10000*time.Hour)
	defer stop()

	select ***REMOVED***
	case x := <-parent.Done():
		t.Errorf("<-parent.Done() == %v want nothing (it should block)", x)
	case x := <-cancelChild.Done():
		t.Errorf("<-cancelChild.Done() == %v want nothing (it should block)", x)
	case x := <-timerChild.Done():
		t.Errorf("<-timerChild.Done() == %v want nothing (it should block)", x)
	case x := <-valueChild.Done():
		t.Errorf("<-valueChild.Done() == %v want nothing (it should block)", x)
	default:
	***REMOVED***

	// The parent's children should contain the two cancelable children.
	pc := parent.(*cancelCtx)
	cc := cancelChild.(*cancelCtx)
	tc := timerChild.(*timerCtx)
	pc.mu.Lock()
	if len(pc.children) != 2 || !pc.children[cc] || !pc.children[tc] ***REMOVED***
		t.Errorf("bad linkage: pc.children = %v, want %v and %v",
			pc.children, cc, tc)
	***REMOVED***
	pc.mu.Unlock()

	if p, ok := parentCancelCtx(cc.Context); !ok || p != pc ***REMOVED***
		t.Errorf("bad linkage: parentCancelCtx(cancelChild.Context) = %v, %v want %v, true", p, ok, pc)
	***REMOVED***
	if p, ok := parentCancelCtx(tc.Context); !ok || p != pc ***REMOVED***
		t.Errorf("bad linkage: parentCancelCtx(timerChild.Context) = %v, %v want %v, true", p, ok, pc)
	***REMOVED***

	cancel()

	pc.mu.Lock()
	if len(pc.children) != 0 ***REMOVED***
		t.Errorf("pc.cancel didn't clear pc.children = %v", pc.children)
	***REMOVED***
	pc.mu.Unlock()

	// parent and children should all be finished.
	check := func(ctx Context, name string) ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		default:
			t.Errorf("<-%s.Done() blocked, but shouldn't have", name)
		***REMOVED***
		if e := ctx.Err(); e != Canceled ***REMOVED***
			t.Errorf("%s.Err() == %v want %v", name, e, Canceled)
		***REMOVED***
	***REMOVED***
	check(parent, "parent")
	check(cancelChild, "cancelChild")
	check(valueChild, "valueChild")
	check(timerChild, "timerChild")

	// WithCancel should return a canceled context on a canceled parent.
	precanceledChild := WithValue(parent, "key", "value")
	select ***REMOVED***
	case <-precanceledChild.Done():
	default:
		t.Errorf("<-precanceledChild.Done() blocked, but shouldn't have")
	***REMOVED***
	if e := precanceledChild.Err(); e != Canceled ***REMOVED***
		t.Errorf("precanceledChild.Err() == %v want %v", e, Canceled)
	***REMOVED***
***REMOVED***

func TestChildFinishesFirst(t *testing.T) ***REMOVED***
	cancelable, stop := WithCancel(Background())
	defer stop()
	for _, parent := range []Context***REMOVED***Background(), cancelable***REMOVED*** ***REMOVED***
		child, cancel := WithCancel(parent)

		select ***REMOVED***
		case x := <-parent.Done():
			t.Errorf("<-parent.Done() == %v want nothing (it should block)", x)
		case x := <-child.Done():
			t.Errorf("<-child.Done() == %v want nothing (it should block)", x)
		default:
		***REMOVED***

		cc := child.(*cancelCtx)
		pc, pcok := parent.(*cancelCtx) // pcok == false when parent == Background()
		if p, ok := parentCancelCtx(cc.Context); ok != pcok || (ok && pc != p) ***REMOVED***
			t.Errorf("bad linkage: parentCancelCtx(cc.Context) = %v, %v want %v, %v", p, ok, pc, pcok)
		***REMOVED***

		if pcok ***REMOVED***
			pc.mu.Lock()
			if len(pc.children) != 1 || !pc.children[cc] ***REMOVED***
				t.Errorf("bad linkage: pc.children = %v, cc = %v", pc.children, cc)
			***REMOVED***
			pc.mu.Unlock()
		***REMOVED***

		cancel()

		if pcok ***REMOVED***
			pc.mu.Lock()
			if len(pc.children) != 0 ***REMOVED***
				t.Errorf("child's cancel didn't remove self from pc.children = %v", pc.children)
			***REMOVED***
			pc.mu.Unlock()
		***REMOVED***

		// child should be finished.
		select ***REMOVED***
		case <-child.Done():
		default:
			t.Errorf("<-child.Done() blocked, but shouldn't have")
		***REMOVED***
		if e := child.Err(); e != Canceled ***REMOVED***
			t.Errorf("child.Err() == %v want %v", e, Canceled)
		***REMOVED***

		// parent should not be finished.
		select ***REMOVED***
		case x := <-parent.Done():
			t.Errorf("<-parent.Done() == %v want nothing (it should block)", x)
		default:
		***REMOVED***
		if e := parent.Err(); e != nil ***REMOVED***
			t.Errorf("parent.Err() == %v want nil", e)
		***REMOVED***
	***REMOVED***
***REMOVED***

func testDeadline(c Context, wait time.Duration, t *testing.T) ***REMOVED***
	select ***REMOVED***
	case <-time.After(wait):
		t.Fatalf("context should have timed out")
	case <-c.Done():
	***REMOVED***
	if e := c.Err(); e != DeadlineExceeded ***REMOVED***
		t.Errorf("c.Err() == %v want %v", e, DeadlineExceeded)
	***REMOVED***
***REMOVED***

func TestDeadline(t *testing.T) ***REMOVED***
	t.Parallel()
	const timeUnit = 500 * time.Millisecond
	c, _ := WithDeadline(Background(), time.Now().Add(1*timeUnit))
	if got, prefix := fmt.Sprint(c), "context.Background.WithDeadline("; !strings.HasPrefix(got, prefix) ***REMOVED***
		t.Errorf("c.String() = %q want prefix %q", got, prefix)
	***REMOVED***
	testDeadline(c, 2*timeUnit, t)

	c, _ = WithDeadline(Background(), time.Now().Add(1*timeUnit))
	o := otherContext***REMOVED***c***REMOVED***
	testDeadline(o, 2*timeUnit, t)

	c, _ = WithDeadline(Background(), time.Now().Add(1*timeUnit))
	o = otherContext***REMOVED***c***REMOVED***
	c, _ = WithDeadline(o, time.Now().Add(3*timeUnit))
	testDeadline(c, 2*timeUnit, t)
***REMOVED***

func TestTimeout(t *testing.T) ***REMOVED***
	t.Parallel()
	const timeUnit = 500 * time.Millisecond
	c, _ := WithTimeout(Background(), 1*timeUnit)
	if got, prefix := fmt.Sprint(c), "context.Background.WithDeadline("; !strings.HasPrefix(got, prefix) ***REMOVED***
		t.Errorf("c.String() = %q want prefix %q", got, prefix)
	***REMOVED***
	testDeadline(c, 2*timeUnit, t)

	c, _ = WithTimeout(Background(), 1*timeUnit)
	o := otherContext***REMOVED***c***REMOVED***
	testDeadline(o, 2*timeUnit, t)

	c, _ = WithTimeout(Background(), 1*timeUnit)
	o = otherContext***REMOVED***c***REMOVED***
	c, _ = WithTimeout(o, 3*timeUnit)
	testDeadline(c, 2*timeUnit, t)
***REMOVED***

func TestCanceledTimeout(t *testing.T) ***REMOVED***
	t.Parallel()
	const timeUnit = 500 * time.Millisecond
	c, _ := WithTimeout(Background(), 2*timeUnit)
	o := otherContext***REMOVED***c***REMOVED***
	c, cancel := WithTimeout(o, 4*timeUnit)
	cancel()
	time.Sleep(1 * timeUnit) // let cancelation propagate
	select ***REMOVED***
	case <-c.Done():
	default:
		t.Errorf("<-c.Done() blocked, but shouldn't have")
	***REMOVED***
	if e := c.Err(); e != Canceled ***REMOVED***
		t.Errorf("c.Err() == %v want %v", e, Canceled)
	***REMOVED***
***REMOVED***

type key1 int
type key2 int

var k1 = key1(1)
var k2 = key2(1) // same int as k1, different type
var k3 = key2(3) // same type as k2, different int

func TestValues(t *testing.T) ***REMOVED***
	check := func(c Context, nm, v1, v2, v3 string) ***REMOVED***
		if v, ok := c.Value(k1).(string); ok == (len(v1) == 0) || v != v1 ***REMOVED***
			t.Errorf(`%s.Value(k1).(string) = %q, %t want %q, %t`, nm, v, ok, v1, len(v1) != 0)
		***REMOVED***
		if v, ok := c.Value(k2).(string); ok == (len(v2) == 0) || v != v2 ***REMOVED***
			t.Errorf(`%s.Value(k2).(string) = %q, %t want %q, %t`, nm, v, ok, v2, len(v2) != 0)
		***REMOVED***
		if v, ok := c.Value(k3).(string); ok == (len(v3) == 0) || v != v3 ***REMOVED***
			t.Errorf(`%s.Value(k3).(string) = %q, %t want %q, %t`, nm, v, ok, v3, len(v3) != 0)
		***REMOVED***
	***REMOVED***

	c0 := Background()
	check(c0, "c0", "", "", "")

	c1 := WithValue(Background(), k1, "c1k1")
	check(c1, "c1", "c1k1", "", "")

	if got, want := fmt.Sprint(c1), `context.Background.WithValue(1, "c1k1")`; got != want ***REMOVED***
		t.Errorf("c.String() = %q want %q", got, want)
	***REMOVED***

	c2 := WithValue(c1, k2, "c2k2")
	check(c2, "c2", "c1k1", "c2k2", "")

	c3 := WithValue(c2, k3, "c3k3")
	check(c3, "c2", "c1k1", "c2k2", "c3k3")

	c4 := WithValue(c3, k1, nil)
	check(c4, "c4", "", "c2k2", "c3k3")

	o0 := otherContext***REMOVED***Background()***REMOVED***
	check(o0, "o0", "", "", "")

	o1 := otherContext***REMOVED***WithValue(Background(), k1, "c1k1")***REMOVED***
	check(o1, "o1", "c1k1", "", "")

	o2 := WithValue(o1, k2, "o2k2")
	check(o2, "o2", "c1k1", "o2k2", "")

	o3 := otherContext***REMOVED***c4***REMOVED***
	check(o3, "o3", "", "c2k2", "c3k3")

	o4 := WithValue(o3, k3, nil)
	check(o4, "o4", "", "c2k2", "")
***REMOVED***

func TestAllocs(t *testing.T) ***REMOVED***
	bg := Background()
	for _, test := range []struct ***REMOVED***
		desc       string
		f          func()
		limit      float64
		gccgoLimit float64
	***REMOVED******REMOVED***
		***REMOVED***
			desc:       "Background()",
			f:          func() ***REMOVED*** Background() ***REMOVED***,
			limit:      0,
			gccgoLimit: 0,
		***REMOVED***,
		***REMOVED***
			desc: fmt.Sprintf("WithValue(bg, %v, nil)", k1),
			f: func() ***REMOVED***
				c := WithValue(bg, k1, nil)
				c.Value(k1)
			***REMOVED***,
			limit:      3,
			gccgoLimit: 3,
		***REMOVED***,
		***REMOVED***
			desc: "WithTimeout(bg, 15*time.Millisecond)",
			f: func() ***REMOVED***
				c, _ := WithTimeout(bg, 15*time.Millisecond)
				<-c.Done()
			***REMOVED***,
			limit:      8,
			gccgoLimit: 16,
		***REMOVED***,
		***REMOVED***
			desc: "WithCancel(bg)",
			f: func() ***REMOVED***
				c, cancel := WithCancel(bg)
				cancel()
				<-c.Done()
			***REMOVED***,
			limit:      5,
			gccgoLimit: 8,
		***REMOVED***,
		***REMOVED***
			desc: "WithTimeout(bg, 100*time.Millisecond)",
			f: func() ***REMOVED***
				c, cancel := WithTimeout(bg, 100*time.Millisecond)
				cancel()
				<-c.Done()
			***REMOVED***,
			limit:      8,
			gccgoLimit: 25,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		limit := test.limit
		if runtime.Compiler == "gccgo" ***REMOVED***
			// gccgo does not yet do escape analysis.
			// TODO(iant): Remove this when gccgo does do escape analysis.
			limit = test.gccgoLimit
		***REMOVED***
		if n := testing.AllocsPerRun(100, test.f); n > limit ***REMOVED***
			t.Errorf("%s allocs = %f want %d", test.desc, n, int(limit))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSimultaneousCancels(t *testing.T) ***REMOVED***
	root, cancel := WithCancel(Background())
	m := map[Context]CancelFunc***REMOVED***root: cancel***REMOVED***
	q := []Context***REMOVED***root***REMOVED***
	// Create a tree of contexts.
	for len(q) != 0 && len(m) < 100 ***REMOVED***
		parent := q[0]
		q = q[1:]
		for i := 0; i < 4; i++ ***REMOVED***
			ctx, cancel := WithCancel(parent)
			m[ctx] = cancel
			q = append(q, ctx)
		***REMOVED***
	***REMOVED***
	// Start all the cancels in a random order.
	var wg sync.WaitGroup
	wg.Add(len(m))
	for _, cancel := range m ***REMOVED***
		go func(cancel CancelFunc) ***REMOVED***
			cancel()
			wg.Done()
		***REMOVED***(cancel)
	***REMOVED***
	// Wait on all the contexts in a random order.
	for ctx := range m ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		case <-time.After(1 * time.Second):
			buf := make([]byte, 10<<10)
			n := runtime.Stack(buf, true)
			t.Fatalf("timed out waiting for <-ctx.Done(); stacks:\n%s", buf[:n])
		***REMOVED***
	***REMOVED***
	// Wait for all the cancel functions to return.
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		wg.Wait()
		close(done)
	***REMOVED***()
	select ***REMOVED***
	case <-done:
	case <-time.After(1 * time.Second):
		buf := make([]byte, 10<<10)
		n := runtime.Stack(buf, true)
		t.Fatalf("timed out waiting for cancel functions; stacks:\n%s", buf[:n])
	***REMOVED***
***REMOVED***

func TestInterlockedCancels(t *testing.T) ***REMOVED***
	parent, cancelParent := WithCancel(Background())
	child, cancelChild := WithCancel(parent)
	go func() ***REMOVED***
		parent.Done()
		cancelChild()
	***REMOVED***()
	cancelParent()
	select ***REMOVED***
	case <-child.Done():
	case <-time.After(1 * time.Second):
		buf := make([]byte, 10<<10)
		n := runtime.Stack(buf, true)
		t.Fatalf("timed out waiting for child.Done(); stacks:\n%s", buf[:n])
	***REMOVED***
***REMOVED***

func TestLayersCancel(t *testing.T) ***REMOVED***
	testLayers(t, time.Now().UnixNano(), false)
***REMOVED***

func TestLayersTimeout(t *testing.T) ***REMOVED***
	testLayers(t, time.Now().UnixNano(), true)
***REMOVED***

func testLayers(t *testing.T, seed int64, testTimeout bool) ***REMOVED***
	rand.Seed(seed)
	errorf := func(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		t.Errorf(fmt.Sprintf("seed=%d: %s", seed, format), a...)
	***REMOVED***
	const (
		timeout   = 200 * time.Millisecond
		minLayers = 30
	)
	type value int
	var (
		vals      []*value
		cancels   []CancelFunc
		numTimers int
		ctx       = Background()
	)
	for i := 0; i < minLayers || numTimers == 0 || len(cancels) == 0 || len(vals) == 0; i++ ***REMOVED***
		switch rand.Intn(3) ***REMOVED***
		case 0:
			v := new(value)
			ctx = WithValue(ctx, v, v)
			vals = append(vals, v)
		case 1:
			var cancel CancelFunc
			ctx, cancel = WithCancel(ctx)
			cancels = append(cancels, cancel)
		case 2:
			var cancel CancelFunc
			ctx, cancel = WithTimeout(ctx, timeout)
			cancels = append(cancels, cancel)
			numTimers++
		***REMOVED***
	***REMOVED***
	checkValues := func(when string) ***REMOVED***
		for _, key := range vals ***REMOVED***
			if val := ctx.Value(key).(*value); key != val ***REMOVED***
				errorf("%s: ctx.Value(%p) = %p want %p", when, key, val, key)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	select ***REMOVED***
	case <-ctx.Done():
		errorf("ctx should not be canceled yet")
	default:
	***REMOVED***
	if s, prefix := fmt.Sprint(ctx), "context.Background."; !strings.HasPrefix(s, prefix) ***REMOVED***
		t.Errorf("ctx.String() = %q want prefix %q", s, prefix)
	***REMOVED***
	t.Log(ctx)
	checkValues("before cancel")
	if testTimeout ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
		case <-time.After(timeout + 100*time.Millisecond):
			errorf("ctx should have timed out")
		***REMOVED***
		checkValues("after timeout")
	***REMOVED*** else ***REMOVED***
		cancel := cancels[rand.Intn(len(cancels))]
		cancel()
		select ***REMOVED***
		case <-ctx.Done():
		default:
			errorf("ctx should be canceled")
		***REMOVED***
		checkValues("after cancel")
	***REMOVED***
***REMOVED***

func TestCancelRemoves(t *testing.T) ***REMOVED***
	checkChildren := func(when string, ctx Context, want int) ***REMOVED***
		if got := len(ctx.(*cancelCtx).children); got != want ***REMOVED***
			t.Errorf("%s: context has %d children, want %d", when, got, want)
		***REMOVED***
	***REMOVED***

	ctx, _ := WithCancel(Background())
	checkChildren("after creation", ctx, 0)
	_, cancel := WithCancel(ctx)
	checkChildren("with WithCancel child ", ctx, 1)
	cancel()
	checkChildren("after cancelling WithCancel child", ctx, 0)

	ctx, _ = WithCancel(Background())
	checkChildren("after creation", ctx, 0)
	_, cancel = WithTimeout(ctx, 60*time.Minute)
	checkChildren("with WithTimeout child ", ctx, 1)
	cancel()
	checkChildren("after cancelling WithTimeout child", ctx, 0)
***REMOVED***
