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

	copyArchive := api.Archive.Create()
	copyArchive.Plan.ID = 2
	copyArchive.Name = "test"
	copyArchive.Source = archive
	err = copyArchive.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = copyArchive.SleepWhileCopying()
	if err != nil {
		t.Fatal(err)
	}
	if copyArchive.Availability != "available" {
		t.Fatalf("Unexpected availability: %s", copyArchive.Availability)
	}
	err = copyArchive.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}
