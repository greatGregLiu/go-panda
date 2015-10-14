package live

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pandastream/go-panda"
)

type Client struct {
	Client *panda.Client
}

func (cl *Client) get(path string, v interface{}) error {
	b, err := cl.Client.Get(path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (cl *Client) ProfilesIDs() (ids []string, err error) {
	if err := cl.get("/v2/profiles.json", &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (cl *Client) Profile(id string) (*Profile, error) {
	profile := Profile{}
	if err := cl.get(fmt.Sprintf("/v2/profiles/%s.json", id), &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (cl *Client) ProfileCreate(p *Profile) (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	b, err = cl.Client.Post("/v2/profiles.json", "application/json", nil, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	return resp.ProfileID, nil
}

func (cl *Client) ProfileDelete(id string) error {
	_, err := cl.Client.Delete(fmt.Sprintf("/v2/profiles/%s.json", id))
	return err
}

func (cl *Client) StreamsIDs() (ids []string, err error) {
	if err := cl.get("/v2/streams.json", &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (cl *Client) Stream(id string) (*Stream, error) {
	stream := Stream{}
	if err := cl.get(fmt.Sprintf("/v2/streams/%s.json", id), &stream); err != nil {
		return nil, err
	}
	return &stream, nil
}

func (cl *Client) StreamCreate(s *Stream) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	b, err = cl.Client.Post("/v2/streams.json", "application/json", nil, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	return resp.StreamID, nil
}

func (cl *Client) StreamCreateProfile(p *Profile) (streamID, profileID string, err error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", "", err
	}
	b, err = cl.Client.Post("/v2/streams/profile.json", "application/json", nil, bytes.NewReader(b))
	if err != nil {
		return "", "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", "", err
	}
	return resp.StreamID, resp.ProfileID, nil
}

func (cl *Client) StreamDuration(id string, dur time.Duration) (streamID string, err error) {
	v := url.Values{}
	v.Add("duration", dur.String())
	b, err := cl.Client.Put(fmt.Sprintf("/v2/streams/%s/duration.json", id), "application/json", v, nil)
	if err != nil {
		return "", err
	}
	resp := postResp{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	return resp.StreamID, nil
}

func (cl *Client) StreamDelete(id string) error {
	_, err := cl.Client.Delete(fmt.Sprintf("/v2/streams/%s.json", id))
	return err
}
