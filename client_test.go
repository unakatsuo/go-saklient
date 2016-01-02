package saklient

import "os"

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
