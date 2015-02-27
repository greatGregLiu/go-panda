# go-panda
Golang library for Panda

# Usage

*Create Panda client*

Client for Panda US
```
cl := panda.NewClient("cloud_id", "access_key", "secret_key")
```

Client for Panda EU
```
cl := panda.NewClientEU("cloud_id", "access_key", "secret_key")
```

*Create manager*
```
m := panda.NewManager(cl)
```

*Create new video*
```
v, err := m.NewVideo(&panda.VideoRequestUrl{
  Url: "www.example.com/myvideo.mp4",
  VideoRequest: panda.VideoRequest: {Profiles: "h264"},
})
```

```
f, err := os.Open("path/to/the/file.mp4")
if err != nil {
  panic(err)
}
v, err := m.NewVideoReader(f, "DesiredFileName", &panda.VideoRequest{
  Profiles: "h264",
})
```

*Get Video by Id*
```
  v, err := m.VideoId("id")
```

*Get Video's Encodings*
```
  es, err := v.Encodings()
```

*Create new profile*
```
	p, err := m.NewProfile(&panda.ProfileRequest{
		Name:         "h666",
		Title:        "myProfile",
		PresetName:   "h264",
		Width:        100,
		Height:       100,
		AddTimestamp: true,
		AspectMode:   panda.Crop,
	})
```

*Edit profile*
```
  p.Title = "MyNewTitle"
  err = p.Update()
```

*Get Clouds*
```
  clouds, err := m.Clouds()
```

*Get Cloud by Id*
```
  cloud, err := m.CloudId("id")
```

