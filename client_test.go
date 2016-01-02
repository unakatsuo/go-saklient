package saklient

import (
	"net/url"
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

func TestClientDebugDumper(t *testing.T) {
	c := newClient(accessToken, accessSecret)
	c.DebugDumper = os.Stderr
	c.BaseURL, _ = url.Parse("https://secure.sakura.ad.jp/cloud/zone/tk1v/api/cloud/1.1/")
	testReq := &struct {
		Limit int `json:"Limit"`
	}{
		Limit: 1,
	}
	c.Request("GET", "disk", testReq, nil)
}
