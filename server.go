package saklient

import (
	"fmt"
	"time"
)

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
	Instance *struct {
		Server struct {
			ID string `json:"ID"`
		} `json:"Server"`
		Status          string `json:"Status"`
		BeforeStatus    string `json:"BeforeStatus"`
		StatusChangedAt string `json:"StatusChangedAt"`
		Host            struct {
			Name string `json:"Name"`
		} `json:"Host"`
		CDROM        *struct{} `json:"CDROM"`
		CDROMStroage *struct{} `json:"CDROMStorage"`
	} `json:"Instance"`
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
		return err
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

func (s *Server) Boot() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	return s.client().Request("PUT", fmt.Sprintf("server/%s/power", s.ID), nil, nil)
}

func (s *Server) Shutdown() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	return s.client().Request("DELETE", fmt.Sprintf("server/%s/power", s.ID), nil, nil)
}

func (s *Server) InstanceStatus() string {
	if s.Instance == nil {
		return ""
	}
	return s.Instance.Status
}

func (s *Server) SleepUntilUp() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	var err error
	for i := 0; i < 1000; i++ {
		if i > 0 {
			time.Sleep(1 * time.Second)
			err = s.Reload()
			if err != nil {
				continue
			}
		}
		if s.InstanceStatus() == "up" {
			return nil
		}
	}
	return err
}
