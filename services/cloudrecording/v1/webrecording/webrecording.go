package webrecording

import (
	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Impl struct {
	Base *baseV1.BaseCollection
}

func NewWebRecording() *Impl {
	return &Impl{}
}

var _ baseV1.WebRecording = (*Impl)(nil)

func (w *Impl) SetBase(base *baseV1.BaseCollection) {
	w.Base = base
}

func (w *Impl) Acquire() baseV1.AcquireWebRecording {
	return &Acquire{BaseAcquire: w.Base.Acquire()}
}

func (w *Impl) Query() baseV1.QueryWebRecording {
	return &Query{BaseQuery: w.Base.Query()}
}

func (w *Impl) Start() baseV1.StartWebRecording {
	return &Starter{BaseStarter: w.Base.Start()}
}

func (w *Impl) Stop() baseV1.StopWebRecording {
	return &Stop{BaseStop: w.Base.Stop()}
}

func (w *Impl) Update() baseV1.UpdateWebRecording {
	return &Update{BaseUpdate: w.Base.Update()}
}
