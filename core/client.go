package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client interface {
	GetAppID() string
	DoREST(ctx context.Context, path string, method string, requestBody interface{}) (*BaseResponse, error)
}

type Config struct {
	AppID      string
	Timeout    time.Duration
	Credential Credential

	RegionCode RegionArea
	Logger     Logger
}

type ClientImpl struct {
	appID      string
	httpClient *http.Client
	timeout    time.Duration
	logger     Logger
	credential Credential

	module     string
	domainPool *DomainPool
}

var _ Client = (*ClientImpl)(nil)

const defaultTimeout = 10 * time.Second

func NewClient(config *Config) *ClientImpl {
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	cc := &http.Client{
		Timeout: config.Timeout,
	}
	if config.Logger == nil {
		config.Logger = defaultLogger
	}
	return &ClientImpl{
		appID:      config.AppID,
		credential: config.Credential,
		httpClient: cc,
		timeout:    config.Timeout,
		logger:     config.Logger,
		module:     "http client",
		domainPool: NewDomainPool(config.RegionCode, config.Logger),
	}
}

func (c *ClientImpl) marshalBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	c.logger.Debugf(context.Background(), c.module, "http request body:%s", jsonBody)
	return bytes.NewReader(jsonBody), nil
}

func (c *ClientImpl) DoREST(ctx context.Context, path string,
	method string, requestBody interface{},
) (*BaseResponse, error) {
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, c.timeout)
	defer func() {
		_ = cancelFunc
	}()

	var (
		resp  *BaseResponse
		err   error
		retry int
	)

	err = RetryDo(func(retryCount int) error {
		var doErr error

		resp, doErr = c.doREST(timeoutCtx, path, method, requestBody)
		if doErr != nil {
			return NewRetryErr(false, doErr)
		}

		statusCode := resp.HttpStatusCode
		switch {
		case statusCode >= 200 && statusCode <= 300:
			return nil
		case statusCode >= 400 && statusCode < 410:
			c.logger.Debugf(ctx, c.module, "http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)
			return &RetryErr{
				false,
				NewInternalErr(fmt.Sprintf("http status code is %d, no retry,http response:%s", statusCode, resp.RawBody)),
			}
		default:
			c.logger.Debugf(ctx, c.module, "http status code is %d, retry,http response:%s", statusCode, resp.RawBody)
			return fmt.Errorf("http status code is %d, retry", resp.RawBody)
		}
	}, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
		}
		return retry >= 3
	}, func(i int) time.Duration {
		return time.Second * time.Duration(i+1)
	}, func(err error) {
		c.logger.Debugf(ctx, c.module, "http request err:%s", err)
		retry++
	})
	if resp != nil {
		c.logger.Debugf(ctx, c.module, "http response:%s", resp.RawBody)
	}
	return resp, err
}

func (c *ClientImpl) doREST(ctx context.Context, path string,
	method string, requestBody interface{},
) (*BaseResponse, error) {
	var (
		err  error
		resp *http.Response
		req  *http.Request
	)

	if err = c.domainPool.SelectBestDomain(ctx); err != nil {
		return nil, err
	}

	doHttpRequest := func() error {
		req, err = c.createRequest(ctx, path, method, requestBody)
		if err != nil {
			return err
		}
		resp, err = c.httpClient.Do(req)
		return err
	}
	err = RetryDo(
		func(retryCount int) error {
			c.logger.Debugf(ctx, c.module, "http retry attempt:%d", retryCount)
			return doHttpRequest()
		},
		func() bool {
			select {
			case <-ctx.Done():
				return true
			default:
			}
			return false
		},
		func(i int) time.Duration {
			if i == 0 {
				return 0 * time.Second
			}
			return 500 * time.Millisecond
		},
		func(err error) {
			c.logger.Debugf(ctx, c.module, "http request err:%s", err)
			c.domainPool.NextRegion()
		},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &BaseResponse{
		RawBody:        body,
		HttpStatusCode: resp.StatusCode,
	}, nil
}

func (c *ClientImpl) addCredential(req *http.Request) {
	if c.credential != nil {
		c.credential.Visit(req)
	}
}

func (c *ClientImpl) GetAppID() string {
	return c.appID
}

func (c *ClientImpl) createRequest(ctx context.Context, path string,
	method string, requestBody interface{},
) (*http.Request, error) {
	url := c.domainPool.GetCurrentUrl() + path

	c.logger.Debugf(ctx, "create request url:%s", url)
	jsonBody, err := c.marshalBody(requestBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, jsonBody)
	if err != nil {
		return nil, err
	}
	c.addCredential(req)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	return req, nil
}
