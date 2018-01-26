package api

// Status can be used to query the Status endpoints
type Status struct ***REMOVED***
	c *Client
***REMOVED***

// Status returns a handle to the status endpoints
func (c *Client) Status() *Status ***REMOVED***
	return &Status***REMOVED***c***REMOVED***
***REMOVED***

// Leader is used to query for a known leader
func (s *Status) Leader() (string, error) ***REMOVED***
	r := s.c.newRequest("GET", "/v1/status/leader")
	_, resp, err := requireOK(s.c.doRequest(r))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer resp.Body.Close()

	var leader string
	if err := decodeBody(resp, &leader); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return leader, nil
***REMOVED***

// Peers is used to query for a known raft peers
func (s *Status) Peers() ([]string, error) ***REMOVED***
	r := s.c.newRequest("GET", "/v1/status/peers")
	_, resp, err := requireOK(s.c.doRequest(r))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	var peers []string
	if err := decodeBody(resp, &peers); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return peers, nil
***REMOVED***
