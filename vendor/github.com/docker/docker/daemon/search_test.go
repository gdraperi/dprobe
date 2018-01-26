package daemon

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
)

type FakeService struct ***REMOVED***
	registry.DefaultService

	shouldReturnError bool

	term    string
	results []registrytypes.SearchResult
***REMOVED***

func (s *FakeService) Search(ctx context.Context, term string, limit int, authConfig *types.AuthConfig, userAgent string, headers map[string][]string) (*registrytypes.SearchResults, error) ***REMOVED***
	if s.shouldReturnError ***REMOVED***
		return nil, errors.New("Search unknown error")
	***REMOVED***
	return &registrytypes.SearchResults***REMOVED***
		Query:      s.term,
		NumResults: len(s.results),
		Results:    s.results,
	***REMOVED***, nil
***REMOVED***

func TestSearchRegistryForImagesErrors(t *testing.T) ***REMOVED***
	errorCases := []struct ***REMOVED***
		filtersArgs       string
		shouldReturnError bool
		expectedError     string
	***REMOVED******REMOVED***
		***REMOVED***
			expectedError:     "Search unknown error",
			shouldReturnError: true,
		***REMOVED***,
		***REMOVED***
			filtersArgs:   "invalid json",
			expectedError: "invalid character 'i' looking for beginning of value",
		***REMOVED***,
		***REMOVED***
			filtersArgs:   `***REMOVED***"type":***REMOVED***"custom":true***REMOVED******REMOVED***`,
			expectedError: "Invalid filter 'type'",
		***REMOVED***,
		***REMOVED***
			filtersArgs:   `***REMOVED***"is-automated":***REMOVED***"invalid":true***REMOVED******REMOVED***`,
			expectedError: "Invalid filter 'is-automated=[invalid]'",
		***REMOVED***,
		***REMOVED***
			filtersArgs:   `***REMOVED***"is-automated":***REMOVED***"true":true,"false":true***REMOVED******REMOVED***`,
			expectedError: "Invalid filter 'is-automated",
		***REMOVED***,
		***REMOVED***
			filtersArgs:   `***REMOVED***"is-official":***REMOVED***"invalid":true***REMOVED******REMOVED***`,
			expectedError: "Invalid filter 'is-official=[invalid]'",
		***REMOVED***,
		***REMOVED***
			filtersArgs:   `***REMOVED***"is-official":***REMOVED***"true":true,"false":true***REMOVED******REMOVED***`,
			expectedError: "Invalid filter 'is-official",
		***REMOVED***,
		***REMOVED***
			filtersArgs:   `***REMOVED***"stars":***REMOVED***"invalid":true***REMOVED******REMOVED***`,
			expectedError: "Invalid filter 'stars=invalid'",
		***REMOVED***,
		***REMOVED***
			filtersArgs:   `***REMOVED***"stars":***REMOVED***"1":true,"invalid":true***REMOVED******REMOVED***`,
			expectedError: "Invalid filter 'stars=invalid'",
		***REMOVED***,
	***REMOVED***
	for index, e := range errorCases ***REMOVED***
		daemon := &Daemon***REMOVED***
			RegistryService: &FakeService***REMOVED***
				shouldReturnError: e.shouldReturnError,
			***REMOVED***,
		***REMOVED***
		_, err := daemon.SearchRegistryForImages(context.Background(), e.filtersArgs, "term", 25, nil, map[string][]string***REMOVED******REMOVED***)
		if err == nil ***REMOVED***
			t.Errorf("%d: expected an error, got nothing", index)
		***REMOVED***
		if !strings.Contains(err.Error(), e.expectedError) ***REMOVED***
			t.Errorf("%d: expected error to contain %s, got %s", index, e.expectedError, err.Error())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSearchRegistryForImages(t *testing.T) ***REMOVED***
	term := "term"
	successCases := []struct ***REMOVED***
		filtersArgs     string
		registryResults []registrytypes.SearchResult
		expectedResults []registrytypes.SearchResult
	***REMOVED******REMOVED***
		***REMOVED***
			filtersArgs:     "",
			registryResults: []registrytypes.SearchResult***REMOVED******REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: "",
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-automated":***REMOVED***"true":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-automated":***REMOVED***"true":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-automated":***REMOVED***"false":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-automated":***REMOVED***"false":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsAutomated: false,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsAutomated: false,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-official":***REMOVED***"true":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-official":***REMOVED***"true":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-official":***REMOVED***"false":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"is-official":***REMOVED***"false":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsOfficial:  false,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					IsOfficial:  false,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"stars":***REMOVED***"0":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					StarCount:   0,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					StarCount:   0,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"stars":***REMOVED***"1":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name",
					Description: "description",
					StarCount:   0,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"stars":***REMOVED***"1":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name0",
					Description: "description0",
					StarCount:   0,
				***REMOVED***,
				***REMOVED***
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			filtersArgs: `***REMOVED***"stars":***REMOVED***"1":true***REMOVED***, "is-official":***REMOVED***"true":true***REMOVED***, "is-automated":***REMOVED***"true":true***REMOVED******REMOVED***`,
			registryResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name0",
					Description: "description0",
					StarCount:   0,
					IsOfficial:  true,
					IsAutomated: true,
				***REMOVED***,
				***REMOVED***
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
					IsOfficial:  false,
					IsAutomated: true,
				***REMOVED***,
				***REMOVED***
					Name:        "name2",
					Description: "description2",
					StarCount:   1,
					IsOfficial:  true,
					IsAutomated: false,
				***REMOVED***,
				***REMOVED***
					Name:        "name3",
					Description: "description3",
					StarCount:   2,
					IsOfficial:  true,
					IsAutomated: true,
				***REMOVED***,
			***REMOVED***,
			expectedResults: []registrytypes.SearchResult***REMOVED***
				***REMOVED***
					Name:        "name3",
					Description: "description3",
					StarCount:   2,
					IsOfficial:  true,
					IsAutomated: true,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for index, s := range successCases ***REMOVED***
		daemon := &Daemon***REMOVED***
			RegistryService: &FakeService***REMOVED***
				term:    term,
				results: s.registryResults,
			***REMOVED***,
		***REMOVED***
		results, err := daemon.SearchRegistryForImages(context.Background(), s.filtersArgs, term, 25, nil, map[string][]string***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			t.Errorf("%d: %v", index, err)
		***REMOVED***
		if results.Query != term ***REMOVED***
			t.Errorf("%d: expected Query to be %s, got %s", index, term, results.Query)
		***REMOVED***
		if results.NumResults != len(s.expectedResults) ***REMOVED***
			t.Errorf("%d: expected NumResults to be %d, got %d", index, len(s.expectedResults), results.NumResults)
		***REMOVED***
		for _, result := range results.Results ***REMOVED***
			found := false
			for _, expectedResult := range s.expectedResults ***REMOVED***
				if expectedResult.Name == result.Name &&
					expectedResult.Description == result.Description &&
					expectedResult.IsAutomated == result.IsAutomated &&
					expectedResult.IsOfficial == result.IsOfficial &&
					expectedResult.StarCount == result.StarCount ***REMOVED***
					found = true
					break
				***REMOVED***
			***REMOVED***
			if !found ***REMOVED***
				t.Errorf("%d: expected results %v, got %v", index, s.expectedResults, results.Results)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
