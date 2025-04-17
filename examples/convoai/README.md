# Conversational AI Service

English | [简体中文](./README_ZH.md)

> This is a sample project for the Agora Conversational AI engine, utilizing the Agora Conversational Engine API to implement conversational AI functionalities.

## Prerequisites

Set up environment variables. The basic environment variables include:

```bash
export APP_ID=<Your App ID>
export BASIC_AUTH_USERNAME=<Your Basic Auth Username>
export BASIC_AUTH_PASSWORD=<Your Basic Auth Password>
export CONVOAI_TOKEN=<Your Agent Token>
export CONVOAI_CHANNEL=<Your Channel Name>
export CONVOAI_AGENT_RTC_UID=<Your Agent RTC UID>
```

You can find the relevant parameters in the [Conversational AI Service Documentation](../../services/convoai/README.md).

When choosing different TTS providers, you need to configure additional environment variables. The currently supported TTS providers are:

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

## Execution

Run the sample project with the following command:

```bash
go run main.go --ttsVendor=<ttsVendor> --serviceRegion=<serviceRegion>
```

`ttsVendor` represents different TTS providers. Choose the appropriate TTS provider based on your requirements.
`serviceRegion` represents chosen service region. The currently supported service regions are:
* `1`:`ChineseMainland`
* `2`:`Global`

