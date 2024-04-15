package mixrecording

import (
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Impl struct {
	Base *baseV1.BaseCollection
}

func NewMixRecording() *Impl {
	return &Impl{}
}

var _ baseV1.MixRecording = (*Impl)(nil)

func (i *Impl) SetBase(base *baseV1.BaseCollection) {
	i.Base = base
}
func (i *Impl) Acquire() baseV1.AcquireMixRecording {
	return &Acquire{BaseAcquire: i.Base.Acquire()}
}

func (i *Impl) Query() baseV1.QueryMixRecording {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Start() baseV1.StartMixRecording {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Stop() baseV1.StopMixRecording {
	//TODO implement me
	panic("implement me")
}

func (i *Impl) Update() baseV1.UpdateMixRecording {
	//TODO implement me
	panic("implement me")
}
