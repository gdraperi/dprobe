package debug

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestEnable(t *testing.T) ***REMOVED***
	defer func() ***REMOVED***
		os.Setenv("DEBUG", "")
		logrus.SetLevel(logrus.InfoLevel)
	***REMOVED***()
	Enable()
	if os.Getenv("DEBUG") != "1" ***REMOVED***
		t.Fatalf("expected DEBUG=1, got %s\n", os.Getenv("DEBUG"))
	***REMOVED***
	if logrus.GetLevel() != logrus.DebugLevel ***REMOVED***
		t.Fatalf("expected log level %v, got %v\n", logrus.DebugLevel, logrus.GetLevel())
	***REMOVED***
***REMOVED***

func TestDisable(t *testing.T) ***REMOVED***
	Disable()
	if os.Getenv("DEBUG") != "" ***REMOVED***
		t.Fatalf("expected DEBUG=\"\", got %s\n", os.Getenv("DEBUG"))
	***REMOVED***
	if logrus.GetLevel() != logrus.InfoLevel ***REMOVED***
		t.Fatalf("expected log level %v, got %v\n", logrus.InfoLevel, logrus.GetLevel())
	***REMOVED***
***REMOVED***

func TestEnabled(t *testing.T) ***REMOVED***
	Enable()
	if !IsEnabled() ***REMOVED***
		t.Fatal("expected debug enabled, got false")
	***REMOVED***
	Disable()
	if IsEnabled() ***REMOVED***
		t.Fatal("expected debug disabled, got true")
	***REMOVED***
***REMOVED***
