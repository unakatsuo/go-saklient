package saklient

import (
	"fmt"
	"testing"
)

func TestDiskService(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	if api.Disk == nil {
		t.Fatal("api.Disk is nil")
	}
}

func TestDiskService_CRUD(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	disk := api.Disk.Create()
	disk.Name = "test"
	disk.Plan.ID = 2
	// Disk.SizeMB must be 20480,40960,61440,81920,102400,256000
	disk.SizeMB = 40960
	err := disk.Save()
	if err != nil {
		t.Fatal(err)
	}
	if disk.ID == "" {
		t.Fatal("ID is unset after create")
	}
	t.Log(fmt.Sprintf("New Disk: %s", disk.ID))
	disk2, err := api.Disk.GetByID(disk.ID)
	if err != nil {
		t.Fatal(err)
	}
	err = disk2.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}
