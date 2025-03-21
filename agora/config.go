package agora

import (
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
)

type Config struct {
	// Agora AppID
	AppID string
	// Timeout for HTTP requests
	HttpTimeout time.Duration
	// Credential for accessing the Agora service.
	//
	// Available credential types:
	//
	//  - BasicAuthCredential: See auth.NewBasicAuthCredential for details
	Credential auth.Credential

	// Domain area for the REST Client. See domain.Area for details.
	DomainArea domain.Area

	// Logger for the REST Client
	//
	// Implement the log.Logger interface in your project to output REST Client logs to your logging component.
	//
	// Alternatively, you can use the default logging component. See log.NewDefaultLogger for details.
	Logger log.Logger
}
