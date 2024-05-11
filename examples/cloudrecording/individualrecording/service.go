package individualrecording

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

// IndividualRecording hls
func IndividualRecording(appId, username, password, token, cname, uid string, storageConfig *v1.StorageConfig, region core.RegionArea) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      appId,
		Credential: core.NewBasicAuthCredential(username, password),
		RegionCode: region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	individualRecordingV1 := cloudrecording.NewAPI(c).V1().IndividualRecording()

	resp, err := individualRecordingV1.Acquire().Do(ctx, cname, uid, false, &v1.AcquireIndividualRecodingClientRequest{})
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Printf("acquire success:%+v", resp.SuccessRes)
	} else {
		log.Fatalf("acquire failed:%+v", resp)
	}

	startResp, err := individualRecordingV1.Start().Do(ctx, resp.SuccessRes.ResourceId, cname, uid, &v1.StartIndividualRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType: 1,
			StreamTypes: 2,
			SubscribeAudioUIDs: []string{
				"22",
				"456",
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
		log.Fatal(err)
	}
	if startResp.IsSuccess() {
		log.Printf("startResp success:%+v", &startResp.SuccessResp)
	} else {
		log.Fatalf("startResp failed:%+v", &startResp.ErrResponse)
	}
	startSuccessResp := startResp.SuccessResp

	defer func() {
		stopResp, err := individualRecordingV1.Stop().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, cname, uid, false)
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
	queryResp, err := individualRecordingV1.Query().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid)
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
	updateResp, err := individualRecordingV1.Update().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, cname, uid, &v1.UpdateIndividualRecordingClientRequest{
		StreamSubscribe: &v1.UpdateStreamSubscribe{
			AudioUidList: &v1.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"999",
				},
			},
			VideoUidList: &v1.UpdateVideoUIDList{
				SubscribeVideoUIDs: []string{
					"999",
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
