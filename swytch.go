package saklient

import "fmt"

type SwytchService struct {
	api *APIService
	basicQuery
}

func newSwytchService(api *APIService) *SwytchService {
	return (&SwytchService{
		api: api,
	}).Reset()
}

func (s *SwytchService) Create() *Swytch {
	return &Swytch{service: s}
}

func (s *SwytchService) Reset() *SwytchService {
	s.basicQuery.Reset()
	return s
}

func (s *SwytchService) GetByID(id string) (*Swytch, error) {
	jsonResp := &struct {
		IsOK    bool    `json:"is_ok"`
		Success bool    `json:"Success"`
		Swytch  *Swytch `json:"Switch"`
	}{
		Swytch: s.Create(),
	}
	err := s.api.client.Request("GET", fmt.Sprintf("switch/%s", id), nil, jsonResp)
	if err != nil {
		return nil, err
	}
	return jsonResp.Swytch, nil
}

type swytchRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

type IPv4Net struct {
	ID             int    `json:"ID"`
	NetworkAddress string `json:"NetworkAddress"`
	NetworkMaskLen int    `json:"NetworkMaskLen"`
	DefaultRoute   string `json:"DefaultRoute"`
	//NextHop        []string          `json:"NextHop"`
	//StaticRoute    []string          `json:"StaticRouter"`
	ServiceClass string            `json:"ServiceClass"`
	IPAddresses  map[string]string `json:"IPAddresses"`
}

type IPv6Net struct {
	ID             int    `json:"ID"`
	ServiceID      string `json:"ServiceID"`
	IPv6Prefix     string `json:"IPv6Prefix"`
	IPv6PrefixLen  int    `json:"IPv6PrefixLen"`
	IPv6PrefixTail string `json:"IPv6PrefixTail"`
	ServiceClass   string `json:"ServiceClass"`
	IPv6Table      struct {
		ID int `json:"ID"`
	} `json:"IPv6Table"`
	NamedIPv6AddrCount int    `json:"NamedIPv6AddrCount"`
	CreatedAt          string `json:"CreatedAt"`
}

type Swytch struct {
	service          *SwytchService `json:"-"`
	ID               string         `json:"ID"`
	Name             string         `json:"Name"`
	Description      string         `json:"Description"`
	ServerCount      int            `json:"ServerCount"`
	ApplianceCount   int            `json:"ApplianceCount"`
	Scope            string         `json:"Scope"`
	UserSubnet       *struct{}      `json:"UserSubnet"`
	HybridConnection *struct{}      `json:"HybridConnection"`
	CreatedAt        string         `json:"CreatedAt"`
	ServiceClass     string         `json:"ServiceClass"`
	Internet         *struct {
		ID            string `json:"ID"`
		Name          string `json:"Name"`
		BandWidthMbps int    `json:"BandWidthMbps"`
		Scope         string `json:"Scope"`
		ServiceClass  string `json:"ServiceClass"`
	} `json:"Internet"`
	IPv4Nets []*IPv4Net `json:"Subnets"`
	IPv6Nets []*IPv6Net `json:"IPv6Nets"`
	Bridge   *struct{}  `json:"Bridge"`
}

func (s *Swytch) Save() error {
	var err error
	if s.ID == "" {
		sr := &swytchRequest{
			Name:        s.Name,
			Description: s.Description,
		}

		req := &struct {
			Swytch *swytchRequest `json:"Switch"`
		}{
			Swytch: sr,
		}
		resp := &struct {
			IsOK    bool    `json:"is_ok"`
			Success bool    `json:"Success"`
			Swytch  *Swytch `json:"Switch"`
		}{
			Swytch: s,
		}
		err = s.client().Request("POST", "switch", req, resp)
	} else {
		panic("Not Implemented")
	}
	return err
}

func (s *Swytch) Destroy() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	err := s.client().Request("DELETE", fmt.Sprintf("switch/%s", s.ID), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Swytch) Reload() error {
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

func (s *Swytch) client() *Client {
	return s.service.api.client
}

func (s *Swytch) Router() (*Router, error) {
	if s.ID == "" {
		return nil, fmt.Errorf("This is not saved")
	}
	if s.Internet.ID == "" {
		return nil, fmt.Errorf("nil")
	}
	return s.service.api.Router.GetByID(s.Internet.ID)
}

func (s *Swytch) AddIPv6Net() (*IPv6Net, error) {
	if s.ID == "" {
		return nil, fmt.Errorf("This is not saved")
	}
	router, err := s.Router()
	if err != nil {
		return nil, err
	}
	ipv6n, err := router.AddIPv6Net()
	if err != nil {
		return nil, err
	}
	err = s.Reload()
	return ipv6n, err
}

func (s *Swytch) RemoveIPv6Net() error {
	if s.ID == "" {
		return fmt.Errorf("This is not saved")
	}
	router, err := s.Router()
	if err != nil {
		return err
	}
	err = router.RemoveIPv6Net()
	if err != nil {
		return err
	}
	return s.Reload()
}
