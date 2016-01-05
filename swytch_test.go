package saklient

import "testing"

func TestSwytchService(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	if api.Swytch == nil {
		t.Fatal("api.Disk is nil")
	}
}

func TestSwytchService_CRUD(t *testing.T) {
	api := API.Authorize(accessToken, accessSecret)
	swytch := api.Swytch.Create()
	swytch.Name = "test"
	swytch.Description = "test"
	err := swytch.Save()
	if err != nil {
		t.Fatal(err)
	}
	if swytch.ID == "" {
		t.Fatal("swytch.ID is empty")
	}

	t.Log("deleting a switch")
	err = swytch.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}
