package coordinate

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// Client manages the estimated network coordinate for a given node, and adjusts
// it as the node observes round trip times and estimated coordinates from other
// nodes. The core algorithm is based on Vivaldi, see the documentation for Config
// for more details.
type Client struct ***REMOVED***
	// coord is the current estimate of the client's network coordinate.
	coord *Coordinate

	// origin is a coordinate sitting at the origin.
	origin *Coordinate

	// config contains the tuning parameters that govern the performance of
	// the algorithm.
	config *Config

	// adjustmentIndex is the current index into the adjustmentSamples slice.
	adjustmentIndex uint

	// adjustment is used to store samples for the adjustment calculation.
	adjustmentSamples []float64

	// latencyFilterSamples is used to store the last several RTT samples,
	// keyed by node name. We will use the config's LatencyFilterSamples
	// value to determine how many samples we keep, per node.
	latencyFilterSamples map[string][]float64

	// mutex enables safe concurrent access to the client.
	mutex sync.RWMutex
***REMOVED***

// NewClient creates a new Client and verifies the configuration is valid.
func NewClient(config *Config) (*Client, error) ***REMOVED***
	if !(config.Dimensionality > 0) ***REMOVED***
		return nil, fmt.Errorf("dimensionality must be >0")
	***REMOVED***

	return &Client***REMOVED***
		coord:                NewCoordinate(config),
		origin:               NewCoordinate(config),
		config:               config,
		adjustmentIndex:      0,
		adjustmentSamples:    make([]float64, config.AdjustmentWindowSize),
		latencyFilterSamples: make(map[string][]float64),
	***REMOVED***, nil
***REMOVED***

// GetCoordinate returns a copy of the coordinate for this client.
func (c *Client) GetCoordinate() *Coordinate ***REMOVED***
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.coord.Clone()
***REMOVED***

// SetCoordinate forces the client's coordinate to a known state.
func (c *Client) SetCoordinate(coord *Coordinate) ***REMOVED***
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.coord = coord.Clone()
***REMOVED***

// ForgetNode removes any client state for the given node.
func (c *Client) ForgetNode(node string) ***REMOVED***
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.latencyFilterSamples, node)
***REMOVED***

// latencyFilter applies a simple moving median filter with a new sample for
// a node. This assumes that the mutex has been locked already.
func (c *Client) latencyFilter(node string, rttSeconds float64) float64 ***REMOVED***
	samples, ok := c.latencyFilterSamples[node]
	if !ok ***REMOVED***
		samples = make([]float64, 0, c.config.LatencyFilterSize)
	***REMOVED***

	// Add the new sample and trim the list, if needed.
	samples = append(samples, rttSeconds)
	if len(samples) > int(c.config.LatencyFilterSize) ***REMOVED***
		samples = samples[1:]
	***REMOVED***
	c.latencyFilterSamples[node] = samples

	// Sort a copy of the samples and return the median.
	sorted := make([]float64, len(samples))
	copy(sorted, samples)
	sort.Float64s(sorted)
	return sorted[len(sorted)/2]
***REMOVED***

// updateVivialdi updates the Vivaldi portion of the client's coordinate. This
// assumes that the mutex has been locked already.
func (c *Client) updateVivaldi(other *Coordinate, rttSeconds float64) ***REMOVED***
	const zeroThreshold = 1.0e-6

	dist := c.coord.DistanceTo(other).Seconds()
	if rttSeconds < zeroThreshold ***REMOVED***
		rttSeconds = zeroThreshold
	***REMOVED***
	wrongness := math.Abs(dist-rttSeconds) / rttSeconds

	totalError := c.coord.Error + other.Error
	if totalError < zeroThreshold ***REMOVED***
		totalError = zeroThreshold
	***REMOVED***
	weight := c.coord.Error / totalError

	c.coord.Error = c.config.VivaldiCE*weight*wrongness + c.coord.Error*(1.0-c.config.VivaldiCE*weight)
	if c.coord.Error > c.config.VivaldiErrorMax ***REMOVED***
		c.coord.Error = c.config.VivaldiErrorMax
	***REMOVED***

	delta := c.config.VivaldiCC * weight
	force := delta * (rttSeconds - dist)
	c.coord = c.coord.ApplyForce(c.config, force, other)
***REMOVED***

// updateAdjustment updates the adjustment portion of the client's coordinate, if
// the feature is enabled. This assumes that the mutex has been locked already.
func (c *Client) updateAdjustment(other *Coordinate, rttSeconds float64) ***REMOVED***
	if c.config.AdjustmentWindowSize == 0 ***REMOVED***
		return
	***REMOVED***

	// Note that the existing adjustment factors don't figure in to this
	// calculation so we use the raw distance here.
	dist := c.coord.rawDistanceTo(other)
	c.adjustmentSamples[c.adjustmentIndex] = rttSeconds - dist
	c.adjustmentIndex = (c.adjustmentIndex + 1) % c.config.AdjustmentWindowSize

	sum := 0.0
	for _, sample := range c.adjustmentSamples ***REMOVED***
		sum += sample
	***REMOVED***
	c.coord.Adjustment = sum / (2.0 * float64(c.config.AdjustmentWindowSize))
***REMOVED***

// updateGravity applies a small amount of gravity to pull coordinates towards
// the center of the coordinate system to combat drift. This assumes that the
// mutex is locked already.
func (c *Client) updateGravity() ***REMOVED***
	dist := c.origin.DistanceTo(c.coord).Seconds()
	force := -1.0 * math.Pow(dist/c.config.GravityRho, 2.0)
	c.coord = c.coord.ApplyForce(c.config, force, c.origin)
***REMOVED***

// Update takes other, a coordinate for another node, and rtt, a round trip
// time observation for a ping to that node, and updates the estimated position of
// the client's coordinate. Returns the updated coordinate.
func (c *Client) Update(node string, other *Coordinate, rtt time.Duration) *Coordinate ***REMOVED***
	c.mutex.Lock()
	defer c.mutex.Unlock()

	rttSeconds := c.latencyFilter(node, rtt.Seconds())
	c.updateVivaldi(other, rttSeconds)
	c.updateAdjustment(other, rttSeconds)
	c.updateGravity()
	return c.coord.Clone()
***REMOVED***

// DistanceTo returns the estimated RTT from the client's coordinate to other, the
// coordinate for another node.
func (c *Client) DistanceTo(other *Coordinate) time.Duration ***REMOVED***
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.coord.DistanceTo(other)
***REMOVED***
