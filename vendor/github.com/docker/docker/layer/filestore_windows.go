package layer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// setOS writes the "os" file to the layer filestore
func (fm *fileMetadataTransaction) setOS(os string) error ***REMOVED***
	if os == "" ***REMOVED***
		return nil
	***REMOVED***
	return fm.ws.WriteFile("os", []byte(os), 0644)
***REMOVED***

// getOS reads the "os" file from the layer filestore
func (fms *fileMetadataStore) getOS(layer ChainID) (string, error) ***REMOVED***
	contentBytes, err := ioutil.ReadFile(fms.getLayerFilename(layer, "os"))
	if err != nil ***REMOVED***
		// For backwards compatibility, the os file may not exist. Default to "windows" if missing.
		if os.IsNotExist(err) ***REMOVED***
			return "windows", nil
		***REMOVED***
		return "", err
	***REMOVED***
	content := strings.TrimSpace(string(contentBytes))

	if content != "windows" && content != "linux" ***REMOVED***
		return "", fmt.Errorf("invalid operating system value: %s", content)
	***REMOVED***

	return content, nil
***REMOVED***
