package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudtranscoder/service"
)

var (
	appId    string
	username string
	password string
	// 选择你的区域，目前支持的区域有：
	// USRegionArea: 北美
	// EURegionArea: 欧洲
	// CNRegionArea: 中国大陆
	// APRegionArea: 亚太
	region core.RegionArea = core.CNRegionArea
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

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

	s := service.New(region, appId)
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
