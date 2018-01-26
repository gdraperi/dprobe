package hcsshim

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

func hnsCall(method, path, request string, returnResponse interface***REMOVED******REMOVED***) error ***REMOVED***
	var responseBuffer *uint16
	logrus.Debugf("[%s]=>[%s] Request : %s", method, path, request)

	err := _hnsCall(method, path, request, &responseBuffer)
	if err != nil ***REMOVED***
		return makeError(err, "hnsCall ", "")
	***REMOVED***
	response := convertAndFreeCoTaskMemString(responseBuffer)

	hnsresponse := &hnsResponse***REMOVED******REMOVED***
	if err = json.Unmarshal([]byte(response), &hnsresponse); err != nil ***REMOVED***
		return err
	***REMOVED***

	if !hnsresponse.Success ***REMOVED***
		return fmt.Errorf("HNS failed with error : %s", hnsresponse.Error)
	***REMOVED***

	if len(hnsresponse.Output) == 0 ***REMOVED***
		return nil
	***REMOVED***

	logrus.Debugf("Network Response : %s", hnsresponse.Output)
	err = json.Unmarshal(hnsresponse.Output, returnResponse)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
