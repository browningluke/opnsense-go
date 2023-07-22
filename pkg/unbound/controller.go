package unbound

import (
	"github.com/browningluke/opnsense-go/pkg/api"
)

const unboundReconfigureEndpoint = "/unbound/service/reconfigure"

// Controller for unbound
type Controller struct {
	Api *api.Client
}

func (c *Controller) Client() *api.Client {
	return c.Api
}

func (c *Controller) Name() string {
	return "unbound"
}
