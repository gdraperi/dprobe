package slack

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
)

const (
	// Add here the defaults in the siten
	DEFAULT_FILES_USER    = ""
	DEFAULT_FILES_CHANNEL = ""
	DEFAULT_FILES_TS_FROM = 0
	DEFAULT_FILES_TS_TO   = -1
	DEFAULT_FILES_TYPES   = "all"
	DEFAULT_FILES_COUNT   = 100
	DEFAULT_FILES_PAGE    = 1
)

// File contains all the information for a file
type File struct ***REMOVED***
	ID        string   `json:"id"`
	Created   JSONTime `json:"created"`
	Timestamp JSONTime `json:"timestamp"`

	Name              string `json:"name"`
	Title             string `json:"title"`
	Mimetype          string `json:"mimetype"`
	ImageExifRotation int    `json:"image_exif_rotation"`
	Filetype          string `json:"filetype"`
	PrettyType        string `json:"pretty_type"`
	User              string `json:"user"`

	Mode         string `json:"mode"`
	Editable     bool   `json:"editable"`
	IsExternal   bool   `json:"is_external"`
	ExternalType string `json:"external_type"`

	Size int `json:"size"`

	URL                string `json:"url"`          // Deprecated - never set
	URLDownload        string `json:"url_download"` // Deprecated - never set
	URLPrivate         string `json:"url_private"`
	URLPrivateDownload string `json:"url_private_download"`

	OriginalH   int    `json:"original_h"`
	OriginalW   int    `json:"original_w"`
	Thumb64     string `json:"thumb_64"`
	Thumb80     string `json:"thumb_80"`
	Thumb160    string `json:"thumb_160"`
	Thumb360    string `json:"thumb_360"`
	Thumb360Gif string `json:"thumb_360_gif"`
	Thumb360W   int    `json:"thumb_360_w"`
	Thumb360H   int    `json:"thumb_360_h"`
	Thumb480    string `json:"thumb_480"`
	Thumb480W   int    `json:"thumb_480_w"`
	Thumb480H   int    `json:"thumb_480_h"`
	Thumb720    string `json:"thumb_720"`
	Thumb720W   int    `json:"thumb_720_w"`
	Thumb720H   int    `json:"thumb_720_h"`
	Thumb960    string `json:"thumb_960"`
	Thumb960W   int    `json:"thumb_960_w"`
	Thumb960H   int    `json:"thumb_960_h"`
	Thumb1024   string `json:"thumb_1024"`
	Thumb1024W  int    `json:"thumb_1024_w"`
	Thumb1024H  int    `json:"thumb_1024_h"`

	Permalink       string `json:"permalink"`
	PermalinkPublic string `json:"permalink_public"`

	EditLink         string `json:"edit_link"`
	Preview          string `json:"preview"`
	PreviewHighlight string `json:"preview_highlight"`
	Lines            int    `json:"lines"`
	LinesMore        int    `json:"lines_more"`

	IsPublic        bool     `json:"is_public"`
	PublicURLShared bool     `json:"public_url_shared"`
	Channels        []string `json:"channels"`
	Groups          []string `json:"groups"`
	IMs             []string `json:"ims"`
	InitialComment  Comment  `json:"initial_comment"`
	CommentsCount   int      `json:"comments_count"`
	NumStars        int      `json:"num_stars"`
	IsStarred       bool     `json:"is_starred"`
***REMOVED***

// FileUploadParameters contains all the parameters necessary (including the optional ones) for an UploadFile() request.
//
// There are three ways to upload a file. You can either set Content if file is small, set Reader if file is large,
// or provide a local file path in File to upload it from your filesystem.
type FileUploadParameters struct ***REMOVED***
	File           string
	Content        string
	Reader         io.Reader
	Filetype       string
	Filename       string
	Title          string
	InitialComment string
	Channels       []string
***REMOVED***

// GetFilesParameters contains all the parameters necessary (including the optional ones) for a GetFiles() request
type GetFilesParameters struct ***REMOVED***
	User          string
	Channel       string
	TimestampFrom JSONTime
	TimestampTo   JSONTime
	Types         string
	Count         int
	Page          int
***REMOVED***

type fileResponseFull struct ***REMOVED***
	File     `json:"file"`
	Paging   `json:"paging"`
	Comments []Comment `json:"comments"`
	Files    []File    `json:"files"`

	SlackResponse
***REMOVED***

// NewGetFilesParameters provides an instance of GetFilesParameters with all the sane default values set
func NewGetFilesParameters() GetFilesParameters ***REMOVED***
	return GetFilesParameters***REMOVED***
		User:          DEFAULT_FILES_USER,
		Channel:       DEFAULT_FILES_CHANNEL,
		TimestampFrom: DEFAULT_FILES_TS_FROM,
		TimestampTo:   DEFAULT_FILES_TS_TO,
		Types:         DEFAULT_FILES_TYPES,
		Count:         DEFAULT_FILES_COUNT,
		Page:          DEFAULT_FILES_PAGE,
	***REMOVED***
***REMOVED***

func fileRequest(ctx context.Context, client HTTPRequester, path string, values url.Values, debug bool) (*fileResponseFull, error) ***REMOVED***
	response := &fileResponseFull***REMOVED******REMOVED***
	err := postForm(ctx, client, SLACK_API+path, values, response, debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return response, nil
***REMOVED***

// GetFileInfo retrieves a file and related comments
func (api *Client) GetFileInfo(fileID string, count, page int) (*File, []Comment, *Paging, error) ***REMOVED***
	return api.GetFileInfoContext(context.Background(), fileID, count, page)
***REMOVED***

// GetFileInfoContext retrieves a file and related comments with a custom context
func (api *Client) GetFileInfoContext(ctx context.Context, fileID string, count, page int) (*File, []Comment, *Paging, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"file":  ***REMOVED***fileID***REMOVED***,
		"count": ***REMOVED***strconv.Itoa(count)***REMOVED***,
		"page":  ***REMOVED***strconv.Itoa(page)***REMOVED***,
	***REMOVED***

	response, err := fileRequest(ctx, api.httpclient, "files.info", values, api.debug)
	if err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***
	return &response.File, response.Comments, &response.Paging, nil
***REMOVED***

// GetFiles retrieves all files according to the parameters given
func (api *Client) GetFiles(params GetFilesParameters) ([]File, *Paging, error) ***REMOVED***
	return api.GetFilesContext(context.Background(), params)
***REMOVED***

// GetFilesContext retrieves all files according to the parameters given with a custom context
func (api *Client) GetFilesContext(ctx context.Context, params GetFilesParameters) ([]File, *Paging, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if params.User != DEFAULT_FILES_USER ***REMOVED***
		values.Add("user", params.User)
	***REMOVED***
	if params.Channel != DEFAULT_FILES_CHANNEL ***REMOVED***
		values.Add("channel", params.Channel)
	***REMOVED***
	if params.TimestampFrom != DEFAULT_FILES_TS_FROM ***REMOVED***
		values.Add("ts_from", strconv.FormatInt(int64(params.TimestampFrom), 10))
	***REMOVED***
	if params.TimestampTo != DEFAULT_FILES_TS_TO ***REMOVED***
		values.Add("ts_to", strconv.FormatInt(int64(params.TimestampTo), 10))
	***REMOVED***
	if params.Types != DEFAULT_FILES_TYPES ***REMOVED***
		values.Add("types", params.Types)
	***REMOVED***
	if params.Count != DEFAULT_FILES_COUNT ***REMOVED***
		values.Add("count", strconv.Itoa(params.Count))
	***REMOVED***
	if params.Page != DEFAULT_FILES_PAGE ***REMOVED***
		values.Add("page", strconv.Itoa(params.Page))
	***REMOVED***

	response, err := fileRequest(ctx, api.httpclient, "files.list", values, api.debug)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return response.Files, &response.Paging, nil
***REMOVED***

// UploadFile uploads a file
func (api *Client) UploadFile(params FileUploadParameters) (file *File, err error) ***REMOVED***
	return api.UploadFileContext(context.Background(), params)
***REMOVED***

// UploadFileContext uploads a file and setting a custom context
func (api *Client) UploadFileContext(ctx context.Context, params FileUploadParameters) (file *File, err error) ***REMOVED***
	// Test if user token is valid. This helps because client.Do doesn't like this for some reason. XXX: More
	// investigation needed, but for now this will do.
	_, err = api.AuthTest()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	response := &fileResponseFull***REMOVED******REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
	***REMOVED***
	if params.Filetype != "" ***REMOVED***
		values.Add("filetype", params.Filetype)
	***REMOVED***
	if params.Filename != "" ***REMOVED***
		values.Add("filename", params.Filename)
	***REMOVED***
	if params.Title != "" ***REMOVED***
		values.Add("title", params.Title)
	***REMOVED***
	if params.InitialComment != "" ***REMOVED***
		values.Add("initial_comment", params.InitialComment)
	***REMOVED***
	if len(params.Channels) != 0 ***REMOVED***
		values.Add("channels", strings.Join(params.Channels, ","))
	***REMOVED***
	if params.Content != "" ***REMOVED***
		values.Add("content", params.Content)
		err = postForm(ctx, api.httpclient, SLACK_API+"files.upload", values, response, api.debug)
	***REMOVED*** else if params.File != "" ***REMOVED***
		err = postLocalWithMultipartResponse(ctx, api.httpclient, SLACK_API+"files.upload", params.File, "file", values, response, api.debug)
	***REMOVED*** else if params.Reader != nil ***REMOVED***
		err = postWithMultipartResponse(ctx, api.httpclient, SLACK_API+"files.upload", params.Filename, "file", values, params.Reader, response, api.debug)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !response.Ok ***REMOVED***
		return nil, errors.New(response.Error)
	***REMOVED***
	return &response.File, nil
***REMOVED***

// DeleteFile deletes a file
func (api *Client) DeleteFile(fileID string) error ***REMOVED***
	return api.DeleteFileContext(context.Background(), fileID)
***REMOVED***

// DeleteFileContext deletes a file with a custom context
func (api *Client) DeleteFileContext(ctx context.Context, fileID string) (err error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"file":  ***REMOVED***fileID***REMOVED***,
	***REMOVED***

	if _, err = fileRequest(ctx, api.httpclient, "files.delete", values, api.debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// RevokeFilePublicURL disables public/external sharing for a file
func (api *Client) RevokeFilePublicURL(fileID string) (*File, error) ***REMOVED***
	return api.RevokeFilePublicURLContext(context.Background(), fileID)
***REMOVED***

// RevokeFilePublicURLContext disables public/external sharing for a file with a custom context
func (api *Client) RevokeFilePublicURLContext(ctx context.Context, fileID string) (*File, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"file":  ***REMOVED***fileID***REMOVED***,
	***REMOVED***

	response, err := fileRequest(ctx, api.httpclient, "files.revokePublicURL", values, api.debug)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &response.File, nil
***REMOVED***

// ShareFilePublicURL enabled public/external sharing for a file
func (api *Client) ShareFilePublicURL(fileID string) (*File, []Comment, *Paging, error) ***REMOVED***
	return api.ShareFilePublicURLContext(context.Background(), fileID)
***REMOVED***

// ShareFilePublicURLContext enabled public/external sharing for a file with a custom context
func (api *Client) ShareFilePublicURLContext(ctx context.Context, fileID string) (*File, []Comment, *Paging, error) ***REMOVED***
	values := url.Values***REMOVED***
		"token": ***REMOVED***api.token***REMOVED***,
		"file":  ***REMOVED***fileID***REMOVED***,
	***REMOVED***

	response, err := fileRequest(ctx, api.httpclient, "files.sharedPublicURL", values, api.debug)
	if err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***
	return &response.File, response.Comments, &response.Paging, nil
***REMOVED***
