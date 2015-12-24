package saklient

import (
	"os"
	"testing"
)

var accessToken string
var accessSecret string

func init() {
	accessToken = os.Getenv("SAKURA_ACCESS_TOKEN")
	if accessToken == "" {
		panic("SAKURA_ACCESS_TOKEN is empty")
	}
	accessSecret = os.Getenv("SAKURA_ACCESS_SECRET")
	if accessSecret == "" {
		panic("SAKURA_ACCESS_SECRET is empty")
	}
}

func TestNew(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	if api.Server == nil {
		t.Fatalf(".Server is nil")
	}
}

func TestServerCreate(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	server, err := api.Server.Create()
	if err != nil {
		t.Fatal(err)
	}
	server.Name = "svr1"
	server.Description = "test"
	_, _, err = server.Save()
	if err != nil {
		t.Fatal(err)
	}

}
