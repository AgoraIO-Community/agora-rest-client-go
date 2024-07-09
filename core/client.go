package core

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client interface {
	GetAppID() string
	GetLogger() Logger
	DoREST(ctx context.Context, path string, method string, requestBody interface{}) (*BaseResponse, error)
}

type Config struct {
	AppID       string
	HttpTimeout time.Duration
	Credential  Credential

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

func (c *ClientImpl) GetLogger() Logger {
	return c.logger
}

var _ Client = (*ClientImpl)(nil)

const defaultHttpTimeout = 10 * time.Second

func NewClient(config *Config) *ClientImpl {
	if config.HttpTimeout == 0 {
		config.HttpTimeout = defaultHttpTimeout
	}
	cc := &http.Client{
		Timeout: config.HttpTimeout,
	}
	if config.Logger == nil {
		config.Logger = defaultLogger
	}
	return &ClientImpl{
		appID:      config.AppID,
		credential: config.Credential,
		httpClient: cc,
		timeout:    config.HttpTimeout,
		logger:     config.Logger,
		module:     "http client",
		domainPool: NewDomainPool(config.RegionCode, config.Logger),
	}
}

func (c *ClientImpl) marshalBody(body interface{}) (io.Reader, error) {
	if IsNil(body) {
		return nil, nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	c.logger.Debugf(context.Background(), c.module, "http request body:%s", jsonBody)

	return bytes.NewReader(jsonBody), nil
}

const doRESTTimeout = 10 * time.Second

func (c *ClientImpl) DoREST(ctx context.Context, path string,
	method string, requestBody interface{},
) (*BaseResponse, error) {
	var (
		err  error
		resp *http.Response
		req  *http.Request
	)
	timeoutCtx, cancel := context.WithTimeout(ctx, doRESTTimeout)
	defer cancel()

	if err = c.domainPool.SelectBestDomain(timeoutCtx); err != nil {
		return nil, err
	}

	doHttpRequest := func() error {
		req, err = c.createRequest(timeoutCtx, path, method, requestBody)
		if err != nil {
			return err
		}

		req.Header.Add("User-Agent", BuildUserAgent())
		resp, err = c.httpClient.Do(req)
		return err
	}
	err = RetryDo(
		func(retryCount int) error {
			c.logger.Debugf(timeoutCtx, c.module, "http retry attempt:%d", retryCount)
			return doHttpRequest()
		},
		func() bool {
			select {
			case <-timeoutCtx.Done():
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
			c.logger.Debugf(timeoutCtx, c.module, "http request err:%s", err)
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
	c.logger.Debugf(ctx, c.module, "http response:%s", body)

	return &BaseResponse{
		RawResponse:    resp,
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
