package exchange

import (
	"context"
	"strings"
	"time"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/events"
	"github.com/containerd/containerd/filters"
	"github.com/containerd/containerd/identifiers"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/typeurl"
	goevents "github.com/docker/go-events"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Exchange broadcasts events
type Exchange struct ***REMOVED***
	broadcaster *goevents.Broadcaster
***REMOVED***

// NewExchange returns a new event Exchange
func NewExchange() *Exchange ***REMOVED***
	return &Exchange***REMOVED***
		broadcaster: goevents.NewBroadcaster(),
	***REMOVED***
***REMOVED***

var _ events.Publisher = &Exchange***REMOVED******REMOVED***
var _ events.Forwarder = &Exchange***REMOVED******REMOVED***
var _ events.Subscriber = &Exchange***REMOVED******REMOVED***

// Forward accepts an envelope to be direcly distributed on the exchange.
//
// This is useful when an event is forwaded on behalf of another namespace or
// when the event is propagated on behalf of another publisher.
func (e *Exchange) Forward(ctx context.Context, envelope *events.Envelope) (err error) ***REMOVED***
	if err := validateEnvelope(envelope); err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		logger := log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"topic": envelope.Topic,
			"ns":    envelope.Namespace,
			"type":  envelope.Event.TypeUrl,
		***REMOVED***)

		if err != nil ***REMOVED***
			logger.WithError(err).Error("error forwarding event")
		***REMOVED*** else ***REMOVED***
			logger.Debug("event forwarded")
		***REMOVED***
	***REMOVED***()

	return e.broadcaster.Write(envelope)
***REMOVED***

// Publish packages and sends an event. The caller will be considered the
// initial publisher of the event. This means the timestamp will be calculated
// at this point and this method may read from the calling context.
func (e *Exchange) Publish(ctx context.Context, topic string, event events.Event) (err error) ***REMOVED***
	var (
		namespace string
		encoded   *types.Any
		envelope  events.Envelope
	)

	namespace, err = namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed publishing event")
	***REMOVED***
	if err := validateTopic(topic); err != nil ***REMOVED***
		return errors.Wrapf(err, "envelope topic %q", topic)
	***REMOVED***

	encoded, err = typeurl.MarshalAny(event)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	envelope.Timestamp = time.Now().UTC()
	envelope.Namespace = namespace
	envelope.Topic = topic
	envelope.Event = encoded

	defer func() ***REMOVED***
		logger := log.G(ctx).WithFields(logrus.Fields***REMOVED***
			"topic": envelope.Topic,
			"ns":    envelope.Namespace,
			"type":  envelope.Event.TypeUrl,
		***REMOVED***)

		if err != nil ***REMOVED***
			logger.WithError(err).Error("error publishing event")
		***REMOVED*** else ***REMOVED***
			logger.Debug("event published")
		***REMOVED***
	***REMOVED***()

	return e.broadcaster.Write(&envelope)
***REMOVED***

// Subscribe to events on the exchange. Events are sent through the returned
// channel ch. If an error is encountered, it will be sent on channel errs and
// errs will be closed. To end the subscription, cancel the provided context.
//
// Zero or more filters may be provided as strings. Only events that match
// *any* of the provided filters will be sent on the channel. The filters use
// the standard containerd filters package syntax.
func (e *Exchange) Subscribe(ctx context.Context, fs ...string) (ch <-chan *events.Envelope, errs <-chan error) ***REMOVED***
	var (
		evch                  = make(chan *events.Envelope)
		errq                  = make(chan error, 1)
		channel               = goevents.NewChannel(0)
		queue                 = goevents.NewQueue(channel)
		dst     goevents.Sink = queue
	)

	closeAll := func() ***REMOVED***
		defer close(errq)
		defer e.broadcaster.Remove(dst)
		defer queue.Close()
		defer channel.Close()
	***REMOVED***

	ch = evch
	errs = errq

	if len(fs) > 0 ***REMOVED***
		filter, err := filters.ParseAll(fs...)
		if err != nil ***REMOVED***
			errq <- errors.Wrapf(err, "failed parsing subscription filters")
			closeAll()
			return
		***REMOVED***

		dst = goevents.NewFilter(queue, goevents.MatcherFunc(func(gev goevents.Event) bool ***REMOVED***
			return filter.Match(adapt(gev))
		***REMOVED***))
	***REMOVED***

	e.broadcaster.Add(dst)

	go func() ***REMOVED***
		defer closeAll()

		var err error
	loop:
		for ***REMOVED***
			select ***REMOVED***
			case ev := <-channel.C:
				env, ok := ev.(*events.Envelope)
				if !ok ***REMOVED***
					// TODO(stevvooe): For the most part, we are well protected
					// from this condition. Both Forward and Publish protect
					// from this.
					err = errors.Errorf("invalid envelope encountered %#v; please file a bug", ev)
					break
				***REMOVED***

				select ***REMOVED***
				case evch <- env:
				case <-ctx.Done():
					break loop
				***REMOVED***
			case <-ctx.Done():
				break loop
			***REMOVED***
		***REMOVED***

		if err == nil ***REMOVED***
			if cerr := ctx.Err(); cerr != context.Canceled ***REMOVED***
				err = cerr
			***REMOVED***
		***REMOVED***

		errq <- err
	***REMOVED***()

	return
***REMOVED***

func validateTopic(topic string) error ***REMOVED***
	if topic == "" ***REMOVED***
		return errors.Wrap(errdefs.ErrInvalidArgument, "must not be empty")
	***REMOVED***

	if topic[0] != '/' ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "must start with '/'")
	***REMOVED***

	if len(topic) == 1 ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "must have at least one component")
	***REMOVED***

	components := strings.Split(topic[1:], "/")
	for _, component := range components ***REMOVED***
		if err := identifiers.Validate(component); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed validation on component %q", component)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func validateEnvelope(envelope *events.Envelope) error ***REMOVED***
	if err := namespaces.Validate(envelope.Namespace); err != nil ***REMOVED***
		return errors.Wrapf(err, "event envelope has invalid namespace")
	***REMOVED***

	if err := validateTopic(envelope.Topic); err != nil ***REMOVED***
		return errors.Wrapf(err, "envelope topic %q", envelope.Topic)
	***REMOVED***

	if envelope.Timestamp.IsZero() ***REMOVED***
		return errors.Wrapf(errdefs.ErrInvalidArgument, "timestamp must be set on forwarded event")
	***REMOVED***

	return nil
***REMOVED***

func adapt(ev interface***REMOVED******REMOVED***) filters.Adaptor ***REMOVED***
	if adaptor, ok := ev.(filters.Adaptor); ok ***REMOVED***
		return adaptor
	***REMOVED***

	return filters.AdapterFunc(func(fieldpath []string) (string, bool) ***REMOVED***
		return "", false
	***REMOVED***)
***REMOVED***
