// Package panda provides a client for PandaStream service
package panda

//go:generate structgen -dir=json -tags=url -o=models.go -pkg=panda -types=created_at:Time,updated_at:Time,height:int,width:int,duration:int,file_size:int64,audio_bitrate:int,audio_channels:int,video_bitrate:int,audio_sample_rate:int,watermark_bottom:int,watermark_height:int,watermark_left:int,watermark_right:int,watermark_top:int,watermark_width:int,keyframe_interval:int,buffer_size:int,max_rate:int,frame_count:int,h264_crf:int,status:Status,page:int,per_page:int,aspect_mode:AspectMode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"time"

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

// Manager provides methods of getting certain data from the Panda Cloud using panda Client
type Manager struct {
	Client *Client
}

func (m *Manager) manageGet(url string, v interface{}, params url.Values) (err error) {
	b, err := m.Client.Get(url, params)
	if err == nil {
		err = json.Unmarshal(b, v)
	}
	return
}

func (m *Manager) managePost(url string, r io.Reader, p, v interface{}) error {
	params, err := query.Values(p)
	if err != nil {
		return err
	}
	b, err := m.Client.Post(url, "", params, r)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// Cloud gets cloud by the given cloud ID
func (m *Manager) Cloud(id string) (*Cloud, error) {
	c := new(Cloud)
	if err := m.manageGet(fmt.Sprintf(cloudsIdPath, id), c, nil); err != nil {
		return nil, err
	}
	return c, nil
}

// Clouds gets all clouds on the given account
func (m *Manager) Clouds() (cs []Cloud, err error) {
	err = m.manageGet(cloudsPath, &cs, nil)
	return
}

// NewEncoding creates a new encoding for the existing video
func (m *Manager) NewEncoding(er *NewEncodingRequest) (*Encoding, error) {
	e := new(Encoding)
	if err := m.managePost(encodingsPath, nil, er, &e); err != nil {
		return nil, err
	}
	return e, nil
}

// Encoding gets encoding object with the given id
func (m *Manager) Encoding(id string) (*Encoding, error) {
	e := new(Encoding)
	if err := m.manageGet(fmt.Sprintf(encodingsIdPath, id), e, nil); err != nil {
		return nil, err
	}
	return e, nil
}

// Encodings get all encodings on the current cloud. Encodings can be filtered by options
// set in the EncodingRequest struct. If EncodingRequest is nil defaults are going to be used
func (m *Manager) Encodings(er *EncodingRequest) (es []Encoding, err error) {
	var params url.Values
	if er != nil {
		params, err = query.Values(er)
		if err != nil {
			return
		}
	}
	err = m.manageGet(encodingsPath, &es, params)
	return
}

// Cancel encoding with the given id
func (m *Manager) Cancel(id string) error {
	_, err := m.Client.Post(fmt.Sprintf(encodingsIdCancelPath, id), "", nil, nil)
	return err
}

// Retry encoding with the given id
func (m *Manager) Retry(id string) error {
	_, err := m.Client.Post(fmt.Sprintf(encodingsIdRetryPath, id), "", nil, nil)
	return err
}

// Delete accepts *Profile, *Video and *Encoding types and deletes those objects
// from Panda's database by their ID
func (m *Manager) Delete(v interface{}) error {
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
	_, err := m.Client.Delete(path)
	return err
}

// NewProfile creates new profile based on profile request object
func (m *Manager) NewProfile(pr *NewProfileRequest) (*Profile, error) {
	p := new(Profile)
	if err := m.managePost(profilesPath, nil, pr, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Profile gets profile with the given ID
func (m *Manager) Profile(id string) (*Profile, error) {
	p := new(Profile)
	if err := m.manageGet(fmt.Sprintf(profilesIdPath, id), p, nil); err != nil {
		return nil, err
	}
	return p, nil
}

// Profiles gets all profiles from the current cloud
func (m *Manager) Profiles(pr *ProfileRequest) (ps []Profile, err error) {
	var params url.Values
	if pr != nil {
		params, err = query.Values(pr)
		if err != nil {
			return
		}
	}
	err = m.manageGet(profilesPath, &ps, params)
	return
}

// Update accepts *Profile and *Notification types and updates records based on the given objects.
// Warning: the given parameter might change if any of the parameters are invalid
func (m *Manager) Update(v interface{}) error {
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
	b, err := m.Client.Put(path, "", params, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// NewVideo gets file name path and potential additional options in a form of
// VideoRequest type and based on this creates a new video in Panda
func (m *Manager) NewVideo(file string, vr *NewVideoRequest) (*Video, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return m.NewVideoReader(f, file, vr)
}

// NewVideoURL creates a new video in Panda based on the given source url and potential
// additional options in form of VideoRequest struct
func (m *Manager) NewVideoURL(URL string, vr *NewVideoRequest) (*Video, error) {
	params, err := query.Values(vr)
	if err != nil {
		return nil, err
	}
	params.Set("source_url", URL)
	b, err := m.Client.Post(videosPath, "", params, nil)
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
func (m *Manager) NewVideoReader(r io.Reader, name string, vr *NewVideoRequest) (*Video, error) {
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
	b, err := m.Client.Post(videosPath, w.FormDataContentType(), params, buf)
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
func (m *Manager) Video(id string) (*Video, error) {
	v := new(Video)
	if err := m.manageGet(fmt.Sprintf(videosIdPath, id), &v, nil); err != nil {
		return nil, err
	}
	return v, nil
}

// Videos gets all the videos from the current cloud. Videos can be filtered by options
// set in the VideoRequest struct. If VideoRequest is nil then defaults are going to be used
func (m *Manager) Videos(vr *VideoRequest) (v []Video, err error) {
	var params url.Values
	if vr != nil {
		params, err = query.Values(vr)
		if err != nil {
			return nil, err
		}
	}
	err = m.manageGet(videosPath, &v, params)
	return
}

// VideoEncoding get alls encodings related to the video with the given id. Encodings can be
// filtered by options set in the EncodingRequest struct.
// If EncodingRequest is nil defaults are going to be used
func (m *Manager) VideoEncodings(id string, er *EncodingRequest) (es []Encoding, err error) {
	var params url.Values
	if er != nil {
		params, err = query.Values(er)
		if err != nil {
			return nil, err
		}
	}
	err = m.manageGet(fmt.Sprintf(videosIdEncodingPath, id), &es, params)
	return
}

// VideoMetaData gets meta data for the video with the given id
func (m *Manager) VideoMetaData(id string) (MetaData, error) {
	md := MetaData{}
	if err := m.manageGet(fmt.Sprintf(videosIdMetaDataPath, id), &md, nil); err != nil {
		return nil, err
	}
	return md, nil
}

// DeleteSource deletes the source video for the given video id
func (m *Manager) DeleteSource(id string) error {
	_, err := m.Client.Delete(fmt.Sprintf(videosIdDeleteSourcePath, id))
	return err
}

// Notifications gets notifications for the current cloud
func (m *Manager) Notifications() (*Notification, error) {
	n := new(Notification)
	if err := m.manageGet(notificationsPath, n, nil); err != nil {
		return nil, err
	}
	return n, nil
}
