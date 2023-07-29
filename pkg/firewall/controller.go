package firewall

import "github.com/browningluke/opnsense-go/pkg/api"

// Controller for firewall
type Controller struct {
	Api *api.Client
}

func (c *Controller) Client() *api.Client {
	return c.Api
}
