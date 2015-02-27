package panda

import "fmt"

type EncodingRequest struct {
	VideoId     string `url:"video_id"`
	ProfileId   string `url:"profile_id,omitempty"`
	ProfileName string `url:"profile_name,omitempty"`
}

type Encoding struct {
	Id                 string   `json:"id"`
	Status             Status   `json:"status"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
	StartedEncoding_at string   `json:"started_encoding_at"`
	EncodingTime       float64  `json:"encoding_time"`
	EncodingProgress   int64    `json:"encoding_progress"`
	VideoId            string   `json:"video_id"`
	ProfileId          string   `json:"profile_id"`
	ProfileName        string   `json:"profile_name"`
	Files              []string `json:"files"`
	MimeType           string   `json:"mime_type"`
	Duration           int64    `json:"duration"`
	Height             int64    `json:"height"`
	Width              int64    `json:"width"`
	Extname            string   `json:"extname"`
	FileSize           float64  `json:"file_size"`
	VideoBitrate       int64    `json:"video_bitrate"`
	AudioBitrate       int64    `json:"audio_bitrate"`
	AudioCodec         string   `json:"audio_codec"`
	VideoCodec         string   `json:"video_codec"`
	Fps                float64  `json:"fps"`
	AudioChannels      int64    `json:"audio_channels"`
	AudioSample_rate   int64    `json:"audio_sample_rate"`
	Path               string   `json:"path"`
	ErrorMessage       string   `json:"error_message"`
	ErrorClass         string   `json:"error_class"`
	ExternalId         string   `json:"external_id"`
	cl                 Client
}

func (e *Encoding) Cancel() (err error) {
	if err = clientError(e.cl); err == nil {
		_, err = e.cl.Post(fmt.Sprintf(encodingsIdCancelPath, e.Id), "", nil, nil)
	}
	return
}

func (e *Encoding) Retry() (err error) {
	if err = clientError(e.cl); err == nil {
		_, err = e.cl.Post(fmt.Sprintf(encodingsIdRetryPath, e.Id), "", nil, nil)
	}
	return
}

func (e *Encoding) Delete() (err error) {
	if err = clientError(e.cl); err == nil {
		_, err = e.cl.Delete(fmt.Sprintf(encodingsIdPath, e.Id))
	}
	return
}

func (e *Encoding) SetClient(cl Client) {
	e.cl = cl
}
