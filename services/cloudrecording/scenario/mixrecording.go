package scenario

import (
	"context"

	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/api"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/req"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/resp"
)

type MixRecording struct {
	acquireAPI      *api.Acquire
	startAPI        *api.Start
	stopAPI         *api.Stop
	queryAPI        *api.Query
	updateLayoutAPI *api.UpdateLayout
	updateAPI       *api.Update
}

func NewMixRecording(
	acquireAPI *api.Acquire,
	startAPI *api.Start,
	stopAPI *api.Stop,
	queryAPI *api.Query,
	updateLayoutAPI *api.UpdateLayout,
	updateAPI *api.Update,
) *MixRecording {
	return &MixRecording{
		acquireAPI:      acquireAPI,
		startAPI:        startAPI,
		stopAPI:         stopAPI,
		queryAPI:        queryAPI,
		updateLayoutAPI: updateLayoutAPI,
		updateAPI:       updateAPI,
	}
}

// @brief Get a resource ID for mix cloud recording.
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
// @param clientRequest The request body.See req.AcquireMixRecodingClientRequest for details.
//
// @return Returns the response *AcquireResp. See api.AcquireResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (m *MixRecording) Acquire(ctx context.Context, cname string, uid string,
	clientRequest *req.AcquireMixRecodingClientRequest,
) (*api.AcquireResp, error) {
	var startParameter *api.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &api.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
		}
	}

	return m.acquireAPI.Do(ctx, &api.AcquireReqBody{
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

// @brief Start mix cloud recording.
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
// @param clientRequest The request body.See req.StartMixRecordingClientRequest for details.
//
// @return Returns the response *StartResp. See api.StartResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (m *MixRecording) Start(ctx context.Context, resourceId string, cname string, uid string,
	clientRequest *req.StartMixRecordingClientRequest,
) (*api.StartResp, error) {
	return m.startAPI.Do(ctx, resourceId, api.MixMode, &api.StartReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StartClientRequest{
			Token:               clientRequest.Token,
			RecordingFileConfig: clientRequest.RecordingFileConfig,
			RecordingConfig:     clientRequest.RecordingConfig,
			StorageConfig:       clientRequest.StorageConfig,
		},
	})
}

// @brief Query the status of mix cloud recording when the video file format is hls.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @return Returns the response *QueryMixRecordingHLSResp. See resp.QueryMixRecordingHLSResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (m *MixRecording) QueryHLS(ctx context.Context, resourceId string, sid string,
) (*resp.QueryMixRecordingHLSResp, error) {
	respData, err := m.queryAPI.Do(ctx, resourceId, sid, api.MixMode)
	if err != nil {
		return nil, err
	}

	var mixResp resp.QueryMixRecordingHLSResp

	mixResp.Response = respData.Response
	if respData.IsSuccess() {
		successResp := respData.SuccessResponse
		mixResp.SuccessResponse = resp.QueryMixRecordingHLSSuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSServerResponse(),
		}
	}

	return &mixResp, nil
}

// @brief Query the status of mix cloud recording when the video file format is hls and mp4.
//
// @since v0.8.0
//
// @param ctx Context to control the request lifecycle.
//
// @param resourceID The resource ID.
//
// @param sid The recording ID, identifying a recording cycle.
//
// @return Returns the response *QueryMixRecordingHLSResp. See resp.QueryMixRecordingHLSResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (m *MixRecording) QueryHLSAndMP4(ctx context.Context, resourceId string, sid string,
) (*resp.QueryMixRecordingHLSAndMP4Resp, error) {
	respData, err := m.queryAPI.Do(ctx, resourceId, sid, api.MixMode)
	if err != nil {
		return nil, err
	}

	var mixResp resp.QueryMixRecordingHLSAndMP4Resp

	mixResp.Response = respData.Response
	if respData.IsSuccess() {
		successResp := respData.SuccessResponse
		mixResp.SuccessResponse = resp.QueryMixRecordingHLSAndMP4SuccessResp{
			ResourceId:     successResp.ResourceId,
			Sid:            successResp.Sid,
			ServerResponse: successResp.GetMixRecordingHLSAndMP4ServerResponse(),
		}
	}

	return &mixResp, nil
}

// @brief Update the mix cloud recording configuration.
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
// @param clientRequest The request body. See req.UpdateMixRecordingClientRequest for details.
//
// @return Returns the response *UpdateResp. See api.UpdateResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (m *MixRecording) Update(ctx context.Context, resourceId string, sid string, cname string, uid string,
	clientRequest *req.UpdateMixRecordingClientRequest,
) (*api.UpdateResp, error) {
	return m.updateAPI.Do(ctx, resourceId, sid, api.MixMode, &api.UpdateReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateClientRequest{
			StreamSubscribe: clientRequest.StreamSubscribe,
		},
	})
}

// @brief Update the mix cloud recording layout.
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
// @param clientRequest The request body. See req.UpdateLayoutUpdateMixRecordingClientRequest for details.
//
// @return Returns the response *UpdateLayoutResp. See api.UpdateLayoutResp for details.
//
// @return Returns an error object. If the request fails, the error object is not nil and contains error information.
func (m *MixRecording) UpdateLayout(ctx context.Context, resourceId string, sid string, cname string, uid string,
	clientRequest *req.UpdateLayoutUpdateMixRecordingClientRequest,
) (*api.UpdateLayoutResp, error) {
	return m.updateLayoutAPI.Do(ctx, resourceId, sid, api.MixMode, &api.UpdateLayoutReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.UpdateLayoutClientRequest{
			MaxResolutionUID:           clientRequest.MaxResolutionUID,
			MixedVideoLayout:           clientRequest.MixedVideoLayout,
			BackgroundColor:            clientRequest.BackgroundColor,
			BackgroundImage:            clientRequest.BackgroundImage,
			DefaultUserBackgroundImage: clientRequest.DefaultUserBackgroundImage,
			LayoutConfig:               clientRequest.LayoutConfig,
			BackgroundConfig:           clientRequest.BackgroundConfig,
		},
	})
}

// @brief Stop mix cloud recording.
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
func (m *MixRecording) Stop(ctx context.Context, resourceId string, sid string, cname string, uid string,
	asyncStop bool,
) (*api.StopResp, error) {
	return m.stopAPI.Do(ctx, resourceId, sid, api.MixMode, &api.StopReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &api.StopClientRequest{
			AsyncStop: asyncStop,
		},
	})
}
