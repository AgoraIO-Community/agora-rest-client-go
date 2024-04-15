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
	// US: 北美
	// EU: 欧洲
	// CN: 中国大陆
	// AP: 亚太
	region core.RegionArea = core.CN
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
	mode := "mix"

	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      appId,
		Credential: core.NewBasicAuthCredential(username, password),
		RegionCode: region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	cloudRecordingAPI := cloudrecording.NewAPI(c)

	resp, err := cloudRecordingAPI.V1().Acquire().Do(ctx, &v1.AcquirerReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &v1.AcquirerClientRequest{
			Scene:               0,
			ResourceExpiredHour: 24,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Printf("start resourceId:%s", resp.SuccessRes.ResourceId)
	} else {
		log.Printf("start resp:%+v", resp.ErrResponse)
	}

	starterResp, err := cloudRecordingAPI.V1().Start().Do(ctx, resp.SuccessRes.ResourceId, mode, &v1.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &v1.StartClientRequest{
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
				},
			},
			StorageConfig: storageConfig,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if starterResp.IsSuccess() {
		log.Printf("success:%+v", &starterResp.SuccessResp)
	} else {
		log.Printf("failed:%+v", &starterResp.ErrResponse)
	}

	startSuccessResp := starterResp.SuccessResp
	defer func() {
		stopResp, err := cloudRecordingAPI.V1().Stop().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.SID, mode, &v1.StopReqBody{
			Cname: cname,
			Uid:   uid,
			ClientRequest: &v1.StopClientRequest{
				AsyncStop: true,
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
		if stopResp.IsSuccess() {
			log.Printf("stop success:%+v", &stopResp.SuccessResp)
		} else {
			log.Fatalf("stop failed:%+v", &stopResp.ErrResponse)
		}
		stopSuccess := stopResp.SuccessResp
		var stopServerResponse interface{}
		switch stopSuccess.GetServerResponseMode() {
		case v1.StopServerResponseUnknownMode:
			log.Fatalln("unknown mode")
		case v1.StopIndividualRecordingServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopIndividualRecordingServerResponseMode)
			stopServerResponse = stopSuccess.GetIndividualRecordingServerResponse()
		case v1.StopIndividualVideoScreenshotServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopIndividualVideoScreenshotServerResponseMode)
			stopServerResponse = stopSuccess.GetIndividualVideoScreenshotServerResponse()
		case v1.StopMixRecordingHlsServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopMixRecordingHlsServerResponseMode)
			stopServerResponse = stopSuccess.GetMixRecordingHLSServerResponse()
		case v1.StopMixRecordingHlsAndMp4ServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopMixRecordingHlsAndMp4ServerResponseMode)
			stopServerResponse = stopSuccess.GetMixRecordingHLSAndMP4ServerResponse()
		case v1.StopWebRecordingServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopWebRecordingServerResponseMode)
			stopServerResponse = stopSuccess.GetWebRecordingServerResponse()
		}
		log.Printf("stopServerResponse:%+v", stopServerResponse)
	}()

	queryResp, err := cloudRecordingAPI.V1().Query().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.SID, mode)
	if err != nil {
		log.Fatalln(err)
	}
	if queryResp.IsSuccess() {
		log.Printf("query success:%+v", queryResp.SuccessResp)
	} else {
		log.Printf("query failed:%+v", queryResp.ErrResponse)
		return
	}

	var queryServerResponse interface{}

	querySuccess := queryResp.SuccessResp
	switch querySuccess.GetServerResponseMode() {
	case v1.QueryServerResponseUnknownMode:
		log.Fatalln("unknown mode")
	case v1.QueryIndividualRecordingServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryIndividualRecordingServerResponseMode)
		queryServerResponse = querySuccess.GetIndividualRecordingServerResponse()
	case v1.QueryIndividualVideoScreenshotServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryIndividualVideoScreenshotServerResponseMode)
		queryServerResponse = querySuccess.GetIndividualVideoScreenshotServerResponse()
	case v1.QueryMixRecordingHlsServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryMixRecordingHlsServerResponseMode)
		queryServerResponse = querySuccess.GetMixRecordingHLSServerResponse()
	case v1.QueryMixRecordingHlsAndMp4ServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryMixRecordingHlsAndMp4ServerResponseMode)
		queryServerResponse = querySuccess.GetMixRecordingHLSAndMP4ServerResponse()
	case v1.QueryWebRecordingServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryWebRecordingServerResponseMode)
		queryServerResponse = querySuccess.GetWebRecording2CDNServerResponse()
	}

	log.Printf("queryServerResponse:%+v", queryServerResponse)

	time.Sleep(3 * time.Second)

	updateLayoutResp, err := cloudRecordingAPI.V1().UpdateLayout().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.SID, mode, &v1.UpdateLayoutReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &v1.UpdateLayoutClientRequest{
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
	})
	if err != nil {
		log.Fatalln(err)
	}
	if updateLayoutResp.IsSuccess() {
		log.Printf("updateLayout success:%+v", updateLayoutResp.SuccessResp)
	} else {
		log.Printf("updateLayout failed:%+v", updateLayoutResp.ErrResponse)
		return
	}

	updateResp, err := cloudRecordingAPI.V1().Update().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.SID, mode, &v1.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &v1.UpdateClientRequest{
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
	mode := "individual"

	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      appId,
		Credential: core.NewBasicAuthCredential(username, password),
		RegionCode: region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	cloudRecordingAPI := cloudrecording.NewAPI(c)

	resp, err := cloudRecordingAPI.V1().Acquire().Do(ctx, &v1.AcquirerReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &v1.AcquirerClientRequest{
			Scene:               0,
			ResourceExpiredHour: 24,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Printf("acquire success:%+v", resp.SuccessRes)
	} else {
		log.Fatalf("acquire failed:%+v", resp)
	}

	starterResp, err := cloudRecordingAPI.V1().Start().Do(ctx, resp.SuccessRes.ResourceId, mode, &v1.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &v1.StartClientRequest{
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
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if starterResp.IsSuccess() {
		log.Printf("starterResp success:%+v", &starterResp.SuccessResp)
	} else {
		log.Fatalf("starterResp failed:%+v", &starterResp.ErrResponse)
	}
	startSuccessResp := starterResp.SuccessResp

	defer func() {
		stopResp, err := cloudRecordingAPI.V1().Stop().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.SID, mode, &v1.StopReqBody{
			Cname: cname,
			Uid:   uid,
			ClientRequest: &v1.StopClientRequest{
				AsyncStop: true,
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
		if stopResp.IsSuccess() {
			log.Printf("stopResp success:%+v", &stopResp.SuccessResp)
		} else {
			log.Fatalf("stopResp failed:%+v", &stopResp.ErrResponse)
		}
		stopSuccess := stopResp.SuccessResp
		var stopServerResponse interface{}
		switch stopSuccess.GetServerResponseMode() {
		case v1.StopServerResponseUnknownMode:
			log.Fatalln("unknown mode")
		case v1.StopIndividualRecordingServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopIndividualRecordingServerResponseMode)
			stopServerResponse = stopSuccess.GetIndividualRecordingServerResponse()
		case v1.StopIndividualVideoScreenshotServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopIndividualVideoScreenshotServerResponseMode)
			stopServerResponse = stopSuccess.GetIndividualVideoScreenshotServerResponse()
		case v1.StopMixRecordingHlsServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopMixRecordingHlsServerResponseMode)
			stopServerResponse = stopSuccess.GetMixRecordingHLSServerResponse()
		case v1.StopMixRecordingHlsAndMp4ServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopMixRecordingHlsAndMp4ServerResponseMode)
			stopServerResponse = stopSuccess.GetMixRecordingHLSAndMP4ServerResponse()
		case v1.StopWebRecordingServerResponseMode:
			log.Printf("serverResponseMode:%d", v1.StopWebRecordingServerResponseMode)
			stopServerResponse = stopSuccess.GetWebRecordingServerResponse()
		}
		log.Printf("stopServerResponse:%+v", stopServerResponse)
	}()
	queryResp, err := cloudRecordingAPI.V1().Query().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.SID, mode)
	if err != nil {
		log.Fatalln(err)
	}
	if queryResp.IsSuccess() {
		log.Printf("queryResp success:%+v", queryResp.SuccessResp)
	} else {
		log.Fatalf("queryResp failed:%+v", queryResp.ErrResponse)
	}

	var queryServerResponse interface{}

	querySuccess := queryResp.SuccessResp
	switch querySuccess.GetServerResponseMode() {
	case v1.QueryServerResponseUnknownMode:
		log.Fatalln("unknown mode")
	case v1.QueryIndividualRecordingServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryIndividualRecordingServerResponseMode)
		queryServerResponse = querySuccess.GetIndividualRecordingServerResponse()
	case v1.QueryIndividualVideoScreenshotServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryIndividualVideoScreenshotServerResponseMode)
		queryServerResponse = querySuccess.GetIndividualVideoScreenshotServerResponse()
	case v1.QueryMixRecordingHlsServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryMixRecordingHlsServerResponseMode)
		queryServerResponse = querySuccess.GetMixRecordingHLSServerResponse()
	case v1.QueryMixRecordingHlsAndMp4ServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryMixRecordingHlsAndMp4ServerResponseMode)
		queryServerResponse = querySuccess.GetMixRecordingHLSAndMP4ServerResponse()
	case v1.QueryWebRecordingServerResponseMode:
		log.Printf("serverResponseMode:%d", v1.QueryWebRecordingServerResponseMode)
		queryServerResponse = querySuccess.GetWebRecording2CDNServerResponse()
	}

	log.Printf("queryServerResponse:%+v", queryServerResponse)

	time.Sleep(3 * time.Second)
	updateResp, err := cloudRecordingAPI.V1().Update().Do(ctx, startSuccessResp.ResourceId, startSuccessResp.SID, mode, &v1.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &v1.UpdateClientRequest{
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
	resp, err := webRecordingV1.Acquire().Do(ctx, cname, uid, &v1.AcquirerWebRecodingClientRequest{
		ResourceExpiredHour: 24,
	})
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
	starterResp, err := webRecordingV1.Start().Do(ctx, resourceId, cname, uid, &v1.StartWebRecordingClientRequest{
		AppsCollection: &v1.AppsCollection{
			CombinationPolicy: "default",
		},
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
	if starterResp.IsSuccess() {
		log.Printf("starterResp success:%+v", &starterResp.SuccessResp)
	} else {
		log.Fatalf("starterResp failed:%+v", &starterResp.ErrResponse)
	}

	startSuccessResp := starterResp.SuccessResp
	sid := startSuccessResp.SID

	defer func() {
		// stop
		stopResp, err := webRecordingV1.Stop().Do(ctx, resourceId, sid, &v1.StopReqBody{
			Cname: cname,
			Uid:   uid,
			ClientRequest: &v1.StopClientRequest{
				AsyncStop: false,
			},
		})
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
