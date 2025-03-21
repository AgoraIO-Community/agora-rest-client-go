package main

import (
	"flag"
	"log"
	"os"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/convoai/service"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai"
)

var (
	appId      string
	username   string
	password   string
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

	ttsVendor := flag.String("ttsVendor", "", "tts vendor, e.g. bytedance,microsoft,tencent,minimax,elevenLabs")
	serviceRegion := flag.Int("serviceRegion", 0, "service region, e.g. 1: ChineseMainland, 2: Global")
	flag.Parse()

	if *serviceRegion == 0 {
		log.Fatalln("serviceRegion should be specified")
	}

	s := service.New(domainArea, appId, convoai.ServiceRegion(*serviceRegion))
	s.SetCredential(username, password)

	switch *ttsVendor {
	case "bytedance":
		if *serviceRegion == 2 {
			log.Fatalln("Bytedance TTS is not supported in Global service region")
		}

		s.RunWithBytedanceTTS()
	case "microsoft":
		if *serviceRegion == 2 {
			log.Fatalln("Microsoft TTS is not supported in ChineseMainland service region")
		}

		s.RunWithMicrosoftTTS()
	case "tencent":
		if *serviceRegion == 2 {
			log.Fatalln("Tencent TTS is not supported in Global service region")
		}

		s.RunWithTencentTTS()
	case "minimax":
		if *serviceRegion == 2 {
			log.Fatalln("Minimax TTS is not supported in Global service region")
		}

		s.RunWithMinimaxTTS()
	case "elevenLabs":
		if *serviceRegion == 2 {
			log.Fatalln("ElevenLabs TTS is not supported in ChineseMainland service region")
		}

		s.RunWithElevenLabsTTS()
	default:
		log.Fatalln("Invalid tts vendor")
	}
}
