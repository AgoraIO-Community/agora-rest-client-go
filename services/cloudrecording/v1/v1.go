package v1

import (
	"github.com/AgoraIO/agora-rest-client-go/core"
)

type Collection struct {
	prefixPath string
	client     core.Client
}

func NewCollection(prefixPath string, client core.Client) *Collection {
	return &Collection{
		prefixPath: "/v1" + prefixPath,
		client:     client,
	}
}

func (c *Collection) Acquire() *Acquire {
	return &Acquire{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *Collection) Start() *Starter {
	return &Starter{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *Collection) Stop() *Stop {
	return &Stop{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *Collection) Query() *Query {
	return &Query{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *Collection) Update() *Update {
	return &Update{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}

func (c *Collection) UpdateLayout() *UpdateLayout {
	return &UpdateLayout{
		client:     c.client,
		prefixPath: c.prefixPath,
	}
}
