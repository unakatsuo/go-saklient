package saklient

import "net/http"

type ServerService struct {
	api *APIService
}

func (s *ServerService) Create() (*Server, error) {
	return &Server{api: s.api}, nil
}

type Server struct {
	api         *APIService
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	Plan        struct {
		ID int `json:"ID"`
	} `json:"Plan"`
	Tags              []string      `json:"Tags"`
	ConnectedSwitches []interface{} `json:"ConnectedSwitches"`
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
