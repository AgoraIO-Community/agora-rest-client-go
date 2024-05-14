package individualrecording

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Service struct {
	region     core.RegionArea
	appId      string
	cname      string
	uid        string
	credential core.Credential
}

func NewService(region core.RegionArea, appId, cname, uid string) *Service {
	return &Service{
		region:     region,
		appId:      appId,
		cname:      cname,
		uid:        uid,
		credential: nil,
	}
}

func (s *Service) SetCredential(username, password string) {
	s.credential = core.NewBasicAuthCredential(username, password)
}

func (s *Service) RunRecording(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	impl := cloudrecording.NewAPI(c).V1().IndividualRecording()
	// acquire
	acquireResp, err := impl.Acquire().Do(ctx, s.cname, s.uid, false, &v1.AcquireIndividualRecodingClientRequest{})
	if err != nil {
		log.Println(err)
		return
	}
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Printf("acquire failed:%+v\n", acquireResp)
		return
	}

	resourceId := acquireResp.SuccessRes.ResourceId
	// start
	startResp, err := impl.Start().Do(ctx, resourceId, s.cname, s.uid, &v1.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType: 1,
			StreamTypes: 2,
			MaxIdleTime: 30,
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeVideoUIDs: []string{
				"#allstream#",
			},
			SubscribeUidGroup: 0,
		},
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
			},
		},
		StorageConfig: storageConfig,
	})
	if err != nil {
		log.Println(err)
		return
	}
	if startResp.IsSuccess() {
		log.Printf("start success:%+v\n", startResp)
	} else {
		log.Printf("start failed:%+v\n", startResp)
		return
	}

	sid := startResp.SuccessResponse.Sid
	// stop
	defer func() {
		stopResp, err := impl.Stop().Do(ctx, resourceId, sid, s.cname, s.uid, false)
		if err != nil {
			log.Println(err)
			return
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v\n", stopResp)
		} else {
			log.Printf("stop failed:%+v\n", stopResp)
			return
		}
	}()

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}

	// update
	updateResp, err := impl.Update().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &v1.UpdateStreamSubscribe{
			AudioUidList: &v1.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &v1.UpdateVideoUIDList{
				SubscribeVideoUIDs: []string{
					"#allstream#",
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	if updateResp.IsSuccess() {
		log.Printf("update success:%+v\n", updateResp)
	} else {
		log.Printf("update failed:%+v\n", updateResp)
		return
	}

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}
}

func (s *Service) RunSnapshot(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	impl := cloudrecording.NewAPI(c).V1().IndividualRecording()
	// acquire
	acquireResp, err := impl.Acquire().Do(ctx, s.cname, s.uid, false, &v1.AcquireIndividualRecodingClientRequest{})
	if err != nil {
		log.Println(err)
		return
	}
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Printf("acquire failed:%+v\n", acquireResp)
		return
	}

	resourceId := acquireResp.SuccessRes.ResourceId
	// start
	startResp, err := impl.Start().Do(ctx, resourceId, s.cname, s.uid, &v1.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType: 1,
			StreamTypes: 2,
			MaxIdleTime: 30,
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeVideoUIDs: []string{
				"#allstream#",
			},
			SubscribeUidGroup: 0,
		},
		SnapshotConfig: &v1.SnapshotConfig{
			CaptureInterval: 5,
			FileType:        []string{"jpg"},
		},
		StorageConfig: storageConfig,
	})
	if err != nil {
		log.Println(err)
		return
	}
	if startResp.IsSuccess() {
		log.Printf("start success:%+v\n", startResp)
	} else {
		log.Printf("start failed:%+v\n", startResp)
		return
	}

	sid := startResp.SuccessResponse.Sid
	// stop
	defer func() {
		stopResp, err := impl.Stop().DoVideoScreenshot(ctx, resourceId, sid, s.cname, s.uid, false)
		if err != nil {
			log.Println(err)
			return
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v\n", stopResp)
		} else {
			log.Printf("stop failed:%+v\n", stopResp)
			return
		}
	}()

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().DoVideoScreenshot(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}

	// update
	updateResp, err := impl.Update().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &v1.UpdateStreamSubscribe{
			AudioUidList: &v1.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &v1.UpdateVideoUIDList{
				SubscribeVideoUIDs: []string{
					"#allstream#",
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	if updateResp.IsSuccess() {
		log.Printf("update success:%+v\n", updateResp)
	} else {
		log.Printf("update failed:%+v\n", updateResp)
		return
	}

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().DoVideoScreenshot(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}
}

func (s *Service) RunRecordingAndSnapshot(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	impl := cloudrecording.NewAPI(c).V1().IndividualRecording()
	// acquire
	acquireResp, err := impl.Acquire().Do(ctx, s.cname, s.uid, false, &v1.AcquireIndividualRecodingClientRequest{})
	if err != nil {
		log.Println(err)
		return
	}
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Printf("acquire failed:%+v\n", acquireResp)
		return
	}

	resourceId := acquireResp.SuccessRes.ResourceId
	// start
	startResp, err := impl.Start().Do(ctx, resourceId, s.cname, s.uid, &v1.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType: 1,
			StreamTypes: 2,
			MaxIdleTime: 30,
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeVideoUIDs: []string{
				"#allstream#",
			},
			SubscribeUidGroup: 0,
		},
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
			},
		},
		SnapshotConfig: &v1.SnapshotConfig{
			CaptureInterval: 5,
			FileType:        []string{"jpg"},
		},
		StorageConfig: storageConfig,
	})
	if err != nil {
		log.Println(err)
		return
	}
	if startResp.IsSuccess() {
		log.Printf("start success:%+v\n", startResp)
	} else {
		log.Printf("start failed:%+v\n", startResp)
		return
	}

	sid := startResp.SuccessResponse.Sid
	// stop
	defer func() {
		stopResp, err := impl.Stop().Do(ctx, resourceId, sid, s.cname, s.uid, false)
		if err != nil {
			log.Println(err)
			return
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v\n", stopResp)
		} else {
			log.Printf("stop failed:%+v\n", stopResp)
			return
		}
	}()

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}

	// update
	updateResp, err := impl.Update().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &v1.UpdateStreamSubscribe{
			AudioUidList: &v1.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &v1.UpdateVideoUIDList{
				SubscribeVideoUIDs: []string{
					"#allstream#",
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	if updateResp.IsSuccess() {
		log.Printf("update success:%+v\n", updateResp)
	} else {
		log.Printf("update failed:%+v\n", updateResp)
		return
	}

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}
}

func (s *Service) RunRecordingAndPostponeTranscoding(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	impl := cloudrecording.NewAPI(c).V1().IndividualRecording()
	// acquire
	acquireResp, err := impl.Acquire().Do(ctx, s.cname, s.uid, true, &v1.AcquireIndividualRecodingClientRequest{})
	if err != nil {
		log.Println(err)
		return
	}
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Printf("acquire failed:%+v\n", acquireResp)
		return
	}

	resourceId := acquireResp.SuccessRes.ResourceId
	// start
	startResp, err := impl.Start().Do(ctx, resourceId, s.cname, s.uid, &v1.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType: 1,
			StreamTypes: 2,
			MaxIdleTime: 30,
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeVideoUIDs: []string{
				"#allstream#",
			},
			SubscribeUidGroup: 0,
		},
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
			},
		},
		StorageConfig: storageConfig,
		AppsCollection: &v1.AppsCollection{
			CombinationPolicy: "postpone_transcoding",
		},
		TranscodeOptions: &v1.TranscodeOptions{
			Container: &v1.Container{
				Format: "mp4",
			},
			TransConfig: &v1.TransConfig{
				TransMode: "postponeTranscoding",
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	if startResp.IsSuccess() {
		log.Printf("start success:%+v\n", startResp)
	} else {
		log.Printf("start failed:%+v\n", startResp)
		return
	}

	sid := startResp.SuccessResponse.Sid
	// stop
	defer func() {
		stopResp, err := impl.Stop().Do(ctx, resourceId, sid, s.cname, s.uid, false)
		if err != nil {
			log.Println(err)
			return
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v\n", stopResp)
		} else {
			log.Printf("stop failed:%+v\n", stopResp)
			return
		}
	}()

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}

	// update
	updateResp, err := impl.Update().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &v1.UpdateStreamSubscribe{
			AudioUidList: &v1.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &v1.UpdateVideoUIDList{
				SubscribeVideoUIDs: []string{
					"#allstream#",
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	if updateResp.IsSuccess() {
		log.Printf("update success:%+v\n", updateResp)
	} else {
		log.Printf("update failed:%+v\n", updateResp)
		return
	}

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}
}

func (s *Service) RunRecordingAndAudioMix(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	impl := cloudrecording.NewAPI(c).V1().IndividualRecording()
	// acquire
	acquireResp, err := impl.Acquire().Do(ctx, s.cname, s.uid, true, &v1.AcquireIndividualRecodingClientRequest{})
	if err != nil {
		log.Println(err)
		return
	}
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Printf("acquire failed:%+v\n", acquireResp)
		return
	}

	resourceId := acquireResp.SuccessRes.ResourceId
	// start
	startResp, err := impl.Start().Do(ctx, resourceId, s.cname, s.uid, &v1.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType: 1,
			StreamTypes: 0,
			StreamMode:  "original",
			MaxIdleTime: 30,
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeUidGroup: 0,
		},
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
			},
		},
		StorageConfig: storageConfig,
		AppsCollection: &v1.AppsCollection{
			CombinationPolicy: "postpone_transcoding",
		},
		TranscodeOptions: &v1.TranscodeOptions{
			Container: &v1.Container{
				Format: "mp3",
			},
			TransConfig: &v1.TransConfig{
				TransMode: "audioMix",
			},
			Audio: &v1.Audio{
				SampleRate: "48000",
				BitRate:    "48000",
				Channels:   "2",
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	if startResp.IsSuccess() {
		log.Printf("start success:%+v\n", startResp)
	} else {
		log.Printf("start failed:%+v\n", startResp)
		return
	}

	sid := startResp.SuccessResponse.Sid
	// stop
	defer func() {
		stopResp, err := impl.Stop().Do(ctx, resourceId, sid, s.cname, s.uid, false)
		if err != nil {
			log.Println(err)
			return
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v\n", stopResp)
		} else {
			log.Printf("stop failed:%+v\n", stopResp)
			return
		}
	}()

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}

	// update
	updateResp, err := impl.Update().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &v1.UpdateStreamSubscribe{
			AudioUidList: &v1.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	if updateResp.IsSuccess() {
		log.Printf("update success:%+v\n", updateResp)
	} else {
		log.Printf("update failed:%+v\n", updateResp)
		return
	}

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().Do(ctx, resourceId, sid)
		if err != nil {
			log.Println(err)
			return
		}
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
			return
		}
		time.Sleep(time.Second * 10)
	}
}
