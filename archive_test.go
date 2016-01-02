package saklient

import (
	"os"
	"testing"
)

func TestArchiveService(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	if api.Archive == nil {
		t.Fatal("api.Archive is nil")
	}
}
func TestArchiveService_CRUD(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	api.client.DebugDumper = os.Stderr
	archives, err := api.Archive.WithNameLike("CentOS 64bit").WithSizeGib(20).WithSharedScope().Limit(1).Find()
	if err != nil {
		t.Fatal(err)
	}
	if len(archives) != 1 {
		t.Fatal("Failed to find the archive 'CentOS 64bit' with 20GB")
	}
	archive := archives[0]
	if archive.ID == "" {
		t.Fatal("ID is empty")
	}

	newArchive := api.Archive.Create()
	newArchive.Plan.ID = 2
	newArchive.Name = "test"
	newArchive.Source = archive
	err = newArchive.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = newArchive.SleepWhileCopying()
	if err != nil {
		t.Fatal(err)
	}
	if newArchive.Availability != "available" {
		t.Fatalf("Unexpected availability: %s", newArchive.Availability)
	}
	err = newArchive.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}
