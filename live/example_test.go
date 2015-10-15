package live_test

import (
	"fmt"
	"time"

	"github.com/pandastream/go-panda/live"
)

func ExampleClient_profileCreate() {
	cl := live.NewClient("token", nil)
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
	profileID, err := cl.ProfileCreate(&profile)
	if err != nil {
		panic(err)
	}
	fmt.Println(profileID)
}

func ExampleClient_streamCreate() {
	cl := live.NewClient("token", nil)
	stream := &live.Stream{
		ProfileID: "999999", // existing profile_id
		Duration:  10,       // 10 minutes
	}
	streamID, err := cl.StreamCreate(stream)
	if err != nil {
		panic(err)
	}
	fmt.Println(streamID)
}

func ExampleClient_streamCreateProfile() {
	cl := live.NewClient("token", nil)
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
	streamID, profileID, err := cl.StreamCreateProfile("stream_name", profile)
	if err != nil {
		panic(err)
	}
	fmt.Println(streamID, profileID)
}

func ExampleClient_streamDuration() {
	cl := live.NewClient("token", nil)
	streamID, err := cl.StreamDuration("999999", time.Minute*5)
	if err != nil {
		panic(err)
	}
	fmt.Println(streamID)
}
