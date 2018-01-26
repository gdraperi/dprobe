package progress

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestOutputOnPrematureClose(t *testing.T) ***REMOVED***
	content := []byte("TESTING")
	reader := ioutil.NopCloser(bytes.NewReader(content))
	progressChan := make(chan Progress, 10)

	pr := NewProgressReader(reader, ChanOutput(progressChan), int64(len(content)), "Test", "Read")

	part := make([]byte, 4)
	_, err := io.ReadFull(pr, part)
	if err != nil ***REMOVED***
		pr.Close()
		t.Fatal(err)
	***REMOVED***

drainLoop:
	for ***REMOVED***
		select ***REMOVED***
		case <-progressChan:
		default:
			break drainLoop
		***REMOVED***
	***REMOVED***

	pr.Close()

	select ***REMOVED***
	case <-progressChan:
	default:
		t.Fatalf("Expected some output when closing prematurely")
	***REMOVED***
***REMOVED***

func TestCompleteSilently(t *testing.T) ***REMOVED***
	content := []byte("TESTING")
	reader := ioutil.NopCloser(bytes.NewReader(content))
	progressChan := make(chan Progress, 10)

	pr := NewProgressReader(reader, ChanOutput(progressChan), int64(len(content)), "Test", "Read")

	out, err := ioutil.ReadAll(pr)
	if err != nil ***REMOVED***
		pr.Close()
		t.Fatal(err)
	***REMOVED***
	if string(out) != "TESTING" ***REMOVED***
		pr.Close()
		t.Fatalf("Unexpected output %q from reader", string(out))
	***REMOVED***

drainLoop:
	for ***REMOVED***
		select ***REMOVED***
		case <-progressChan:
		default:
			break drainLoop
		***REMOVED***
	***REMOVED***

	pr.Close()

	select ***REMOVED***
	case <-progressChan:
		t.Fatalf("Should have closed silently when read is complete")
	default:
	***REMOVED***
***REMOVED***
