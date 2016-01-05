package saklient

import "fmt"

type DiskConfig struct {
	client        *Client `json:"-"`
	diskID        string  `json:"-"`
	HostName      string  `json:"HostName"`
	Password      string  `json:"Password"`
	UserIPAddress string  `json:"UserIPAddress"`
	SSHKey        struct {
		ID        string `json:"ID"`
		PublicKey string `json:"PublicKey"`
	} `json:"SSHKey"`
	UserSubnet struct {
		DefaultRoute   string `json:"DefaultRoute"`
		NetworkMaskLen string `json:"NetworkMaskLen"`
	} `json:"UserSubnet"`
}

func (d *DiskConfig) DiskID() string {
	return d.diskID
}

func (d *DiskConfig) Write() error {
	if d.diskID == "" {
		return fmt.Errorf("Disk.ID is unset")
	}
	return d.client.Request("PUT", fmt.Sprintf("disk/%s/config", d.diskID), d, nil)
}
