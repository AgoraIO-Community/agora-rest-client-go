package webrecording

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

// WebRecording hls&mp4
func WebRecording(appId, username, password, token, cname, uid string, storageConfig *v1.StorageConfig, region core.RegionArea) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      appId,
		Credential: core.NewBasicAuthCredential(username, password),
		RegionCode: region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	webRecordingV1 := cloudrecording.NewAPI(c).V1().WebRecording()

	// acquire
	resp, err := webRecordingV1.Acquire().Do(ctx, cname, uid, &v1.AcquireWebRecodingClientRequest{})
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Printf("acquire success:%+v", resp.SuccessRes)
	} else {
		log.Fatalf("acquire failed:%+v", resp)
	}

	resourceId := resp.SuccessRes.ResourceId

	// start
	startResp, err := webRecordingV1.Start().Do(ctx, resourceId, cname, uid, &v1.StartWebRecordingClientRequest{
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
				"mp4",
			},
		},
		StorageConfig: storageConfig,
		ExtensionServiceConfig: &v1.ExtensionServiceConfig{
			ErrorHandlePolicy: "error_abort",
			ExtensionServices: []v1.ExtensionService{
				{
					ServiceName:       "web_recorder_service",
					ErrorHandlePolicy: "error_abort",
					ServiceParam: &v1.WebRecordingServiceParam{
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
		log.Fatal(err)
	}
	if startResp.IsSuccess() {
		log.Printf("startResp success:%+v", &startResp.SuccessResp)
	} else {
		log.Fatalf("startResp failed:%+v", &startResp.ErrResponse)
	}

	startSuccessResp := startResp.SuccessResp
	sid := startSuccessResp.Sid

	defer func() {
		// stop
		stopResp, err := webRecordingV1.Stop().Do(ctx, resourceId, sid, cname, uid, false)
		if err != nil {
			log.Fatalln(err)
		}
		if stopResp.IsSuccess() {
			log.Printf("stopResp success:%+v", &stopResp.SuccessResp)
		} else {
			log.Fatalf("stopResp failed:%+v", &stopResp.ErrResponse)
		}
		log.Printf("stopServerResponse:%+v", stopResp.SuccessResp.ServerResponse)
	}()

	// query
	queryResp, err := webRecordingV1.Query().Do(ctx, resourceId, sid)
	if err != nil {
		log.Fatalln(err)
	}
	if queryResp.IsSuccess() {
		log.Printf("queryResp success:%+v", queryResp.SuccessResp)
	} else {
		log.Fatalf("queryResp failed:%+v", queryResp.ErrResponse)
	}
	log.Printf("queryServerResponse:%+v", queryResp.SuccessResp.ServerResponse)

	time.Sleep(3 * time.Second)

	// update
	updateResp, err := webRecordingV1.Update().Do(ctx, resourceId, sid, cname, uid, &v1.UpdateWebRecordingClientRequest{
		WebRecordingConfig: &v1.UpdateWebRecordingConfig{
			Onhold: false,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if updateResp.IsSuccess() {
		log.Printf("updateResp success:%+v", updateResp.SuccessResp)
	} else {
		log.Fatalf("updateResp failed:%+v", updateResp.ErrResponse)
	}

	time.Sleep(3 * time.Second)
}
