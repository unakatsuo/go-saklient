package saklient

import (
	"encoding/base64"
	"fmt"

	"github.com/dghubble/sling"
)

type Client struct {
	sling *sling.Sling
}

func (c *Client) NewSling() *sling.Sling {
	return c.sling.New()
}

func newClient(token string, secret string) *Client {
	basicToken := fmt.Sprintf("%s:%s", token, secret)
	sl := sling.New().Base("https://secure.sakura.ad.jp/cloud/zone/tk1v/api/cloud/1.1/").Add(
		"Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(basicToken))))

	c := &Client{
		sling: sl,
	}

	return c
}

type APIError struct {
}

type APIService struct {
	client *Client
	Server *ServerService
}

type AuthorizeFunc func(token string, secret string) *APIService

func BasicAuthorize(token string, secret string) *APIService {
	api := &APIService{
		client: newClient(token, secret),
	}
	api.Server = &ServerService{api: api}
	return api
}

var API struct {
	Authorize AuthorizeFunc
}

func init() {
	API.Authorize = BasicAuthorize
}
