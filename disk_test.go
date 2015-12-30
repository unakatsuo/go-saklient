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

func createTestDisk(name string) (*Disk, error) {
	api := API.Authorize(accessToken, accessSecret)
	disk := api.Disk.Create()
	disk.Name = name
	disk.Plan.ID = 2
	// Disk.SizeMB must be 20480,40960,61440,81920,102400,256000
	disk.SizeMB = 40960
	err := disk.Save()
	if err != nil {
		return nil, err
	}
	return disk, nil
}

func TestDiskServiceFilter(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	d1, err := createTestDisk("d1")
	d2, err := createTestDisk("d2")
	defer func() {
		d1.Destroy()
		d2.Destroy()
		t.Log("Destroied d1 and d2 disks")
	}()
	t.Log(fmt.Sprintf("d1.ID=%s, d2.ID=%s", d1.ID, d2.ID))

	var disks []*Disk

	disks, err = api.Disk.Limit(1).Offset(0).Find()
	if err != nil {
		t.Error(err)
	}
	if len(disks) != 1 {
		t.Error("Failed to filter the disks. Limit(1).Offset(0)")
	}
	t.Log(len(disks))
	disks, err = api.Disk.Reset().WithNameLike("d1").Find()
	if len(disks) != 1 {
		t.Error("Failed to filter the disks. WithNameLike('d1')")
	} else if disks[0].ID != d1.ID {
		t.Error("Wrong disk resource. disks[0].ID != d1.ID")
	}
	t.Log(len(disks))
}
