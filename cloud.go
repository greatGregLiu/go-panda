package panda

import (
	"encoding/json"
	"fmt"

	"github.com/ernesto-jimenez/go-querystring/query"
)

type Cloud struct {
	Id              string `json:"id" url:"id"`
	Name            string `json:"name" url:"name"`
	S3VideosBucket  string `json:"s3_videos_bucket" url:"s3_videos_bucket"`
	S3PrivateAccess bool   `json:"s3_private_access" url:"s3_private_access"`
	CreatedAt       string `json:"created_at" url:"created_at"`
	UpdatedAt       string `json:"updated_at" url:"updated_at"`
	Url             string `json:"url" url:"url"`
	cl              Client
}

func (c *Cloud) Update(accessKey, secretKey string) (err error) {
	if err = clientError(c.cl); err != nil {
		return
	}
	params, err := query.Values(c)
	if err != nil {
		return
	}
	params.Set("aws_access_key", accessKey)
	params.Set("aws_secret_key", secretKey)
	b, err := c.cl.Put(fmt.Sprintf(cloudsIdPath, c.Id), "", params, nil)
	if err == nil {
		err = json.Unmarshal(b, c)
	}
	return
}

func (c *Cloud) SetClient(cl Client) {
	c.cl = cl
}
