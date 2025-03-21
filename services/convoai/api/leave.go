package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/retry"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type Leave struct {
	baseHandler
}

// NewLeave Creates a new Leave instance
func NewLeave(module string, logger log.Logger, client client.Client, prefixPath string) *Leave {
	return &Leave{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /api/conversational-ai-agent/v2/projects/{appid}/agents/{agentId}/leave
func (d *Leave) buildPath(agentId string) string {
	return d.prefixPath + "/agents/" + agentId + "/leave"
}

func (d *Leave) Do(ctx context.Context, agentId string) (*resp.LeaveResp, error) {
	path := d.buildPath(agentId)
	responseData, err := d.doRESTWithRetry(ctx, path, http.MethodPost, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var response resp.LeaveResp

	response.BaseResponse = responseData

	if responseData.HttpStatusCode != http.StatusOK {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		response.ErrResponse = errResponse
	}

	return &response, nil
}

const leaveRetryCount = 3

func (d *Leave) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error) {
	var (
		response   *agora.BaseResponse
		err        error
		retryCount int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		response, doErr = d.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := response.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			d.logger.Debugf(ctx, d.module, "http status code is %d, no retry,http response:%s", statusCode, response.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, response.RawBody)),
			)
		default:
			d.logger.Debugf(ctx, d.module, "http status code is %d, retry,http response:%s", statusCode, response.RawBody)
			return fmt.Errorf("http status code is %d, retry", response.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retryCount >= leaveRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		d.logger.Debugf(ctx, d.module, "http request err:%s", err)
		retryCount++
	})

	return response, err
}
