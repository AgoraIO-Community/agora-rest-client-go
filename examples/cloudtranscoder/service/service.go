package service

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	agoraLogger "github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudtranscoder"
	cloudTranscoderAPI "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudtranscoder/api"
)

type Service struct {
	domainArea domain.Area
	appId      string
	credential auth.Credential
}

func New(region domain.Area, appId string) *Service {
	return &Service{
		domainArea: region,
		appId:      appId,
		credential: nil,
	}
}

func (s *Service) SetCredential(username string, password string) {
	s.credential = auth.NewBasicAuthCredential(username, password)
}

func (s *Service) acquireResource(ctx context.Context, client *cloudtranscoder.Client, instanceId string, createBody *cloudTranscoderAPI.CreateReqServices) string {
	acquireResp, err := client.Acquire(ctx, &cloudTranscoderAPI.AcquireReqBody{
		InstanceId: instanceId,
		Services:   createBody,
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
	inputUID1 := os.Getenv("INPUT_UID_1")
	if inputUID1 == "" {
		panic("INPUT_UID_1 is required")
	}

	inputUID1Int, err := strconv.Atoi(inputUID1)
	if err != nil {
		panic(err)
	}

	inputUID2 := os.Getenv("INPUT_UID_2")
	if inputUID2 == "" {
		panic("INPUT_UID_2 is required")
	}

	inputUID2Int, err := strconv.Atoi(inputUID2)
	if err != nil {
		panic(err)
	}

	inputChannelName := os.Getenv("INPUT_CHANNEL_NAME")
	if inputChannelName == "" {
		panic("INPUT_CHANNEL_NAME is required")
	}

	inputToken1 := os.Getenv("INPUT_TOKEN_1")
	if inputToken1 == "" {
		panic("INPUT_TOKEN_1 is required")
	}

	inputToken2 := os.Getenv("INPUT_TOKEN_2")
	if inputToken2 == "" {
		panic("INPUT_TOKEN_2 is required")
	}

	outputChannelName := os.Getenv("OUTPUT_CHANNEL_NAME")
	if outputChannelName == "" {
		panic("OUTPUT_CHANNEL_NAME is required")
	}

	outputUID := os.Getenv("OUTPUT_UID")
	if outputUID == "" {
		panic("OUTPUT_UID is required")
	}

	outputUIDInt, err := strconv.Atoi(outputUID)
	if err != nil {
		panic(err)
	}

	outputToken := os.Getenv("OUTPUT_TOKEN")
	if outputToken == "" {
		panic("OUTPUT_TOKEN is required")
	}

	updateInputUID3 := os.Getenv("UPDATE_INPUT_UID_3")
	if updateInputUID3 == "" {
		panic("UPDATE_INPUT_UID_3 is required")
	}

	updateInputUID3Int, err := strconv.Atoi(updateInputUID3)
	if err != nil {
		panic(err)
	}

	updateInputToken3 := os.Getenv("UPDATE_INPUT_TOKEN_3")
	if updateInputToken3 == "" {
		panic("UPDATE_INPUT_TOKEN_3 is required")
	}

	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.appId,
		Credential: s.credential,
		DomainArea: s.domainArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	cloudTranscoderClient, err := cloudtranscoder.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	createBody := &cloudTranscoderAPI.CreateReqBody{
		Services: &cloudTranscoderAPI.CreateReqServices{
			CloudTranscoder: &cloudTranscoderAPI.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &cloudTranscoderAPI.CloudTranscoderConfig{
					Transcoder: &cloudTranscoderAPI.CloudTranscoderConfigPayload{
						IdleTimeout: 100,
						AudioInputs: []cloudTranscoderAPI.CloudTranscoderAudioInput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID1Int,
									RtcToken:   inputToken1,
								},
							},
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID2Int,
									RtcToken:   inputToken2,
								},
							},
						},
						VideoInputs: []cloudTranscoderAPI.CloudTranscoderVideoInput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID1Int,
									RtcToken:   inputToken1,
								},
								Region: &cloudTranscoderAPI.CloudTranscoderRegion{
									X:      0,
									Y:      0,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID2Int,
									RtcToken:   inputToken2,
								},
								Region: &cloudTranscoderAPI.CloudTranscoderRegion{
									X:      0,
									Y:      240,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
						},
						Canvas: &cloudTranscoderAPI.CloudTranscoderCanvas{
							Width:  1280,
							Height: 720,
							Color:  255,
						},
						Outputs: []cloudTranscoderAPI.CloudTranscoderOutput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: outputChannelName,
									RtcUID:     outputUIDInt,
									RtcToken:   outputToken,
								},
								AudioOption: &cloudTranscoderAPI.CloudTranscoderOutputAudioOption{
									ProfileType: "AUDIO_PROFILE_MUSIC_STANDARD",
								},
								VideoOption: &cloudTranscoderAPI.CloudTranscoderOutputVideoOption{
									FPS:     15,
									Codec:   "H264",
									Bitrate: 1500,
									Width:   1280,
									Height:  720,
								},
							},
						},
					},
				},
			},
		},
	}

	tokenName := s.acquireResource(ctx, cloudTranscoderClient, instanceId, createBody.Services)
	log.Printf("tokenName:%s\n", tokenName)

	createResp, err := cloudTranscoderClient.Create(context.Background(), tokenName, createBody)
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
		deleteResp, err := cloudTranscoderClient.Delete(ctx, taskId, tokenName)
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
		queryResp, err := cloudTranscoderClient.Query(ctx, taskId, tokenName)
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

	updateResp, err := cloudTranscoderClient.Update(ctx, taskId, tokenName, 1, "services.cloudTranscoder.config", &cloudTranscoderAPI.UpdateReqBody{
		Services: &cloudTranscoderAPI.CreateReqServices{
			CloudTranscoder: &cloudTranscoderAPI.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &cloudTranscoderAPI.CloudTranscoderConfig{
					Transcoder: &cloudTranscoderAPI.CloudTranscoderConfigPayload{
						IdleTimeout: 100,
						AudioInputs: []cloudTranscoderAPI.CloudTranscoderAudioInput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID1Int,
									RtcToken:   inputToken1,
								},
							},
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID2Int,
									RtcToken:   inputToken2,
								},
							},
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     updateInputUID3Int,
									RtcToken:   updateInputToken3,
								},
							},
						},
						VideoInputs: []cloudTranscoderAPI.CloudTranscoderVideoInput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID1Int,
									RtcToken:   inputToken1,
								},
								Region: &cloudTranscoderAPI.CloudTranscoderRegion{
									X:      0,
									Y:      0,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID2Int,
									RtcToken:   inputToken2,
								},
								Region: &cloudTranscoderAPI.CloudTranscoderRegion{
									X:      0,
									Y:      240,
									Width:  480,
									Height: 360,
									ZOrder: 2,
								},
							},
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     updateInputUID3Int,
									RtcToken:   updateInputToken3,
								},
								Region: &cloudTranscoderAPI.CloudTranscoderRegion{
									X:      240,
									Y:      240,
									Width:  240,
									Height: 240,
									ZOrder: 2,
								},
							},
						},
						Canvas: &cloudTranscoderAPI.CloudTranscoderCanvas{
							Width:  1280,
							Height: 720,
							Color:  255,
						},
						Outputs: []cloudTranscoderAPI.CloudTranscoderOutput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: outputChannelName,
									RtcUID:     outputUIDInt,
									RtcToken:   outputToken,
								},
								AudioOption: &cloudTranscoderAPI.CloudTranscoderOutputAudioOption{
									ProfileType: "AUDIO_PROFILE_MUSIC_STANDARD",
								},
								VideoOption: &cloudTranscoderAPI.CloudTranscoderOutputVideoOption{
									FPS:     15,
									Codec:   "H264",
									Bitrate: 1500,
									Width:   1280,
									Height:  720,
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
		queryResp, err := cloudTranscoderClient.Query(ctx, taskId, tokenName)
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
	inputUID1 := os.Getenv("INPUT_UID_1")
	if inputUID1 == "" {
		panic("INPUT_UID_1 is required")
	}

	inputUID1Int, err := strconv.Atoi(inputUID1)
	if err != nil {
		panic(err)
	}

	inputToken1 := os.Getenv("INPUT_TOKEN_1")
	if inputToken1 == "" {
		panic("INPUT_TOKEN_1 is required")
	}

	inputChannelName := os.Getenv("INPUT_CHANNEL_NAME")
	if inputChannelName == "" {
		panic("INPUT_CHANNEL_NAME is required")
	}

	outputChannelName := os.Getenv("OUTPUT_CHANNEL_NAME")
	if outputChannelName == "" {
		panic("OUTPUT_CHANNEL_NAME is required")
	}

	outputUID := os.Getenv("OUTPUT_UID")
	if outputUID == "" {
		panic("OUTPUT_UID is required")
	}

	outputUIDInt, err := strconv.Atoi(outputUID)
	if err != nil {
		panic(err)
	}

	outputToken := os.Getenv("OUTPUT_TOKEN")
	if outputToken == "" {
		panic("OUTPUT_TOKEN is required")
	}

	ctx := context.Background()
	config := &agora.Config{
		AppID:      s.appId,
		Credential: s.credential,
		DomainArea: s.domainArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	cloudTranscoderClient, err := cloudtranscoder.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	createBody := &cloudTranscoderAPI.CreateReqBody{
		Services: &cloudTranscoderAPI.CreateReqServices{
			CloudTranscoder: &cloudTranscoderAPI.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &cloudTranscoderAPI.CloudTranscoderConfig{
					Transcoder: &cloudTranscoderAPI.CloudTranscoderConfigPayload{
						IdleTimeout: 300,
						AudioInputs: []cloudTranscoderAPI.CloudTranscoderAudioInput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID1Int,
									RtcToken:   inputToken1,
								},
							},
						},
						Outputs: []cloudTranscoderAPI.CloudTranscoderOutput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: outputChannelName,
									RtcUID:     outputUIDInt,
									RtcToken:   outputToken,
								},
								AudioOption: &cloudTranscoderAPI.CloudTranscoderOutputAudioOption{
									ProfileType: "AUDIO_PROFILE_MUSIC_STANDARD",
								},
							},
						},
					},
				},
			},
		},
	}

	tokenName := s.acquireResource(ctx, cloudTranscoderClient, instanceId, createBody.Services)
	log.Printf("tokenName:%s\n", tokenName)

	createResp, err := cloudTranscoderClient.Create(ctx, tokenName, createBody)
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
		deleteResp, err := cloudTranscoderClient.Delete(ctx, taskId, tokenName)
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
		queryResp, err := cloudTranscoderClient.Query(ctx, taskId, tokenName)
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

	updateResp, err := cloudTranscoderClient.Update(ctx, taskId, tokenName, 1, "services.cloudTranscoder.config", &cloudTranscoderAPI.UpdateReqBody{
		Services: &cloudTranscoderAPI.CreateReqServices{
			CloudTranscoder: &cloudTranscoderAPI.CloudTranscoderPayload{
				ServiceType: "cloudTranscoderV2",
				Config: &cloudTranscoderAPI.CloudTranscoderConfig{
					Transcoder: &cloudTranscoderAPI.CloudTranscoderConfigPayload{
						IdleTimeout: 300,
						AudioInputs: []cloudTranscoderAPI.CloudTranscoderAudioInput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: inputChannelName,
									RtcUID:     inputUID1Int,
									RtcToken:   inputToken1,
								},
							},
						},
						Outputs: []cloudTranscoderAPI.CloudTranscoderOutput{
							{
								Rtc: &cloudTranscoderAPI.CloudTranscoderRtc{
									RtcChannel: outputChannelName,
									RtcUID:     outputUIDInt,
									RtcToken:   outputToken,
								},
								AudioOption: &cloudTranscoderAPI.CloudTranscoderOutputAudioOption{
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
		queryResp, err := cloudTranscoderClient.Query(ctx, taskId, tokenName)
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
