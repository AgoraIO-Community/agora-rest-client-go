# Agora REST Client for Go
<p>
<img alt="GitHub License" src="https://img.shields.io/github/license/AgoraIO-Community/agora-rest-client-go">
<a href="https://pkg.go.dev/github.com/AgoraIO-Community/agora-rest-client-go"><img src="https://pkg.go.dev/badge/github.com/AgoraIO-Community/agora-rest-client-go.svg" alt="Go Reference"></a>
<a href="https://github.com/seymourtang/agora-rest-client-go/actions/workflows/go.yml"><img src="https://github.com/AgoraIO-Community/agora-rest-client-go/actions/workflows/go.yml/badge.svg" alt="Go Actions"></a>
<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/AgoraIO-Community/agora-rest-client-go">
<img alt="GitHub" src="https://img.shields.io/github/v/release/AgoraIO-Community/agora-rest-client-go">
<img alt="GitHub Issues or Pull Requests" src="https://img.shields.io/github/issues-pr/AgoraIO-Community/agora-rest-client-go">
</p>

`agora-rest-client-go`是用Go语言编写的一个开源项目，专门为 Agora REST API设计。它包含了 Agora 官方提供的REST API接口的包装和内部实现，可以帮助开发者更加方便的集成服务端Agora REST API。

> 注意： 该SDK经过一些测试以确保其基本功能正常运作。然而，由于软件开发的复杂性，我们无法保证它是完全没有缺陷的。
>
>该SDK目前可能存在一些潜在的BUG或不稳定性。我们鼓励社区的开发者和用户积极参与，共同改进这个项目。

## 特性
* 封装了Agora REST API的请求和响应处理，简化与Agora REST API 的通信流程
* 当遇到 DNS 解析失败、网络错误或者请求超时等问题的时候，提供了自动切换最佳域名的能力，以保障请求 REST API 服务的可用性
* 提供了易于使用的API，可轻松地实现调用 Agora REST API 的常见功能，如开启云录制、停止云录制等
* 基于Go语言，具有高效性、并发性和可扩展性

## 支持的服务
* [云端录制 Cloud Recording ](./services/cloudrecording/README.md)

## 环境准备
* [Go 1.18 或以上版本](https://go.dev/)
* 在声网 [Console 平台](https://console.shengwang.cn/)申请的 App ID 和 App Certificate
* 在声网 [Console 平台](https://console.shengwang.cn/)的 Basic Auth 认证信息
* 在声网 [Console 平台](https://console.shengwang.cn/)开启相关的服务能力

## 安装
使用以下命令从 GitHub 安装依赖：
```shell
go get -u github.com/AgoraIO-Community/agora-rest-client-go
```
## 使用示例
以调用云录制服务为例：
```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

const (
	appId    = "<your appId>"
	cname    = "<your cname>"
	uid      = "<your uid>"
	username = "<the username of basic auth credential>"
	password = "<the password of basic auth credential>"
	token    = "<your token>"
)

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
	// 初始化Agora REST API客户端
	client := core.NewClient(&core.Config{
		AppID:      appId,
		Credential: core.NewBasicAuthCredential(username, password),
		// 指定服务器所在的区域，可选值有CN, NA, EU, AP,client 将会根据配置的区域自动切换使用最佳的域名
		RegionCode: core.CNRegionArea,
		// 指定日志输出的级别，可选值有LogDebug, LogInfo, LogWarn, LogError
		// 如果要关闭日志输出，可将 logger 设置为 DiscardLogger
		Logger: core.NewDefaultLogger(core.LogDebug),
	})

	// 初始化云端录制服务 API
	impl := cloudrecording.NewAPI(client).V1().MixRecording()

	// 调用云端录制服务 API 的 Acquire 接口
	acquireResp, err := impl.Acquire().Do(context.TODO(), cname, uid, &v1.AcquireMixRecodingClientRequest{})
	// 处理非业务错误
	if err != nil {
		log.Fatal(err)
	}

	// 处理业务响应
	if acquireResp.IsSuccess() {
		log.Printf("acquire success:%+v\n", acquireResp)
	} else {
		log.Fatalf("acquire failed:%+v\n", acquireResp)
	}

	// 调用云端录制服务 API 的 Start 接口
	resourceId := acquireResp.SuccessRes.ResourceId
	startResp, err := impl.Start().Do(context.TODO(), resourceId, cname, uid, &v1.StartMixRecordingClientRequest{
		Token: token,
		RecordingConfig: &v1.RecordingConfig{
			ChannelType:  1,
			StreamTypes:  2,
			MaxIdleTime:  30,
			AudioProfile: 2,
			TranscodingConfig: &v1.TranscodingConfig{
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
		RecordingFileConfig: &v1.RecordingFileConfig{
			AvFileType: []string{
				"hls",
				"mp4",
			},
		},
		StorageConfig: storageConfig,
	})
	// 处理非业务错误
	if err != nil {
		log.Fatal(err)
	}

	// 处理业务响应
	if startResp.IsSuccess() {
		log.Printf("start success:%+v\n", startResp)
	} else {
		log.Fatalf("start failed:%+v\n", startResp)
	}

	sid := startResp.SuccessResponse.Sid
	// query
	for i := 0; i < 6; i++ {
		queryResp, err := impl.Query().DoHLSAndMP4(context.TODO(), resourceId, sid)
		// 处理非业务错误
		if err != nil {
			log.Fatal(err)
		}

		// 处理业务响应
		if queryResp.IsSuccess() {
			log.Printf("query success:%+v\n", queryResp)
		} else {
			log.Printf("query failed:%+v\n", queryResp)
		}
		time.Sleep(time.Second * 10)
	}

	// 调用云端录制服务 API 的 Stop 接口
	stopResp, err := impl.Stop().DoHLSAndMP4(context.TODO(), resourceId, sid, cname, uid, false)
	// 处理非业务错误
	if err != nil {
		log.Fatal(err)
	}

	// 处理业务响应
	if stopResp.IsSuccess() {
		log.Printf("stop success:%+v\n", stopResp)
	} else {
		log.Printf("stop failed:%+v\n", stopResp)
	}
}
```
更多的示例可在[Example](./examples) 查看

## 集成遇到困难，该如何联系声网获取协助

> 方案1：如果您已经在使用声网服务或者在对接中，可以直接联系对接的销售或服务
>
> 方案2：发送邮件给 [support@agora.io](mailto:support@agora.io) 咨询
>
> 方案3：扫码加入我们的微信交流群提问
>
> <img src="https://download.agora.io/demo/release/SDHY_QA.jpg" width="360" height="360">
---

## 贡献
本项目欢迎并接受贡献。如果您在使用中遇到问题或有改进建议，请提出issue或向我们提交Pull Request。

# SemVer 版本规范
本项目使用语义化版本号规范 (SemVer) 来管理版本。格式为 MAJOR.MINOR.PATCH。

* MAJOR 版本号表示不向后兼容的重大更改。
* MINOR 版本号表示向后兼容的新功能或增强。
* PATCH 版本号表示向后兼容的错误修复和维护。
有关详细信息，请参阅 [语义化版本](https://semver.org/lang/zh-CN/) 规范。

## 参考
* [Agora API 文档](https://doc.shengwang.cn/)

## 许可证
该项目使用MIT许可证，详细信息请参阅LICENSE文件。