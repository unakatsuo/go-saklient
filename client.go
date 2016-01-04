package saklient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient  *http.Client
	BaseURL     *url.URL
	Headers     http.Header
	DebugDumper io.Writer
}

// Allow to set custom http.Client
func (c *Client) HttpClient(client *http.Client) *Client {
	c.httpClient = client
	return c
}

func newClient(token string, secret string) *Client {
	basicToken := fmt.Sprintf("%s:%s", token, secret)
	c := &Client{
		httpClient: http.DefaultClient,
		Headers:    make(http.Header),
	}
	c.BaseURL, _ = url.Parse("https://secure.sakura.ad.jp/cloud/zone/tk1v/api/cloud/1.1/")
	c.Headers.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(basicToken))))
	return c
}

func (c *Client) Request(method string, path string, reqParams interface{}, respSuccess interface{}) error {
	pathURL, err := url.Parse(path)
	if err != nil {
		return err
	}
	destURL := c.BaseURL.ResolveReference(pathURL)

	jsonReq := new(bytes.Buffer)
	if reqParams != nil {
		err = json.NewEncoder(jsonReq).Encode(reqParams)
		if err != nil {
			return err
		}
	}

	// Build http.Request
	var reqBody io.Reader
	if method == "GET" {
		// Embed JSON string at query string.
		// ** url.QueryEscape() converts " " to + but
		// API Server does not recognize it so it had to be changed to "%20".
		// This is important for filter queries like .WithNameLike().
		destURL.RawQuery = strings.Replace(
			url.QueryEscape(jsonReq.String()),
			"+", "%20", -1)
	} else {
		reqBody = jsonReq
	}

	httpReq, err := http.NewRequest(method, destURL.String(), reqBody)
	if err != nil {
		return err
	}
	httpReq.Header = c.Headers
	if c.DebugDumper != nil {
		buf, err := httputil.DumpRequestOut(httpReq, true)
		if err != nil {
			return err
		}
		c.DebugDumper.Write(buf)
	}
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if c.DebugDumper != nil {
		buf, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		c.DebugDumper.Write(buf)
	}
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		if code := resp.StatusCode; 200 <= code && code < 300 {
			if respSuccess != nil {
				err = json.NewDecoder(resp.Body).Decode(respSuccess)
			}
		} else {
			apiErr := new(APIError)
			err = json.NewDecoder(resp.Body).Decode(apiErr)
			if err != nil {
				return err
			}
			apiErr.HTTPCode = resp.StatusCode
			err = apiErr
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type APIError struct {
	HTTPCode  int
	Fatal     bool   `json:"is_fatal"`
	Serial    string `json:"serial"`
	Status    string `json:"status"`
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP Status Code: %d, Error Code: %s, Error Message: %s",
		e.HTTPCode, e.ErrorCode, e.ErrorMsg)
}

type APIService struct {
	client  *Client
	Server  *ServerService
	Disk    *DiskService
	Archive *ArchiveService
	Iface   *IfaceService
}

type AuthorizeFunc func(token string, secret string) *APIService

func BasicAuthorize(token string, secret string) *APIService {
	api := &APIService{
		client: newClient(token, secret),
	}
	api.Server = newServerService(api)
	api.Disk = newDiskService(api)
	api.Archive = newArchiveService(api)
	api.Iface = newIfaceService(api)
	return api
}

var API struct {
	Authorize AuthorizeFunc
}

func init() {
	API.Authorize = BasicAuthorize
}

// Handle common filter operations for collection API.
// i.e. "GET /cloud/1.1/server" or "GET /cloud/1.1/archive"
type basicQuery struct {
	Offset   int                    `json:"From,omitempty"`
	Limit    int                    `json:"Count,omitempty"`
	tags     []string               `json:"-"`
	Filter   map[string]interface{} `json:"Filter,omitempty"`
	SortKeys []string               `json:"Sort,omitempty"`
}

func (b *basicQuery) SortBy(key string, reverse bool) {
	if reverse {
		key = "-" + key
	}
	b.SortKeys = append(b.SortKeys, key)
}

func (b *basicQuery) Reset() {
	b.Offset = 0
	b.Limit = 0
	b.Filter = make(map[string]interface{}, 0)
	b.SortKeys = make([]string, 0)
	b.tags = make([]string, 0)
}

func (b *basicQuery) WithTag(tag string) {
	b.tags = append(b.tags, tag)
}

func (b *basicQuery) WithTags(tags []string) {
	for _, t := range tags {
		b.tags = append(b.tags, t)
	}
}

type SourceResource interface {
	SourceID() string
}

type subResourceID struct {
	ID string `json:"ID"`
}

func sleepUntil(eval func() bool) error {
	sleepSec := 2 * time.Second
	timeout := 5 * time.Minute
	startAt := time.Now()
	for i := 0; time.Since(startAt) <= timeout; i++ {
		if eval() {
			return nil
		}
		time.Sleep(sleepSec)
	}

	return fmt.Errorf("Timed out")
}
