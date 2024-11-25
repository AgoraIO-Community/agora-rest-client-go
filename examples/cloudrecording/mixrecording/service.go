package mixrecording

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	agoraLogger "github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/region"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/base"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	base2 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario/mixrecording"
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

func (s *Service) RunHLS(token string, storageConfig *base2.StorageConfig) {
	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		RegionCode: s.RegionArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	mixRecordingClient := cloudrecording.NewClient(config).MixRecording()
	// acquire
	acquireResp, err := mixRecordingClient.Acquire().Do(ctx, s.Cname, s.Uid, &mixrecording.AcquireMixRecodingClientRequest{})
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
	startResp, err := mixRecordingClient.Start().Do(ctx, resourceId, s.Cname, s.Uid, &mixrecording.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &base2.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			MaxIdleTime:  30,
			AudioProfile: 2,
			TranscodingConfig: &base2.TranscodingConfig{
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
		RecordingFileConfig: &base2.RecordingFileConfig{
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
		stopResp, err := mixRecordingClient.Stop().DoHLS(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := mixRecordingClient.Query().DoHLS(ctx, resourceId, sid)
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
	updateResp, err := mixRecordingClient.Update().Do(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateMixRecordingClientRequest{
		StreamSubscribe: &base2.UpdateStreamSubscribe{
			AudioUidList: &base2.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &base2.UpdateVideoUIDList{
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
	updateLayoutResp, err := mixRecordingClient.UpdateLayout().Do(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateLayoutUpdateMixRecordingClientRequest{
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
		queryResp, err := mixRecordingClient.Query().DoHLS(ctx, resourceId, sid)
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

func (s *Service) RunHLSAndMP4(token string, storageConfig *base2.StorageConfig) {
	ctx := context.Background()

	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		RegionCode: s.RegionArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	mixRecordingClient := cloudrecording.NewClient(config).MixRecording()

	// acquire
	acquireResp, err := mixRecordingClient.Acquire().Do(ctx, s.Cname, s.Uid, &mixrecording.AcquireMixRecodingClientRequest{})
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
	startResp, err := mixRecordingClient.Start().Do(ctx, resourceId, s.Cname, s.Uid, &mixrecording.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &base2.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			MaxIdleTime:  30,
			AudioProfile: 2,
			TranscodingConfig: &base2.TranscodingConfig{
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
		RecordingFileConfig: &base2.RecordingFileConfig{
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

	sid := startResp.SuccessResponse.Sid
	// stop
	defer func() {
		stopResp, err := mixRecordingClient.Stop().DoHLSAndMP4(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := mixRecordingClient.Query().DoHLSAndMP4(ctx, resourceId, sid)
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
	updateResp, err := mixRecordingClient.Update().Do(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateMixRecordingClientRequest{
		StreamSubscribe: &base2.UpdateStreamSubscribe{
			AudioUidList: &base2.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &base2.UpdateVideoUIDList{
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
	updateLayoutResp, err := mixRecordingClient.UpdateLayout().Do(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateLayoutUpdateMixRecordingClientRequest{
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
		queryResp, err := mixRecordingClient.Query().DoHLSAndMP4(ctx, resourceId, sid)
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
