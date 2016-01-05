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
		Swytch  *Swytch `json:"Swytch"`
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

type Swytch struct {
	service     *SwytchService `json:"-"`
	ID          string         `json:"ID"`
	Name        string         `json:"Name"`
	Description string         `json:"Description"`
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
