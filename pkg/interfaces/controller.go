package interfaces

import (
	"github.com/browningluke/opnsense-go/pkg/api"
)

// Controller for interfaces
type Controller struct {
	Api *api.Client
}

func (c *Controller) Client() *api.Client {
	return c.Api
}
