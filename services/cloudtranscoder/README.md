# 云端转码服务

## 服务简介

声网的云端转码服务是专为实时互动产品中的直播场景而开发。云端转码服务允许你在服务器端获取 RTC 频道中主播的音视频源流，并对其进行转码、混音、合图等音视频处理，最后将处理后的音视频流发布到声网的 RTC 频道，供观众端订阅。通过使用云端转码服务，观众无需订阅多个主播的音视频流，从而减轻了下行带宽压力和客户端设备的性能消耗。

## 环境准备

- 获取声网App ID -------- [声网Agora - 文档中心 - 如何获取 App ID](https://docs.agora.io/cn/Agora%20Platform/get_appid_token?platform=All%20Platforms#%E8%8E%B7%E5%8F%96-app-id)

  > - 点击创建应用
  >
  >   ![](https://accktvpic.oss-cn-beijing.aliyuncs.com/pic/github_readme/create_app_1.jpg)
  >
  > - 选择你要创建的应用类型
  >
  >   ![](https://accktvpic.oss-cn-beijing.aliyuncs.com/pic/github_readme/create_app_2.jpg)

- 获取App 证书 ----- [声网Agora - 文档中心 - 获取 App 证书](https://docs.agora.io/cn/Agora%20Platform/get_appid_token?platform=All%20Platforms#%E8%8E%B7%E5%8F%96-app-%E8%AF%81%E4%B9%A6)

  > 在声网控制台的项目管理页面，找到你的项目，点击配置。
  > ![](https://fullapp.oss-cn-beijing.aliyuncs.com/scenario_api/callapi/config/1641871111769.png)
  > 点击主要证书下面的复制图标，即可获取项目的 App 证书。
  > ![](https://fullapp.oss-cn-beijing.aliyuncs.com/scenario_api/callapi/config/1637637672988.png)

- 开启云端转码服务
请联系[声网技术支持](https://docportal.shengwang.cn/cn/Agora%20Platform/ticket?platform=Android&_gl=1%2a19d2qxx%2a_gcl_au%2aMTg0ODkxMDM3My4xNzIwNTExNDM3%2a_ga%2aMTI2ODMxNDY2OC4xNjg0MjkxODI0%2a_ga_BFVGG7E02W%2aMTcyMDUxMTIyMC4zMDIuMS4xNzIwNTExNjA0LjAuMC4w)

## API V1 接口调用示例

### 获取云端转码资源
>
> 在开始创建云端转码之前，你需要调用 acquire 方法获取一个 tokenName。一个 builderToken 仅可用于一次云端转码任务。

需要设置的参数有：

- appId: 声网的项目 AppID
- username: 声网的Basic Auth认证的用户名
- password: 声网的Basic Auth认证的密码
- instanceId: 用户指定的实例 ID

通过调用`Acquire().Do`方法来实现获取云端转码资源

```go

 ctx := context.Background()
 c := core.NewClient(&core.Config{
  AppID:      appId,
  Credential: core.NewBasicAuthCredential(username, password),
  RegionCode: core.CN,
  Logger:     core.NewDefaultLogger(core.LogDebug),
 })

 v1Impl := cloudtranscoder.NewAPI(c).V1()
    
 acquireResp, err := v1Impl.Acquire().Do(ctx, &v1.AcquireReqBody{
  InstanceId: instanceId,
 })
 if err != nil {
  log.Fatalln(err)
 }
 if acquireResp.IsSuccess() {
  log.Printf("acquire success:%+v\n", acquireResp)
 } else {
  log.Fatalf("acquire failed:%+v\n", acquireResp)
 }
```

### 开启云端转码
>
> 通过 acquire 方法获取云端转码资源后，调用 create 方法开始云端转码。

需要设置的参数有：

- builderToken： 通过 acquire 方法获取的 tokenName
- 更多 Body中的参数见[Create](https://doc.shengwang.cn/doc/cloud-transcoder/restful/cloud-transcoder/operations/post-v1-projects-appId-rtsc-cloud-transcoder-tasks)接口文档

通过调用`Create().Do`方法来实现创建云端转码

```go
 createResp, err := v1Impl.Create().Do(ctx, tokenName, &v1.CreateReqBody{
  Services: &v1.CreateReqServices{
   CloudTranscoder: &v1.CloudTranscoderPayload{
    ServiceType: "cloudTranscoderV2",
    Config: &v1.CloudTranscoderConfig{
     Transcoder: &v1.CloudTranscoderConfigPayload{
      IdleTimeout: 300,
      AudioInputs: []v1.CloudTranscoderAudioInput{
       {
        Rtc: &v1.CloudTranscoderRtc{
         RtcChannel: "test-abc",
         RtcUID:     123,
         RtcToken:   "xxxxxx",
        },
       },
      },
      Outputs: []v1.CloudTranscoderOutput{
       {
        Rtc: &v1.CloudTranscoderRtc{
         RtcChannel: "test-efg",
         RtcUID:     456,
         RtcToken:   "xxxxx",
        },
        AudioOption: &v1.CloudTranscoderOutputAudioOption{
         ProfileType: "AUDIO_PROFILE_MUSIC_STANDARD",
        },
       },
      },
     },
    },
   },
  },
 })
 if err != nil {
  log.Fatalln(err)
 }

 if createResp.IsSuccess() {
  log.Printf("create success:%+v\n", createResp)
 } else {
  log.Printf("create failed:%+v\n", createResp)
  return
 }
```

### 查询云端转码
>
> 开始云端转码后，你可以调用 query 方法查询云端转码状态。

需要设置的参数有：

- taskId: 从 Create 方法获取到的 taskId
- builderToken： 通过 acquire 方法获取的 tokenName

通过调用`Query().Do`方法来实现查询云端转码状态：

```go
  queryResp, err := v1Impl.Query().Do(ctx, taskId, tokenName)
  if err != nil {
   log.Println(err)
   return
  }

  if queryResp.IsSuccess() {
   log.Printf("query success:%+v\n", queryResp)
  } else {
   log.Printf("query failed:%+v\n", queryResp)
   return
  }
```

### 更新云端转码
>
> 开始云端转码后，你可以调用 update 方法更新云端转码状态。

需要设置的参数有：

- taskId: 从 Create 方法获取到的 taskId
- builderToken： 通过 acquire 方法获取的 tokenName
- sequenceId: Update 请求的序列号。取值需要大于或等于 0。请确保后一次 Update 请求的序列号大于前一次 Update 请求的序列号。序列号可以确保声网服务器按照你指定的最新配置来更新 cloud transcoder。

- 更多 Body中的参数见[Update](https://doc.shengwang.cn/doc/cloud-transcoder/restful/cloud-transcoder/operations/patch-v1-projects-appId-rtsc-cloud-transcoder-tasks-taskId)接口文档

通过调用`Update().Do`方法来实现更新云端转码

```go
 updateResp, err := v1Impl.Update().Do(ctx, taskId, tokenName, 1, &v1.UpdateReqBody{
  Services: &v1.CreateReqServices{
   CloudTranscoder: &v1.CloudTranscoderPayload{
    ServiceType: "cloudTranscoderV2",
    Config: &v1.CloudTranscoderConfig{
     Transcoder: &v1.CloudTranscoderConfigPayload{
      IdleTimeout: 300,
      AudioInputs: []v1.CloudTranscoderAudioInput{
       {
        Rtc: &v1.CloudTranscoderRtc{
         RtcChannel: "test-abc",
         RtcUID:     123,
         RtcToken:   "xxxxxx",
        },
       },
      },
      Outputs: []v1.CloudTranscoderOutput{
       {
        Rtc: &v1.CloudTranscoderRtc{
         RtcChannel: "test-efg",
         RtcUID:     456,
         RtcToken:   "xxxxx",
        },
        AudioOption: &v1.CloudTranscoderOutputAudioOption{
         ProfileType: "AUDIO_PROFILE_MUSIC_HIGH_QUALITY_STEREO",
        },
       },
      },
     },
    },
   },
  },
 })
 if err != nil {
  log.Println(err)
  return
 }

 if updateResp.IsSuccess() {
  log.Printf("update success:%+v\n", updateResp)
 } else {
  log.Printf("update failed:%+v\n", updateResp)
  return
 }
```

### 停止云端转码
>
> 如果你不再需要云端转码，你可以发起 Delete 请求销毁。

需要设置的参数有：

- taskId: 从 Create 方法获取到的 taskId
- builderToken： 通过 acquire 方法获取的 tokenName

通过调用`Delete().Do`方法来实现停止云端转码

```go
  deleteResp, err := v1Impl.Delete().Do(ctx, taskId, tokenName)
  if err != nil {
   log.Println(err)
   return
  }
  if deleteResp.IsSuccess() {
   log.Printf("delete success:%+v\n", deleteResp)
  } else {
   log.Printf("delete failed:%+v\n", deleteResp)
   return
  }
```

## 错误码和响应状态码处理

具体的业务响应码请参考[响应状态码](https://doc.shengwang.cn/doc/cloud-transcoder/restful/response-code)文档
