package saklient

type ServerService struct {
	api *APIService
}

func (s *ServerService) Create() *Server {
	return &Server{client: s.api.client}
}

type Server struct {
	client      *Client
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

func (s *Server) Save() error {
	req := &struct {
		Server *Server `json:"Server"`
	}{
		Server: s,
	}

	respServer := new(ServerResponse)
	err := s.client.Request("POST", "server", req, respServer)
	return err
}
