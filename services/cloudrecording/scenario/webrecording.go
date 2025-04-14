package scenario

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/resp"
)

type WebRecording struct {
	acquireAPI *api.Acquire
	startAPI   *api.Start
	stopAPI    *api.Stop
	queryAPI   *api.Query
	updateAPI  *api.Update
}

func NewWebRecording(
	acquireAPI *api.Acquire,
	startAPI *api.Start,
	stopAPI *api.Stop,
	queryAPI *api.Query,
	updateAPI *api.Update,
) *WebRecording {
	return &WebRecording{
		acquireAPI: acquireAPI,
		startAPI:   startAPI,
		stopAPI:    stopAPI,
		queryAPI:   queryAPI,
		updateAPI:  updateAPI,
	}
}

// @brief Get a resource ID for web recording.
//
// @since v0.8.0
//
// @post After receiving the resource ID, call the Start API to start cloud recording.
//
// @param ctx Context to control the request lifecycle.
//
// @param cname The name of the channel to be recorded.
//
// @param uid The user ID used by the cloud recording service in the RTC channel to identify the recording service in the channel.
//
// @param clientRequest The request body. See req.AcquireWebRecodingClientRequest for details.
//
// @return Returns the response *AcquireResp. See api.AcquireResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (w *WebRecording) Acquire(ctx context.Context, cname string, uid string, clientRequest *req.AcquireWebRecodingClientRequest) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			RecordingFileConfig:    clientRequest.StartParameter.RecordingFileConfig,
			StorageConfig:          clientRequest.StartParameter.StorageConfig,
			ExtensionServiceConfig: clientRequest.StartParameter.ExtensionServiceConfig,
		}
	}

	return w.acquireAPI.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               1,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}

// @brief Start web recording.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param cname Channel name.
//
// @param uid User ID.
//
// @param clientRequest The request body. See req.StartWebRecordingClientRequest for details.
//
// @return Returns the response *StartResp. See api.StartResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (w *WebRecording) Start(ctx context.Context, resourceID string, cname string, uid string, clientRequest *req.StartWebRecordingClientRequest) (*api.StartResp, error) {
	return w.startAPI.Do(ctx, resourceID, api.WebMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			RecordingFileConfig:    clientRequest.RecordingFileConfig,
			StorageConfig:          clientRequest.StorageConfig,
			ExtensionServiceConfig: clientRequest.ExtensionServiceConfig,
		},
	})
}

// @brief Query the status of web recording.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @return Returns the response *QueryWebRecordingResp. See resp.QueryWebRecordingResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (w *WebRecording) Query(ctx context.Context, resourceID string, sid string) (*resp.QueryWebRecordingResp, error) {
	respData, err := w.queryAPI.Do(ctx, resourceID, sid, api.WebMode)
	if err != nil {
		return nil, err
	}

	var webResp resp.QueryWebRecordingResp

	webResp.Response = respData.Response
	if respData.IsSuccess() {
		successResp := respData.SuccessResponse
		webResp.SuccessResponse = resp.QueryWebRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetWebRecording2CDNServerResponse(),
		}
	}

	return &webResp, nil
}

// @brief Query the status of pushing web page recording to the CDN.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @return Returns the response *QueryRtmpPublishResp. See resp.QueryRtmpPublishResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (w *WebRecording) QueryRtmpPublish(ctx context.Context, resourceID string, sid string) (*resp.QueryRtmpPublishResp, error) {
	respData, err := w.queryAPI.Do(ctx, resourceID, sid, api.WebMode)
	if err != nil {
		return nil, err
	}

	var webResp resp.QueryRtmpPublishResp

	webResp.Response = respData.Response
	if respData.IsSuccess() {
		successResp := respData.SuccessResponse
		webResp.SuccessResponse = resp.QueryRtmpPublishSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetRtmpPublishServiceServerResponse(),
		}
	}

	return &webResp, nil
}

// @brief Update web recording configuration.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @param cname The name of the channel to be recorded.
//
// @param uid The user ID used by the cloud recording service in the RTC channel to identify the recording service in the channel.
//
// @param clientRequest The request body. See req.UpdateWebRecordingClientRequest for details.
//
// @return Returns the response *UpdateResp. See api.UpdateResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (w *WebRecording) Update(ctx context.Context, resourceID string, sid string, cname string, uid string, clientRequest *req.UpdateWebRecordingClientRequest) (*api.UpdateResp, error) {
	return w.updateAPI.Do(ctx, resourceID, sid, api.WebMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			WebRecordingConfig: clientRequest.WebRecordingConfig,
			RtmpPublishConfig:  clientRequest.RtmpPublishConfig,
		},
	})
}

// @brief Stop web recording.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @param cname The name of the channel to be recorded.
//
// @param uid The user ID used by the cloud recording service in the RTC channel to identify the recording service in the channel.
//
// @param asyncStop Whether to stop the recording asynchronously.
//   - true: Stop the recording asynchronously.
//   - false: Stop the recording synchronously.
//
// @return Returns the response *StopResp. See api.StopResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (w *WebRecording) Stop(ctx context.Context, resourceID string, sid string, cname string, uid string, asyncStop bool) (*api.StopResp, error) {
	return w.stopAPI.Do(ctx, resourceID, sid, api.WebMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
}
