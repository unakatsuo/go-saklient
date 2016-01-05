package saklient

import "testing"

func TestRouterService(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	if api.Router == nil {
		t.Fatal("api.Router is nil")
	}
}

func TestRouterService_CRUD(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)

	t.Log("creating a router and a switch")
	router := api.Router.Create()
	router.Name = "test"
	router.BandWidthMbps = 100
	router.NetworkMaskLen = 28
	err := router.Save()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("waiting for the router to become ready")
	err = router.SleepWhileCreating()
	if err != nil {
		t.Fatal(err)
	}
}
