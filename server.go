package saklient

import "fmt"

type ServerService struct {
	api *APIService
	basicQuery
}

func newServerService(api *APIService) *ServerService {
	return (&ServerService{
		api: api,
	}).Reset()
}

func (s *ServerService) Create() *Server {
	return &Server{service: s}
}

func (s *ServerService) Reset() *ServerService {
	s.basicQuery.Reset()
	return s
}

func (s *ServerService) GetByID(id string) (*Server, error) {
	jsonResp := &struct {
		IsOK    bool    `json:"is_ok"`
		Success bool    `json:"Success"`
		Server  *Server `json:"Server"`
	}{
		Server: s.Create(),
	}
	err := s.api.client.Request("GET", fmt.Sprintf("server/%s", id), nil, jsonResp)
	if err != nil {
		return nil, err
	}
	return jsonResp.Server, nil
}

type serverRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	ServerPlan  struct {
		ID int `json:"ID"`
	} `json:"ServerPlan"`
	Tags              []string      `json:"Tags,omitempty"`
	ConnectedSwitches []interface{} `json:"ConnectedSwitches,omitempty"`
}

type Server struct {
	service     *ServerService `json:"-"`
	ID          string         `json:"ID"`
	Name        string         `json:"Name"`
	Description string         `json:"Description"`
	ServerPlan  struct {
		ID           int    `json:"ID"`
		Name         string `json:"Name"`
		CPU          int    `json:"CPU"`
		MemoryMB     int    `json:"MemoryMB"`
		ServiceClass string `json:"ServiceClass"`
	} `json:"ServerPlan"`
	Tags              []string      `json:"Tags"`
	ConnectedSwitches []interface{} `json:"ConnectedSwitches"`
}

func (s *Server) Save() error {
	var err error
	if s.ID == "" {
		sr := &serverRequest{
			Name: s.Name,
		}
		sr.ServerPlan.ID = s.ServerPlan.ID

		req := &struct {
			Server *serverRequest `json:"Server"`
		}{
			Server: sr,
		}
		resp := &struct {
			IsOK    bool    `json:"is_ok"`
			Success bool    `json:"Success"`
			Server  *Server `json:"Server"`
		}{
			Server: s,
		}
		err = s.client().Request("POST", "server", req, resp)
	} else {
		panic("Not Implemented")
	}
	return err
}

func (s *Server) Destroy() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	err := s.client().Request("DELETE", fmt.Sprintf("server/%s", s.ID), nil, nil)
	if err != nil {
		return nil
	}
	return nil
}

func (s *Server) Reload() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	n, err := s.service.GetByID(s.ID)
	if err != nil {
		return err
	}
	*s = *n
	return nil
}

func (s *Server) client() *Client {
	return s.service.api.client
}
