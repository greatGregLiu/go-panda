package panda

import (
	"encoding/json"
	"fmt"

	"github.com/ernesto-jimenez/go-querystring/query"
)

type AspectMode string

const (
	LetterBox AspectMode = "letterbox"
	Preserve             = "preserve"
	Constrain            = "constrain"
	Pad                  = "pad"
	Crop                 = "crop"
)

type ProfileRequest struct {
	Name             string     `json:"name,omitempty" url:"name,omitempty"`
	Title            string     `json:"title,omitempty" url:"title,omitempty"`
	Extname          string     `json:"extname,omitempty" url:"extname,omitempty"`
	Width            int64      `json:"width,omitempty" url:"width,omitempty"`
	Height           int64      `json:"height,omitempty" url:"height,omitempty"`
	Upscale          bool       `json:"upscale,omitempty" url:"upscale,omitempty"`
	AspectMode       AspectMode `json:"aspect_mode,omitempty" url:"aspect_mode,omitempty"`
	TwoPass          string     `json:"two_pass,omitempty" url:"two_pass,omitempty"`
	VideoBitrate     int64      `json:"video_bitrate,omitempty" url:"video_bitrate,omitempty"`
	Fps              float64    `json:"fps,omitempty" url:"fps,omitempty"`
	KeyFrameInterval int64      `json:"keyframe_interval,omitempty" url:"keyframe_interval,omitempty"`
	KeyFrameRate     float64    `json:"Keyframe_rate,omitempty" url:"Keyframe_rate,omitempty"`
	AudioBitrate     int64      `json:"audio_bitrate,omitempty" url:"audio_bitrate,omitempty"`
	AudioSampleRate  int64      `json:"audio_sample_rate,omitempty" url:"audio_sample_rate,omitempty"`
	AudioChannels    int64      `json:"audio_channels,omitempty" url:"audio_channels,omitempty"`
	ClipLength       string     `json:"clip_length,omitempty" url:"clip_length,omitempty"`
	ClipOffset       string     `json:"clip_offset,omitempty" url:"clip_offset,omitempty"`
	WatermarkUrl     string     `json:"watermark_url,omitempty" url:"watermark_url,omitempty"`
	WatermarkLeft    int64      `json:"watermark_left,omitempty" url:"watermark_left,omitempty"`
	WatermarkRight   int64      `json:"watermark_right,omitempty" url:"watermark_right,omitempty"`
	WatermarkTop     int64      `json:"watermark_top,omitempty" url:"watermark_top,omitempty"`
	WatermarkBottom  int64      `json:"watermark_bottom,omitempty" url:"watermark_bottom,omitempty"`
	WatermarkWidth   int64      `json:"watermark_width,omitempty" url:"watermark_width,omitempty"`
	WatermarkHeight  int64      `json:"watermark_height,omitempty" url:"watermark_height,omitempty"`
	FrameCount       int64      `json:"frame_count,omitempty" url:"frame_count,omitempty"`
	FrameOffsets     string     `json:"frame_offsets,omitempty" url:"frame_offsets,omitempty"`
	FrameInterval    string     `json:"frame_interval,omitempty" url:"frame_interval,omitempty"`
	Command          string     `json:"command,omitempty" url:"command,omitempty"`
	PresetName       string     `json:"preset_name,omitempty" url:"preset_name,omitempty"`
	H264Level        string     `json:"h264_level,omitempty" url:"h264_level,omitempty"`
	H264Profile      string     `json:"h264_profile,omitempty" url:"h264_profile,omitempty"`
	H264Tune         string     `json:"h264_tune,omitempty" url:"h264_tune,omitempty"`
	H264Crf          int64      `json:"h264_crf,omitempty" url:"h264_crf,omitempty"`
	Encryption       bool       `json:"encryption,omitempty" url:"encryption,omitempty"`
	EncryptionKeyUrl string     `json:"encryption_key_url,omitempty" url:"encryption_key_url,omitempty"`
	EncryptionKey    string     `json:"encryption_key,omitempty" url:"encryption_key,omitempty"`
	EncryptionIv     string     `json:"encryption_iv,omitempty" url:"encryption_iv,omitempty"`
	AddTimestamp     bool       `json:"add_timestamp,omitempty" url:"add_timestamp,omitempty"`
	Deinterlace      string     `json:"deinterlace,omitempty" url:"deinterlace,omitempty"`
	Stack            string     `json:"stack,omitempty" url:"stack,omitempty"`
	BufferSize       int64      `json:"buffer_size,omitempty" url:"buffer_size,omitempty"`
	MaxRate          int64      `json:"max_rate,omitempty" url:"max_rate,omitempty"`
}

type Profile struct {
	ProfileRequest
	Id          string `json:"id" url:"id"`
	CreatedAt   string `json:"created_at" url:"created_at"`
	UpdatedAt   string `json:"updated_at" url:"updated_at"`
	Variants    string `json:"variants,omitempty" url:"variants,omitempty"`
	FullCommand string `json:"full_command,omitempty" url:"full_command,omitempty"`
	cl          Client
}

func (p *Profile) Delete() (err error) {
	if err = clientError(p.cl); err == nil {
		_, err = p.cl.Delete(fmt.Sprintf(profilesIdPath, p.Id))
	}
	return
}

func (p *Profile) Update() (err error) {
	if err = clientError(p.cl); err != nil {
		return
	}
	params, err := query.Values(p)
	if err != nil {
		return err
	}
	b, err := p.cl.Put(fmt.Sprintf(profilesIdPath, p.Id), "", params, nil)
	if err == nil {
		err = json.Unmarshal(b, p)
	}
	return
}

func (p *Profile) SetClient(cl Client) {
	p.cl = cl
}
