package saklient

import "fmt"

type RouterService struct {
	api *APIService
	basicQuery
}

func newRouterService(api *APIService) *RouterService {
	return (&RouterService{
		api: api,
	}).Reset()
}

func (s *RouterService) Create() *Router {
	return &Router{service: s}
}

func (s *RouterService) Reset() *RouterService {
	s.basicQuery.Reset()
	return s
}

func (s *RouterService) GetByID(id string) (*Router, error) {
	jsonResp := &struct {
		IsOK    bool    `json:"is_ok"`
		Success bool    `json:"Success"`
		Router  *Router `json:"Internet"`
	}{
		Router: s.Create(),
	}
	err := s.api.client.Request("GET", fmt.Sprintf("internet/%s", id), nil, jsonResp)
	if err != nil {
		return nil, err
	}
	return jsonResp.Router, nil
}

type routerRequest struct {
	Name           string `json:"Name"`
	Description    string `json:"Description,omitempty"`
	NetworkMaskLen int    `json:"NetworkMaskLen"`
	BandWidthMbps  int    `json:"BandWidthMbps"`
}

type Router struct {
	service        *RouterService `json:"-"`
	ID             string         `json:"ID"`
	Name           string         `json:"Name"`
	Description    string         `json:"Description"`
	NetworkMaskLen int            `json:"NetworkMaskLen"`
	BandWidthMbps  int            `json:"BandWidthMbps"`
	Scope          string         `json:"Scope"`
	ServiceClass   string         `json:"ServiceClass"`
	Swytch         *Swytch        `json:"Switch"`
}

func (s *Router) Save() error {
	var err error
	if s.ID == "" {
		sr := &routerRequest{
			Name:           s.Name,
			NetworkMaskLen: s.NetworkMaskLen,
			BandWidthMbps:  s.BandWidthMbps,
		}

		req := &struct {
			Router *routerRequest `json:"Internet"`
		}{
			Router: sr,
		}
		resp := &struct {
			IsOK    bool    `json:"is_ok"`
			Success bool    `json:"Success"`
			Router  *Router `json:"Internet"`
		}{
			Router: s,
		}
		err = s.client().Request("POST", "internet", req, resp)
	} else {
		panic("Not Implemented")
	}
	return err
}

func (s *Router) Destroy() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	err := s.client().Request("DELETE", fmt.Sprintf("internet/%s", s.ID), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Router) Reload() error {
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

func (s *Router) client() *Client {
	return s.service.api.client
}

func (l *Router) Exists() (bool, error) {
	if l.ID == "" {
		return false, fmt.Errorf("This is not saved yet")
	}
	query := basicQuery{
		Filter:  map[string]interface{}{"ID": l.ID},
		Include: []string{"ID"},
	}
	resp := &struct {
		Count int `json:"Count"`
	}{}
	err := l.client().Request("GET", "internet", query, resp)
	if err != nil {
		return false, err
	}
	return (resp.Count == 1), nil
}

func (l *Router) SleepWhileCreating() error {
	if l.ID == "" {
		return fmt.Errorf("This is not saved yet")
	}
	return sleepUntil(func() bool {
		ok, err := l.Exists()
		if err == nil && ok {
			l.Reload()
			return true
		}
		return false
	})
}

func (l *Router) GetSwytch() (*Swytch, error) {
	if l.ID == "" {
		return nil, fmt.Errorf("This is not saved yet")
	}
	if l.Swytch == nil || l.Swytch.ID == "" {
		return nil, fmt.Errorf("Invalid Switch")
	}
	sw, err := l.service.api.Swytch.GetByID(l.Swytch.ID)
	if err != nil {
		return nil, err
	}
	l.Swytch = sw
	return sw, nil
}

func (l *Router) AddIPv6Net() (*IPv6Net, error) {
	if l.ID == "" {
		return nil, fmt.Errorf("This is not saved yet")
	}
	resp := &struct {
		IsOK    bool     `json:"is_ok"`
		IPv6Net *IPv6Net `json:"IPv6Net"`
	}{}
	err := l.client().Request("POST", fmt.Sprintf("internet/%s/ipv6net", l.ID), nil, resp)
	if err != nil {
		return nil, err
	}
	err = l.Reload()
	if err != nil {
		return nil, err
	}
	return resp.IPv6Net, nil
}

func (l *Router) RemoveIPv6Net() error {
	if l.ID == "" {
		return fmt.Errorf("This is not saved yet")
	}
	if l.Swytch == nil || l.Swytch.ID == "" {
		return fmt.Errorf("Invalid .Swytch")
	}
	if len(l.Swytch.IPv6Nets) < 1 {
		return fmt.Errorf("No IPv6 network assignment")
	}
	err := l.client().Request("DELETE", fmt.Sprintf("internet/%s/ipv6net/%d", l.ID, l.Swytch.IPv6Nets[0].ID), nil, nil)
	if err != nil {
		return err
	}
	return l.Reload()
}
