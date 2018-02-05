package slack

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

const (
	DEFAULT_LOGINS_COUNT = 100
	DEFAULT_LOGINS_PAGE  = 1
)

type TeamResponse struct ***REMOVED***
	Team TeamInfo `json:"team"`
	SlackResponse
***REMOVED***

type TeamInfo struct ***REMOVED***
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	EmailDomain string                 `json:"email_domain"`
	Icon        map[string]interface***REMOVED******REMOVED*** `json:"icon"`
***REMOVED***

type LoginResponse struct ***REMOVED***
	Logins []Login `json:"logins"`
	Paging `json:"paging"`
	SlackResponse
***REMOVED***

type Login struct ***REMOVED***
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	DateFirst int    `json:"date_first"`
	DateLast  int    `json:"date_last"`
	Count     int    `json:"count"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	ISP       string `json:"isp"`
	Country   string `json:"country"`
	Region    string `json:"region"`
***REMOVED***

type BillableInfoResponse struct ***REMOVED***
	BillableInfo map[string]BillingActive `json:"billable_info"`
	SlackResponse
***REMOVED***

type BillingActive struct ***REMOVED***
	BillingActive bool `json:"billing_active"`
***REMOVED***

// AccessLogParameters contains all the parameters necessary (including the optional ones) for a GetAccessLogs() request
type AccessLogParameters struct ***REMOVED***
	Count int
	Page  int
***REMOVED***

// NewAccessLogParameters provides an instance of AccessLogParameters with all the sane default values set
func NewAccessLogParameters() AccessLogParameters ***REMOVED***
	return AccessLogParameters***REMOVED***
		Count: DEFAULT_LOGINS_COUNT,
		Page:  DEFAULT_LOGINS_PAGE,
	***REMOVED***
***REMOVED***

func teamRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*TeamResponse, error) ***REMOVED***
	response := &TeamResponse***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***

	return response, nil
***REMOVED***

func billableInfoRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (map[string]BillingActive, error) ***REMOVED***
	response := &BillableInfoResponse***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***

	return response.BillableInfo, nil
***REMOVED***

func accessLogsRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*LoginResponse, error) ***REMOVED***
	response := &LoginResponse***REMOVED******REMOVED***
	err := post(ctx, client, path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// GetTeamInfo gets the Team Information of the user
func (api *Client) GetTeamInfo() (*TeamInfo, error) ***REMOVED***
	return api.GetTeamInfoContext(context.Background())
***REMOVED***

// GetTeamInfoContext gets the Team Information of the user with a custom context
func (api *Client) GetTeamInfoContext(ctx context.Context) (*TeamInfo, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***

	response, err := teamRequest(ctx, api.httpclient, "team.info", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.Team, nil
***REMOVED***

// GetAccessLogs retrieves a page of logins according to the parameters given
func (api *Client) GetAccessLogs(params AccessLogParameters) ([]Login, *Paging, error) ***REMOVED***
	return api.GetAccessLogsContext(context.Background(), params)
***REMOVED***

// GetAccessLogsContext retrieves a page of logins according to the parameters given with a custom context
func (api *Client) GetAccessLogsContext(ctx context.Context, params AccessLogParameters) ([]Login, *Paging, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if params.Count != DEFAULT_LOGINS_COUNT ***REMOVED***
		values.Add("count", strconv.Itoa(params.Count))
	***REMOVED***
	if params.Page != DEFAULT_LOGINS_PAGE ***REMOVED***
		values.Add("page", strconv.Itoa(params.Page))
	***REMOVED***

	response, err := accessLogsRequest(ctx, api.httpclient, "team.accessLogs", values, api.debug)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return response.Logins, &response.Paging, nil
***REMOVED***

func (api *Client) GetBillableInfo(user string) (map[string]BillingActive, error) ***REMOVED***
	return api.GetBillableInfoContext(context.Background(), user)
***REMOVED***

func (api *Client) GetBillableInfoContext(ctx context.Context, user string) (map[string]BillingActive, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"user":  ***REMOVED***user***REMOVED***,
	***REMOVED***

	return billableInfoRequest(ctx, api.httpclient, "team.billableInfo", values, api.debug)
***REMOVED***

// GetBillableInfoForTeam returns the billing_active status of all users on the team.
func (api *Client) GetBillableInfoForTeam() (map[string]BillingActive, error) ***REMOVED***
	return api.GetBillableInfoForTeamContext(context.Background())
***REMOVED***

// GetBillableInfoForTeamContext returns the billing_active status of all users on the team with a custom context
func (api *Client) GetBillableInfoForTeamContext(ctx context.Context) (map[string]BillingActive, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***

	return billableInfoRequest(ctx, api.httpclient, "team.billableInfo", values, api.debug)
***REMOVED***
