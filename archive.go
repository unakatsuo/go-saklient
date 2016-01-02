package saklient

import (
	"fmt"
	"time"
)

type ArchiveService struct {
	api      *APIService
	offset   int
	limit    int
	tags     []string
	filter   map[string]interface{}
	sortKeys []string
}

func (l *ArchiveService) Offset(offset int) *ArchiveService {
	l.offset = offset
	return l
}

func (l *ArchiveService) Limit(limit int) *ArchiveService {
	l.limit = limit
	return l
}

func (l *ArchiveService) SortByName(reverse bool) *ArchiveService {
	key := "Archive.Name"
	if reverse {
		key = "-" + key
	}
	l.sortKeys = append(l.sortKeys, key)
	return l
}

func (l *ArchiveService) SortBySize(reverse bool) *ArchiveService {
	key := "Archive.SizeMB"
	if reverse {
		key = "-" + key
	}
	l.sortKeys = append(l.sortKeys, key)
	return l
}

func (l *ArchiveService) FilterBy(key string, value interface{}, multiple bool) *ArchiveService {
	// TODO: multipe case
	l.filter[key] = value
	return l
}

func (l *ArchiveService) WithNameLike(name string) *ArchiveService {
	return l.FilterBy("Name", name, false)
}

func (l *ArchiveService) WithSharedScope() *ArchiveService {
	return l.FilterBy("Scope", "shared", false)
}

func (l *ArchiveService) WithTag(tag string) *ArchiveService {
	l.tags = append(l.tags, tag)
	return l
}

func (l *ArchiveService) WithTags(tags []string) *ArchiveService {
	for _, t := range tags {
		l.tags = append(l.tags, t)
	}
	return l
}

func (l *ArchiveService) WithSizeGib(size int) *ArchiveService {
	return l.FilterBy("SizeMB", size*1024, false)
}

func (l *ArchiveService) Reset() *ArchiveService {
	l.limit = 0
	l.offset = 0
	l.tags = []string{}
	l.sortKeys = []string{}
	l.filter = map[string]interface{}{}
	return l
}

func (l *ArchiveService) Find() ([]*Archive, error) {
	jsonResp := &struct {
		Total    int        `json:"Total`
		From     int        `json:"From"`
		Count    int        `json:"Count"`
		Archives []*Archive `json:"Archives"`
	}{}
	getReq := &struct {
		From   int                    `json:"From,omitempty"`
		Count  int                    `json:"Count,omitempty"`
		Sort   []string               `json:"Sort,omitempty"`
		Filter map[string]interface{} `json:"Filter,omitempty"`
	}{
		From:   l.offset,
		Count:  l.limit,
		Filter: l.filter,
		Sort:   l.sortKeys,
	}
	err := l.api.client.Request("GET", "archive", getReq, jsonResp)
	if err != nil {
		return nil, err
	}
	print(jsonResp.Total)
	for _, d := range jsonResp.Archives {
		d.client = l.api.client
	}
	return jsonResp.Archives, nil
}

func (l *ArchiveService) Create() *Archive {
	return &Archive{
		client: l.api.client,
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

type subResourceID struct {
	ID string `json:"ID"`
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

type SourceResource interface {
	SourceID() string
}

type Archive struct {
	client       *Client `json:"-"`
	ID           string  `json:"ID"`
	DisplayOrder string  `json:"DisplayOrder"`
	Name         string  `json:"Name"`
	Description  string  `json:"Description"`
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
			switch t := l.Source.(type) {
			case *Archive:
				dr.SourceArchive = &subResourceID{ID: l.Source.SourceID()}
			//case *Disk:
			//dr.SourceDisk = &subResourceID{ID: l.Source.SourceID()}
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

		err = l.client.Request("POST", "archive", postReq, postResp)
	} else {
		putResp := &struct {
			IsOK    bool   `json:"is_ok"`
			Success string `json:"Success"`
		}{
			IsOK: false,
		}
		err = l.client.Request("PUT", "archive", l, putResp)
	}

	return err
}

func (l *Archive) Destroy() error {
	if l.ID == "" {
		return fmt.Errorf("is not saved yet")
	}
	return l.client.Request("DELETE", fmt.Sprintf("archive/%s", l.ID), nil, nil)
}

func (l *Archive) Reload() error {
	if l.ID == "" {
		return fmt.Errorf("This is not saved yet")
	}
	jsonResp := &struct {
		Archive *Archive `json:"Archive"`
	}{
		Archive: &Archive{client: l.client},
	}
	err := l.client.Request("GET", fmt.Sprintf("archive/%s", l.ID), nil, jsonResp)
	if err != nil {
		return err
	}
	*l = *jsonResp.Archive
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
