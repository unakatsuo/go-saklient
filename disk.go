package saklient

import (
	"fmt"
	"net/http"
)

type DiskService struct {
	api    *APIService
	offset int
	limit  int
}

func (l *DiskService) Offset(offset int) *DiskService {
	l.offset = offset
	return l
}

func (l *DiskService) Limit(limit int) *DiskService {
	l.limit = limit
	return l
}

func (l *DiskService) Reset() *DiskService {
	l.limit = 0
	l.offset = 0
	return l
}

func (l *DiskService) Create() *Disk {
	return &Disk{
		client: l.api.client,
	}
}

func (l *DiskService) GetByID(id string) (*Disk, error) {
	jsonResp := &struct {
		Disk *Disk `json:"Disk"`
	}{
		Disk: l.Create(),
	}
	apiErr := new(APIError)
	resp, err := l.api.client.NewSling().Get(fmt.Sprintf("disk/%s", id)).Receive(jsonResp, apiErr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, apiErr
	}
	return jsonResp.Disk, nil
}

type diskRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	Plan        struct {
		ID int `json:"ID"`
	} `json:"Plan"`
	SizeMB     int    `json:"SizeMB,omitempty"`
	Connection string `json:"Connection,omitempty"`
	SourceDisk *struct {
		ID int `json:"ID,omitempty"`
	} `json:"SourceDisk,omitempty"`
	SourceArchive *struct {
		ID int `json:"ID,omitempty"`
	} `json:"SourceArchive,omitempty"`
}

type Disk struct {
	client      *Client
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Plan        struct {
		ID           int    `json:"ID"`
		Name         string `json:"Name"`
		StorageClass string `json:"StroageClass"`
	} `json:"Plan"`
	SizeMB     int    `json:"SizeMB"`
	Connection string `json:"Connection"`
	SourceDisk *struct {
		ID int `json:"ID"`
	} `json:"SourceDisk"`
	SourceArchive *struct {
		ID int `json:"ID"`
	} `json:"SourceArchive"`
	Availability    string `json:"Availability"`
	ConnectionOrder string `json:"ConnectionOrder"`
	ReinstallCount  int    `json:"ReinstallCount"`
	MigratedMB      int    `json:"MigratedMB"`
	WaitingJobCount int    `json:"WaitingJobCount"`
	JobStatus       struct {
	} `json:"JobStatus"`
	ServiceClass string `json:"ServiceClass"`
	BundleID     string `json:"BundleID"`
	CreatedAt    string `json:"CreatedAt"`
	Icon         string `json:"Icon"`
	Storage      struct {
		ID          string `json:"ID"`
		Class       string `json:"Class"`
		Name        string `json:"Name"`
		Description string `json:"Description"`
		Zone        struct {
			ID           int    `json:"ID"`
			DisplayOrder int    `json:"DisplayOrder"`
			Name         string `json:"Name"`
			Description  string `json:"Description"`
			IsDummy      bool   `json:"IsDummy"`
			VNCProxy     struct {
				HostName  string `json:"HostName"`
				IPAddress string `json:"IPAddress"`
			}
			FTPServer struct {
				HostName  string `json:"HostName"`
				IPAddress string `json:"IPAddress"`
			}
			Settings struct {
				Subnet struct {
					Plan struct {
						Member []int `json:"Member"`
						Staff  []int `json:"Staff"`
					} `json:"Plan"`
				} `json:"Subnet"`
			} `json:"Setteings"`
			Region struct {
				ID          int      `json:"ID"`
				Name        string   `json:"Name"`
				Description string   `json:"Description"`
				NameServers []string `json:"NameServers"`
			} `json:"Region"`
		} `json:"Zone"`
		DiskPlan struct {
			ID           int    `json:"ID"`
			StorageClass string `json:"StorageClass"`
			Name         string `json:"Name"`
			Capacity     []int  `json:"Capacity"`
		} `json:"DiskPlan"`
	} `json:"Storage"`
	Appliance *struct{} `json:"Appliance"`
	Server    *struct{} `json:"Server"`
	Tags      []string  `json:"Tags"`
}

func (l *Disk) Save() error {
	sling := l.client.NewSling()
	apiErr := new(APIError)
	var resp *http.Response
	var err error
	if l.ID == "" {
		postResp := &struct {
			IsOK    bool   `json:"is_ok"`
			Success string `json:"Success"`
			Disk    *Disk  `json:"Disk"`
		}{
			IsOK: false,
			Disk: l,
		}

		dr := &diskRequest{
			Name: l.Name,
		}
		dr.Plan.ID = l.Plan.ID
		dr.SizeMB = l.SizeMB

		postReq := &struct {
			Disk *diskRequest `json:"Disk"`
		}{
			Disk: dr,
		}

		resp, err = sling.Post("disk").BodyJSON(postReq).Receive(postResp, apiErr)
	} else {
		putResp := &struct {
			IsOK    bool   `json:"is_ok"`
			Success string `json:"Success"`
		}{
			IsOK: false,
		}
		resp, err = sling.Put("disk").BodyJSON(l).Receive(putResp, apiErr)
	}

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		apiErr.HTTPCode = resp.StatusCode
		return apiErr
	}
	return nil
}

func (l *Disk) Destroy() error {
	sling := l.client.NewSling()
	if l.ID == "" {
		return fmt.Errorf("is not saved yet")
	}
	apiErr := new(APIError)

	resp, err := sling.Delete(fmt.Sprintf("disk/%s", l.ID)).Receive(nil, apiErr)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return apiErr
	}
	return nil
}
