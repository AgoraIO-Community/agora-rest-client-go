package mixrecording

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

func (s *Service) RunHLS(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	impl := cloudrecording.NewAPI(c).V1().MixRecording()
	// acquire
	acquireResp, err := impl.Acquire().Do(ctx, s.cname, s.uid, &v1.AcquireMixRecodingClientRequest{})
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
	startResp, err := impl.Start().Do(ctx, resourceId, s.cname, s.uid, &v1.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			MaxIdleTime:  30,
			AudioProfile: 2,
			TranscodingConfig: &v1.TranscodingConfig{
				Width:            640,
				Height:           640,
				FPS:              15,
				BitRate:          800,
				MixedVideoLayout: 0,
				BackgroundColor:  "#000000",
			},
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeVideoUIDs: []string{
				"#allstream#",
			},
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

	sid := startResp.SuccessResp.Sid
	// stop
	defer func() {
		stopResp, err := impl.Stop().DoHLS(ctx, resourceId, sid, s.cname, s.uid, false)
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
		queryResp, err := impl.Query().DoHLS(ctx, resourceId, sid)
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
	updateResp, err := impl.Update().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateMixRecordingClientRequest{
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

	// updateLayout
	updateLayoutResp, err := impl.UpdateLayout().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateLayoutUpdateMixRecordingClientRequest{
		MixedVideoLayout: 1,
		BackgroundColor:  "#FF0000",
	},
	)
	if err != nil {
		log.Println(err)
		return
	}
	if updateLayoutResp.IsSuccess() {
		log.Printf("updateLayout success:%+v\n", updateLayoutResp)
	} else {
		log.Printf("updateLayout failed:%+v\n", updateLayoutResp)
		return
	}

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().DoHLS(ctx, resourceId, sid)
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

func (s *Service) RunHLSAndMP4(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	impl := cloudrecording.NewAPI(c).V1().MixRecording()
	// acquire
	acquireResp, err := impl.Acquire().Do(ctx, s.cname, s.uid, &v1.AcquireMixRecodingClientRequest{})
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
	startResp, err := impl.Start().Do(ctx, resourceId, s.cname, s.uid, &v1.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			MaxIdleTime:  30,
			AudioProfile: 2,
			TranscodingConfig: &v1.TranscodingConfig{
				Width:            640,
				Height:           640,
				FPS:              15,
				BitRate:          800,
				MixedVideoLayout: 0,
				BackgroundColor:  "#000000",
			},
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeVideoUIDs: []string{
				"#allstream#",
			},
		},
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
				"mp4",
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

	sid := startResp.SuccessResp.Sid
	// stop
	defer func() {
		stopResp, err := impl.Stop().DoHLSAndMP4(ctx, resourceId, sid, s.cname, s.uid, false)
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
		queryResp, err := impl.Query().DoHLSAndMP4(ctx, resourceId, sid)
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
	updateResp, err := impl.Update().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateMixRecordingClientRequest{
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

	// updateLayout
	updateLayoutResp, err := impl.UpdateLayout().Do(ctx, resourceId, sid, s.cname, s.uid, &v1.UpdateLayoutUpdateMixRecordingClientRequest{
		MixedVideoLayout: 1,
		BackgroundColor:  "#FF0000",
	},
	)
	if err != nil {
		log.Println(err)
		return
	}
	if updateLayoutResp.IsSuccess() {
		log.Printf("updateLayout success:%+v\n", updateLayoutResp)
	} else {
		log.Printf("updateLayout failed:%+v\n", updateLayoutResp)
		return
	}

	// query
	for i := 0; i < 3; i++ {
		queryResp, err := impl.Query().DoHLSAndMP4(ctx, resourceId, sid)
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
