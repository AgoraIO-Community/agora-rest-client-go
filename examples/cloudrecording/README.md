# CloudRecording Example

> 这是 Agora Cloud Recording 的一个示例项目，使用了 Agora Cloud Recording RESTful API，实现了频道录制的功能。本示例支持合流录制、单流录制和页面录制三种模式。

## 运行示例项目

### 前提条件

配置环境变量，环境变量包括以下参数内容：

```bash
export APP_ID=<Your App ID>
export CNAME=<Your Channel Name>
export USER_ID=<Your User ID>
export BASIC_AUTH_USERNAME=<Your Basic Auth Username>
export BASIC_AUTH_PASSWORD=<Your Basic Auth Password>
export TOKEN=<Your Token>
export STORAGE_CONFIG_VENDOR=<Your Storage Vendor>
export STORAGE_CONFIG_REGION=<Your Storage Region>
export STORAGE_CONFIG_BUCKET=<Your Storage Bucket>
export STORAGE_CONFIG_ACCESS_KEY=<Your Storage Access Key>
export STORAGE_CONFIG_SECRET_KEY=<Your Storage Secret Key>
```

相关的参数可以通过可在 [CloudRecording 服务说明](../../services/cloudrecording/README.md) 查看

### 执行

通过下面的命令来运行示例项目：

```bash
go run main.go -mode=mix -mix_scene=<scene>
go run main.go -mode=individual -individual_scene=<scene>
go run main.go -mode=web -web_scene=<scene>
```

其中 `mode` 表示云录制模式：

* mix: 合流录制
* individual: 单流录制
* web: 页面录制

其中 `mix_scene` 表示合流录制场景：

* hls: 录制hls格式
* hls_and_mp4: 录制hls和mp4格式

其中 `individual_scene` 表示单流录制场景：

* recording: 仅录制
* snapshot: 仅截图
* recording_and_snapshot: 录制+截图
* recording_and_postpone_transcoding: 录制+延时转码
* recording_and_audio_mix: 录制+延时混音

其中 `web_scene` 表示页面录制场景：

* web_recorder: 页面录制
* web_recorder_and_rtmp_publish: 页面录制+转推到CDN
