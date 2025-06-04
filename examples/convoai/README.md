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
go run main.go --ttsVendor=<ttsVendor> --serviceRegion=2
```

`ttsVendor` represents different TTS providers, currently supported TTS providers are:

-   `microsoft`
-   `elevenLabs`

Choose the appropriate TTS provider based on your requirements.
