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
}

type AuthorizeFunc func(token string, secret string) *APIService

func BasicAuthorize(token string, secret string) *APIService {
	api := &APIService{
		client: newClient(token, secret),
	}
	api.Server = &ServerService{api: api}
	api.Disk = &DiskService{
		api: api,
	}
	api.Archive = &ArchiveService{api: api}
	api.Archive.Reset()
	return api
}

var API struct {
	Authorize AuthorizeFunc
}

func init() {
	API.Authorize = BasicAuthorize
}
