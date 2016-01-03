package saklient

import (
	"fmt"
	"testing"
)

func TestServerService(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	if api.Server == nil {
		t.Fatal("api.Server is nil")
	}
}

func TestServerService_CRUD(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)

	archives, err := api.Archive.WithNameLike("CentOS 64bit").WithSizeGib(20).WithSharedScope().Limit(1).Find()
	if err != nil {
		t.Fatal(err)
	}
	archive := archives[0]

	disk := api.Disk.Create()
	disk.Name = "test"
	disk.Plan.ID = 4
	disk.Source = archive
	err = disk.Save()
	if err != nil {
		t.Fatal(err)
	}

	server := api.Server.Create()
	server.Name = "svr1"
	server.Description = "test"
	server.ServerPlan.ID = 1001
	err = server.Save()
	if err != nil {
		t.Fatal(err)
	}

	err = disk.SleepWhileCopying()
	if err != nil {
		t.Fatal(err)
	}
	disk.Reload()
	if !(disk.SizeMB == archive.SizeMB && disk.SourceArchive.ID == archive.ID) {
		t.Fatal(fmt.Sprintf("Invalid Disk %s info", disk.ID))
	}
}
