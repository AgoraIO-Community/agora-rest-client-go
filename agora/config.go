package agora

import (
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/region"
)

type Config struct {
	AppID       string
	HttpTimeout time.Duration
	Credential  auth.Credential

	RegionCode region.Area
	Logger     log.Logger
}
