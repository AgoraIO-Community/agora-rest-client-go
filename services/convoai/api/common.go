package api

import (
	"context"
	"fmt"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/retry"
)

func doRESTWithRetry(ctx context.Context, module string, logger log.Logger, handleRetryCount int, client client.Client, path string, method string, requestBody any) (*agora.BaseResponse, error) {
	var (
		baseResponse *agora.BaseResponse
		err          error
		retryCount   int
	)

	err = retry.Do(func(retryCount int) error {
		var doErr error

		baseResponse, doErr = client.DoREST(ctx, path, method, requestBody)
		if doErr != nil {
			return agora.NewRetryErr(false, doErr)
		}

		statusCode := baseResponse.HttpStatusCode
		switch {
		case statusCode == 200 || statusCode == 201:
			return nil
		case statusCode >= 400 && statusCode < 499:
			logger.Debugf(ctx, module, "http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)
			return agora.NewRetryErr(
				false,
				agora.NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, baseResponse.RawBody)),
			)
		default:
			logger.Debugf(ctx, module, "http status code is %d, retry,http response:%s", statusCode, baseResponse.RawBody)
			return fmt.Errorf("http status code is %d, retry", baseResponse.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retryCount >= handleRetryCount
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		logger.Debugf(ctx, module, "http request err:%s", err)
		retryCount++
	})

	return baseResponse, err
}
