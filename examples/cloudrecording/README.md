# CloudRecording Example

English | [简体中文](./README_ZH.md)

> This is a sample project for Agora Cloud Recording that utilizes the Agora Cloud Recording RESTful API to implement channel recording functionality. This example supports three recording modes: composite recording, individual recording, and web recording.

## Run

### Prerequisites

Configure the environment variables with the following parameters:

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

Relevant parameters can be found in the [CloudRecording Service Documentation](../../services/cloudrecording/README.md)

### Execution

Run the example project using the following commands:

```bash
go run main.go -mode=mix -mix_scene=<scene>
go run main.go -mode=individual -individual_scene=<scene>
go run main.go -mode=web -web_scene=<scene>
```

Where `mode` indicates the cloud recording mode:

* mix: Composite recording
* individual: Individual recording
* web: Web recording

Where `mix_scene` indicates the composite recording scenario:

* hls: Recording in HLS format
* hls_and_mp4: Recording in both HLS and MP4 formats

Where `individual_scene` indicates the individual recording scenario:

* recording: Recording only
* snapshot: Screenshot only
* recording_and_snapshot: Recording + Screenshot

Where `web_scene` indicates the web recording scenario:

* web_recorder: Web recording
* web_recorder_and_rtmp_publish: Web recording + RTMP streaming to CDN
