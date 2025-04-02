package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudtranscoder/service"
)

var (
	appId    string
	username string
	password string
	// Choose your region, currently supported regions are:
	// US: North America
	// EU: Europe
	// CN: China
	// AP: Asia Pacific
	domainArea = domain.CN
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

	s := service.New(domainArea, appId)
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
