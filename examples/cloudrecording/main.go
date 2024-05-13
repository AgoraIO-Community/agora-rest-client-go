package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

var (
	appId    string
	username string
	password string
	token    string
	cname    string
	uid      string
	// 选择你的区域，目前支持的区域有：
	// USRegionArea: 北美
	// EURegionArea: 欧洲
	// CNRegionArea: 中国大陆
	// APRegionArea: 亚太
	region core.RegionArea = core.CNRegionArea
)

// 选择你的存储配置 第三方云存储地区说明详情见 https://doc.shengwang.cn/api-ref/cloud-recording/restful/region-vendor
// 配置存储需要的参数
var storageConfig = &v1.StorageConfig{
	Vendor:    0,
	Region:    0,
	Bucket:    "",
	AccessKey: "",
	SecretKey: "",
	FileNamePrefix: []string{
		"",
	},
}

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
		panic("BASIC_PASSWORD is required")
	}
	token = os.Getenv("TOKEN")
	cname = os.Getenv("CNAME")
	if cname == "" {
		panic("CNAME is required")
	}
	uid = os.Getenv("UID")
	if uid == "" {
		panic("UID is required")
	}

	vendorStr := os.Getenv("STORAGE_CONFIG_VENDOR")
	storageVendor, err := strconv.Atoi(vendorStr)
	if err != nil {
		panic(err)
	}
	storageRegionStr := os.Getenv("STORAGE_CONFIG_REGION")
	storageRegion, err := strconv.Atoi(storageRegionStr)
	if err != nil {
		panic(err)
	}

	storageConfig.Vendor = storageVendor
	storageConfig.Region = storageRegion

	storageConfig.Bucket = os.Getenv("STORAGE_CONFIG_BUCKET")
	storageConfig.AccessKey = os.Getenv("STORAGE_CONFIG_ACCESS_KEY")
	storageConfig.SecretKey = os.Getenv("STORAGE_CONFIG_SECRET_KEY")

	mode := flag.Int("mode", 3, "1: mix, 2: individual, 3: webRecording")
	flag.Parse()

	switch *mode {
	case 1:
		MixRecording()
	case 2:
		IndividualRecording()
	case 3:
		WebRecording()
	default:
		panic("mode is required, 1: mix, 2: individual, 3: webRecording")
	}
}

// MixRecording hls&mp4
func MixRecording() {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      appId,
		Credential: core.NewBasicAuthCredential(username, password),
		RegionCode: region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	mixRecordingV1 := cloudrecording.NewAPI(c).V1().MixRecording()
	resp, err := mixRecordingV1.Acquire().WithForwardRegion(core.CNForwardedReginPrefix).Do(ctx, cname, uid, &v1.AcquireMixRecodingClientRequest{})
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Printf("start resourceId:%s", resp.SuccessRes.ResourceId)
	} else {
		log.Printf("start resp:%+v", resp.ErrResponse)
	}

	startResp, err := mixRecordingV1.Start().Do(ctx, resp.SuccessRes.ResourceId, cname, uid, &v1.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			AudioProfile: 2,
			MaxIdleTime:  30,
			TranscodingConfig: &v1.TranscodingConfig{
				Width:            640,
				Height:           260,
				FPS:              15,
				BitRate:          500,
				MixedVideoLayout: 0,
				BackgroundColor:  "#000000",
			},
			SubscribeAudioUIDs: []string{
				"#allstream#",
			},
			SubscribeVideoUIDs: []string{
				"#allstream#",
			},
		},
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
				"mp4",
			},
		},
		StorageConfig: storageConfig,
	})
	if err != nil {
		log.Fatalln(err)
	}
	if startResp.IsSuccess() {
		log.Printf("success:%+v", &startResp.SuccessResp)
	} else {
		log.Printf("failed:%+v", &startResp.ErrResponse)
	}

	startSuccessResp := startResp.SuccessResp
	defer func() {
		stopResp, err := mixRecordingV1.Stop().DoHLSAndMP4(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, cname, uid, false)
		if err != nil {
			log.Fatalln(err)
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v", &stopResp.SuccessResp)
		} else {
			log.Fatalf("stop failed:%+v", &stopResp.ErrResponse)
		}
		log.Printf("stopServerResponse:%+v", stopResp.SuccessResp.ServerResponse)
	}()

	queryResp, err := mixRecordingV1.Query().DoHLSAndMP4(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid)
	if err != nil {
		log.Fatalln(err)
	}
	if queryResp.IsSuccess() {
		log.Printf("query success:%+v", queryResp.SuccessResp)
	} else {
		log.Printf("query failed:%+v", queryResp.ErrResponse)
		return
	}

	log.Printf("queryServerResponse:%+v", queryResp.SuccessResp.ServerResponse)

	time.Sleep(3 * time.Second)

	updateLayoutResp, err := mixRecordingV1.UpdateLayout().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, cname, uid, &v1.UpdateLayoutUpdateMixRecordingClientRequest{
		MixedVideoLayout: 3,
		BackgroundColor:  "#FF0000",
		LayoutConfig: []v1.UpdateLayoutConfig{
			{
				UID:        "22",
				XAxis:      0.1,
				YAxis:      0.1,
				Width:      0.1,
				Height:     0.1,
				Alpha:      1,
				RenderMode: 1,
			},
			{
				UID:        "2",
				XAxis:      0.2,
				YAxis:      0.2,
				Width:      0.1,
				Height:     0.1,
				Alpha:      1,
				RenderMode: 1,
			},
		},
	},
	)
	if err != nil {
		log.Fatalln(err)
	}
	if updateLayoutResp.IsSuccess() {
		log.Printf("updateLayout success:%+v", updateLayoutResp.SuccessResp)
	} else {
		log.Printf("updateLayout failed:%+v", updateLayoutResp.ErrResponse)
		return
	}
	time.Sleep(3 * time.Second)

	updateResp, err := mixRecordingV1.Update().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.Sid, cname, uid, &v1.UpdateMixRecordingClientRequest{
		StreamSubscribe: &v1.UpdateStreamSubscribe{
			AudioUidList: &v1.UpdateAudioUIDList{
				SubscribeAudioUIDs: []string{
					"#allstream#",
				},
			},
			VideoUidList: &v1.UpdateVideoUIDList{
				SubscribeVideoUIDs: []string{
					"#allstream#",
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

// IndividualRecording hls
func IndividualRecording() {
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

// WebRecording hls&mp4
func WebRecording() {
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
					ServiceParam: &v1.ServiceParam{
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
