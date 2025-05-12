# Agora REST Client for Go

<p>
<img alt="GitHub License" src="https://img.shields.io/github/license/AgoraIO-Community/agora-rest-client-go">
<a href="https://pkg.go.dev/github.com/AgoraIO-Community/agora-rest-client-go"><img src="https://pkg.go.dev/badge/github.com/AgoraIO-Community/agora-rest-client-go.svg" alt="Go Reference"></a>
<a href="https://github.com/AgoraIO-Community/agora-rest-client-go/actions/workflows/go.yml"><img src="https://github.com/AgoraIO-Community/agora-rest-client-go/actions/workflows/go.yml/badge.svg" alt="Go Actions"></a>
<a href="https://github.com/AgoraIO-Community/agora-rest-client-go/actions/workflows/gitee-sync.yml"><img alt="gitee-sync" src="https://github.com/AgoraIO-Community/agora-rest-client-go/actions/workflows/gitee-sync.yml/badge.svg?branch=main"></a>
<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/AgoraIO-Community/agora-rest-client-go">
<img alt="GitHub" src="https://img.shields.io/github/v/release/AgoraIO-Community/agora-rest-client-go">
<img alt="GitHub Issues or Pull Requests" src="https://img.shields.io/github/issues-pr/AgoraIO-Community/agora-rest-client-go">
</p>

English | [简体中文](./README_ZH.md)

`agora-rest-client-go` is an open-source project written in Go, specifically designed for the Agora REST API. It includes wrappers and internal implementations of the official Agora REST API interfaces, making it easier for developers to integrate the server-side Agora REST API.

> [!IMPORTANT]
> This SDK has undergone some testing to ensure its basic functionality works correctly. However, due to the complexity of software development, we cannot guarantee it is completely free of defects. We encourage community developers and users to actively participate and help improve this project.

## Features

-   Encapsulates the request and response handling of the Agora REST API, simplifying the communication process with the Agora REST API.
-   Provides automatic switching to the best domain in case of DNS resolution failure, network errors, or request timeouts, ensuring the availability of the REST API service.
-   Offers easy-to-use APIs to easily implement common functions of the Agora REST API, such as starting and stopping cloud recording.
-   Based on Go language, it is efficient, concurrent, and scalable.

## Supported Services

-   [Cloud Recording](./services/cloudrecording/README.md)
-   [Cloud Transcoder](./services/cloudtranscoder/README.md)
-   [Conversational AI Engine](./services/convoai/README.md)

## Environment Setup

-   [Go 1.18 or later](https://go.dev/)
-   App ID and App Certificate obtained from the [Agora Console](https://console.agora.io/v2)
-   Basic Auth credentials from the [Agora Console](https://console.agora.io/v2)
-   Enable the relevant service capabilities on the [Agora Console](https://console.agora.io/v2)

## Installation

Install the dependency from GitHub using the following command:

```shell
go get -u github.com/AgoraIO-Community/agora-rest-client-go
```

## Usage Example

Here is an example of calling the cloud recording service:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	agoraLogger "github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	cloudRecordingAPI "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/req"
)

const (
	appId    = "<your appId>"
	cname    = "<your cname>"
	uid      = "<your uid>"
	username = "<the username of basic auth credential>"
	password = "<the password of basic auth credential>"
	token    = "<your token>"
)

var storageConfig = &cloudRecordingAPI.StorageConfig{
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
	// Initialize Cloud Recording Config
	config := &cloudrecording.Config{
		AppID:      appId,
		Credential: auth.NewBasicAuthCredential(username, password),
		// Specify the region where the server is located. Options include CN, EU, AP, US.
		// The client will automatically switch to use the best domain based on the configured region.
		DomainArea: domain.CN,
		// Specify the log output level. Options include DebugLevel, InfoLevel, WarningLevel, ErrLevel.
		// To disable log output, set logger to DiscardLogger.
		Logger: agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	// Initialize the cloud recording service client
	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Call the Acquire API of the cloud recording service client
	acquireResp, err := cloudRecordingClient.MixRecording().
		Acquire(context.TODO(), cname, uid, &req.AcquireMixRecodingClientRequest{})
		// Handle non-business errors
	if err != nil {
		log.Fatal(err)
	}

	// Handle business response
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Fatalf("acquire failed:%+v\n", acquireResp)
	}

	// Call the Start API of the cloud recording service client
	resourceId := acquireResp.SuccessRes.ResourceId
	startResp, err := cloudRecordingClient.MixRecording().
		Start(context.TODO(), resourceId, cname, uid, &req.StartMixRecordingClientRequest{
			Token: token,
			RecordingConfig: &cloudRecordingAPI.RecordingConfig{
				ChannelType:  1,
				StreamTypes:  2,
				MaxIdleTime:  30,
				AudioProfile: 2,
				TranscodingConfig: &cloudRecordingAPI.TranscodingConfig{
					Width:            640,
					Height:           640,
					FPS:              15,
					BitRate:          800,
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
			RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
				AvFileType: []string{
					"hls",
					"mp4",
				},
			},
			StorageConfig: storageConfig,
		})
		// Handle non-business errors
	if err != nil {
		log.Fatal(err)
	}

	// Handle business response
	if startResp.IsSuccess() {
		log.Printf("start success:%+v\n", startResp)
	} else {
		log.Fatalf("start failed:%+v\n", startResp)
	}

	sid := startResp.SuccessResponse.Sid
	// Query
	for i := 0; i < 6; i++ {
		queryResp, err := cloudRecordingClient.MixRecording().
			QueryHLSAndMP4(context.TODO(), resourceId, sid)
			// Handle non-business errors
		if err != nil {
			log.Fatal(err)
		}

		// Handle business response
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
		}
		time.Sleep(time.Second * 10)
	}

	// Call the Stop API of the cloud recording service client
	stopResp, err := cloudRecordingClient.MixRecording().
		Stop(context.TODO(), resourceId, sid, cname, uid, true)
		// Handle non-business errors
	if err != nil {
		log.Fatal(err)
	}

	// Handle business response
	if stopResp.IsSuccess() {
		log.Printf("stop success:%+v\n", stopResp)
	} else {
		log.Printf("stop failed:%+v\n", stopResp)
	}
}


```

For more examples, see [Example](./examples).

## Contribution

This project welcomes and accepts contributions. If you encounter any issues or have suggestions for improvements, please open an issue or submit a Pull Request.

# SemVer Versioning

This project uses Semantic Versioning (SemVer) to manage versions. The format is MAJOR.MINOR.PATCH.

-   MAJOR version indicates incompatible changes.
-   MINOR version indicates backward-compatible new features or enhancements.
-   PATCH version indicates backward-compatible bug fixes and maintenance.
    For more details, please refer to the [Semantic Versioning](https://semver.org) specification.

## References

-   [Agora API Documentation](https://docs.agora.io/en/)

## License

This project is licensed under the MIT License. For more details, please refer to the LICENSE file.
