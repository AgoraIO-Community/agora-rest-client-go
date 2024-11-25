package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/region"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudtranscoder/service"
)

var (
	appId    string
	username string
	password string
	// 选择你的区域，目前支持的区域有：
	// USArea: 北美
	// EUArea: 欧洲
	// CNArea: 中国大陆
	// APArea: 亚太
	regionArea = region.CNArea
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	appId = os.Getenv("APP_ID")
	if appId == "" {
		panic("APP_ID is required")
	}

	username = os.Getenv("BASIC_AUTH_USERNAME")
	if username == "" {
		panic("BASIC_AUTH_USERNAME is required")
	}

	password = os.Getenv("BASIC_AUTH_PASSWORD")
	if password == "" {
		panic("BASIC_AUTH_PASSWORD is required")
	}

	scene := flag.String("scene", "", "scene name")
	instanceId := flag.String("instanceId", uuid.NewString(), "instanceId for cloudTransCoder service")
	flag.Parse()

	s := service.New(regionArea, appId)
	s.SetCredential(username, password)

	switch *scene {
	case "single_channel_rtc_pull_mixer_rtc_push":
		s.RunSingleChannelRtcPullMixerRtcPush(*instanceId)

	case "single_channel_rtc_pull_fullchannel_audiomixer_rtc_push":
		s.RunSingleChannelRtcPullFullChannelAudioMixerRtcPush(*instanceId)

	default:
		panic("invalid scene name")
	}
}
