package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/individualrecording"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/mixrecording"
	"github.com/AgoraIO-Community/agora-rest-client-go/examples/cloudrecording/webrecording"
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
		mixrecording.MixRecording(appId, username, password, token, cname, uid, storageConfig, region)
	case 2:
		individualrecording.IndividualRecording(appId, username, password, token, cname, uid, storageConfig, region)
	case 3:
		webrecording.WebRecording(appId, username, password, token, cname, uid, storageConfig, region)
	default:
		panic("mode is required, 1: mix, 2: individual, 3: webRecording")
	}
}
