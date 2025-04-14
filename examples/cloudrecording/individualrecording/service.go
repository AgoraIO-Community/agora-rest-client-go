package individualrecording

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/base"
	cloudRecordingAPI "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/req"
)

type Scenario struct {
	*base.Service
}

func NewScenario(s *base.Service) *Scenario {
	return &Scenario{
		Service: s,
	}
}

func (s *Scenario) RunRecording(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()

	// acquire
	acquireResp, err := s.CloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, false, &req.AcquireIndividualRecordingClientRequest{})
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
	startResp, err := s.CloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &req.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &cloudRecordingAPI.RecordingConfig{
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
		RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
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
		stopResp, err := s.CloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := s.CloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &cloudRecordingAPI.UpdateStreamSubscribe{
			AudioUidList: &cloudRecordingAPI.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &cloudRecordingAPI.UpdateVideoUIDList{
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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

func (s *Scenario) RunSnapshot(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()

	// acquire
	acquireResp, err := s.CloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, false, &req.AcquireIndividualRecordingClientRequest{})
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
	startResp, err := s.CloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &req.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &cloudRecordingAPI.RecordingConfig{
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
		SnapshotConfig: &cloudRecordingAPI.SnapshotConfig{
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
		stopResp, err := s.CloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().QueryVideoScreenshot(ctx, resourceId, sid)
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
	updateResp, err := s.CloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &cloudRecordingAPI.UpdateStreamSubscribe{
			AudioUidList: &cloudRecordingAPI.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &cloudRecordingAPI.UpdateVideoUIDList{
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().QueryVideoScreenshot(ctx, resourceId, sid)
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

func (s *Scenario) RunRecordingAndSnapshot(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()
	// acquire
	acquireResp, err := s.CloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, false, &req.AcquireIndividualRecordingClientRequest{})
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
	startResp, err := s.CloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &req.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &cloudRecordingAPI.RecordingConfig{
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
		RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
			AvFileType: []string{
				"hls",
			},
		},
		SnapshotConfig: &cloudRecordingAPI.SnapshotConfig{
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
		stopResp, err := s.CloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := s.CloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &cloudRecordingAPI.UpdateStreamSubscribe{
			AudioUidList: &cloudRecordingAPI.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &cloudRecordingAPI.UpdateVideoUIDList{
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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

func (s *Scenario) RunRecordingAndPostponeTranscoding(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()

	// acquire
	acquireResp, err := s.CloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, true, &req.AcquireIndividualRecordingClientRequest{})
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
	startResp, err := s.CloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &req.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &cloudRecordingAPI.RecordingConfig{
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
		RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
			AvFileType: []string{
				"hls",
			},
		},
		StorageConfig: storageConfig,
		AppsCollection: &cloudRecordingAPI.AppsCollection{
			CombinationPolicy: "postpone_transcoding",
		},
		TranscodeOptions: &cloudRecordingAPI.TranscodeOptions{
			Container: &cloudRecordingAPI.Container{
				Format: "mp4",
			},
			TransConfig: &cloudRecordingAPI.TransConfig{
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
		stopResp, err := s.CloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := s.CloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &cloudRecordingAPI.UpdateStreamSubscribe{
			AudioUidList: &cloudRecordingAPI.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &cloudRecordingAPI.UpdateVideoUIDList{
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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

func (s *Scenario) RunRecordingAndAudioMix(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()

	// acquire
	acquireResp, err := s.CloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, true, &req.AcquireIndividualRecordingClientRequest{})
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
	startResp, err := s.CloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &req.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &cloudRecordingAPI.RecordingConfig{
			ChannelType: 1,
			StreamTypes: 0,
			StreamMode:  "original",
			MaxIdleTime: 30,
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeUidGroup: 0,
		},
		RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
			AvFileType: []string{
				"hls",
			},
		},
		StorageConfig: storageConfig,
		AppsCollection: &cloudRecordingAPI.AppsCollection{
			CombinationPolicy: "postpone_transcoding",
		},
		TranscodeOptions: &cloudRecordingAPI.TranscodeOptions{
			Container: &cloudRecordingAPI.Container{
				Format: "mp3",
			},
			TransConfig: &cloudRecordingAPI.TransConfig{
				TransMode: "audioMix",
			},
			Audio: &cloudRecordingAPI.Audio{
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
		stopResp, err := s.CloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := s.CloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &cloudRecordingAPI.UpdateStreamSubscribe{
			AudioUidList: &cloudRecordingAPI.UpdateAudioUIDList{
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
		queryResp, err := s.CloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
