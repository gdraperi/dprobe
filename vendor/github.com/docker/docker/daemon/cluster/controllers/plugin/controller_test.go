package plugin

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/docker/distribution/reference"
	enginetypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm/runtime"
	"github.com/docker/docker/pkg/pubsub"
	"github.com/docker/docker/plugin"
	"github.com/docker/docker/plugin/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	pluginTestName          = "test"
	pluginTestRemote        = "testremote"
	pluginTestRemoteUpgrade = "testremote2"
)

func TestPrepare(t *testing.T) ***REMOVED***
	b := newMockBackend()
	c := newTestController(b, false)
	ctx := context.Background()

	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if b.p == nil ***REMOVED***
		t.Fatal("pull not performed")
	***REMOVED***

	c = newTestController(b, false)
	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if b.p == nil ***REMOVED***
		t.Fatal("unexpected nil")
	***REMOVED***
	if b.p.PluginObj.PluginReference != pluginTestRemoteUpgrade ***REMOVED***
		t.Fatal("upgrade not performed")
	***REMOVED***

	c = newTestController(b, false)
	c.serviceID = "1"
	if err := c.Prepare(ctx); err == nil ***REMOVED***
		t.Fatal("expected error on prepare")
	***REMOVED***
***REMOVED***

func TestStart(t *testing.T) ***REMOVED***
	b := newMockBackend()
	c := newTestController(b, false)
	ctx := context.Background()

	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if !b.p.IsEnabled() ***REMOVED***
		t.Fatal("expected plugin to be enabled")
	***REMOVED***

	c = newTestController(b, true)
	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if b.p.IsEnabled() ***REMOVED***
		t.Fatal("expected plugin to be disabled")
	***REMOVED***

	c = newTestController(b, false)
	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !b.p.IsEnabled() ***REMOVED***
		t.Fatal("expected plugin to be enabled")
	***REMOVED***
***REMOVED***

func TestWaitCancel(t *testing.T) ***REMOVED***
	b := newMockBackend()
	c := newTestController(b, true)
	ctx := context.Background()
	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ctxCancel, cancel := context.WithCancel(ctx)
	chErr := make(chan error)
	go func() ***REMOVED***
		chErr <- c.Wait(ctxCancel)
	***REMOVED***()
	cancel()
	select ***REMOVED***
	case err := <-chErr:
		if err != context.Canceled ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for cancelation")
	***REMOVED***
***REMOVED***

func TestWaitDisabled(t *testing.T) ***REMOVED***
	b := newMockBackend()
	c := newTestController(b, true)
	ctx := context.Background()
	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	chErr := make(chan error)
	go func() ***REMOVED***
		chErr <- c.Wait(ctx)
	***REMOVED***()

	if err := b.Enable("test", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	select ***REMOVED***
	case err := <-chErr:
		if err == nil ***REMOVED***
			t.Fatal("expected error")
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for event")
	***REMOVED***

	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ctxWaitReady, cancelCtxWaitReady := context.WithTimeout(ctx, 30*time.Second)
	c.signalWaitReady = cancelCtxWaitReady
	defer cancelCtxWaitReady()

	go func() ***REMOVED***
		chErr <- c.Wait(ctx)
	***REMOVED***()

	chEvent, cancel := b.SubscribeEvents(1)
	defer cancel()

	if err := b.Disable("test", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	select ***REMOVED***
	case <-chEvent:
		<-ctxWaitReady.Done()
		if err := ctxWaitReady.Err(); err == context.DeadlineExceeded ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		select ***REMOVED***
		case <-chErr:
			t.Fatal("wait returned unexpectedly")
		default:
			// all good
		***REMOVED***
	case <-chErr:
		t.Fatal("wait returned unexpectedly")
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for event")
	***REMOVED***

	if err := b.Remove("test", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	select ***REMOVED***
	case err := <-chErr:
		if err == nil ***REMOVED***
			t.Fatal("expected error")
		***REMOVED***
		if !strings.Contains(err.Error(), "removed") ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for event")
	***REMOVED***
***REMOVED***

func TestWaitEnabled(t *testing.T) ***REMOVED***
	b := newMockBackend()
	c := newTestController(b, false)
	ctx := context.Background()
	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	chErr := make(chan error)
	go func() ***REMOVED***
		chErr <- c.Wait(ctx)
	***REMOVED***()

	if err := b.Disable("test", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	select ***REMOVED***
	case err := <-chErr:
		if err == nil ***REMOVED***
			t.Fatal("expected error")
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for event")
	***REMOVED***

	if err := c.Start(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ctxWaitReady, ctxWaitCancel := context.WithCancel(ctx)
	c.signalWaitReady = ctxWaitCancel
	defer ctxWaitCancel()

	go func() ***REMOVED***
		chErr <- c.Wait(ctx)
	***REMOVED***()

	chEvent, cancel := b.SubscribeEvents(1)
	defer cancel()

	if err := b.Enable("test", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	select ***REMOVED***
	case <-chEvent:
		<-ctxWaitReady.Done()
		if err := ctxWaitReady.Err(); err == context.DeadlineExceeded ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		select ***REMOVED***
		case <-chErr:
			t.Fatal("wait returned unexpectedly")
		default:
			// all good
		***REMOVED***
	case <-chErr:
		t.Fatal("wait returned unexpectedly")
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for event")
	***REMOVED***

	if err := b.Remove("test", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	select ***REMOVED***
	case err := <-chErr:
		if err == nil ***REMOVED***
			t.Fatal("expected error")
		***REMOVED***
		if !strings.Contains(err.Error(), "removed") ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for event")
	***REMOVED***
***REMOVED***

func TestRemove(t *testing.T) ***REMOVED***
	b := newMockBackend()
	c := newTestController(b, false)
	ctx := context.Background()

	if err := c.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := c.Shutdown(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	c2 := newTestController(b, false)
	if err := c2.Prepare(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := c.Remove(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if b.p == nil ***REMOVED***
		t.Fatal("plugin removed unexpectedly")
	***REMOVED***
	if err := c2.Shutdown(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := c2.Remove(ctx); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if b.p != nil ***REMOVED***
		t.Fatal("expected plugin to be removed")
	***REMOVED***
***REMOVED***

func newTestController(b Backend, disabled bool) *Controller ***REMOVED***
	return &Controller***REMOVED***
		logger:  &logrus.Entry***REMOVED***Logger: &logrus.Logger***REMOVED***Out: ioutil.Discard***REMOVED******REMOVED***,
		backend: b,
		spec: runtime.PluginSpec***REMOVED***
			Name:     pluginTestName,
			Remote:   pluginTestRemote,
			Disabled: disabled,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func newMockBackend() *mockBackend ***REMOVED***
	return &mockBackend***REMOVED***
		pub: pubsub.NewPublisher(0, 0),
	***REMOVED***
***REMOVED***

type mockBackend struct ***REMOVED***
	p   *v2.Plugin
	pub *pubsub.Publisher
***REMOVED***

func (m *mockBackend) Disable(name string, config *enginetypes.PluginDisableConfig) error ***REMOVED***
	m.p.PluginObj.Enabled = false
	m.pub.Publish(plugin.EventDisable***REMOVED******REMOVED***)
	return nil
***REMOVED***

func (m *mockBackend) Enable(name string, config *enginetypes.PluginEnableConfig) error ***REMOVED***
	m.p.PluginObj.Enabled = true
	m.pub.Publish(plugin.EventEnable***REMOVED******REMOVED***)
	return nil
***REMOVED***

func (m *mockBackend) Remove(name string, config *enginetypes.PluginRmConfig) error ***REMOVED***
	m.p = nil
	m.pub.Publish(plugin.EventRemove***REMOVED******REMOVED***)
	return nil
***REMOVED***

func (m *mockBackend) Pull(ctx context.Context, ref reference.Named, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, privileges enginetypes.PluginPrivileges, outStream io.Writer, opts ...plugin.CreateOpt) error ***REMOVED***
	m.p = &v2.Plugin***REMOVED***
		PluginObj: enginetypes.Plugin***REMOVED***
			ID:              "1234",
			Name:            name,
			PluginReference: ref.String(),
		***REMOVED***,
	***REMOVED***
	return nil
***REMOVED***

func (m *mockBackend) Upgrade(ctx context.Context, ref reference.Named, name string, metaHeaders http.Header, authConfig *enginetypes.AuthConfig, privileges enginetypes.PluginPrivileges, outStream io.Writer) error ***REMOVED***
	m.p.PluginObj.PluginReference = pluginTestRemoteUpgrade
	return nil
***REMOVED***

func (m *mockBackend) Get(name string) (*v2.Plugin, error) ***REMOVED***
	if m.p == nil ***REMOVED***
		return nil, errors.New("not found")
	***REMOVED***
	return m.p, nil
***REMOVED***

func (m *mockBackend) SubscribeEvents(buffer int, events ...plugin.Event) (eventCh <-chan interface***REMOVED******REMOVED***, cancel func()) ***REMOVED***
	ch := m.pub.SubscribeTopicWithBuffer(nil, buffer)
	cancel = func() ***REMOVED*** m.pub.Evict(ch) ***REMOVED***
	return ch, cancel
***REMOVED***
