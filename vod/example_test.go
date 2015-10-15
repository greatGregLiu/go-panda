package vod_test

import (
	"fmt"

	"github.com/pandastream/go-panda/client"
	"github.com/pandastream/go-panda/vod"
)

func ExampleManager_newEncoding() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	encoding, err := cl.NewEncoding(&vod.NewEncodingRequest{
		ProfileID: "profile_id",
		VideoID:   "video_id",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(encoding)
}

func ExampleManager_encodings() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	encodings, err := cl.Encodings(&vod.EncodingRequest{
		Status:  vod.StatusProcessing,
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(encodings)
}

func ExampleManager_newProfile() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	profile, err := cl.NewProfile(&vod.NewProfileRequest{
		PresetName:   "h264",
		AspectMode:   vod.ModeLetterBox,
		Fps:          29.7,
		Name:         "My_New_Profile",
		Width:        640,
		Height:       480,
		AddTimestamp: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(profile)
}

func ExampleManager_delete() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	profile, err := cl.Profile("profile_id")
	if err != nil {
		panic(err)
	}
	if err := cl.Delete(&profile); err != nil {
		panic(err)
	}
}

func ExampleManager_update() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	profile, err := cl.Profile("profile_id")
	if err != nil {
		panic(err)
	}
	// ...
	profile.AspectMode = vod.ModePreserve
	// ...
	if err := cl.Update(&profile); err != nil {
		panic(err)
	}
}

func ExampleManager_newVideo() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	video, err := cl.NewVideo("filepath", &vod.NewVideoRequest{
		Profiles: []string{"Profile1", "Profile2"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(video)
}

func ExampleManager_newVideoURL() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	video, err := cl.NewVideoURL("source_url", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(video)
}

func ExampleManager_videos() {
	cl := vod.NewClient(client.HostGCE, "token", nil)
	videos, err := cl.Videos(&vod.VideoRequest{
		Status:  vod.StatusProcessing,
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(videos)
}
