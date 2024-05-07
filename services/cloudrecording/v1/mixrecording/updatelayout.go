package mixrecording

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type UpdateLayout struct {
	Base *baseV1.UpdateLayout
}

var _ baseV1.UpdateLayoutMixRecording = (*UpdateLayout)(nil)

func (u *UpdateLayout) WithForwardRegion(prefix core.ForwardedReginPrefix) baseV1.UpdateLayoutMixRecording {
	u.Base.WithForwardRegion(prefix)

	return u
}

func (u *UpdateLayout) Do(ctx context.Context, resourceID string, sid string, cname string, uid string,
	clientRequest *baseV1.UpdateLayoutUpdateMixRecordingClientRequest,
) (*baseV1.UpdateLayoutResp, error) {
	return u.Base.Do(ctx, resourceID, sid, baseV1.MixMode, &baseV1.UpdateLayoutReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.UpdateLayoutClientRequest{
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
