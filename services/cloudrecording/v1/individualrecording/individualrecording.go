package individualrecording

import (
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Impl struct {
	Base *baseV1.BaseCollection
}

var _ baseV1.IndividualRecording = (*Impl)(nil)

func NewIndividualRecording() *Impl {
	return &Impl{}
}

func (i *Impl) SetBase(base *baseV1.BaseCollection) {
	i.Base = base
}

func (i *Impl) Acquire() baseV1.AcquireIndividualRecording {
	return &Acquire{Base: i.Base.Acquire()}
}

func (i *Impl) Start() baseV1.StartIndividualRecording {
	return &Start{Base: i.Base.Start()}
}

func (i *Impl) Query() baseV1.QueryIndividualRecording {
	return &Query{Base: i.Base.Query()}
}

func (i *Impl) Update() baseV1.UpdateIndividualRecording {
	return &Update{Base: i.Base.Update()}
}

func (i *Impl) Stop() baseV1.StopIndividualRecording {
	return &Stop{BaseStop: i.Base.Stop()}
}
