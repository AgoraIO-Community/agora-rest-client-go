package cloudtrasncoder

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

const projectName = "rtsc/cloud-transcoder"

type API struct {
	client core.Client
}

func NewAPI(client core.Client) *API {
	return &API{client: client}
}

func (a *API) buildPrefixPath() string {
	return "/projects/" + a.client.GetAppID() + "/" + projectName
}
