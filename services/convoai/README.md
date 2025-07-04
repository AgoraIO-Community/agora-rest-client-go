# Conversational AI Engine Service

English | [简体中文](./README_ZH.md)

## Service Overview

Agora's Conversational AI Engine redefines human-computer interaction, breaking through traditional text interactions to achieve highly realistic, natural, and smooth real-time voice conversations, allowing AI to truly "speak." It is suitable for innovative scenarios such as intelligent assistants, emotional companionship, spoken language training, intelligent customer service, smart hardware, and immersive game NPCs.

## Environment Setup

-   Obtain Agora App ID -------- [Agora Console](https://console.agora.io/v2)

    > -   Click Create Application
    >
    >     ![](../../assets/imges/EN/create_app_1.png)
    >
    > -   Select the type of application you want to create
    >
    >     ![](../../assets/imges/EN/create_app_2.png)

-   Obtain App Certificate ----- [Agora Console](https://console.agora.io/v2)

    > In the project management page of the Agora Console, find your project and click Configure.
    > ![](../../assets/imges/EN/config_app.png)
    > Click the copy icon under Primary Certificate to obtain the App Certificate for your project.
    > ![](../../assets/imges/EN/copy_app_cert.png)

-   Enable Conversational AI Engine Service ----- [Enable Service](https://docs.agora.io/en/conversational-ai/get-started/manage-agora-account)
    > ![](../../assets/imges/EN/open_convo_ai.png)

## API Definition

For more api details, please refer to the [API Documentation](https://docs.agora.io/en/conversational-ai/rest-api/join)

## API Call Examples

### Initialize Conversational AI Engine Client

```go
    const (
        appId                 = "<your appId>"
        username              = "<the username of basic auth credential>"
        password              = "<the password of basic auth credential>"
	)
	// Initialize Conversational AI Config
	config := &convoai.Config{
		AppID:      appId,
		Credential: auth.NewBasicAuthCredential(username, password),
		// Specify the region where the server is located. Options include CN, EU, AP, US.
		// The client will automatically switch to use the best domain based on the configured region.
		DomainArea: domain.US,
		// Specify the log output level. Options include DebugLevel, InfoLevel, WarningLevel, ErrLevel.
		// To disable log output, set logger to DiscardLogger.
		Logger: agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
		ServiceRegion: convoai.GlobalServiceRegion,
	}

	// Initialize the Conversational AI service client
	convoaiClient, err := convoai.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
```

### Create Conversational Agent

> Create a Conversational AI agent instance and join an RTC channel.

Parameters to set: LLM, TTS, and Agent related parameters.

Call the `Join` method to create a conversational agent, using Microsoft TTS as an example:

```go
    const (
		appId                 = "<your appId>"
        cname                 = "<your cname>"
        agentRtcUid           = "<your agent rtc uid>"
        username              = "<the username of basic auth credential>"
        password              = "<the password of basic auth credential>"
        agentRtcToken         = "<your agent rtc token>"
        llmURL                = "<your LLM URL>"
        llmAPIKey             = "<your LLM API Key>"
        llmModel              = "<your LLM model>"
		ttsMicrosoftKey       = "<your microsoft tts key>"
		ttsMicrosoftRegion    = "<your microsoft tts region>"
		ttsMicrosoftVoiceName = "<your microsoft tts voice name>"
    )
	// Start agent
	name := appId + ":" + cname
	joinResp, err := convoaiClient.Join(context.Background(), name, &req.JoinPropertiesReqBody{
		Token:           agentRtcToken,
		Channel:         cname,
		AgentRtcUId:     agentRtcUid,
		RemoteRtcUIds:   []string{"*"},
		EnableStringUId: agoraUtils.Ptr(false),
		IdleTimeout:     agoraUtils.Ptr(120),
		LLM: &req.JoinPropertiesCustomLLMBody{
			Url:    llmURL,
			APIKey: llmAPIKey,
			SystemMessages: []map[string]any{
				{
					"role":    "system",
					"content": "You are a helpful chatbot.",
				},
			},
			Params: map[string]any{
				"model":      llmModel,
				"max_tokens": 1024,
			},
			MaxHistory:      agoraUtils.Ptr(30),
			GreetingMessage: "Hello, how can I help you?",
		},
		TTS: &req.JoinPropertiesTTSBody{
			Vendor: req.MicrosoftTTSVendor,
			Params: &req.TTSMicrosoftVendorParams{
				Key:        ttsMicrosoftKey,
				Region:     ttsMicrosoftRegion,
				VoiceName:  ttsMicrosoftVoiceName,
				Speed:      1.0,
				Volume:     70,
				SampleRate: 24000,
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	if joinResp.IsSuccess() {
		log.Printf("Join success:%+v", joinResp)
	} else {
		log.Printf("Join failed:%+v", joinResp)
	}

```

### Stop Conversational Agent

> Stop the conversational agent and leave the RTC channel.

Parameters to set:

-   AgentId returned by the `Join` interface

```go
    // Leave agent
    leaveResp, err := convoaiClient.Leave(ctx, agentId)
	if err != nil {
		log.Fatalln(err)
	}

	if leaveResp.IsSuccess() {
		log.Printf("Leave success:%+v", leaveResp)
	} else {
		log.Printf("Leave failed:%+v", leaveResp)
	}
```

### Update Agent Configuration

> Currently, only the Token information of a running conversational agent can be updated.

Parameters to set:

-   AgentId returned by the `Join` interface
-   Token to be updated

```go
    // Update agent
	updateResp, err := convoaiClient.Update(ctx, agentId, &req.UpdateReqBody{
		Token: updateToken,
	})
	if err != nil {
		log.Fatalln(err)
	}

	if updateResp.IsSuccess() {
		log.Printf("Update success:%+v", updateResp)
	} else {
		log.Printf("Update failed:%+v", updateResp)
	}
```

### Query Agent Status

> Query the status of the conversational agent.

Parameters to set:

-   AgentId returned by the `Join` interface

```go
    // Query agent
	queryResp, err := convoaiClient.Query(ctx, agentId)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if queryResp.IsSuccess() {
		log.Printf("Query success:%+v", queryResp)
	} else {
		log.Printf("Query failed:%+v", queryResp)
	}
```

## Retrieves a list of agents

> Retrieves a list of agents that meet the specified criteria.

Parameters to set:

-   AgentId returned by the `Join` interface

```go
     // List agent
	listResp, err := convoaiClient.List(ctx,
		req.WithState(2),
		req.WithLimit(10),
	)
	if err != nil {
		log.Fatalln(err)
	}

	if listResp.IsSuccess() {
		log.Printf("List success:%+v", listResp)
	} else {
		log.Printf("List failed:%+v", listResp)
	}
```

## Error Codes and Response Status Codes Handling

For specific business response codes, please refer to the [Business Response Codes](https://docs.agora.io/en/conversational-ai/rest-api/reference) documentation.
