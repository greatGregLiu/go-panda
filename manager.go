package panda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/ernesto-jimenez/go-querystring/query"
)

type Status string

const (
	Success    Status = "success"
	Fail              = "fail"
	Processing        = "processing"
)

type Manager struct {
	cl Client
}

func NewManager(cl Client) *Manager {
	return &Manager{
		cl: cl,
	}
}

func (m *Manager) manageGet(url string, v interface{}) (err error) {
	b, err := m.cl.Get(url, nil)
	if err == nil {
		err = json.Unmarshal(b, v)
	}
	return
}

func (m *Manager) managePost(url string, r io.Reader, p, v interface{}) (err error) {
	params, err := query.Values(p)
	if err != nil {
		return
	}
	b, err := m.cl.Post(url, "", params, r)
	if err == nil {
		err = json.Unmarshal(b, v)
	}
	return
}

func (m *Manager) CloudId(id string) (c Cloud, err error) {
	if err = m.manageGet(fmt.Sprintf(cloudsIdPath, id), &c); err == nil {
		c.cl = m.cl
	}
	return
}

func (m *Manager) Clouds() (cs []Cloud, err error) {
	if err = m.manageGet(cloudsPath, &cs); err != nil {
		for i := range cs {
			cs[i].cl = m.cl
		}
	}
	return
}

func (m *Manager) NewEncoding(er *EncodingRequest) (e Encoding, err error) {
	if err = m.managePost(encodingsPath, nil, er, &e); err == nil {
		e.cl = m.cl
	}
	return
}

func (m *Manager) EncodingId(id string) (e Encoding, err error) {
	if err = m.manageGet(fmt.Sprintf(encodingsIdPath, id), &e); err == nil {
		e.cl = m.cl
	}
	return
}

func (m *Manager) NewProfile(pr *ProfileRequest) (p Profile, err error) {
	if err = m.managePost(profilesPath, nil, pr, &p); err == nil {
		p.cl = m.cl
	}
	return
}

func (m *Manager) ProfileId(id string) (p Profile, err error) {
	b, err := m.cl.Get(fmt.Sprintf(profilesIdPath, id), nil)
	if err == nil {
		err = json.Unmarshal(b, &p)
	}
	p.cl = m.cl
	return
}

func (m *Manager) Profiles() (ps []Profile, err error) {
	if err = m.manageGet(profilesPath, &ps); err == nil {
		for i := range ps {
			ps[i].cl = m.cl
		}
	}
	return
}

func (m *Manager) NewVideo(vr *VideoRequestUrl) (v Video, err error) {
	if err = m.managePost(videosPath, nil, vr, &v); err == nil {
		v.cl = m.cl
	}
	return
}

func (m *Manager) NewVideoReader(r io.Reader, name string, vr *VideoRequest) (v Video, err error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	if err = w.SetBoundary("--panda--"); err != nil {
		return
	}
	p, err := w.CreateFormFile("file", name)
	if err != nil {
		return
	}
	if _, err = io.Copy(p, r); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}
	var params url.Values
	if vr != nil {
		if params, err = query.Values(vr); err != nil {
			return
		}
	}
	b, err := m.cl.Post(videosPath, w.FormDataContentType(), params, buf)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &v); err == nil {
		v.cl = m.cl
	}
	return
}

func (m *Manager) Videos() (vs []Video, err error) {
	if err = m.manageGet(videosPath, &vs); err == nil {
		for i := range vs {
			vs[i].cl = m.cl
		}
	}
	return
}

func (m *Manager) VideoId(id string) (v Video, err error) {
	if err = m.manageGet(fmt.Sprintf(videosIdPath, id), &v); err == nil {
		v.cl = m.cl
	}
	return
}

func (m *Manager) Notifications() (n Notifications, err error) {
	if err = m.manageGet(notificationsPath, &n); err == nil {
		n.cl = m.cl
	}
	return
}

func (m *Manager) NotificationsId(id string) (n Notifications, err error) {
	if err = m.manageGet(fmt.Sprintf(notificationsIdPath, id), &n); err == nil {
		n.cl = m.cl
	}
	return
}
