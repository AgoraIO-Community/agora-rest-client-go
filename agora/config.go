package agora

import (
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Config struct {
	AppID       string
	HttpTimeout time.Duration
	Credential  auth.Credential

	DomainArea domain.Area
	Logger     log.Logger
}
