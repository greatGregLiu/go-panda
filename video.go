package panda

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type MetaData map[string]interface{}

type VideoRequest struct {
	Profiles   []string `url:"profiles,omitempty,comma"`
	PathFormat string   `url:"path_format,omitempty"`
	Payload    string   `url:"payload,omitempty"`
}

type VideoRequestUrl struct {
	VideoRequest
	Url string `url:"source_url"`
}

type Video struct {
	Id               string  `json:"id"`
	Status           Status  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	MimeType         string  `json:"mime_type"`
	OriginalFilename string  `json:"original_filename"`
	SourceUrl        string  `json:"source_url"`
	Duration         float64 `json:"duration"`
	Height           int64   `json:"height"`
	Width            int64   `json:"width"`
	Extname          string  `json:"extname"`
	FileSize         float64 `json:"file_size"`
	VideoBitrate     int64   `json:"video_bitrate"`
	AudioBitrate     int64   `json:"audio_bitrate"`
	AudioCodec       string  `json:"audio_codec"`
	VideoCodec       string  `json:"video_codec"`
	Fps              float64 `json:"fps"`
	AudioChannels    int64   `json:"audio_channels"`
	AudioSampleRate  int64   `json:"audio_sample_rate"`
	Path             string  `json:"path"`
	ErrorMessage     string  `json:"error_message"`
	ErrorClass       string  `json:"error_class"`
	Payload          string  `json:"payload"`
	cl               Client
}

func (v *Video) Encodings() (es []Encoding, err error) {
	if err = clientError(v.cl); err != nil {
		return
	}
	b, err := v.cl.Get(fmt.Sprintf(videosIdEncodingPath, v.Id), nil)
	if err == nil {
		err = json.Unmarshal(b, &es)
	}
	for i := range es {
		es[i].cl = v.cl
	}
	return
}

func (v *Video) EncodingsStatus(s Status) (es []Encoding, err error) {
	val := url.Values{}
	val.Set("status", strings.ToLower(string(s)))
	b, err := v.cl.Get(fmt.Sprintf(videosIdEncodingPath, v.Id), val)
	if err == nil {
		err = json.Unmarshal(b, &es)
	}
	for i := range es {
		es[i].cl = v.cl
	}
	return
}

func (v *Video) MetaData() (md MetaData, err error) {
	if err = clientError(v.cl); err != nil {
		return
	}
	md = MetaData{}
	b, err := v.cl.Get(fmt.Sprintf(videosIdMetaDataPath, v.Id), nil)
	if err == nil {
		err = json.Unmarshal(b, &md)
	}
	return
}

func (v *Video) Delete() (err error) {
	if err = clientError(v.cl); err == nil {
		_, err = v.cl.Delete(fmt.Sprintf(videosIdDeletePath, v.Id))
	}
	return
}

func (v *Video) DeleteSource() (err error) {
	if err = clientError(v.cl); err == nil {
		_, err = v.cl.Delete(fmt.Sprintf(videosIdDeleteSourcePath, v.Id))
	}
	return
}

func (v *Video) SetClient(cl Client) {
	v.cl = cl
}
