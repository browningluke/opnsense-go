package diagnostics

import "github.com/browningluke/opnsense-go/pkg/api"

// Controller for diagnostics
type Controller struct {
	Api *api.Client
}

func (c *Controller) Client() *api.Client {
	return c.Api
}
