package live_test

import (
	"fmt"
	"time"

	"github.com/pandastream/go-panda"
	"github.com/pandastream/go-panda/live"
)

func ExampleClient_profileCreate() {
	client := live.Client{
		Client: &panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
				Namespace: "live",
			},
		},
	}

	nodes := live.Nodes{
		"ingester": live.Node{
			Type: "rtmp_ingest",
			Config: map[string]interface{}{
				"app": "in",
			},
		},
		"hls_variant": live.Node{
			Type:    "hls",
			Sources: []string{"ingester"},
			Config: map[string]interface{}{
				"app":       "hello1",
				"bandwidth": 1000,
			},
		},
		"hls_top": live.Node{
			Type:    "hls_master",
			Sources: []string{"hls_variant"},
			Config: map[string]interface{}{
				"app": "hello",
			},
		},
	}

	profile := live.Profile{
		Nodes: nodes,
	}
	profileID, err := client.ProfileCreate(&profile)
	if err != nil {
		panic(err)
	}
	fmt.Println(profileID)
}

func ExampleClient_streamCreate() {
	client := live.Client{
		Client: &panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
				Namespace: "live",
			},
		},
	}

	stream := &live.Stream{
		ProfileID: "999999", // existing profile_id
		Duration:  10,       // 10 minutes
	}

	streamID, err := client.StreamCreate(stream)
	if err != nil {
		panic(err)
	}
	fmt.Println(streamID)
}

func ExampleClient_streamCreateProfile() {
	client := live.Client{
		Client: &panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
				Namespace: "live",
			},
		},
	}

	nodes := live.Nodes{
		"ingester": live.Node{
			Type: "rtmp_ingest",
			Config: map[string]interface{}{
				"app": "in",
			},
		},
		"hls_variant": live.Node{
			Type:    "hls",
			Sources: []string{"ingester"},
			Config: map[string]interface{}{
				"app":       "hello1",
				"bandwidth": 1000,
			},
		},
		"hls_top": live.Node{
			Type:    "hls_master",
			Sources: []string{"hls_variant"},
			Config: map[string]interface{}{
				"app": "hello",
			},
		},
	}

	profile := &live.Profile{
		Nodes: nodes,
	}

	streamID, profileID, err := client.StreamCreateProfile("stream_name", profile)
	if err != nil {
		panic(err)
	}
	fmt.Println(streamID, profileID)
}

func ExampleClient_streamDuration() {
	client := live.Client{
		Client: &panda.Client{
			Host: panda.HostGCE,
			Options: &panda.ClientOptions{
				AccessKey: "access_key",
				SecretKey: "secret_key",
				CloudID:   "cloud_id",
				Namespace: "live",
			},
		},
	}

	streamID, err := client.StreamDuration("999999", time.Minute*5)
	if err != nil {
		panic(err)
	}
	fmt.Println(streamID)
}
