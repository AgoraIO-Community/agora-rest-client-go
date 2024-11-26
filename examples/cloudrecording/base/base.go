package base

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/region"
)

type Service struct {
	RegionArea region.Area
	AppId      string
	Cname      string
	Uid        string
	Credential auth.Credential
}

func (s *Service) SetCredential(username, password string) {
	s.Credential = auth.NewBasicAuthCredential(username, password)
}
