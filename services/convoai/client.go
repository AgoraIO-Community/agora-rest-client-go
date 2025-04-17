package convoai

import (
	"context"
	"errors"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/agora"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/auth"
	agoraClient "github.com/AgoraIO-Community/agora-rest-client-go/agora/client"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/domain"
	"github.com/AgoraIO-Community/agora-rest-client-go/agora/log"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/convoai/resp"
)

const projectName = "conversational-ai-agent"

type Client struct {
	joinAPI      *api.Join
	leaveAPI     *api.Leave
	listAPI      *api.List
	queryAPI     *api.Query
	updateAPI    *api.Update
	interruptAPI *api.Interrupt
	historyAPI   *api.History
	speakAPI     *api.Speak
}

// @brief ServiceRegion represents the region of the Conversational AI engine service
//
// @note The service in Chinese mainland and the global region are two different services
//
// @since v0.7.0
type ServiceRegion int

const (
	UnknownServiceRegion ServiceRegion = iota
	// ChineseMainlandServiceRegion represents the Conversational AI engine service in Chinese mainland
	ChineseMainlandServiceRegion
	// GlobalServiceRegion represents the Conversational AI engine service in the global region, except Chinese mainland
	GlobalServiceRegion
)

// @brief Defines the configuration for the Conversational AI engine client
//
// @since v0.7.0
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

	// Service version. See ServiceRegion for details.
	ServiceRegion ServiceRegion
}

// NewClient
//
// @brief Creates a Conversational AI engine client with the specified configuration
//
// @param config Configuration of the Conversational AI engine client. See Config for details.
//
// @return Returns the Conversational AI engine client.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
//
// @since v0.7.0
func NewClient(config *Config) (*Client, error) {
	var prefixPath string

	switch config.ServiceRegion {
	case ChineseMainlandServiceRegion:
		prefixPath = "/cn/api/" + projectName + "/v2/projects/" + config.AppID
	case GlobalServiceRegion:
		prefixPath = "/api/" + projectName + "/v2/projects/" + config.AppID
	default:
		return nil, errors.New("ServiceRegion should not be Unknown")
	}

	c, err := agoraClient.New(&agora.Config{
		AppID:       config.AppID,
		HttpTimeout: config.HttpTimeout,
		Credential:  config.Credential,
		DomainArea:  config.DomainArea,
		Logger:      config.Logger,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		joinAPI:      api.NewJoin("convoai:join", config.Logger, c, prefixPath),
		leaveAPI:     api.NewLeave("convoai:leave", config.Logger, c, prefixPath),
		listAPI:      api.NewList("convoai:list", config.Logger, c, prefixPath),
		queryAPI:     api.NewQuery("convoai:query", config.Logger, c, prefixPath),
		updateAPI:    api.NewUpdate("convoai:update", config.Logger, c, prefixPath),
		interruptAPI: api.NewInterrupt("convoai:interrupt", config.Logger, c, prefixPath),
		historyAPI:   api.NewHistory("convoai:history", config.Logger, c, prefixPath),
		speakAPI:     api.NewSpeak("convoai:speak", config.Logger, c, prefixPath),
	}, nil
}

// Join
// @brief Creates an agent instance and joins the specified RTC channel
//
// @since v0.7.0
//
// @example Use this to create an agent instance in an RTC channel.
//
// @post After successful execution, the agent will join the specified channel. You can perform subsequent operations using the returned agent ID.
//
// @param ctx Context to control the request lifecycle.
//
// @param name Unique identifier for the agent. The same identifier cannot be used repeatedly.
//
// @param propertiesBody Configuration properties of the agent, including channel information, token, LLM settings, TTS settings, etc. See api.JoinPropertiesReqBody for details.
//
// @return Returns the response *JoinResp. See api.JoinResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) Join(ctx context.Context, name string, payload *req.JoinPropertiesReqBody) (*resp.JoinResp, error) {
	return c.joinAPI.Do(ctx, name, payload)
}

// Leave
//
// @brief Stops the specified agent instance and leaves the RTC channel
//
// @since v0.7.0
//
// @example Use this to stop an agent instance.
//
// @post After successful execution, the agent will be stopped and leave the RTC channel
//
// @note Ensure the agent ID has been obtained by calling the Join API before using this method.
//
// @param ctx Context to control the request lifecycle.
//
// @param agentId Agent ID.
//
// @return Returns the response *LeaveResp. See api.LeaveResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) Leave(ctx context.Context, agentId string) (*resp.LeaveResp, error) {
	return c.leaveAPI.Do(ctx, agentId)
}

// Query
//
// @brief Query the current status of the specified agent instance
//
// @since v0.7.0
//
// @example Use this to get the current status of the specified agent instance.
//
// @post After successful execution, the current status of the specified agent instance will be retrieved.
//
// @note Ensure the agent ID has been obtained by calling the Join API before using this method.
//
// @param ctx Context to control the request lifecycle.
//
// @param agentId Agent ID.
//
// @return Returns the response *QueryResp. See api.QueryResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) Query(ctx context.Context, agentId string) (*resp.QueryResp, error) {
	return c.queryAPI.Do(ctx, agentId)
}

// List
// @brief Retrieves a list of agents that meet the specified criteria
//
// @since v0.7.0
//
// @example Use this to get a list of agents that meet the specified criteria.
//
// @post After successful execution, a list of agents that meet the specified criteria will be retrieved.
//
// @param ctx Context to control the request lifecycle.
//
// @param ListOption Query parameters. See api.ListOption for details.
//
// @return Returns the response *ListResp. See api.ListResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) List(ctx context.Context, options ...req.ListOption) (*resp.ListResp, error) {
	return c.listAPI.Do(ctx, options...)
}

// Update
//
// @brief Adjusts the agent's parameters at runtime
//
// @since v0.7.0
//
// @example Use this to adjust the agent's parameters at runtime.
//
// @post After successful execution, the agent's parameters will be adjusted.
//
// @note Ensure the agent ID has been obtained by calling the Join API before using this method.
//
// @param ctx Context to control the request lifecycle.
//
// @param agentId Agent ID.
//
// @param payload Parameters to be adjusted. See api.UpdateReqBody for details.
//
// @return Returns the response *UpdateResp. See api.UpdateResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) Update(ctx context.Context, agentId string, payload *req.UpdateReqBody) (*resp.UpdateResp, error) {
	return c.updateAPI.Do(ctx, agentId, payload)
}

// Interrupt
//
// @brief Interrupts the specified agent instance
//
// @since v0.9.0
//
// @example Use this method to interrupt the specified agent instance.
//
// @param ctx Context to control the request lifecycle.
//
// @param agentId Agent ID.
//
// @return Returns the response *InterruptResp. See api.InterruptResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) Interrupt(ctx context.Context, agentId string) (*resp.InterruptResp, error) {
	return c.interruptAPI.Do(ctx, agentId)
}

// GetHistory
//
// @brief Acquires the short-term memory of the specified agent instance
//
// @since v0.9.0
//
// @example Use this method to acquire the short-term memory of the specified agent instance.
//
// @param ctx Context to control the request lifecycle.
//
// @param agentId Agent ID.
//
// @return Returns the response *HistoryResp. See api.HistoryResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) GetHistory(ctx context.Context, agentId string) (*resp.HistoryResp, error) {
	return c.historyAPI.Do(ctx, agentId)
}

// Speak
//
// @brief Speaks a custom message for the specified agent instance
//
// @since v0.9.0
//
// @param ctx Context to control the request lifecycle.
//
// @param agentId Agent ID.
//
// @param payload Request body for the specified agent to speak a custom message. See api.SpeakBody for details.
//
// @return Returns the response *SpeakResp. See api.SpeakResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (c *Client) Speak(ctx context.Context, agentId string, payload *req.SpeakBody) (*resp.SpeakResp, error) {
	return c.speakAPI.Do(ctx, agentId, payload)
}
