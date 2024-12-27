package webrecording

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
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario/webrecording"
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

func (s *Service) RunWebRecorder(storageConfig *cloudRecordingAPI.StorageConfig) {
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
	acquireResp, err := cloudRecordingClient.WebRecording().Acquire(ctx, s.Cname, s.Uid, &webrecording.AcquireWebRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.WebRecording().Start(ctx, resourceId, s.Cname, s.Uid, &webrecording.StartWebRecordingClientRequest{
		RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
			AvFileType: []string{
				"hls",
				"mp4",
			},
		},
		StorageConfig: storageConfig,
		ExtensionServiceConfig: &cloudRecordingAPI.ExtensionServiceConfig{
			ErrorHandlePolicy: "error_abort",
			ExtensionServices: []cloudRecordingAPI.ExtensionService{
				{
					ServiceName:       "web_recorder_service",
					ErrorHandlePolicy: "error_abort",
					ServiceParam: &cloudRecordingAPI.WebRecordingServiceParam{
						URL:              "https://live.bilibili.com/",
						AudioProfile:     2,
						VideoWidth:       1280,
						VideoHeight:      720,
						MaxRecordingHour: 1,
					},
				},
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
		stopResp, err := cloudRecordingClient.WebRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, true)
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
		queryResp, err := cloudRecordingClient.WebRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.WebRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &webrecording.UpdateWebRecordingClientRequest{
		WebRecordingConfig: &cloudRecordingAPI.UpdateWebRecordingConfig{
			Onhold: false,
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
		queryResp, err := cloudRecordingClient.WebRecording().Query(ctx, resourceId, sid)
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

func (s *Service) RunWebRecorderAndRtmpPublish(storageConfig *cloudRecordingAPI.StorageConfig) {
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
	acquireResp, err := cloudRecordingClient.WebRecording().Acquire(ctx, s.Cname, s.Uid, &webrecording.AcquireWebRecodingClientRequest{})
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
	startResp, err := cloudRecordingClient.WebRecording().Start(ctx, resourceId, s.Cname, s.Uid, &webrecording.StartWebRecordingClientRequest{
		RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
			AvFileType: []string{
				"hls",
				"mp4",
			},
		},
		StorageConfig: storageConfig,
		ExtensionServiceConfig: &cloudRecordingAPI.ExtensionServiceConfig{
			ErrorHandlePolicy: "error_abort",
			ExtensionServices: []cloudRecordingAPI.ExtensionService{
				{
					ServiceName:       "web_recorder_service",
					ErrorHandlePolicy: "error_abort",
					ServiceParam: &cloudRecordingAPI.WebRecordingServiceParam{
						URL:              "https://live.bilibili.com/",
						AudioProfile:     2,
						VideoWidth:       1280,
						VideoHeight:      720,
						MaxRecordingHour: 1,
					},
				},
				{
					ServiceName:       "rtmp_publish_service",
					ErrorHandlePolicy: "error_ignore",
					ServiceParam: &cloudRecordingAPI.RtmpPublishServiceParam{
						Outputs: []cloudRecordingAPI.Outputs{
							{
								RtmpURL: "rtmp://xxx.xxx.xxx.xxx:1935/live/test",
							},
						},
					},
				},
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
		stopResp, err := cloudRecordingClient.WebRecording().Stop(ctx, resourceId, sid, s.Cname, s.Uid, true)
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
		queryResp, err := cloudRecordingClient.WebRecording().Query(ctx, resourceId, sid)
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
	updateResp, err := cloudRecordingClient.WebRecording().Update(ctx, resourceId, sid, s.Cname, s.Uid, &webrecording.UpdateWebRecordingClientRequest{
		WebRecordingConfig: &cloudRecordingAPI.UpdateWebRecordingConfig{
			Onhold: false,
		},
		RtmpPublishConfig: &cloudRecordingAPI.UpdateRtmpPublishConfig{
			Outputs: []cloudRecordingAPI.UpdateOutput{
				{
					RtmpURL: "rtmp://yyy.yyy.yyy.yyy:1935/live/test",
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
		queryResp, err := cloudRecordingClient.WebRecording().Query(ctx, resourceId, sid)
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
