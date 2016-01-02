package saklient

import "testing"

func TestServerService(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	if api.Server == nil {
		t.Fatal("api.Server is nil")
	}
}

func TestServerService_CRUD(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	server := api.Server.Create()
	server.Name = "svr1"
	server.Description = "test"
	server.Plan.ID = 1
	err := server.Save()
	if err != nil {
		t.Fatal(err)
	}
}
