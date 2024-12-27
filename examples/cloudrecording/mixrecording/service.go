package mixrecording

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	agoraLogger "github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/base"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	cloudRecordingAPI "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario/mixrecording"
)

type Service struct {
	base.Service
}

func NewService(region domain.Area, appId, cname, uid string) *Service {
	return &Service{
		base.Service{
			DomainArea: region,
			AppId:      appId,
			Cname:      cname,
			Uid:        uid,
			Credential: nil,
		},
	}
}

func (s *Service) RunHLS(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		DomainArea: s.DomainArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	// acquire
	acquireResp, err := cloudRecordingClient.MixRecording().
		Acquire(ctx, s.Cname, s.Uid, &mixrecording.AcquireMixRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.MixRecording().Start(ctx, resourceId, s.Cname, s.Uid, &mixrecording.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &cloudRecordingAPI.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			MaxIdleTime:  30,
			AudioProfile: 2,
			TranscodingConfig: &cloudRecordingAPI.TranscodingConfig{
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
		stopResp, err := cloudRecordingClient.MixRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := cloudRecordingClient.MixRecording().QueryHLS(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.MixRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateMixRecordingClientRequest{
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

	// updateLayout
	updateLayoutResp, err := cloudRecordingClient.MixRecording().UpdateLayout(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateLayoutUpdateMixRecordingClientRequest{
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
		queryResp, err := cloudRecordingClient.MixRecording().QueryHLS(ctx, resourceId, sid)
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

func (s *Service) RunHLSAndMP4(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()

	config := &agora.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		DomainArea: s.DomainArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}
	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	// acquire
	acquireResp, err := cloudRecordingClient.MixRecording().Acquire(ctx, s.Cname, s.Uid, &mixrecording.AcquireMixRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.MixRecording().Start(ctx, resourceId, s.Cname, s.Uid, &mixrecording.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &cloudRecordingAPI.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			MaxIdleTime:  30,
			AudioProfile: 2,
			TranscodingConfig: &cloudRecordingAPI.TranscodingConfig{
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
		RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
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
		stopResp, err := cloudRecordingClient.MixRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := cloudRecordingClient.MixRecording().QueryHLSAndMP4(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.MixRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateMixRecordingClientRequest{
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

	// updateLayout
	updateLayoutResp, err := cloudRecordingClient.MixRecording().UpdateLayout(ctx, resourceId, sid, s.Cname, s.Uid, &mixrecording.UpdateLayoutUpdateMixRecordingClientRequest{
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
		queryResp, err := cloudRecordingClient.MixRecording().QueryHLSAndMP4(ctx, resourceId, sid)
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
