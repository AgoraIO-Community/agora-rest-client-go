package v1

import "github.com/AgoraIO-Community/agora-rest-client-go/core"

type BaseCollection struct {
	prefixPath string
	client     core.Client
}

func NewCollection(prefixPath string, client core.Client) *BaseCollection {
	b := &BaseCollection{
		prefixPath: "/v1" + prefixPath,
		client:     client,
	}

	return b
}

func (b *BaseCollection) Acquire() *Acquire {
	return &Acquire{
		module:     "cloudTranscoder:acquire",
		logger:     b.client.GetLogger(),
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Create() *Create {
	return &Create{
		module:     "cloudTranscoder:create",
		logger:     b.client.GetLogger(),
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Query() *Query {
	return &Query{
		module:     "cloudTranscoder:query",
		logger:     b.client.GetLogger(),
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Delete() *Delete {
	return &Delete{
		module:     "cloudTranscoder:delete",
		logger:     b.client.GetLogger(),
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Update() *Update {
	return &Update{
		module:     "cloudTranscoder:update",
		logger:     b.client.GetLogger(),
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}
