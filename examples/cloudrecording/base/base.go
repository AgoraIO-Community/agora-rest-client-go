package base

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	agoraLogger "github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording"
)

type Service struct {
	DomainArea           domain.Area
	AppId                string
	Cname                string
	Uid                  string
	Credential           auth.Credential
	CloudRecordingClient *cloudrecording.Client
}

func NewService(region domain.Area, appId, cname, uid string, username, password string) (*Service, error) {
	s := &Service{
		DomainArea: region,
		AppId:      appId,
		Cname:      cname,
		Uid:        uid,
		Credential: auth.NewBasicAuthCredential(username, password),
	}

	cloudRecordingClient, err := cloudrecording.NewClient(&cloudrecording.Config{
		AppID:      s.AppId,
		Credential: s.Credential,
		DomainArea: s.DomainArea,
		Logger:     agoraLogger.NewDefaultLogger(agoraLogger.DebugLevel),
	})
	if err != nil {
		return nil, err
	}

	s.CloudRecordingClient = cloudRecordingClient

	return s, nil
}
