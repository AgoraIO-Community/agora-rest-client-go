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

func (s *Service) MixRecording(token string, storageConfig *v1.StorageConfig) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	mixRecordingV1 := cloudrecording.NewAPI(c).V1().MixRecording()
	resp, err := mixRecordingV1.Acquire().WithForwardRegion(core.CNForwardedReginPrefix).Do(ctx, s.cname, s.uid, &v1.AcquireMixRecodingClientRequest{})
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Printf("start resourceId:%s", resp.SuccessRes.ResourceId)
	} else {
		log.Printf("start resp:%+v", resp.ErrResponse)
	}

	startResp, err := mixRecordingV1.Start().Do(ctx, resp.SuccessRes.ResourceId, s.cname, s.uid, &v1.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			AudioProfile: 2,
			MaxIdleTime:  30,
			TranscodingConfig: &v1.TranscodingConfig{
				Width:            640,
				Height:           260,
				FPS:              15,
				BitRate:          500,
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
		log.Fatalln(err)
	}
	if startResp.IsSuccess() {
		log.Printf("success:%+v", &startResp.SuccessResp)
	} else {
		log.Printf("failed:%+v", &startResp.ErrResponse)
	}

	startSuccessResp := startResp.SuccessResp
	defer func() {
		stopResp, err := mixRecordingV1.Stop().DoHLSAndMP4(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, s.cname, s.uid, false)
		if err != nil {
			log.Fatalln(err)
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v", &stopResp.SuccessResp)
		} else {
			log.Fatalf("stop failed:%+v", &stopResp.ErrResponse)
		}
		log.Printf("stopServerResponse:%+v", stopResp.SuccessResp.ServerResponse)
	}()

	queryResp, err := mixRecordingV1.Query().DoHLSAndMP4(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid)
	if err != nil {
		log.Fatalln(err)
	}
	if queryResp.IsSuccess() {
		log.Printf("query success:%+v", queryResp.SuccessResp)
	} else {
		log.Printf("query failed:%+v", queryResp.ErrResponse)
		return
	}

	log.Printf("queryServerResponse:%+v", queryResp.SuccessResp.ServerResponse)

	time.Sleep(3 * time.Second)

	updateLayoutResp, err := mixRecordingV1.UpdateLayout().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, s.cname, s.uid, &v1.UpdateLayoutUpdateMixRecordingClientRequest{
		MixedVideoLayout: 3,
		BackgroundColor:  "#FF0000",
		LayoutConfig: []v1.UpdateLayoutConfig{
			{
				UID:        "22",
				XAxis:      0.1,
				YAxis:      0.1,
				Width:      0.1,
				Height:     0.1,
				Alpha:      1,
				RenderMode: 1,
			},
			{
				UID:        "2",
				XAxis:      0.2,
				YAxis:      0.2,
				Width:      0.1,
				Height:     0.1,
				Alpha:      1,
				RenderMode: 1,
			},
		},
	},
	)
	if err != nil {
		log.Fatalln(err)
	}
	if updateLayoutResp.IsSuccess() {
		log.Printf("updateLayout success:%+v", updateLayoutResp.SuccessResp)
	} else {
		log.Printf("updateLayout failed:%+v", updateLayoutResp.ErrResponse)
		return
	}
	time.Sleep(3 * time.Second)

	updateResp, err := mixRecordingV1.Update().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, s.cname, s.uid, &v1.UpdateMixRecordingClientRequest{
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
		log.Fatalln(err)
	}
	if updateResp.IsSuccess() {
		log.Printf("update success:%+v", updateResp.SuccessResp)
	} else {
		log.Printf("update failed:%+v", updateResp.ErrResponse)
		return
	}
	time.Sleep(2 * time.Second)
}
