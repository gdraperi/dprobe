package logrus_test

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Example_basic() ***REMOVED***
	var log = logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter)                     //default
	log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	log.Level = logrus.DebugLevel
	log.Out = os.Stdout

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil ***REMOVED***
	// 	log.Out = file
	// ***REMOVED*** else ***REMOVED***
	// 	log.Info("Failed to log to file, using default stderr")
	// ***REMOVED***

	defer func() ***REMOVED***
		err := recover()
		if err != nil ***REMOVED***
			entry := err.(*logrus.Entry)
			log.WithFields(logrus.Fields***REMOVED***
				"omg":         true,
				"err_animal":  entry.Data["animal"],
				"err_size":    entry.Data["size"],
				"err_level":   entry.Level,
				"err_message": entry.Message,
				"number":      100,
			***REMOVED***).Error("The ice breaks!") // or use Fatal() to force the process to exit with a nonzero code
		***REMOVED***
	***REMOVED***()

	log.WithFields(logrus.Fields***REMOVED***
		"animal": "walrus",
		"number": 8,
	***REMOVED***).Debug("Started observing beach")

	log.WithFields(logrus.Fields***REMOVED***
		"animal": "walrus",
		"size":   10,
	***REMOVED***).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields***REMOVED***
		"omg":    true,
		"number": 122,
	***REMOVED***).Warn("The group's number increased tremendously!")

	log.WithFields(logrus.Fields***REMOVED***
		"temperature": -4,
	***REMOVED***).Debug("Temperature changes")

	log.WithFields(logrus.Fields***REMOVED***
		"animal": "orca",
		"size":   9009,
	***REMOVED***).Panic("It's over 9000!")

	// Output:
	// level=debug msg="Started observing beach" animal=walrus number=8
	// level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
	// level=warning msg="The group's number increased tremendously!" number=122 omg=true
	// level=debug msg="Temperature changes" temperature=-4
	// level=panic msg="It's over 9000!" animal=orca size=9009
	// level=error msg="The ice breaks!" err_animal=orca err_level=panic err_message="It's over 9000!" err_size=9009 number=100 omg=true
***REMOVED***
