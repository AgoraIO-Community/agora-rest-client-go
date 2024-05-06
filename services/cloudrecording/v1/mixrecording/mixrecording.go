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
	return &Acquire{Base: i.Base.Acquire()}
}

func (i *Impl) Query() baseV1.QueryMixRecording {
	return &Query{Base: i.Base.Query()}
}

func (i *Impl) Start() baseV1.StartMixRecording {
	return &Start{Base: i.Base.Start()}
}

func (i *Impl) Stop() baseV1.StopMixRecording {
	return &Stop{Base: i.Base.Stop()}
}

func (i *Impl) Update() baseV1.UpdateMixRecording {
	return &Update{Base: i.Base.Update()}
}

func (i *Impl) UpdateLayout() baseV1.UpdateLayoutMixRecording {
	return &UpdateLayout{Base: i.Base.UpdateLayout()}
}
