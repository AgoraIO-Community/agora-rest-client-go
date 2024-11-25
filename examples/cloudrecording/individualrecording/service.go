package individualrecording

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	agoraLogger "github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/region"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/base"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	cloudRecordingAPI "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario/individualrecording"
)

type Service struct {
	base.Service
}

func NewService(region region.Area, appId, cname, uid string) *Service {
	return &Service{
		base.Service{
			RegionArea: region,
			AppId:      appId,
			Cname:      cname,
			Uid:        uid,
			Credential: nil,
		},
	}
}

func (s *Service) RunRecording(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		RegionCode: s.RegionArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Println(err)
		return
	}

	// acquire
	acquireResp, err := cloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, false, &individualrecording.AcquireIndividualRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &individualrecording.StartIndividualRecordingClientRequest{
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
		stopResp, err := cloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &individualrecording.UpdateIndividualRecordingClientRequest{
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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

func (s *Service) RunSnapshot(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		RegionCode: s.RegionArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}
	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	// acquire
	acquireResp, err := cloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, false, &individualrecording.AcquireIndividualRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &individualrecording.StartIndividualRecordingClientRequest{
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
		stopResp, err := cloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := cloudRecordingClient.IndividualRecording().QueryVideoScreenshot(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &individualrecording.UpdateIndividualRecordingClientRequest{
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
		queryResp, err := cloudRecordingClient.IndividualRecording().QueryVideoScreenshot(ctx, resourceId, sid)
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

func (s *Service) RunRecordingAndSnapshot(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()
	c := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		RegionCode: s.RegionArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}
	cloudRecordingClient, err := cloudrecording.NewClient(c)
	if err != nil {
		log.Fatalln(err)
	}
	// acquire
	acquireResp, err := cloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, false, &individualrecording.AcquireIndividualRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &individualrecording.StartIndividualRecordingClientRequest{
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
		stopResp, err := cloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &individualrecording.UpdateIndividualRecordingClientRequest{
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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

func (s *Service) RunRecordingAndPostponeTranscoding(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		RegionCode: s.RegionArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	// acquire
	acquireResp, err := cloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, true, &individualrecording.AcquireIndividualRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &individualrecording.StartIndividualRecordingClientRequest{
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
		stopResp, err := cloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &individualrecording.UpdateIndividualRecordingClientRequest{
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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

func (s *Service) RunRecordingAndAudioMix(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		RegionCode: s.RegionArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	// acquire
	acquireResp, err := cloudRecordingClient.IndividualRecording().Acquire(ctx, s.Cname, s.Uid, true, &individualrecording.AcquireIndividualRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.IndividualRecording().Start(ctx, resourceId, s.Cname, s.Uid, &individualrecording.StartIndividualRecordingClientRequest{
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
		stopResp, err := cloudRecordingClient.IndividualRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.IndividualRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &individualrecording.UpdateIndividualRecordingClientRequest{
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
		queryResp, err := cloudRecordingClient.IndividualRecording().Query(ctx, resourceId, sid)
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
