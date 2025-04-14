package mixrecording

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

func (s *Scenario) RunHLS(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()

	// acquire
	acquireResp, err := s.CloudRecordingClient.MixRecording().
		Acquire(ctx, s.Cname, s.Uid, &req.AcquireMixRecodingClientRequest{})
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
	startResp, err := s.CloudRecordingClient.MixRecording().Start(ctx, resourceId, s.Cname, s.Uid, &req.StartMixRecordingClientRequest{
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
		stopResp, err := s.CloudRecordingClient.MixRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := s.CloudRecordingClient.MixRecording().QueryHLS(ctx, resourceId, sid)
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
	updateResp, err := s.CloudRecordingClient.MixRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateMixRecordingClientRequest{
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
	updateLayoutResp, err := s.CloudRecordingClient.MixRecording().UpdateLayout(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateLayoutUpdateMixRecordingClientRequest{
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
		queryResp, err := s.CloudRecordingClient.MixRecording().QueryHLS(ctx, resourceId, sid)
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

func (s *Scenario) RunHLSAndMP4(token string, storageConfig *cloudRecordingAPI.StorageConfig) {
	ctx := context.Background()

	// acquire
	acquireResp, err := s.CloudRecordingClient.MixRecording().Acquire(ctx, s.Cname, s.Uid, &req.AcquireMixRecodingClientRequest{})
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
	startResp, err := s.CloudRecordingClient.MixRecording().Start(ctx, resourceId, s.Cname, s.Uid, &req.StartMixRecordingClientRequest{
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
		stopResp, err := s.CloudRecordingClient.MixRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, false)
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
		queryResp, err := s.CloudRecordingClient.MixRecording().QueryHLSAndMP4(ctx, resourceId, sid)
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
	updateResp, err := s.CloudRecordingClient.MixRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateMixRecordingClientRequest{
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
	updateLayoutResp, err := s.CloudRecordingClient.MixRecording().UpdateLayout(ctx, resourceId, sid, s.Cname, s.Uid, &req.UpdateLayoutUpdateMixRecordingClientRequest{
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
		queryResp, err := s.CloudRecordingClient.MixRecording().QueryHLSAndMP4(ctx, resourceId, sid)
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
