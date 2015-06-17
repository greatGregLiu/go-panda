package panda

import (
	"net/url"
	"testing"
	"time"
)

// Taken from "Worked example" - https://www.pandastream.com/docs/#api-authentication
func TestBuildSignature(t *testing.T) {
	expected := "kVnZs/NX13ldKPdhFYoVnoclr8075DwiZF0TGgIbMsc="
	cl := Client{
		Host: "api.pandastream.com:8080",
		Options: &ClientOptions{
			CloudID:   "123456789",
			AccessKey: "abcdefgh",
			SecretKey: "ijklmnop",
		},
	}
	d := time.Date(2011, 3, 1, 15, 39, 10, 260762000, time.UTC)
	v := url.Values{}
	cl.addAuthParams(v, d)
	sign, err := cl.buildSignature(v, "GET", videosPath)
	if err != nil {
		t.Fatalf("want err=nil; got %v", err)
	}
	if sign != expected {
		t.Errorf("want signature=%s; got %s", expected, sign)
	}
}

func TestBuildUrl(t *testing.T) {
	cases := []struct {
		cl  Client
		v   url.Values
		url string
		exp string
	}{
		{
			Client{
				Host:    "localhost:9999",
				Options: &ClientOptions{},
			},
			url.Values{},
			"/videos.json",
			"http://localhost:9999/v2/videos.json",
		},
		{
			Client{
				Host:    "localhost:80",
				Options: &ClientOptions{},
			},
			url.Values{"cloud_id": []string{"12345"}},
			"/videos.json",
			"http://localhost:80/v2/videos.json?cloud_id=12345",
		},
		{
			Client{
				Host:    "localhost",
				Options: &ClientOptions{},
			},
			url.Values{"cloud_id": []string{"12345"}},
			"/videos.json",
			"http://localhost/v2/videos.json?cloud_id=12345",
		},
		{
			Client{
				Host:    "localhost:443",
				Options: &ClientOptions{},
			},
			url.Values{},
			"/videos.json",
			"https://localhost:443/v2/videos.json",
		},
	}
	for i, cas := range cases {
		if res := cas.cl.buildURL(cas.v, cas.url).String(); res != cas.exp {
			t.Errorf("expected uri=%s; got %s (i=%d)", cas.exp, res, i)
		}
	}
}
