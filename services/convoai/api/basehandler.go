package api

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type baseHandler struct {
	module     string
	logger     log.Logger
	client     client.Client
	prefixPath string // /api/conversational-ai-agent/v2/projects/{appid}
}
