package cloudrecording

import (
	"context"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	agoraClient "github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/scenario"
)

const projectName = "cloud_recording"

type Client struct {
	acquireAPI      *api.Acquire
	startAPI        *api.Start
	stopAPI         *api.Stop
	queryAPI        *api.Query
	updateLayoutAPI *api.UpdateLayout
	updateAPI       *api.Update

	individualRecordingScenario *scenario.IndividualRecording
	webRecordingScenario        *scenario.WebRecording
	mixRecordingScenario        *scenario.MixRecording
}

// @brief Defines the configuration for the Cloud Recording client
//
// @since v0.8.0
type Config struct {
	// Agora AppID
	AppID string
	// Timeout for HTTP requests
	HttpTimeout time.Duration
	// Credential for accessing the Agora service.
	//
	// Available credential types:
	//
	//  - BasicAuthCredential: See auth.NewBasicAuthCredential for details
	Credential auth.Credential

	// Domain area for the REST Client. See domain.Area for details.
	DomainArea domain.Area

	// Logger for the REST Client
	//
	// Implement the log.Logger interface in your project to output REST Client logs to your logging component.
	//
	// Alternatively, you can use the default logging component. See log.NewDefaultLogger for details.
	Logger log.Logger
}

var RetryCount = 3

// NewClient
//
// @brief Creates a Cloud Recording client with the specified configuration
//
// @param config Configuration of the Cloud Recording client instance. See Config for details.
//
// @return Returns the Cloud Recording client instance.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
//
// @since v0.8.0
func NewClient(config *Config) (*Client, error) {
	prefixPath := "/v1/apps/" + config.AppID + "/" + projectName

	agoraClient, err := agoraClient.New(&agora.Config{
		AppID:       config.AppID,
		HttpTimeout: config.HttpTimeout,
		Credential:  config.Credential,
		DomainArea:  config.DomainArea,
		Logger:      config.Logger,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{
		acquireAPI:      api.NewAcquire("cloudRecording:acquire", config.Logger, RetryCount, agoraClient, prefixPath),
		startAPI:        api.NewStart("cloudRecording:start", config.Logger, RetryCount, agoraClient, prefixPath),
		stopAPI:         api.NewStop("cloudRecording:stop", config.Logger, RetryCount, agoraClient, prefixPath),
		queryAPI:        api.NewQuery("cloudRecording:query", config.Logger, RetryCount, agoraClient, prefixPath),
		updateLayoutAPI: api.NewUpdateLayout("cloudRecording:updateLayout", config.Logger, RetryCount, agoraClient, prefixPath),
		updateAPI:       api.NewUpdate("cloudRecording:update", config.Logger, RetryCount, agoraClient, prefixPath),
	}

	c.individualRecordingScenario = scenario.NewIndividualRecording(c.acquireAPI, c.startAPI, c.stopAPI, c.queryAPI, c.updateAPI)
	c.webRecordingScenario = scenario.NewWebRecording(c.acquireAPI, c.startAPI, c.stopAPI, c.queryAPI, c.updateAPI)
	c.mixRecordingScenario = scenario.NewMixRecording(c.acquireAPI, c.startAPI, c.stopAPI, c.queryAPI, c.updateLayoutAPI, c.updateAPI)

	return c, nil
}

func (c *Client) Acquire(ctx context.Context, payload *api.AcquireReqBody) (*api.AcquireResp, error) {
	return c.acquireAPI.Do(ctx, payload)
}

func (c *Client) Start(ctx context.Context, resourceID string, mode string, payload *api.StartReqBody) (*api.StartResp, error) {
	return c.startAPI.Do(ctx, resourceID, mode, payload)
}

func (c *Client) Stop(ctx context.Context, resourceID string, sid string, mode string, payload *api.StopReqBody) (*api.StopResp, error) {
	return c.stopAPI.Do(ctx, resourceID, sid, mode, payload)
}

func (c *Client) Query(ctx context.Context, resourceID string, sid string, mode string) (*api.QueryResp, error) {
	return c.queryAPI.Do(ctx, resourceID, sid, mode)
}

func (c *Client) Update(ctx context.Context, resourceID string, sid string, mode string, payload *api.UpdateReqBody) (*api.UpdateResp, error) {
	return c.updateAPI.Do(ctx, resourceID, sid, mode, payload)
}

func (c *Client) UpdateLayout(ctx context.Context, resourceID string, sid string, mode string, payload *api.UpdateLayoutReqBody) (*api.UpdateLayoutResp, error) {
	return c.updateLayoutAPI.Do(ctx, resourceID, sid, mode, payload)
}

// @brief Returns the individual recording scenario instance.
//
// @return Returns the individual recording scenario instance. See scenario.IndividualRecording for details.
//
// @since v0.8.0
func (c *Client) IndividualRecording() *scenario.IndividualRecording {
	return c.individualRecordingScenario
}

// @brief Returns the web recording scenario instance.
//
// @return Returns the web recording scenario instance. See scenario.WebRecording for details.
//
// @since v0.8.0
func (c *Client) WebRecording() *scenario.WebRecording {
	return c.webRecordingScenario
}

// @brief Returns the mix recording scenario instance.
//
// @return Returns the mix recording scenario instance. See scenario.MixRecording for details.
//
// @since v0.8.0
func (c *Client) MixRecording() *scenario.MixRecording {
	return c.mixRecordingScenario
}
