package saklient

import (
	"encoding/base64"
	"fmt"
	"net/http"

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

type ServerPlan struct {
	ID int `json:"ID"`
}

type Server struct {
	api               *APIService
	Name              string `json:"Name"`
	Description       string `json:"Description,omitempty"`
	PlanID            int
	Plan              ServerPlan    `json:"Plan"`
	Tags              []string      `json:"Tags"`
	ConnectedSwitches []interface{} `json:"ConnectedSwitches"`
}

type APIError struct {
}

type ServerResponse struct {
}

func (s *Server) Save() (*ServerResponse, *http.Response, error) {
	req := struct {
		Server *Server `json:"Server"`
		Count  int     `json:"Count"`
	}{
		Server: s,
	}

	respServer := new(ServerResponse)
	apiErr := new(APIError)
	resp, err := s.api.client.NewSling().Post("server").BodyJSON(&req).Receive(respServer, apiErr)
	return respServer, resp, err
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

type ServerService struct {
	api *APIService
}

func (s *ServerService) Create() (*Server, error) {
	return &Server{api: s.api}, nil
}
