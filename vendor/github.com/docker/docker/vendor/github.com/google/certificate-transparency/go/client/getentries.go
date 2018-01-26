package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	ct "github.com/google/certificate-transparency/go"
	"golang.org/x/net/context"
)

// LeafEntry respresents a JSON leaf entry.
type LeafEntry struct ***REMOVED***
	LeafInput []byte `json:"leaf_input"`
	ExtraData []byte `json:"extra_data"`
***REMOVED***

// GetEntriesResponse respresents the JSON response to the CT get-entries method.
type GetEntriesResponse struct ***REMOVED***
	Entries []LeafEntry `json:"entries"` // the list of returned entries
***REMOVED***

// GetRawEntries exposes the /ct/v1/get-entries result with only the JSON parsing done.
func GetRawEntries(ctx context.Context, httpClient *http.Client, logURL string, start, end int64) (*GetEntriesResponse, error) ***REMOVED***
	if end < 0 ***REMOVED***
		return nil, errors.New("end should be >= 0")
	***REMOVED***
	if end < start ***REMOVED***
		return nil, errors.New("start should be <= end")
	***REMOVED***

	baseURL, err := url.Parse(strings.TrimRight(logURL, "/") + GetEntriesPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	baseURL.RawQuery = url.Values***REMOVED***
		"start": []string***REMOVED***strconv.FormatInt(start, 10)***REMOVED***,
		"end":   []string***REMOVED***strconv.FormatInt(end, 10)***REMOVED***,
	***REMOVED***.Encode()

	var resp GetEntriesResponse
	err = fetchAndParse(context.TODO(), httpClient, baseURL.String(), &resp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &resp, nil
***REMOVED***

// GetEntries attempts to retrieve the entries in the sequence [|start|, |end|] from the CT log server. (see section 4.6.)
// Returns a slice of LeafInputs or a non-nil error.
func (c *LogClient) GetEntries(start, end int64) ([]ct.LogEntry, error) ***REMOVED***
	resp, err := GetRawEntries(context.TODO(), c.httpClient, c.uri, start, end)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	entries := make([]ct.LogEntry, len(resp.Entries))
	for index, entry := range resp.Entries ***REMOVED***
		leaf, err := ct.ReadMerkleTreeLeaf(bytes.NewBuffer(entry.LeafInput))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		entries[index].Leaf = *leaf

		var chain []ct.ASN1Cert
		switch leaf.TimestampedEntry.EntryType ***REMOVED***
		case ct.X509LogEntryType:
			chain, err = ct.UnmarshalX509ChainArray(entry.ExtraData)

		case ct.PrecertLogEntryType:
			chain, err = ct.UnmarshalPrecertChainArray(entry.ExtraData)

		default:
			return nil, fmt.Errorf("saw unknown entry type: %v", leaf.TimestampedEntry.EntryType)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		entries[index].Chain = chain
		entries[index].Index = start + int64(index)
	***REMOVED***
	return entries, nil
***REMOVED***
