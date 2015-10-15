// package vod provides a client for PandaStream file transcoding service
package vod

//go:generate structgen -dir=json -tags=url -o=models.go -pkg=panda -types=created_at:Time,updated_at:Time,height:int,width:int,duration:int,file_size:int64,audio_bitrate:int,audio_channels:int,video_bitrate:int,audio_sample_rate:int,watermark_bottom:int,watermark_height:int,watermark_left:int,watermark_right:int,watermark_top:int,watermark_width:int,keyframe_interval:int,buffer_size:int,max_rate:int,frame_count:int,h264_crf:int,status:Status,page:int,per_page:int,aspect_mode:AspectMode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pandastream/go-panda/client"

	"github.com/ernesto-jimenez/go-querystring/query"
)

// Status describes current state of the video or encoding
type Status string

const (
	StatusSuccess    = Status("success")
	StatusFail       = Status("fail")
	StatusProcessing = Status("processing")
)

type AspectMode string

const (
	ModeLetterBox = AspectMode("letterbox")
	ModePreserve  = AspectMode("preserve")
	ModeConstrain = AspectMode("constrain")
	ModePad       = AspectMode("pad")
	ModeCrop      = AspectMode("crop")
)

// MetaData holds all the video's data which ffprobe was able to get
type MetaData map[string]interface{}

const timeFormat = "2006/01/02 15:04:05 -0700"

// Time holds time.Time and is capable of marshalling and unmarshalling Panda's timestamp
// in the correct way
type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(`"` + timeFormat + `"`)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	time, err := time.Parse(`"`+timeFormat+`"`, string(data))
	if err != nil {
		return err
	}
	*t = Time(time)
	return nil
}

// Client provides methods of getting certain data from the Panda Cloud using panda Client
type Client struct {
	PandaClient *client.Client
}

func NewClient(host, token string, httpClient *http.Client) *Client {
	return &Client{
		PandaClient: &client.Client{
			Host: host,
			Options: &client.ClientOptions{
				Token: token,
			},
			HTTPClient: httpClient,
		},
	}
}

func (c *Client) get(url string, v interface{}, params url.Values) (err error) {
	b, err := c.PandaClient.Get(url, params)
	if err == nil {
		err = json.Unmarshal(b, v)
	}
	return
}

func (c *Client) post(url string, r io.Reader, p, v interface{}) error {
	params, err := query.Values(p)
	if err != nil {
		return err
	}
	b, err := c.PandaClient.Post(url, "", params, r)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// Cloud gets cloud by the given cloud ID
func (c *Client) Cloud(id string) (*Cloud, error) {
	cl := new(Cloud)
	if err := c.get(fmt.Sprintf(cloudsIdPath, id), cl, nil); err != nil {
		return nil, err
	}
	return cl, nil
}

// Clouds gets all clouds on the given account
func (c *Client) Clouds() (cs []Cloud, err error) {
	err = c.get(cloudsPath, &cs, nil)
	return
}

// NewEncoding creates a new encoding for the existing video
func (c *Client) NewEncoding(er *NewEncodingRequest) (*Encoding, error) {
	e := new(Encoding)
	if err := c.post(encodingsPath, nil, er, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// Encoding gets encoding object with the given id
func (c *Client) Encoding(id string) (*Encoding, error) {
	e := new(Encoding)
	if err := c.get(fmt.Sprintf(encodingsIdPath, id), e, nil); err != nil {
		return nil, err
	}
	return e, nil
}

// Encodings gets all encodings on the current cloud. Encodings can be filtered by options
// set in the EncodingRequest struct. If EncodingRequest is nil defaults are going to be used
func (c *Client) Encodings(er *EncodingRequest) (es []Encoding, err error) {
	var params url.Values
	if er != nil {
		params, err = query.Values(er)
		if err != nil {
			return
		}
	}
	err = c.get(encodingsPath, &es, params)
	return
}

// Cancel encoding with the given id
func (c *Client) Cancel(id string) error {
	_, err := c.PandaClient.Post(fmt.Sprintf(encodingsIdCancelPath, id), "", nil, nil)
	return err
}

// Retry encoding with the given id
func (c *Client) Retry(id string) error {
	_, err := c.PandaClient.Post(fmt.Sprintf(encodingsIdRetryPath, id), "", nil, nil)
	return err
}

// Delete accepts *Profile, *Video and *Encoding types and deletes those objects
// from Panda's database by their ID
func (c *Client) Delete(v interface{}) error {
	var path string
	switch t := v.(type) {
	case *Profile:
		path = fmt.Sprintf(profilesIdPath, t.ID)
	case *Video:
		path = fmt.Sprintf(videosIdDeletePath, t.ID)
	case *Encoding:
		path = fmt.Sprintf(encodingsIdPath, t.ID)
	default:
		panic("Invalid type")
	}
	_, err := c.PandaClient.Delete(path)
	return err
}

// NewProfile creates new profile based on profile request object
func (c *Client) NewProfile(pr *NewProfileRequest) (*Profile, error) {
	p := new(Profile)
	if err := c.post(profilesPath, nil, pr, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Profile gets profile with the given ID
func (c *Client) Profile(id string) (*Profile, error) {
	p := new(Profile)
	if err := c.get(fmt.Sprintf(profilesIdPath, id), p, nil); err != nil {
		return nil, err
	}
	return p, nil
}

// Profiles gets all profiles from the current cloud
func (c *Client) Profiles(pr *ProfileRequest) (ps []Profile, err error) {
	var params url.Values
	if pr != nil {
		params, err = query.Values(pr)
		if err != nil {
			return
		}
	}
	err = c.get(profilesPath, &ps, params)
	return
}

// Update accepts *Profile and *Notification types and updates records based on the given objects.
// Warning: the given parameter might change if any of the parameters are invalid
func (c *Client) Update(v interface{}) error {
	var path string
	switch t := v.(type) {
	case *Profile:
		path = fmt.Sprintf(profilesIdPath, t.ID)
	case *Notification:
		path = fmt.Sprintf(notificationsPath)
	default:
		panic("Invalid type")
	}
	params, err := query.Values(v)
	if err != nil {
		return err
	}
	b, err := c.PandaClient.Put(path, "", params, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// NewVideo gets file name path and potential additional options in a form of
// VideoRequest type and based on this creates a new video in Panda
func (c *Client) NewVideo(file string, vr *NewVideoRequest) (*Video, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return c.NewVideoReader(f, file, vr)
}

// NewVideoURL creates a new video in Panda based on the given source url and potential
// additional options in form of VideoRequest struct
func (c *Client) NewVideoURL(URL string, vr *NewVideoRequest) (*Video, error) {
	params, err := query.Values(vr)
	if err != nil {
		return nil, err
	}
	params.Set("source_url", URL)
	b, err := c.PandaClient.Post(videosPath, "", params, nil)
	if err != nil {
		return nil, err
	}
	v := new(Video)
	if err = json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}

// NewVideoReader creates a new video in Panda with the given name and reading from the given reader
func (c *Client) NewVideoReader(r io.Reader, name string, vr *NewVideoRequest) (*Video, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	if err := w.SetBoundary("--panda--"); err != nil {
		return nil, err
	}
	p, err := w.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(p, r); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	var params url.Values
	if vr != nil {
		if params, err = query.Values(vr); err != nil {
			return nil, err
		}
	}
	b, err := c.PandaClient.Post(videosPath, w.FormDataContentType(), params, buf)
	if err != nil {
		return nil, err
	}
	v := new(Video)
	if err = json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}

// Video gets video with the given id
func (c *Client) Video(id string) (*Video, error) {
	v := new(Video)
	if err := c.get(fmt.Sprintf(videosIdPath, id), &v, nil); err != nil {
		return nil, err
	}
	return v, nil
}

// Videos gets all the videos from the current cloud. Videos can be filtered by options
// set in the VideoRequest struct. If VideoRequest is nil then defaults are going to be used
func (c *Client) Videos(vr *VideoRequest) (v []Video, err error) {
	var params url.Values
	if vr != nil {
		params, err = query.Values(vr)
		if err != nil {
			return nil, err
		}
	}
	err = c.get(videosPath, &v, params)
	return
}

// VideoEncoding gets all encodings related to the video with the given id. Encodings can be
// filtered by options set in the EncodingRequest struct.
// If EncodingRequest is nil defaults are going to be used
func (c *Client) VideoEncodings(id string, er *EncodingRequest) (es []Encoding, err error) {
	var params url.Values
	if er != nil {
		params, err = query.Values(er)
		if err != nil {
			return nil, err
		}
	}
	err = c.get(fmt.Sprintf(videosIdEncodingPath, id), &es, params)
	return
}

// VideoMetaData gets meta data for the video with the given id
func (c *Client) VideoMetaData(id string) (MetaData, error) {
	md := MetaData{}
	if err := c.get(fmt.Sprintf(videosIdMetaDataPath, id), &md, nil); err != nil {
		return nil, err
	}
	return md, nil
}

// DeleteSource deletes the source video for the given video id
func (c *Client) DeleteSource(id string) error {
	_, err := c.PandaClient.Delete(fmt.Sprintf(videosIdDeleteSourcePath, id))
	return err
}

// Notifications gets notifications for the current cloud
func (c *Client) Notifications() (*Notification, error) {
	n := new(Notification)
	if err := c.get(notificationsPath, n, nil); err != nil {
		return nil, err
	}
	return n, nil
}
