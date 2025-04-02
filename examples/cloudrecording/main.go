package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/individualrecording"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/mixrecording"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/webrecording"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

var (
	appId    string
	cname    string
	uid      string
	username string
	password string
	token    string
	// Choose your region, currently supported regions are:
	// US: North America
	// EU: Europe
	// CN: China
	// AP: Asia Pacific
	domainArea = domain.CN
)

// Select your storage configuration. For more details on third-party cloud storage regions,
// Parameters needed for storage configuration
var storageConfig = &v1.StorageConfig{
	Vendor:    0,
	Region:    0,
	Bucket:    "",
	AccessKey: "",
	SecretKey: "",
	FileNamePrefix: []string{
		"recordings",
	},
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	appId = os.Getenv("APP_ID")
	if appId == "" {
		panic("APP_ID is required")
	}

	cname = os.Getenv("CNAME")
	if cname == "" {
		panic("CNAME is required")
	}

	uid = os.Getenv("USER_ID")
	if uid == "" {
		panic("USER_ID is required")
	}

	username = os.Getenv("BASIC_AUTH_USERNAME")
	if username == "" {
		panic("BASIC_AUTH_USERNAME is required")
	}

	password = os.Getenv("BASIC_AUTH_PASSWORD")
	if password == "" {
		panic("BASIC_AUTH_PASSWORD is required")
	}

	token = os.Getenv("TOKEN")

	storageVendorStr := os.Getenv("STORAGE_CONFIG_VENDOR")
	storageVendor, err := strconv.Atoi(storageVendorStr)
	if err != nil {
		panic(err)
	}
	storageConfig.Vendor = storageVendor

	storageRegionStr := os.Getenv("STORAGE_CONFIG_REGION")
	storageRegion, err := strconv.Atoi(storageRegionStr)
	if err != nil {
		panic(err)
	}
	storageConfig.Region = storageRegion

	storageConfig.Bucket = os.Getenv("STORAGE_CONFIG_BUCKET")
	if storageConfig.Bucket == "" {
		panic("STORAGE_CONFIG_BUCKET is required")
	}

	storageConfig.AccessKey = os.Getenv("STORAGE_CONFIG_ACCESS_KEY")
	if storageConfig.AccessKey == "" {
		panic("STORAGE_CONFIG_ACCESS_KEY is required")
	}

	storageConfig.SecretKey = os.Getenv("STORAGE_CONFIG_SECRET_KEY")
	if storageConfig.SecretKey == "" {
		panic("STORAGE_CONFIG_SECRET_KEY is required")
	}

	mode := flag.String("mode", "mix", "recording mode, options is mix/individual/web")
	mix_scene := flag.String("mix_scene", "hls", "scene for mix mode, options is hls/hls_and_mp4")
	individual_scene := flag.String("individual_scene", "recording", "scene for individual mode, options is recording/snapshot/recording_and_snapshot/recording_and_postpone_transcoding/recording_and_audio_mix")
	web_scene := flag.String("web_scene", "web_recorder", "scene for web mode, options is web_recorder/web_recorder_and_rtmp_publish")
	flag.Parse()

	switch *mode {
	case "mix":
		service := mixrecording.NewService(domainArea, appId, cname, uid)
		service.SetCredential(username, password)

		switch *mix_scene {
		case "hls":
			service.RunHLS(token, storageConfig)
		case "hls_and_mp4":
			service.RunHLSAndMP4(token, storageConfig)
		default:
			panic("invalid mix_scene")
		}
	case "individual":
		service := individualrecording.NewService(domainArea, appId, cname, uid)
		service.SetCredential(username, password)

		switch *individual_scene {
		case "recording":
			service.RunRecording(token, storageConfig)
		case "snapshot":
			service.RunSnapshot(token, storageConfig)
		case "recording_and_snapshot":
			service.RunRecordingAndSnapshot(token, storageConfig)
		case "recording_and_postpone_transcoding":
			service.RunRecordingAndPostponeTranscoding(token, storageConfig)
		case "recording_and_audio_mix":
			service.RunRecordingAndAudioMix(token, storageConfig)
		default:
			panic("invalid individual_scene")
		}
	case "web":
		service := webrecording.NewService(domainArea, appId, cname, uid)
		service.SetCredential(username, password)

		switch *web_scene {
		case "web_recorder":
			service.RunWebRecorder(storageConfig)
		case "web_recorder_and_rtmp_publish":
			service.RunWebRecorderAndRtmpPublish(storageConfig)
		default:
			panic("invalid web_scene")
		}
	default:
		panic("invalid mode")
	}
}
