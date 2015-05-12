package panda_test

import (
	"fmt"

	panda "github.com/pandastream/go-panda"
)

func ExampleManager_newEncoding() {
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	encoding, err := m.NewEncoding(&panda.NewEncodingRequest{
		ProfileID: "profile_id",
		VideoID:   "video_id",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(encoding)
}

func ExampleManager_encodings() {
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	encodings, err := m.Encodings(&panda.EncodingRequest{
		Status:  panda.StatusProcessing,
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(encodings)
}

func ExampleManager_newProfile() {
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	profile, err := m.NewProfile(&panda.NewProfileRequest{
		PresetName:   "h264",
		AspectMode:   panda.ModeLetterBox,
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
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	profile, err := m.Profile("profile_id")
	if err != nil {
		panic(err)
	}
	if err := m.Delete(&profile); err != nil {
		panic(err)
	}
}

func ExampleManager_update() {
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	profile, err := m.Profile("profile_id")
	if err != nil {
		panic(err)
	}
	// ...
	profile.AspectMode = panda.ModePreserve
	// ...
	if err := m.Update(&profile); err != nil {
		panic(err)
	}
}

func ExampleManager_newVideo() {
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	video, err := m.NewVideo("filepath", &panda.NewVideoRequest{
		Profiles: []string{"Profile1", "Profile2"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(video)
}

func ExampleManager_newVideoURL() {
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	video, err := m.NewVideoURL("source_url", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(video)
}

func ExampleManager_videos() {
	m := panda.Manager{
		&panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
			},
		},
	}

	videos, err := m.Videos(&panda.VideoRequest{
		Status:  panda.StatusProcessing,
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(videos)
}
