# Cloud Recording Service

English | [简体中文](./README_ZH.md)

## Service Introduction
Cloud Recording is a recording component developed by Agora for audio and video calls and live broadcasts. It provides RESTful APIs for developers to implement recording functionality and store recording files in third-party cloud storage. Cloud Recording offers advantages such as stability, reliability, ease of use, cost control, flexible solutions, and support for private deployment, making it an ideal recording solution for online education, video conferences, financial supervision, and customer service scenarios.

## Environment Preparation

- Obtain Agora App ID -------- [Agora Console](https://console.agora.io/v2)

  > - Click Create Application
  >
  >   ![](../../assets/imges/EN/create_app_1.png)
  >
  > - Select the type of application you want to create
  >
  >   ![](../../assets/imges/EN/create_app_2.png)

- Obtain App Certificate ----- [Agora Console](https://console.agora.io/v2)

  > In the project management page of the Agora Console, find your project and click Configure.
  > ![](../../assets/imges/EN/config_app.png)
  > Click the copy icon under Primary Certificate to obtain the App Certificate for your project.
  > ![](../../assets/imges/EN/copy_app_cert.png)

- Check the status of the recording service
  > ![](../../assets/imges/EN/open_cloud_recording.png)

## API Call Examples
### Acquire Cloud Recording Resources
> Before starting cloud recording, you need to call the acquire method to obtain a Resource ID. A Resource ID can only be used for one cloud recording service.

Parameters that need to be set:
- appId: Agora project AppID
- username: Agora Basic Auth authentication username
- password: Agora Basic Auth authentication password
- cname: Channel name
- uid: User UID
- For more parameters in clientRequest, see the [Acquire](https://docs.agora.io/en/cloud-recording/reference/restful-api#acquire) API documentation

Implement acquiring cloud recording resources by calling the `Acquire` method
```go
    appId := "xxxx"
    username := "xxxx"
    password := "xxxx"
    credential := auth.NewBasicAuthCredential(username, password)
	config := &agora.Config{
		AppID:      appId,
		Credential: credential,
        DomainArea: domain.CN,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	}

	cloudRecordingClient, err := cloudrecording.NewClient(config)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := cloudRecordingClient.Acquire(context.TODO(), &cloudRecordingAPI.AcquireReqBody{
		Cname: "12321",
		Uid:   "43434",
		ClientRequest: &cloudRecordingAPI.AcquireClientRequest{
			Scene:               0,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Printf("resourceId:%s", resp.SuccessRes.ResourceId)
	} else {
		log.Printf("resp:%+v", resp)
	}
```

### Start Cloud Recording
> After acquiring cloud recording resources through the acquire method, call the start method to begin cloud recording.

Parameters that need to be set:
- cname: Channel name
- uid: User UID
- token: Token corresponding to the user UID
- resourceId: Cloud recording resource ID
- mode: Cloud recording mode
- For more parameters in clientRequest, see the [Start](https://docs.agora.io/en/cloud-recording/reference/restful-api#start) API documentation

Implement starting cloud recording by calling the `Start` method
```go
	startResp, err := cloudRecordingClient.Start(ctx, resp.SuccessRes.ResourceId, mode, &cloudRecordingAPI.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &cloudRecordingAPI.StartClientRequest{
			Token: token,
			RecordingConfig: &cloudRecordingAPI.RecordingConfig{
				ChannelType:  1,
				StreamTypes:  2,
				AudioProfile: 2,
				MaxIdleTime:  30,
				TranscodingConfig: &cloudRecordingAPI.TranscodingConfig{
					Width:            640,
					Height:           260,
					FPS:              15,
					BitRate:          500,
					MixedVideoLayout: 0,
					BackgroundColor:  "#000000",
				},
				SubscribeAudioUIDs: []string{
					"22",
					"456",
				},
				SubscribeVideoUIDs: []string{
					"22",
					"456",
				},
			},
			RecordingFileConfig: &cloudRecordingAPI.RecordingFileConfig{
				AvFileType: []string{
					"hls",
				},
			},
			StorageConfig: &cloudRecordingAPI.StorageConfig{
				Vendor:    2,
				Region:    3,
				Bucket:    "xxx",
				AccessKey: "xxxx",
				SecretKey: "xxx",
				FileNamePrefix: []string{
					"xx1",
					"xx2",
				},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if startResp.IsSuccess() {
		log.Printf("success:%+v", &startResp.SuccessResponse)
	} else {
		log.Printf("failed:%+v", &startResp.ErrResponse)
	}
```

### Stop Cloud Recording
> After starting recording, you can call the stop method to leave the channel and stop recording. If you need to record again after stopping, you must call the acquire method again to request a new Resource ID.

Parameters that need to be set:
- cname: Channel name
- uid: User ID
- resourceId: Cloud recording resource ID
- sid: Session ID
- mode: Cloud recording mode
- For more parameters in clientRequest, see the [Stop](https://docs.agora.io/en/cloud-recording/reference/restful-api#stop) API documentation

Since the Stop interface does not return a fixed structure, you need to determine the specific return type based on the serverResponseMode returned

Implement stopping cloud recording by calling the `Stop` method
```go
    stopResp, err := cloudRecordingClient.Stop(ctx, resourceId, sid, mode, &cloudRecordingAPI.cloudRecordingAPI{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &cloudRecordingAPI.StopClientRequest{
			AsyncStop: true,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if stopResp.IsSuccess() {
		log.Printf("stop success:%+v", &stopResp.SuccessResponse)
	} else {
		log.Fatalf("stop failed:%+v", &stopResp.ErrResponse)
	}
```

### Query Cloud Recording Status
> After starting recording, you can call the query method to check the recording status.

Parameters that need to be set:
- cname: Channel name
- uid: User ID
- resourceId: Cloud recording resource ID
- sid: Session ID
- mode: Cloud recording mode
- For more parameters in clientRequest, see the [Query](https://docs.agora.io/en/cloud-recording/reference/restful-api#query) API documentation

Since the Query interface does not return a fixed structure, you need to determine the specific return type based on the serverResponseMode returned

Implement querying cloud recording status by calling the `Query` method
```go
	queryResp, err := cloudRecordingClient.Query(ctx, resourceId, sid, mode)
	if err != nil {
		log.Fatalln(err)
	}
	if queryResp.IsSuccess() {
		log.Printf("query success:%+v", queryResp.SuccessResponse)
	} else {
		log.Printf("query failed:%+v", queryResp.ErrResponse)
		return
	}

	var queryServerResponse interface{}

	querySuccess := queryResp.SuccessResponse
	switch querySuccess.GetServerResponseMode() {
	case cloudRecordingAPI.QueryServerResponseUnknownMode:
		log.Fatalln("unknown mode")
	case cloudRecordingAPI.QueryIndividualRecordingServerResponseMode:
		log.Printf("serverResponseMode:%d", cloudRecordingAPI.QueryIndividualRecordingServerResponseMode)
		queryServerResponse = querySuccess.GetIndividualRecordingServerResponse()
	case cloudRecordingAPI.QueryIndividualVideoScreenshotServerResponseMode:
		log.Printf("serverResponseMode:%d", cloudRecordingAPI.QueryIndividualVideoScreenshotServerResponseMode)
		queryServerResponse = querySuccess.GetIndividualVideoScreenshotServerResponse()
	case cloudRecordingAPI.QueryMixRecordingHlsServerResponseMode:
		log.Printf("serverResponseMode:%d", cloudRecordingAPI.QueryMixRecordingHlsServerResponseMode)
		queryServerResponse = querySuccess.GetMixRecordingHLSServerResponse()
	case cloudRecordingAPI.QueryMixRecordingHlsAndMp4ServerResponseMode:
		log.Printf("serverResponseMode:%d", cloudRecordingAPI.QueryMixRecordingHlsAndMp4ServerResponseMode)
		queryServerResponse = querySuccess.GetMixRecordingHLSAndMP4ServerResponse()
	case cloudRecordingAPI.QueryWebRecordingServerResponseMode:
		log.Printf("serverResponseMode:%d", cloudRecordingAPI.QueryWebRecordingServerResponseMode)
		queryServerResponse = querySuccess.GetWebRecording2CDNServerResponse()
	}

	log.Printf("queryServerResponse:%+v", queryServerResponse)
```

### Update Cloud Recording Settings
> After starting recording, you can call the update method to update the following recording configurations:
> * For individual recording and composite recording, update the subscription list.
> * For web recording, set pause/resume web recording, or update the streaming URL for pushing web recording to CDN.

Parameters that need to be set:
- cname: Channel name
- uid: User UID
- resourceId: Cloud recording resource ID
- sid: Session ID
- mode: Cloud recording mode
- For more parameters in clientRequest, see the [Update](https://docs.agora.io/en/cloud-recording/reference/restful-api#update) API documentation

Implement updating cloud recording settings by calling the `Update` method
```go
	updateResp, err := cloudRecordingClient.Update(ctx, resourceId, sid, mode, &cloudRecordingAPI.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &cloudRecordingAPI.UpdateClientRequest{
			StreamSubscribe: &cloudRecordingAPI.UpdateStreamSubscribe{
				AudioUidList: &cloudRecordingAPI.UpdateAudioUIDList{
					SubscribeAudioUIDs: []string{
						"999",
					},
				},
				VideoUidList: &cloudRecordingAPI.UpdateVideoUIDList{
					SubscribeVideoUIDs: []string{
						"999",
					},
				},
			},
		},
	})

	if err != nil {
		log.Fatalln(err)
	}
	if updateResp.IsSuccess() {
		log.Printf("update success:%+v", updateResp.SuccessResponse)
	} else {
		log.Printf("update failed:%+v", updateResp.ErrResponse)
		return
	}
```

### Update Cloud Recording Composite Layout
> After starting recording, you can call the updateLayout method to update the composite layout.
> Each call to this method will overwrite the previous layout settings. For example, if you set backgroundColor to "#FF0000" (red) when starting recording, if you call the updateLayout method to update the composite layout without setting the backgroundColor field again, the background color will change to black (default value).

Parameters that need to be set:
- cname: Channel name
- uid: User UID
- resourceId: Cloud recording resource ID
- sid: Session ID
- mode: Cloud recording mode
- For more parameters in clientRequest, see the [UpdateLayout](https://docs.agora.io/en/cloud-recording/reference/restful-api#updatelayout) API documentation

Implement updating cloud recording composite layout by calling the `UpdateLayout` method
```go
	updateLayoutResp, err := cloudRecordingClient.UpdateLayout(ctx, resourceId, sid, mode, &cloudRecordingAPI.UpdateLayoutReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &cloudRecordingAPI.UpdateLayoutClientRequest{
			MixedVideoLayout: 3,
			BackgroundColor:  "#FF0000",
			LayoutConfig: []cloudRecordingAPI.UpdateLayoutConfig{
				{
					UID:        "22",
					XAxis:      0.1,
					YAxis:      0.1,
					Width:      0.1,
					Height:     0.1,
					Alpha:      1,
					RenderMode: 1,
				},
				{
					UID:        "2",
					XAxis:      0.2,
					YAxis:      0.2,
					Width:      0.1,
					Height:     0.1,
					Alpha:      1,
					RenderMode: 1,
				},
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if updateLayoutResp.IsSuccess() {
		log.Printf("updateLayout success:%+v", updateLayoutResp.SuccessResponse)
	} else {
		log.Printf("updateLayout failed:%+v", updateLayoutResp.ErrResponse)
		return
	}
```

## Error Codes and Response Status Codes
For specific business response codes, please refer to the [Business Response Codes](https://docs.agora.io/en/cloud-recording/reference/common-errors) documentation