package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/retry"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/utils"
)

type Client interface {
	GetAppID() string
	GetLogger() log.Logger
	DoREST(ctx context.Context, path string, method string, requestBody interface{}) (*agora.BaseResponse, error)
}

type Impl struct {
	appID      string
	httpClient *http.Client
	timeout    time.Duration
	logger     log.Logger
	credential auth.Credential

	module     string
	domainPool *domain.Pool
}

func (c *Impl) GetLogger() log.Logger {
	return c.logger
}

var _ Client = (*Impl)(nil)

const defaultHttpTimeout = 10 * time.Second

func New(config *agora.Config) (*Impl, error) {
	if config.HttpTimeout == 0 {
		config.HttpTimeout = defaultHttpTimeout
	}
	cc := &http.Client{
		Timeout: config.HttpTimeout,
	}
	if config.Logger == nil {
		config.Logger = log.DefaultLogger
	}

	domainPool, err := domain.NewPool(config.DomainArea, config.Logger)
	if err != nil {
		return nil, err
	}

	return &Impl{
		appID:      config.AppID,
		credential: config.Credential,
		httpClient: cc,
		timeout:    config.HttpTimeout,
		logger:     config.Logger,
		module:     "http client",
		domainPool: domainPool,
	}, nil
}

func (c *Impl) marshalBody(body interface{}) (io.Reader, error) {
	if utils.IsNil(body) {
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

func (c *Impl) DoREST(ctx context.Context, path string,
	method string, requestBody interface{},
) (*agora.BaseResponse, error) {
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

		req.Header.Add("User-Agent", agora.BuildUserAgent())
		resp, err = c.httpClient.Do(req)
		return err
	}
	err = retry.Do(
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

	return &agora.BaseResponse{
		RawResponse:    resp,
		RawBody:        body,
		HttpStatusCode: resp.StatusCode,
	}, nil
}

func (c *Impl) addCredential(req *http.Request) {
	if c.credential != nil {
		c.credential.SetAuth(req)
	}
}

func (c *Impl) GetAppID() string {
	return c.appID
}

func (c *Impl) createRequest(ctx context.Context, path string,
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
