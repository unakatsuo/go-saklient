package saklient

import "fmt"

type IfaceService struct {
	api *APIService
	basicQuery
}

func newIfaceService(api *APIService) *IfaceService {
	return (&IfaceService{
		api: api,
	}).Reset()
}

func (s *IfaceService) Create() *Iface {
	return &Iface{service: s}
}

func (s *IfaceService) Reset() *IfaceService {
	s.basicQuery.Reset()
	return s
}

func (s *IfaceService) GetByID(id string) (*Iface, error) {
	jsonResp := &struct {
		IsOK    bool   `json:"is_ok"`
		Success bool   `json:"Success"`
		Iface   *Iface `json:"Iface"`
	}{
		Iface: s.Create(),
	}
	err := s.api.client.Request("GET", fmt.Sprintf("interface/%s", id), nil, jsonResp)
	if err != nil {
		return nil, err
	}
	return jsonResp.Iface, nil
}

type ifaceRequest struct {
	Server struct {
		ID string `json:"ID"`
	} `json:"Server"`
}

type Iface struct {
	service       *IfaceService `json:"-"`
	ID            string        `json:"ID"`
	MACAddress    string        `json:"MACAddress"`
	IPAddress     string        `json:"IPAddress"`
	UserIPAddress string        `json:"UserIPAddress"`
	Server        struct {
		ID string `json:"ID"`
	} `json:"Server"`
}

func (s *Iface) Save() error {
	var err error
	if s.ID == "" {
		sr := &ifaceRequest{}
		sr.Server.ID = s.Server.ID
		req := &struct {
			Iface *ifaceRequest `json:"Interface"`
		}{
			Iface: sr,
		}
		resp := &struct {
			IsOK    bool   `json:"is_ok"`
			Success bool   `json:"Success"`
			Iface   *Iface `json:"Interface"`
		}{
			Iface: s,
		}
		err = s.client().Request("POST", "interface", req, resp)
	} else {
		panic("Not Implemented")
	}
	return err
}

func (s *Iface) Destroy() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	err := s.client().Request("DELETE", fmt.Sprintf("interface/%s", s.ID), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Iface) Reload() error {
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

func (s *Iface) client() *Client {
	return s.service.api.client
}

func (s *Iface) ConnectToSharedSegment() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	return s.client().Request("PUT", fmt.Sprintf("interface/%s/to/switch/shared", s.ID), nil, nil)
}

func (s *Iface) DisconnectFromSwytch() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	return s.client().Request("DELETE", fmt.Sprintf("interface/%s/to/switch", s.ID), nil, nil)
}

func (s *Iface) ConnectToSwytch(swytch *Swytch) error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	if swytch.ID == "" {
		return fmt.Errorf("swytch is invalid")
	}
	return s.client().Request("PUT", fmt.Sprintf("interface/%s/to/switch/%s", s.ID, swytch.ID), nil, nil)
}
