package diagnose

import "fmt"

// StringInterface interface that has to be implemented by messages
type StringInterface interface ***REMOVED***
	String() string
***REMOVED***

// CommandSucceed creates a success message
func CommandSucceed(result StringInterface) *HTTPResult ***REMOVED***
	return &HTTPResult***REMOVED***
		Message: "OK",
		Details: result,
	***REMOVED***
***REMOVED***

// FailCommand creates a failure message with error
func FailCommand(err error) *HTTPResult ***REMOVED***
	return &HTTPResult***REMOVED***
		Message: "FAIL",
		Details: &ErrorCmd***REMOVED***Error: err.Error()***REMOVED***,
	***REMOVED***
***REMOVED***

// WrongCommand creates a wrong command response
func WrongCommand(message, usage string) *HTTPResult ***REMOVED***
	return &HTTPResult***REMOVED***
		Message: message,
		Details: &UsageCmd***REMOVED***Usage: usage***REMOVED***,
	***REMOVED***
***REMOVED***

// HTTPResult Diagnose Server HTTP result operation
type HTTPResult struct ***REMOVED***
	Message string          `json:"message"`
	Details StringInterface `json:"details"`
***REMOVED***

func (h *HTTPResult) String() string ***REMOVED***
	rsp := h.Message
	if h.Details != nil ***REMOVED***
		rsp += "\n" + h.Details.String()
	***REMOVED***
	return rsp
***REMOVED***

// UsageCmd command with usage field
type UsageCmd struct ***REMOVED***
	Usage string `json:"usage"`
***REMOVED***

func (u *UsageCmd) String() string ***REMOVED***
	return "Usage: " + u.Usage
***REMOVED***

// StringCmd command with info string
type StringCmd struct ***REMOVED***
	Info string `json:"info"`
***REMOVED***

func (s *StringCmd) String() string ***REMOVED***
	return s.Info
***REMOVED***

// ErrorCmd command with error
type ErrorCmd struct ***REMOVED***
	Error string `json:"error"`
***REMOVED***

func (e *ErrorCmd) String() string ***REMOVED***
	return "Error: " + e.Error
***REMOVED***

// TableObj network db table object
type TableObj struct ***REMOVED***
	Length   int               `json:"size"`
	Elements []StringInterface `json:"entries"`
***REMOVED***

func (t *TableObj) String() string ***REMOVED***
	output := fmt.Sprintf("total entries: %d\n", t.Length)
	for _, e := range t.Elements ***REMOVED***
		output += e.String()
	***REMOVED***
	return output
***REMOVED***

// PeerEntryObj entry in the networkdb peer table
type PeerEntryObj struct ***REMOVED***
	Index int    `json:"-"`
	Name  string `json:"-=name"`
	IP    string `json:"ip"`
***REMOVED***

func (p *PeerEntryObj) String() string ***REMOVED***
	return fmt.Sprintf("%d) %s -> %s\n", p.Index, p.Name, p.IP)
***REMOVED***

// TableEntryObj network db table entry object
type TableEntryObj struct ***REMOVED***
	Index int    `json:"-"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Owner string `json:"owner"`
***REMOVED***

func (t *TableEntryObj) String() string ***REMOVED***
	return fmt.Sprintf("%d) k:`%s` -> v:`%s` owner:`%s`\n", t.Index, t.Key, t.Value, t.Owner)
***REMOVED***

// TableEndpointsResult fully typed message for proper unmarshaling on the client side
type TableEndpointsResult struct ***REMOVED***
	TableObj
	Elements []TableEntryObj `json:"entries"`
***REMOVED***

// TablePeersResult fully typed message for proper unmarshaling on the client side
type TablePeersResult struct ***REMOVED***
	TableObj
	Elements []PeerEntryObj `json:"entries"`
***REMOVED***
