package vod

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pandastream/go-panda/client"
)

func mustWrite(w http.ResponseWriter, b []byte) {
	if _, err := w.Write(b); err != nil {
		panic(err)
	}
}

func newClient(addr string, t *testing.T) *Client {
	URL, err := url.Parse(addr)
	if err != nil {
		t.Fatal(err)
	}
	return NewClient(URL.Host, "", nil)
}

func TestDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/profiles/"):
		case strings.Contains(r.URL.Path, "/videos/"):
		case strings.Contains(r.URL.Path, "/encodings/"):
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			mustWrite(w, []byte("Invalid path"))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	m := newClient(ts.URL, t)
	cases := []struct {
		obj   interface{}
		valid bool
		err   error
	}{
		{
			&Profile{},
			true,
			nil,
		},
		{
			&Video{},
			true,
			nil,
		},
		{
			&Encoding{},
			true,
			&client.Error{Code: http.StatusBadRequest},
		},
		{
			&Notification{},
			false,
			nil,
		},
	}
	for i := range cases {
		if cases[i].valid {
			if err := m.Delete(cases[i].obj); !reflect.DeepEqual(err, cases[i].err) {
				t.Errorf("want err=%v; got %v (i=%d)", cases[i].err, err, i)
			}
			continue
		}
		done := make(chan struct{})
		go func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected Delete to panic (i=%d)", i)
				}
				done <- struct{}{}
			}()
			m.Delete(cases[i].obj)
		}()
		<-done
	}
}

func TestUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b []byte
		var err error
		switch {
		case strings.Contains(r.URL.Path, "/profiles/"):
			b, err = json.Marshal(&Profile{Name: "ProfileName"})
		case strings.Contains(r.URL.Path, "/notifications.json"):
			b, err = json.Marshal(&Notification{Delay: 5})
		default:
			mustWrite(w, []byte("Invalid path"))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			mustWrite(w, []byte(err.Error()))
		}
		mustWrite(w, b)
	}))
	defer ts.Close()
	m := newClient(ts.URL, t)
	zeroTime, err := time.Parse(timeFormat, "0001/01/01 00:00:00 +0000")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		obj      interface{}
		valid    bool
		expected interface{}
	}{
		{
			&Notification{},
			true,
			&Notification{Delay: 5},
		},
		{
			&Profile{},
			true,
			&Profile{
				Name:      "ProfileName",
				CreatedAt: Time(zeroTime),
				UpdatedAt: Time(zeroTime),
			},
		},
		{
			&Encoding{},
			false,
			nil,
		},
	}
	for i := range cases {
		if cases[i].valid {
			if err := m.Update(cases[i].obj); err != nil {
				t.Errorf("want err=nil; got %v (i=%d)", err, i)
			}
			if !reflect.DeepEqual(cases[i].obj, cases[i].expected) {
				t.Errorf("want %#v; got %#v", cases[i].expected, cases[i].obj)
			}
			continue
		}
		done := make(chan struct{})
		go func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected Update to panic (i=%d)", i)
				}
				done <- struct{}{}
			}()
			m.Update(cases[i].obj)
		}()
		<-done
	}
}
