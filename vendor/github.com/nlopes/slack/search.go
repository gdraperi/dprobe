package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

const (
	DEFAULT_SEARCH_SORT      = "score"
	DEFAULT_SEARCH_SORT_DIR  = "desc"
	DEFAULT_SEARCH_HIGHLIGHT = false
	DEFAULT_SEARCH_COUNT     = 100
	DEFAULT_SEARCH_PAGE      = 1
)

type SearchParameters struct ***REMOVED***
	Sort          string
	SortDirection string
	Highlight     bool
	Count         int
	Page          int
***REMOVED***

type CtxChannel struct ***REMOVED***
	ID   string `json:"id"`
	Name string `json:"name"`
***REMOVED***

type CtxMessage struct ***REMOVED***
	User      string `json:"user"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"`
	Type      string `json:"type"`
***REMOVED***

type SearchMessage struct ***REMOVED***
	Type      string     `json:"type"`
	Channel   CtxChannel `json:"channel"`
	User      string     `json:"user"`
	Username  string     `json:"username"`
	Timestamp string     `json:"ts"`
	Text      string     `json:"text"`
	Permalink string     `json:"permalink"`
	Previous  CtxMessage `json:"previous"`
	Previous2 CtxMessage `json:"previous_2"`
	Next      CtxMessage `json:"next"`
	Next2     CtxMessage `json:"next_2"`
***REMOVED***

type SearchMessages struct ***REMOVED***
	Matches    []SearchMessage `json:"matches"`
	Paging     `json:"paging"`
	Pagination `json:"pagination"`
	Total      int `json:"total"`
***REMOVED***

type SearchFiles struct ***REMOVED***
	Matches    []File `json:"matches"`
	Paging     `json:"paging"`
	Pagination `json:"pagination"`
	Total      int `json:"total"`
***REMOVED***

type searchResponseFull struct ***REMOVED***
	Query          string `json:"query"`
	SearchMessages `json:"messages"`
	SearchFiles    `json:"files"`
	SlackResponse
***REMOVED***

func NewSearchParameters() SearchParameters ***REMOVED***
	return SearchParameters***REMOVED***
		Sort:          DEFAULT_SEARCH_SORT,
		SortDirection: DEFAULT_SEARCH_SORT_DIR,
		Highlight:     DEFAULT_SEARCH_HIGHLIGHT,
		Count:         DEFAULT_SEARCH_COUNT,
		Page:          DEFAULT_SEARCH_PAGE,
	***REMOVED***
***REMOVED***

func (api *Client) _search(ctx context.Context, path, query string, params SearchParameters, files, messages bool) (response *searchResponseFull, error error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"query": ***REMOVED***query***REMOVED***,
	***REMOVED***
	if params.Sort != DEFAULT_SEARCH_SORT ***REMOVED***
		values.Add("sort", params.Sort)
	***REMOVED***
	if params.SortDirection != DEFAULT_SEARCH_SORT_DIR ***REMOVED***
		values.Add("sort_dir", params.SortDirection)
	***REMOVED***
	if params.Highlight != DEFAULT_SEARCH_HIGHLIGHT ***REMOVED***
		values.Add("highlight", strconv.Itoa(1))
	***REMOVED***
	if params.Count != DEFAULT_SEARCH_COUNT ***REMOVED***
		values.Add("count", strconv.Itoa(params.Count))
	***REMOVED***
	if params.Page != DEFAULT_SEARCH_PAGE ***REMOVED***
		values.Add("page", strconv.Itoa(params.Page))
	***REMOVED***

	response = &searchResponseFull***REMOVED******REMOVED***
	err := post(ctx, api.httpclient, path, values, response, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil

***REMOVED***

func (api *Client) Search(query string, params SearchParameters) (*SearchMessages, *SearchFiles, error) ***REMOVED***
	return api.SearchContext(context.Background(), query, params)
***REMOVED***

func (api *Client) SearchContext(ctx context.Context, query string, params SearchParameters) (*SearchMessages, *SearchFiles, error) ***REMOVED***
	response, err := api._search(ctx, "search.all", query, params, true, true)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return &response.SearchMessages, &response.SearchFiles, nil
***REMOVED***

func (api *Client) SearchFiles(query string, params SearchParameters) (*SearchFiles, error) ***REMOVED***
	return api.SearchFilesContext(context.Background(), query, params)
***REMOVED***

func (api *Client) SearchFilesContext(ctx context.Context, query string, params SearchParameters) (*SearchFiles, error) ***REMOVED***
	response, err := api._search(ctx, "search.files", query, params, true, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.SearchFiles, nil
***REMOVED***

func (api *Client) SearchMessages(query string, params SearchParameters) (*SearchMessages, error) ***REMOVED***
	return api.SearchMessagesContext(context.Background(), query, params)
***REMOVED***

func (api *Client) SearchMessagesContext(ctx context.Context, query string, params SearchParameters) (*SearchMessages, error) ***REMOVED***
	response, err := api._search(ctx, "search.messages", query, params, false, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.SearchMessages, nil
***REMOVED***
