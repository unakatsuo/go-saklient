package saklient

import (
	"fmt"
	"testing"
)

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
	swytch, err := router.GetSwytch()
	if err != nil {
		t.Fatal(err)
	}
	if len(swytch.IPv4Nets) < 1 {
		t.Fatal("No IPv4Nets")
	}

	if len(swytch.IPv6Nets) > 0 {
		t.Log("Deregistering IPv6 network from the switch")
		err := swytch.RemoveIPv6Net()
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err = swytch.AddIPv6Net()
	if err != nil {
		t.Fatal(err)
	}
	if len(swytch.IPv6Nets) != 1 {
		t.Log(fmt.Sprintf("IPv6Nets len: %d, %v", len(swytch.IPv6Nets), swytch))
		t.Error("No IPv6Nets")
	}

	router.Reload()
	t.Log("Removing IPv6Net from the switch")
	err = swytch.RemoveIPv6Net()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Destroying the switch + router")
	err = router.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}
