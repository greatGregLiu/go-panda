package panda

import (
	"net/url"
	"testing"
	"time"
)

// Taken from "Worked example" - https://www.pandastream.com/docs/api#api_authentication
func TestBuildSignature(t *testing.T) {
	expected := "kVnZs/NX13ldKPdhFYoVnoclr8075DwiZF0TGgIbMsc="
	cl := client{
		host:      "api.pandastream.com",
		cloudId:   "123456789",
		accessKey: "abcdefgh",
		secretKey: []byte("ijklmnop"),
	}
	d := time.Date(2011, 3, 1, 15, 39, 10, 260762000, time.UTC)
	v := url.Values{}
	cl.addAuthParams(v, d)
	sign, err := cl.buildSignature(v, "GET", videosPath, d)
	if err != nil {
		t.Fatalf("want err=nil; got %v", err)
	}
	if sign != expected {
		t.Errorf("want signature=%s; got %s", expected, sign)
	}
}

func TestBuildUrl(t *testing.T) {
	cases := []struct {
		cl  client
		v   url.Values
		url string
		exp string
	}{
		{
			client{
				host: "localhost",
				port: 9999,
				ver:  "v2",
			},
			url.Values{},
			"/videos.json",
			"http://localhost:9999/v2/videos.json",
		},
		{
			client{
				host: "localhost",
				port: 80,
				ver:  "v1",
			},
			url.Values{"cloud_id": []string{"12345"}},
			"/videos.json",
			"http://localhost:80/v1/videos.json?cloud_id=12345",
		},
		{
			client{
				host: "localhost",
				port: 443,
				ver:  "v2",
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
