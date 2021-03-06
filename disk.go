package saklient

import "fmt"

type DiskService struct {
	api *APIService
	basicQuery
}

func newDiskService(api *APIService) *DiskService {
	return (&DiskService{
		api: api,
	}).Reset()
}

func (l *DiskService) Offset(offset int) *DiskService {
	l.basicQuery.Offset = offset
	return l
}

func (l *DiskService) Limit(limit int) *DiskService {
	l.basicQuery.Limit = limit
	return l
}

func (l *DiskService) SortByName(reverse bool) *DiskService {
	l.basicQuery.SortBy("Disk.Name", reverse)
	return l
}

func (l *DiskService) SortBySize(reverse bool) *DiskService {
	l.basicQuery.SortBy("Disk.SizeMB", reverse)
	return l
}

func (l *DiskService) FilterBy(key string, value interface{}, multiple bool) *DiskService {
	// TODO: multipe case
	l.basicQuery.Filter[key] = value
	return l
}

func (l *DiskService) WithNameLike(name string) *DiskService {
	return l.FilterBy("Name", name, false)
}

func (l *DiskService) WithTag(tag string) *DiskService {
	l.basicQuery.WithTag(tag)
	return l
}

func (l *DiskService) WithTags(tags []string) *DiskService {
	l.basicQuery.WithTags(tags)
	return l
}

func (l *DiskService) WithSizeGib(size int) *DiskService {
	return l.FilterBy("SizeMB", size*1024, false)
}

func (l *DiskService) WithServerID(serverID string) *DiskService {
	return l.FilterBy("Server.ID", serverID, false)
}

func (l *DiskService) Reset() *DiskService {
	l.basicQuery.Reset()
	return l
}

func (l *DiskService) Find() ([]*Disk, error) {
	jsonResp := &struct {
		Total int     `json:"Total`
		From  int     `json:"From"`
		Count int     `json:"Count"`
		Disks []*Disk `json:"Disks"`
	}{}
	err := l.api.client.Request("GET", "disk", l.basicQuery, jsonResp)
	if err != nil {
		return nil, err
	}

	for _, d := range jsonResp.Disks {
		d.service = l
	}
	return jsonResp.Disks, nil
}

func (l *DiskService) Create() *Disk {
	return &Disk{
		service: l,
	}
}

func (l *DiskService) GetByID(id string) (*Disk, error) {
	jsonResp := &struct {
		Disk *Disk `json:"Disk"`
	}{
		Disk: l.Create(),
	}
	err := l.api.client.Request("GET", fmt.Sprintf("disk/%s", id), nil, jsonResp)
	if err != nil {
		return nil, err
	}
	return jsonResp.Disk, nil
}

type diskRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	Plan        struct {
		ID int `json:"ID"`
	} `json:"Plan"`
	SizeMB        int            `json:"SizeMB,omitempty"`
	Connection    string         `json:"Connection,omitempty"`
	SourceDisk    *subResourceID `json:"SourceDisk,omitempty"`
	SourceArchive *subResourceID `json:"SourceArchive,omitempty"`
}

type Disk struct {
	service     *DiskService `json:"-"`
	ID          string       `json:"ID"`
	Name        string       `json:"Name"`
	Description string       `json:"Description"`
	Plan        struct {
		ID           int    `json:"ID"`
		Name         string `json:"Name"`
		StorageClass string `json:"StroageClass"`
	} `json:"Plan"`
	SizeMB          int            `json:"SizeMB"`
	Connection      string         `json:"Connection"`
	Source          SourceResource `json:"-"`
	SourceDisk      *Disk          `json:"SourceDisk"`
	SourceArchive   *Archive       `json:"SourceArchive"`
	Availability    string         `json:"Availability"`
	ConnectionOrder string         `json:"ConnectionOrder"`
	ReinstallCount  int            `json:"ReinstallCount"`
	MigratedMB      int            `json:"MigratedMB"`
	WaitingJobCount int            `json:"WaitingJobCount"`
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

func (l *Disk) SourceID() string {
	return l.ID
}

func (l *Disk) Save() error {
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
		if l.Source != nil {
			subRes := &subResourceID{ID: l.Source.SourceID()}
			switch t := l.Source.(type) {
			case *Archive:
				dr.SourceArchive = subRes
			case *Disk:
				dr.SourceDisk = subRes
			default:
				return fmt.Errorf("Unsupported .Source type: %T", t)
			}
		} else {
			dr.SizeMB = l.SizeMB
		}
		postReq := &struct {
			Disk *diskRequest `json:"Disk"`
		}{
			Disk: dr,
		}

		err = l.client().Request("POST", "disk", postReq, postResp)
	} else {
		putResp := &struct {
			IsOK    bool   `json:"is_ok"`
			Success string `json:"Success"`
		}{
			IsOK: false,
		}
		err = l.client().Request("PUT", fmt.Sprintf("disk/%s", l.ID), l, putResp)
	}

	return err
}

func (l *Disk) Destroy() error {
	if l.ID == "" {
		return fmt.Errorf("is not saved yet")
	}
	return l.client().Request("DELETE", fmt.Sprintf("disk/%s", l.ID), nil, nil)
}

func (l *Disk) Reload() error {
	if l.ID == "" {
		return fmt.Errorf("This is not saved yet")
	}
	n, err := l.service.GetByID(l.ID)
	if err != nil {
		return err
	}
	*l = *n
	return nil
}

func (l *Disk) SleepWhileCopying() error {
	if l.ID == "" {
		return fmt.Errorf("This is not saved yet")
	}
	return sleepUntil(func() bool {
		err := l.Reload()
		if err != nil {
			return false
		}
		return (l.Availability == "available")
	})
}

func (l *Disk) client() *Client {
	return l.service.api.client
}

func (l *Disk) ConnectTo(server *Server) error {
	if l.ID == "" {
		return fmt.Errorf("Disk is not saved")
	}
	if server.ID == "" {
		return fmt.Errorf("Server is invalid")
	}
	return l.client().Request("PUT", fmt.Sprintf("disk/%s/to/server/%s", l.ID, server.ID), nil, nil)
}

func (l *Disk) Disconnect() error {
	if l.ID == "" {
		return fmt.Errorf("Disk is not saved")
	}
	return l.client().Request("DELETE", fmt.Sprintf("disk/%s/to/server", l.ID), nil, nil)
}

func (l *Disk) CreateConfig() *DiskConfig {
	return &DiskConfig{
		client: l.client(),
		diskID: l.ID,
	}
}
