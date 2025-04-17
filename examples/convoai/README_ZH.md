# 对话式 AI 引擎服务

 [English](./README.md) |  简体中文

> 这是一个对话式 AI 引擎的示例项目，利用对话式 AI 引擎 API 实现对话式 AI 功能。

## 前提条件

设置环境变量。基本环境变量包括：

```bash
export APP_ID=<Your App ID>
export BASIC_AUTH_USERNAME=<您的基本认证用户名>
export BASIC_AUTH_PASSWORD=<您的基本认证密码>
export CONVOAI_TOKEN=<您的代理令牌>
export CONVOAI_CHANNEL=<您的频道名称>
export CONVOAI_AGENT_RTC_UID=<您的代理 RTC UID>
```

您可以在 [对话式 AI 服务文档](../../services/convoai/README_ZH.md) 中找到相关参数。

选择不同的 TTS 提供商时，需要配置额外的环境变量。目前支持的 TTS 提供商有：

### bytedance

```bash
export CONVOAI_TTS_BYTEDANCE_TOKEN=<Your tts bytedance token>
export CONVOAI_TTS_BYTEDANCE_APP_ID=<Your tts bytedance app id>
export CONVOAI_TTS_BYTEDANCE_CLUSTER=<Your tts bytedance cluster>
export CONVOAI_TTS_BYTEDANCE_VOICE_TYPE=<Your tts bytedance voice type>
```

### tencent

```bash
export CONVOAI_TTS_TENCENT_APP_ID=<Your tts tencent app id>
export CONVOAI_TTS_TENCENT_SECRET_ID=<Your tts tencent secret id>
export CONVOAI_TTS_TENCENT_SECRET_KEY=<Your tts tencent secret key>
```

### minimax

```bash
export CONVOAI_TTS_MINIMAX_GROUP_ID=<Your tts minimax group id>
export CONVOAI_TTS_MINIMAX_GROUP_KEY=<Your tts minimax group key>
export CONVOAI_TTS_MINIMAX_GROUP_MODEL=<Your tts minimax group model>
```

### microsoft

```bash
export CONVOAI_TTS_MICROSOFT_KEY=<Your tts microsoft key>
export CONVOAI_TTS_MICROSOFT_REGION=<Your tts microsoft region>
export CONVOAI_TTS_MICROSOFT_VOICE_NAME=<Your tts microsoft voice name>
```

### elevenLabs

```bash
export CONVOAI_TTS_ELEVENLABS_API_KEY=<Your tts elevenLabs api key>
export CONVOAI_TTS_ELEVENLABS_MODEL_ID=<Your tts elevenLabs model id>
export CONVOAI_TTS_ELEVENLABS_VOICE_ID=<Your tts elevenLabs voice id>
```

## 执行

使用以下命令运行示例项目：

```bash
go run main.go --ttsVendor=<ttsVendor> --serviceRegion=<serviceRegion>
```

`ttsVendor` 代表不同的 TTS 提供商。根据您的需求选择合适的 TTS 提供商。
`serviceRegion` 代表选择的服务区域。目前支持的服务区域有：
* `1`:`ChineseMainland`
* `2`:`Global`

