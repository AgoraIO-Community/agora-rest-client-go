package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/retry"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type List struct {
	baseHandler
}

func NewList(module string, logger log.Logger, client client.Client, prefixPath string) *List {
	return &List{
		baseHandler: baseHandler{
			module:     module,
			logger:     logger,
			client:     client,
			prefixPath: prefixPath,
		},
	}
}

// buildPath returns the request path.
// /api/conversational-ai-agent/v2/projects/{appid}/agents?limit=10&state=2&from_time=1733013296&to_time=1734016896
func (q *List) buildPath(queryFields map[string]any) string {
	return q.prefixPath + "/agents?" + buildQuery(queryFields)
}

func buildQuery(queryFields map[string]any) string {
	urlValues := make(url.Values)
	for key, value := range queryFields {
		urlValues[key] = []string{fmt.Sprintf("%v", value)}
	}

	return urlValues.Encode()
}

// buildQueryFields converts ListOptions to a query field map
func buildQueryFields(options ...req.ListOption) map[string]any {
	opts := req.ListOptions{}

	for _, option := range options {
		option(&opts)
	}

	queryFields := make(map[string]any)

	if opts.Limit != nil {
		queryFields["limit"] = *opts.Limit
	}

	if opts.State != nil {
		queryFields["state"] = *opts.State
	}

	if opts.FromTime != nil {
		queryFields["from_time"] = *opts.FromTime
	}

	if opts.ToTime != nil {
		queryFields["to_time"] = *opts.ToTime
	}

	if opts.Cursor != nil {
		queryFields["cursor"] = *opts.Cursor
	}

	if opts.Channel != nil {
		queryFields["channel"] = *opts.Channel
	}

	return queryFields
}

func (q *List) Do(ctx context.Context, options ...req.ListOption) (*resp.ListResp, error) {
	queryFields := buildQueryFields(options...)
	path := q.buildPath(queryFields)
	responseData, err := q.doRESTWithRetry(ctx, path, http.MethodGet, nil)
	if err != nil {
		var internalErr *agora.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	var listResp resp.ListResp

	listResp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse resp.ListSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, err
		}
		listResp.SuccessRes = successResponse
	} else {
		var errResponse resp.ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, err
		}
		listResp.ErrResponse = errResponse
	}

	return &listResp, nil
}

const listRetryCount = 3

func (q *List) doRESTWithRetry(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error) {
	var (
		baseResponse *agora.BaseResponse
		err          error
		retryCount   int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		baseResponse, doErr = q.client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := baseResponse.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			q.logger.Debugf(ctx, q.module, "http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)),
			)
		default:
			q.logger.Debugf(ctx, q.module, "http status code is %d, retry,http response:%s", statusCode, baseResponse.RawBody)
			return fmt.Errorf("http status code is %d, retry", baseResponse.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retryCount >= listRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		q.logger.Debugf(ctx, q.module, "http request err:%s", err)
		retryCount++
	})

	return baseResponse, err
}
