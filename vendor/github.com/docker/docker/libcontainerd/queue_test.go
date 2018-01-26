package libcontainerd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) ***REMOVED***
	var (
		q             queue
		serialization = 1
	)

	q.append("aaa", func() ***REMOVED***
		//simulate a long time task
		time.Sleep(10 * time.Millisecond)
		require.EqualValues(t, serialization, 1)
		serialization = 2
	***REMOVED***)
	q.append("aaa", func() ***REMOVED***
		require.EqualValues(t, serialization, 2)
		serialization = 3
	***REMOVED***)
	q.append("aaa", func() ***REMOVED***
		require.EqualValues(t, serialization, 3)
		serialization = 4
	***REMOVED***)
	time.Sleep(20 * time.Millisecond)
***REMOVED***
