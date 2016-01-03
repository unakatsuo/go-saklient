package saklient

import (
	"fmt"
	"time"
)

type ArchiveService struct {
	api *APIService
	basicQuery
}

func newArchiveService(api *APIService) *ArchiveService {
	return (&ArchiveService{
		api: api,
	}).Reset()
}
func (l *ArchiveService) Offset(offset int) *ArchiveService {
	l.basicQuery.Offset = offset
	return l
}

func (l *ArchiveService) Limit(limit int) *ArchiveService {
	l.basicQuery.Limit = limit
	return l
}

func (l *ArchiveService) SortByName(reverse bool) *ArchiveService {
	l.basicQuery.SortBy("Archive.Name", reverse)
	return l
}

func (l *ArchiveService) SortBySize(reverse bool) *ArchiveService {
	l.basicQuery.SortBy("Archive.SizeMB", reverse)
	return l
}

func (l *ArchiveService) FilterBy(key string, value interface{}, multiple bool) *ArchiveService {
	// TODO: multipe case
	l.basicQuery.Filter[key] = value
	return l
}

func (l *ArchiveService) WithNameLike(name string) *ArchiveService {
	return l.FilterBy("Name", name, false)
}

func (l *ArchiveService) WithSharedScope() *ArchiveService {
	return l.FilterBy("Scope", "shared", false)
}

func (l *ArchiveService) WithTag(tag string) *ArchiveService {
	l.basicQuery.WithTag(tag)
	return l
}

func (l *ArchiveService) WithTags(tags []string) *ArchiveService {
	l.basicQuery.WithTags(tags)
	return l
}

func (l *ArchiveService) WithSizeGib(size int) *ArchiveService {
	return l.FilterBy("SizeMB", size*1024, false)
}

func (l *ArchiveService) Reset() *ArchiveService {
	l.basicQuery.Reset()
	return l
}

func (l *ArchiveService) Find() ([]*Archive, error) {
	jsonResp := &struct {
		Total    int        `json:"Total`
		From     int        `json:"From"`
		Count    int        `json:"Count"`
		Archives []*Archive `json:"Archives"`
	}{}
	err := l.api.client.Request("GET", "archive", l.basicQuery, jsonResp)
	if err != nil {
		return nil, err
	}
	for _, d := range jsonResp.Archives {
		d.service = l
	}
	return jsonResp.Archives, nil
}

func (l *ArchiveService) Create() *Archive {
	return &Archive{
		service: l,
	}
}

func (l *ArchiveService) GetByID(id string) (*Archive, error) {
	jsonResp := &struct {
		Archive *Archive `json:"Archive"`
	}{
		Archive: l.Create(),
	}
	err := l.api.client.Request("GET", fmt.Sprintf("archive/%s", id), nil, jsonResp)
	if err != nil {
		return nil, err
	}
	return jsonResp.Archive, nil
}

type archiveRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	Plan        struct {
		ID int `json:"ID"`
	} `json:"Plan"`
	SizeMB        int            `json:"SizeMB,omitempty"`
	SourceDisk    *subResourceID `json:"SourceDisk,omitempty"`
	SourceArchive *subResourceID `json:"SourceArchive,omitempty"`
}

type Archive struct {
	service      *ArchiveService `json:"-"`
	ID           string          `json:"ID"`
	DisplayOrder string          `json:"DisplayOrder"`
	Name         string          `json:"Name"`
	Description  string          `json:"Description"`
	Plan         struct {
		ID           int    `json:"ID"`
		Name         string `json:"Name"`
		StorageClass string `json:"StroageClass"`
	} `json:"Plan"`
	SizeMB          int            `json:"SizeMB"`
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
	OriginalArchive struct {
		ID string `json:"ID"`
	} `json:"OriginalArchive"`
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
	Tags []string `json:"Tags"`
}

// SourceResource interface.
func (l *Archive) SourceID() string {
	return l.ID
}

func (l *Archive) Save() error {
	var err error
	if l.ID == "" {
		postResp := &struct {
			IsOK    bool     `json:"is_ok"`
			Success bool     `json:"Success"`
			Archive *Archive `json:"Archive"`
		}{
			IsOK:    false,
			Archive: l,
		}

		dr := &archiveRequest{
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
			Archive *archiveRequest `json:"Archive"`
		}{
			Archive: dr,
		}

		err = l.client().Request("POST", "archive", postReq, postResp)
	} else {
		putResp := &struct {
			IsOK    bool   `json:"is_ok"`
			Success string `json:"Success"`
		}{
			IsOK: false,
		}
		err = l.client().Request("PUT", fmt.Sprintf("archive/%s", l.ID), l, putResp)
	}

	return err
}

func (l *Archive) Destroy() error {
	if l.ID == "" {
		return fmt.Errorf("is not saved yet")
	}
	return l.client().Request("DELETE", fmt.Sprintf("archive/%s", l.ID), nil, nil)
}

func (l *Archive) Reload() error {
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

func (l *Archive) SleepWhileCopying() error {
	if l.ID == "" {
		return fmt.Errorf("This is not saved yet")
	}
	var err error
	for i := 0; i < 1000; i++ {
		if i > 0 {
			time.Sleep(1 * time.Second)
			err = l.Reload()
			if err != nil {
				continue
			}
		}
		if l.Availability == "available" {
			break
		}
	}
	return err
}

func (l *Archive) client() *Client {
	return l.service.api.client
}
