package scenario

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/resp"
)

type IndividualRecording struct {
	acquireAPI *api.Acquire
	startAPI   *api.Start
	stopAPI    *api.Stop
	queryAPI   *api.Query
	updateAPI  *api.Update
}

func NewIndividualRecording(
	acquireAPI *api.Acquire,
	startAPI *api.Start,
	stopAPI *api.Stop,
	queryAPI *api.Query,
	updateAPI *api.Update,
) *IndividualRecording {
	return &IndividualRecording{
		acquireAPI: acquireAPI,
		startAPI:   startAPI,
		stopAPI:    stopAPI,
		queryAPI:   queryAPI,
		updateAPI:  updateAPI,
	}
}

// @brief Get a resource ID for individual cloud recording.
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
// @param clientRequest The request body.
//
// @return Returns the response *AcquireResp. See api.AcquireResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (i *IndividualRecording) Acquire(ctx context.Context, cname string, uid string,
	clientRequest *req.AcquireIndividualRecordingClientRequest,
) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			SnapshotConfig:      clientRequest.StartParameter.SnapshotConfig,
		}
	}

	return i.acquireAPI.Do(ctx, &api.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.AcquireClientRequest{
			Scene:               0,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}

// @brief Start individual cloud recording.
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
// @param clientRequest The request body. See req.StartIndividualRecordingClientRequest for details.
//
// @return Returns the response *StartResp. See api.StartResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (i *IndividualRecording) Start(ctx context.Context, resourceId string, cname string, uid string,
	clientRequest *req.StartIndividualRecordingClientRequest,
) (*api.StartResp, error) {
	return i.startAPI.Do(ctx, resourceId, api.IndividualMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			Token:               clientRequest.Token,
			RecordingConfig:     clientRequest.RecordingConfig,
			RecordingFileConfig: clientRequest.RecordingFileConfig,
			SnapshotConfig:      clientRequest.SnapshotConfig,
			StorageConfig:       clientRequest.StorageConfig,
		},
	})
}

// @brief Query the status of individual cloud recording when video screenshot capture is turned off.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @return Returns the response *QueryIndividualRecordingResp. See resp.QueryIndividualRecordingResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (i *IndividualRecording) Query(ctx context.Context, resourceId string, sid string) (*resp.QueryIndividualRecordingResp, error) {
	respData, err := i.queryAPI.Do(ctx, resourceId, sid, api.IndividualMode)
	if err != nil {
		return nil, err
	}

	var individualResp resp.QueryIndividualRecordingResp

	individualResp.Response = respData.Response
	if respData.IsSuccess() {
		successResp := respData.SuccessResponse
		individualResp.SuccessRes = resp.QueryIndividualRecordingSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualRecordingServerResponse(),
		}
	}

	return &individualResp, nil
}

// @brief Query the status of the individual cloud recording when the video screenshot capture is turned on.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @return Returns the response *QueryIndividualRecordingVideoScreenshotResp. See resp.QueryIndividualRecordingVideoScreenshotResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (i *IndividualRecording) QueryVideoScreenshot(ctx context.Context, resourceId string, sid string) (*resp.QueryIndividualRecordingVideoScreenshotResp, error) {
	respData, err := i.queryAPI.Do(ctx, resourceId, sid, api.IndividualMode)
	if err != nil {
		return nil, err
	}

	var individualResp resp.QueryIndividualRecordingVideoScreenshotResp

	individualResp.Response = respData.Response
	if respData.IsSuccess() {
		successResp := respData.SuccessResponse
		individualResp.SuccessRes = resp.QueryIndividualRecordingVideoScreenshotSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetIndividualVideoScreenshotServerResponse(),
		}
	}

	return &individualResp, nil
}

// @brief Update the individual cloud recording configuration.
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
// @param clientRequest The request body. See req.UpdateIndividualRecordingClientRequest for details.
//
// @return Returns the response *UpdateResp. See api.UpdateResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (i *IndividualRecording) Update(ctx context.Context, resourceId string, sid string, cname string, uid string,
	clientRequest *req.UpdateIndividualRecordingClientRequest,
) (*api.UpdateResp, error) {
	return i.updateAPI.Do(ctx, resourceId, sid, api.IndividualMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			StreamSubscribe: clientRequest.StreamSubscribe,
		},
	})
}

// @brief Stop individual cloud recording.
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
func (i *IndividualRecording) Stop(ctx context.Context, resourceId string, sid string, cname string, uid string, asyncStop bool) (*api.StopResp, error) {
	return i.stopAPI.Do(ctx, resourceId, sid, api.IndividualMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
}
