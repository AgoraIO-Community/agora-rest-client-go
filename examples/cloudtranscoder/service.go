package main

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudtranscoder"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudtranscoder/v1"
)

type Service struct {
	region     core.RegionArea
	appId      string
	credential core.Credential
}

func NewService(region core.RegionArea) *Service {
	return &Service{
		region:     region,
		appId:      appId,
		credential: nil,
	}
}

func (s *Service) SetCredential(credential core.Credential) {
	s.credential = credential
}

func (s *Service) acquireResource(ctx context.Context, v1Impl *v1.BaseCollection, instanceId string) string {
	acquireResp, err := v1Impl.Acquire().Do(ctx, &v1.AcquireReqBody{
		InstanceId: instanceId,
	})
	if err != nil {
		log.Fatalln(err)
	}
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Fatalf("acquire failed:%+v\n", acquireResp)
	}

	tokenName := acquireResp.SuccessResp.TokenName

	if tokenName == "" {
		log.Fatalln("tokenName is empty")
	}
	return tokenName
}

func (s *Service) RunSingleChannelRtcPullMixerRtcPush(instanceId string) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	v1Impl := cloudtranscoder.NewAPI(c).V1()

	tokenName := s.acquireResource(ctx, v1Impl, instanceId)
	log.Printf("tokenName:%s\n", tokenName)

	createResp, err := v1Impl.Create().Do(context.Background(), tokenName, &v1.CreateReqBody{
		Services: v1.CreateReqServices{
			CloudTranscoder: &v1.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &v1.CloudTranscoderConfig{
					Transcoder: &v1.CloudTranscoderConfigPayload{
						IdleTimeout: 100,
						AudioInputs: []v1.CloudTranscoderAudioInput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-abc",
									RtcUID:     0,
									RtcToken:   "",
								},
							},
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-abc",
									RtcUID:     18,
									RtcToken:   "",
								},
							},
						},
						VideoInputs: []v1.CloudTranscoderVideoInput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "",
									RtcUID:     0,
									RtcToken:   "",
								},
								Region: &v1.CloudTranscoderRegion{
									X:      0,
									Y:      0,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "",
									RtcUID:     0,
									RtcToken:   "",
								},
								Region: &v1.CloudTranscoderRegion{
									X:      0,
									Y:      240,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
						},
						Canvas: &v1.CloudTranscoderCanvas{
							Width:  1280,
							Height: 720,
							Color:  255,
						},
						Outputs: []v1.CloudTranscoderOutput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-efg",
									RtcUID:     0,
									RtcToken:   "",
								},
								AudioOption: &v1.CloudTranscoderOutputAudioOption{
									ProfileType: "AUDIO_PROFILE_MUSIC_STANDARD",
								},
								VideoOption: &v1.CloudTranscoderOutputVideoOption{
									FPS:                   15,
									Codec:                 "H264",
									Bitrate:               1500,
									Width:                 1280,
									Height:                720,
									LowBitrateHighQuality: false,
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	if createResp.IsSuccess() {
		log.Printf("create success:%+v\n", createResp)
	} else {
		log.Printf("create failed:%+v\n", createResp)
		return
	}

	taskId := createResp.SuccessResp.TaskID

	defer func() {
		deleteResp, err := v1Impl.Delete().Do(ctx, taskId, tokenName)
		if err != nil {
			log.Println(err)
			return
		}
		if deleteResp.IsSuccess() {
			log.Printf("delete success:%+v\n", deleteResp)
		} else {
			log.Printf("delete failed:%+v\n", deleteResp)
			return
		}
	}()

	for i := 0; i < 3; i++ {
		queryResp, err := v1Impl.Query().Do(ctx, taskId, tokenName)
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

	updateResp, err := v1Impl.Update().Do(ctx, taskId, taskId, 1, &v1.UpdateReqBody{
		Services: v1.CreateReqServices{
			CloudTranscoder: &v1.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &v1.CloudTranscoderConfig{
					Transcoder: &v1.CloudTranscoderConfigPayload{
						IdleTimeout: 100,
						AudioInputs: []v1.CloudTranscoderAudioInput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-abc",
									RtcUID:     0,
									RtcToken:   "",
								},
							},
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-abc",
									RtcUID:     18,
									RtcToken:   "",
								},
							},
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "",
									RtcUID:     0,
									RtcToken:   "",
								},
							},
						},
						VideoInputs: []v1.CloudTranscoderVideoInput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "",
									RtcUID:     0,
									RtcToken:   "",
								},
								Region: &v1.CloudTranscoderRegion{
									X:      0,
									Y:      0,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "",
									RtcUID:     0,
									RtcToken:   "",
								},
								Region: &v1.CloudTranscoderRegion{
									X:      0,
									Y:      240,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
						},
						Canvas: &v1.CloudTranscoderCanvas{
							Width:  1280,
							Height: 720,
							Color:  255,
						},
						Outputs: []v1.CloudTranscoderOutput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-efg",
									RtcUID:     0,
									RtcToken:   "",
								},
								AudioOption: &v1.CloudTranscoderOutputAudioOption{
									ProfileType: "AUDIO_PROFILE_MUSIC_STANDARD",
								},
								VideoOption: &v1.CloudTranscoderOutputVideoOption{
									FPS:                   15,
									Codec:                 "H264",
									Bitrate:               1500,
									Width:                 1280,
									Height:                720,
									LowBitrateHighQuality: false,
								},
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

	if updateResp.IsSuccess() {
		log.Printf("update success:%+v\n", updateResp)
	} else {
		log.Printf("update failed:%+v\n", updateResp)
		return
	}

	for i := 0; i < 3; i++ {
		queryResp, err := v1Impl.Query().Do(ctx, taskId, tokenName)
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

func (s *Service) RunSingleChannelRtcPullFullChannelAudioMixerRtcPush(instanceId string) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	v1Impl := cloudtranscoder.NewAPI(c).V1()

	tokenName := s.acquireResource(ctx, v1Impl, instanceId)
	log.Printf("tokenName:%s\n", tokenName)

	createResp, err := v1Impl.Create().Do(ctx, tokenName, &v1.CreateReqBody{
		Services: v1.CreateReqServices{
			CloudTranscoder: &v1.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &v1.CloudTranscoderConfig{
					Transcoder: &v1.CloudTranscoderConfigPayload{
						IdleTimeout: 300,
						AudioInputs: []v1.CloudTranscoderAudioInput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-abc",
									RtcUID:     0,
									RtcToken:   "",
								},
							},
						},
						Outputs: []v1.CloudTranscoderOutput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-efg",
									RtcUID:     888,
									RtcToken:   "",
								},
								AudioOption: &v1.CloudTranscoderOutputAudioOption{
									ProfileType: "AUDIO_PROFILE_MUSIC_HIGH_QUALITY_STEREO",
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	if createResp.IsSuccess() {
		log.Printf("create success:%+v\n", createResp)
	} else {
		log.Printf("create failed:%+v\n", createResp)
		return
	}

	taskId := createResp.SuccessResp.TaskID

	defer func() {
		deleteResp, err := v1Impl.Delete().Do(ctx, taskId, tokenName)
		if err != nil {
			log.Println(err)
			return
		}
		if deleteResp.IsSuccess() {
			log.Printf("delete success:%+v\n", deleteResp)
		} else {
			log.Printf("delete failed:%+v\n", deleteResp)
			return
		}
	}()

	for i := 0; i < 3; i++ {
		queryResp, err := v1Impl.Query().Do(ctx, taskId, tokenName)
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

	updateResp, err := v1Impl.Update().Do(ctx, taskId, taskId, 1, &v1.UpdateReqBody{
		Services: v1.CreateReqServices{
			CloudTranscoder: &v1.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &v1.CloudTranscoderConfig{
					Transcoder: &v1.CloudTranscoderConfigPayload{
						IdleTimeout: 300,
						AudioInputs: []v1.CloudTranscoderAudioInput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-abc",
									RtcUID:     0,
									RtcToken:   "",
								},
							},
						},
						Outputs: []v1.CloudTranscoderOutput{
							{
								Rtc: &v1.CloudTranscoderRtc{
									RtcChannel: "test-efg",
									RtcUID:     888,
									RtcToken:   "",
								},
								AudioOption: &v1.CloudTranscoderOutputAudioOption{
									ProfileType: "AUDIO_PROFILE_MUSIC_HIGH_QUALITY_STEREO",
								},
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

	if updateResp.IsSuccess() {
		log.Printf("update success:%+v\n", updateResp)
	} else {
		log.Printf("update failed:%+v\n", updateResp)
		return
	}

	for i := 0; i < 3; i++ {
		queryResp, err := v1Impl.Query().Do(ctx, taskId, tokenName)
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
