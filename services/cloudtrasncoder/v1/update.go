package v1

import "github.com/AgoraIO-Community/agora-rest-client-go/core"

type Update struct {
	client     core.Client
	prefixPath string // /v1/apps/{appid}/cloud_recording
}

// buildPath returns the request path.
// /v1/projects/{appid}/rtsc/cloud-transcoder/tasks/{taskId}
func (u *Update) buildPath(taskId string) string {
	return u.prefixPath + "/tasks/" + taskId
}
