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

	t.Log("creating a switch")
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

	t.Log("creating a server")
	server := api.Server.Create()
	server.Name = "test"
	server.ServerPlan.ID = 1001
	err = server.Save()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("adding an interface to the server")
	iface, err := server.AddIface()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("attaching the interface to the switch")
	err = iface.ConnectToSwytch(swytch)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("detaching the interface from the switch")
	err = iface.DisconnectFromSwytch()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("destroying the server")
	err = server.Destroy()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("deleting the switch")
	err = swytch.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}
