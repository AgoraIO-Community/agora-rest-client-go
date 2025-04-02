# CloudTranscoder Example

 [English](./README.md) |  简体中文

> 这是 Agora Cloud Transcoder 的一个示例项目，使用了 Agora Cloud Transcoder RESTful API，实现了云端转码的功能。本示例支持`单频道+RTC拉流合图转码+输出到RTC`和`单频道+RTC拉流全频道混音转码+输出到RTC`两种场景。

## 运行示例项目

### 前提条件

配置环境变量，环境变量包括以下参数内容：

```bash
export APP_ID=<Your App ID>
export BASIC_AUTH_USERNAME=<Your Basic Auth Username>
export BASIC_AUTH_PASSWORD=<Your Basic Auth Password>
export INPUT_UID_1=<Your First Input UID>
export INPUT_UID_2=<Your Second Input UID>
export INPUT_CHANNEL_NAME=<Input Channel Name>
export INPUT_TOKEN_1=<Your First Input Token>
export INPUT_TOKEN_2=<Your Second Input Token>
export UPDATE_INPUT_UID_3=<Your Third Update  UID>
export UPDATE_INPUT_TOKEN_3=<Your Third Update Token>
export OUTPUT_UID=<Your Output UID>
export OUTPUT_TOKEN=<Your Output Token>
export OUTPUT_CHANNEL_NAME=<Your Output Channel Name>
```

相关的参数可以通过可在 [CloudTranscoder 服务说明](../../services/cloudtranscoder/README_ZH.md) 查看

### 执行

通过下面的命令来运行示例项目：

```bash
go run main.go -scene=single_channel_rtc_pull_mixer_rtc_push 
go run main.go -scene=single_channel_rtc_pull_fullchannel_audiomixer_rtc_push
```

其中 `scene` 表示云端转码的场景：

* single_channel_rtc_pull_mixer_rtc_push: 单频道+RTC拉流合图转码+输出到RTC
* single_channel_rtc_pull_fullchannel_audiomixer_rtc_push: 单频道+RTC拉流全频道混音转码+输出到RTC
