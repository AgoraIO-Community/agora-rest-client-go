package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

type List struct {
	baseHandler
}

func NewList(module string, logger log.Logger, retryCount int, client client.Client, prefixPath string) *List {
	return &List{
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
// /api/conversational-ai-agent/v2/projects/{appid}/agents?limit=10&state=2&from_time=1733013296&to_time=1734016896
func (l *List) buildPath(queryFields map[string]any) string {
	return l.prefixPath + "/agents?" + buildQuery(queryFields)
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

func (l *List) Do(ctx context.Context, options ...req.ListOption) (*resp.ListResp, error) {
	queryFields := buildQueryFields(options...)
	path := l.buildPath(queryFields)
	responseData, err := doRESTWithRetry(ctx, l.module, l.logger, l.retryCount, l.client, path, http.MethodGet, nil)
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
