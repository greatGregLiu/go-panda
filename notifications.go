package panda

import (
	"encoding/json"

	"github.com/ernesto-jimenez/go-querystring/query"
)

type Notifications struct {
	Url    string `json:"url" url:"url"`
	Events Events `json:"events" url:"events"`
	Delay  int    `json:"delay" url:"delay"`
	cl     Client
}

type Events struct {
	VideoCreated      bool `json:"video_created" url:"video_created"`
	VideoEncoded      bool `json:"video_encoded" url:"video_encoded"`
	EncodingProgress  bool `json:"encoding_progress" url:"encoding_progress"`
	EncodingCompleted bool `json:"encoding_completed" url:"encoding_completed"`
}

func (n *Notifications) Update() (err error) {
	if err = clientError(n.cl); err != nil {
		return
	}
	params, err := query.Values(n)
	if err != nil {
		return
	}
	b, err := n.cl.Put(notificationsPath, "", params, nil)
	if err == nil {
		err = json.Unmarshal(b, n)
	}
	return
}

func (n *Notifications) SetClient(cl Client) {
	n.cl = cl
}
