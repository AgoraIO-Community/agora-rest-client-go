package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type Speak struct {
	baseHandler
}

func NewSpeak(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *Speak {
	return &Speak{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			retryCount: retryCount,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
//
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}/speak
func (s *Speak) buildPath(agentId string) string {
	return s.prefixPath + "/agents/" + agentId + "/speak"
}

func (s *Speak) Do(ctx context.Context, agentId string, body *req.SpeakBody) (*resp.SpeakResp, error) {
	path := s.buildPath(agentId)
	responseData, err := doRESTWithRetry(ctx, s.module, s.logger, s.retryCount, s.client, path, http.MethodPost, body)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var speakResp resp.SpeakResp

	speakResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.SpeakSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		speakResp.SuccessRes = successResponse
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		speakResp.ErrResponse = errResponse
	}

	return &speakResp, nil
}
