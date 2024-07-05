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
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Create() *Create {
	return &Create{
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Query() *Query {
	return &Query{
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Delete() *Delete {
	return &Delete{
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}

func (b *BaseCollection) Update() *Update {
	return &Update{
		client:     b.client,
		prefixPath: b.prefixPath,
	}
}
