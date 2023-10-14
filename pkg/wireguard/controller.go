package wireguard

import "github.com/browningluke/opnsense-go/pkg/api"

const wireguardReconfigureEndpoint = "/wireguard/service/reconfigure"

// Controller for wireguard
type Controller struct {
	Api *api.Client
}

func (c *Controller) Client() *api.Client {
	return c.Api
}
