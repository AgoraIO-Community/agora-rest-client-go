package service

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	agoraLogger "github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	agoraUtils "github.com/AgoraIO-Community/agora-rest-client-go/agora/utils"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/req"
)

type Service struct {
	domainArea    domain.Area
	appId         string
	credential    auth.Credential
	serviceRegion convoai.ServiceRegion
}

func New(region domain.Area, appId string, serviceRegion convoai.ServiceRegion) *Service {
	return &Service{
		domainArea:    region,
		appId:         appId,
		credential:    nil,
		serviceRegion: serviceRegion,
	}
}

func (s *Service) SetCredential(username string, password string) {
	s.credential = auth.NewBasicAuthCredential(username, password)
}

func (s *Service) RunWithCustomTTS(ttsParam req.TTSVendorParamsInterface) {
	ctx := context.Background()
	config := &convoai.Config{
		AppID:         s.appId,
		HttpTimeout:   20 * time.Second,
		Credential:    s.credential,
		DomainArea:    s.domainArea,
		Logger:        agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
		ServiceRegion: s.serviceRegion,
	}

	convoaiClient, err := convoai.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	token := os.Getenv("CONVOAI_TOKEN")

	updateToken := os.Getenv("CONVOAI_UPDATE_TOKEN")

	channel := os.Getenv("CONVOAI_CHANNEL")
	if channel == "" {
		log.Fatalln("CONVOAI_CHANNEL is required")
	}

	agentRtcUId := os.Getenv("CONVOAI_AGENT_RTC_UID")
	if agentRtcUId == "" {
		log.Fatalln("CONVOAI_AGENT_RTC_UID is required")
	}

	remoteRtcUIds := os.Getenv("CONVOAI_REMOTE_RTC_UIDS")
	if remoteRtcUIds == "" {
		log.Fatalln("CONVOAI_REMOTE_RTC_UIDS is required")
	}

	llmURL := os.Getenv("CONVOAI_LLM_URL")
	if llmURL == "" {
		log.Fatalln("CONVOAI_LLM_URL is required")
	}

	llmAPIKey := os.Getenv("CONVOAI_LLM_API_KEY")
	if llmAPIKey == "" {
		log.Fatalln("CONVOAI_LLM_API_KEY is required")
	}

	llmModel := os.Getenv("CONVOAI_LLM_MODEL")
	if llmModel == "" {
		log.Fatalln("CONVOAI_LLM_MODEL is required")
	}

	name := s.appId + ":" + channel

	joinResp, err := convoaiClient.Join(ctx, name, &req.JoinPropertiesReqBody{
		Token:           token,
		Channel:         channel,
		AgentRtcUId:     agentRtcUId,
		RemoteRtcUIds:   []string{remoteRtcUIds},
		EnableStringUId: agoraUtils.Ptr(false),
		IdleTimeout:     agoraUtils.Ptr(120),
		AdvancedFeatures: &req.JoinPropertiesAdvancedFeaturesBody{
			EnableAIVad: agoraUtils.Ptr(true),
		},
		LLM: &req.JoinPropertiesCustomLLMBody{
			Url:    llmURL,
			APIKey: llmAPIKey,
			SystemMessages: []map[string]any{
				{
					"role":    "system",
					"content": "You are a helpful chatbot。",
				},
			},
			Params: map[string]any{
				"model":      llmModel,
				"max_tokens": 1024,
				"username":   "Jack",
			},
			MaxHistory:      agoraUtils.Ptr(30),
			GreetingMessage: "Hello,how can I help you?",
		},
		TTS: &req.JoinPropertiesTTSBody{
			Vendor: ttsParam.GetVendorType(),
			Params: ttsParam,
		},
		Vad: &req.JoinPropertiesVadBody{
			InterruptDurationMs: agoraUtils.Ptr(160),
			PrefixPaddingMs:     agoraUtils.Ptr(300),
			SilenceDurationMs:   agoraUtils.Ptr(480),
			Threshold:           agoraUtils.Ptr(0.5),
		},
		Asr: &req.JoinPropertiesAsrBody{
			Language: "zh-CN",
		},
		Parameters: &req.Parameters{
			FixedParams: &req.FixedParams{
				SilenceConfig: &req.SilenceConfig{
					TimeoutMs: agoraUtils.Ptr(1200),
					Action:    agoraUtils.Ptr("speak"),
					Content:   agoraUtils.Ptr("Hello, how can I help you?"),
				},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	if joinResp.IsSuccess() {
		log.Printf("Join success:%+v", joinResp)
	} else {
		log.Printf("Join failed:%+v", joinResp)
		return
	}

	agentId := joinResp.SuccessResp.AgentId

	defer func() {
		leaveResp, err := convoaiClient.Leave(ctx, agentId)
		if err != nil {
			log.Fatalln(err)
		}

		if leaveResp.IsSuccess() {
			log.Printf("Leave success:%+v", leaveResp)
		} else {
			log.Printf("Leave failed:%+v", leaveResp)
			return
		}
	}()

	for i := 0; i < 3; i++ {
		queryResp, err := convoaiClient.Query(ctx, agentId)
		if err != nil {
			log.Fatalln(err)
			return
		}

		if queryResp.IsSuccess() {
			log.Printf("Query success:%+v", queryResp)
		} else {
			log.Printf("Query failed:%+v", queryResp)
			return
		}
		time.Sleep(time.Second * 3)
	}

	interruptResp, err := convoaiClient.Interrupt(ctx, agentId)
	if err != nil {
		log.Fatalln(err)
	}

	if interruptResp.IsSuccess() {
		log.Printf("Interrupt success:%+v", interruptResp)
	} else {
		log.Printf("Interrupt failed:%+v", interruptResp)
		return
	}
	time.Sleep(time.Second * 1)

	historyResp, err := convoaiClient.GetHistory(ctx, agentId)
	if err != nil {
		log.Fatalln(err)
	}

	if historyResp.IsSuccess() {
		log.Printf("History success:%+v", historyResp)
	} else {
		log.Printf("History failed:%+v", historyResp)
		return
	}

	time.Sleep(time.Second * 1)

	speakResp, err := convoaiClient.Speak(ctx, agentId, &req.SpeakBody{
		Text:          agoraUtils.Ptr("Hello, how can I help you?"),
		Priority:      agoraUtils.Ptr("INTERRUPT"),
		Interruptable: agoraUtils.Ptr(true),
	})
	if err != nil {
		log.Fatalln(err)
	}

	if speakResp.IsSuccess() {
		log.Printf("Speak success:%+v", speakResp)
	} else {
		log.Printf("Speak failed:%+v", speakResp)
		return
	}
	time.Sleep(time.Second * 1)

	updateResp, err := convoaiClient.Update(ctx, agentId, &req.UpdateReqBody{
		Token: updateToken,
		LLM: &req.UpdateLLMBody{
			SystemMessages: []map[string]any{
				{
					"role":    "system",
					"content": "You are a helpful chatbot, and you are a new assistant.",
				},
			},
			Params: map[string]any{
				"model":      llmModel,
				"max_tokens": 2048,
				"username":   "Tom",
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	if updateResp.IsSuccess() {
		log.Printf("Update success:%+v", updateResp)
	} else {
		log.Printf("Update failed:%+v", updateResp)
		return
	}

	for i := 0; i < 3; i++ {
		queryRes, err := convoaiClient.Query(ctx, agentId)
		if err != nil {
			log.Fatalln(err)
		}

		if queryRes.IsSuccess() {
			log.Printf("Query success:%+v", queryRes)
		} else {
			log.Printf("Query failed:%+v", queryRes)
			return
		}
		time.Sleep(time.Second * 3)
	}

	listResp, err := convoaiClient.List(ctx,
		req.WithState(2),
		req.WithLimit(10),
	)
	if err != nil {
		log.Fatalln(err)
	}

	if listResp.IsSuccess() {
		log.Printf("List success:%+v", listResp)
	} else {
		log.Printf("List failed:%+v", listResp)
		return
	}

	if listResp.SuccessRes.Data.Count > 0 {
		for _, agent := range listResp.SuccessRes.Data.List {
			log.Printf("Agent:%+v\n", agent)
		}
	} else {
		log.Printf("No agent found\n")
	}
}

func (s *Service) RunWithBytedanceTTS() {
	ttsToken := os.Getenv("CONVOAI_TTS_BYTEDANCE_TOKEN")
	if ttsToken == "" {
		log.Fatalln("CONVOAI_TTS_TOKEN is required")
	}

	ttsAppId := os.Getenv("CONVOAI_TTS_BYTEDANCE_APP_ID")
	if ttsAppId == "" {
		log.Fatalln("CONVOAI_TTS_APP_ID is required")
	}

	ttsCluster := os.Getenv("CONVOAI_TTS_BYTEDANCE_CLUSTER")
	if ttsCluster == "" {
		log.Fatalln("CONVOAI_TTS_CLUSTER is required")
	}

	ttsVoiceType := os.Getenv("CONVOAI_TTS_BYTEDANCE_VOICE_TYPE")
	if ttsVoiceType == "" {
		log.Fatalln("CONVOAI_TTS_BYTEDANCE_VOICE_TYPE is required")
	}

	ttsParam := req.TTSBytedanceVendorParams{
		Token:       ttsToken,
		AppId:       ttsAppId,
		Cluster:     ttsCluster,
		VoiceType:   ttsVoiceType,
		SpeedRatio:  1.0,
		VolumeRatio: 1.0,
		PitchRatio:  1.0,
		Emotion:     "happy",
	}
	s.RunWithCustomTTS(ttsParam)
}

func (s *Service) RunWithTencentTTS() {
	ttsAppId := os.Getenv("CONVOAI_TTS_TENCENT_APP_ID")
	if ttsAppId == "" {
		log.Fatalln("CONVOAI_TTS_TENCENT_APP_ID is required")
	}

	ttsSecretId := os.Getenv("CONVOAI_TTS_TENCENT_SECRET_ID")
	if ttsSecretId == "" {
		log.Fatalln("CONVOAI_TTS_TENCENT_SECRET_ID is required")
	}

	ttsSecretKey := os.Getenv("CONVOAI_TTS_TENCENT_SECRET_KEY")
	if ttsSecretKey == "" {
		log.Fatalln("CONVOAI_TTS_TENCENT_SECRET_KEY is required")
	}

	ttsParam := req.TTSTencentVendorParams{
		AppId:            ttsAppId,
		SecretId:         ttsSecretId,
		SecretKey:        ttsSecretKey,
		VoiceType:        601005,
		Volume:           0,
		Speed:            0,
		EmotionCategory:  "happy",
		EmotionIntensity: 100,
	}
	s.RunWithCustomTTS(ttsParam)
}

func (s *Service) RunWithMinimaxTTS() {
	ttsGroupId := os.Getenv("CONVOAI_TTS_MINIMAX_GROUP_ID")
	if ttsGroupId == "" {
		log.Fatalln("CONVOAI_TTS_MINIMAX_GROUP_ID is required")
	}

	ttsGroupKey := os.Getenv("CONVOAI_TTS_MINIMAX_GROUP_KEY")
	if ttsGroupKey == "" {
		log.Fatalln("CONVOAI_TTS_MINIMAX_GROUP_KEY is required")
	}

	ttsGroupModel := os.Getenv("CONVOAI_TTS_MINIMAX_GROUP_MODEL")
	if ttsGroupModel == "" {
		log.Fatalln("CONVOAI_TTS_MINIMAX_GROUP_MODEL is required")
	}

	ttsParam := req.TTSMinimaxVendorParams{
		GroupId: ttsGroupId,
		Key:     ttsGroupKey,
		Model:   ttsGroupModel,
		VoiceSetting: &req.TTSMinimaxVendorVoiceSettingParam{
			VoiceId: "female-shaonv",
			Speed:   1,
			Vol:     1,
			Pitch:   0,
			Emotion: "happy",
		},
		AudioSetting: &req.TTSMinimaxVendorAudioSettingParam{
			SampleRate: 16000,
		},
	}

	s.RunWithCustomTTS(ttsParam)
}

func (s *Service) RunWithMicrosoftTTS() {
	ttsKey := os.Getenv("CONVOAI_TTS_MICROSOFT_KEY")
	if ttsKey == "" {
		log.Fatalln("CONVOAI_TTS_MICROSOFT_KEY is required")
	}

	ttsRegion := os.Getenv("CONVOAI_TTS_MICROSOFT_REGION")
	if ttsRegion == "" {
		log.Fatalln("CONVOAI_TTS_MICROSOFT_REGION is required")
	}

	ttsVoiceName := os.Getenv("CONVOAI_TTS_MICROSOFT_VOICE_NAME")
	if ttsVoiceName == "" {
		log.Fatalln("CONVOAI_TTS_MICROSOFT_VOICE_NAME is required")
	}

	ttsParam := req.TTSMicrosoftVendorParams{
		Key:        ttsKey,
		Region:     ttsRegion,
		VoiceName:  ttsVoiceName,
		Speed:      1.0,
		Volume:     70,
		SampleRate: 24000,
	}

	s.RunWithCustomTTS(ttsParam)
}

func (s *Service) RunWithElevenLabsTTS() {
	ttsApiKey := os.Getenv("CONVOAI_TTS_ELEVENLABS_API_KEY")
	if ttsApiKey == "" {
		log.Fatalln("CONVOAI_TTS_ELEVENLABS_API_KEY is required")
	}

	ttsModelId := os.Getenv("CONVOAI_TTS_ELEVENLABS_MODEL_ID")
	if ttsModelId == "" {
		log.Fatalln("CONVOAI_TTS_ELEVENLABS_MODEL_ID is required")
	}

	ttsVoiceId := os.Getenv("CONVOAI_TTS_ELEVENLABS_VOICE_ID")
	if ttsVoiceId == "" {
		log.Fatalln("CONVOAI_TTS_ELEVENLABS_VOICE_ID is required")
	}

	ttsParam := req.TTSElevenLabsVendorParams{
		Key:     ttsApiKey,
		ModelId: ttsModelId,
		VoiceId: ttsVoiceId,
	}

	s.RunWithCustomTTS(ttsParam)
}

func (s *Service) RunWithCartesiaTTS() {
	ttsApiKey := os.Getenv("CONVOAI_TTS_CARTESIA_API_KEY")
	if ttsApiKey == "" {
		log.Fatalln("CONVOAI_TTS_CARTESIA_API_KEY is required")
	}

	ttsModelId := os.Getenv("CONVOAI_TTS_CARTESIA_MODEL_ID")
	if ttsModelId == "" {
		log.Fatalln("CONVOAI_TTS_CARTESIA_MODEL_ID is required")
	}

	ttsVoiceMode := os.Getenv("CONVOAI_TTS_CARTESIA_VOICE_MODE")
	if ttsVoiceMode == "" {
		log.Fatalln("CONVOAI_TTS_CARTESIA_VOICE_MODE is required")
	}

	ttsVoiceId := os.Getenv("CONVOAI_TTS_CARTESIA_VOICE_ID")
	if ttsVoiceId == "" {
		log.Fatalln("CONVOAI_TTS_CARTESIA_VOICE_ID is required")
	}

	ttsParam := req.TTSCartesiaVendorParams{
		APIKey:  ttsApiKey,
		ModelId: ttsModelId,
		Voice: &req.TTSCartesiaVendorVoice{
			Mode: ttsVoiceMode,
			Id:   ttsVoiceId,
		},
	}

	s.RunWithCustomTTS(ttsParam)
}

func (s *Service) RunWithOpenAITTS() {
	ttsApiKey := os.Getenv("CONVOAI_TTS_OPENAI_API_KEY")
	if ttsApiKey == "" {
		log.Fatalln("CONVOAI_TTS_OPENAI_API_KEY is required")
	}

	ttsModel := os.Getenv("CONVOAI_TTS_OPENAI_MODEL")
	if ttsModel == "" {
		log.Fatalln("CONVOAI_TTS_OPENAI_MODEL is required")
	}

	ttsVoice := os.Getenv("CONVOAI_TTS_OPENAI_VOICE")
	if ttsVoice == "" {
		log.Fatalln("CONVOAI_TTS_OPENAI_VOICE is required")
	}

	ttsInstructions := os.Getenv("CONVOAI_TTS_OPENAI_INSTRUCTIONS")
	if ttsInstructions == "" {
		log.Fatalln("CONVOAI_TTS_OPENAI_INSTRUCTIONS is required")
	}

	ttsSpeed := os.Getenv("CONVOAI_TTS_OPENAI_SPEED")
	if ttsSpeed == "" {
		log.Fatalln("CONVOAI_TTS_OPENAI_SPEED is required")
	}

	speed, err := strconv.ParseFloat(ttsSpeed, 32)
	if err != nil {
		log.Fatalln("CONVOAI_TTS_OPENAI_SPEED is not a valid float")
	}

	ttsParam := req.TTSOpenAIVendorParams{
		APIKey:       ttsApiKey,
		Model:        ttsModel,
		Voice:        ttsVoice,
		Instructions: ttsInstructions,
		Speed:        float32(speed),
	}

	s.RunWithCustomTTS(ttsParam)
}
