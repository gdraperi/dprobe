package coordinate

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// GenerateClients returns a slice with nodes number of clients, all with the
// given config.
func GenerateClients(nodes int, config *Config) ([]*Client, error) ***REMOVED***
	clients := make([]*Client, nodes)
	for i, _ := range clients ***REMOVED***
		client, err := NewClient(config)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		clients[i] = client
	***REMOVED***
	return clients, nil
***REMOVED***

// GenerateLine returns a truth matrix as if all the nodes are in a straight linke
// with the given spacing between them.
func GenerateLine(nodes int, spacing time.Duration) [][]time.Duration ***REMOVED***
	truth := make([][]time.Duration, nodes)
	for i := range truth ***REMOVED***
		truth[i] = make([]time.Duration, nodes)
	***REMOVED***

	for i := 0; i < nodes; i++ ***REMOVED***
		for j := i + 1; j < nodes; j++ ***REMOVED***
			rtt := time.Duration(j-i) * spacing
			truth[i][j], truth[j][i] = rtt, rtt
		***REMOVED***
	***REMOVED***
	return truth
***REMOVED***

// GenerateGrid returns a truth matrix as if all the nodes are in a two dimensional
// grid with the given spacing between them.
func GenerateGrid(nodes int, spacing time.Duration) [][]time.Duration ***REMOVED***
	truth := make([][]time.Duration, nodes)
	for i := range truth ***REMOVED***
		truth[i] = make([]time.Duration, nodes)
	***REMOVED***

	n := int(math.Sqrt(float64(nodes)))
	for i := 0; i < nodes; i++ ***REMOVED***
		for j := i + 1; j < nodes; j++ ***REMOVED***
			x1, y1 := float64(i%n), float64(i/n)
			x2, y2 := float64(j%n), float64(j/n)
			dx, dy := x2-x1, y2-y1
			dist := math.Sqrt(dx*dx + dy*dy)
			rtt := time.Duration(dist * float64(spacing))
			truth[i][j], truth[j][i] = rtt, rtt
		***REMOVED***
	***REMOVED***
	return truth
***REMOVED***

// GenerateSplit returns a truth matrix as if half the nodes are close together in
// one location and half the nodes are close together in another. The lan factor
// is used to separate the nodes locally and the wan factor represents the split
// between the two sides.
func GenerateSplit(nodes int, lan time.Duration, wan time.Duration) [][]time.Duration ***REMOVED***
	truth := make([][]time.Duration, nodes)
	for i := range truth ***REMOVED***
		truth[i] = make([]time.Duration, nodes)
	***REMOVED***

	split := nodes / 2
	for i := 0; i < nodes; i++ ***REMOVED***
		for j := i + 1; j < nodes; j++ ***REMOVED***
			rtt := lan
			if (i <= split && j > split) || (i > split && j <= split) ***REMOVED***
				rtt += wan
			***REMOVED***
			truth[i][j], truth[j][i] = rtt, rtt
		***REMOVED***
	***REMOVED***
	return truth
***REMOVED***

// GenerateCircle returns a truth matrix for a set of nodes, evenly distributed
// around a circle with the given radius. The first node is at the "center" of the
// circle because it's equidistant from all the other nodes, but we place it at
// double the radius, so it should show up above all the other nodes in height.
func GenerateCircle(nodes int, radius time.Duration) [][]time.Duration ***REMOVED***
	truth := make([][]time.Duration, nodes)
	for i := range truth ***REMOVED***
		truth[i] = make([]time.Duration, nodes)
	***REMOVED***

	for i := 0; i < nodes; i++ ***REMOVED***
		for j := i + 1; j < nodes; j++ ***REMOVED***
			var rtt time.Duration
			if i == 0 ***REMOVED***
				rtt = 2 * radius
			***REMOVED*** else ***REMOVED***
				t1 := 2.0 * math.Pi * float64(i) / float64(nodes)
				x1, y1 := math.Cos(t1), math.Sin(t1)
				t2 := 2.0 * math.Pi * float64(j) / float64(nodes)
				x2, y2 := math.Cos(t2), math.Sin(t2)
				dx, dy := x2-x1, y2-y1
				dist := math.Sqrt(dx*dx + dy*dy)
				rtt = time.Duration(dist * float64(radius))
			***REMOVED***
			truth[i][j], truth[j][i] = rtt, rtt
		***REMOVED***
	***REMOVED***
	return truth
***REMOVED***

// GenerateRandom returns a truth matrix for a set of nodes with normally
// distributed delays, with the given mean and deviation. The RNG is re-seeded
// so you always get the same matrix for a given size.
func GenerateRandom(nodes int, mean time.Duration, deviation time.Duration) [][]time.Duration ***REMOVED***
	rand.Seed(1)

	truth := make([][]time.Duration, nodes)
	for i := range truth ***REMOVED***
		truth[i] = make([]time.Duration, nodes)
	***REMOVED***

	for i := 0; i < nodes; i++ ***REMOVED***
		for j := i + 1; j < nodes; j++ ***REMOVED***
			rttSeconds := rand.NormFloat64()*deviation.Seconds() + mean.Seconds()
			rtt := time.Duration(rttSeconds * secondsToNanoseconds)
			truth[i][j], truth[j][i] = rtt, rtt
		***REMOVED***
	***REMOVED***
	return truth
***REMOVED***

// Simulate runs the given number of cycles using the given list of clients and
// truth matrix. On each cycle, each client will pick a random node and observe
// the truth RTT, updating its coordinate estimate. The RNG is re-seeded for
// each simulation run to get deterministic results (for this algorithm and the
// underlying algorithm which will use random numbers for position vectors when
// starting out with everything at the origin).
func Simulate(clients []*Client, truth [][]time.Duration, cycles int) ***REMOVED***
	rand.Seed(1)

	nodes := len(clients)
	for cycle := 0; cycle < cycles; cycle++ ***REMOVED***
		for i, _ := range clients ***REMOVED***
			if j := rand.Intn(nodes); j != i ***REMOVED***
				c := clients[j].GetCoordinate()
				rtt := truth[i][j]
				node := fmt.Sprintf("node_%d", j)
				clients[i].Update(node, c, rtt)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stats is returned from the Evaluate function with a summary of the algorithm
// performance.
type Stats struct ***REMOVED***
	ErrorMax float64
	ErrorAvg float64
***REMOVED***

// Evaluate uses the coordinates of the given clients to calculate estimated
// distances and compares them with the given truth matrix, returning summary
// stats.
func Evaluate(clients []*Client, truth [][]time.Duration) (stats Stats) ***REMOVED***
	nodes := len(clients)
	count := 0
	for i := 0; i < nodes; i++ ***REMOVED***
		for j := i + 1; j < nodes; j++ ***REMOVED***
			est := clients[i].DistanceTo(clients[j].GetCoordinate()).Seconds()
			actual := truth[i][j].Seconds()
			error := math.Abs(est-actual) / actual
			stats.ErrorMax = math.Max(stats.ErrorMax, error)
			stats.ErrorAvg += error
			count += 1
		***REMOVED***
	***REMOVED***

	stats.ErrorAvg /= float64(count)
	fmt.Printf("Error avg=%9.6f max=%9.6f\n", stats.ErrorAvg, stats.ErrorMax)
	return
***REMOVED***
