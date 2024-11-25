package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
)

type UpdateLayout struct {
	Base *api.UpdateLayout
}

func NewUpdateLayout(base *api.UpdateLayout) *UpdateLayout {
	return &UpdateLayout{Base: base}
}

func (u *UpdateLayout) Do(ctx context.Context, resourceID string, sid string, cname string, uid string,
	clientRequest *UpdateLayoutUpdateMixRecordingClientRequest,
) (*api.UpdateLayoutResp, error) {
	return u.Base.Do(ctx, resourceID, sid, api.MixMode, &api.UpdateLayoutReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateLayoutClientRequest{
			MaxResolutionUID:           clientRequest.MaxResolutionUID,
			MixedVideoLayout:           clientRequest.MixedVideoLayout,
			BackgroundColor:            clientRequest.BackgroundColor,
			BackgroundImage:            clientRequest.BackgroundImage,
			DefaultUserBackgroundImage: clientRequest.DefaultUserBackgroundImage,
			LayoutConfig:               clientRequest.LayoutConfig,
			BackgroundConfig:           clientRequest.BackgroundConfig,
		},
	})
}
