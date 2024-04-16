package v1

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type BaseCollection struct {
	prefixPath          string
	client              core.Client
	webRecording        WebRecording
	mixRecording        MixRecording
	individualRecording IndividualRecording
}

func NewCollection(prefixPath string, client core.Client, webRecording WebRecording, mixRecording MixRecording, individualRecording IndividualRecording) *BaseCollection {
	b := &BaseCollection{
		prefixPath:          "/v1" + prefixPath,
		client:              client,
		webRecording:        webRecording,
		mixRecording:        mixRecording,
		individualRecording: individualRecording,
	}
	b.webRecording.SetBase(b)
	b.mixRecording.SetBase(b)
	b.individualRecording.SetBase(b)

	return b
}

func (c *BaseCollection) Acquire() *Acquire {
	return &Acquire{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *BaseCollection) Start() *Starter {
	return &Starter{
		module:     "cloudRecording:start",
		logger:     c.client.GetLogger(),
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *BaseCollection) Stop() *Stop {
	return &Stop{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *BaseCollection) Query() *Query {
	return &Query{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *BaseCollection) Update() *Update {
	return &Update{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *BaseCollection) UpdateLayout() *UpdateLayout {
	return &UpdateLayout{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *BaseCollection) WebRecording() WebRecording {
	return c.webRecording
}

func (c *BaseCollection) MixRecording() MixRecording {
	return c.mixRecording
}

func (c *BaseCollection) IndividualRecording() IndividualRecording {
	return c.individualRecording
}
