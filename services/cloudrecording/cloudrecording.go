package cloudrecording

import (
	"github.com/AgoraIO/agora-rest-client-go/core"
	v1 "github.com/AgoraIO/agora-rest-client-go/services/cloudrecording/v1"
)

const projectName = "cloud_recording"

type API struct {
	client core.Client
}

func NewAPI(client core.Client) *API {
	return &API{client: client}
}

func (a *API) buildPrefixPath() string {
	return "/apps/" + a.client.GetAppID() + "/" + projectName
}

func (a *API) V1() *v1.Collection {
	return v1.NewCollection(a.buildPrefixPath(), a.client)
}
