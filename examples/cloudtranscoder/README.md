# CloudTranscoder Example

English | [简体中文](./README_ZH.md)

> This is a sample project for Agora Cloud Transcoder that utilizes the Agora Cloud Transcoder RESTful API to implement cloud transcoding functionality. This example supports two scenarios: `Single Channel + RTC Pull Stream Composite Transcoding + Output to RTC` and `Single Channel + RTC Pull Stream Full Channel Audio Mixing Transcoding + Output to RTC`.

## Run

### Prerequisites

Configure the environment variables with the following parameters:

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

Relevant parameters can be found in the [CloudTranscoder Service Documentation](../../services/cloudtranscoder/README.md)

### Execution

Run the example project using the following commands:

```bash
go run main.go -scene=single_channel_rtc_pull_mixer_rtc_push 
go run main.go -scene=single_channel_rtc_pull_fullchannel_audiomixer_rtc_push
```

Where `scene` indicates the cloud transcoding scenario:

* single_channel_rtc_pull_mixer_rtc_push: Single Channel + RTC Pull Stream Composite Transcoding + Output to RTC
* single_channel_rtc_pull_fullchannel_audiomixer_rtc_push: Single Channel + RTC Pull Stream Full Channel Audio Mixing Transcoding + Output to RTC
