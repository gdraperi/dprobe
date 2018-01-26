package bridge

type setupStep func(*networkConfiguration, *bridgeInterface) error

type bridgeSetup struct ***REMOVED***
	config *networkConfiguration
	bridge *bridgeInterface
	steps  []setupStep
***REMOVED***

func newBridgeSetup(c *networkConfiguration, i *bridgeInterface) *bridgeSetup ***REMOVED***
	return &bridgeSetup***REMOVED***config: c, bridge: i***REMOVED***
***REMOVED***

func (b *bridgeSetup) apply() error ***REMOVED***
	for _, fn := range b.steps ***REMOVED***
		if err := fn(b.config, b.bridge); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (b *bridgeSetup) queueStep(step setupStep) ***REMOVED***
	b.steps = append(b.steps, step)
***REMOVED***
