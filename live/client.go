// package vod provides a client for Live transcoding service
package live

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pandastream/go-panda/client"
)

type Client struct {
	PandaClient *client.Client
}

func NewClient(token string, httpClient *http.Client) *Client {
	return &Client{
		PandaClient: &client.Client{
			Host: client.HostGCE,
			Options: &client.ClientOptions{
				Token:     token,
				Namespace: "live",
			},
			HTTPClient: httpClient,
		},
	}
}

func (c *Client) get(path string, v interface{}) error {
	b, err := c.PandaClient.Get(path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// ProfilesIDs lists all existing profiles.
func (c *Client) ProfilesIDs() (ids []string, err error) {
	if err := c.get("/v2/profiles.json", &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// Profile retrieves Profile object with the given id.
func (c *Client) Profile(id string) (*Profile, error) {
	profile := Profile{}
	if err := c.get(fmt.Sprintf("/v2/profiles/%s.json", id), &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// ProfileCreate creates new profile based on the given object.
func (c *Client) ProfileCreate(p *Profile) (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	b, err = c.PandaClient.Post("/v2/profiles.json", "application/json", nil, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	return resp.ProfileID, nil
}

// ProfileDelete deletes profile with the given id.
func (c *Client) ProfileDelete(id string) error {
	_, err := c.PandaClient.Delete(fmt.Sprintf("/v2/profiles/%s.json", id))
	return err
}

// StreamsIDs lists all existing streams.
func (c *Client) StreamsIDs() (ids []string, err error) {
	if err := c.get("/v2/streams.json", &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// Stream retrieves a stream object with the given id.
func (c *Client) Stream(id string) (*Stream, error) {
	stream := Stream{}
	if err := c.get(fmt.Sprintf("/v2/streams/%s.json", id), &stream); err != nil {
		return nil, err
	}
	return &stream, nil
}

// StreamCreate creates a stream based on the given object.
func (c *Client) StreamCreate(s *Stream) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	b, err = c.PandaClient.Post("/v2/streams.json", "application/json", nil, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	return resp.StreamID, nil
}

// StreamCreateProfile creates a stream and a profile based on the given profile object.
// If stream name argument is an empty string then the stream's name will be set to it's ID.
func (c *Client) StreamCreateProfile(streamName string, p *Profile) (streamID, profileID string, err error) {
	b, err := json.Marshal(reqStreamProfile{Profile: p, StreamName: streamName})
	if err != nil {
		return "", "", err
	}
	b, err = c.PandaClient.Post("/v2/streams/profile.json", "application/json", nil, bytes.NewReader(b))
	if err != nil {
		return "", "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", "", err
	}
	return resp.StreamID, resp.ProfileID, nil
}

// StreamDuration updates stream's duration.
func (c *Client) StreamDuration(id string, dur time.Duration) (streamID string, err error) {
	v := url.Values{}
	v.Add("duration", dur.String())
	b, err := c.PandaClient.Put(fmt.Sprintf("/v2/streams/%s/duration.json", id), "application/json", v, nil)
	if err != nil {
		return "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	return resp.StreamID, nil
}

// StreamDelete deletes a stream with the given id.
func (c *Client) StreamDelete(id string) error {
	_, err := c.PandaClient.Delete(fmt.Sprintf("/v2/streams/%s.json", id))
	return err
}
